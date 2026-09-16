package agent

import (
	"context"
	"errors"
	"log"
	"sort"
	"strings"

	"html-ppt/backend/internal/authctx"
	"html-ppt/backend/internal/config"
	"html-ppt/backend/internal/cryptox"
	"html-ppt/backend/internal/service/deck"
	"html-ppt/backend/internal/store"
	"html-ppt/backend/internal/trace"
	"html-ppt/backend/internal/vision"

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
	// MaxPerRun 工具名 → 一次 run 里的最大执行次数（0/缺失 = 不限）。
	// 它是"审查→修复→再审"这类无收敛循环的唯一硬闸门：提示词是软约束，
	// 实测的 22 轮 run 就是靠 maxTurns 兜底才停下来的。见 Tool.MaxPerRun 与 execTool。
	MaxPerRun map[string]int
	st          *store.Store // BYOK：查用户密钥密文
	box         *cryptox.Box // BYOK：解密；nil = BYOK 关闭
	// CustomCSS = features.custom_css：关闭时 buildTools 不挂载自定义样式三件套
	CustomCSS bool
	// AssetsDir 共享静态资源目录（web/assets）：read_component 从这里读组件库，只读
	AssetsDir string

	serverKey        string // 服务器默认apiKey
	WebSearch        bool   // 是否开启联网搜索
	AnthropicBaseURL string // deepseek 只支持 anthropic格式的联网搜索

	Vision        bool
	ChromePath    string
	VisionBaseURL string
	VisionGrants  *vision.Grants

	// TraceCfg 观测配置（features.trace + trace.*）。关掉时循环里一切 Emit 都是空操作，
	// 连上下文都不会被序列化（见 trace.Active 的用途），所以关掉是真零开销。
	TraceCfg trace.Config
}

var agentServer *AgentService

func InitAgentModel(cfg config.LLM, st *store.Store, box *cryptox.Box, deckService *deck.Service, features config.Features, assetsDir string, traceCfg trace.Config) error {

	if cfg.APIKey == "" {
		return errors.New("LLM api_key 未配置（config.yaml 或环境变量 LLM_API_KEY）")
	}

	client := openai.NewClient(
		option.WithAPIKey(cfg.APIKey),
		option.WithBaseURL(cfg.BaseURL),
		option.WithMaxRetries(3),
	)

	agentServer = &AgentService{
		ModelClient:      &client,
		BaseURL:          cfg.BaseURL,
		ModelID:          cfg.ModelID,
		MaxToken:         cfg.MaxToken,
		DeckService:      deckService,
		Exec:             make(map[string]ToolFunc),
		st:               st,
		box:              box,
		CustomCSS:        features.CustomCSS,
		AssetsDir:        assetsDir,
		serverKey:        cfg.APIKey,
		WebSearch:        features.WebSearch,
		AnthropicBaseURL: cfg.AnthropicBaseURL,
		Vision:           features.Vision,
		// ChromePath: ,
		TraceCfg: traceCfg,
	}

	tools := agentServer.buildTools()
	for name, t := range tools {
		agentServer.Tools = append(agentServer.Tools, t.Definition)
		agentServer.Exec[name] = t.Execute
		if t.MaxPerRun > 0 {
			if agentServer.MaxPerRun == nil {
				agentServer.MaxPerRun = make(map[string]int)
			}
			agentServer.MaxPerRun[name] = t.MaxPerRun
		}
	}
	// 启动时把挂载了哪些工具打出来：features 开关的效果要可见
	//（关掉 custom_css 后这里就不该再有 read/update_custom_css 和 read_component）
	names := make([]string, 0, len(tools))
	for name := range tools {
		names = append(names, name)
	}
	sort.Strings(names)
	log.Printf("[info] agent 挂载 %d 个工具: %s", len(names), strings.Join(names, " "))
	log.Printf("[info] 组件库目录（read_component 只读）: %s", assetsDir)
	log.Printf("[info] 联网搜索（DeepSeek 服务端 web_search，走 Anthropic 兼容入口）: %v", features.WebSearch)
	log.Printf("[info] agent 观测（trace）: %v，落盘目录 %s（截图 %v，每会话保留 %d 个 run）",
		traceCfg.Enabled, traceCfg.Dir, traceCfg.CaptureImages, traceCfg.RetainRuns)

	return nil
}

// resolveKey 返回本次会话使用的key
func (as *AgentService) resolveKey(ctx context.Context) (key string, byok bool) {
	if as.st == nil || as.box == nil {
		return as.serverKey, false
	}

	uid, ok := authctx.UserID(ctx)
	if !ok {
		return as.serverKey, false
	}

	var u store.User
	if err := as.st.DB.Model(&store.User{}).Where("id = ?", uid).Select("api_key_enc").First(&u).Error; err != nil || u.APIKeyEnc == "" {
		return as.serverKey, false
	}

	plain, err := as.box.Open(u.APIKeyEnc)
	if err != nil {
		log.Printf("[warn] 解密用户 %d 的 API Key 失败，回退服务器默认 key: %v", uid, err)
		return as.serverKey, false
	}

	return plain, true
}

// clientFor 返回本次对话使用的 LLM 客户端：用户配置了自带 API Key（BYOK）
// 就解密用户的来用，否则用服务器默认 key。每个请求解密一次，
// 明文 key 不缓存在内存里，也不打日志。
func (as *AgentService) clientFor(ctx context.Context) *openai.Client {
	key, byok := as.resolveKey(ctx)

	if !byok {
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
