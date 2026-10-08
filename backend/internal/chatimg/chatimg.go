// Package chatimg 用户聊天附图的共享能力：校验归一、多模态消息构造、
// 观测脱敏、模型上下文的滑动窗口。
//
// 两方共用：usertpl（定制模板对话，「只看一轮」策略）与 agent（文稿对话，
// 图随会话存库、模型侧滑动窗口）。SDK 变体的读取 API 不稳定，回放与窗口
// 变换一律走 wire JSON 往返——marshal 形状即 API 协议形状，最稳。
package chatimg

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"regexp"
	"strings"

	"github.com/openai/openai-go/v3"
	_ "golang.org/x/image/webp"
)

const (
	// MaxPerMessage 单条消息附图上限（前端 MAX_CHAT_IMAGES 与此对齐）
	MaxPerMessage = 3
	// MaxBytes 单张解码后上限（前端压到长边 1568 的 JPEG，正常远够）
	MaxBytes = 4 << 20
	// MaxEdgePX 尺寸上限：防解码炸弹，也防极端长图的 token 失控
	MaxEdgePX = 8192
	// KeepRecentImages 文稿对话模型可见的滑动窗口：历史里最近 N 张图可见，
	// 更早的在请求时降级为文字占位（数据库与回放不受影响）
	KeepRecentImages = 4
)

// Image 一张校验通过的用户附图。
type Image struct {
	DataURL string // data:image/<png|jpeg|webp>;base64,...（mime 按魔数归一，不信任客户端自报）
	Bytes   []byte // 解码后的原始字节（观测落盘用）
}

// Normalize 校验并归一用户附图。在服务层调用（不是只靠 handler）：
// 同步端点、SSE 端点、以及将来任何调用方共用同一道闸。
func Normalize(images []string) ([]Image, error) {
	if len(images) > MaxPerMessage {
		return nil, fmt.Errorf("一次最多发 %d 张图（收到 %d 张）", MaxPerMessage, len(images))
	}
	out := make([]Image, 0, len(images))
	for i, raw := range images {
		img, err := normalizeOne(raw)
		if err != nil {
			return nil, fmt.Errorf("第 %d 张图：%w", i+1, err)
		}
		out = append(out, img)
	}
	return out, nil
}

func normalizeOne(raw string) (Image, error) {
	const prefix = "data:image/"
	if !strings.HasPrefix(raw, prefix) || !strings.Contains(raw, ";base64,") {
		return Image{}, fmt.Errorf("不是 data:image/*;base64 形态的图片")
	}
	b64 := raw[strings.Index(raw, ";base64,")+len(";base64,"):]
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return Image{}, fmt.Errorf("base64 解码失败")
	}
	if len(data) == 0 {
		return Image{}, fmt.Errorf("图片内容为空")
	}
	if len(data) > MaxBytes {
		return Image{}, fmt.Errorf("超过 %dMB（发图前请压缩）", MaxBytes>>20)
	}
	kind := sniffKind(data)
	if kind == "" {
		return Image{}, fmt.Errorf("不是 PNG/JPEG/WebP 图片")
	}
	// DecodeConfig 只读文件头，不整图解码：拿到尺寸做炸弹防护就够了。
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return Image{}, fmt.Errorf("图片头解析失败")
	}
	if cfg.Width <= 0 || cfg.Height <= 0 || cfg.Width > MaxEdgePX || cfg.Height > MaxEdgePX {
		return Image{}, fmt.Errorf("尺寸 %dx%d 超界（最长边 ≤%d）", cfg.Width, cfg.Height, MaxEdgePX)
	}
	return Image{
		DataURL: "data:image/" + kind + ";base64," + base64.StdEncoding.EncodeToString(data),
		Bytes:   data,
	}, nil
}

// sniffKind 魔数嗅探真实格式，返回 png/jpeg/webp 或空串。
func sniffKind(d []byte) string {
	switch {
	case len(d) >= 8 && d[0] == 0x89 && d[1] == 'P' && d[2] == 'N' && d[3] == 'G':
		return "png"
	case len(d) >= 3 && d[0] == 0xFF && d[1] == 0xD8 && d[2] == 0xFF:
		return "jpeg"
	case len(d) >= 12 && string(d[0:4]) == "RIFF" && string(d[8:12]) == "WEBP":
		return "webp"
	}
	return ""
}

// BuildUserMessage 组装带图的用户消息：文本可空（纯图消息），图片按序追加。
// 无图时走纯文本形态——既有请求体形状（string content）不被扰动，测试断言也稳。
func BuildUserMessage(message string, imgs []Image) openai.ChatCompletionMessageParamUnion {
	if len(imgs) == 0 {
		return openai.UserMessage(message)
	}
	parts := make([]openai.ChatCompletionContentPartUnionParam, 0, len(imgs)+1)
	if strings.TrimSpace(message) != "" {
		parts = append(parts, openai.TextContentPart(message))
	}
	for _, im := range imgs {
		parts = append(parts, openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{
			URL:    im.DataURL,
			Detail: "auto",
		}))
	}
	return openai.UserMessage(parts)
}

// imageDataURLRe 从序列化后的消息 JSON 里抠出图片 data URL（观测脱敏用）。
// 必须匹配到引号边界：替换值本身也是合法 JSON 字符串，抠完文件仍是合法 JSON。
var imageDataURLRe = regexp.MustCompile(`"data:image/[a-zA-Z0-9.+-]+;base64,[A-Za-z0-9+/=]+"`)

// RedactDataURLs 把消息 JSON 里的图片 data URL 换成短占位——原图有几百 KB，
// 内联进观测 JSONL 只会把逐行拉取的文件撑爆。替换值仍是被引号包住的字符串，
// 输出保持合法 JSON。
func RedactDataURLs(raw []byte) []byte {
	return imageDataURLRe.ReplaceAll(raw, []byte(`"data:image/*;base64,[用户附图已脱敏，原图见 img/]"`))
}
