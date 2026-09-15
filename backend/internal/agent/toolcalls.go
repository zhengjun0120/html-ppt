package agent

import (
	"context"
	"encoding/json"
	"errors"
	"html-ppt/backend/internal/store"
	"log"

	"github.com/openai/openai-go/v3"
)

// 本文件守的是一条**协议不变式**：每条带 tool_calls 的 assistant 消息，后面必须紧跟齐
// 全部对应 tool_call_id 的 tool 消息，中间不能夹别的角色。
//
// 违反它的后果不是"这轮报错"，而是**这个会话废了**：API 直接返回
//   400 An assistant message with 'tool_calls' must be followed by tool messages
//       responding to each 'tool_call_id'
// 而这个错误不会自愈——坏形状的消息已经被持久化进 ChatSession.Messages 了，
// 之后每次请求都带着它，连"老老实实回答问题"这条路也一起 400。
//
// 历史上真实发生过：agent 提问后用户没回答、直接又发了一条消息，就正好构造出这个形状
//（见 NewStreamChat 里那道闸门的注释）。所以这里做两件事：
//   - 把要追加消息的入口全部拦住（guardNewMessage）
//   - 把已经被弄坏的会话在加载时修回来（repairToolCalls）

// ErrPendingAsk 会话还在等用户回答 ask_user 的提问，此时不接受新的用户消息。
//
// 用哨兵错误而不是普通 error：handler 要据此把"用户的待办"和"服务故障"分开讲，
// 前者说"请先回答提问"，后者才说"稍后重试"。混在一起会让用户一直重试一个
// 永远不可能成功的请求。
var ErrPendingAsk = errors.New("agent: waiting for the answer of a pending ask_user")

// pendingAsk 会话的暂停态。存在 ChatSession.PendingAsk 里的一小段 JSON。
type pendingAsk struct {
	ToolCallID string `json:"tool_call_id"` // 在等回答的那次 ask_user 调用
	RunID      string `json:"run_id"`       // 它属于哪个观测 run（老数据里没有）
}

// parsePendingAsk 解出暂停态。
//
// 刻意**不返回错误**：暂停态坏了不该让整个会话不可用——调用方拿到空 ToolCallID
// 就当作"没有暂停"，而修复流程本身更不该因为暂停态损坏而失败。
func parsePendingAsk(sess *store.ChatSession) pendingAsk {
	var p pendingAsk
	if sess == nil || sess.PendingAsk == "" {
		return p
	}
	if err := json.Unmarshal([]byte(sess.PendingAsk), &p); err != nil {
		log.Printf("[warn] 会话 %d 的暂停态解析失败（当作没有暂停）err: %v", sess.ID, err)
		return pendingAsk{}
	}
	return p
}

func encodePendingAsk(p pendingAsk) string {
	raw, err := json.Marshal(p)
	if err != nil {
		// 三个字符串字段的 struct 序列化不会失败；真失败了也不该让对话挂掉
		log.Printf("[warn] 序列化暂停态失败 err: %v", err)
		return ""
	}
	return string(raw)
}

// guardNewMessage 判断"现在能不能往这个会话里追加一条新的用户消息"。
//
// 拦住的原因见文件开头：暂停态的消息数组结尾是一条带 tool_calls 的 assistant 消息，
// 它的 tool 应答要等用户回答才补上。此时追加用户消息就构造出那个必然 400 的形状，
// 而且会被持久化、把这个会话永久弄坏。
//
// 为什么不是"当作提问作废、直接继续"：那需要替用户猜他到底是想答题还是想改需求。
// 猜错的代价是模型把"算了，先加一页"当成对问题的回答——比多一次点击更让人困惑。
func guardNewMessage(sess *store.ChatSession) error {
	if parsePendingAsk(sess).ToolCallID != "" {
		return ErrPendingAsk
	}
	return nil
}

// toolCallIDsOf 取一条 assistant 消息声明的 tool_call_ids（不是 assistant 或没调用工具就返回 nil）。
func toolCallIDsOf(m openai.ChatCompletionMessageParamUnion) []string {
	if m.OfAssistant == nil || len(m.OfAssistant.ToolCalls) == 0 {
		return nil
	}
	ids := make([]string, 0, len(m.OfAssistant.ToolCalls))
	for _, tc := range m.OfAssistant.ToolCalls {
		if id := tc.GetID(); id != nil && *id != "" {
			ids = append(ids, *id)
		}
	}
	return ids
}

// toolMsgCallID 取一条 tool 消息回应的 tool_call_id（不是 tool 消息就返回空串）。
func toolMsgCallID(m openai.ChatCompletionMessageParamUnion) string {
	if m.OfTool == nil {
		return ""
	}
	return m.OfTool.ToolCallID
}

// repairToolCalls 把消息数组修成满足协议不变式的形状，返回修了几处。
//
// pendingID 是会话暂停态记录的那个 tool_call_id。它单独给一句更准确的兜底说明：
// 其余缺失的 tool 应答是"会话中断"，而这一个的缺失是**预期之中的**（正在等用户回答）。
//
// 返回值 n 是**修了几处**，并且刻意保证 n > 0 当且仅当消息确实变了——
// 调用方写的是 `if n > 0 { messages = repaired }`，所以任何一种改动漏掉计数，
// 那种修复就会被整个丢掉（踩过一次：只数了"补齐"没数"归位"，于是被 400 卡死、
// 恰好需要归位的会话压根没被修，日志上什么都看不出来）。
//
// 四件事一起做：
//  1. 补齐：assistant 声明了 tool_call 却没有对应 tool 消息 → 紧跟其后补一条说明
//  2. 归位：tool 消息存在但位置不对（中间夹了别的角色）→ 移到它的 assistant 消息后面
//  3. 去孤：没有对应 tool_call 的 tool 消息丢掉
//     （API 同样拒绝"tool 消息不回应任何 tool_calls"，留着就是把会话卡死）
//  4. 去重：同一个 id 有两条应答只留先出现的
func repairToolCalls(messages []openai.ChatCompletionMessageParamUnion, pendingID string) ([]openai.ChatCompletionMessageParamUnion, int) {
	if len(messages) == 0 {
		return messages, 0
	}

	// 全部被声明过的 id
	declared := map[string]bool{}
	for _, m := range messages {
		for _, id := range toolCallIDsOf(m) {
			declared[id] = true
		}
	}

	// 把 tool 消息按 id 收上来，并标记原位置需要跳过（第 3 步统一在 assistant 后面重建）
	answers := map[string]openai.ChatCompletionMessageParamUnion{}
	answerAt := map[string]int{} // 原位置，用来判断"归位"算不算一次修复
	skip := make([]bool, len(messages))
	fixed := 0
	for i, m := range messages {
		id := toolMsgCallID(m)
		if id == "" {
			continue
		}
		skip[i] = true
		if !declared[id] {
			fixed++ // 孤儿：丢掉
			continue
		}
		if _, dup := answers[id]; dup {
			fixed++ // 同 id 的第二条应答：丢掉
			continue
		}
		answers[id] = m
		answerAt[id] = i
	}

	out := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages))
	expected := 0 // 其中有多少处是"预期之内"的（正在等用户回答的那个提问）
	for i, m := range messages {
		if skip[i] {
			continue
		}
		out = append(out, m)
		ids := toolCallIDsOf(m)
		if len(ids) == 0 {
			continue
		}
		// 紧跟这条 assistant 消息，按它自己声明的顺序把应答摆齐
		for k, id := range ids {
			if a, ok := answers[id]; ok {
				// 位置对不上（应答原先被夹在了别的角色之后）算一次归位。
				// 健康对话里第 k 条应答原本就在 i+1+k，所以这里不会误报
				if answerAt[id] != i+1+k {
					fixed++
				}
				out = append(out, a)
				continue
			}
			out = append(out, openai.ToolMessage(repairNote(id, pendingID), id))
			fixed++
			if id == pendingID {
				// 暂停态**本来**就有一条没应答的 ask_user 调用（正在等用户回答），
				// 所以这里的补齐是预期之内的：紧接着 setToolAnswer 就会把它换成
				// 用户的真实回答。不算异常，否则每一次正常的"提问 → 回答"都会刷一条
				// 修复警告，把真正需要看的那种淹掉。
				expected++
			}
		}
	}

	if fixed > expected {
		log.Printf("[warn] 会话消息修复：消息里的工具调用不成对，改了 %d 处（其中 %d 处是待答提问的兜底）", fixed, expected)
	}
	return out, fixed
}

// repairNote 给"缺失的工具结果"编一句如实的话。
// 内容要明确、不能像个真结果——模型把兜底说明当成工具真的返回了东西，
// 就会基于不存在的输出继续往下编。
func repairNote(id, pendingID string) string {
	if id != "" && id == pendingID {
		return `{"note":"用户没有回答这个提问，会话在这里中断过，该提问作废。不要继续等待回答，也不要假设用户答了什么；按已有信息继续，必要时重新提问。"}`
	}
	return `{"note":"这次工具调用没有留下结果（会话中断），视作未执行。如果仍然需要，请重新调用。"}`
}

// setToolAnswer 把某个 tool_call_id 的应答内容设成 content：已有就**替换**，没有才追加。
//
// 用"设置"而不是"追加"是必要的：重复回答（用户点了两次提交、或上一次回答因 400 失败过）
// 会让同一个 tool_call_id 出现两条应答，而 API 同样会拒。
func setToolAnswer(messages []openai.ChatCompletionMessageParamUnion, toolCallID, content string) []openai.ChatCompletionMessageParamUnion {
	if toolCallID == "" {
		return messages
	}
	out := make([]openai.ChatCompletionMessageParamUnion, len(messages))
	copy(out, messages)
	for i := range out {
		if toolMsgCallID(out[i]) == toolCallID {
			out[i] = openai.ToolMessage(content, toolCallID)
			return out
		}
	}
	return append(out, openai.ToolMessage(content, toolCallID))
}

// askUserArgsOf 取某次工具调用发出去的参数原文（对 ask_user 就是它的问题列表 JSON）。
//
// 用途是让提问卡片能被重建：页面一刷新卡片就没了，而硬拦截又只接受"回答"，
// 不能重建的话那个会话就永远答不上——用户只能把它弃掉。
func askUserArgsOf(messages []openai.ChatCompletionMessageParamUnion, toolCallID string) string {
	for _, m := range messages {
		if m.OfAssistant == nil {
			continue
		}
		for _, tc := range m.OfAssistant.ToolCalls {
			if id := tc.GetID(); id == nil || *id != toolCallID {
				continue
			}
			if fn := tc.GetFunction(); fn != nil {
				return fn.Arguments
			}
		}
	}
	return ""
}

// PendingAskQuestions 返回暂停中那条提问的原始参数；没有暂停时返回空串。
// handler 据此让页面在刷新后把提问卡片重建出来。
func (as *AgentService) PendingAskQuestions(ctx context.Context, userID, sessionID uint) (string, error) {
	sess, messages, err := as.loadSession(ctx, userID, sessionID)
	if err != nil {
		return "", err
	}
	pending := parsePendingAsk(sess)
	if pending.ToolCallID == "" {
		return "", nil
	}
	return askUserArgsOf(messages, pending.ToolCallID), nil
}
