package handler

import (
	"twin-os/backend/internal/services"
)

// Handler 处理器聚合
type Handler struct {
	WeChat           *WeChatHandlerV2
	Analysis         *AnalysisHandler
	EnhancedAnalysis *EnhancedAnalysisHandler
	Data             *DataHandler
	Settings         *SettingsHandler
	Backup           *BackupHandler
	RealtimeSync     *RealtimeSyncHandler
	AutoSync         *services.AutoSyncManager
}

// New 创建处理器聚合
func New(services *services.Service) *Handler {
	return &Handler{
		WeChat:           NewWeChatHandlerV2(services.WeChat),
		Analysis:         NewAnalysisHandler(services.Analysis),
		EnhancedAnalysis: NewEnhancedAnalysisHandler(services.EnhancedAnalysis),
		Data:             NewDataHandler(services.Data),
		Settings:         NewSettingsHandler(services.Settings),
		Backup:           NewBackupHandler(services.Backup),
		RealtimeSync:     NewRealtimeSyncHandler(services.RealtimeSync),
		AutoSync:         services.AutoSync,
	}
}