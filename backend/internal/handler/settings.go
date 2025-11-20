package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"twin-os/backend/internal/services"
)

// SettingsHandler 设置处理器
type SettingsHandler struct {
	service *services.SettingsService
}

// NewSettingsHandler 创建设置处理器
func NewSettingsHandler(service *services.SettingsService) *SettingsHandler {
	return &SettingsHandler{service: service}
}

// GetSettings 获取设置
func (h *SettingsHandler) GetSettings(c *gin.Context) {
	settings, err := h.service.GetSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get settings: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   settings,
	})
}

// UpdateSettings 更新设置
func (h *SettingsHandler) UpdateSettings(c *gin.Context) {
	var settings map[string]string

	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	result, err := h.service.UpdateSettings(settings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update settings: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}