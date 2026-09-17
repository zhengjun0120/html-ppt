package handler

import (
	"context"

	"html-ppt/backend/internal/agent"
	"html-ppt/backend/internal/service/auth"
	"html-ppt/backend/internal/service/deck"
	"html-ppt/backend/internal/service/template"
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
	// traces 观测记录的读取端（写端在 agent 里）。**不受 features.trace 影响也要装配**：
	// 关掉的是"继续记录"，已经落盘的记录应该照样能看，
	// 否则调一次开关就会把之前跑出来的东西变成读不到的孤儿文件。
	traces *trace.Store
}

// Exporter 导出服务的最小接口。返回值用 any：export.Result 的形状 handler 不关心
// （原样序列化给前端），避免 handler 直接 import export 包造成的装配环。
type Exporter interface {
	Export(ctx context.Context, uid uint, deckID, format string) (any, error)
}

func New(st *store.Store, decks *deck.Service, agentSvc *agent.AgentService, authSvc *auth.Service, renderGrantsSvc *vision.Grants, traces *trace.Store, templates *template.Registry, exporter Exporter) *Handler {
	return &Handler{st: st, decks: decks, agent: agentSvc, auth: authSvc, renderGrants: renderGrantsSvc, traces: traces, templates: templates, exporter: exporter}
}

// TemplatesAvailable 模板库是否可用（router 据此决定挂不挂预览静态路由）。
func (h *Handler) TemplatesAvailable() bool { return h.templates != nil }
