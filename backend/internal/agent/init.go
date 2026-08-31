package agent

import (
	"html-ppt/backend/internal/config"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type AgentServer struct {
	ModelClient *openai.Client
	ModelID string
	MaxToken int64
}

var agentServer *AgentServer

func InitAgentModel(cfg config.LLM)error{

	if cfg.APIKey == ""{
		panic("默认apiKey为空")
	} 

	client := openai.NewClient(
		option.WithAPIKey(cfg.APIKey),
		option.WithBaseURL(cfg.BaseURL),
		option.WithMaxRetries(3),
		option.WithRequestTimeout(30 *time.Second),
	)

	agentServer  = &AgentServer{
		ModelClient: &client,
		ModelID: cfg.ModelID,
		MaxToken: cfg.MaxToken,
	}
	return nil
}