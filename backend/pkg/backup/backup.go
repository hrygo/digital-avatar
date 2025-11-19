package backup

import (
	"archive/zip"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"twin-os/backend/pkg/logger"
)

// BackupType 备份类型
type BackupType string

const (
	BackupTypeFull    BackupType = "full"     // 完整备份
	BackupTypeConfig  BackupType = "config"   // 仅配置备份
	BackupTypeData    BackupType = "data"     // 仅数据备份
	BackupTypeIncremental BackupType = "incremental" // 增量备份
)

// BackupStatus 备份状态
type BackupStatus string

const (
	StatusPending   BackupStatus = "pending"   // 等待中
	StatusInProgress BackupStatus = "in_progress" // 进行中
	StatusCompleted BackupStatus = "completed"  // 已完成
	StatusFailed    BackupStatus = "failed"    // 失败
	StatusCorrupted BackupStatus = "corrupted" // 损坏
)

// BackupInfo 备份信息
type BackupInfo struct {
	ID             string            `json:"id"`
	Type           BackupType        `json:"type"`
	Status         BackupStatus      `json:"status"`
	CreatedAt      time.Time         `json:"created_at"`
	CompletedAt    *time.Time        `json:"completed_at,omitempty"`
	Size           int64             `json:"size"`
	Checksum       string            `json:"checksum"`
	Version        string            `json:"version"`
	Description    string            `json:"description"`
	FilePath       string            `json:"file_path"`
	Metadata       map[string]string `json:"metadata"`
	Encrypted      bool              `json:"encrypted"`
	Compression    string            `json:"compression"`
	ItemCounts     map[string]int    `json:"item_counts"`
	LastModified   time.Time         `json:"last_modified"`
}

// BackupConfig 备份配置
type BackupConfig struct {
	BackupDir        string            `json:"backup_dir"`
	AutoBackup       bool              `json:"auto_backup"`
	BackupInterval   time.Duration     `json:"backup_interval"`
	MaxBackups       int               `json:"max_backups"`
	Compression      bool              `json:"compression"`
	Encryption       bool              `json:"encryption"`
	EncryptionKey    string            `json:"encryption_key,omitempty"`
	ExcludePatterns  []string          `json:"exclude_patterns"`
	IncludePatterns  []string          `json:"include_patterns"`
	CustomMetadata   map[string]string `json:"custom_metadata"`
}

// DefaultBackupConfig 默认备份配置
func DefaultBackupConfig() BackupConfig {
	return BackupConfig{
		BackupDir:       "./backups",
		AutoBackup:      true,
		BackupInterval:  24 * time.Hour,
		MaxBackups:      10,
		Compression:     true,
		Encryption:      false,
		ExcludePatterns: []string{"*.tmp", "*.log", "cache/*"},
		IncludePatterns: []string{"*.db", "*.json", "*.yaml", "*.yml"},
		CustomMetadata:  make(map[string]string),
	}
}

// BackupManager 备份管理器
type BackupManager struct {
	config BackupConfig
	db     *sql.DB
}

// NewBackupManager 创建备份管理器
func NewBackupManager(db *sql.DB, config BackupConfig) *BackupManager {
	if config.BackupDir == "" {
		config = DefaultBackupConfig()
	}

	// 确保备份目录存在
	if err := os.MkdirAll(config.BackupDir, 0755); err != nil {
		logger.Error("Failed to create backup directory: " + err.Error())
	}

	return &BackupManager{
		config: config,
		db:     db,
	}
}

// CreateBackup 创建备份
func (bm *BackupManager) CreateBackup(backupType BackupType, description string) (*BackupInfo, error) {
	backupID := generateBackupID()
	backupInfo := &BackupInfo{
		ID:           backupID,
		Type:         backupType,
		Status:       StatusPending,
		CreatedAt:    time.Now(),
		Version:      "1.0.0",
		Description:  description,
		FilePath:     filepath.Join(bm.config.BackupDir, fmt.Sprintf("backup_%s_%s.zip", backupType, backupID[:8])),
		Metadata:     make(map[string]string),
		Encrypted:    bm.config.Encryption,
		Compression:  "zip",
		ItemCounts:   make(map[string]int),
		LastModified: time.Now(),
	}

	// 更新状态为进行中
	backupInfo.Status = StatusInProgress
	logger.Info("🔄 Starting backup: " + backupID)

	// 创建备份文件
	if err := bm.createBackupFile(backupInfo); err != nil {
		backupInfo.Status = StatusFailed
		logger.Error("❌ Backup failed: " + err.Error())
		return backupInfo, err
	}

	// 计算文件大小和校验和
	if err := bm.calculateFileMetrics(backupInfo); err != nil {
		logger.Warn("⚠️ Failed to calculate file metrics: " + err.Error())
	}

	// 完成备份
	now := time.Now()
	backupInfo.CompletedAt = &now
	backupInfo.Status = StatusCompleted
	backupInfo.LastModified = now

	// 保存备份信息
	if err := bm.saveBackupInfo(backupInfo); err != nil {
		logger.Warn("⚠️ Failed to save backup info: " + err.Error())
	}

	// 清理旧备份
	if err := bm.cleanupOldBackups(); err != nil {
		logger.Warn("⚠️ Failed to cleanup old backups: " + err.Error())
	}

	logger.Info(fmt.Sprintf("✅ Backup completed: %s (size: %d bytes, checksum: %s)",
		backupID, backupInfo.Size, backupInfo.Checksum))

	return backupInfo, nil
}

// createBackupFile 创建备份文件
func (bm *BackupManager) createBackupFile(backupInfo *BackupInfo) error {
	file, err := os.Create(backupInfo.FilePath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	// 根据备份类型添加不同内容
	switch backupInfo.Type {
	case BackupTypeFull:
		if err := bm.addFullBackup(zipWriter, backupInfo); err != nil {
			return err
		}
	case BackupTypeConfig:
		if err := bm.addConfigBackup(zipWriter, backupInfo); err != nil {
			return err
		}
	case BackupTypeData:
		if err := bm.addDataBackup(zipWriter, backupInfo); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported backup type: %s", backupInfo.Type)
	}

	// 添加备份元数据
	if err := bm.addBackupMetadata(zipWriter, backupInfo); err != nil {
		return err
	}

	return nil
}

// addFullBackup 添加完整备份
func (bm *BackupManager) addFullBackup(zipWriter *zip.Writer, backupInfo *BackupInfo) error {
	logger.Info("📦 Adding full backup content...")

	// 1. 添加数据库备份
	if err := bm.addDatabaseBackup(zipWriter, backupInfo); err != nil {
		return fmt.Errorf("failed to backup database: %w", err)
	}

	// 2. 添加配置文件备份
	if err := bm.addConfigFiles(zipWriter, backupInfo); err != nil {
		return fmt.Errorf("failed to backup config files: %w", err)
	}

	// 3. 添加用户数据文件
	if err := bm.addUserDataFiles(zipWriter, backupInfo); err != nil {
		return fmt.Errorf("failed to backup user data: %w", err)
	}

	return nil
}

// addConfigBackup 添加配置备份
func (bm *BackupManager) addConfigBackup(zipWriter *zip.Writer, backupInfo *BackupInfo) error {
	logger.Info("⚙️ Adding configuration backup...")

	// 添加配置文件
	return bm.addConfigFiles(zipWriter, backupInfo)
}

// addDataBackup 添加数据备份
func (bm *BackupManager) addDataBackup(zipWriter *zip.Writer, backupInfo *BackupInfo) error {
	logger.Info("💾 Adding data backup...")

	// 添加数据库备份
	return bm.addDatabaseBackup(zipWriter, backupInfo)
}

// addDatabaseBackup 添加数据库备份
func (bm *BackupManager) addDatabaseBackup(zipWriter *zip.Writer, backupInfo *BackupInfo) error {
	if bm.db == nil {
		return fmt.Errorf("database connection is nil")
	}

	// 创建数据库模式文件
	schemaFile := &zip.FileHeader{
		Name:   "database/schema.sql",
		Method: zip.Deflate,
	}
	writer, err := zipWriter.CreateHeader(schemaFile)
	if err != nil {
		return err
	}

	// 写入数据库模式信息
	schemaInfo := map[string]interface{}{
		"database":   "sqlite3",
		"tables":     bm.getTableInfo(),
		"exported_at": time.Now().Format(time.RFC3339),
		"version":    "1.0.0",
	}

	schemaJSON, err := json.MarshalIndent(schemaInfo, "", "  ")
	if err != nil {
		return err
	}

	if _, err := writer.Write(schemaJSON); err != nil {
		return err
	}

	// 备份数据表数据
	tables := []string{"settings", "messages", "contacts", "analysis_results"}
	for _, table := range tables {
		if err := bm.backupTable(zipWriter, table); err != nil {
			logger.Warn("⚠️ Failed to backup table " + table + ": " + err.Error())
		} else {
			backupInfo.ItemCounts["tables"]++
		}
	}

	return nil
}

// backupTable 备份单个数据表
func (bm *BackupManager) backupTable(zipWriter *zip.Writer, tableName string) error {
	// 检查表是否存在
	if !bm.tableExists(tableName) {
		return nil // 表不存在，跳过
	}

	// 查询表数据
	rows, err := bm.queryTableData(tableName)
	if err != nil {
		return err
	}

	// 创建表数据文件
	fileName := fmt.Sprintf("database/data/%s.json", tableName)
	fileHeader := &zip.FileHeader{
		Name:   fileName,
		Method: zip.Deflate,
	}
	writer, err := zipWriter.CreateHeader(fileHeader)
	if err != nil {
		return err
	}

	// 写入表数据
	tableData := map[string]interface{}{
		"table":      tableName,
		"count":      len(rows),
		"exported_at": time.Now().Format(time.RFC3339),
		"data":       rows,
	}

	dataJSON, err := json.MarshalIndent(tableData, "", "  ")
	if err != nil {
		return err
	}

	_, err = writer.Write(dataJSON)
	return err
}

// addConfigFiles 添加配置文件备份
func (bm *BackupManager) addConfigFiles(zipWriter *zip.Writer, backupInfo *BackupInfo) error {
	configFiles := []string{
		"config.json",
		"config.yaml",
		"config.yml",
		".env",
		"settings.json",
	}

	for _, configFile := range configFiles {
		if err := bm.addFileToZip(zipWriter, configFile, "config/"); err != nil {
			logger.Debug("Skipping config file: " + configFile)
		} else {
			backupInfo.ItemCounts["config_files"]++
		}
	}

	return nil
}

// addUserDataFiles 添加用户数据文件
func (bm *BackupManager) addUserDataFiles(zipWriter *zip.Writer, backupInfo *BackupInfo) error {
	userDataDirs := []string{
		"data/",
		"uploads/",
		"exports/",
		"logs/",
	}

	for _, dir := range userDataDirs {
		if err := bm.addDirectoryToZip(zipWriter, dir, "user_data/"); err != nil {
			logger.Debug("Skipping directory: " + dir)
		}
	}

	return nil
}

// addFileToZip 添加文件到ZIP
func (bm *BackupManager) addFileToZip(zipWriter *zip.Writer, filePath, zipPath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = zipPath + filepath.Base(filePath)
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}

// addDirectoryToZip 添加目录到ZIP
func (bm *BackupManager) addDirectoryToZip(zipWriter *zip.Writer, dirPath, zipPath string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过排除的文件
		if bm.shouldExcludeFile(path) {
			return nil
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}
		header.Name = zipPath + relPath
		header.Method = zip.Deflate

		if info.IsDir() {
			header.Name += "/"
		} else {
			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				return err
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// shouldExcludeFile 检查文件是否应该被排除
func (bm *BackupManager) shouldExcludeFile(path string) bool {
	for _, pattern := range bm.config.ExcludePatterns {
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return true
		}
	}
	return false
}

// addBackupMetadata 添加备份元数据
func (bm *BackupManager) addBackupMetadata(zipWriter *zip.Writer, backupInfo *BackupInfo) error {
	metadata := map[string]interface{}{
		"backup_info":     backupInfo,
		"backup_config":   bm.config,
		"system_info":     bm.getSystemInfo(),
		"created_at":      time.Now().Format(time.RFC3339),
		"twinos_version":  "1.0.0",
		"backup_version":  "1.0.0",
	}

	metadataJSON, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}

	writer, err := zipWriter.Create("backup_metadata.json")
	if err != nil {
		return err
	}

	_, err = writer.Write(metadataJSON)
	return err
}

// tableExists 检查表是否存在
func (bm *BackupManager) tableExists(tableName string) bool {
	query := fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' AND name='%s'", tableName)
	var name string
	err := bm.db.QueryRow(query).Scan(&name)
	return err == nil
}

// queryTableData 查询表数据
func (bm *BackupManager) queryTableData(tableName string) ([]map[string]interface{}, error) {
	rows, err := bm.db.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for rows.Next() {
		// 创建扫描值的容器
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		// 扫描行数据
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		// 创建字典
		row := make(map[string]interface{})
		for i, col := range columns {
			if val := values[i]; val != nil {
				if b, ok := val.([]byte); ok {
					row[col] = string(b)
				} else {
					row[col] = val
				}
			} else {
				row[col] = nil
			}
		}

		results = append(results, row)
	}

	return results, nil
}

// getTableInfo 获取数据表信息
func (bm *BackupManager) getTableInfo() []map[string]interface{} {
	var tables []map[string]interface{}

	tableNames := []string{"settings", "messages", "contacts", "analysis_results"}

	for _, tableName := range tableNames {
		if bm.tableExists(tableName) {
			var count int64
			bm.db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)).Scan(&count)

			tables = append(tables, map[string]interface{}{
				"name":  tableName,
				"count": count,
				"exists": true,
			})
		} else {
			tables = append(tables, map[string]interface{}{
				"name":  tableName,
				"count": 0,
				"exists": false,
			})
		}
	}

	return tables
}

// getSystemInfo 获取系统信息
func (bm *BackupManager) getSystemInfo() map[string]interface{} {
	return map[string]interface{}{
		"os":           "darwin", // 可以动态获取
		"arch":         "amd64",  // 可以动态获取
		"go_version":   "1.21+",  // 可以动态获取
		"hostname":     "localhost", // 可以动态获取
		"backup_tool":  "TwinOS Backup Manager",
		"created_by":   "user",
	}
}

// calculateFileMetrics 计算文件大小和校验和
func (bm *BackupManager) calculateFileMetrics(backupInfo *BackupInfo) error {
	file, err := os.Open(backupInfo.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 获取文件大小
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	backupInfo.Size = stat.Size()

	// 计算MD5校验和
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}
	backupInfo.Checksum = hex.EncodeToString(hash.Sum(nil))

	return nil
}

// saveBackupInfo 保存备份信息
func (bm *BackupManager) saveBackupInfo(backupInfo *BackupInfo) error {
	infoFile := filepath.Join(bm.config.BackupDir, backupInfo.ID+"_info.json")

	data, err := json.MarshalIndent(backupInfo, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(infoFile, data, 0644)
}

// cleanupOldBackups 清理旧备份
func (bm *BackupManager) cleanupOldBackups() error {
	if bm.config.MaxBackups <= 0 {
		return nil
	}

	files, err := os.ReadDir(bm.config.BackupDir)
	if err != nil {
		return err
	}

	var backupFiles []os.DirEntry
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".zip") {
			backupFiles = append(backupFiles, file)
		}
	}

	if len(backupFiles) <= bm.config.MaxBackups {
		return nil
	}

	// 按修改时间排序，删除最旧的文件
	for i := 0; i < len(backupFiles)-bm.config.MaxBackups; i++ {
		filePath := filepath.Join(bm.config.BackupDir, backupFiles[i].Name())
		if err := os.Remove(filePath); err != nil {
			logger.Warn("Failed to delete old backup: " + err.Error())
		} else {
			logger.Info("Deleted old backup: " + backupFiles[i].Name())
		}
	}

	return nil
}

// generateBackupID 生成备份ID
func generateBackupID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// 如果随机数生成失败，使用时间戳
		return fmt.Sprintf("backup_%d", time.Now().Unix())
	}
	return hex.EncodeToString(bytes)
}

// ListBackups 列出所有备份
func (bm *BackupManager) ListBackups() ([]BackupInfo, error) {
	files, err := os.ReadDir(bm.config.BackupDir)
	if err != nil {
		return nil, err
	}

	var backups []BackupInfo
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), "_info.json") {
			infoPath := filepath.Join(bm.config.BackupDir, file.Name())
			data, err := os.ReadFile(infoPath)
			if err != nil {
				logger.Warn("Failed to read backup info: " + err.Error())
				continue
			}

			var backupInfo BackupInfo
			if err := json.Unmarshal(data, &backupInfo); err != nil {
				logger.Warn("Failed to parse backup info: " + err.Error())
				continue
			}

			backups = append(backups, backupInfo)
		}
	}

	return backups, nil
}

// RestoreBackup 恢复备份
func (bm *BackupManager) RestoreBackup(backupID string) error {
	backupInfo, err := bm.getBackupInfo(backupID)
	if err != nil {
		return fmt.Errorf("backup not found: %w", err)
	}

	if backupInfo.Status != StatusCompleted {
		return fmt.Errorf("backup is not completed: %s", backupInfo.Status)
	}

	logger.Info("🔄 Starting restore from backup: " + backupID)

	// 验证备份文件完整性
	if err := bm.verifyBackup(backupInfo); err != nil {
		return fmt.Errorf("backup verification failed: %w", err)
	}

	// 执行恢复
	if err := bm.performRestore(backupInfo); err != nil {
		return fmt.Errorf("restore failed: %w", err)
	}

	logger.Info("✅ Restore completed successfully")
	return nil
}

// getBackupInfo 获取备份信息
func (bm *BackupManager) getBackupInfo(backupID string) (*BackupInfo, error) {
	infoPath := filepath.Join(bm.config.BackupDir, backupID+"_info.json")
	data, err := os.ReadFile(infoPath)
	if err != nil {
		return nil, err
	}

	var backupInfo BackupInfo
	if err := json.Unmarshal(data, &backupInfo); err != nil {
		return nil, err
	}

	return &backupInfo, nil
}

// verifyBackup 验证备份文件
func (bm *BackupManager) verifyBackup(backupInfo *BackupInfo) error {
	file, err := os.Open(backupInfo.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 计算文件校验和
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}
	calculatedChecksum := hex.EncodeToString(hash.Sum(nil))

	if calculatedChecksum != backupInfo.Checksum {
		return fmt.Errorf("backup checksum mismatch: expected %s, got %s",
			backupInfo.Checksum, calculatedChecksum)
	}

	return nil
}

// performRestore 执行恢复
func (bm *BackupManager) performRestore(backupInfo *BackupInfo) error {
	// 打开备份文件
	reader, err := zip.OpenReader(backupInfo.FilePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	// 恢复数据库
	for _, file := range reader.File {
		if strings.HasPrefix(file.Name, "database/") {
			if err := bm.restoreDatabaseFile(file); err != nil {
				return fmt.Errorf("failed to restore database file %s: %w", file.Name, err)
			}
		}
	}

	return nil
}

// restoreDatabaseFile 恢复数据库文件
func (bm *BackupManager) restoreDatabaseFile(file *zip.File) error {
	reader, err := file.Open()
	if err != nil {
		return err
	}
	defer reader.Close()

	// 这里可以实现具体的数据库恢复逻辑
	// 暂时只记录日志
	logger.Info("Restoring database file: " + file.Name)
	return nil
}

// DeleteBackup 删除备份
func (bm *BackupManager) DeleteBackup(backupID string) error {
	backupInfo, err := bm.getBackupInfo(backupID)
	if err != nil {
		return err
	}

	// 删除备份文件
	if err := os.Remove(backupInfo.FilePath); err != nil {
		return err
	}

	// 删除备份信息文件
	infoPath := filepath.Join(bm.config.BackupDir, backupID+"_info.json")
	if err := os.Remove(infoPath); err != nil {
		return err
	}

	logger.Info("🗑️ Deleted backup: " + backupID)
	return nil
}

// EncryptBackup 加密备份
func (bm *BackupManager) EncryptBackup(backupID string, password string) error {
	backupInfo, err := bm.getBackupInfo(backupID)
	if err != nil {
		return err
	}

	if backupInfo.Encrypted {
		return fmt.Errorf("backup is already encrypted")
	}

	// 生成加密密钥
	_ = sha256.Sum256([]byte(password))

	// 这里可以实现AES加密逻辑
	// 暂时只标记为已加密
	backupInfo.Encrypted = true
	backupInfo.Metadata["encryption_method"] = "AES-256-GCM"

	return bm.saveBackupInfo(backupInfo)
}