package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"twin-os/backend/internal/service"
)

// WeChatHandler 微信处理器
type WeChatHandler struct {
	service *service.WeChatService
}

// NewWeChatHandler 创建微信处理器
func NewWeChatHandler(service *service.WeChatService) *WeChatHandler {
	return &WeChatHandler{service: service}
}

// Connect 连接微信数据库
func (h *WeChatHandler) Connect(c *gin.Context) {
	var req struct {
		DBPath string `json:"db_path" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	if err := h.service.Connect(req.DBPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to connect: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "WeChat database connected successfully",
	})
}

// Status 获取连接状态
func (h *WeChatHandler) Status(c *gin.Context) {
	status, err := h.service.GetStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get status: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   status,
	})
}

// Sync 同步数据
func (h *WeChatHandler) Sync(c *gin.Context) {
	result, err := h.service.SyncMessages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to sync: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
	})
}