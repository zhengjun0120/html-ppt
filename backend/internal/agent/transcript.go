package agent

import (
	"context"
	"errors"
	"fmt"
	"html-ppt/backend/internal/store"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/param"
)

// 对话历史的读取端：把 ChatSession.Messages（OpenAI 协议形状的 JSON 整体）
// 投影成前端能直接渲染的"可回放消息列表"。
//
// 为什么是投影而不是把 messages 原样吐出去：协议形状是给模型看的——system
// 混在最前面、tool 消息只带 tool_call_id 没有名字、assistant 可能只有 tool_calls
// 没有正文。前端要的却是"用户气泡 / 助手气泡 / 工具卡片（名字+参数+结果）"。
// 在这里投影一次，渲染层只认 TranscriptMessage 一个形状，协议演化不外溢。
//
// 一致性目标：同一次对话，回放出来的顺序和语义与直播（SSE 事件流）一致。
// 流式期间的临时态（think、tool 参数的增量）不落库，回放不含它们——
// 回放展示的是"协议里留了痕的东西"。

// ErrSessionStoreUnavailable 数据库不可用（降级模式），会话历史功能关闭。
// 独立成哨兵错误：handler 要把"功能关了"（503）和"查询失败"（500）分开讲，
// 混在一起会让前端在降级模式下不停重试一个永远不可能成功的请求。
var ErrSessionStoreUnavailable = errors.New("数据库不可用，会话功能关闭")

// SessionSummary 会话摘要。列表页只需要知道"有哪些对话、各自聊到哪了"，
// 完整消息列表（longtext，单个会话可能上百 KB）不随列表拉取。
type SessionSummary struct {
	ID        uint      `json:"id"`
	DeckID    string    `json:"deck_id"`
	Title     string    `json:"title"`
	Pending   bool      `json:"pending"` // true = 有 ask_user 提问还没回答
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TranscriptToolCall 一次工具调用的展示形状（来自 assistant 消息的声明）。
type TranscriptToolCall struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// TranscriptMessage 一条可回放消息。
//
// 工具调用挂在产生它的 assistant 消息上（ToolCalls）；工具结果单独一条
// role=tool 的消息，用 ToolCallID 关联。ToolName 是顺手解好的名字——
// 协议里 tool 消息只有 id，不解好的话每个前端都得自己往前翻 assistant 声明。
type TranscriptMessage struct {
	Seq        int64                `json:"seq"`
	Role       string               `json:"role"` // user | assistant | tool
	Content    string               `json:"content"`
	ToolCalls  []TranscriptToolCall `json:"tool_calls,omitempty"`
	ToolCallID string               `json:"tool_call_id,omitempty"`
	ToolName   string               `json:"tool_name,omitempty"`
}

// Transcript 一个会话的完整回放：会话元信息 + 投影后的消息列表。
type Transcript struct {
	Session  SessionSummary      `json:"session"`
	Messages []TranscriptMessage `json:"messages"`
}

// DeckSessions 列出某个 deck 名下的全部会话，按最近活跃倒序。
// 前端进工作台时取第一条当"当前对话"，其余留给"历史会话"切换。
func (as *AgentService) DeckSessions(ctx context.Context, userID uint, deckID string) ([]SessionSummary, error) {
	if as.st == nil {
		return nil, ErrSessionStoreUnavailable
	}
	var rows []store.ChatSession
	// 刻意不 SELECT messages：那是 longtext，列表页用不上，
	// 不加限制的话"看一眼有哪些对话"要把每个会话的完整上下文从库里搬出来
	err := as.st.DB.WithContext(ctx).
		Select("id", "user_id", "deck_id", "pending_ask", "title", "created_at", "updated_at").
		Where("deck_id = ? AND user_id = ?", deckID, userID).
		Order("updated_at DESC").
		Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("查询会话列表: %w", err)
	}
	out := make([]SessionSummary, 0, len(rows)) // 保证空列表序列化成 [] 而不是 null
	for i := range rows {
		out = append(out, sessionSummary(&rows[i]))
	}
	return out, nil
}

// RecentSessionSummary 最近会话摘要 + 所属 deck 的阶段/格式。
// DeckStage 为空 = 会话还没落到 deck（澄清中、大纲未产出）——这类会话
// 在文稿列表里不可见，只能靠这里被发现。
type RecentSessionSummary struct {
	SessionSummary
	DeckStage  string `json:"deck_stage"`
	DeckFormat string `json:"deck_format"`
}

// RecentSessions 用户跨 deck 的最近会话（不含消息体，按活跃时间倒序）。
// /new 的"继续上次对话"横幅靠它发现没走完向导的澄清/大纲会话：
// 按 deck 查会话的接口（DeckSessions）在没有 deck 上下文的页面上无从下手。
func (as *AgentService) RecentSessions(ctx context.Context, userID uint, limit int) ([]RecentSessionSummary, error) {
	if as.st == nil {
		return nil, ErrSessionStoreUnavailable
	}
	if limit <= 0 || limit > 50 {
		limit = 10
	}
	var rows []store.ChatSession
	err := as.st.DB.WithContext(ctx).
		Select("id", "user_id", "deck_id", "pending_ask", "title", "created_at", "updated_at").
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Limit(limit).
		Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("查询最近会话: %w", err)
	}
	// deck 阶段批量补齐（会话→deck 是多对一，去重后一次查询）
	deckIDs := make([]string, 0, len(rows))
	for i := range rows {
		if rows[i].DeckID != "" {
			deckIDs = append(deckIDs, rows[i].DeckID)
		}
	}
	deckByID := map[string]store.Deck{}
	if len(deckIDs) > 0 {
		var decks []store.Deck
		if err := as.st.DB.WithContext(ctx).
			Select("id", "format", "stage").
			Where("id IN ?", deckIDs).
			Find(&decks).Error; err != nil {
			return nil, fmt.Errorf("查询会话所属 deck: %w", err)
		}
		for _, d := range decks {
			deckByID[d.ID] = d
		}
	}
	out := make([]RecentSessionSummary, 0, len(rows))
	for i := range rows {
		s := RecentSessionSummary{SessionSummary: sessionSummary(&rows[i])}
		if d, ok := deckByID[rows[i].DeckID]; ok {
			s.DeckStage = d.Stage
			s.DeckFormat = d.Format
		}
		out = append(out, s)
	}
	return out, nil
}

func sessionSummary(sess *store.ChatSession) SessionSummary {
	return SessionSummary{
		ID:        sess.ID,
		DeckID:    sess.DeckID,
		Title:     sess.Title,
		Pending:   parsePendingAsk(sess).ToolCallID != "",
		CreatedAt: sess.CreatedAt,
		UpdatedAt: sess.UpdatedAt,
	}
}

// SessionTranscript 读一个会话的可回放消息列表。
// afterSeq > 0 时只返回 seq 更大的消息（增量拉取）。
func (as *AgentService) SessionTranscript(ctx context.Context, userID, sessionID uint, afterSeq int64) (*Transcript, error) {
	if as.st == nil {
		return nil, ErrSessionStoreUnavailable
	}
	// loadSession 顺带做了归属校验（别人的会话一律"不存在"）和坏形状修复：
	// 回放展示"修好之后"的形状，和下一次真正发给模型的东西保持一致。
	sess, messages, err := as.loadSession(ctx, userID, sessionID)
	if err != nil {
		return nil, err
	}
	pending := parsePendingAsk(sess)
	t := &Transcript{
		Session:  sessionSummary(sess),
		Messages: projectTranscript(messages, pending.ToolCallID, afterSeq),
	}
	if t.Messages == nil {
		t.Messages = []TranscriptMessage{}
	}
	return t, nil
}

// projectTranscript 把协议消息数组投影成回放列表（纯函数，DB 之外可测）。
//
// system 不进对话流——它是提示词不是对话内容。
//
// pendingToolCallID 是暂停中那次 ask_user 的 tool_call_id（空串 = 没有暂停）。
// loadSession 的修复流程会给它补一条兜底应答（那是发给模型用的），展示时必须
// 丢掉：用户还没回答，回放里不能假装他答了。其余修复补的"会话中断"说明保留——
// 如实展示"这次调用没有结果"。
//
// afterSeq 只影响返回哪些消息，不影响投影本身：工具名的解析要靠读全量
// （tool 消息的名字来自它前面的 assistant 声明），增量拉取也得先投影再裁剪。
//
// seq 是消息在数组里的下标。健康会话的消息只追加、不改动，seq 因此稳定，
// after_seq 的增量口径才成立。例外是坏形状会话被修复时可能删掉/挪动消息，
// 那种会话本来就是异常态，以修复后为准重新全量拉一次即可。
func projectTranscript(messages []openai.ChatCompletionMessageParamUnion, pendingToolCallID string, afterSeq int64) []TranscriptMessage {
	out := make([]TranscriptMessage, 0, len(messages))
	toolNames := map[string]string{} // tool_call_id → 工具名，从 assistant 的声明解出来
	for i, m := range messages {
		seq := int64(i)
		switch {
		case m.OfUser != nil:
			out = append(out, TranscriptMessage{Seq: seq, Role: "user", Content: contentString(m.OfUser.Content.OfString)})
		case m.OfAssistant != nil:
			tm := TranscriptMessage{Seq: seq, Role: "assistant", Content: contentString(m.OfAssistant.Content.OfString)}
			for _, tc := range m.OfAssistant.ToolCalls {
				id, fn := tc.GetID(), tc.GetFunction()
				if id == nil || fn == nil {
					continue // 半截的工具调用（理论上有修复流程兜底），投影不出可展示的卡片
				}
				if *id != "" {
					toolNames[*id] = fn.Name
				}
				tm.ToolCalls = append(tm.ToolCalls, TranscriptToolCall{ID: *id, Name: fn.Name, Arguments: fn.Arguments})
			}
			out = append(out, tm)
		case m.OfTool != nil:
			if pendingToolCallID != "" && m.OfTool.ToolCallID == pendingToolCallID {
				continue // 修复流程给"还没回答的提问"补的占位应答，不属于对话历史
			}
			out = append(out, TranscriptMessage{
				Seq:        seq,
				Role:       "tool",
				Content:    contentString(m.OfTool.Content.OfString),
				ToolCallID: m.OfTool.ToolCallID,
				ToolName:   toolNames[m.OfTool.ToolCallID],
			})
		default:
			// system 等其他角色跳过，但 seq 仍按数组下标计，保持增量口径与库里的位置一致
		}
	}
	// 裁剪放在投影**之后**：tool 消息的名字来自它前面的 assistant 声明，
	// 那条 assistant 往往在 after_seq 之前——先裁剪的话增量拉取会丢名字
	if afterSeq <= 0 {
		return out
	}
	trimmed := make([]TranscriptMessage, 0, len(out))
	for _, m := range out {
		if m.Seq > afterSeq {
			trimmed = append(trimmed, m)
		}
	}
	return trimmed
}

// contentString 取消息正文的字符串形态。本系统写路径只用纯文本构造消息
// （openai.UserMessage / ToolMessage / ToParam），OfString 必然有值；
// 内容分片数组（图片输入那种）不是会话历史的形状，真出现了就当空串，
// 不让一个边缘形状把整个回放接口打挂。
func contentString(opt param.Opt[string]) string {
	if param.IsOmitted(opt) {
		return ""
	}
	return opt.Value
}
