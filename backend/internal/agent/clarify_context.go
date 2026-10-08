package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"html-ppt/backend/internal/store"

	"github.com/openai/openai-go/v3"
)

// 模板推荐的诉求来源：把 deck 名下会话里的真实用户消息拼成一段文本，
// 喂给推荐 LLM 当「用户想要什么」的上下文。
//
// 拼接顺序刻意用会话创建时间正序：澄清会话（意图的出处）永远排最前，
// 后续大纲修订会话的碎语排在后面。maxChars 截断保头部——超限时丢的是
// 末尾的修订碎语，而不是开头的核心诉求。

// ClarifyUserMessages 收集 deck 名下全部会话中 role=user 的真实消息
//（滤掉后端合成的工作流指令，口径同回放投影），按会话先后拼成一段文本。
// maxChars <= 0 时用 2000。没有会话/没有用户消息返回空串 + nil（不是错误）。
func (as *AgentService) ClarifyUserMessages(ctx context.Context, userID uint, deckID string, maxChars int) (string, error) {
	if as.st == nil {
		return "", ErrSessionStoreUnavailable
	}
	if maxChars <= 0 {
		maxChars = 2000
	}
	// 这次必须 SELECT messages：诉求就在消息体里（DeckSessions 刻意不选它，
	// 那是列表页的口径；这里是内容消费方）
	var rows []store.ChatSession
	if err := as.st.DB.WithContext(ctx).
		Where("deck_id = ? AND user_id = ?", deckID, userID).
		Order("created_at ASC, id ASC").
		Find(&rows).Error; err != nil {
		return "", fmt.Errorf("查询 deck 会话: %w", err)
	}

	var b strings.Builder
	for i := range rows {
		var messages []openai.ChatCompletionMessageParamUnion
		if rows[i].Messages == "" {
			continue
		}
		if err := json.Unmarshal([]byte(rows[i].Messages), &messages); err != nil {
			continue // 坏形状会话不阻断推荐——诉求是尽力而为的上下文
		}
		for _, m := range messages {
			if m.OfUser == nil {
				continue
			}
			content := strings.TrimSpace(contentString(m.OfUser.Content.OfString))
			if content == "" || isSyntheticUserMsg(content) {
				continue
			}
			if b.Len() > 0 {
				b.WriteByte('\n')
			}
			b.WriteString("- ")
			b.WriteString(content)
			if b.Len() >= maxChars {
				goto done
			}
		}
	}
done:
	s := b.String()
	if len(s) > maxChars {
		// 按 rune 收尾，避免把多字节字符腰斩成乱码
		runes := []rune(s)
		if len(runes) > maxChars {
			s = string(runes[:maxChars])
		}
	}
	return s, nil
}
