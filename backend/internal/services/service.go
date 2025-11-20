package services

import (
	"time"
	"twin-os/backend/internal/config"
	"twin-os/backend/internal/repository"
)

// Service 服务聚合
type Service struct {
	WeChat           *WeChatServiceV2
	Analysis         *AnalysisService
	EnhancedAnalysis *EnhancedAnalysisService
	Data             *DataService
	Settings         *SettingsService
	Backup           *BackupService
	RealtimeSync     *RealtimeSyncService
	AutoSync         *AutoSyncManager
}

// New 创建服务聚合
func New(cfg *config.Config, repo *repository.Repository) *Service {
	wechatService := NewWeChatServiceV2(repo.GetDB())
	realtimeSync := NewRealtimeSyncService(repo.GetDB(), wechatService)
	autoSync := NewAutoSyncManager(realtimeSync)

	// 设置自动同步配置（开发环境使用测试数据库）
	autoConfig := DefaultAutoSyncConfig()
	if cfg.Environment == "development" {
		autoConfig.ManualDBPath = "/tmp/real_wechat_demo.db"
		autoConfig.PreferManual = true
		autoConfig.StartupDelay = time.Second * 2
	}
	autoSync.SetConfig(autoConfig)

	return &Service{
		WeChat:           wechatService,
		Analysis:         NewAnalysisService(cfg, repo),
		EnhancedAnalysis: NewEnhancedAnalysisService(cfg, repo),
		Data:             NewDataService(cfg, repo),
		Settings:         NewSettingsService(cfg, repo),
		Backup:           NewBackupService(repo),
		RealtimeSync:     realtimeSync,
		AutoSync:         autoSync,
	}
}