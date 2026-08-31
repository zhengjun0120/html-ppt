// 最小示例：连接模型 → 发送请求 → 接收输出
// 运行前设置环境变量 OPENAI_API_KEY（第三方厂商另设 OPENAI_BASE_URL 或用 WithBaseURL）
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func main() {
	// 1. 连接：构造客户端
	client := openai.NewClient(
		option.WithAPIKey("b18435c2ac38428893828c20b5e61c6c.NobEtbpjyrKG1VJC"),
		// 适配第三方兼容厂商时改端点，例如：
		option.WithBaseURL("https://open.bigmodel.cn/api/paas/v4/"),
	)

	// 总超时由 context 控制（覆盖所有重试）
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 2. 发送请求：行为参数都在请求结构体里
	completion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.ChatModel("glm-5.3-flash"), // 任意字符串也行，如 openai.ChatModel("deepseek-chat")
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("你是一个简洁的助手，用一句话回答。"),
			openai.UserMessage("用一句话介绍 Go 语言"),
		},
		Temperature: openai.Float(0.7),
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			fmt.Printf("API 错误 %d: %v\n", apierr.StatusCode, apierr)
		} else {
			fmt.Println("网络或其他错误:", err)
		}
		os.Exit(1)
	}

	// 3. 接收输出
	fmt.Println("回复:", completion.Choices[0].Message.Content)
	fmt.Printf("结束原因: %s, token: 输入 %d + 输出 %d\n",
		completion.Choices[0].FinishReason,
		completion.Usage.PromptTokens, completion.Usage.CompletionTokens)
}
