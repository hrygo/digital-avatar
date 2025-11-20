package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"twin-os/backend/internal/services"
)

// EnhancedAnalysisHandler 增强分析处理器
type EnhancedAnalysisHandler struct {
	service *services.EnhancedAnalysisService
}

// NewEnhancedAnalysisHandler 创建增强分析处理器
func NewEnhancedAnalysisHandler(service *services.EnhancedAnalysisService) *EnhancedAnalysisHandler {
	return &EnhancedAnalysisHandler{service: service}
}

// AnalyzeMessage 分析单个消息
// @Summary 分析消息情感、意图、主题等
// @Description 对单个消息进行全面的AI分析
// @Tags AI Analysis
// @Accept json
// @Produce json
// @Param request body services.MessageAnalysisRequest true "消息分析请求"
// @Success 200 {object} services.MessageAnalysisResponse "分析结果"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/v1/analysis/message [post]
func (h *EnhancedAnalysisHandler) AnalyzeMessage(c *gin.Context) {
	var req services.MessageAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format: " + err.Error(),
		})
		return
	}

	result, err := h.service.AnalyzeMessage(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to analyze message: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
	})
}

// AnalyzeConversation 分析整个对话
// @Summary 分析对话历史
// @Description 对整个对话历史进行综合分析
// @Tags AI Analysis
// @Accept json
// @Produce json
// @Param userId query string false "用户ID"
// @Param timeRange query string false "时间范围 (1h, 24h, 7d)"
// @Param limit query int false "消息数量限制" default(100)
// @Success 200 {object} services.ConversationAnalysisResponse "分析结果"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/v1/analysis/conversation [get]
func (h *EnhancedAnalysisHandler) AnalyzeConversation(c *gin.Context) {
	userID := c.Query("userId")
	timeRange := c.Query("timeRange")
	limitStr := c.DefaultQuery("limit", "100")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 100
	}

	result, err := h.service.AnalyzeConversation(userID, timeRange, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to analyze conversation: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
	})
}

// GetEmotionTrend 获取情感趋势
// @Summary 获取情感变化趋势
// @Description 分析指定时间范围内的情感变化趋势
// @Tags AI Analysis
// @Accept json
// @Produce json
// @Param userId query string false "用户ID"
// @Param timeRange query string false "时间范围 (1h, 24h, 7d, 30d)" default(24h)
// @Success 200 {object} services.EmotionTrendResponse "情感趋势结果"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/v1/analysis/emotion/trend [get]
func (h *EnhancedAnalysisHandler) GetEmotionTrend(c *gin.Context) {
	userID := c.Query("userId")
	timeRange := c.DefaultQuery("timeRange", "24h")

	result, err := h.service.GetEmotionTrend(userID, timeRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get emotion trend: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
	})
}

// GetTopicAnalysis 获取主题分析
// @Summary 获取对话主题分析
// @Description 分析对话中的主要主题和分布
// @Tags AI Analysis
// @Accept json
// @Produce json
// @Param userId query string false "用户ID"
// @Param timeRange query string false "时间范围 (1h, 24h, 7d, 30d)" default(24h)
// @Success 200 {object} services.TopicAnalysisResponse "主题分析结果"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/v1/analysis/topics [get]
func (h *EnhancedAnalysisHandler) GetTopicAnalysis(c *gin.Context) {
	userID := c.Query("userId")
	timeRange := c.DefaultQuery("timeRange", "24h")

	result, err := h.service.GetTopicAnalysis(userID, timeRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get topic analysis: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
	})
}

// GetIntentDistribution 获取意图分布
// @Summary 获取意图分布统计
// @Description 统计对话中各种意图类型的分布情况
// @Tags AI Analysis
// @Accept json
// @Produce json
// @Param userId query string false "用户ID"
// @Param timeRange query string false "时间范围 (1h, 24h, 7d, 30d)" default(24h)
// @Success 200 {object} services.IntentDistributionResponse "意图分布结果"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/v1/analysis/intents [get]
func (h *EnhancedAnalysisHandler) GetIntentDistribution(c *gin.Context) {
	userID := c.Query("userId")
	timeRange := c.DefaultQuery("timeRange", "24h")

	result, err := h.service.GetIntentDistribution(userID, timeRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get intent distribution: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
	})
}

// GetConversationSummary 获取对话摘要
// @Summary 获取对话摘要报告
// @Description 生成对话的全面摘要报告
// @Tags AI Analysis
// @Accept json
// @Produce json
// @Param userId query string false "用户ID"
// @Param timeRange query string false "时间范围 (1h, 24h, 7d, 30d)" default(24h)
// @Success 200 {object} services.ConversationSummaryResponse "对话摘要结果"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/v1/analysis/summary [get]
func (h *EnhancedAnalysisHandler) GetConversationSummary(c *gin.Context) {
	userID := c.Query("userId")
	timeRange := c.DefaultQuery("timeRange", "24h")

	result, err := h.service.GetConversationSummary(userID, timeRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get conversation summary: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
	})
}

// BatchAnalyzeMessages 批量分析消息
// @Summary 批量分析多条消息
// @Description 对多条消息进行批量AI分析
// @Tags AI Analysis
// @Accept json
// @Produce json
// @Param request body services.BatchAnalysisRequest true "批量分析请求"
// @Success 200 {object} services.BatchAnalysisResponse "批量分析结果"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/v1/analysis/batch [post]
func (h *EnhancedAnalysisHandler) BatchAnalyzeMessages(c *gin.Context) {
	var req services.BatchAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format: " + err.Error(),
		})
		return
	}

	result, err := h.service.BatchAnalyzeMessages(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to batch analyze messages: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
	})
}

// GetAnalysisStatus 获取分析状态
// @Summary 获取AI分析系统状态
// @Description 检查AI分析模块的运行状态和配置
// @Tags AI Analysis
// @Accept json
// @Produce json
// @Success 200 {object} services.AnalysisStatusResponse "分析状态结果"
// @Router /api/v1/analysis/status [get]
func (h *EnhancedAnalysisHandler) GetAnalysisStatus(c *gin.Context) {
	result, err := h.service.GetAnalysisStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get analysis status: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
	})
}

// ConfigureAnalysis 配置分析参数
// @Summary 配置AI分析参数
// @Description 动态配置AI分析的各种参数
// @Tags AI Analysis
// @Accept json
// @Produce json
// @Param request body services.AnalysisConfigRequest true "分析配置请求"
// @Success 200 {object} map[string]string "配置成功"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/v1/analysis/config [put]
func (h *EnhancedAnalysisHandler) ConfigureAnalysis(c *gin.Context) {
	var req services.AnalysisConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format: " + err.Error(),
		})
		return
	}

	err := h.service.ConfigureAnalysis(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to configure analysis: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Analysis configuration updated successfully",
	})
}