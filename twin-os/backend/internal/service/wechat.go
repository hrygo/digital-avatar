package service

import (
	"fmt"
	"time"

	"twin-os/backend/internal/config"
	"twin-os/backend/internal/model"
	"twin-os/backend/internal/repository"
	"twin-os/backend/pkg/crypto"
	"twin-os/backend/pkg/logger"
	"twin-os/backend/pkg/wechat"
)

// WeChatService 微信服务
type WeChatService struct {
	config     *config.Config
	repository *repository.Repository
	reader     *wechat.Reader
	detector   *crypto.PIIDetector
}

// NewWeChatService 创建微信服务
func NewWeChatService(cfg *config.Config, repo *repository.Repository) *WeChatService {
	return &WeChatService{
		config:     cfg,
		repository: repo,
		detector:   crypto.NewPIIDetector(),
	}
}

// Connect 连接微信数据库
func (s *WeChatService) Connect(dbPath string) error {
	if dbPath == "" {
		// 尝试自动查找微信数据库
		foundPath, err := wechat.FindWeChatDB()
		if err != nil {
			return fmt.Errorf("failed to find wechat database: %w", err)
		}
		dbPath = foundPath
	}

	// 验证数据库有效性
	if !wechat.IsValidWeChatDB(dbPath) {
		return fmt.Errorf("invalid wechat database file: %s", dbPath)
	}

	// 创建读取器
	s.reader = wechat.NewReader(dbPath)

	// 打开数据库
	if err := s.reader.Open(); err != nil {
		return fmt.Errorf("failed to open wechat database: %w", err)
	}

	// 保存配置到数据库
	if err := s.repository.Setting.Set("wechat_db_path", dbPath, "微信数据库路径"); err != nil {
		logger.Warn("⚠️ Failed to save wechat db path: " + err.Error())
	}

	logger.Info("✅ WeChat database connected: " + dbPath)
	return nil
}

// Disconnect 断开微信数据库连接
func (s *WeChatService) Disconnect() error {
	if s.reader != nil {
		if err := s.reader.Close(); err != nil {
			return fmt.Errorf("failed to close wechat database: %w", err)
		}
		s.reader = nil
	}

	logger.Info("✅ WeChat database disconnected")
	return nil
}

// GetStatus 获取连接状态
func (s *WeChatService) GetStatus() (*WeChatStatus, error) {
	status := &WeChatStatus{
		IsConnected: s.reader != nil && s.reader.IsOpened(),
		LastSync:    time.Now(),
	}

	if status.IsConnected && s.reader != nil {
		// 获取一些基本信息
		if messages, err := s.reader.GetMessages(1, 0); err == nil && len(messages) > 0 {
			status.MessageCount = &messages[0].ID // 简化处理
		}

		if contacts, err := s.reader.GetContacts(); err == nil {
			status.ContactCount = len(contacts)
		}

		// 获取配置中的数据库路径
		if dbPath, err := s.repository.Setting.Get("wechat_db_path"); err == nil {
			status.DBPath = dbPath
		}
	}

	return status, nil
}

// SyncMessages 同步消息
func (s *WeChatService) SyncMessages() (*SyncResult, error) {
	if s.reader == nil || !s.reader.IsOpened() {
		return nil, fmt.Errorf("wechat database not connected")
	}

	result := &SyncResult{
		StartTime: time.Now(),
	}

	// 获取最后同步时间 (临时禁用)
	// lastSyncTimeStr, err := s.repository.Setting.Get("last_sync_time")
	// var lastSyncTime time.Time
	// if err == nil && lastSyncTimeStr != "" {
	// 	if timestamp, parseErr := time.Parse(time.RFC3339, lastSyncTimeStr); parseErr == nil {
	// 		lastSyncTime = timestamp
	// 	}
	// }

	// 临时调试：获取所有消息
	messages, err := s.reader.GetMessages(1000, 0)

	// 读取微信消息（最近24小时）
	// now := time.Now()
	// startTime := lastSyncTime
	// if startTime.IsZero() {
	// 	startTime = now.Add(-24 * time.Hour)
	// }

	// messages, err := s.reader.GetMessagesByTimeRange(startTime, now)
	if err != nil {
		return nil, fmt.Errorf("failed to read messages: %w", err)
	}

	logger.Info(fmt.Sprintf("📊 Found %d messages to sync", len(messages)))

	// 处理每条消息
	for _, msg := range messages {
		// 清理消息内容
		msg.Content = wechat.CleanMessageContent(msg.Content)

		// 临时禁用PII脱敏处理以测试同步
		// if msg.Content != "" {
		//	processedContent, report := s.detector.DetectAndReplace(msg.Content)
		//	if report.TotalDetections > 0 {
		//		logger.Debug(fmt.Sprintf("🔒 PII detected in message %s: %d items", msg.MessageID, report.TotalDetections))
		//	}
		//	msg.Content = processedContent
		// }

		// 保存到数据库
		if err := s.repository.Message.CreateOrUpdate(msg); err != nil {
			logger.Warn("⚠️ Failed to save message " + msg.MessageID + ": " + err.Error())
			result.Errors++
		} else {
			result.SuccessCount++
		}
	}

	// 同步联系人
	if err := s.syncContacts(); err != nil {
		logger.Warn("⚠️ Failed to sync contacts: " + err.Error())
		result.Errors++
	}

	// 更新最后同步时间
	now := time.Now()
	nowStr := now.Format(time.RFC3339)
	if err := s.repository.Setting.Set("last_sync_time", nowStr, "最后同步时间"); err != nil {
		logger.Warn("⚠️ Failed to update last sync time: " + err.Error())
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	logger.Info(fmt.Sprintf("✅ Sync completed: %d messages, %d errors, %v",
		result.SuccessCount, result.Errors, result.Duration))

	return result, nil
}

// syncContacts 同步联系人
func (s *WeChatService) syncContacts() error {
	contacts, err := s.reader.GetContacts()
	if err != nil {
		return fmt.Errorf("failed to read contacts: %w", err)
	}

	for _, contact := range contacts {
		if err := s.repository.Contact.CreateOrUpdate(contact); err != nil {
			logger.Warn("⚠️ Failed to save contact " + contact.UserName + ": " + err.Error())
		}
	}

	logger.Info(fmt.Sprintf("📇 Synced %d contacts", len(contacts)))
	return nil
}

// GetMessages 获取消息列表
func (s *WeChatService) GetMessages(limit, offset int) ([]*model.Message, error) {
	return s.repository.Message.GetList(limit, offset)
}

// GetMessagesByTalker 根据聊天对象获取消息
func (s *WeChatService) GetMessagesByTalker(talkerID string, limit int) ([]*model.Message, error) {
	return s.repository.Message.GetByTalker(talkerID, limit)
}

// SearchMessages 搜索消息
func (s *WeChatService) SearchMessages(query string, limit int) ([]*model.Message, error) {
	return s.repository.Message.Search(query, limit)
}

// GetContacts 获取联系人列表
func (s *WeChatService) GetContacts() ([]*model.Contact, error) {
	return s.repository.Contact.GetAll()
}

// GetContactByUserName 根据用户名获取联系人
func (s *WeChatService) GetContactByUserName(userName string) (*model.Contact, error) {
	return s.repository.Contact.GetByUserName(userName)
}

// WeChatStatus 微信状态
type WeChatStatus struct {
	IsConnected   bool      `json:"is_connected"`
	DBPath        string    `json:"db_path"`
	MessageCount  *int64    `json:"message_count,omitempty"`
	ContactCount  int       `json:"contact_count"`
	LastSync      time.Time `json:"last_sync"`
}

// SyncResult 同步结果
type SyncResult struct {
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Duration     time.Duration `json:"duration"`
	SuccessCount int       `json:"success_count"`
	Errors       int       `json:"errors"`
}