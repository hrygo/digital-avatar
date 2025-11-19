package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// 服务器配置
	Port        string `json:"port"`
	Environment string `json:"environment"`
	LogLevel    string `json:"log_level"`

	// 数据库配置
	DatabasePath string `json:"database_path"`

	// 微信数据库配置
	WeChatDBPath string `json:"wechat_db_path"`

	// AI配置
	DeepSeekAPIKey string `json:"deepseek_api_key"`
	DeepSeekBaseURL string `json:"deepseek_base_url"`

	// 加密配置
	EncryptionKey string `json:"encryption_key"`

	// 安全配置
	MaxReadFrequency int `json:"max_read_frequency"`
}

func Load() *Config {
	// 加载.env文件
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	cfg := &Config{
		Port:        getEnv("PORT", "1234"),
		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),

		DatabasePath: getEnv("DATABASE_PATH", "./data/twin-os.db"),
		WeChatDBPath: getEnv("WECHAT_DB_PATH", ""),

		DeepSeekAPIKey:  getEnv("DEEPSEEK_API_KEY", ""),
		DeepSeekBaseURL: getEnv("DEEPSEEK_BASE_URL", "https://api.deepseek.com/v1"),

		EncryptionKey: getEnv("ENCRYPTION_KEY", "twin-os-default-key-change-in-production"),

		MaxReadFrequency: getEnvAsInt("MAX_READ_FREQUENCY", 5),
	}

	// 创建数据目录
	os.MkdirAll("./data", 0755)

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}