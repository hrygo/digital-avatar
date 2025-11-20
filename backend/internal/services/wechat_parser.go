package services

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/hrygo/log"
	"twin-os/backend/internal/models"
)

// WeChatParser 微信数据库解析器
type WeChatParser struct {
	dbPath       string
	db           *sql.DB
	contacts     map[string]*models.WeChatContact
	messages     []*models.WeChatMessage
	chats        map[string]*models.WeChatChat
}

// NewWeChatParser 创建微信解析器
func NewWeChatParser(dbPath string) (*WeChatParser, error) {
	// 检查文件是否存在
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("微信数据库文件不存在: %s", dbPath)
	}

	// 打开数据库
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("打开微信数据库失败: %w", err)
	}

	return &WeChatParser{
		dbPath:   dbPath,
		db:       db,
		contacts: make(map[string]*models.WeChatContact),
		messages: make([]*models.WeChatMessage, 0),
		chats:    make(map[string]*models.WeChatChat),
	}, nil
}

// ParseDatabase 解析微信数据库
func (p *WeChatParser) ParseDatabase() error {
	log.Info("开始解析微信数据库")

	// 1. 解析联系人信息
	if err := p.parseContacts(); err != nil {
		return fmt.Errorf("解析联系人失败: %w", err)
	}
	log.Infof("✅ 解析联系人完成，共 %d 个联系人", len(p.contacts))

	// 2. 解析消息记录
	if err := p.parseMessages(); err != nil {
		return fmt.Errorf("解析消息失败: %w", err)
	}
	log.Infof("✅ 解析消息完成，共 %d 条消息", len(p.messages))

	// 3. 构建聊天会话
	if err := p.buildChats(); err != nil {
		return fmt.Errorf("构建聊天会话失败: %w", err)
	}
	log.Infof("✅ 构建聊天会话完成，共 %d 个会话", len(p.chats))

	log.Info("✅ 微信数据库解析完成")
	return nil
}

// parseContacts 解析联系人信息
func (p *WeChatParser) parseContacts() error {
	query := `
		SELECT username, nickname, remark, type, headImgUrl
		FROM rcontact
		WHERE type & 1 = 1
		AND username NOT LIKE 'gh_%'
		AND username NOT LIKE 'filehelper'
		AND username NOT LIKE 'weixin'
		AND username NOT LIKE 'fmessage'
		AND username NOT LIKE 'medianote'
		AND username NOT LIKE 'floatbottle'
	`

	rows, err := p.db.Query(query)
	if err != nil {
		return fmt.Errorf("查询联系人失败: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var contact models.WeChatContact
		var avatar sql.NullString

		if err := rows.Scan(&contact.Username, &contact.Nickname,
			&contact.Remark, &contact.Type, &avatar); err != nil {
			log.Warnf("扫描联系人数据失败: %v", err)
			continue
		}

		// 设置头像路径
		if avatar.Valid {
			contact.Avatar = avatar.String
		}

		// 确定聊天类型
		if strings.Contains(contact.Username, "@chatroom") {
			contact.ChatType = models.ChatTypeGroup
		} else {
			contact.ChatType = models.ChatTypePrivate
		}

		// 获取群成员数量（如果是群聊）
		if contact.ChatType == models.ChatTypeGroup {
			contact.MemberCount = p.getGroupMemberCount(contact.Username)
		}

		contact.CreatedAt = time.Now()
		contact.UpdatedAt = time.Now()

		p.contacts[contact.Username] = &contact
	}

	return nil
}

// parseMessages 解析消息记录
func (p *WeChatParser) parseMessages() error {
	query := `
		SELECT msgId, msgSvrId, type, subType, isSender, createTime,
			   talker, content, displayContent, status, lvbuffer
		FROM message
		ORDER BY createTime ASC
	`

	rows, err := p.db.Query(query)
	if err != nil {
		return fmt.Errorf("查询消息失败: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var msg models.WeChatMessage
		var displayContent sql.NullString

		if err := rows.Scan(&msg.ID, &msg.SvrID, &msg.Type, &msg.SubType,
			&msg.IsSender, &msg.CreateTime, &msg.Talker, &msg.Content,
			&displayContent, &msg.Status, &msg.LVBuffer); err != nil {
			log.Warnf("扫描消息数据失败: %v", err)
			continue
		}

		// 设置显示内容
		if displayContent.Valid {
			msg.DisplayContent = displayContent.String
		} else {
			msg.DisplayContent = msg.Content
		}

		// 确定聊天类型
		if strings.Contains(msg.Talker, "@chatroom") {
			msg.ChatType = models.ChatTypeGroup
		} else {
			msg.ChatType = models.ChatTypePrivate
		}

		// 解析消息内容
		if err := p.parseMessageContent(&msg); err != nil {
			log.Warnf("解析消息内容失败 (ID: %d): %v", msg.ID, err)
		}

		// 提取@提及信息
		if msg.ChatType == models.ChatTypeGroup {
			msg.Mentions = p.extractMentions(msg.Content)
		}

		msg.CreatedAt = time.Now()
		msg.UpdatedAt = time.Now()

		p.messages = append(p.messages, &msg)
	}

	return nil
}

// parseMessageContent 解析消息内容
func (p *WeChatParser) parseMessageContent(msg *models.WeChatMessage) error {
	switch models.WeChatMessageType(msg.Type) {
	case models.MessageTypeText:
		// 文本消息
		if msg.DisplayContent == "" {
			msg.DisplayContent = msg.Content
		}

	case models.MessageTypeImage:
		// 图片消息
		msg.DisplayContent = "[图片]"
		if len(msg.LVBuffer) > 0 {
			msg.MediaPath = p.extractImagePath(msg.LVBuffer)
		}

	case models.MessageTypeVoice:
		// 语音消息
		msg.DisplayContent = "[语音]"
		if len(msg.LVBuffer) > 0 {
			msg.MediaPath = p.extractVoicePath(msg.LVBuffer)
		}

	case models.MessageTypeVideo:
		// 视频消息
		msg.DisplayContent = "[视频]"
		if len(msg.LVBuffer) > 0 {
			msg.MediaPath = p.extractVideoPath(msg.LVBuffer)
		}

	case models.MessageTypeEmoji:
		// 表情消息
		msg.DisplayContent = "[表情]"
		if len(msg.LVBuffer) > 0 {
			msg.MediaPath = p.extractEmojiPath(msg.LVBuffer)
		}

	case models.MessageTypeLocation:
		// 位置消息
		msg.DisplayContent = "[位置]"
		if len(msg.Content) > 0 {
			// 解析位置信息
			if location := p.parseLocationContent(msg.Content); location != "" {
				msg.DisplayContent = fmt.Sprintf("[位置] %s", location)
			}
		}

	case models.MessageTypeLink:
		// 链接或文件消息
		if len(msg.Content) > 0 {
			if linkInfo := p.parseLinkContent(msg.Content); linkInfo != nil {
				msg.DisplayContent = fmt.Sprintf("[%s]", linkInfo.Title)
				msg.LinkTitle = linkInfo.Title
				msg.LinkDescription = linkInfo.Description
				msg.LinkURL = linkInfo.URL
			} else {
				msg.DisplayContent = "[文件]"
			}
		}

	case models.MessageTypeSystem:
		// 系统消息
		msg.DisplayContent = p.parseSystemMessage(msg.Content)

	default:
		msg.DisplayContent = "[未知消息类型]"
	}

	return nil
}

// buildChats 构建聊天会话
func (p *WeChatParser) buildChats() error {
	// 按聊天者分组消息
	messageGroups := make(map[string][]*models.WeChatMessage)
	for _, msg := range p.messages {
		messageGroups[msg.Talker] = append(messageGroups[msg.Talker], msg)
	}

	// 为每个聊天者创建会话
	for talker, msgs := range messageGroups {
		chat := &models.WeChatChat{
			ChatID:   talker,
			ChatType: msgs[0].ChatType,
		}

		// 设置会话名称
		if contact, exists := p.contacts[talker]; exists {
			if contact.Remark != "" {
				chat.Name = contact.Remark
			} else {
				chat.Name = contact.Nickname
			}
			chat.Avatar = contact.Avatar
		} else {
			chat.Name = talker
		}

		// 统计消息和最后一条消息
		chat.MessageCount = len(msgs)
		if len(msgs) > 0 {
			lastMsg := msgs[len(msgs)-1]
			chat.LastMessage = lastMsg.DisplayContent
			chat.LastTime = lastMsg.CreateTime
		}

		chat.CreatedAt = time.Now()
		chat.UpdatedAt = time.Now()

		p.chats[talker] = chat
	}

	return nil
}

// getGroupMemberCount 获取群成员数量
func (p *WeChatParser) getGroupMemberCount(groupID string) int {
	query := `
		SELECT COUNT(*)
		FROM chatroom
		WHERE chatroomname = ?
	`

	var count int
	err := p.db.QueryRow(query, groupID).Scan(&count)
	if err != nil {
		log.Warnf("获取群成员数量失败: %v", err)
		return 0
	}

	return count
}

// extractMentions 提取@提及的用户
func (p *WeChatParser) extractMentions(content string) models.JSON {
	mentions := make(models.JSON)

	// 匹配@用户名格式
	pattern := regexp.MustCompile(`@([^\s@]+)`)
	matches := pattern.FindAllStringSubmatch(content, -1)

	mentionedUsers := make([]string, 0)
	for _, match := range matches {
		if len(match) > 1 {
			mentionedUsers = append(mentionedUsers, match[1])
		}
	}

	if len(mentionedUsers) > 0 {
		mentions["users"] = mentionedUsers
	}

	return mentions
}

// extractImagePath 提取图片路径
func (p *WeChatParser) extractImagePath(lvBuffer []byte) string {
	if len(lvBuffer) < 100 {
		return ""
	}

	// 从LVBuffer中提取图片文件路径
	// 这里需要根据微信的实际数据格式进行解析
	// 简化处理，假设路径在固定位置
	start := 20
	end := start + 100
	if end > len(lvBuffer) {
		end = len(lvBuffer)
	}

	pathBytes := lvBuffer[start:end]
	return strings.TrimRight(string(pathBytes), "\x00")
}

// extractVoicePath 提取语音路径
func (p *WeChatParser) extractVoicePath(lvBuffer []byte) string {
	// 类似图片路径提取的逻辑
	return p.extractImagePath(lvBuffer)
}

// extractVideoPath 提取视频路径
func (p *WeChatParser) extractVideoPath(lvBuffer []byte) string {
	// 类似图片路径提取的逻辑
	return p.extractImagePath(lvBuffer)
}

// extractEmojiPath 提取表情路径
func (p *WeChatParser) extractEmojiPath(lvBuffer []byte) string {
	// 类似图片路径提取的逻辑
	return p.extractImagePath(lvBuffer)
}

// parseLocationContent 解析位置内容
func (p *WeChatParser) parseLocationContent(content string) string {
	// 简单的位置信息解析
	if strings.Contains(content, "lat:") && strings.Contains(content, "lng:") {
		// 提取经纬度信息
		parts := strings.Split(content, ";")
		for _, part := range parts {
			if strings.Contains(part, "title:") {
				return strings.TrimPrefix(part, "title:")
			}
		}
	}
	return ""
}

// parseLinkContent 解析链接内容
func (p *WeChatParser) parseLinkContent(content string) *LinkInfo {
	linkInfo := &LinkInfo{}

	// 简单的链接解析逻辑
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "title:") {
			linkInfo.Title = strings.TrimPrefix(line, "title:")
		} else if strings.HasPrefix(line, "des:") {
			linkInfo.Description = strings.TrimPrefix(line, "des:")
		} else if strings.HasPrefix(line, "url:") {
			linkInfo.URL = strings.TrimPrefix(line, "url:")
		}
	}

	if linkInfo.Title == "" {
		return nil
	}

	return linkInfo
}

// parseSystemMessage 解析系统消息
func (p *WeChatParser) parseSystemMessage(content string) string {
	// 常见系统消息类型
	if strings.Contains(content, "加入了群聊") {
		return "[加入群聊]"
	} else if strings.Contains(content, "退出了群聊") {
		return "[退出群聊]"
	} else if strings.Contains(content, "修改群名为") {
		return "[修改群名]"
	} else if strings.Contains(content, "邀请") && strings.Contains(content, "加入了群聊") {
		return "[邀请加入群聊]"
	} else if strings.Contains(content, "已撤回") {
		return "[撤回消息]"
	} else if strings.Contains(content, "拍了拍") {
		return "[拍了拍]"
	}

	return "[系统消息]"
}

// GetContacts 获取解析出的联系人
func (p *WeChatParser) GetContacts() []*models.WeChatContact {
	contacts := make([]*models.WeChatContact, 0, len(p.contacts))
	for _, contact := range p.contacts {
		contacts = append(contacts, contact)
	}
	return contacts
}

// GetMessages 获取解析出的消息
func (p *WeChatParser) GetMessages() []*models.WeChatMessage {
	return p.messages
}

// GetChats 获取解析出的聊天会话
func (p *WeChatParser) GetChats() []*models.WeChatChat {
	chats := make([]*models.WeChatChat, 0, len(p.chats))
	for _, chat := range p.chats {
		chats = append(chats, chat)
	}
	return chats
}

// Close 关闭解析器
func (p *WeChatParser) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

// LinkInfo 链接信息
type LinkInfo struct {
	Title       string
	Description string
	URL         string
}

// FindWeChatDB 查找微信数据库文件
func FindWeChatDB(basePath string) ([]string, error) {
	var dbFiles []string

	// 微信数据库可能的路径
	possiblePaths := []string{
		filepath.Join(basePath, "Documents", "MicroMsg"),
		filepath.Join(basePath, "Library", "Containers", "com.tencent.xin", "Data", "Documents", "MicroMsg"),
		filepath.Join(basePath, "AppData", "Roaming", "Tencent", "MicroMsg"),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		// 查找所有以MD5命名的文件夹
		entries, err := os.ReadDir(path)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			// 检查是否为32位MD5
			if len(entry.Name()) == 32 {
				dbPath := filepath.Join(path, entry.Name(), "DB", "MicroMsg.db")
				if _, err := os.Stat(dbPath); err == nil {
					dbFiles = append(dbFiles, dbPath)
				}
			}
		}
	}

	if len(dbFiles) == 0 {
		return nil, fmt.Errorf("未找到微信数据库文件")
	}

	return dbFiles, nil
}