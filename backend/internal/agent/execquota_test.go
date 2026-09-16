package agent

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
)

// TestExecToolQuotaBlocksAfterMaxPerRun 配额是"审查→修复→再审"死循环的唯一硬闸门：
// 前 max 次放行执行，第 max+1 次起不执行、把"配额用完"的指引当结果还给模型。
// 它同时守两件事：真的拦住了（执行次数不超），且拦截不是故障（callErr 为空，
// 走 tool_call 而不是 tool_error——模型需要的是"接下来怎么办"，不是红色报错）。
func TestExecToolQuotaBlocksAfterMaxPerRun(t *testing.T) {
	calls := 0
	as := &AgentService{
		Exec: map[string]ToolFunc{
			"review_slides": func(ctx context.Context, arguments string) (string, error) {
				calls++
				return fmt.Sprintf(`{"exec":%d}`, calls), nil
			},
		},
		MaxPerRun: map[string]int{"review_slides": 3},
	}
	rr := newRunRecorder()
	tool := openai.ChatCompletionMessageToolCallUnion{
		ID:       "call_1",
		Type:     "function",
		Function: openai.ChatCompletionMessageFunctionToolCallFunction{Name: "review_slides", Arguments: "{}"},
	}
	emit := func(StreamEvent) error { return nil }

	for i := 0; i < 3; i++ {
		res, err := as.execTool(context.Background(), tool, emit, rr)
		if err != nil || res != fmt.Sprintf(`{"exec":%d}`, i+1) {
			t.Fatalf("第 %d 次调用应正常执行: res=%q err=%v", i+1, res, err)
		}
	}

	res, err := as.execTool(context.Background(), tool, emit, rr)
	if err != nil {
		t.Fatalf("配额拦截不应向循环返回错误: %v", err)
	}
	if !strings.Contains(res, "配额用完") || strings.Contains(res, `"exec"`) {
		t.Errorf("第 4 次调用应被拦下并给出指引, got %q", res)
	}
	if calls != 3 {
		t.Errorf("实际执行了 %d 次, 要 3", calls)
	}

	// 拦截后的重试同样被拦，且措辞稳定——模型可能会再试一次
	res2, _ := as.execTool(context.Background(), tool, emit, rr)
	if !strings.Contains(res2, "配额用完") {
		t.Errorf("配额用尽后的重试也应被拦下, got %q", res2)
	}
}

// TestExecToolNoQuotaMeansUnlimited 没配额的工具（MaxPerRun 缺省/为 0）不受影响——
// 配额机制只该作用于声明了配额的工具，不能顺手把别的工具也限了。
func TestExecToolNoQuotaMeansUnlimited(t *testing.T) {
	calls := 0
	as := &AgentService{
		Exec: map[string]ToolFunc{
			"read_slide": func(ctx context.Context, arguments string) (string, error) {
				calls++
				return "ok", nil
			},
		},
		// MaxPerRun 为 nil：read_slide 没声明配额
	}
	rr := newRunRecorder()
	tool := openai.ChatCompletionMessageToolCallUnion{
		ID:       "call_1",
		Type:     "function",
		Function: openai.ChatCompletionMessageFunctionToolCallFunction{Name: "read_slide", Arguments: "{}"},
	}
	for i := 0; i < 5; i++ {
		if _, err := as.execTool(context.Background(), tool, func(StreamEvent) error { return nil }, rr); err != nil {
			t.Fatalf("第 %d 次调用不应出错: %v", i+1, err)
		}
	}
	if calls != 5 {
		t.Errorf("无配额工具被执行了 %d 次, 要 5", calls)
	}
}
