package agent

import (
	"context"
	"errors"
	"log"

	"html-ppt/backend/internal/authctx"
	"html-ppt/backend/internal/config"
	"html-ppt/backend/internal/cryptox"
	"html-ppt/backend/internal/service/deck"
	"html-ppt/backend/internal/store"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type AgentService struct {
	ModelClient *openai.Client // 服务器默认 key（免费额度模式 / BYOK 失败兜底）
	BaseURL     string         // BYOK 复用同一供应商入口（用户自选供应商留到阶段5.5）
	ModelID     string
	MaxToken    int64
	DeckService *deck.Service
	Tools       []openai.ChatCompletionToolUnionParam
	Exec        map[string]ToolFunc
	st          *store.Store // BYOK：查用户密钥密文
	box         *cryptox.Box // BYOK：解密；nil = BYOK 关闭
}

var agentServer *AgentService

func InitAgentModel(cfg config.LLM, st *store.Store, box *cryptox.Box, deckService *deck.Service) error {

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
		BaseURL:     cfg.BaseURL,
		ModelID:     cfg.ModelID,
		MaxToken:    cfg.MaxToken,
		DeckService: deckService,
		Exec:        make(map[string]ToolFunc),
		st:          st,
		box:         box,
	}

	tools := agentServer.buildTools()
	for name, t := range tools {
		agentServer.Tools = append(agentServer.Tools, t.Definition)
		agentServer.Exec[name] = t.Execute
	}

	return nil
}

// clientFor 返回本次对话使用的 LLM 客户端：用户配置了自带 API Key（BYOK）
// 就解密用户的来用，否则用服务器默认 key。每个请求解密一次，
// 明文 key 不缓存在内存里，也不打日志。
func (as *AgentService) clientFor(ctx context.Context) *openai.Client {
	if as.st == nil || as.box == nil {
		return as.ModelClient
	}
	uid, ok := authctx.UserID(ctx)
	if !ok {
		return as.ModelClient
	}
	var u store.User
	if err := as.st.DB.Select("api_key_enc").First(&u, uid).Error; err != nil || u.APIKeyEnc == "" {
		return as.ModelClient
	}
	key, err := as.box.Open(u.APIKeyEnc)
	if err != nil {
		log.Printf("[warn] 解密用户 %d 的 API Key 失败，回退服务器默认 key: %v", uid, err)
		return as.ModelClient
	}
	client := openai.NewClient(
		option.WithAPIKey(key),
		option.WithBaseURL(as.BaseURL),
		option.WithMaxRetries(3),
	)
	return &client
}

func GetAgentService() *AgentService {
	return agentServer
}
