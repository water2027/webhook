package config

import (
	"github.com/joho/godotenv"
	"os"
)

func Get(key string) string {
	return os.Getenv(key)
}

func GetWithDefault(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}

func InitConfig() {
	_ = godotenv.Load()
}
