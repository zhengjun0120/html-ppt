package agent

// ClarifyUserMessages 的回归：按会话正序拼接真实用户消息、滤合成指令、
// 保头部截断、坏形状会话跳过、降级模式报哨兵错误。
//
// 内存库是进程级共享的（store.OpenMemory cache=shared），deck_id 用本文件
// 专属的值隔离，不与其它测试的 deck 串味。

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"html-ppt/backend/internal/store"

	"github.com/openai/openai-go/v3"
)

func mustMessagesJSON(t *testing.T, msgs ...openai.ChatCompletionMessageParamUnion) string {
	t.Helper()
	b, err := json.Marshal(msgs)
	if err != nil {
		t.Fatalf("序列化消息: %v", err)
	}
	return string(b)
}

func TestClarifyUserMessages(t *testing.T) {
	ctx := context.Background()
	st, err := store.OpenMemory(ctx)
	if err != nil {
		t.Fatalf("内存库: %v", err)
	}
	const uid = uint(42)
	const deckID = "deck-clarify-ctx"

	// 会话 1（较早）：真实诉求 + 一条合成的开工指令 + 一条助手消息
	sess1 := store.ChatSession{UserID: uid, DeckID: deckID, Title: "澄清"}
	sess1.Messages = mustMessagesJSON(t,
		openai.UserMessage("我要给投资人讲商业模式"),
		openai.AssistantMessage("好的，先确认几个问题"),
		openai.UserMessage(kickoffGenerateMsg), // 合成指令：必须被滤掉
		openai.UserMessage("预算大概五十万"),
	)
	// 会话 2（较晚）：大纲修订 + 一条空消息
	sess2 := store.ChatSession{UserID: uid, DeckID: deckID, Title: "修订"}
	sess2.Messages = mustMessagesJSON(t,
		openai.UserMessage("第3页改成数据页"),
		openai.UserMessage("   "), // 纯空白：滤掉
	)
	// 坏形状会话：JSON 解析失败，跳过不阻断
	sess3 := store.ChatSession{UserID: uid, DeckID: deckID, Title: "坏形状"}
	sess3.Messages = "这不是JSON{{{"
	// 别人的 deck：不该出现在结果里
	other := store.ChatSession{UserID: uid + 1, DeckID: deckID, Title: "别人的"}
	other.Messages = mustMessagesJSON(t, openai.UserMessage("越权内容"))
	// 无关 deck：同理
	unrelated := store.ChatSession{UserID: uid, DeckID: "deck-other", Title: "无关"}
	unrelated.Messages = mustMessagesJSON(t, openai.UserMessage("无关内容"))

	for _, s := range []*store.ChatSession{&sess1, &sess2, &sess3, &other, &unrelated} {
		if err := st.DB.Create(s).Error; err != nil {
			t.Fatalf("造会话 %s: %v", s.Title, err)
		}
	}

	as := &AgentService{st: st}

	t.Run("拼接与过滤", func(t *testing.T) {
		got, err := as.ClarifyUserMessages(ctx, uid, deckID, 0)
		if err != nil {
			t.Fatalf("ClarifyUserMessages: %v", err)
		}
		for _, want := range []string{"我要给投资人讲商业模式", "预算大概五十万", "第3页改成数据页"} {
			if !strings.Contains(got, want) {
				t.Errorf("结果缺 %q，got: %q", want, got)
			}
		}
		for _, banned := range []string{"开始生成", "越权内容", "无关内容"} {
			if strings.Contains(got, banned) {
				t.Errorf("结果不该含 %q", banned)
			}
		}
		// 会话正序：澄清会话的诉求在修订消息之前
		if strings.Index(got, "投资人") > strings.Index(got, "第3页") {
			t.Errorf("顺序不对（澄清应在前）: %q", got)
		}
	})

	t.Run("头部截断不腰斩rune", func(t *testing.T) {
		got, err := as.ClarifyUserMessages(ctx, uid, deckID, 30)
		if err != nil {
			t.Fatalf("ClarifyUserMessages: %v", err)
		}
		if got == "" {
			t.Fatal("截断后不应为空")
		}
		if len([]rune(got)) > 30 {
			t.Fatalf("超长: %d runes", len([]rune(got)))
		}
		if !strings.Contains(got, "投资人") {
			t.Errorf("保头部截断应保留开头的核心诉求: %q", got)
		}
	})

	t.Run("降级模式报哨兵错误", func(t *testing.T) {
		if _, err := (&AgentService{st: nil}).ClarifyUserMessages(ctx, uid, deckID, 0); !errors.Is(err, ErrSessionStoreUnavailable) {
			t.Fatalf("期望 ErrSessionStoreUnavailable, got %v", err)
		}
	})

	t.Run("没有会话返回空串", func(t *testing.T) {
		got, err := as.ClarifyUserMessages(ctx, uid, "deck-none", 0)
		if err != nil || got != "" {
			t.Fatalf("期望空串+nil, got %q err %v", got, err)
		}
	})
}
