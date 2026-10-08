package handler

import (
	"context"

	"html-ppt/backend/internal/agent"
	"html-ppt/backend/internal/service/auth"
	"html-ppt/backend/internal/service/deck"
	"html-ppt/backend/internal/service/template"
	"html-ppt/backend/internal/service/tplsuggest"
	"html-ppt/backend/internal/service/tplthumb"
	"html-ppt/backend/internal/service/usage"
	"html-ppt/backend/internal/service/usertpl"
	"html-ppt/backend/internal/store"
	"html-ppt/backend/internal/trace"
	"html-ppt/backend/internal/vision"
)

// Handler 持有所有路由处理函数需要的依赖（struct-based handler 模式）。
// 装配点唯一：新增依赖只改这里和 New，路由注册一行不用动。
// 规矩：结构体里只放依赖，且 handler 只调 service，不越层摸 store / 文件系统。
type Handler struct {
	decks        *deck.Service
	st           *store.Store // 可能为 nil（数据库降级模式），使用处需判空
	agent        *agent.AgentService
	auth         *auth.Service
	renderGrants *vision.Grants
	// templates deck-v2 模板注册表（可能为 nil：模板库损坏时不阻塞整个服务，
	// 只有模板相关接口不可用）。detail 同上。
	templates *template.Registry
	// exporter deck-v2 导出服务（可能为 nil：Chrome 不可用时导出不可用）
	exporter Exporter
	// usertpl 用户自定义模板服务（可能为 nil：数据库降级模式下不可用）
	usertpl UserTemplateService
	// traces 观测记录的读取端（写端在 agent 里）。**不受 features.trace 影响也要装配**：
	// 关掉的是"继续记录"，已经落盘的记录应该照样能看，
	// 否则调一次开关就会把之前跑出来的东西变成读不到的孤儿文件。
	traces *trace.Store
	// suggest 模板推荐服务（New 里从 decks+templates 装配；templates 为 nil 时
	// 其接口返回"模板库不可用"）
	suggest *tplsuggest.Service
	// thumbs 模板缩略图服务（可能为 nil：Chrome 不可用或未装配——thumb 端点
	// 返回 404，前端回退活 iframe）
	thumbs *tplthumb.Service
	// usage 用量账本服务（可能为 nil：数据库降级模式下不可用，端点返回 503）
	usage *usage.Service
}

// Exporter 导出服务的最小接口。返回值用 any：export.Result 的形状 handler 不关心
// （原样序列化给前端），避免 handler 直接 import export 包造成的装配环。
type Exporter interface {
	Export(ctx context.Context, uid uint, deckID, format string) (any, error)
	// EnsureThumbs 缩略图缓存就绪（过期/缺失时整本重渲一次），返回缩略图目录。
	EnsureThumbs(ctx context.Context, uid uint, deckID string) (string, error)
}

// UserTemplateService 用户自定义模板的最小接口（避免 handler 直接依赖 usertpl 包；
// 返回的具体行类型是 store.UserTemplate——handler 本就依赖 store）。
type UserTemplateService interface {
	Fork(userID uint, baseID, name string) (*store.UserTemplate, error)
	List(userID uint) ([]store.UserTemplate, error)
	Community() ([]map[string]any, error)
	GetReadable(userID uint, id string) (*store.UserTemplate, error)
	UpdateMeta(userID uint, id, name, description string) error
	Delete(userID uint, id string) error
	Publish(ctx context.Context, userID uint, id string) (*usertpl.PublishReport, error)
	Unpublish(userID uint, id string) error
	// Checkup 质量体检（2026-09-28 从发布门禁降级而来）：渲染量测，只报告不拦发布
	Checkup(ctx context.Context, userID uint, id string) (*usertpl.PublishReport, error)
	// Customize emit 为 SSE 事件出口；同步调用传 nil。images 是用户附图（data URL，服务层校验）
	Customize(ctx context.Context, userID uint, id, message string, images []string, llm usertpl.LLM, emit func(usertpl.CustEvent) error) (string, error)
	// 编辑器（deck-editor-plan §4.2）：编辑态读取要 Dir 定位 demo 文件
	Dir(id string) string
	GetOwned(userID uint, id string) (*store.UserTemplate, error)
	SaveIndexHTML(userID uint, id, html string) error
	SaveStyleCSS(userID uint, id, css string) error
	// 结构契约（版式面板）：读挂载态契约 / 改版式元数据（roles/名称/用途）
	StructureContract(userID uint, id string) (*usertpl.StructureContractView, error)
	UpdateLayoutMeta(userID uint, id, layoutID string, patch usertpl.LayoutMetaPatch) (string, error)
	// Peek 无归属读行（受控公开端点先看状态再决定鉴权，user-template-history-plan.md §3.5）
	Peek(id string) (*store.UserTemplate, error)
	// 历史版本（user-template-history-plan.md §4）：列表/回滚/删单版/清空
	ListUTVersions(userID uint, id string) ([]usertpl.UTVersionMeta, error)
	RestoreUTVersion(userID uint, id, version string) error
	DeleteUTVersion(userID uint, id, version string) error
	ClearUTHistory(userID uint, id string) (int, error)
}

func New(st *store.Store, decks *deck.Service, agentSvc *agent.AgentService, authSvc *auth.Service, renderGrantsSvc *vision.Grants, traces *trace.Store, templates *template.Registry, exporter Exporter, utpl UserTemplateService, traceCfg trace.Config, thumbs *tplthumb.Service, usageSvc *usage.Service) *Handler {
	return &Handler{st: st, decks: decks, agent: agentSvc, auth: authSvc, renderGrants: renderGrantsSvc, traces: traces, templates: templates, exporter: exporter, usertpl: utpl, suggest: tplsuggest.New(decks, templates, traceCfg), thumbs: thumbs, usage: usageSvc}
}

// TemplatesAvailable 模板库是否可用（router 据此决定挂不挂预览静态路由）。
func (h *Handler) TemplatesAvailable() bool { return h.templates != nil }

// UsertplAvailable 用户模板子系统是否可用（router 据此挂受控伺服路由）。
func (h *Handler) UsertplAvailable() bool { return h.usertpl != nil }
