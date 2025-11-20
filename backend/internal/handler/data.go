package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"twin-os/backend/internal/services"
)

// DataHandler 数据处理器
type DataHandler struct {
	service *services.DataService
}

// NewDataHandler 创建数据处理器
func NewDataHandler(service *services.DataService) *DataHandler {
	return &DataHandler{service: service}
}

// GetMessages 获取消息列表
func (h *DataHandler) GetMessages(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	messages, err := h.service.GetMessages(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get messages: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   messages,
	})
}

// GetContacts 获取联系人列表
func (h *DataHandler) GetContacts(c *gin.Context) {
	contacts, err := h.service.GetContacts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get contacts: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   contacts,
	})
}

// ExportData 导出数据
func (h *DataHandler) ExportData(c *gin.Context) {
	result, err := h.service.ExportData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to export data: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
	})
}