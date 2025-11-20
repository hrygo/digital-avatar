package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hrygo/log"
	"twin-os/backend/internal/models"
)

// WeChatServiceV2 微信数据服务 - 使用原生SQL
type WeChatServiceV2 struct {
	db *sql.DB
}

// NewWeChatServiceV2 创建微信数据服务
func NewWeChatServiceV2(db *sql.DB) *WeChatServiceV2 {
	return &WeChatServiceV2{
		db: db,
	}
}

// ImportFromDatabase 从微信数据库导入数据
func (s *WeChatServiceV2) ImportFromDatabase(dbPath string) (*models.WeChatSyncRecord, error) {
	log.Infof("开始从微信数据库导入数据: %s", dbPath)

	// 创建同步记录
	syncRecord := &models.WeChatSyncRecord{
		SyncType:  "database",
		Source:    dbPath,
		Status:    "running",
		Progress:  0,
		StartTime: time.Now(),
	}

	// 插入同步记录
	recordID, err := s.insertSyncRecord(syncRecord)
	if err != nil {
		return nil, fmt.Errorf("创建同步记录失败: %w", err)
	}
	syncRecord.ID = int64(recordID)

	// 创建解析器
	parser, err := NewWeChatParser(dbPath)
	if err != nil {
		s.updateSyncRecord(recordID, "failed", 0, fmt.Sprintf("创建解析器失败: %v", err))
		return syncRecord, fmt.Errorf("创建微信解析器失败: %w", err)
	}
	defer parser.Close()

	// 解析数据库
	if err := parser.ParseDatabase(); err != nil {
		s.updateSyncRecord(recordID, "failed", 0, fmt.Sprintf("解析数据库失败: %v", err))
		return syncRecord, fmt.Errorf("解析微信数据库失败: %w", err)
	}

	// 更新进度
	s.updateSyncRecord(recordID, "running", 0.2, "")

	// 导入联系人
	contacts := parser.GetContacts()
	newContacts, updatedContacts, err := s.importContactsV2(contacts)
	if err != nil {
		s.updateSyncRecord(recordID, "failed", 0.5, fmt.Sprintf("导入联系人失败: %v", err))
		return syncRecord, fmt.Errorf("导入联系人失败: %w", err)
	}

	// 更新进度
	s.updateSyncRecord(recordID, "running", 0.5, "")

	// 导入消息
	messages := parser.GetMessages()
	newMessages, updatedMessages, err := s.importMessagesV2(messages)
	if err != nil {
		s.updateSyncRecord(recordID, "failed", 0.8, fmt.Sprintf("导入消息失败: %v", err))
		return syncRecord, fmt.Errorf("导入消息失败: %w", err)
	}

	// 更新进度
	s.updateSyncRecord(recordID, "running", 0.9, "")

	// 导入聊天会话
	chats := parser.GetChats()
	if err := s.importChatsV2(chats); err != nil {
		s.updateSyncRecord(recordID, "failed", 0.9, fmt.Sprintf("导入聊天会话失败: %v", err))
		return syncRecord, fmt.Errorf("导入聊天会话失败: %w", err)
	}

	// 完成同步
	endTime := time.Now()
	duration := endTime.Sub(syncRecord.StartTime).Seconds()

	s.updateSyncRecord(recordID, "completed", 1.0, "")

	// 更新同步记录完成状态
	updateSQL := `
		UPDATE wechat_sync_records
		SET status = ?, progress = ?, total_contacts = ?, new_contacts = ?,
		    updated_contacts = ?, total_messages = ?, new_messages = ?,
		    updated_messages = ?, end_time = ?, duration = ?, error_message = ?
		WHERE id = ?
	`
	_, err = s.db.Exec(updateSQL, "completed", 1.0, len(contacts), newContacts,
		updatedContacts, len(messages), newMessages, updatedMessages,
		endTime, int64(duration), "", recordID)

	if err != nil {
		log.Errorf("更新同步记录失败: %v", err)
	}

	log.Infof("✅ 微信数据导入完成: 联系人 %d 个 (%d 新增, %d 更新), 消息 %d 条 (%d 新增, %d 更新), 会话 %d 个, 耗时 %.2f 秒",
		len(contacts), newContacts, updatedContacts,
		len(messages), newMessages, updatedMessages,
		len(chats), duration)

	return syncRecord, nil
}

// GetWeChatStatus 获取微信数据状态
func (s *WeChatServiceV2) GetWeChatStatus() (*WeChatStatus, error) {
	var status WeChatStatus

	// 统计消息数量
	err := s.db.QueryRow("SELECT COUNT(*) FROM wechat_messages").Scan(&status.MessageCount)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("统计消息数量失败: %w", err)
	}

	// 统计联系人数量
	err = s.db.QueryRow("SELECT COUNT(*) FROM wechat_contacts").Scan(&status.ContactCount)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("统计联系人数量失败: %w", err)
	}

	// 统计会话数量
	var chatCount int64
	err = s.db.QueryRow("SELECT COUNT(*) FROM wechat_chats").Scan(&chatCount)
	status.ContactCount = int(chatCount) // 临时使用ContactCount字段存储聊天数量
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("统计会话数量失败: %w", err)
	}

	// 获取最后同步记录
	var lastSyncTime sql.NullInt64
	err = s.db.QueryRow("SELECT start_time FROM wechat_sync_records ORDER BY created_at DESC LIMIT 1").
		Scan(&lastSyncTime)
	if err == nil {
		status.LastSync = time.Unix(lastSyncTime.Int64, 0)
	}

	// 检查是否有微信数据
	status.IsConnected = status.MessageCount != nil && *status.MessageCount > 0

	return &status, nil
}

// GetMessages 获取消息列表
func (s *WeChatServiceV2) GetMessages(page, pageSize int, chatType, talker string) ([]*models.WeChatMessage, int64, error) {
	var messages []*models.WeChatMessage
	var total int64

	// 构建查询条件
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if chatType != "" {
		whereClause += fmt.Sprintf(" AND chat_type = $%d", argIndex)
		args = append(args, chatType)
		argIndex++
	}
	if talker != "" {
		whereClause += fmt.Sprintf(" AND talker = $%d", argIndex)
		args = append(args, talker)
		argIndex++
	}

	// 统计总数
	countSQL := "SELECT COUNT(*) FROM wechat_messages " + whereClause
	err := s.db.QueryRow(countSQL, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("统计消息总数失败: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	querySQL := fmt.Sprintf(`
		SELECT id, svr_id, create_time, talker, type, sub_type, is_sender, seq, flag, status,
		       content, display_content, chat_type, group_id, mentions, reply_to, created_at, updated_at
		FROM wechat_messages %s
		ORDER BY create_time DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIndex, argIndex+1)

	args = append(args, pageSize, offset)
	rows, err := s.db.Query(querySQL, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("查询消息失败: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		msg := &models.WeChatMessage{}
		var mentions sql.NullString
		var replyTo sql.NullString

		err := rows.Scan(&msg.ID, &msg.SvrID, &msg.CreateTime, &msg.Talker, &msg.Type,
			&msg.SubType, &msg.IsSender, &msg.Seq, &msg.Flag, &msg.Status,
			&msg.Content, &msg.DisplayContent, &msg.ChatType, &msg.GroupID,
			&mentions, &replyTo, &msg.CreatedAt, &msg.UpdatedAt)
		if err != nil {
			log.Warnf("扫描消息数据失败: %v", err)
			continue
		}

		if mentions.Valid {
			msg.Mentions = models.JSON{}
			json.Unmarshal([]byte(mentions.String), &msg.Mentions)
		}
		if replyTo.Valid {
			msg.ReplyTo = replyTo.String
		}

		messages = append(messages, msg)
	}

	return messages, total, nil
}

// GetContacts 获取联系人列表
func (s *WeChatServiceV2) GetContacts(page, pageSize int, chatType string) ([]*models.WeChatContact, int64, error) {
	var contacts []*models.WeChatContact
	var total int64

	// 构建查询条件
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if chatType != "" {
		whereClause += fmt.Sprintf(" AND chat_type = $%d", argIndex)
		args = append(args, chatType)
		argIndex++
	}

	// 统计总数
	countSQL := "SELECT COUNT(*) FROM wechat_contacts " + whereClause
	err := s.db.QueryRow(countSQL, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("统计联系人总数失败: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	querySQL := fmt.Sprintf(`
		SELECT id, username, nickname, remark, avatar, type, chat_type, member_count, owner,
		       notice, message_count, last_active, is_blocked, privacy_level, created_at, updated_at
		FROM wechat_contacts %s
		ORDER BY last_active DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIndex, argIndex+1)

	args = append(args, pageSize, offset)
	rows, err := s.db.Query(querySQL, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("查询联系人失败: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		contact := &models.WeChatContact{}
		err := rows.Scan(&contact.ID, &contact.Username, &contact.Nickname, &contact.Remark,
			&contact.Avatar, &contact.Type, &contact.ChatType, &contact.MemberCount,
			&contact.Owner, &contact.Notice, &contact.MessageCount, &contact.LastActive,
			&contact.IsBlocked, &contact.PrivacyLevel, &contact.CreatedAt, &contact.UpdatedAt)
		if err != nil {
			log.Warnf("扫描联系人数据失败: %v", err)
			continue
		}

		contacts = append(contacts, contact)
	}

	return contacts, total, nil
}

// GetSyncRecords 获取同步记录
func (s *WeChatServiceV2) GetSyncRecords(page, pageSize int) ([]*models.WeChatSyncRecord, int64, error) {
	var records []*models.WeChatSyncRecord
	var total int64

	// 统计总数
	err := s.db.QueryRow("SELECT COUNT(*) FROM wechat_sync_records").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("统计同步记录总数失败: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	querySQL := `
		SELECT id, sync_type, source, total_messages, new_messages, updated_messages,
		       total_contacts, new_contacts, status, progress, error_message,
		       start_time, end_time, duration, created_at, updated_at
		FROM wechat_sync_records
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := s.db.Query(querySQL, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("查询同步记录失败: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		record := &models.WeChatSyncRecord{}
		var endTime sql.NullTime

		err := rows.Scan(&record.ID, &record.SyncType, &record.Source, &record.TotalMessages,
			&record.NewMessages, &record.UpdatedMessages, &record.TotalContacts,
			&record.NewContacts, &record.Status, &record.Progress, &record.ErrorMessage,
			&record.StartTime, &endTime, &record.Duration, &record.CreatedAt, &record.UpdatedAt)
		if err != nil {
			log.Warnf("扫描同步记录数据失败: %v", err)
			continue
		}

		if endTime.Valid {
			record.EndTime = &endTime.Time
		}

		records = append(records, record)
	}

	return records, total, nil
}

// insertSyncRecord 插入同步记录
func (s *WeChatServiceV2) insertSyncRecord(record *models.WeChatSyncRecord) (uint, error) {
	var id uint
	err := s.db.QueryRow(`
		INSERT INTO wechat_sync_records (sync_type, source, status, progress, start_time, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, record.SyncType, record.Source, record.Status, record.Progress,
		record.StartTime, record.CreatedAt, record.UpdatedAt).Scan(&id)

	return id, err
}

// updateSyncRecord 更新同步记录
func (s *WeChatServiceV2) updateSyncRecord(recordID uint, status string, progress float64, errorMsg string) {
	updateSQL := `
		UPDATE wechat_sync_records
		SET status = $1, progress = $2, error_message = $3, updated_at = $4
		WHERE id = $5
	`
	_, err := s.db.Exec(updateSQL, status, progress, errorMsg, time.Now(), recordID)
	if err != nil {
		log.Errorf("更新同步记录失败: %v", err)
	}
}

// importContactsV2 导入联系人
func (s *WeChatServiceV2) importContactsV2(contacts []*models.WeChatContact) (int, int, error) {
	var newCount, updatedCount int

	for _, contact := range contacts {
		var existingID sql.NullInt64
		err := s.db.QueryRow("SELECT id FROM wechat_contacts WHERE username = $1", contact.Username).Scan(&existingID)

		if err == sql.ErrNoRows {
			// 新增联系人
			insertSQL := `
				INSERT INTO wechat_contacts (username, nickname, remark, avatar, type, chat_type,
					member_count, owner, notice, message_count, last_active, is_blocked,
					privacy_level, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
			`
			_, err = s.db.Exec(insertSQL, contact.Username, contact.Nickname, contact.Remark,
				contact.Avatar, contact.Type, contact.ChatType, contact.MemberCount,
				contact.Owner, contact.Notice, contact.MessageCount, contact.LastActive,
				contact.IsBlocked, contact.PrivacyLevel, contact.CreatedAt, contact.UpdatedAt)
			if err != nil {
				log.Warnf("创建联系人失败 (%s): %v", contact.Username, err)
				continue
			}
			newCount++
		} else if err == nil {
			// 更新联系人
			updateSQL := `
				UPDATE wechat_contacts SET nickname = $1, remark = $2, avatar = $3, type = $4,
					chat_type = $5, member_count = $6, owner = $7, notice = $8, message_count = $9,
					last_active = $10, is_blocked = $11, privacy_level = $12, updated_at = $13
				WHERE id = $14
			`
			_, err = s.db.Exec(updateSQL, contact.Nickname, contact.Remark, contact.Avatar,
				contact.Type, contact.ChatType, contact.MemberCount, contact.Owner, contact.Notice,
				contact.MessageCount, contact.LastActive, contact.IsBlocked,
				contact.PrivacyLevel, contact.UpdatedAt, existingID.Int64)
			if err != nil {
				log.Warnf("更新联系人失败 (%s): %v", contact.Username, err)
				continue
			}
			updatedCount++
		} else {
			log.Warnf("查询联系人失败 (%s): %v", contact.Username, err)
			continue
		}
	}

	return newCount, updatedCount, nil
}

// importMessagesV2 导入消息
func (s *WeChatServiceV2) importMessagesV2(messages []*models.WeChatMessage) (int, int, error) {
	var newCount, updatedCount int

	for _, msg := range messages {
		var existingID sql.NullInt64
		err := s.db.QueryRow("SELECT id FROM wechat_messages WHERE svr_id = $1", msg.SvrID).Scan(&existingID)

		mentionsJSON, _ := json.Marshal(msg.Mentions)

		if err == sql.ErrNoRows {
			// 新增消息
			insertSQL := `
				INSERT INTO wechat_messages (svr_id, create_time, talker, type, sub_type, is_sender,
					seq, flag, status, content, display_content, chat_type, group_id, mentions,
					reply_to, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
			`
			_, err = s.db.Exec(insertSQL, msg.SvrID, msg.CreateTime, msg.Talker, msg.Type,
				msg.SubType, msg.IsSender, msg.Seq, msg.Flag, msg.Status, msg.Content,
				msg.DisplayContent, msg.ChatType, msg.GroupID, string(mentionsJSON),
				msg.ReplyTo, msg.CreatedAt, msg.UpdatedAt)
			if err != nil {
				log.Warnf("创建消息失败 (ID: %d): %v", msg.ID, err)
				continue
			}
			newCount++
		} else if err == nil {
			// 更新消息
			updateSQL := `
				UPDATE wechat_messages SET create_time = $1, talker = $2, type = $3, sub_type = $4,
					is_sender = $5, seq = $6, flag = $7, status = $8, content = $9,
					display_content = $10, chat_type = $11, group_id = $12, mentions = $13,
					reply_to = $14, updated_at = $15
				WHERE id = $16
			`
			_, err = s.db.Exec(updateSQL, msg.CreateTime, msg.Talker, msg.Type, msg.SubType,
				msg.IsSender, msg.Seq, msg.Flag, msg.Status, msg.Content, msg.DisplayContent,
				msg.ChatType, msg.GroupID, string(mentionsJSON), msg.ReplyTo,
				msg.UpdatedAt, existingID.Int64)
			if err != nil {
				log.Warnf("更新消息失败 (ID: %d): %v", msg.ID, err)
				continue
			}
			updatedCount++
		} else {
			log.Warnf("查询消息失败 (ID: %d): %v", msg.ID, err)
			continue
		}
	}

	return newCount, updatedCount, nil
}

// importChatsV2 导入聊天会话
func (s *WeChatServiceV2) importChatsV2(chats []*models.WeChatChat) error {
	for _, chat := range chats {
		var existingID sql.NullInt64
		err := s.db.QueryRow("SELECT id FROM wechat_chats WHERE chat_id = $1", chat.ChatID).Scan(&existingID)

		if err == sql.ErrNoRows {
			// 新增会话
			insertSQL := `
				INSERT INTO wechat_chats (chat_id, chat_type, name, avatar, message_count, unread_count,
					last_message, last_time, is_pinned, is_muted, is_archived, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
			`
			_, err = s.db.Exec(insertSQL, chat.ChatID, chat.ChatType, chat.Name, chat.Avatar,
				chat.MessageCount, chat.UnreadCount, chat.LastMessage, chat.LastTime,
				chat.IsPinned, chat.IsMuted, chat.IsArchived, chat.CreatedAt, chat.UpdatedAt)
			if err != nil {
				log.Warnf("创建会话失败 (%s): %v", chat.ChatID, err)
				continue
			}
		} else if err == nil {
			// 更新会话
			updateSQL := `
				UPDATE wechat_chats SET chat_type = $1, name = $2, avatar = $3, message_count = $4,
					unread_count = $5, last_message = $6, last_time = $7, is_pinned = $8,
					is_muted = $9, is_archived = $10, updated_at = $11
				WHERE id = $12
			`
			_, err = s.db.Exec(updateSQL, chat.ChatType, chat.Name, chat.Avatar, chat.MessageCount,
				chat.UnreadCount, chat.LastMessage, chat.LastTime, chat.IsPinned,
				chat.IsMuted, chat.IsArchived, chat.UpdatedAt, existingID.Int64)
			if err != nil {
				log.Warnf("更新会话失败 (%s): %v", chat.ChatID, err)
				continue
			}
		} else {
			log.Warnf("查询会话失败 (%s): %v", chat.ChatID, err)
			continue
		}
	}

	return nil
}

// GetChats 获取聊天会话列表 - 简化版本
func (s *WeChatServiceV2) GetChats(page, pageSize int, chatType string) ([]*models.WeChatChat, int64, error) {
	var chats []*models.WeChatChat
	var total int64

	// 构建查询条件
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if chatType != "" {
		whereClause += fmt.Sprintf(" AND chat_type = $%d", argIndex)
		args = append(args, chatType)
		argIndex++
	}

	// 统计总数
	countSQL := "SELECT COUNT(*) FROM wechat_chats " + whereClause
	err := s.db.QueryRow(countSQL, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("统计会话总数失败: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	querySQL := fmt.Sprintf(`
		SELECT id, chat_id, chat_type, name, avatar, message_count, unread_count,
		       last_message, last_time, is_pinned, is_muted, is_archived, created_at, updated_at
		FROM wechat_chats %s
		ORDER BY last_time DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIndex, argIndex+1)

	args = append(args, pageSize, offset)
	rows, err := s.db.Query(querySQL, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("查询会话失败: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		chat := &models.WeChatChat{}
		err := rows.Scan(&chat.ID, &chat.ChatID, &chat.ChatType, &chat.Name, &chat.Avatar,
			&chat.MessageCount, &chat.UnreadCount, &chat.LastMessage, &chat.LastTime,
			&chat.IsPinned, &chat.IsMuted, &chat.IsArchived, &chat.CreatedAt, &chat.UpdatedAt)
		if err != nil {
			log.Warnf("扫描会话数据失败: %v", err)
			continue
		}

		chats = append(chats, chat)
	}

	return chats, total, nil
}

// IncrementalSync 增量同步 - 简化版本
func (s *WeChatServiceV2) IncrementalSync(dbPath string, lastSyncTime int64) (*models.WeChatSyncRecord, error) {
	log.Infof("开始增量同步微信数据: %s", dbPath)

	// 创建同步记录
	syncRecord := &models.WeChatSyncRecord{
		SyncType:  "incremental",
		Source:    dbPath,
		Status:    "running",
		Progress:  0,
		StartTime: time.Now(),
	}

	recordID, err := s.insertSyncRecord(syncRecord)
	if err != nil {
		return nil, fmt.Errorf("创建同步记录失败: %w", err)
	}
	syncRecord.ID = int64(recordID)

	// 创建解析器
	parser, err := NewWeChatParser(dbPath)
	if err != nil {
		s.updateSyncRecord(recordID, "failed", 0, fmt.Sprintf("创建解析器失败: %v", err))
		return syncRecord, fmt.Errorf("创建微信解析器失败: %w", err)
	}
	defer parser.Close()

	// 解析数据库
	if err := parser.ParseDatabase(); err != nil {
		s.updateSyncRecord(recordID, "failed", 0, fmt.Sprintf("解析数据库失败: %v", err))
		return syncRecord, fmt.Errorf("解析微信数据库失败: %w", err)
	}

	// 获取指定时间之后的消息
	allMessages := parser.GetMessages()
	var newMessages []*models.WeChatMessage
	for _, msg := range allMessages {
		if msg.CreateTime > lastSyncTime {
			newMessages = append(newMessages, msg)
		}
	}

	// 导入新消息
	newCount, updatedCount, err := s.importMessagesV2(newMessages)
	if err != nil {
		s.updateSyncRecord(recordID, "failed", 0.8, fmt.Sprintf("导入新消息失败: %v", err))
		return syncRecord, fmt.Errorf("导入新消息失败: %w", err)
	}

	// 完成同步
	endTime := time.Now()
	duration := endTime.Sub(syncRecord.StartTime).Seconds()

	s.updateSyncRecord(recordID, "completed", 1.0, "")

	// 更新同步记录完成状态
	updateSQL := `
		UPDATE wechat_sync_records
		SET status = ?, progress = ?, total_messages = ?, new_messages = ?,
		    updated_messages = ?, end_time = ?, duration = ?, error_message = ?
		WHERE id = ?
	`
	_, err = s.db.Exec(updateSQL, "completed", 1.0, len(newMessages), newCount,
		updatedCount, endTime, int64(duration), "", recordID)

	if err != nil {
		log.Errorf("更新同步记录失败: %v", err)
	}

	log.Infof("✅ 增量同步完成: 新消息 %d 条 (%d 新增, %d 更新), 耗时 %.2f 秒",
		len(newMessages), newCount, updatedCount, duration)

	return syncRecord, nil
}

// 其他必要方法的简化实现
func (s *WeChatServiceV2) SearchMessages(keyword string, page, pageSize int) ([]*models.WeChatMessage, int64, error) {
	// 简化搜索实现
	return s.GetMessages(page, pageSize, "", "")
}

func (s *WeChatServiceV2) GetChatMessages(chatID string, page, pageSize int) ([]*models.WeChatMessage, int64, error) {
	// 简化聊天消息获取实现
	return s.GetMessages(page, pageSize, "", chatID)
}

func (s *WeChatServiceV2) GetMessageStatistics() (*MessageStatistics, error) {
	stats := &MessageStatistics{
		ByType: make(map[string]int64),
		ByDay:  make(map[string]int64),
	}

	// 统计总消息数
	s.db.QueryRow("SELECT COUNT(*) FROM wechat_messages").Scan(&stats.TotalMessages)

	// 简化的类型统计
	typeRows, _ := s.db.Query("SELECT type, COUNT(*) FROM wechat_messages GROUP BY type")
	if typeRows != nil {
		defer typeRows.Close()
		for typeRows.Next() {
			var msgType int
			var count int64
			typeRows.Scan(&msgType, &count)
			stats.ByType[strconv.Itoa(msgType)] = count
		}
	}

	return stats, nil
}

// MessageStatistics 消息统计信息
type MessageStatistics struct {
	TotalMessages int64            `json:"total_messages"`
	ByType        map[string]int64 `json:"by_type"`
	ByDay         map[string]int64 `json:"by_day"`
}


// Sync 同步微信数据
func (s *WeChatServiceV2) Sync(dbPath string, mode string) (*models.WeChatSyncRecord, error) {
	if mode == "incremental" {
		// 获取最后同步时间
		lastSyncTime := int64(0)
		s.db.QueryRow("SELECT start_time FROM wechat_sync_records ORDER BY created_at DESC LIMIT 1").Scan(&lastSyncTime)
		return s.IncrementalSync(dbPath, lastSyncTime)
	}
	return s.ImportFromDatabase(dbPath)
}

// GetStatistics 获取统计信息
func (s *WeChatServiceV2) GetStatistics() (*MessageStatistics, error) {
	return s.GetMessageStatistics()
}