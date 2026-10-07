package chatimg

// 滑动窗口与回放提取：都走 wire JSON 往返。
//
// openai-go 的 union 变体读取 API 在小版本间不稳定，而「marshal 出来的形状 =
// API 协议形状」是稳定的契约——对协议形状做解析/改写，再 unmarshal 回 union，
// 比追 SDK 内部类型省心且不会漏字段。整段消息数组里非 user 消息的
// json.RawMessage 原样透传，零失真。

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/openai/openai-go/v3"
)

// windowedMsg 解析消息数组时的最小形状；其余字段以 RawMessage 原样透传。
type windowedMsg struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content,omitempty"`
}

type wirePart struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	ImageURL *struct {
		URL string `json:"url"`
	} `json:"image_url,omitempty"`
}

// WindowOldestImages 滑动窗口：把 messages 里超出 keep 的最老图片 parts 替换成
// 文字占位，返回新消息数组。纯请求时变换——数据库里的完整历史、刷新回放的
// 缩略图都不受影响；下一个窗口里的图照常重发。
//
// changed=false 表示没有任何图片或未超窗，返回原切片（零拷贝快速路径）。
func WindowOldestImages(messages []openai.ChatCompletionMessageParamUnion, keep int) ([]openai.ChatCompletionMessageParamUnion, bool, error) {
	if len(messages) == 0 || keep <= 0 {
		return messages, false, nil
	}
	raw, err := json.Marshal(messages)
	if err != nil {
		return messages, false, fmt.Errorf("序列化消息失败: %w", err)
	}
	if !bytes.Contains(raw, []byte(`"image_url"`)) {
		return messages, false, nil // 快速路径：整段历史没有图
	}
	var msgs []windowedMsg
	if err := json.Unmarshal(raw, &msgs); err != nil {
		return messages, false, fmt.Errorf("解析消息失败: %w", err)
	}

	// 按时间序收集全部图片 part 的位置
	type at struct{ mi, pi int }
	var all []at
	for mi, m := range msgs {
		if m.Role != "user" {
			continue
		}
		var parts []wirePart
		if err := json.Unmarshal(m.Content, &parts); err != nil {
			continue // 纯文本 content 或解析失败——不是带图消息，跳过
		}
		for pi, p := range parts {
			if p.Type == "image_url" {
				all = append(all, at{mi, pi})
			}
		}
	}
	if len(all) <= keep {
		return messages, false, nil
	}

	stale := make(map[at]struct{}, len(all)-keep)
	for _, pos := range all[:len(all)-keep] {
		stale[pos] = struct{}{}
	}
	const staleStub = "（用户发的这张参考图已超出可见窗口，不再随请求重发；需要时让用户重新发送）"
	for mi, m := range msgs {
		if m.Role != "user" {
			continue
		}
		var parts []wirePart
		if err := json.Unmarshal(m.Content, &parts); err != nil {
			continue
		}
		touched := false
		for pi := range parts {
			if parts[pi].Type == "image_url" {
				if _, ok := stale[at{mi, pi}]; !ok {
					continue
				}
				parts[pi] = wirePart{Type: "text", Text: staleStub}
				touched = true
			}
		}
		if touched {
			fresh, err := json.Marshal(parts)
			if err != nil {
				return messages, false, fmt.Errorf("序列化窗口变换失败: %w", err)
			}
			msgs[mi].Content = fresh
		}
	}
	freshRaw, err := json.Marshal(msgs)
	if err != nil {
		return messages, false, fmt.Errorf("序列化窗口变换失败: %w", err)
	}
	var out []openai.ChatCompletionMessageParamUnion
	if err := json.Unmarshal(freshRaw, &out); err != nil {
		return messages, false, fmt.Errorf("回读窗口变换失败: %w", err)
	}
	return out, true, nil
}

// UserTextAndImages 从 user 消息的 wire JSON（marshal 过的 union 值）里提取
// 正文文本与图片 data URL 列表。会话回放用：parts 数组形态的消息从此不再
// 被当空串，刷新恢复能同时拿回文字与缩略图。
func UserTextAndImages(userMsgWire json.RawMessage) (text string, images []string) {
	var m struct {
		Content json.RawMessage `json:"content"`
	}
	if err := json.Unmarshal(userMsgWire, &m); err != nil {
		return "", nil
	}
	var s string
	if err := json.Unmarshal(m.Content, &s); err == nil {
		return s, nil // 纯文本形态
	}
	var parts []wirePart
	if err := json.Unmarshal(m.Content, &parts); err != nil {
		return "", nil
	}
	var b strings.Builder
	for _, p := range parts {
		switch p.Type {
		case "text":
			if p.Text != "" {
				if b.Len() > 0 {
					b.WriteByte('\n')
				}
				b.WriteString(p.Text)
			}
		case "image_url":
			if p.ImageURL != nil && p.ImageURL.URL != "" {
				images = append(images, p.ImageURL.URL)
			}
		}
	}
	return b.String(), images
}
