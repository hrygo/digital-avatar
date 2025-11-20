package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"twin-os/backend/pkg/backup"
	"twin-os/backend/pkg/logger"
	"twin-os/backend/internal/services"
)

// BackupHandler 备份处理器
type BackupHandler struct {
	backupService *services.BackupService
}

// NewBackupHandler 创建备份处理器
func NewBackupHandler(backupService *services.BackupService) *BackupHandler {
	return &BackupHandler{
		backupService: backupService,
	}
}

// CreateBackup 创建备份
// @Summary 创建备份
// @Description 创建新的数据备份
// @Tags backup
// @Accept json
// @Produce json
// @Param request body services.CreateBackupRequest true "创建备份请求"
// @Success 200 {object} services.CreateBackupResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/backup/create [post]
func (h *BackupHandler) CreateBackup(c *gin.Context) {
	var req services.CreateBackupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("❌ Invalid backup request: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// 验证备份类型
	validTypes := []backup.BackupType{
		backup.BackupTypeFull,
		backup.BackupTypeConfig,
		backup.BackupTypeData,
		backup.BackupTypeIncremental,
	}

	isValidType := false
	for _, validType := range validTypes {
		if req.Type == validType {
			isValidType = true
			break
		}
	}

	if !isValidType {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid backup type",
			"valid_types": []string{
				string(backup.BackupTypeFull),
				string(backup.BackupTypeConfig),
				string(backup.BackupTypeData),
				string(backup.BackupTypeIncremental),
			},
		})
		return
	}

	// 设置默认描述
	if req.Description == "" {
		req.Description = "Backup created at " + c.GetHeader("User-Agent")
	}

	response, err := h.backupService.CreateBackup(&req)
	if err != nil {
		logger.Error("❌ Failed to create backup: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create backup",
			"details": err.Error(),
		})
		return
	}

	logger.Info("✅ Backup created successfully: " + response.BackupID)
	c.JSON(http.StatusOK, response)
}

// ListBackups 列出所有备份
// @Summary 列出备份
// @Description 获取所有备份的列表
// @Tags backup
// @Produce json
// @Param type query string false "备份类型过滤 (full, config, data, incremental)"
// @Param status query string false "状态过滤 (pending, in_progress, completed, failed, corrupted)"
// @Param limit query int false "限制返回数量" default(50)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/backup/list [get]
func (h *BackupHandler) ListBackups(c *gin.Context) {
	backups, err := h.backupService.ListBackups()
	if err != nil {
		logger.Error("❌ Failed to list backups: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list backups",
			"details": err.Error(),
		})
		return
	}

	// 应用过滤
	filteredBackups := h.filterBackups(backups, c.Query("type"), c.Query("status"))

	// 应用限制
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}

	if len(filteredBackups) > limit {
		filteredBackups = filteredBackups[:limit]
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"data":      filteredBackups,
		"total":     len(backups),
		"filtered":  len(filteredBackups),
		"timestamp": c.GetHeader("X-Request-Time"),
	})
}

// filterBackups 过滤备份列表
func (h *BackupHandler) filterBackups(backups []backup.BackupInfo, typeFilter, statusFilter string) []backup.BackupInfo {
	var filtered []backup.BackupInfo

	for _, backup := range backups {
		// 类型过滤
		if typeFilter != "" && string(backup.Type) != typeFilter {
			continue
		}

		// 状态过滤
		if statusFilter != "" && string(backup.Status) != statusFilter {
			continue
		}

		filtered = append(filtered, backup)
	}

	return filtered
}

// GetBackup 获取备份详情
// @Summary 获取备份详情
// @Description 根据备份ID获取备份的详细信息
// @Tags backup
// @Produce json
// @Param id path string true "备份ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/backup/{id} [get]
func (h *BackupHandler) GetBackup(c *gin.Context) {
	backupID := c.Param("id")
	if backupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Backup ID is required",
		})
		return
	}

	backupInfo, err := h.backupService.GetBackupInfo(backupID)
	if err != nil {
		logger.Error("❌ Failed to get backup info: " + err.Error())
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Backup not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"data":      backupInfo,
		"timestamp": c.GetHeader("X-Request-Time"),
	})
}

// RestoreBackup 恢复备份
// @Summary 恢复备份
// @Description 从指定备份恢复数据
// @Tags backup
// @Accept json
// @Produce json
// @Param id path string true "备份ID"
// @Param request body services.RestoreBackupRequest false "恢复选项"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/backup/{id}/restore [post]
func (h *BackupHandler) RestoreBackup(c *gin.Context) {
	backupID := c.Param("id")
	if backupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Backup ID is required",
		})
		return
	}

	var req services.RestoreBackupRequest
	req.BackupID = backupID

	// 可选地解析请求体以获取密码等选项
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Warn("⚠️ Failed to parse restore request, using default options: " + err.Error())
		}
	}

	// 添加确认检查
	confirmHeader := c.GetHeader("X-Confirm-Restore")
	if confirmHeader != "true" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":          "Confirmation required",
			"message":        "Please add 'X-Confirm-Restore: true' header to confirm restore operation",
			"warning":        "This operation will overwrite current data",
			"confirm_header": "X-Confirm-Restore: true",
		})
		return
	}

	logger.Info("🔄 Starting backup restore: " + backupID)
	if err := h.backupService.RestoreBackup(&req); err != nil {
		logger.Error("❌ Failed to restore backup: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to restore backup",
			"details": err.Error(),
		})
		return
	}

	logger.Info("✅ Backup restored successfully: " + backupID)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Backup restored successfully",
		"backup_id": backupID,
	})
}

// DeleteBackup 删除备份
// @Summary 删除备份
// @Description 删除指定的备份文件
// @Tags backup
// @Accept json
// @Produce json
// @Param id path string true "备份ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/backup/{id} [delete]
func (h *BackupHandler) DeleteBackup(c *gin.Context) {
	backupID := c.Param("id")
	if backupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Backup ID is required",
		})
		return
	}

	req := services.DeleteBackupRequest{BackupID: backupID}

	logger.Info("🗑️ Deleting backup: " + backupID)
	if err := h.backupService.DeleteBackup(&req); err != nil {
		logger.Error("❌ Failed to delete backup: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete backup",
			"details": err.Error(),
		})
		return
	}

	logger.Info("✅ Backup deleted successfully: " + backupID)
	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   "Backup deleted successfully",
		"backup_id": backupID,
	})
}

// GetBackupStats 获取备份统计
// @Summary 获取备份统计
// @Description 获取备份系统的统计信息
// @Tags backup
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/backup/stats [get]
func (h *BackupHandler) GetBackupStats(c *gin.Context) {
	stats, err := h.backupService.GetBackupStats()
	if err != nil {
		logger.Error("❌ Failed to get backup stats: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get backup statistics",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"data":      stats,
		"timestamp": c.GetHeader("X-Request-Time"),
	})
}

// ScheduleBackup 计划备份
// @Summary 计划备份
// @Description 设置自动备份计划
// @Tags backup
// @Accept json
// @Produce json
// @Param request body services.ScheduleBackupRequest true "备份计划配置"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/backup/schedule [post]
func (h *BackupHandler) ScheduleBackup(c *gin.Context) {
	var req services.ScheduleBackupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("❌ Invalid schedule request: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// 验证间隔格式
	if req.Interval != "" && !h.isValidInterval(req.Interval) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid interval format",
			"examples": []string{"1h", "24h", "1d", "1w", "1M"},
		})
		return
	}

	if err := h.backupService.ScheduleBackup(&req); err != nil {
		logger.Error("❌ Failed to schedule backup: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to schedule backup",
			"details": err.Error(),
		})
		return
	}

	status := "disabled"
	if req.Enabled {
		status = "enabled"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Backup schedule " + status,
		"config":  req,
	})
}

// isValidInterval 验证时间间隔格式
func (h *BackupHandler) isValidInterval(interval string) bool {
	validIntervals := []string{
		"1h", "2h", "6h", "12h", "24h", // 小时
		"1d", "7d", "30d",             // 天数
		"1w", "2w", "4w",               // 周数
		"1M", "3M", "6M", "12M",        // 月数
	}

	for _, valid := range validIntervals {
		if interval == valid {
			return true
		}
	}

	return false
}

// ExportBackup 导出备份
// @Summary 导出备份
// @Description 将备份导出为指定格式
// @Tags backup
// @Accept json
// @Produce json
// @Param id path string true "备份ID"
// @Param request body services.ExportBackupRequest true "导出配置"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/backup/{id}/export [post]
func (h *BackupHandler) ExportBackup(c *gin.Context) {
	backupID := c.Param("id")
	if backupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Backup ID is required",
		})
		return
	}

	var req services.ExportBackupRequest
	req.BackupID = backupID

	if err := c.ShouldBindJSON(&req); err != nil {
		req.Format = "zip" // 默认格式
	}

	// 验证导出格式
	validFormats := []string{"zip", "tar", "json"}
	isValidFormat := false
	for _, format := range validFormats {
		if req.Format == format {
			isValidFormat = true
			break
		}
	}

	if !isValidFormat {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":        "Invalid export format",
			"valid_formats": validFormats,
		})
		return
	}

	exportPath, err := h.backupService.ExportBackup(&req)
	if err != nil {
		logger.Error("❌ Failed to export backup: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to export backup",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "success",
		"message":     "Backup exported successfully",
		"backup_id":   backupID,
		"export_path": exportPath,
		"format":      req.Format,
	})
}

// ImportBackup 导入备份
// @Summary 导入备份
// @Description 从外部文件导入备份
// @Tags backup
// @Accept json
// @Produce json
// @Param request body services.ImportBackupRequest true "导入配置"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/backup/import [post]
func (h *BackupHandler) ImportBackup(c *gin.Context) {
	var req services.ImportBackupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("❌ Invalid import request: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	if req.FilePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File path is required",
		})
		return
	}

	// 安全检查：防止路径遍历攻击
	if strings.Contains(req.FilePath, "..") || strings.HasPrefix(req.FilePath, "/") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid file path",
		})
		return
	}

	backupInfo, err := h.backupService.ImportBackup(&req)
	if err != nil {
		logger.Error("❌ Failed to import backup: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to import backup",
			"details": err.Error(),
		})
		return
	}

	logger.Info("✅ Backup imported successfully: " + backupInfo.ID)
	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   "Backup imported successfully",
		"backup":    backupInfo,
	})
}

// VerifyBackup 验证备份
// @Summary 验证备份
// @Description 验证备份文件的完整性
// @Tags backup
// @Produce json
// @Param id path string true "备份ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/backup/{id}/verify [get]
func (h *BackupHandler) VerifyBackup(c *gin.Context) {
	backupID := c.Param("id")
	if backupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Backup ID is required",
		})
		return
	}

	if err := h.backupService.VerifyBackup(backupID); err != nil {
		logger.Error("❌ Backup verification failed: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"error":   "Backup verification failed",
			"details": err.Error(),
		})
		return
	}

	logger.Info("✅ Backup verification passed: " + backupID)
	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   "Backup verification passed",
		"backup_id": backupID,
	})
}