package usertpl

// 用户附图（customize_images.go）的回归：
//   - 校验表：形态/魔数/大小/尺寸/条数，错误粒度到「第 N 张」；
//   - 消息形状：带图请求的 user content 是数组（text + image_url），无图请求保持纯文本；
//   - 只看一轮：本轮工具循环内每轮都带图，下一轮 Customize 的请求里图降级为文字占位；
//   - 观测：run_start 落 img/u001.png，llm_request 的 JSONL 里 base64 已脱敏。

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"

	"html-ppt/backend/internal/chatimg"
)

func testPNG(t *testing.T, w, h int) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := png.Encode(&buf, image.NewRGBA(image.Rect(0, 0, w, h))); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func pngDataURL(t *testing.T, w, h int) string {
	t.Helper()
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(testPNG(t, w, h))
}

// newCapturingFakeLLM 与 newFakeLLM 同款脚本回放，额外把每次请求体存下来供断言。
func newCapturingFakeLLM(t *testing.T, rounds [][]string) (LLM, *[]string) {
	t.Helper()
	var calls atomic.Int32
	var mu sync.Mutex
	var bodies []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		mu.Lock()
		bodies = append(bodies, string(raw))
		mu.Unlock()
		n := int(calls.Add(1)) - 1
		if n >= len(rounds) {
			http.Error(w, "no script", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		for _, frame := range rounds[n] {
			_, _ = w.Write([]byte(frame))
		}
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
	}))
	t.Cleanup(srv.Close)
	client := openai.NewClient(option.WithAPIKey("test-key"), option.WithBaseURL(srv.URL))
	return LLM{Client: &client, Model: "fake-model"}, &bodies
}

func TestNormalizeUserImages(t *testing.T) {
	t.Run("空输入放行", func(t *testing.T) {
		imgs, err := chatimg.Normalize(nil)
		if err != nil || len(imgs) != 0 {
			t.Fatalf("nil 应放行: %v %v", imgs, err)
		}
	})

	t.Run("合法png且mime按魔数归一", func(t *testing.T) {
		// 客户端谎报 jpeg，服务端按魔数改回 png
		url := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(testPNG(t, 2, 3))
		imgs, err := chatimg.Normalize([]string{url})
		if err != nil {
			t.Fatalf("合法 png 被拒: %v", err)
		}
		if !strings.HasPrefix(imgs[0].DataURL, "data:image/png;base64,") {
			t.Fatalf("mime 未按魔数归一: %q", imgs[0].DataURL[:40])
		}
	})

	t.Run("jpeg放行", func(t *testing.T) {
		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, image.NewRGBA(image.Rect(0, 0, 6, 6)), nil); err != nil {
			t.Fatal(err)
		}
		imgs, err := chatimg.Normalize([]string{"data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())})
		if err != nil || !strings.HasPrefix(imgs[0].DataURL, "data:image/jpeg;base64,") {
			t.Fatalf("合法 jpeg 应放行: %v", err)
		}
	})

	cases := []struct {
		name string
		in   string
		want string // 报错里应含的片段
	}{
		{"不是data URL", "https://example.com/a.png", "data:image"},
		{"坏base64", "data:image/png;base64,@@@@@@@@", "base64"},
		{"不是图片", "data:image/png;base64," + base64.StdEncoding.EncodeToString([]byte("just some text, not an image at all")), "PNG/JPEG/WebP"},
		{"尺寸超界", "data:image/png;base64," + base64.StdEncoding.EncodeToString(testPNG(t, 9000, 1)), "尺寸"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := chatimg.Normalize([]string{c.in})
			if err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("want err 含 %q, got %v", c.want, err)
			}
		})
	}

	t.Run("超过单张大小", func(t *testing.T) {
		big := base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{0}, chatimg.MaxBytes+1))
		_, err := chatimg.Normalize([]string{"data:image/png;base64," + big})
		if err == nil || !strings.Contains(err.Error(), "MB") {
			t.Fatalf("超大应拒: %v", err)
		}
	})

	t.Run("超过条数", func(t *testing.T) {
		four := make([]string, chatimg.MaxPerMessage+1)
		for i := range four {
			four[i] = pngDataURL(t, 1, 1)
		}
		_, err := chatimg.Normalize(four)
		if err == nil || !strings.Contains(err.Error(), "最多") {
			t.Fatalf("超条数应拒: %v", err)
		}
	})
}

func TestCustomizeImageMessageShape(t *testing.T) {
	s, _, id := newTraceService(t, "img-shape")
	llm, bodies := newCapturingFakeLLM(t, [][]string{
		toolRound("write_tokens", `{"tokens":{"--accent":"#123456"}}`),
		toolRound("finish", `{"reply":"参考图的橙色已经落到主色上。"}`),
		toolRound("finish", `{"reply":"收到，继续调。"}`),
	})

	img := pngDataURL(t, 4, 4)
	if _, err := s.Customize(t.Context(), 9, id, "参考这张图的配色换主色", []string{img}, llm, nil); err != nil {
		t.Fatalf("第一轮: %v", err)
	}
	if len(*bodies) != 2 {
		t.Fatalf("第一轮应发 2 次请求, got %d", len(*bodies))
	}
	// 本轮工具循环内：每次请求的 user 消息都带图
	for i := 0; i < 2; i++ {
		b := (*bodies)[i]
		if !strings.Contains(b, `"type":"image_url"`) || !strings.Contains(b, `"data:image/png;base64,`) {
			t.Fatalf("第 %d 次请求应含图片部分（本轮内可见）: %s", i+1, clip(b, 300))
		}
		if !strings.Contains(b, `参考这张图的配色换主色`) {
			t.Fatalf("第 %d 次请求应含文本部分", i+1)
		}
	}

	// 第二轮：图降级为占位，base64 不再重发
	if _, err := s.Customize(t.Context(), 9, id, "再深一点", nil, llm, nil); err != nil {
		t.Fatalf("第二轮: %v", err)
	}
	third := (*bodies)[2]
	if strings.Contains(third, ";base64,") {
		t.Fatalf("下一轮请求不应再带图: %s", clip(third, 300))
	}
	if !strings.Contains(third, "发过 1 张图") {
		t.Fatalf("下一轮应含降级占位: %s", clip(third, 300))
	}
}

func TestCustomizePureImageMessage(t *testing.T) {
	s, _, id := newTraceService(t, "img-pure")
	llm, bodies := newCapturingFakeLLM(t, [][]string{
		toolRound("finish", `{"reply":"看到了你的图。"}`),
	})
	if _, err := s.Customize(t.Context(), 9, id, "", []string{pngDataURL(t, 1, 1)}, llm, nil); err != nil {
		t.Fatalf("纯图消息应放行: %v", err)
	}
	if !strings.Contains((*bodies)[0], `"type":"image_url"`) {
		t.Fatalf("纯图消息应只含图片部分: %s", clip((*bodies)[0], 300))
	}
}

func TestCustomizeImageValidationSurfaces(t *testing.T) {
	s, _, id := newTraceService(t, "img-bad")
	llm, _ := newCapturingFakeLLM(t, [][]string{toolRound("finish", `{"reply":"x"}`)})
	_, err := s.Customize(t.Context(), 9, id, "看图", []string{"data:image/png;base64,!!!!"}, llm, nil)
	if err == nil || !strings.Contains(err.Error(), "第 1 张图") {
		t.Fatalf("报错应定位到第几张: %v", err)
	}
}

func TestCustomizeTraceUserImages(t *testing.T) {
	s, _, id := newTraceService(t, "img-trace")
	s.traceCfg.CaptureImages = true
	llm, _ := newCapturingFakeLLM(t, [][]string{
		toolRound("finish", `{"reply":"图收到。"}`),
	})
	png := testPNG(t, 5, 5)
	if _, err := s.Customize(t.Context(), 9, id, "看这张图", []string{"data:image/png;base64," + base64.StdEncoding.EncodeToString(png)}, llm, nil); err != nil {
		t.Fatalf("customize: %v", err)
	}

	// trace 的落点：<sess>/<run_id>.jsonl 与 <sess>/<run_id>/img/
	var imgPath, jsonl string
	err := filepath.Walk(s.traceCfg.Dir, func(p string, info os.FileInfo, werr error) error {
		if werr != nil || info.IsDir() {
			return werr
		}
		if filepath.Base(p) == "u001.png" {
			imgPath = p
		}
		if strings.HasSuffix(p, ".jsonl") {
			jsonl = p
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if imgPath == "" {
		t.Fatal("用户附图未落盘（应为 <run>/img/u001.png）")
	}
	got, err := os.ReadFile(imgPath)
	if err != nil || !bytes.Equal(got, png) {
		t.Fatalf("落盘字节与原图不一致: %v", err)
	}

	raw, err := os.ReadFile(jsonl)
	if err != nil {
		t.Fatal(err)
	}
	// 占位文案本身含 ";base64,"，负向断言盯的是长 base64 载荷而非这个前缀
	var payloadRe = regexp.MustCompile(`;base64,[A-Za-z0-9+/=]{64,}`)
	lines := strings.Split(string(raw), "\n")
	var sawStart, sawLLMReq bool
	for _, ln := range lines {
		if payloadRe.MatchString(ln) {
			t.Fatalf("JSONL 里出现未脱敏的图片数据: %s", clip(ln, 200))
		}
		if strings.Contains(ln, `"kind":"run_start"`) {
			sawStart = true
			if !strings.Contains(ln, `"name":"u001.png"`) || !strings.Contains(ln, "用户附图 1") {
				t.Fatalf("run_start 应带附图元数据: %s", clip(ln, 300))
			}
			if !strings.Contains(ln, `img/u001.png`) {
				t.Fatalf("run_start 的附图应有可读 URL: %s", clip(ln, 300))
			}
		}
		if strings.Contains(ln, `"kind":"llm_request"`) {
			sawLLMReq = true
			if !strings.Contains(ln, "已脱敏") {
				t.Fatalf("llm_request 应出现脱敏占位: %s", clip(ln, 300))
			}
		}
	}
	if !sawStart || !sawLLMReq {
		t.Fatalf("缺关键事件: run_start=%v llm_request=%v", sawStart, sawLLMReq)
	}
}

func clip(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n]) + "…"
	}
	return s
}
