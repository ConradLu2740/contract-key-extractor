package config

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	AIService AIServiceConfig `yaml:"ai_service"`
	Upload    UploadConfig    `yaml:"upload"`
	Output    OutputConfig    `yaml:"output"`
	LLM       LLMConfig       `yaml:"llm"`
	Logging   LoggingConfig   `yaml:"logging"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type AIServiceConfig struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	Timeout int    `yaml:"timeout"`
}

type UploadConfig struct {
	Path    string `yaml:"path"`
	MaxSize int64  `yaml:"max_size"`
}

type OutputConfig struct {
	Path string `yaml:"path"`
}

type LLMConfig struct {
	Provider string `yaml:"provider"`
	APIKey   string `yaml:"api_key"`
	Model    string `yaml:"model"`
	OCRModel string `yaml:"ocr_model"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

var globalConfig *Config

func Load(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = "./configs/config.yaml"
	}

	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	cfg.expandEnvVars()
	globalConfig = &cfg

	return &cfg, nil
}

func (c *Config) expandEnvVars() {
	if c.LLM.APIKey != "" && len(c.LLM.APIKey) > 3 && c.LLM.APIKey[0] == '$' {
		envName := c.LLM.APIKey[2 : len(c.LLM.APIKey)-1]
		c.LLM.APIKey = os.Getenv(envName)
	}
}

func Get() *Config {
	return globalConfig
}

func InitLogger(cfg *LoggingConfig) (*zap.Logger, error) {
	var config zap.Config
	if cfg.Format == "json" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	switch cfg.Level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	return config.Build()
}
