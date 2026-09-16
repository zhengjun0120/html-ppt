package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 是整个后端的配置总入口。
// 优先级：代码默认值 < config.yaml < 环境变量。
type Config struct {
	Server Server `yaml:"server"`
	Data   Data   `yaml:"data"`
	Assets Assets `yaml:"assets"`
	DB     DB     `yaml:"db"`
	LLM    LLM    `yaml:"llm"`
	Auth   Auth   `yaml:"auth"`
	SMTP   SMTP   `yaml:"smtp"`
	Redis  Redis  `yaml:"redis"`
	Crypto Crypto `yaml:"crypto"`
	// Features 功能开关。零值 = 全关（托管部署不配就是安全默认）；本地开发在
	// config.yaml 里显式打开。
	Features Features `yaml:"features"`
	Vision   Vision   `yaml:"vision"`
	Trace    Trace    `yaml:"trace"`
}

// Features 控制"能力型"功能的挂载。关掉某项 = agent 连对应工具都看不到，
// 而不是"看得到但一律拒绝"——后者只会让模型浪费轮次去试。
type Features struct {
	// CustomCSS 允许 agent 用 update_custom_css / read_custom_css 操作该 deck 的
	// 自定义样式槽（<style id="deck-custom">）。写入是"整体替换 + 清洗"，
	// 且整份 deck 的历史快照覆盖它——写坏了能回滚。
	CustomCSS bool `yaml:"custom_css"`

	// 联网搜索
	WebSearch bool `yaml:"web_search"`

	//视觉审查
	Vision bool `yaml:"vision"`

	// Trace 打开 agent 观测（每轮请求/工具参数与返回/工具内部过程/分项 token 落盘）。
	// 与别的开关不同，它**不影响 agent 的行为**，只影响"记录多少"——
	// 所以关掉它是纯粹的"不记录"，不会让模型少一个工具。
	// 记录量按 MB 计（全量上下文 + 逐页截图），托管部署建议关。
	Trace bool `yaml:"trace"`
}

type Server struct {
	Addr         string   `yaml:"addr"`
	AllowOrigins []string `yaml:"allow_origins"` // CORS 白名单
}

type Data struct {
	Dir string `yaml:"dir"` // 会话数据根目录
}

type Assets struct {
	Dir string `yaml:"dir"` // reveal.js 等静态资源目录
}

type DB struct {
	Enabled  bool   `yaml:"enabled"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type LLM struct {
	BaseURL  string `yaml:"base_url"`
	APIKey   string `yaml:"api_key"`
	ModelID  string `yaml:"model_id"`
	MaxToken int64  `yaml:"max_token"`

	AnthropicBaseURL string `yaml:"anthropic_base_url"`
}

// Auth 登录态配置。JWTSecret 是会话签名密钥：泄露 = 任何人可伪造登录态，
// 只放 config.yaml（gitignored）或环境变量，绝不入库存明文。
type Auth struct {
	JWTSecret string `yaml:"jwt_secret"`
	TokenTTL  string `yaml:"token_ttl"` // 如 "720h"（30 天）；解析失败回落默认值
}

// SMTP 邮件服务（QQ 邮箱示例）：password 填"授权码"而非邮箱登录密码，
// 在 QQ 邮箱 设置→账户→POP3/SMTP 服务 里生成。
type SMTP struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
}

type Redis struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
}

// Crypto AESKey 是 API Key 可逆加密的主密钥（base64 编码的 32 字节）。
// 与 JWTSecret 同级敏感：丢了它库里的密文全部作废，泄露 = 密钥全泄。
type Crypto struct {
	AESKey string `yaml:"aes_key"`
}

type Vision struct {
	ChromePath string `yaml:"chrome_path"`
}

// Trace 是 agent 的观测配置（把每轮请求、工具参数与返回、工具内部过程、
// 分项 token 落成 JSONL + 截图）。开关在 Features.Trace，这里只管"存到哪、存多少"——
// 与 features.vision / vision.chrome_path 是同一种分工。
type Trace struct {
	// Dir 落盘根目录。留空 = <data.dir>/traces（与 deck 同源，整目录备份/清理方便）。
	Dir string `yaml:"dir"`
	// CaptureImages 是否把视觉审查的每页截图存下来。
	// 关掉之后审查报告仍然完整，但"报告说的对不对"就无从核对了——那是这个工具
	// 最初要解决的问题之一，所以默认开。
	CaptureImages bool `yaml:"capture_images"`
	// RetainRunsPerSession 每个会话最多保留几个 run。0 = 不限。
	// 全量上下文是按 MB 计的，不裁剪迟早撑爆 data 目录。
	RetainRunsPerSession int `yaml:"retain_runs_per_session"`
	// MaxFieldBytes 单个长字段（工具参数/返回、上下文、报告）的字节上限。0 = 不截断。
	// 默认不截断是有意的：排查"模型为什么这么答"时，答案经常就在被截掉的那一段里。
	MaxFieldBytes int `yaml:"max_field_bytes"`
}

// 默认值。注意 Trace 只能在这里给"非零"的默认，因为 Load 用 yaml.Unmarshal
// 覆盖到已填好默认的结构体上：yaml 里没写的键保留默认，显式写了 0 就是 0。
const (
	// DefaultRetainRunsPerSession 见 Trace.RetainRunsPerSession
	DefaultRetainRunsPerSession = 20
	// DefaultCaptureImages 见 Trace.CaptureImages
	DefaultCaptureImages = true
)

// TTL 返回解析后的 token 有效期；配置缺失/写坏时回落 720h（30 天——
// 配合 auth.Refresh 的"每天一次静默续期"，正常用户基本不会再见到登录页）。
// 宽松降级：token 过期顶多要重新登录，不值得为此拒绝启动。
func (a Auth) TTL() time.Duration {
	d, err := time.ParseDuration(a.TokenTTL)
	if err != nil || d <= 0 {
		return 720 * time.Hour
	}
	return d
}

// Load 读取配置文件。配置文件不存在不是错误（用默认值跑），
// 但文件存在而格式错误是错误——配置写坏了应该立刻暴露，而不是静默用默认值。
func Load(path string) (*Config, error) {
	cfg := defaultConfig()

	resolved, err := resolvePath(path)
	if errors.Is(err, os.ErrNotExist) {
		return cfg, nil // 找不到配置文件 → 用默认值跑
	}
	if err != nil {
		return nil, fmt.Errorf("resolve %s: %w", path, err)
	}

	data, err := os.ReadFile(resolved)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", resolved, err)
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse %s: %w", resolved, err)
	}

	// 契约：配置里的相对路径一律相对 config.yaml 所在目录，与启动时的 cwd 解耦。
	// 否则从 backend/ 或 cmd/server/ 启动会得到两套 data 目录、assets 404。
	baseDir := filepath.Dir(resolved)
	if !filepath.IsAbs(cfg.Data.Dir) {
		cfg.Data.Dir = filepath.Join(baseDir, cfg.Data.Dir)
	}
	if !filepath.IsAbs(cfg.Assets.Dir) {
		cfg.Assets.Dir = filepath.Join(baseDir, cfg.Assets.Dir)
	}
	// 再取一次绝对路径：当 config.yaml 本身是用相对路径加载时，上面的 baseDir 是 "."，
	// Join 出来仍是相对路径——那"与 cwd 解耦"的契约就没真正成立（换个目录启动就找不到
	// 文件了：静态资源 404、read_component 读不到组件库）。这里补实这个契约。
	if abs, err := filepath.Abs(cfg.Data.Dir); err == nil {
		cfg.Data.Dir = abs
	}
	if abs, err := filepath.Abs(cfg.Assets.Dir); err == nil {
		cfg.Assets.Dir = abs
	}

	// trace.dir 必须在这里一起归一化。它和上面两个字段是同一个坑：
	// 只判断"是不是绝对路径"就 Join，config.yaml 用相对路径加载时 baseDir 是 "."，
	// 结果是"换个目录启动就写到别处去了"。而这个失败特别难发现——trace 会照常
	// 写进一个新目录、页面照常是空的，看起来像"观测没生效"而不像"路径配错了"。
	// 空值表示默认落在 data 目录下（此时 Data.Dir 已经是绝对路径）。
	cfg.Trace.Dir = strings.TrimSpace(cfg.Trace.Dir)
	if cfg.Trace.Dir == "" {
		cfg.Trace.Dir = filepath.Join(cfg.Data.Dir, "traces")
	} else {
		if !filepath.IsAbs(cfg.Trace.Dir) {
			cfg.Trace.Dir = filepath.Join(baseDir, cfg.Trace.Dir)
		}
		if abs, err := filepath.Abs(cfg.Trace.Dir); err == nil {
			cfg.Trace.Dir = abs
		}
	}

	if v := os.Getenv("SERVER_ADDR"); v != "" {
		cfg.Server.Addr = v
	}
	if v := os.Getenv("LLM_API_KEY"); v != "" {
		cfg.LLM.APIKey = v
	}
	// 联网搜索走的是 Anthropic 兼容入口，基址可以单独覆盖（自己前面挂网关时用）
	if v := os.Getenv("LLM_ANTHROPIC_BASE_URL"); v != "" {
		cfg.LLM.AnthropicBaseURL = v
	}
	if v := os.Getenv("AUTH_JWT_SECRET"); v != "" {
		cfg.Auth.JWTSecret = v
	}
	if v := os.Getenv("CRYPTO_AES_KEY"); v != "" {
		cfg.Crypto.AESKey = v
	}
	if v := os.Getenv("SMTP_PASSWORD"); v != "" {
		cfg.SMTP.Password = v
	}
	return cfg, nil
}

// resolvePath 解析配置文件路径：先用调用方给定的路径（相对 cwd），
// 找不到时从 cwd 向上逐级查找同名文件——和 git 向上找 .git 是同一个思路。
// 解决"从 backend/ 启动"和"从 cmd/server/ 启动"时相对路径对不上的问题。
func resolvePath(name string) (string, error) {
	if _, err := os.Stat(name); err == nil {
		return name, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for i := 0; i < 8; i++ { // 最多向上 8 级，防呆
		candidate := filepath.Join(dir, name)
		if st, err := os.Stat(candidate); err == nil && !st.IsDir() {
			return candidate, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break // 已到根目录
		}
		dir = parent
	}
	return "", os.ErrNotExist
}

func defaultConfig() *Config {
	return &Config{
		Server: Server{
			Addr:         ":8080",
			AllowOrigins: []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		},
		Data:   Data{Dir: "./data"},
		Assets: Assets{Dir: "./web/assets"},
		DB: DB{
			Enabled: true,
			Host:    "127.0.0.1",
			Port:    3306,
			User:    "root",
			Name:    "html_ppt",
		},
		LLM: LLM{
			BaseURL: "https://api.deepseek.com",
			ModelID: "deepseek-chat",
		},
		// Trace 的默认值只能在这里给（Load 是"先填默认、再用 yaml 覆盖"）：
		// 留到使用处判零值的话，用户显式写 retain_runs_per_session: 0（表示不限）
		// 会被当成"没配"而塞回默认值——一个永远不生效的配置项。
		Trace: Trace{
			CaptureImages:        DefaultCaptureImages,
			RetainRunsPerSession: DefaultRetainRunsPerSession,
		},
	}
}
