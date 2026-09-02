// ssedemo：Gin + SSE 流式接口，把模型输出实时转发给前端。
//
// 运行（在 backend 目录）：
//
//	export LLM_API_KEY=sk-xxx            # 或 OPENAI_API_KEY
//	export LLM_BASE_URL=https://api.deepseek.com/v1   # 可选
//	go run ./cmd/ssedemo
//
// 测试：
//
//	curl -N "http://localhost:8081/api/chat/stream?msg=用一句话介绍SSE"
//
// -N 禁用 curl 自身的缓冲，否则看不到逐字输出。
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// sseEvent 推给前端的事件结构。
// 用 JSON 帧而不是裸文本，前端好按 type 分发（打字机/结束/报错）。
type sseEvent struct {
	Type        string `json:"type"`                   // delta | done | error
	Text        string `json:"text,omitempty"`         // delta 的增量文本 / error 的信息
	TotalTokens int64  `json:"total_tokens,omitempty"` // done 时带上用量
}

// writeSSE 写一帧数据并立刻冲刷。
// SSE 帧格式固定：data: <payload>\n\n —— 必须以空行(\n\n)结尾，客户端才认为一帧结束。
// 返回 false 表示客户端已断开，调用方应立即停止。
func writeSSE(c *gin.Context, payload any) bool {
	raw, err := json.Marshal(payload)
	if err != nil {
		return false
	}
	if _, err := fmt.Fprintf(c.Writer, "data: %s\n\n", raw); err != nil {
		return false // 写失败 ≈ 客户端断开
	}
	c.Writer.Flush() // 关键：gin 默认有缓冲，不 Flush 的话前端会攒一大坨才收到
	return true
}

func main() {
	apiKey := os.Getenv("LLM_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	if apiKey == "" {
		panic("请先设置环境变量 LLM_API_KEY（或 OPENAI_API_KEY）")
	}

	opts := []option.RequestOption{option.WithAPIKey(apiKey)}
	if base := os.Getenv("LLM_BASE_URL"); base != "" {
		opts = append(opts, option.WithBaseURL(base))
	}
	model := os.Getenv("LLM_MODEL")
	if model == "" {
		model = "deepseek-chat"
	}
	client := openai.NewClient(opts...) // 值类型，见 agentdemo 里的说明

	r := gin.Default()

	r.GET("/api/chat/stream", func(c *gin.Context) {
		msg := c.Query("msg")
		if msg == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 msg 参数"})
			return
		}

		// ---- 第 1 步：声明这是 SSE 响应 ----
		// text/event-stream 是协议规定的 Content-Type；
		// no-cache 防中间层缓存；X-Accel-Buffering=no 专门关 nginx 的响应缓冲，
		// 不加的话部署在 nginx 后面时前端会等到请求结束才一次性收到所有数据。
		w := c.Writer
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("X-Accel-Buffering", "no")
		w.WriteHeader(http.StatusOK)

		// ---- 第 2 步：把 gin 请求的 context 传给 SDK ----
		// 客户端关掉页面/断开连接时，net/http 会取消这个 context，
		// SDK 的流会立刻停止——否则模型还在生成、你在白烧 token。
		ctx := c.Request.Context()

		params := openai.ChatCompletionNewParams{
			Model: openai.ChatModel(model),
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(msg),
			},
			StreamOptions: openai.ChatCompletionStreamOptionsParam{
				IncludeUsage: openai.Bool(true), // 流式下默认没有 usage
			},
		}

		// ---- 第 3 步：边收边转发，同时累加 ----
		stream := client.Chat.Completions.NewStreaming(ctx, params)
		defer stream.Close()

		acc := openai.ChatCompletionAccumulator{}
		for stream.Next() {
			chunk := stream.Current()
			if !acc.AddChunk(chunk) {
				continue
			}
			// 只有文本增量要转发；tool_calls 增量由 AddChunk 静默拼接（本 demo 无工具）
			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
				ok := writeSSE(c, sseEvent{Type: "delta", Text: chunk.Choices[0].Delta.Content})
				if !ok {
					return // 客户端断开，结束处理（ctx 已取消，SDK 流也会停）
				}
			}
		}

		// ---- 第 4 步：流结束，区分"客户端断开"和"模型端出错" ----
		if err := stream.Err(); err != nil {
			if ctx.Err() != nil {
				return // 客户端主动断开导致的取消，不是错误，也不用再写帧
			}
			_ = writeSSE(c, sseEvent{Type: "error", Text: err.Error()})
			return
		}

		// 累加器里有完整结果和用量（IncludeUsage 的成果）
		_ = writeSSE(c, sseEvent{Type: "done", TotalTokens: acc.Usage.TotalTokens})
		// [DONE] 帧沿用 OpenAI 的结束约定，前端见到它就知道整个流结束了
		fmt.Fprint(w, "data: [DONE]\n\n")
		w.Flush()
	})

	addr := os.Getenv("SSE_ADDR")
	if addr == "" {
		addr = ":8081"
	}
	fmt.Printf("SSE demo 监听 %s，测试：curl -N \"http://localhost%s/api/chat/stream?msg=你好\"\n", addr, addr)
	_ = r.Run(addr)
}
