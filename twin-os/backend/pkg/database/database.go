package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"twin-os/backend/pkg/logger"
)

var DB *sql.DB

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

	// 设置连接池
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(time.Hour)

	logger.Info("✅ Database connected successfully")

	// 创建表
	if err = createTables(); err != nil {
		logger.Fatal("❌ Failed to create tables: " + err.Error())
	}

	logger.Info("✅ Database tables created/verified")

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

	// 创建索引
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_messages_talker_id ON messages(talker_id)",
		"CREATE INDEX IF NOT EXISTS idx_messages_timestamp ON messages(timestamp)",
		"CREATE INDEX IF NOT EXISTS idx_analysis_results_type ON analysis_results(type)",
		"CREATE INDEX IF NOT EXISTS idx_analysis_results_created_at ON analysis_results(created_at)",
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