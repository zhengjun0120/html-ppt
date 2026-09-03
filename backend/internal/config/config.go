package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

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
	BaseURL string `yaml:"base_url"`
	APIKey  string `yaml:"api_key"`
	ModelID   string `yaml:"model_id"`
	MaxToken int64 `yaml:"max_token"`
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

	if v := os.Getenv("SERVER_ADDR"); v != "" {
		cfg.Server.Addr = v
	}
	if v := os.Getenv("LLM_API_KEY"); v != "" {
		cfg.LLM.APIKey = v
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
			ModelID:   "deepseek-chat",
		},
	}
}
