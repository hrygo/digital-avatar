package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"twin-os/backend/internal/services"
)

// AnalysisHandler 分析处理器
type AnalysisHandler struct {
	service *services.AnalysisService
}

// NewAnalysisHandler 创建分析处理器
func NewAnalysisHandler(service *services.AnalysisService) *AnalysisHandler {
	return &AnalysisHandler{service: service}
}

// GenerateBriefing 生成情报简报
func (h *AnalysisHandler) GenerateBriefing(c *gin.Context) {
	result, err := h.service.GenerateBriefing()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate briefing: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
	})
}

// ExtractTodos 提取待办事项
func (h *AnalysisHandler) ExtractTodos(c *gin.Context) {
	result, err := h.service.ExtractTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to extract todos: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
	})
}

// AnalyzeConnections 分析人脉关系
func (h *AnalysisHandler) AnalyzeConnections(c *gin.Context) {
	result, err := h.service.AnalyzeConnections()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to analyze connections: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
	})
}