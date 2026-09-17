package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"html-ppt/backend/internal/agent"
	"html-ppt/backend/internal/export"
	"html-ppt/backend/internal/config"
	"html-ppt/backend/internal/cryptox"
	"html-ppt/backend/internal/handler"
	"html-ppt/backend/internal/router"
	"html-ppt/backend/internal/service/auth"
	"html-ppt/backend/internal/service/deck"
	"html-ppt/backend/internal/service/template"
	"html-ppt/backend/internal/store"
	"html-ppt/backend/internal/trace"
	"html-ppt/backend/internal/vision"
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

	deckSvc := deck.New(cfg.Data.Dir, cfg.Assets.Dir, st)

	// deck-v2 模板注册表：启动时一次加载+校验。失败不阻塞整个服务
	//（旧管线与登录等照常工作），但模板相关接口不可用、日志里有全貌——
	// 半成品模板让生成管线拿到坏契约的代价，远高于"模板库暂时不可用"。
	var templateReg *template.Registry
	if reg, err := template.NewRegistry(cfg.Templates.Dir, cfg.Assets.Dir); err != nil {
		log.Printf("[warn] 模板库加载有问题: %v", err)
		templateReg = reg // 部分模板可能仍注册成功
	} else {
		templateReg = reg
	}
	log.Printf("[info] 模板库: %d 个模板就绪（%s）", templateReg.Count(), cfg.Templates.Dir)
	deckSvc.WithTemplateRegistry(templateReg)
	// AI 味 lint 词表：config 可整体覆盖，缺省用代码内置词表
	deck.WithLintConfig(deck.LintCfg{
		CJKBanned: cfg.DeckV2.Lint.CJK_Banned,
		ENBanned:  cfg.DeckV2.Lint.EN_Banned,
		TitleMax:  cfg.DeckV2.Lint.TitleMaxChars,
		BulletMax: cfg.DeckV2.Lint.BulletMaxChars,
		EyebrowLimit: cfg.DeckV2.Lint.EyebrowPerPages,
	})

	// 视觉审查的一次性门票表：agent 发票（Issue）、handler 收票（Take），
	// 必须是**同一个实例**——各建一份的话，票发出去永远换不回来（而且不报错，只是 404）。
	visionGrants := &vision.Grants{}

	if err := agent.InitAgentModel(cfg.LLM, st, box, deckSvc, cfg.Features, cfg.Assets.Dir, trace.Config{
		Enabled:       cfg.Features.Trace,
		Dir:           cfg.Trace.Dir,
		CaptureImages: cfg.Trace.CaptureImages,
		RetainRuns:    cfg.Trace.RetainRunsPerSession,
		MaxFieldBytes: cfg.Trace.MaxFieldBytes,
	}); err != nil {
		return fmt.Errorf("初始化 agent: %w", err)
	}
	agentSvc := agent.GetAgentService()
	// deck-v2：模板注册表 + 阶段轮数预算（config 缺省走代码内默认）
	agentSvc.Templates = templateReg
	agentSvc.StageMaxTurns = map[string]int{
		"clarifying":         cfg.DeckV2.MaxTurns.Clarify,
		"outlining":          cfg.DeckV2.MaxTurns.Outline,
		"outline_review":     cfg.DeckV2.MaxTurns.OutlineReview,
		"generating":         cfg.DeckV2.MaxTurns.Generate,
		"iterating":          cfg.DeckV2.MaxTurns.Iterate,
	}

	// 审查的三个运行时依赖在 init 之后补：buildTools 只读开关（features.vision），
	// 而门票表 / 回环地址 / Chrome 路径只在"真的要审查那一刻"才被读到，
	// 所以这样接不必改 InitAgentModel 的签名（它已经有 6 个参数了）。
	var loopback string
	if cfg.Features.Vision || cfg.Features.Export {
		if base, err := loopbackBase(cfg.Server.Addr); err != nil {
			// 不返回错误：审查是增强信息，配错了不该让整个服务起不来（fail-open，
			// 与变量契约同一条原则——它防的是静默无效，不是安全问题）
			log.Printf("[warn] 视觉/导出：%v，相关功能不可用", err)
		} else {
			loopback = base
		}
	}
	if cfg.Features.Vision && loopback != "" {
		agentSvc.VisionGrants = visionGrants
		agentSvc.VisionBaseURL = loopback
		agentSvc.ChromePath = cfg.Vision.ChromePath
		log.Printf("[info] 视觉审查已开启：无头浏览器走 %s 取页", loopback)
	}

	// deck-v2 导出：复用视觉审查的同一张门票表与 Chrome；超时来自 config（默认 60s）
	var exportSvc *export.Service
	if cfg.Features.Export && loopback != "" {
		timeout := cfg.Export.TimeoutSeconds
		if timeout <= 0 {
			timeout = 60
		}
		exportSvc = export.New(deckSvc, loopback, cfg.Vision.ChromePath, time.Duration(timeout)*time.Second, visionGrants)
		log.Printf("[info] deck-v2 导出已开启（pdf/png/html，超时 %ds）", timeout)
	}

	h := handler.New(st, deckSvc, agentSvc, authSvc, visionGrants, trace.NewStore(cfg.Trace.Dir), templateReg, exportSvc)
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

// loopbackBase 由监听地址推出"后端自己怎么访问自己"。
//
// 为什么需要它：视觉审查要起一个无头浏览器去打开 deck 页，而那条一次性通道
// （/api/render/<nonce>）和 /assets/* 都挂在本服务上，所以必须告诉浏览器一个能用的地址。
// 监听写成 ":8080" 或 "0.0.0.0:8080" 时一律换成 127.0.0.1——0.0.0.0 不是一个可连接的目标。
func loopbackBase(addr string) (string, error) {
	_, port, err := net.SplitHostPort(addr)
	if err != nil || port == "" {
		return "", fmt.Errorf("从监听地址 %q 里解析不出端口", addr)
	}
	return "http://127.0.0.1:" + port, nil
}
