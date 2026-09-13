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

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"html-ppt/backend/internal/agent"
	"html-ppt/backend/internal/config"
	"html-ppt/backend/internal/cryptox"
	"html-ppt/backend/internal/handler"
	"html-ppt/backend/internal/router"
	"html-ppt/backend/internal/service/auth"
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
	// 注意：引入用户系统后，降级模式只剩 health 可用（其余接口都要求登录）。
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

	// Redis：验证码存储。连不上只降级"验证码注册"，登录不受影响。
	var rdb *redis.Client
	if cfg.Redis.Addr != "" {
		rdb = redis.NewClient(&redis.Options{Addr: cfg.Redis.Addr, Password: cfg.Redis.Password})
		pctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		err := rdb.Ping(pctx).Err()
		cancel()
		if err != nil {
			log.Printf("[warn] Redis 连接失败，邮箱验证码注册不可用: %v", err)
			rdb = nil
		}
	}

	// AES 主密钥：没配置则 BYOK（用户自带 API Key）功能关闭，其余不受影响
	var box *cryptox.Box
	if cfg.Crypto.AESKey != "" {
		box, err = cryptox.NewBox(cfg.Crypto.AESKey)
		if err != nil {
			return fmt.Errorf("crypto.aes_key 配置错误: %w", err) // 主密钥写坏必须启动即暴露
		}
	}

	// auth 服务的 JWT 密钥同理：没配置则拒绝启动（否则 token 可被伪造）
	if cfg.Auth.JWTSecret == "" {
		return errors.New("auth.jwt_secret 未配置（config.yaml 或环境变量 AUTH_JWT_SECRET）")
	}

	var db *gorm.DB
	if st != nil {
		db = st.DB
	}
	authSvc := auth.New(db, rdb, box, cfg.Auth.JWTSecret, cfg.Auth.TTL(), cfg.SMTP)

	deckSvc := deck.New(cfg.Data.Dir, st)

	if err := agent.InitAgentModel(cfg.LLM, st, box, deckSvc, cfg.Features); err != nil {
		return fmt.Errorf("初始化 agent: %w", err)
	}
	agentSvc := agent.GetAgentService()

	h := handler.New(st, deckSvc, agentSvc, authSvc)
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
