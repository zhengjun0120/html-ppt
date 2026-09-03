package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"html-ppt/backend/internal/agent"
	"html-ppt/backend/internal/config"
	"html-ppt/backend/internal/handler"
	"html-ppt/backend/internal/router"
	"html-ppt/backend/internal/service/deck"
	"html-ppt/backend/internal/store"
)

// main 保持极薄：真正的启动逻辑在 run() 里，方便统一处理错误退出码。
func main() {
	if err := run(); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}

func run() error {
	// —— 1. 配置 ——
	// Load 内部会向上查找 config.yaml，从任何子目录启动都能找到
	cfg, err := config.Load("config.yaml")
	if err != nil {
		return err
	}

	// —— 2. 依赖装配（依赖注入的核心：每层只接收自己需要的东西，不读全局）——
	// store 层可选：数据库连不上时降级为无 DB 模式，开发体验优先。
	var st *store.Store
	if cfg.DB.Enabled {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		st, err = store.Open(ctx, cfg.DB)
		cancel()
		if err != nil {
			log.Printf("[warn] 数据库连接失败，以无数据库模式运行: %v", err)
			st = nil
		} else {
			log.Printf("[info] 数据库已连接: %s@%s:%d/%s", cfg.DB.User, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)
		}
	}

	

	deckSvc := deck.New(cfg.Data.Dir)

	if err := agent.InitAgentModel(cfg.LLM, deckSvc); err != nil {
		return fmt.Errorf("初始化 agent: %w", err)
	}
	agentSvc := agent.GetAgentService()

	h := handler.New(st, deckSvc,agentSvc)
	engine := router.New(cfg, h)
	srv := &http.Server{Addr: cfg.Server.Addr, Handler: engine}

	// —— 3. 启动 ——
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("[info] HTTP server listening on %s", cfg.Server.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()

	// —— 4. 优雅停机：Ctrl+C 后给在途请求最多 5s 收尾 ——
	<-ctx.Done()
	log.Println("[info] shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
