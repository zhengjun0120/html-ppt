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
// 用 insert_slide（没有退还语义的普通配额工具）测通用闸门；review_slides 的
// "只数真看图"语义在下面单独测。
func TestExecToolQuotaBlocksAfterMaxPerRun(t *testing.T) {
	calls := 0
	as := &AgentService{
		Exec: map[string]ToolFunc{
			"insert_slide": func(ctx context.Context, arguments string) (string, error) {
				calls++
				return fmt.Sprintf(`{"exec":%d}`, calls), nil
			},
		},
		MaxPerRun: map[string]int{"insert_slide": 3},
	}
	rr := newRunRecorder()
	tool := openai.ChatCompletionMessageToolCallUnion{
		ID:       "call_1",
		Type:     "function",
		Function: openai.ChatCompletionMessageFunctionToolCallFunction{Name: "insert_slide", Arguments: "{}"},
	}
	emit := func(StreamEvent) error { return nil }

	for i := 0; i < 3; i++ {
		res, err := as.execTool(context.Background(), tool, emit, rr, nil)
		if err != nil || res != fmt.Sprintf(`{"exec":%d}`, i+1) {
			t.Fatalf("第 %d 次调用应正常执行: res=%q err=%v", i+1, res, err)
		}
	}

	res, err := as.execTool(context.Background(), tool, emit, rr, nil)
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
	res2, _ := as.execTool(context.Background(), tool, emit, rr, nil)
	if !strings.Contains(res2, "配额用完") {
		t.Errorf("配额用尽后的重试也应被拦下, got %q", res2)
	}
}

// TestExecToolReviewQuotaOnlyCountsRealReviews review_slides 的配额只数
// "真的看了图"的调用（结果带 visionReportMarker）：pages 留空的数字复查、
// 看图失败的降级结果都退还计数。实测（trace 1789613158323-30a4）没有这条时，
// 两次免费复查加一次看图失败就把配额耗尽，模型想重试真审查时被拒。
func TestExecToolReviewQuotaOnlyCountsRealReviews(t *testing.T) {
	as := &AgentService{
		Exec: map[string]ToolFunc{
			"review_slides": func(ctx context.Context, arguments string) (string, error) {
				// 假执行器按参数模拟两种结果：pages 留空 = 只回数字；带了 pages = 看图成功
				if strings.Contains(arguments, "pages") {
					return "xx" + visionReportMarker + "第 1 页｜【硬】｜测试", nil
				}
				return "共 3 页，画布 1244*700 ……", nil
			},
		},
		MaxPerRun: map[string]int{"review_slides": 3},
	}
	rr := newRunRecorder()
	emit := func(StreamEvent) error { return nil }
	call := func(args string) string {
		tool := openai.ChatCompletionMessageToolCallUnion{
			ID: "call_1", Type: "function",
			Function: openai.ChatCompletionMessageFunctionToolCallFunction{Name: "review_slides", Arguments: args},
		}
		res, err := as.execTool(context.Background(), tool, emit, rr, nil)
		if err != nil {
			t.Fatalf("执行不应出错: %v", err)
		}
		return res
	}

	// 两次数字复查：执行但不占配额
	for i := 0; i < 2; i++ {
		if res := call(`{"deck_id":"deck-0001"}`); strings.Contains(res, visionReportMarker) {
			t.Fatal("数字复查不该带看图标记")
		}
	}
	// 一次真看图：占配额
	if res := call(`{"deck_id":"deck-0001","pages":[1,2]}`); !strings.Contains(res, visionReportMarker) {
		t.Fatalf("看图结果应带标记: %q", res)
	}
	// 再来两次真看图：到上限
	call(`{"deck_id":"deck-0001","pages":[1]}`)
	if res := call(`{"deck_id":"deck-0001","pages":[2]}`); !strings.Contains(res, visionReportMarker) {
		t.Fatalf("第 3 次真看图应执行: %q", res)
	}
	// 第 4 次真看图：被拦
	if res := call(`{"deck_id":"deck-0001","pages":[3]}`); !strings.Contains(res, "配额用完") {
		t.Errorf("真看图超过 3 次应被拦下, got %q", res)
	}
	// 此时来一次数字复查：仍应放行（退还语义保证"免费复查随时可调"的承诺）
	if res := call(`{"deck_id":"deck-0001"}`); strings.Contains(res, "配额用完") {
		t.Errorf("数字复查不应被配额拦截, got %q", res)
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
		if _, err := as.execTool(context.Background(), tool, func(StreamEvent) error { return nil }, rr, nil); err != nil {
			t.Fatalf("第 %d 次调用不应出错: %v", i+1, err)
		}
	}
	if calls != 5 {
		t.Errorf("无配额工具被执行了 %d 次, 要 5", calls)
	}
}
