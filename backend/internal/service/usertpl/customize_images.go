package usertpl

// 定制对话的用户附图：把参考图带进模型上下文（多模态消息）。
//
// 生命周期刻意收窄成「只看一轮」：图随本轮 customizeLoop 的每次模型调用可见，
// 本轮成功结束后把该条用户消息降级回纯文本 + 占位说明。聊天补全协议下历史
// 每轮整体重发，一张压缩图 1-2k token，长对话费用会滚涨；模型当轮已经把图
// 消化进回复与文件改动，占位足够后续引用。本轮失败时不降级——图留着，重试
// 那轮模型还能看到（没看过的图说"看过"是撒谎）。
//
// 校验/消息构造/观测脱敏的共享实现在 internal/chatimg（文稿对话与这里有
// 不同的留存策略：图存库、模型侧滑动窗口），本文件只留定制侧的降级策略。

import (
	"fmt"
	"strings"

	"github.com/openai/openai-go/v3"
)

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
