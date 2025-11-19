package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"twin-os/backend/internal/model"
)

// ContactRepository 联系人仓储
type ContactRepository struct {
	db *sql.DB
}

// NewContactRepository 创建联系人仓储
func NewContactRepository(db *sql.DB) *ContactRepository {
	return &ContactRepository{db: db}
}

// CreateOrUpdate 创建或更新联系人 - v0.3.0性能优化
func (r *ContactRepository) CreateOrUpdate(ctx context.Context, contact *model.Contact) error {
	query := `
		INSERT OR REPLACE INTO contacts (
			user_name, nick_name, remark, type, updated_at
		) VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		contact.UserName,
		contact.NickName,
		contact.Remark,
		contact.Type,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to create/update contact: %w", err)
	}

	return nil
}

// GetAll 获取所有联系人
func (r *ContactRepository) GetAll() ([]*model.Contact, error) {
	query := `
		SELECT id, user_name, nick_name, remark, type, created_at, updated_at
		FROM contacts
		ORDER BY nick_name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get contacts: %w", err)
	}
	defer rows.Close()

	var contacts []*model.Contact
	for rows.Next() {
		contact := &model.Contact{}

		err := rows.Scan(
			&contact.ID,
			&contact.UserName,
			&contact.NickName,
			&contact.Remark,
			&contact.Type,
			&contact.CreatedAt,
			&contact.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan contact: %w", err)
		}

		contacts = append(contacts, contact)
	}

	return contacts, nil
}

// GetByUserName 根据用户名获取联系人
func (r *ContactRepository) GetByUserName(userName string) (*model.Contact, error) {
	query := `
		SELECT id, user_name, nick_name, remark, type, created_at, updated_at
		FROM contacts
		WHERE user_name = ?
	`

	row := r.db.QueryRow(query, userName)
	contact := &model.Contact{}

	err := row.Scan(
		&contact.ID,
		&contact.UserName,
		&contact.NickName,
		&contact.Remark,
		&contact.Type,
		&contact.CreatedAt,
		&contact.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("contact not found")
		}
		return nil, fmt.Errorf("failed to get contact: %w", err)
	}

	return contact, nil
}