package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hrygo/log"
)

// AutoSyncManager 自动同步管理器
type AutoSyncManager struct {
	syncService *RealtimeSyncService
	config      *AutoSyncConfig
}

// AutoSyncConfig 自动同步配置
type AutoSyncConfig struct {
	// 自动启动同步服务
	AutoStart      bool          `json:"auto_start"`       // 应用启动时自动启动同步
	AutoDetect     bool          `json:"auto_detect"`      // 自动检测微信数据库
	StartupDelay   time.Duration `json:"startup_delay"`    // 启动延迟时间

	// 手动指定路径
	ManualDBPath   string        `json:"manual_db_path"`   // 手动指定的数据库路径
	PreferManual   bool          `json:"prefer_manual"`    // 优先使用手动路径

	// 搜索路径配置
	SearchPaths    []string      `json:"search_paths"`     // 搜索路径列表
	DefaultInterval int          `json:"default_interval"` // 默认同步间隔(秒)

	// 故障恢复配置
	MaxRetries     int           `json:"max_retries"`      // 最大重试次数
	RetryDelay     time.Duration `json:"retry_delay"`      // 重试延迟

	// 通知配置
	EnableNotifications bool       `json:"enable_notifications"` // 启用通知
	NotificationChannel   string     `json:"notification_channel"`   // 通知渠道
}

// DefaultAutoSyncConfig 默认自动同步配置
func DefaultAutoSyncConfig() *AutoSyncConfig {
	return &AutoSyncConfig{
		AutoStart:            true,
		AutoDetect:           true,
		StartupDelay:         time.Second * 5,
		ManualDBPath:         "",
		PreferManual:         false,
		DefaultInterval:      60,
		MaxRetries:           3,
		RetryDelay:           time.Second * 3,
		EnableNotifications:  true,
		NotificationChannel:  "log",
	}
}

// NewAutoSyncManager 创建自动同步管理器
func NewAutoSyncManager(syncService *RealtimeSyncService) *AutoSyncManager {
	return &AutoSyncManager{
		syncService: syncService,
		config:      DefaultAutoSyncConfig(),
	}
}

// SetConfig 设置自动同步配置
func (m *AutoSyncManager) SetConfig(config *AutoSyncConfig) {
	m.config = config
}

// GetConfig 获取当前配置
func (m *AutoSyncManager) GetConfig() *AutoSyncConfig {
	return m.config
}

// StartOnStartup 应用启动时启动自动同步
func (m *AutoSyncManager) StartOnStartup() error {
	if !m.config.AutoStart {
		log.Info("📝 自动同步已禁用，跳过启动")
		return nil
	}

	log.Info("🚀 启动应用启动自动同步管理器")

	// 延迟启动，确保系统稳定
	if m.config.StartupDelay > 0 {
		log.Infof("⏳ 等待 %v 后启动自动同步", m.config.StartupDelay)
		time.Sleep(m.config.StartupDelay)
	}

	// 检测微信数据库路径
	wechatDBPath, err := m.detectWeChatDatabase()
	if err != nil {
		return fmt.Errorf("检测微信数据库失败: %w", err)
	}

	// 启动同步服务
	return m.startSyncService(wechatDBPath)
}

// detectWeChatDatabase 检测微信数据库路径
func (m *AutoSyncManager) detectWeChatDatabase() (string, error) {
	// 如果有手动指定的路径且优先使用，直接返回
	if m.config.PreferManual && m.config.ManualDBPath != "" {
		if _, err := os.Stat(m.config.ManualDBPath); err == nil {
			log.Infof("✅ 使用手动指定路径: %s", m.config.ManualDBPath)
			return m.config.ManualDBPath, nil
		} else {
			log.Warnf("⚠️ 手动指定路径不存在: %s", m.config.ManualDBPath)
		}
	}

	if !m.config.AutoDetect {
		if m.config.ManualDBPath != "" {
			if _, err := os.Stat(m.config.ManualDBPath); err == nil {
				log.Infof("✅ 使用手动指定路径: %s", m.config.ManualDBPath)
				return m.config.ManualDBPath, nil
			}
		}
		return "", fmt.Errorf("自动检测已禁用且未提供手动路径")
	}

	log.Info("🔍 开始自动检测微信数据库")

	// 首先尝试检测常见的微信数据库路径
	searchPaths := m.getSearchPaths()
	var detectedPaths []string

	// 添加手动路径到搜索列表中
	if m.config.ManualDBPath != "" {
		searchPaths = append([]string{m.config.ManualDBPath}, searchPaths...)
	}

	for _, path := range searchPaths {
		// 如果是精确路径（不包含通配符），直接检查
		if !containsWildcard(path) {
			if _, err := os.Stat(path); err == nil {
				detectedPaths = append(detectedPaths, path)
				log.Debugf("发现微信数据库: %s", path)
			}
			continue
		}

		// 处理通配符路径
		matches, err := filepath.Glob(path)
		if err != nil {
			log.Debugf("搜索路径失败 %s: %v", path, err)
			continue
		}

		for _, match := range matches {
			if _, err := os.Stat(match); err == nil {
				detectedPaths = append(detectedPaths, match)
				log.Debugf("发现微信数据库: %s", match)
			}
		}
	}

	if len(detectedPaths) == 0 {
		return "", fmt.Errorf("未检测到微信数据库文件")
	}

	// 选择最新的数据库文件
	latestPath := m.selectLatestDatabase(detectedPaths)
	log.Infof("✅ 检测到微信数据库: %s", latestPath)

	return latestPath, nil
}

// containsWildcard 检查路径是否包含通配符
func containsWildcard(path string) bool {
	return strings.Contains(path, "*") || strings.Contains(path, "?")
}

// getSearchPaths 获取搜索路径
func (m *AutoSyncManager) getSearchPaths() []string {
	if len(m.config.SearchPaths) > 0 {
		return m.config.SearchPaths
	}

	// 默认搜索路径
	homeDir, _ := os.UserHomeDir()
	return []string{
		filepath.Join(homeDir, "Library/Containers/com.tencent.xin/Data/Library/Application Support/com.tencent.xin/*/Message/MessageTemp/*/DB/EnMicroMsg.db"),
		filepath.Join(homeDir, "Documents/WeChat Files/*/Msg/*/EnMicroMsg.db"),
		filepath.Join(homeDir, "Library/Application Support/Tencent/MicroMsg/*/DB/MicroMsg.db"),
		filepath.Join(homeDir, "AppData/Roaming/Tencent/MicroMsg/*/DB/MicroMsg.db"),
	}
}

// selectLatestDatabase 选择最新的数据库文件
func (m *AutoSyncManager) selectLatestDatabase(paths []string) string {
	var latestPath string
	var latestTime time.Time

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}

		if info.ModTime().After(latestTime) {
			latestTime = info.ModTime()
			latestPath = path
		}
	}

	if latestPath == "" {
		return paths[0] // 返回第一个可用路径
	}

	return latestPath
}

// startSyncService 启动同步服务
func (m *AutoSyncManager) startSyncService(wechatDBPath string) error {
	log.Info("🚀 启动自动同步服务")

	// 创建同步配置
	syncConfig := &SyncConfig{
		WeChatDBPath:     wechatDBPath,
		AutoDetectPath:   false, // 已经检测到路径
		SyncInterval:     time.Duration(m.config.DefaultInterval) * time.Second,
		EnableMonitoring: true,
		MaxSyncRetries:   m.config.MaxRetries,
		SyncBatchSize:    1000,
		MinMessageAge:    time.Second * 10,
		WorkerCount:      2,
		BufferSize:       10000,
	}

	// 设置配置
	m.syncService.SetConfig(syncConfig)

	// 重试启动同步服务
	var lastErr error
	for i := 0; i < m.config.MaxRetries; i++ {
		if i > 0 {
			log.Infof("🔄 重试启动同步服务 (%d/%d)", i+1, m.config.MaxRetries)
			time.Sleep(m.config.RetryDelay)
		}

		err := m.syncService.Start()
		if err == nil {
			log.Info("✅ 自动同步服务启动成功")
			m.sendNotification("success", "自动同步服务启动成功", fmt.Sprintf("监控数据库: %s", wechatDBPath))
			return nil
		}

		lastErr = err
		log.Warnf("❌ 启动同步服务失败: %v", err)
	}

	// 发送失败通知
	errorMsg := fmt.Sprintf("自动同步服务启动失败: %v", lastErr)
	m.sendNotification("error", "自动同步启动失败", errorMsg)

	return fmt.Errorf("启动同步服务失败，已重试 %d 次: %w", m.config.MaxRetries, lastErr)
}

// sendNotification 发送通知
func (m *AutoSyncManager) sendNotification(level, title, message string) {
	if !m.config.EnableNotifications {
		return
	}

	switch m.config.NotificationChannel {
	case "system":
		log.Infof("🔔 通知 [%s]: %s - %s", level, title, message)
	case "log":
		switch level {
		case "success":
			log.Infof("✅ %s: %s", title, message)
		case "error":
			log.Errorf("❌ %s: %s", title, message)
		case "warning":
			log.Warnf("⚠️ %s: %s", title, message)
		default:
			log.Infof("📢 %s: %s", title, message)
		}
	}
}

// StopOnShutdown 应用关闭时停止自动同步
func (m *AutoSyncManager) StopOnShutdown() error {
	log.Info("🛑 停止应用关闭自动同步管理器")

	err := m.syncService.Stop()
	if err != nil {
		log.Warnf("停止同步服务失败: %v", err)
		return err
	}

	log.Info("✅ 自动同步服务已停止")
	return nil
}

// GetStatus 获取自动同步状态
func (m *AutoSyncManager) GetStatus() map[string]interface{} {
	syncStatus := m.syncService.GetStatus()

	return map[string]interface{}{
		"auto_start_enabled":   m.config.AutoStart,
		"auto_detect_enabled":  m.config.AutoDetect,
		"sync_service_status": syncStatus,
		"config":              m.config,
	}
}