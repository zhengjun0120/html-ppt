package chatimg

// 共享附图能力的回归：窗口变换（滑动窗口语义）、回放提取、脱敏。

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
)

func testImgURL(tag string) string {
	// 1x1 PNG，内容无所谓——窗口变换只看形状不看像素
	return "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg==" + tag
}

// userWithImages 构造一条带 N 张图（可夹文本）的 user 消息。
func userWithImages(text string, n int) openai.ChatCompletionMessageParamUnion {
	imgs := make([]Image, n)
	for i := range imgs {
		imgs[i] = Image{DataURL: testImgURL(string(rune('a' + i)))}
	}
	return BuildUserMessage(text, imgs)
}

func countImages(t *testing.T, msgs []openai.ChatCompletionMessageParamUnion) int {
	t.Helper()
	raw, err := json.Marshal(msgs)
	if err != nil {
		t.Fatal(err)
	}
	// 每张图恰好一个 data URL（"image_url" 字符串会同时出现在 type 值与字段名里，不能数它）
	return strings.Count(string(raw), `"data:image/`)
}

func TestWindowOldestImages(t *testing.T) {
	t.Run("无图快速路径", func(t *testing.T) {
		msgs := []openai.ChatCompletionMessageParamUnion{openai.SystemMessage("s"), openai.UserMessage("你好")}
		got, changed, err := WindowOldestImages(msgs, KeepRecentImages)
		if err != nil || changed {
			t.Fatalf("无图应零改动: changed=%v err=%v", changed, err)
		}
		if &got[0] != &msgs[0] && len(got) != len(msgs) {
			t.Fatal("快速路径应返回原切片")
		}
	})

	t.Run("未超窗零改动", func(t *testing.T) {
		msgs := []openai.ChatCompletionMessageParamUnion{userWithImages("图一", 2), userWithImages("图二", 2)}
		_, changed, err := WindowOldestImages(msgs, KeepRecentImages)
		if err != nil || changed {
			t.Fatalf("4 张恰好在窗内: changed=%v err=%v", changed, err)
		}
	})

	t.Run("超窗裁最老的", func(t *testing.T) {
		msgs := []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("s"),
			userWithImages("早轮文字", 3),           // 3 张（最老）
			openai.UserMessage("中间一轮纯文本"),
			userWithImages("近轮文字", 3),           // 3 张（最新）
		}
		got, changed, err := WindowOldestImages(msgs, KeepRecentImages)
		if err != nil || !changed {
			t.Fatalf("6 张超窗应有变换: changed=%v err=%v", changed, err)
		}
		if got := countImages(t, got); got != KeepRecentImages {
			t.Fatalf("窗口后应剩 %d 张图, got %d", KeepRecentImages, got)
		}
		// 最老消息：前 2 张被换成占位文本，原文字与第 3 张图保留
		raw, _ := json.Marshal(got[1])
		var m struct {
			Content []wirePart `json:"content"`
		}
		if err := json.Unmarshal(raw, &m); err != nil {
			t.Fatal(err)
		}
		if len(m.Content) != 4 {
			t.Fatalf("被改消息应有 4 个 part（文字+2占位+1图）: %d", len(m.Content))
		}
		if m.Content[0].Type != "text" || m.Content[0].Text != "早轮文字" {
			t.Fatalf("原文字应保留: %+v", m.Content[0])
		}
		if m.Content[1].Type != "text" || !strings.Contains(m.Content[1].Text, "超出可见窗口") {
			t.Fatalf("第 1 张图应降级为占位: %+v", m.Content[1])
		}
		if m.Content[2].Type != "text" || !strings.Contains(m.Content[2].Text, "超出可见窗口") {
			t.Fatalf("第 2 张图应降级为占位: %+v", m.Content[2])
		}
		if m.Content[3].Type != "image_url" || !strings.HasSuffix(m.Content[3].ImageURL.URL, "c") {
			t.Fatalf("窗口内最新的图应保留: %+v", m.Content[3])
		}
		// 最新消息 3 张图原样保留
		if got := countImages(t, []openai.ChatCompletionMessageParamUnion{got[3]}); got != 3 {
			t.Fatalf("最新消息的图不应被动: %d", got)
		}
		// 纯文本消息原样
		if raw, _ := json.Marshal(got[2]); strings.Contains(string(raw), "超出可见窗口") {
			t.Fatal("纯文本消息不应被改写")
		}
	})

	t.Run("变换结果能被 SDK union 读回", func(t *testing.T) {
		msgs := []openai.ChatCompletionMessageParamUnion{userWithImages("x", 5)}
		got, changed, err := WindowOldestImages(msgs, 2)
		if err != nil || !changed {
			t.Fatalf("changed=%v err=%v", changed, err)
		}
		// projectTranscript 对变换后的消息也能走（wire 往返不炸）；
		// 占位文本随 parts 一起被提取是预期行为（回放读的是库里未变换的消息）
		raw, _ := json.Marshal(got[0].OfUser)
		text, imgs := UserTextAndImages(raw)
		if !strings.HasPrefix(text, "x") || len(imgs) != 2 {
			t.Fatalf("往返后 text=%q imgs=%d", text, len(imgs))
		}
	})
}

func TestUserTextAndImages(t *testing.T) {
	t.Run("纯文本", func(t *testing.T) {
		raw, _ := json.Marshal(openai.UserMessage("只说话"))
		text, imgs := UserTextAndImages(raw)
		if text != "只说话" || imgs != nil {
			t.Fatalf("text=%q imgs=%v", text, imgs)
		}
	})
	t.Run("文本加图", func(t *testing.T) {
		raw, _ := json.Marshal(userWithImages("看这两张", 2))
		text, imgs := UserTextAndImages(raw)
		if text != "看这两张" || len(imgs) != 2 {
			t.Fatalf("text=%q imgs=%d", text, len(imgs))
		}
		if !strings.HasPrefix(imgs[0], "data:image/png;base64,") {
			t.Fatalf("图片应是 data URL: %q", imgs[0][:30])
		}
	})
	t.Run("纯图无文本", func(t *testing.T) {
		raw, _ := json.Marshal(userWithImages("", 1))
		text, imgs := UserTextAndImages(raw)
		if text != "" || len(imgs) != 1 {
			t.Fatalf("text=%q imgs=%d", text, len(imgs))
		}
	})
}

func TestRedactDataURLs(t *testing.T) {
	msgs := []openai.ChatCompletionMessageParamUnion{userWithImages("含图", 1)}
	raw, _ := json.Marshal(msgs)
	out := RedactDataURLs(raw)
	if strings.Contains(string(out), "iVBOR") {
		t.Fatal("base64 载荷未脱敏")
	}
	if !strings.Contains(string(out), "已脱敏") {
		t.Fatal("应有脱敏占位")
	}
	var back []map[string]any
	if err := json.Unmarshal(out, &back); err != nil {
		t.Fatalf("脱敏后应是合法 JSON: %v", err)
	}
}
