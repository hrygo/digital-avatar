package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hrygo/log"
	"twin-os/backend/internal/services"
)

// RealtimeSyncHandler 实时同步处理器
type RealtimeSyncHandler struct {
	syncService *services.RealtimeSyncService
}

// NewRealtimeSyncHandler 创建实时同步处理器
func NewRealtimeSyncHandler(syncService *services.RealtimeSyncService) *RealtimeSyncHandler {
	return &RealtimeSyncHandler{
		syncService: syncService,
	}
}

// StartSyncRequest 启动同步请求
type StartSyncRequest struct {
	WeChatDBPath   string        `json:"wechat_db_path" binding:"required"` // 微信数据库路径
	AutoDetect     bool          `json:"auto_detect"`                       // 自动检测路径
	SyncInterval   int           `json:"sync_interval"`                     // 同步间隔(秒)
	EnableMonitoring bool        `json:"enable_monitoring"`                // 启用文件监控
}

// ConfigRequest 配置更新请求
type ConfigRequest struct {
	SyncInterval     int      `json:"sync_interval"`      // 同步间隔(秒)
	MaxSyncRetries   int      `json:"max_sync_retries"`   // 最大重试次数
	SyncBatchSize    int      `json:"sync_batch_size"`    // 批量同步大小
	EnableMonitoring bool     `json:"enable_monitoring"`  // 启用文件监控
	MinMessageAge    int      `json:"min_message_age"`    // 最小消息年龄(秒)
	ChatFilters      []string `json:"chat_filters"`       // 聊天过滤器
	WorkerCount      int      `json:"worker_count"`      // 工作协程数量
	BufferSize       int      `json:"buffer_size"`       // 缓冲区大小
}

// StartSync 启动实时同步
func (h *RealtimeSyncHandler) StartSync(c *gin.Context) {
	var req StartSyncRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	log.Infof("收到启动实时同步请求: db_path=%s, auto_detect=%t, interval=%ds",
		req.WeChatDBPath, req.AutoDetect, req.SyncInterval)

	// 创建同步配置
	config := &services.SyncConfig{
		WeChatDBPath:     req.WeChatDBPath,
		AutoDetectPath:   req.AutoDetect,
		SyncInterval:     time.Duration(req.SyncInterval) * time.Second,
		EnableMonitoring: req.EnableMonitoring,
	}

	// 设置配置
	h.syncService.SetConfig(config)

	// 启动同步服务
	if err := h.syncService.Start(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "启动同步服务失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "实时同步服务已启动",
	})
}

// StopSync 停止实时同步
func (h *RealtimeSyncHandler) StopSync(c *gin.Context) {
	log.Info("收到停止实时同步请求")

	if err := h.syncService.Stop(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "停止同步服务失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "实时同步服务已停止",
	})
}

// GetSyncStatus 获取同步状态
func (h *RealtimeSyncHandler) GetSyncStatus(c *gin.Context) {
	status := h.syncService.GetStatus()

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   status,
	})
}

// UpdateConfig 更新同步配置
func (h *RealtimeSyncHandler) UpdateConfig(c *gin.Context) {
	var req ConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	log.Infof("收到更新同步配置请求: interval=%ds, monitoring=%t",
		req.SyncInterval, req.EnableMonitoring)

	// 获取当前配置
	config := h.syncService.GetConfig()

	// 更新配置
	if req.SyncInterval > 0 {
		config.SyncInterval = time.Duration(req.SyncInterval) * time.Second
	}
	if req.MaxSyncRetries > 0 {
		config.MaxSyncRetries = req.MaxSyncRetries
	}
	if req.SyncBatchSize > 0 {
		config.SyncBatchSize = req.SyncBatchSize
	}
	config.EnableMonitoring = req.EnableMonitoring
	if req.MinMessageAge > 0 {
		config.MinMessageAge = time.Duration(req.MinMessageAge) * time.Second
	}
	if len(req.ChatFilters) > 0 {
		config.ChatFilters = req.ChatFilters
	}
	if req.WorkerCount > 0 {
		config.WorkerCount = req.WorkerCount
	}
	if req.BufferSize > 0 {
		config.BufferSize = req.BufferSize
	}

	// 应用新配置
	h.syncService.SetConfig(config)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "同步配置已更新",
		"data":    config,
	})
}

// GetConfig 获取当前配置
func (h *RealtimeSyncHandler) GetConfig(c *gin.Context) {
	config := h.syncService.GetConfig()

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   config,
	})
}

// TriggerManualSync 手动触发同步
func (h *RealtimeSyncHandler) TriggerManualSync(c *gin.Context) {
	log.Info("收到手动触发同步请求")

	if err := h.syncService.TriggerManualSync(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "手动同步失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "手动同步已触发",
	})
}

// DetectWeChatDB 检测微信数据库路径
func (h *RealtimeSyncHandler) DetectWeChatDB(c *gin.Context) {
	paths := []string{}

	// 手动执行检测逻辑
	homeDir, _ := os.UserHomeDir()
	searchPaths := []string{
		filepath.Join(homeDir, "Library/Containers/com.tencent.xin/Data/Library/Application Support/com.tencent.xin/*/Message/MessageTemp/*/DB/EnMicroMsg.db"),
		filepath.Join(homeDir, "Documents/WeChat Files/*/Msg/*/EnMicroMsg.db"),
		filepath.Join(homeDir, "Library/Application Support/Tencent/MicroMsg/*/DB/MicroMsg.db"),
	}

	for _, pattern := range searchPaths {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}
		paths = append(paths, matches...)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"detected_paths": paths,
			"count":          len(paths),
		},
	})
}

// GetSyncHistory 获取同步历史
func (h *RealtimeSyncHandler) GetSyncHistory(c *gin.Context) {
	page := 1
	pageSize := 20

	// 解析分页参数
	if p, exists := c.GetQuery("page"); exists {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if p, exists := c.GetQuery("page_size"); exists {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	// 这里需要访问数据库获取同步历史
	// 由于这个handler没有直接访问数据库，我们需要通过service获取
	// 暂时返回示例数据
	history := []gin.H{
		{
			"id":              1,
			"sync_type":       "incremental",
			"status":          "completed",
			"total_messages":  150,
			"new_messages":    5,
			"new_contacts":    1,
			"start_time":      time.Now().Add(-5 * time.Minute).Unix(),
			"duration":        2.5,
			"error_message":   "",
		},
		{
			"id":              2,
			"sync_type":       "incremental",
			"status":          "completed",
			"total_messages":  145,
			"new_messages":    3,
			"new_contacts":    0,
			"start_time":      time.Now().Add(-10 * time.Minute).Unix(),
			"duration":        1.8,
			"error_message":   "",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"history":      history,
			"page":         page,
			"page_size":    pageSize,
			"total":        len(history),
			"total_pages":  1,
		},
	})
}