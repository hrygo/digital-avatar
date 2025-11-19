package handler

import (
	"twin-os/backend/internal/service"
)

// Handler 处理器聚合
type Handler struct {
	WeChat           *WeChatHandler
	Analysis         *AnalysisHandler
	EnhancedAnalysis *EnhancedAnalysisHandler
	Data             *DataHandler
	Settings         *SettingsHandler
	Backup           *BackupHandler
}

// New 创建处理器聚合
func New(services *service.Service) *Handler {
	return &Handler{
		WeChat:           NewWeChatHandler(services.WeChat),
		Analysis:         NewAnalysisHandler(services.Analysis),
		EnhancedAnalysis: NewEnhancedAnalysisHandler(services.EnhancedAnalysis),
		Data:             NewDataHandler(services.Data),
		Settings:         NewSettingsHandler(services.Settings),
		Backup:           NewBackupHandler(services.Backup),
	}
}