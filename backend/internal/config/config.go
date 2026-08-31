package config

import (
	"errors"
	"fmt"
	"os"

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

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return cfg, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}

	if v := os.Getenv("SERVER_ADDR"); v != "" {
		cfg.Server.Addr = v
	}
	if v := os.Getenv("LLM_API_KEY"); v != "" {
		cfg.LLM.APIKey = v
	}
	return cfg, nil
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
