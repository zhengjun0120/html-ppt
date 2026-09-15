package handler

import (
	"html-ppt/backend/internal/agent"
	"html-ppt/backend/internal/service/auth"
	"html-ppt/backend/internal/service/deck"
	"html-ppt/backend/internal/store"
	"html-ppt/backend/internal/vision"
)

// Handler 持有所有路由处理函数需要的依赖（struct-based handler 模式）。
// 装配点唯一：新增依赖只改这里和 New，路由注册一行不用动。
// 规矩：结构体里只放依赖，且 handler 只调 service，不越层摸 store / 文件系统。
type Handler struct {
	decks *deck.Service
	st    *store.Store // 可能为 nil（数据库降级模式），使用处需判空
	agent *agent.AgentService
	auth  *auth.Service
	renderGrants *vision.Grants
}

func New(st *store.Store, decks *deck.Service, agentSvc *agent.AgentService, authSvc *auth.Service,renderGrantsSvc *vision.Grants) *Handler {
	return &Handler{st: st, decks: decks, agent: agentSvc, auth: authSvc,renderGrants: renderGrantsSvc}
}
