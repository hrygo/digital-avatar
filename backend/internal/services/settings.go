package services

import (
	"twin-os/backend/internal/config"
	"twin-os/backend/internal/repository"
)

// SettingsService 设置服务
type SettingsService struct {
	config     *config.Config
	repository *repository.Repository
}

// NewSettingsService 创建设置服务
func NewSettingsService(cfg *config.Config, repo *repository.Repository) *SettingsService {
	return &SettingsService{
		config:     cfg,
		repository: repo,
	}
}

// GetSettings 获取设置
func (s *SettingsService) GetSettings() (interface{}, error) {
	return s.repository.Setting.GetAll()
}

// UpdateSettings 更新设置
func (s *SettingsService) UpdateSettings(settings map[string]string) (interface{}, error) {
	for key, value := range settings {
		err := s.repository.Setting.Set(key, value, "")
		if err != nil {
			return nil, err
		}
	}

	return map[string]interface{}{
		"status":  "success",
		"message": "Settings updated successfully",
	}, nil
}