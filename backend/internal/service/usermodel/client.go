package usermodel

import (
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// NewClient 按用户给的入口构造 OpenAI 兼容客户端。key 可空（本地网关）。
// 重试与 agent 平台客户端同款（3 次）：生成链路一轮几十秒，网络抖动不该直接判死。
// 导出给 agent.resolveLLM 复用——自选模型的调用客户端只有这一处构造。
func NewClient(baseURL, apiKey string, retries int) *openai.Client {
	opts := []option.RequestOption{option.WithBaseURL(baseURL)}
	if apiKey != "" {
		opts = append(opts, option.WithAPIKey(apiKey))
	}
	if retries > 0 {
		opts = append(opts, option.WithMaxRetries(retries))
	}
	c := openai.NewClient(opts...)
	return &c
}

// newTestRequest 最小的连通性探针。
func newTestRequest(modelID string) openai.ChatCompletionNewParams {
	return openai.ChatCompletionNewParams{
		Model: openai.ChatModel(modelID),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("连通性测试。只回复两个字母：ok"),
			openai.UserMessage("ping"),
		},
	}
}
