package service

import (
	"twin-os/backend/internal/config"
	"twin-os/backend/internal/repository"
)

// Service 服务聚合
type Service struct {
	WeChat           *WeChatService
	Analysis         *AnalysisService
	EnhancedAnalysis *EnhancedAnalysisService
	Data             *DataService
	Settings         *SettingsService
	Backup           *BackupService
}

// New 创建服务聚合
func New(cfg *config.Config, repo *repository.Repository) *Service {
	return &Service{
		WeChat:           NewWeChatService(cfg, repo),
		Analysis:         NewAnalysisService(cfg, repo),
		EnhancedAnalysis: NewEnhancedAnalysisService(cfg, repo),
		Data:             NewDataService(cfg, repo),
		Settings:         NewSettingsService(cfg, repo),
		Backup:           NewBackupService(repo),
	}
}