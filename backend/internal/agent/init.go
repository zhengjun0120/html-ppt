package agent

import (
	"errors"

	"html-ppt/backend/internal/config"
	"html-ppt/backend/internal/service/deck"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type AgentService struct {
	ModelClient *openai.Client
	ModelID     string
	MaxToken    int64
	DeckService *deck.Service
	Tools []openai.ChatCompletionToolUnionParam
	Exec map[string]ToolFunc
}

var agentServer *AgentService

func InitAgentModel(cfg config.LLM,deckService *deck.Service) error {

	if cfg.APIKey == "" {
		return errors.New("LLM api_key 未配置（config.yaml 或环境变量 LLM_API_KEY）")
	}

	client := openai.NewClient(
		option.WithAPIKey(cfg.APIKey),
		option.WithBaseURL(cfg.BaseURL),
		option.WithMaxRetries(3),
	)

	agentServer = &AgentService{
		ModelClient: &client,
		ModelID:     cfg.ModelID,
		MaxToken:    cfg.MaxToken,
		DeckService: deckService,
		Exec: make(map[string]ToolFunc),
	}

	tools := agentServer.buildTools()
	for name,t := range tools{
		agentServer.Tools = append(agentServer.Tools,t.Definition)
		agentServer.Exec[name] = t.Execute
	}

	return nil
}

func GetAgentService() *AgentService{
	return agentServer
}
