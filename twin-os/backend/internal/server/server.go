package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"twin-os/backend/internal/config"
	"twin-os/backend/internal/handler"
	"twin-os/backend/internal/repository"
	"twin-os/backend/internal/service"
	"twin-os/backend/pkg/database"
	"twin-os/backend/pkg/logger"
)

type Server struct {
	config       *config.Config
	router       *gin.Engine
	httpServer   *http.Server
	handlers     *handler.Handler
	services     *service.Service
	repositories *repository.Repository
}

func New(cfg *config.Config) *Server {
	// 初始化数据库
	db := database.Init(cfg.DatabasePath)

	// 初始化仓储层
	repos := repository.New(db)

	// 初始化服务层
	srvs := service.New(cfg, repos)

	// 初始化处理器层
	hdlrs := handler.New(srvs)

	// 初始化路由
	router := setupRouter(cfg, hdlrs)

	return &Server{
		config:       cfg,
		router:       router,
		handlers:     hdlrs,
		services:     srvs,
		repositories: repos,
	}
}

func (s *Server) Start() error {
	// 设置 Gin 模式
	if s.config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建 HTTP 服务器
	s.httpServer = &http.Server{
		Addr:    ":" + s.config.Port,
		Handler: s.router,
	}

	logger.Info("🌐 HTTP Server starting on port " + s.config.Port)

	// 启动服务器
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("❌ Server failed to start: " + err.Error())
		}
	}()

	logger.Info("✅ TwinOS Backend started successfully")

	// 等待中断信号来优雅地关闭服务器
	return s.gracefulShutdown()
}

func (s *Server) gracefulShutdown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("🛑 Shutting down server...")

	// 创建一个带超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭 HTTP 服务器
	if err := s.httpServer.Shutdown(ctx); err != nil {
		logger.Error("❌ Server forced to shutdown: " + err.Error())
		return err
	}

	// 关闭数据库连接
	if err := s.repositories.Close(); err != nil {
		logger.Error("❌ Database close error: " + err.Error())
		return err
	}

	logger.Info("✅ Server shutdown completed")
	return nil
}

func setupRouter(cfg *config.Config, h *handler.Handler) *gin.Engine {
	r := gin.New()

	// 中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"version":   "0.1.0",
		})
	})

	// API 路由组
	api := r.Group("/api/v1")
	{
		// 微信数据相关
		api.POST("/wechat/connect", h.WeChat.Connect)
		api.GET("/wechat/status", h.WeChat.Status)
		api.POST("/wechat/sync", h.WeChat.Sync)

		// 数据分析相关
		api.POST("/analysis/briefing", h.Analysis.GenerateBriefing)
		api.POST("/analysis/todos", h.Analysis.ExtractTodos)
		api.POST("/analysis/connections", h.Analysis.AnalyzeConnections)

		// 用户数据相关
		api.GET("/data/messages", h.Data.GetMessages)
		api.GET("/data/contacts", h.Data.GetContacts)
		api.POST("/data/export", h.Data.ExportData)

		// 设置相关
		api.GET("/settings", h.Settings.GetSettings)
		api.PUT("/settings", h.Settings.UpdateSettings)
	}

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}