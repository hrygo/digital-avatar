package main

import (
	"log"
	"os"

	"twin-os/backend/internal/config"
	"twin-os/backend/internal/server"
	"twin-os/backend/pkg/logger"
)

func main() {
	// 初始化配置
	cfg := config.Load()

	// 初始化日志
	logger.Init(cfg.LogLevel)
	log.Println("🚀 TwinOS Backend Starting...")
	log.Printf("📊 Environment: %s", cfg.Environment)

	// 启动服务器
	srv := server.New(cfg)
	if err := srv.Start(); err != nil {
		log.Printf("❌ Server failed to start: %v", err)
		os.Exit(1)
	}
}