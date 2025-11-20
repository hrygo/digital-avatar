package services

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/hrygo/log"
	"github.com/rjeczalik/notify"
)

// RealtimeSyncService 微信数据库准实时同步服务
type RealtimeSyncService struct {
	db            *sql.DB
	wechatService *WeChatServiceV2

	// 同步配置
	config        *SyncConfig

	// 状态管理
	isRunning     bool
	isMonitoring  bool
	lastSyncTime  int64
	syncInterval  time.Duration

	// 控制通道
	stopChan      chan struct{}
	syncChan      chan struct{}

	// 互斥锁
	mu            sync.RWMutex

	// 文件监控
	watcher       notifyWatcher
}

// SyncConfig 同步配置
type SyncConfig struct {
	// 基本配置
	WeChatDBPath      string        `json:"wechat_db_path"`      // 微信数据库路径
	AutoDetectPath    bool          `json:"auto_detect_path"`     // 自动检测路径
	SyncInterval      time.Duration `json:"sync_interval"`        // 同步间隔

	// 高级配置
	MaxSyncRetries    int           `json:"max_sync_retries"`     // 最大重试次数
	SyncBatchSize     int           `json:"sync_batch_size"`      // 批量同步大小
	EnableMonitoring  bool          `json:"enable_monitoring"`    // 启用文件监控

	// 过滤配置
	MinMessageAge     time.Duration `json:"min_message_age"`      // 最小消息年龄（避免刚发送的消息）
	ChatFilters       []string      `json:"chat_filters"`         // 聊天过滤器

	// 性能配置
	WorkerCount       int           `json:"worker_count"`         // 工作协程数量
	BufferSize        int           `json:"buffer_size"`          // 缓冲区大小
}

// DefaultSyncConfig 默认同步配置
func DefaultSyncConfig() *SyncConfig {
	return &SyncConfig{
		AutoDetectPath:   true,
		SyncInterval:     time.Second * 30,
		MaxSyncRetries:   3,
		SyncBatchSize:    1000,
		EnableMonitoring: true,
		MinMessageAge:    time.Second * 10,
		WorkerCount:      2,
		BufferSize:       10000,
	}
}

// notifyWatcher 文件监控接口
type notifyWatcher interface {
	Watch(path string, eventChan chan<- notify.EventInfo, events ...notify.Event) error
	Stop(eventChan chan<- notify.EventInfo) error
}

// NewRealtimeSyncService 创建实时同步服务
func NewRealtimeSyncService(db *sql.DB, wechatService *WeChatServiceV2) *RealtimeSyncService {
	return &RealtimeSyncService{
		db:            db,
		wechatService: wechatService,
		config:        DefaultSyncConfig(),
		syncInterval:  time.Second * 30,
		stopChan:      make(chan struct{}),
		syncChan:      make(chan struct{}, 1),
		watcher:       &notifyWrapper{},
	}
}

// SetConfig 设置同步配置
func (s *RealtimeSyncService) SetConfig(config *SyncConfig) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config = config
	s.syncInterval = config.SyncInterval
}

// GetConfig 获取当前配置
func (s *RealtimeSyncService) GetConfig() *SyncConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.config
}

// Start 启动实时同步服务
func (s *RealtimeSyncService) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return fmt.Errorf("实时同步服务已经在运行")
	}

	log.Info("🚀 启动微信数据库准实时同步服务")

	// 1. 检查并获取微信数据库路径
	wechatDBPath, err := s.getWeChatDBPath()
	if err != nil {
		return fmt.Errorf("获取微信数据库路径失败: %w", err)
	}
	s.config.WeChatDBPath = wechatDBPath

	// 2. 初始化最后同步时间
	if err := s.initLastSyncTime(); err != nil {
		log.Warnf("初始化最后同步时间失败: %v", err)
	}

	// 3. 启动文件监控（如果启用）
	if s.config.EnableMonitoring {
		if err := s.startFileMonitoring(); err != nil {
			log.Warnf("启动文件监控失败: %v", err)
		}
	}

	// 4. 启动定时同步
	s.isRunning = true
	go s.runSyncScheduler()

	log.Infof("✅ 实时同步服务已启动，数据库路径: %s", s.config.WeChatDBPath)
	log.Infof("📊 同步间隔: %v, 文件监控: %v", s.syncInterval, s.config.EnableMonitoring)

	return nil
}

// Stop 停止实时同步服务
func (s *RealtimeSyncService) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return nil
	}

	log.Info("🛑 停止微信数据库准实时同步服务")

	// 1. 停止同步调度器
	s.isRunning = false
	close(s.stopChan)

	// 2. 停止文件监控
	if s.isMonitoring {
		if err := s.stopFileMonitoring(); err != nil {
			log.Warnf("停止文件监控失败: %v", err)
		}
	}

	log.Info("✅ 实时同步服务已停止")
	return nil
}

// getWeChatDBPath 获取微信数据库路径
func (s *RealtimeSyncService) getWeChatDBPath() (string, error) {
	if s.config.WeChatDBPath != "" {
		// 验证指定路径是否存在
		if _, err := os.Stat(s.config.WeChatDBPath); err != nil {
			return "", fmt.Errorf("指定的微信数据库文件不存在: %s", s.config.WeChatDBPath)
		}
		return s.config.WeChatDBPath, nil
	}

	if s.config.AutoDetectPath {
		// 自动检测微信数据库路径
		paths := s.detectWeChatDBPaths()
		if len(paths) == 0 {
			return "", fmt.Errorf("未检测到微信数据库文件，请手动指定路径")
		}

		// 返回最新找到的数据库
		return paths[0], nil
	}

	return "", fmt.Errorf("未指定微信数据库路径且未启用自动检测")
}

// detectWeChatDBPaths 检测微信数据库路径
func (s *RealtimeSyncService) detectWeChatDBPaths() []string {
	var paths []string

	// 常见的微信数据路径
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

	return paths
}

// initLastSyncTime 初始化最后同步时间
func (s *RealtimeSyncService) initLastSyncTime() error {
	var lastSync sql.NullInt64
	err := s.db.QueryRow(`
		SELECT MAX(start_time) FROM wechat_sync_records
		WHERE status = 'completed' AND sync_type = 'incremental'
	`).Scan(&lastSync)

	if err == nil && lastSync.Valid {
		s.lastSyncTime = lastSync.Int64
		log.Infof("📅 初始化最后同步时间: %s", time.Unix(s.lastSyncTime, 0).Format("2006-01-02 15:04:05"))
	} else {
		// 设置为24小时前
		s.lastSyncTime = time.Now().Add(-24 * time.Hour).Unix()
		log.Infof("📅 设置初始同步时间为24小时前")
	}

	return nil
}

// startFileMonitoring 启动文件监控
func (s *RealtimeSyncService) startFileMonitoring() error {
	dbDir := filepath.Dir(s.config.WeChatDBPath)

	eventChan := make(chan notify.EventInfo, 1)
	err := s.watcher.Watch(dbDir, eventChan, notify.Write)
	if err != nil {
		return fmt.Errorf("启动文件监控失败: %w", err)
	}

	s.isMonitoring = true
	go s.handleFileEvents(eventChan)

	log.Infof("👁️ 启动文件监控，监控目录: %s", dbDir)
	return nil
}

// stopFileMonitoring 停止文件监控
func (s *RealtimeSyncService) stopFileMonitoring() error {
	if !s.isMonitoring {
		return nil
	}

	s.isMonitoring = false
	eventChan := make(chan notify.EventInfo)
	return s.watcher.Stop(eventChan)
}

// handleFileEvents 处理文件事件
func (s *RealtimeSyncService) handleFileEvents(eventChan chan notify.EventInfo) {
	for event := range eventChan {
		if !s.isMonitoring {
			break
		}

		// 只关心微信数据库文件的变化
		if filepath.Base(event.Path()) != filepath.Base(s.config.WeChatDBPath) {
			continue
		}

		log.Debugf("📝 检测到微信数据库文件变化: %s", event.Path())

		// 防抖：短时间内多次变化只触发一次同步
		select {
		case s.syncChan <- struct{}{}:
			log.Debug("🔄 触发增量同步")
		default:
			// 同步信号已存在，忽略
		}
	}
}

// runSyncScheduler 运行同步调度器
func (s *RealtimeSyncService) runSyncScheduler() {
	ticker := time.NewTicker(s.syncInterval)
	defer ticker.Stop()

	log.Infof("⏰ 启动定时同步调度器，间隔: %v", s.syncInterval)

	// 首次同步
	s.triggerSync()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			s.triggerSync()
		case <-s.syncChan:
			s.triggerSync()
		}
	}
}

// triggerSync 触发同步
func (s *RealtimeSyncService) triggerSync() {
	if !s.isRunning {
		return
	}

	// 防止并发同步
	s.mu.Lock()

	go func() {
		defer s.mu.Unlock()

		log.Debug("🚀 开始执行增量同步")

		// 添加最小消息年龄延迟
		minTime := time.Now().Add(-s.config.MinMessageAge).Unix()
		syncTime := minTime
		if s.lastSyncTime > minTime {
			syncTime = s.lastSyncTime
		}

		record, err := s.wechatService.IncrementalSync(s.config.WeChatDBPath, syncTime)
		if err != nil {
			log.Errorf("💥 增量同步失败: %v", err)
			return
		}

		if record != nil {
			s.lastSyncTime = record.StartTime.Unix()
			log.Infof("✅ 增量同步完成，新增消息: %d, 新增联系人: %d, 耗时: %.2fs",
				record.NewMessages, record.NewContacts, record.Duration)
		}
	}()
}

// GetStatus 获取同步服务状态
func (s *RealtimeSyncService) GetStatus() *SyncStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return &SyncStatus{
		IsRunning:        s.isRunning,
		IsMonitoring:     s.isMonitoring,
		WeChatDBPath:     s.config.WeChatDBPath,
		LastSyncTime:     s.lastSyncTime,
		SyncInterval:     s.syncInterval.String(),
		NextSyncTime:     time.Now().Add(s.syncInterval).Unix(),
		TotalSyncs:       s.getTotalSyncCount(),
		LastSyncSuccess:  s.getLastSyncSuccess(),
	}
}

// SyncStatus 同步状态
type SyncStatus struct {
	IsRunning        bool   `json:"is_running"`
	IsMonitoring     bool   `json:"is_monitoring"`
	WeChatDBPath     string `json:"wechat_db_path"`
	LastSyncTime     int64  `json:"last_sync_time"`
	SyncInterval     string `json:"sync_interval"`
	NextSyncTime     int64  `json:"next_sync_time"`
	TotalSyncs       int64  `json:"total_syncs"`
	LastSyncSuccess  bool   `json:"last_sync_success"`
}

// getTotalSyncCount 获取总同步次数
func (s *RealtimeSyncService) getTotalSyncCount() int64 {
	var count int64
	s.db.QueryRow("SELECT COUNT(*) FROM wechat_sync_records").Scan(&count)
	return count
}

// getLastSyncSuccess 获取最后一次同步是否成功
func (s *RealtimeSyncService) getLastSyncSuccess() bool {
	var status string
	err := s.db.QueryRow(`
		SELECT status FROM wechat_sync_records
		ORDER BY created_at DESC LIMIT 1
	`).Scan(&status)

	return err == nil && status == "completed"
}

// TriggerManualSync 手动触发同步
func (s *RealtimeSyncService) TriggerManualSync() error {
	if !s.isRunning {
		return fmt.Errorf("同步服务未运行")
	}

	select {
	case s.syncChan <- struct{}{}:
		log.Info("🔄 手动触发同步")
		return nil
	default:
		return fmt.Errorf("同步正在进行中，请稍后再试")
	}
}

// notifyWrapper notify包装器
type notifyWrapper struct{}

func (w *notifyWrapper) Watch(path string, eventChan chan<- notify.EventInfo, events ...notify.Event) error {
	return notify.Watch(path, eventChan, events...)
}

func (w *notifyWrapper) Stop(eventChan chan<- notify.EventInfo) error {
	notify.Stop(eventChan)
	return nil
}