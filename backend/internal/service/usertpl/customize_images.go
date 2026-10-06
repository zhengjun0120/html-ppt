package usertpl

// 定制对话的用户附图：把参考图带进模型上下文（多模态消息）。
//
// 生命周期刻意收窄成「只看一轮」：图随本轮 customizeLoop 的每次模型调用可见，
// 本轮成功结束后把该条用户消息降级回纯文本 + 占位说明。聊天补全协议下历史
// 每轮整体重发，一张压缩图 1-2k token，长对话费用会滚涨；模型当轮已经把图
// 消化进回复与文件改动，占位足够后续引用。本轮失败时不降级——图留着，重试
// 那轮模型还能看到（没看过的图说"看过"是撒谎）。

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"strings"

	"github.com/openai/openai-go/v3"
	_ "golang.org/x/image/webp"
)

const (
	maxUserImages  = 3
	maxImageBytes  = 4 << 20 // 单张解码后上限（前端会压到长边 1568 的 JPEG，正常远够）
	maxImageEdgePX = 8192    // 最长边上限：防解码炸弹，也防极端长图的 token 失控
)

// userImage 一张校验通过的用户附图。
type userImage struct {
	DataURL string // data:image/<png|jpeg|webp>;base64,...（mime 按魔数归一，不信任客户端自报）
	Bytes   []byte // 解码后的原始字节（观测落盘用）
}

// normalizeUserImages 校验并归一用户附图。校验在服务层做（不是只靠 handler）：
// 同步端点、SSE 端点、以及将来任何调用方共用同一道闸。
func normalizeUserImages(images []string) ([]userImage, error) {
	if len(images) > maxUserImages {
		return nil, fmt.Errorf("一次最多发 %d 张图（收到 %d 张）", maxUserImages, len(images))
	}
	out := make([]userImage, 0, len(images))
	for i, raw := range images {
		img, err := normalizeOneImage(raw)
		if err != nil {
			return nil, fmt.Errorf("第 %d 张图：%w", i+1, err)
		}
		out = append(out, img)
	}
	return out, nil
}

func normalizeOneImage(raw string) (userImage, error) {
	const prefix = "data:image/"
	if !strings.HasPrefix(raw, prefix) || !strings.Contains(raw, ";base64,") {
		return userImage{}, fmt.Errorf("不是 data:image/*;base64 形态的图片")
	}
	b64 := raw[strings.Index(raw, ";base64,")+len(";base64,"):]
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return userImage{}, fmt.Errorf("base64 解码失败")
	}
	if len(data) == 0 {
		return userImage{}, fmt.Errorf("图片内容为空")
	}
	if len(data) > maxImageBytes {
		return userImage{}, fmt.Errorf("超过 %dMB（发图前请压缩）", maxImageBytes>>20)
	}
	kind := sniffImageKind(data)
	if kind == "" {
		return userImage{}, fmt.Errorf("不是 PNG/JPEG/WebP 图片")
	}
	// DecodeConfig 只读文件头，不整图解码：拿到尺寸做炸弹防护就够了。
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return userImage{}, fmt.Errorf("图片头解析失败")
	}
	if cfg.Width <= 0 || cfg.Height <= 0 || cfg.Width > maxImageEdgePX || cfg.Height > maxImageEdgePX {
		return userImage{}, fmt.Errorf("尺寸 %dx%d 超界（最长边 ≤%d）", cfg.Width, cfg.Height, maxImageEdgePX)
	}
	return userImage{
		DataURL: "data:image/" + kind + ";base64," + base64.StdEncoding.EncodeToString(data),
		Bytes:   data,
	}, nil
}

// sniffImageKind 魔数嗅探真实格式，返回 png/jpeg/webp 或空串。
func sniffImageKind(d []byte) string {
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

// buildUserMessage 组装带图的用户消息：文本可空（纯图消息），图片按序追加。
// 无图时走纯文本形态——既有请求体形状（string content）不被扰动，测试断言也稳。
func buildUserMessage(message string, imgs []userImage) openai.ChatCompletionMessageParamUnion {
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

// stubImageNote 历史降级的占位文案（附在原文本后面）。
func stubImageNote(n int) string {
	if n == 1 {
		return "\n\n（用户随这条消息发过 1 张图，已在当时看过，此处不再重发）"
	}
	return fmt.Sprintf("\n\n（用户随这条消息发过 %d 张图，已在当时看过，此处不再重发）", n)
}

// degradeUserMessage 把带图的历史用户消息重建为纯文本 + 占位。idx 是消息在
// sess.messages 里的下标，n 是当初附图张数；n=0 或下标越界都是无害空操作。
func degradeUserMessage(sess *customizeSession, idx, n int, originalText string) {
	if n == 0 || idx < 0 || idx >= len(sess.messages) {
		return
	}
	text := strings.TrimSpace(originalText)
	if text == "" {
		text = "（用户发了一张图）"
	}
	sess.messages[idx] = openai.UserMessage(text + stubImageNote(n))
}
