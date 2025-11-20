package database

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"twin-os/backend/pkg/logger"
)

var DB *sql.DB

// PerformanceStats 数据库性能统计
type PerformanceStats struct {
	TotalQueries       int64 `json:"total_queries"`
	SlowQueries        int64 `json:"slow_queries"`
	TotalConnections   int64 `json:"total_connections"`
	ActiveConnections  int64 `json:"active_connections"`
	AvgResponseTime    int64 `json:"avg_response_time_micros"`
	MaxResponseTime    int64 `json:"max_response_time_micros"`
}

var (
	perfStats = &PerformanceStats{}
	statsMutex = &sync.RWMutex{}
	slowQueryThreshold = 100 * time.Millisecond
)

// GetPerformanceStats 获取数据库性能统计
func GetPerformanceStats() PerformanceStats {
	statsMutex.RLock()
	defer statsMutex.RUnlock()
	return *perfStats
}

// trackQueryPerformance 追踪查询性能
func trackQueryPerformance(ctx context.Context, query string, start time.Time, args ...interface{}) {
	duration := time.Since(start)
	micros := duration.Microseconds()

	// 原子操作更新统计
	atomic.AddInt64(&perfStats.TotalQueries, 1)
	atomic.AddInt64(&perfStats.AvgResponseTime, micros)

	// 更新最大响应时间
	for {
		currentMax := atomic.LoadInt64(&perfStats.MaxResponseTime)
		if micros <= currentMax || atomic.CompareAndSwapInt64(&perfStats.MaxResponseTime, currentMax, micros) {
			break
		}
	}

	// 慢查询检测和记录
	if duration > slowQueryThreshold {
		atomic.AddInt64(&perfStats.SlowQueries, 1)
		logger.Warn(fmt.Sprintf("🐌 Slow Query (%v): %s", duration, query))
	}
}

// Init 初始化数据库
func Init(dbPath string) *sql.DB {
	var err error

	// 打开数据库连接
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.Fatal("❌ Failed to connect to database: " + err.Error())
	}

	// 测试连接
	if err = DB.Ping(); err != nil {
		logger.Fatal("❌ Database ping failed: " + err.Error())
	}

	// 优化连接池配置 - v0.3.0性能优化
	// 根据系统负载和并发需求调整连接池参数
	DB.SetMaxOpenConns(25)           // 增加最大连接数以支持更高并发
	DB.SetMaxIdleConns(10)           // 增加空闲连接数以减少连接建立开销
	DB.SetConnMaxLifetime(30 * time.Minute)  // 减少连接生命周期避免长连接问题
	DB.SetConnMaxIdleTime(5 * time.Minute)   // 设置空闲连接超时，回收空闲连接

	logger.Info("✅ Database connected successfully")

	// 启用性能监控
	if err = enablePerformanceMonitoring(); err != nil {
		logger.Warn("⚠️ Failed to enable performance monitoring: " + err.Error())
	}

	// 创建表
	if err = createTables(); err != nil {
		logger.Fatal("❌ Failed to create tables: " + err.Error())
	}

	logger.Info("✅ Database tables created/verified")

	// 启动健康检查
	go startHealthCheck()

	return DB
}

// createTables 创建数据库表
func createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			message_id TEXT UNIQUE NOT NULL,
			talker_id TEXT NOT NULL,
			type INTEGER NOT NULL,
			content TEXT,
			timestamp DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS contacts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_name TEXT UNIQUE NOT NULL,
			nick_name TEXT,
			remark TEXT,
			type INTEGER NOT NULL DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS analysis_results (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			type TEXT NOT NULL,
			message_ids TEXT NOT NULL,
			content TEXT NOT NULL,
			metadata TEXT,
			confidence REAL DEFAULT 0.0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			processed_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS system_settings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			key TEXT UNIQUE NOT NULL,
			value TEXT NOT NULL,
			description TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query %s: %w", query, err)
		}
	}

	// 创建复合索引 - v0.3.0性能优化
	indexes := []string{
		// 基础索引
		"CREATE INDEX IF NOT EXISTS idx_messages_talker_id ON messages(talker_id)",
		"CREATE INDEX IF NOT EXISTS idx_messages_timestamp ON messages(timestamp)",
		"CREATE INDEX IF NOT EXISTS idx_messages_message_id ON messages(message_id)",
		"CREATE INDEX IF NOT EXISTS idx_analysis_results_type ON analysis_results(type)",
		"CREATE INDEX IF NOT EXISTS idx_analysis_results_created_at ON analysis_results(created_at)",

		// 复合索引 - 优化常见查询组合
		"CREATE INDEX IF NOT EXISTS idx_messages_talker_timestamp ON messages(talker_id, timestamp DESC)",
		"CREATE INDEX IF NOT EXISTS idx_messages_timestamp_content ON messages(timestamp DESC) WHERE content != ''",
		"CREATE INDEX IF NOT EXISTS idx_analysis_results_type_created ON analysis_results(type, created_at DESC)",

		// 联系人索引
		"CREATE INDEX IF NOT EXISTS idx_contacts_user_name ON contacts(user_name)",
		"CREATE INDEX IF NOT EXISTS idx_contacts_type ON contacts(type)",

		// 设置索引
		"CREATE INDEX IF NOT EXISTS idx_system_settings_key ON system_settings(key)",
	}

	for _, index := range indexes {
		if _, err := DB.Exec(index); err != nil {
			return fmt.Errorf("failed to create index %s: %w", index, err)
		}
	}

	// 插入默认系统设置
	defaultSettings := map[string]string{
		"wechat_db_path":          "",
		"auto_sync_interval":      "300", // 5分钟
		"analysis_enabled":        "true",
		"pii_detection_enabled":   "true",
		"deepseek_api_key":        "",
		"encryption_enabled":      "true",
		"last_sync_time":          "0",
		"theme":                   "dark",
		"language":                "zh-CN",
	}

	for key, value := range defaultSettings {
		query := `INSERT OR IGNORE INTO system_settings (key, value, description) VALUES (?, ?, ?)`
		description := getSettingDescription(key)
		if _, err := DB.Exec(query, key, value, description); err != nil {
			logger.Warn("⚠️ Failed to insert default setting " + key + ": " + err.Error())
		}
	}

	return nil
}

// getSettingDescription 获取设置项的描述
func getSettingDescription(key string) string {
	descriptions := map[string]string{
		"wechat_db_path":          "微信数据库路径",
		"auto_sync_interval":      "自动同步间隔(秒)",
		"analysis_enabled":        "是否启用分析",
		"pii_detection_enabled":   "是否启用PII检测",
		"deepseek_api_key":        "DeepSeek API密钥",
		"encryption_enabled":      "是否启用加密",
		"last_sync_time":          "最后同步时间",
		"theme":                   "界面主题",
		"language":                "界面语言",
	}

	if desc, exists := descriptions[key]; exists {
		return desc
	}
	return ""
}

// Close 关闭数据库连接
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// enablePerformanceMonitoring 启用数据库性能监控
func enablePerformanceMonitoring() error {
	// SQLite 特定的性能优化设置
	queries := []string{
		"PRAGMA journal_mode = WAL",           // 写前日志模式，提升并发性能
		"PRAGMA synchronous = NORMAL",         // 平衡性能和安全性
		"PRAGMA cache_size = 10000",          // 增大缓存到10MB
		"PRAGMA temp_store = MEMORY",         // 临时表存储在内存中
		"PRAGMA mmap_size = 268435456",       // 256MB内存映射
		"PRAGMA optimize",                    // 自动优化查询计划
	}

	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			logger.Warn("⚠️ Failed to set PRAGMA " + query + ": " + err.Error())
			return err
		}
	}

	logger.Info("✅ Database performance monitoring enabled")
	return nil
}

// startHealthCheck 启动数据库健康检查
func startHealthCheck() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		stats := DB.Stats()

		// 更新连接统计
		atomic.StoreInt64(&perfStats.ActiveConnections, int64(stats.OpenConnections))

		// 检查连接池健康状态
		if stats.OpenConnections > 20 {
			logger.Warn(fmt.Sprintf("⚠️ High connection usage: %d/%d",
				stats.OpenConnections, stats.MaxOpenConnections))
		}

		// 执行健康检查查询
		start := time.Now()
		if err := DB.Ping(); err != nil {
			logger.Error("❌ Database health check failed: " + err.Error())
		} else {
			trackQueryPerformance(context.Background(), "HEALTH_CHECK", start)
		}
	}
}

// GetDatabaseHealth 获取数据库健康状态
func GetDatabaseHealth() map[string]interface{} {
	stats := DB.Stats()
	perfStats := GetPerformanceStats()

	health := map[string]interface{}{
		"status": "healthy",
		"connections": map[string]interface{}{
			"open":       stats.OpenConnections,
			"in_use":     stats.InUse,
			"idle":       stats.Idle,
			"max_open":   stats.MaxOpenConnections,
		},
		"performance": map[string]interface{}{
			"total_queries":    perfStats.TotalQueries,
			"slow_queries":     perfStats.SlowQueries,
			"avg_response_ms":  perfStats.AvgResponseTime / 1000,
			"max_response_ms":  perfStats.MaxResponseTime / 1000,
		},
		"timestamp": time.Now().Unix(),
	}

	// 判断健康状态
	if stats.OpenConnections > stats.MaxOpenConnections*9/10 {
		health["status"] = "degraded"
	}

	if perfStats.SlowQueries > perfStats.TotalQueries/10 {
		health["status"] = "poor"
	}

	return health
}

// ResetPerformanceStats 重置性能统计
func ResetPerformanceStats() {
	atomic.StoreInt64(&perfStats.TotalQueries, 0)
	atomic.StoreInt64(&perfStats.SlowQueries, 0)
	atomic.StoreInt64(&perfStats.AvgResponseTime, 0)
	atomic.StoreInt64(&perfStats.MaxResponseTime, 0)
}