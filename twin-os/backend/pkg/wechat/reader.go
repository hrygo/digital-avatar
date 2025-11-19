package wechat

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"twin-os/backend/internal/model"
	"twin-os/backend/pkg/logger"
)

// Reader 微信数据库读取器
type Reader struct {
	dbPath   string
	db       *sql.DB
	isOpened bool
}

// NewReader 创建微信数据库读取器
func NewReader(dbPath string) *Reader {
	return &Reader{
		dbPath: dbPath,
	}
}

// Open 打开微信数据库
func (r *Reader) Open() error {
	if r.isOpened {
		return nil
	}

	// 检查数据库文件是否存在
	if _, err := os.Stat(r.dbPath); os.IsNotExist(err) {
		return fmt.Errorf("wechat database file not found: %s", r.dbPath)
	}

	// 打开数据库
	db, err := sql.Open("sqlite3", r.dbPath)
	if err != nil {
		return fmt.Errorf("failed to open wechat database: %w", err)
	}

	// 测试连接
	if err = db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("failed to ping wechat database: %w", err)
	}

	r.db = db
	r.isOpened = true

	logger.Info("✅ WeChat database opened successfully")
	return nil
}

// IsOpened 检查数据库是否已打开
func (r *Reader) IsOpened() bool {
	return r.isOpened
}

// Close 关闭数据库连接
func (r *Reader) Close() error {
	if r.db != nil && r.isOpened {
		err := r.db.Close()
		r.isOpened = false
		return err
	}
	return nil
}

// GetMessages 获取消息列表
func (r *Reader) GetMessages(limit int, offset int) ([]*model.Message, error) {
	if !r.isOpened {
		return nil, fmt.Errorf("database not opened")
	}

	query := `
		SELECT
			msgId,
			talkerId,
			type,
			content,
			createTime
		FROM message
		ORDER BY createTime DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		msg := &model.Message{}
		var createTime int64

		err := rows.Scan(
			&msg.MessageID,
			&msg.TalkerID,
			&msg.Type,
			&msg.Content,
			&createTime,
		)
		if err != nil {
			logger.Warn("⚠️ Failed to scan message row: " + err.Error())
			continue
		}

		// 转换时间戳
		msg.Timestamp = time.Unix(createTime, 0)
		msg.CreatedAt = time.Now()
		msg.UpdatedAt = time.Now()

		messages = append(messages, msg)
	}

	return messages, nil
}

// GetMessagesByTimeRange 按时间范围获取消息
func (r *Reader) GetMessagesByTimeRange(startTime, endTime time.Time) ([]*model.Message, error) {
	if !r.isOpened {
		return nil, fmt.Errorf("database not opened")
	}

	query := `
		SELECT
			msgId,
			talkerId,
			type,
			content,
			createTime
		FROM message
		WHERE createTime >= ? AND createTime <= ?
		ORDER BY createTime DESC
	`

	startTs := startTime.Unix()
	endTs := endTime.Unix()

	rows, err := r.db.Query(query, startTs, endTs)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages by time range: %w", err)
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		msg := &model.Message{}
		var createTime int64

		err := rows.Scan(
			&msg.MessageID,
			&msg.TalkerID,
			&msg.Type,
			&msg.Content,
			&createTime,
		)
		if err != nil {
			logger.Warn("⚠️ Failed to scan message row: " + err.Error())
			continue
		}

		msg.Timestamp = time.Unix(createTime, 0)
		msg.CreatedAt = time.Now()
		msg.UpdatedAt = time.Now()

		messages = append(messages, msg)
	}

	return messages, nil
}

// GetContacts 获取联系人列表
func (r *Reader) GetContacts() ([]*model.Contact, error) {
	if !r.isOpened {
		return nil, fmt.Errorf("database not opened")
	}

	query := `
		SELECT
			username,
			nickname,
			remark,
			type
		FROM rcontact
		WHERE type & 1 = 1 AND type != 4
		ORDER BY nickname
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query contacts: %w", err)
	}
	defer rows.Close()

	var contacts []*model.Contact
	for rows.Next() {
		contact := &model.Contact{}

		err := rows.Scan(
			&contact.UserName,
			&contact.NickName,
			&contact.Remark,
			&contact.Type,
		)
		if err != nil {
			logger.Warn("⚠️ Failed to scan contact row: " + err.Error())
			continue
		}

		contact.CreatedAt = time.Now()
		contact.UpdatedAt = time.Now()

		// 如果有备注，优先显示备注
		if contact.Remark != "" {
			contact.NickName = contact.Remark
		}

		contacts = append(contacts, contact)
	}

	return contacts, nil
}

// GetChatRooms 获取群聊列表
func (r *Reader) GetChatRooms() ([]*model.Contact, error) {
	if !r.isOpened {
		return nil, fmt.Errorf("database not opened")
	}

	query := `
		SELECT
			username,
			nickname,
			remark,
			type
		FROM rcontact
		WHERE type & 2 = 2
		ORDER BY nickname
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query chat rooms: %w", err)
	}
	defer rows.Close()

	var chatRooms []*model.Contact
	for rows.Next() {
		contact := &model.Contact{}

		err := rows.Scan(
			&contact.UserName,
			&contact.NickName,
			&contact.Remark,
			&contact.Type,
		)
		if err != nil {
			logger.Warn("⚠️ Failed to scan chat room row: " + err.Error())
			continue
		}

		contact.CreatedAt = time.Now()
		contact.UpdatedAt = time.Now()

		chatRooms = append(chatRooms, contact)
	}

	return chatRooms, nil
}

// FindWeChatDB 查找微信数据库路径
func FindWeChatDB() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	// 常见的微信数据路径
	possiblePaths := []string{
		filepath.Join(homeDir, "Library/Containers/com.tencent.xin/Data/Library/Application Support/com.tencent.xin/*/Message/MessageTemp/*/DB/EnMicroMsg.db"),
		filepath.Join(homeDir, "Documents/WeChat Files/*/Msg/*/EnMicroMsg.db"),
	}

	for _, pattern := range possiblePaths {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}

		if len(matches) > 0 {
			// 返回最新找到的数据库
			return matches[0], nil
		}
	}

	return "", fmt.Errorf("wechat database not found in common locations")
}

// IsValidWeChatDB 检查是否为有效的微信数据库
func IsValidWeChatDB(dbPath string) bool {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return false
	}
	defer db.Close()

	// 检查必要的表是否存在
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='message'").Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

// ExtractUserInfo 从talkerId中提取用户信息
func ExtractUserInfo(talkerID string) (isGroup bool, userID string) {
	// 群聊ID通常以@@开头
	if strings.HasPrefix(talkerID, "@@") {
		return true, talkerID[2:]
	}
	return false, talkerID
}

// CleanMessageContent 清理消息内容
func CleanMessageContent(content string) string {
	if content == "" {
		return ""
	}

	// 移除XML标签
	xmlTagPattern := regexp.MustCompile(`<[^>]+>`)
	content = xmlTagPattern.ReplaceAllString(content, "")

	// 移除多余的空白字符
	spacePattern := regexp.MustCompile(`\s+`)
	content = spacePattern.ReplaceAllString(content, " ")

	// 去除首尾空白
	content = strings.TrimSpace(content)

	return content
}