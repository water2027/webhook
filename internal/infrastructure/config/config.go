package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	DatabaseURL string
	Port        string
	FeishuAppID     string
	FeishuAppSecret string
	FeishuOpenID    string
	LogLevel        string
	LogFormat       string
}

var GlobalConfig *Config

func Load() {
	_ = godotenv.Load()

	GlobalConfig = &Config{
		DatabaseURL: getWithDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/webhook?sslmode=disable"),
		Port:        getWithDefault("PORT", "8080"),
		FeishuAppID:     os.Getenv("FEISHU_APP_ID"),
		FeishuAppSecret: os.Getenv("FEISHU_APP_SECRET"),
		FeishuOpenID:    os.Getenv("FEISHU_OPEN_ID"),
		LogLevel:        getWithDefault("LOG_LEVEL", "info"),
		LogFormat:       getWithDefault("LOG_FORMAT", "json"),
	}
}

func getWithDefault(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}
