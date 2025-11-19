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
	"twin-os/backend/pkg/cache"
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
	cacheManager *cache.CacheManager
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

	// 初始化缓存管理器
	cacheMgr := cache.NewCacheManager()

	// 初始化路由
	router := setupRouter(cfg, hdlrs, cacheMgr)

	return &Server{
		config:        cfg,
		router:        router,
		handlers:      hdlrs,
		services:      srvs,
		repositories:  repos,
		cacheManager:  cacheMgr,
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

func setupRouter(cfg *config.Config, h *handler.Handler, cacheMgr *cache.CacheManager) *gin.Engine {
	r := gin.New()

	// 中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())

	// 获取API缓存实例
	apiCache := cacheMgr.GetCache("api")

	// 创建缓存配置
	shortCacheConfig := cache.DefaultCacheConfig()
	shortCacheConfig.TTL = 2 * time.Minute // 短期缓存2分钟

	longCacheConfig := cache.DefaultCacheConfig()
	longCacheConfig.TTL = 10 * time.Minute // 长期缓存10分钟

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

		// 缓存状态查询 - 短期缓存
		api.GET("/wechat/status", cache.CacheMiddleware(apiCache, shortCacheConfig), h.WeChat.Status)
		api.POST("/wechat/sync", h.WeChat.Sync)

		// 数据分析相关 (不缓存 - 这些是计算密集型操作，结果应该实时)
		api.POST("/analysis/briefing", h.Analysis.GenerateBriefing)
		api.POST("/analysis/todos", h.Analysis.ExtractTodos)
		api.POST("/analysis/connections", h.Analysis.AnalyzeConnections)

		// 增强AI分析相关 (不缓存 - 实时AI计算)
		api.POST("/analysis/message", h.EnhancedAnalysis.AnalyzeMessage)
		api.GET("/analysis/conversation", h.EnhancedAnalysis.AnalyzeConversation)
		api.GET("/analysis/emotion/trend", h.EnhancedAnalysis.GetEmotionTrend)
		api.GET("/analysis/topics", h.EnhancedAnalysis.GetTopicAnalysis)
		api.GET("/analysis/intents", h.EnhancedAnalysis.GetIntentDistribution)
		api.GET("/analysis/summary", h.EnhancedAnalysis.GetConversationSummary)
		api.POST("/analysis/batch", h.EnhancedAnalysis.BatchAnalyzeMessages)
		api.GET("/analysis/status", h.EnhancedAnalysis.GetAnalysisStatus)
		api.PUT("/analysis/config", h.EnhancedAnalysis.ConfigureAnalysis)

		// 用户数据相关 - 缓存数据查询
		api.GET("/data/messages", cache.CacheMiddleware(apiCache, shortCacheConfig), h.Data.GetMessages)
		api.GET("/data/contacts", cache.CacheMiddleware(apiCache, longCacheConfig), h.Data.GetContacts)
		api.POST("/data/export", h.Data.ExportData)

		// 设置相关 - 缓存设置查询
		api.GET("/settings", cache.CacheMiddleware(apiCache, longCacheConfig), h.Settings.GetSettings)
		api.PUT("/settings", h.Settings.UpdateSettings)

		// 备份管理路由
		api.POST("/backup/create", h.Backup.CreateBackup)
		api.GET("/backup/list", h.Backup.ListBackups)
		api.GET("/backup/stats", h.Backup.GetBackupStats)
		api.POST("/backup/schedule", h.Backup.ScheduleBackup)
		api.POST("/backup/import", h.Backup.ImportBackup)
		api.GET("/backup/:id", h.Backup.GetBackup)
		api.POST("/backup/:id/restore", h.Backup.RestoreBackup)
		api.POST("/backup/:id/export", h.Backup.ExportBackup)
		api.GET("/backup/:id/verify", h.Backup.VerifyBackup)
		api.DELETE("/backup/:id", h.Backup.DeleteBackup)

		// 缓存管理路由
		api.GET("/cache/stats", func(c *gin.Context) {
			stats := cacheMgr.GetAllStats()
			c.JSON(200, gin.H{
				"status": "success",
				"data":   stats,
				"timestamp": time.Now().Unix(),
			})
		})
		api.POST("/cache/clear", func(c *gin.Context) {
			cacheMgr.ClearAll()
			c.JSON(200, gin.H{
				"status": "success",
				"message": "All caches cleared",
			})
		})

		// 错误报告路由
		api.POST("/error/report", func(c *gin.Context) {
			// 处理前端错误报告
			var errorReport struct {
				Message       string `json:"message"`
				Stack         string `json:"stack"`
				ComponentStack string `json:"componentStack"`
				Timestamp     string `json:"timestamp"`
				UserAgent     string `json:"userAgent"`
				URL           string `json:"url"`
			}

			if err := c.ShouldBindJSON(&errorReport); err != nil {
				c.JSON(400, gin.H{"error": "Invalid error report format"})
				return
			}

			// 记录错误到日志
			logger.Error("🚨 Frontend Error Report:",
				"message", errorReport.Message,
				"stack", errorReport.Stack,
				"url", errorReport.URL,
				"timestamp", errorReport.Timestamp)

			c.JSON(200, gin.H{
				"status": "success",
				"message": "Error report received",
			})
		})
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