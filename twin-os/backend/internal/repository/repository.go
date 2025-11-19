package repository

import (
	"database/sql"

	"twin-os/backend/pkg/database"
)

// Repository 仓储聚合
type Repository struct {
	Message  *MessageRepository
	Contact  *ContactRepository
	Analysis *AnalysisRepository
	Setting  *SettingRepository
}

// New 创建仓储聚合
func New(db *sql.DB) *Repository {
	return &Repository{
		Message:  NewMessageRepository(db),
		Contact:  NewContactRepository(db),
		Analysis: NewAnalysisRepository(db),
		Setting:  NewSettingRepository(db),
	}
}

// Close 关闭所有仓储连接
func (r *Repository) Close() error {
	return database.Close()
}