package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hrygo/log"
	"twin-os/backend/internal/services"
)

// WeChatHandlerV2 微信数据处理器 v0.4.0
type WeChatHandlerV2 struct {
	wechatService *services.WeChatServiceV2
}

// NewWeChatHandlerV2 创建微信数据处理器
func NewWeChatHandlerV2(wechatService *services.WeChatServiceV2) *WeChatHandlerV2 {
	return &WeChatHandlerV2{
		wechatService: wechatService,
	}
}

// ImportRequest 导入请求
type ImportRequest struct {
	Type string `json:"type" binding:"required,oneof=database backup file"` // 导入类型
	Path string `json:"path" binding:"required"`                           // 数据源路径
	Mode string `json:"mode" binding:"oneof=full incremental"`             // 同步模式
}

// ImportResponse 导入响应
type ImportResponse struct {
	SyncID int64  `json:"sync_id"`
	Status string `json:"status"`
}

// Connect 连接微信数据源
func (h *WeChatHandlerV2) Connect(c *gin.Context) {
	var req ImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	log.Infof("收到微信数据导入请求: type=%s, path=%s, mode=%s", req.Type, req.Path, req.Mode)

	if req.Type != "database" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "不支持的导入类型: " + req.Type,
		})
		return
	}

	// 默认使用全量同步
	if req.Mode == "" {
		req.Mode = "full"
	}

	syncRecord, err := h.wechatService.Sync(req.Path, req.Mode)
	if err != nil {
		log.Errorf("微信数据导入失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "导入失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": ImportResponse{
			SyncID: syncRecord.ID,
			Status: syncRecord.Status,
		},
	})
}

// Status 获取微信数据状态
func (h *WeChatHandlerV2) Status(c *gin.Context) {
	status, err := h.wechatService.GetWeChatStatus()
	if err != nil {
		log.Errorf("获取微信状态失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取状态失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": status})
}

// Sync 同步微信数据
func (h *WeChatHandlerV2) Sync(c *gin.Context) {
	var req ImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	log.Infof("收到微信数据同步请求: type=%s, path=%s, mode=%s", req.Type, req.Path, req.Mode)

	if req.Type != "database" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "不支持的同步类型: " + req.Type,
		})
		return
	}

	// 默认使用增量同步
	if req.Mode == "" {
		req.Mode = "incremental"
	}

	syncRecord, err := h.wechatService.Sync(req.Path, req.Mode)
	if err != nil {
		log.Errorf("微信数据同步失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "同步失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": ImportResponse{
			SyncID: syncRecord.ID,
			Status: syncRecord.Status,
		},
	})
}

// GetMessages 获取消息列表
func (h *WeChatHandlerV2) GetMessages(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	chatType := c.Query("chat_type")
	talker := c.Query("talker")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	messages, total, err := h.wechatService.GetMessages(page, pageSize, chatType, talker)
	if err != nil {
		log.Errorf("获取消息列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取消息失败: " + err.Error(),
		})
		return
	}

	response := gin.H{
		"data":       messages,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// GetContacts 获取联系人列表
func (h *WeChatHandlerV2) GetContacts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	chatType := c.Query("chat_type")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	contacts, total, err := h.wechatService.GetContacts(page, pageSize, chatType)
	if err != nil {
		log.Errorf("获取联系人列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取联系人失败: " + err.Error(),
		})
		return
	}

	response := gin.H{
		"data":       contacts,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// GetChats 获取聊天会话列表
func (h *WeChatHandlerV2) GetChats(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	chatType := c.Query("chat_type")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	chats, total, err := h.wechatService.GetChats(page, pageSize, chatType)
	if err != nil {
		log.Errorf("获取聊天会话列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取聊天会话失败: " + err.Error(),
		})
		return
	}

	response := gin.H{
		"data":       chats,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// GetSyncRecords 获取同步记录
func (h *WeChatHandlerV2) GetSyncRecords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	records, total, err := h.wechatService.GetSyncRecords(page, pageSize)
	if err != nil {
		log.Errorf("获取同步记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取同步记录失败: " + err.Error(),
		})
		return
	}

	response := gin.H{
		"data":       records,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// SearchMessages 搜索消息
func (h *WeChatHandlerV2) SearchMessages(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "搜索关键词不能为空",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	messages, total, err := h.wechatService.SearchMessages(keyword, page, pageSize)
	if err != nil {
		log.Errorf("搜索消息失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "搜索失败: " + err.Error(),
		})
		return
	}

	response := gin.H{
		"data":       messages,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		"keyword":    keyword,
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// GetChatMessages 获取聊天消息
func (h *WeChatHandlerV2) GetChatMessages(c *gin.Context) {
	chatID := c.Param("chat_id")
	if chatID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "聊天ID不能为空",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	messages, total, err := h.wechatService.GetChatMessages(chatID, page, pageSize)
	if err != nil {
		log.Errorf("获取聊天消息失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取聊天消息失败: " + err.Error(),
		})
		return
	}

	response := gin.H{
		"data":       messages,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		"chat_id":    chatID,
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// GetStatistics 获取统计信息
func (h *WeChatHandlerV2) GetStatistics(c *gin.Context) {
	stats, err := h.wechatService.GetStatistics()
	if err != nil {
		log.Errorf("获取统计信息失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取统计信息失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}