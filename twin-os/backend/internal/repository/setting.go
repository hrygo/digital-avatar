package repository

import (
	"database/sql"
	"fmt"
	"time"
)

// SettingRepository 设置仓储
type SettingRepository struct {
	db *sql.DB
}

// NewSettingRepository 创建设置仓储
func NewSettingRepository(db *sql.DB) *SettingRepository {
	return &SettingRepository{db: db}
}

// Get 获取设置值
func (r *SettingRepository) Get(key string) (string, error) {
	query := `
		SELECT value FROM system_settings WHERE key = ?
	`

	var value string
	err := r.db.QueryRow(query, key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("setting not found")
		}
		return "", fmt.Errorf("failed to get setting: %w", err)
	}

	return value, nil
}

// Set 设置配置值
func (r *SettingRepository) Set(key, value, description string) error {
	query := `
		INSERT OR REPLACE INTO system_settings (
			key, value, description, updated_at
		) VALUES (?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, key, value, description, time.Now())
	if err != nil {
		return fmt.Errorf("failed to set setting: %w", err)
	}

	return nil
}

// GetAll 获取所有设置
func (r *SettingRepository) GetAll() (map[string]string, error) {
	query := `
		SELECT key, value FROM system_settings
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}
	defer rows.Close()

	settings := make(map[string]string)
	for rows.Next() {
		var key, value string
		err := rows.Scan(&key, &value)
		if err != nil {
			return nil, fmt.Errorf("failed to scan setting: %w", err)
		}
		settings[key] = value
	}

	return settings, nil
}