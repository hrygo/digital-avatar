package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"twin-os/backend/internal/model"
)

// MessageRepository 消息仓储
type MessageRepository struct {
	db *sql.DB
}

// NewMessageRepository 创建消息仓储
func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

// CreateOrUpdate 创建或更新消息
func (r *MessageRepository) CreateOrUpdate(message *model.Message) error {
	query := `
		INSERT OR REPLACE INTO messages (
			message_id, talker_id, type, content, timestamp, updated_at
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		message.MessageID,
		message.TalkerID,
		message.Type,
		message.Content,
		message.Timestamp.Unix(),
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to create/update message: %w", err)
	}

	return nil
}

// GetByID 根据ID获取消息
func (r *MessageRepository) GetByID(id int64) (*model.Message, error) {
	query := `
		SELECT id, message_id, talker_id, type, content, timestamp, created_at, updated_at
		FROM messages
		WHERE id = ?
	`

	row := r.db.QueryRow(query, id)
	message := &model.Message{}

		err := row.Scan(
		&message.ID,
		&message.MessageID,
		&message.TalkerID,
		&message.Type,
		&message.Content,
		&message.Timestamp,
		&message.CreatedAt,
		&message.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("message not found")
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

		return message, nil
}

// GetByMessageID 根据消息ID获取消息
func (r *MessageRepository) GetByMessageID(messageID string) (*model.Message, error) {
	query := `
		SELECT id, message_id, talker_id, type, content, timestamp, created_at, updated_at
		FROM messages
		WHERE message_id = ?
	`

	row := r.db.QueryRow(query, messageID)
	message := &model.Message{}

	err := row.Scan(
		&message.ID,
		&message.MessageID,
		&message.TalkerID,
		&message.Type,
		&message.Content,
		&message.Timestamp,
		&message.CreatedAt,
		&message.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("message not found")
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

		return message, nil
}

// GetList 获取消息列表
func (r *MessageRepository) GetList(limit, offset int) ([]*model.Message, error) {
	query := `
		SELECT id, message_id, talker_id, type, content, timestamp, created_at, updated_at
		FROM messages
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		message := &model.Message{}
		
		err := rows.Scan(
			&message.ID,
			&message.MessageID,
			&message.TalkerID,
			&message.Type,
			&message.Content,
			&message.Timestamp,
			&message.CreatedAt,
			&message.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

				messages = append(messages, message)
	}

	return messages, nil
}

// GetByTalker 根据聊天对象获取消息
func (r *MessageRepository) GetByTalker(talkerID string, limit int) ([]*model.Message, error) {
	query := `
		SELECT id, message_id, talker_id, type, content, timestamp, created_at, updated_at
		FROM messages
		WHERE talker_id = ?
		ORDER BY timestamp DESC
		LIMIT ?
	`

	rows, err := r.db.Query(query, talkerID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages by talker: %w", err)
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		message := &model.Message{}
		
		err := rows.Scan(
			&message.ID,
			&message.MessageID,
			&message.TalkerID,
			&message.Type,
			&message.Content,
			&message.Timestamp,
			&message.CreatedAt,
			&message.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

				messages = append(messages, message)
	}

	return messages, nil
}

// GetRecentMessages 获取最近的消息
func (r *MessageRepository) GetRecentMessages(duration time.Duration, limit int) ([]*model.Message, error) {
	since := time.Now().Add(-duration).Unix()

	query := `
		SELECT id, message_id, talker_id, type, content, timestamp, created_at, updated_at
		FROM messages
		WHERE timestamp >= ?
		AND content != ''
		ORDER BY timestamp DESC
		LIMIT ?
	`

	rows, err := r.db.Query(query, since, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent messages: %w", err)
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		message := &model.Message{}
		
		err := rows.Scan(
			&message.ID,
			&message.MessageID,
			&message.TalkerID,
			&message.Type,
			&message.Content,
			&message.Timestamp,
			&message.CreatedAt,
			&message.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

				messages = append(messages, message)
	}

	return messages, nil
}

// Search 搜索消息
func (r *MessageRepository) Search(query string, limit int) ([]*model.Message, error) {
	searchQuery := "%" + query + "%"

	sql := `
		SELECT id, message_id, talker_id, type, content, timestamp, created_at, updated_at
		FROM messages
		WHERE content LIKE ?
		ORDER BY timestamp DESC
		LIMIT ?
	`

	rows, err := r.db.Query(sql, searchQuery, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search messages: %w", err)
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		message := &model.Message{}
		
		err := rows.Scan(
			&message.ID,
			&message.MessageID,
			&message.TalkerID,
			&message.Type,
			&message.Content,
			&message.Timestamp,
			&message.CreatedAt,
			&message.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

				messages = append(messages, message)
	}

	return messages, nil
}

// GetTotalCount 获取总消息数
func (r *MessageRepository) GetTotalCount() (int64, error) {
	query := "SELECT COUNT(*) FROM messages"

	var count int64
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get message count: %w", err)
	}

	return count, nil
}

// GetCountByTalker 获取指定聊天对象的消息数
func (r *MessageRepository) GetCountByTalker(talkerID string) (int64, error) {
	query := "SELECT COUNT(*) FROM messages WHERE talker_id = ?"

	var count int64
	err := r.db.QueryRow(query, talkerID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get message count by talker: %w", err)
	}

	return count, nil
}

// Delete 删除消息
func (r *MessageRepository) Delete(id int64) error {
	query := "DELETE FROM messages WHERE id = ?"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("message not found")
	}

	return nil
}

// DeleteBefore 删除指定时间之前的消息
func (r *MessageRepository) DeleteBefore(before time.Time) (int64, error) {
	query := "DELETE FROM messages WHERE timestamp < ?"

	result, err := r.db.Exec(query, before.Unix())
	if err != nil {
		return 0, fmt.Errorf("failed to delete old messages: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

// GetTalkers 获取所有聊天对象
func (r *MessageRepository) GetTalkers() ([]string, error) {
	query := `
		SELECT DISTINCT talker_id
		FROM messages
		ORDER BY talker_id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get talkers: %w", err)
	}
	defer rows.Close()

	var talkers []string
	for rows.Next() {
		var talkerID string
		err := rows.Scan(&talkerID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan talker: %w", err)
		}
		talkers = append(talkers, talkerID)
	}

	return talkers, nil
}

// GetMessagesByIDs 根据消息ID列表获取消息
func (r *MessageRepository) GetMessagesByIDs(messageIDs []string) ([]*model.Message, error) {
	if len(messageIDs) == 0 {
		return []*model.Message{}, nil
	}

	// 构建IN查询
	placeholders := make([]string, len(messageIDs))
	args := make([]interface{}, len(messageIDs))
	for i, id := range messageIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT id, message_id, talker_id, type, content, timestamp, created_at, updated_at
		FROM messages
		WHERE message_id IN (%s)
		ORDER BY timestamp DESC
	`, strings.Join(placeholders, ","))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages by IDs: %w", err)
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		message := &model.Message{}
		
		err := rows.Scan(
			&message.ID,
			&message.MessageID,
			&message.TalkerID,
			&message.Type,
			&message.Content,
			&message.Timestamp,
			&message.CreatedAt,
			&message.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

				messages = append(messages, message)
	}

	return messages, nil
}