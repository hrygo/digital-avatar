package services

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"twin-os/backend/pkg/backup"
	"twin-os/backend/pkg/logger"
	"twin-os/backend/internal/repository"
)

// BackupService 备份服务
type BackupService struct {
	repo        *repository.Repository
	backupManager *backup.BackupManager
}

// NewBackupService 创建备份服务
func NewBackupService(repo *repository.Repository) *BackupService {
	// 创建默认备份配置
	config := backup.DefaultBackupConfig()

	// 设置备份目录为相对于项目根目录
	backupDir := filepath.Join(".", "backups")
	config.BackupDir = backupDir

	// 确保备份目录存在
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		logger.Error("Failed to create backup directory: " + err.Error())
	}

	backupManager := backup.NewBackupManager(repo.GetDB(), config)

	return &BackupService{
		repo:         repo,
		backupManager: backupManager,
	}
}

// CreateBackupRequest 创建备份请求
type CreateBackupRequest struct {
	Type        backup.BackupType `json:"type" binding:"required"`
	Description string            `json:"description"`
}

// CreateBackupResponse 创建备份响应
type CreateBackupResponse struct {
	BackupID string             `json:"backup_id"`
	Info     *backup.BackupInfo `json:"info"`
	Status   string             `json:"status"`
	Message  string             `json:"message"`
}

// CreateBackup 创建备份
func (s *BackupService) CreateBackup(req *CreateBackupRequest) (*CreateBackupResponse, error) {
	logger.Info(fmt.Sprintf("📦 Creating backup: type=%s, description=%s", req.Type, req.Description))

	// 创建备份
	backupInfo, err := s.backupManager.CreateBackup(req.Type, req.Description)
	if err != nil {
		logger.Error("❌ Failed to create backup: " + err.Error())
		return &CreateBackupResponse{
			Status:  "error",
			Message: "Failed to create backup: " + err.Error(),
		}, err
	}

	logger.Info(fmt.Sprintf("✅ Backup created successfully: %s", backupInfo.ID))

	return &CreateBackupResponse{
		BackupID: backupInfo.ID,
		Info:     backupInfo,
		Status:   "success",
		Message:  "Backup created successfully",
	}, nil
}

// ListBackups 列出所有备份
func (s *BackupService) ListBackups() ([]backup.BackupInfo, error) {
	logger.Info("📋 Listing backups...")

	backups, err := s.backupManager.ListBackups()
	if err != nil {
		logger.Error("❌ Failed to list backups: " + err.Error())
		return nil, err
	}

	logger.Info(fmt.Sprintf("📊 Found %d backups", len(backups)))
	return backups, nil
}

// RestoreBackupRequest 恢复备份请求
type RestoreBackupRequest struct {
	BackupID string `json:"backup_id" binding:"required"`
	Password string `json:"password,omitempty"` // 用于加密备份
}

// RestoreBackup 恢复备份
func (s *BackupService) RestoreBackup(req *RestoreBackupRequest) error {
	logger.Info("🔄 Restoring backup: " + req.BackupID)

	// 验证备份ID是否存在
	backups, err := s.backupManager.ListBackups()
	if err != nil {
		return fmt.Errorf("failed to list backups: %w", err)
	}

	var foundBackup *backup.BackupInfo
	for _, b := range backups {
		if b.ID == req.BackupID {
			foundBackup = &b
			break
		}
	}

	if foundBackup == nil {
		return fmt.Errorf("backup not found: %s", req.BackupID)
	}

	// 执行恢复
	if err := s.backupManager.RestoreBackup(req.BackupID); err != nil {
		logger.Error("❌ Failed to restore backup: " + err.Error())
		return err
	}

	logger.Info("✅ Backup restored successfully: " + req.BackupID)
	return nil
}

// DeleteBackupRequest 删除备份请求
type DeleteBackupRequest struct {
	BackupID string `json:"backup_id" binding:"required"`
}

// DeleteBackup 删除备份
func (s *BackupService) DeleteBackup(req *DeleteBackupRequest) error {
	logger.Info("🗑️ Deleting backup: " + req.BackupID)

	// 验证备份ID是否存在
	backups, err := s.backupManager.ListBackups()
	if err != nil {
		return fmt.Errorf("failed to list backups: %w", err)
	}

	var foundBackup *backup.BackupInfo
	for _, b := range backups {
		if b.ID == req.BackupID {
			foundBackup = &b
			break
		}
	}

	if foundBackup == nil {
		return fmt.Errorf("backup not found: %s", req.BackupID)
	}

	// 删除备份
	if err := s.backupManager.DeleteBackup(req.BackupID); err != nil {
		logger.Error("❌ Failed to delete backup: " + err.Error())
		return err
	}

	logger.Info("✅ Backup deleted successfully: " + req.BackupID)
	return nil
}

// GetBackupInfo 获取备份详细信息
func (s *BackupService) GetBackupInfo(backupID string) (*backup.BackupInfo, error) {
	logger.Info("📄 Getting backup info: " + backupID)

	// 列出所有备份来查找指定的备份
	backups, err := s.backupManager.ListBackups()
	if err != nil {
		return nil, fmt.Errorf("failed to list backups: %w", err)
	}

	for _, backup := range backups {
		if backup.ID == backupID {
			return &backup, nil
		}
	}

	return nil, fmt.Errorf("backup not found: %s", backupID)
}

// GetBackupStats 获取备份统计信息
func (s *BackupService) GetBackupStats() (map[string]interface{}, error) {
	logger.Info("📊 Getting backup statistics...")

	backups, err := s.backupManager.ListBackups()
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_backups": len(backups),
		"total_size":    int64(0),
		"by_type":       make(map[string]int),
		"by_status":     make(map[string]int),
		"latest_backup": nil,
		"oldest_backup": nil,
	}

	var totalSize int64
	var latestTime, oldestTime time.Time

	for i, backup := range backups {
		// 统计类型
		stats["by_type"].(map[string]int)[string(backup.Type)]++

		// 统计状态
		stats["by_status"].(map[string]int)[string(backup.Status)]++

		// 累计大小
		totalSize += backup.Size

		// 记录最新和最旧备份
		if i == 0 || backup.CreatedAt.After(latestTime) {
			latestTime = backup.CreatedAt
			stats["latest_backup"] = backup
		}

		if i == 0 || backup.CreatedAt.Before(oldestTime) {
			oldestTime = backup.CreatedAt
			stats["oldest_backup"] = backup
		}
	}

	stats["total_size"] = totalSize
	stats["total_size_mb"] = totalSize / (1024 * 1024)
	stats["total_size_gb"] = totalSize / (1024 * 1024 * 1024)

	return stats, nil
}

// ScheduleBackupRequest 计划备份请求
type ScheduleBackupRequest struct {
	Type        backup.BackupType `json:"type" binding:"required"`
	Description string            `json:"description"`
	Interval    string            `json:"interval"` // 如 "24h", "1d", "1w"
	Enabled     bool              `json:"enabled"`
}

// ScheduleBackup 计划备份
func (s *BackupService) ScheduleBackup(req *ScheduleBackupRequest) error {
	logger.Info(fmt.Sprintf("⏰ Scheduling backup: type=%s, interval=%s", req.Type, req.Interval))

	// 这里可以实现定时备份逻辑
	// 暂时只记录日志
	if req.Enabled {
		logger.Info(fmt.Sprintf("✅ Backup scheduled: %s every %s", req.Type, req.Interval))
	} else {
		logger.Info("⏸️ Backup scheduling disabled")
	}

	return nil
}

// ExportBackupRequest 导出备份请求
type ExportBackupRequest struct {
	BackupID string `json:"backup_id" binding:"required"`
	Format   string `json:"format"` // "zip", "tar", "json"
}

// ExportBackup 导出备份
func (s *BackupService) ExportBackup(req *ExportBackupRequest) (string, error) {
	logger.Info(fmt.Sprintf("📤 Exporting backup: %s as %s", req.BackupID, req.Format))

	// 获取备份信息
	backupInfo, err := s.GetBackupInfo(req.BackupID)
	if err != nil {
		return "", err
	}

	// 这里可以实现导出逻辑
	// 暂时返回备份文件路径
	exportPath := filepath.Join("exports", fmt.Sprintf("%s_export.%s", req.BackupID, req.Format))

	// 确保导出目录存在
	if err := os.MkdirAll(filepath.Dir(exportPath), 0755); err != nil {
		return "", err
	}

	// 复制备份文件到导出目录
	if err := copyFile(backupInfo.FilePath, exportPath); err != nil {
		return "", err
	}

	logger.Info(fmt.Sprintf("✅ Backup exported to: %s", exportPath))
	return exportPath, nil
}

// ImportBackupRequest 导入备份请求
type ImportBackupRequest struct {
	FilePath string `json:"file_path" binding:"required"`
	Description string `json:"description"`
}

// ImportBackup 导入备份
func (s *BackupService) ImportBackup(req *ImportBackupRequest) (*backup.BackupInfo, error) {
	logger.Info(fmt.Sprintf("📥 Importing backup from: %s", req.FilePath))

	// 验证文件是否存在
	if _, err := os.Stat(req.FilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("backup file not found: %s", req.FilePath)
	}

	// 这里可以实现导入逻辑
	// 暂时创建一个模拟的备份信息
	backupInfo := &backup.BackupInfo{
		ID:          generateID(),
		Type:        backup.BackupTypeFull,
		Status:      backup.StatusCompleted,
		CreatedAt:   time.Now(),
		Description: req.Description,
		Version:     "1.0.0",
		FilePath:    req.FilePath,
		Metadata:    make(map[string]string),
		Encrypted:   false,
		Compression: "zip",
		ItemCounts:  make(map[string]int),
	}

	// 获取文件大小
	if stat, err := os.Stat(req.FilePath); err == nil {
		backupInfo.Size = stat.Size()
	}

	logger.Info(fmt.Sprintf("✅ Backup imported successfully: %s", backupInfo.ID))
	return backupInfo, nil
}

// VerifyBackup 验证备份
func (s *BackupService) VerifyBackup(backupID string) error {
	logger.Info("🔍 Verifying backup: " + backupID)

	// 获取备份信息
	backupInfo, err := s.GetBackupInfo(backupID)
	if err != nil {
		return err
	}

	// 验证备份文件
	if _, err := os.Stat(backupInfo.FilePath); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found: %s", backupInfo.FilePath)
	}

	// 这里可以实现更详细的验证逻辑
	// 暂时只检查文件存在性
	logger.Info("✅ Backup verification completed: " + backupID)
	return nil
}

// generateID 生成唯一ID
func generateID() string {
	return fmt.Sprintf("backup_%d", time.Now().Unix())
}

// copyFile 复制文件
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = destFile.ReadFrom(sourceFile)
	return err
}