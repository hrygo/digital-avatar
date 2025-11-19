package service

import (
	"twin-os/backend/internal/config"
	"twin-os/backend/internal/repository"
)

// DataService 数据服务
type DataService struct {
	config     *config.Config
	repository *repository.Repository
}

// NewDataService 创建数据服务
func NewDataService(cfg *config.Config, repo *repository.Repository) *DataService {
	return &DataService{
		config:     cfg,
		repository: repo,
	}
}

// GetMessages 获取消息列表
func (s *DataService) GetMessages(limit, offset int) (interface{}, error) {
	return s.repository.Message.GetList(limit, offset)
}

// GetContacts 获取联系人列表
func (s *DataService) GetContacts() (interface{}, error) {
	return s.repository.Contact.GetAll()
}

// ExportData 导出数据
func (s *DataService) ExportData() (interface{}, error) {
	// 实现数据导出逻辑
	return map[string]interface{}{
		"status":  "success",
		"message": "Data export not implemented yet",
	}, nil
}