package migrations

import (
	"github.com/hrygo/log"
	"gorm.io/gorm"
	"twin-os/backend/internal/models"
)

// Migration002 创建微信相关数据表
type Migration002 struct{}

func (m *Migration002) Migrate(db *gorm.DB) error {
	log.Info("开始执行 Migration002: 创建微信相关数据表")

	// 创建微信消息表
	if err := db.AutoMigrate(&models.WeChatMessage{}); err != nil {
		log.Errorf("创建微信消息表失败: %v", err)
		return err
	}
	log.Info("✅ 微信消息表创建成功")

	// 创建微信联系人表
	if err := db.AutoMigrate(&models.WeChatContact{}); err != nil {
		log.Errorf("创建微信联系人表失败: %v", err)
		return err
	}
	log.Info("✅ 微信联系人表创建成功")

	// 创建微信聊天会话表
	if err := db.AutoMigrate(&models.WeChatChat{}); err != nil {
		log.Errorf("创建微信聊天会话表失败: %v", err)
		return err
	}
	log.Info("✅ 微信聊天会话表创建成功")

	// 创建微信媒体文件表
	if err := db.AutoMigrate(&models.WeChatMedia{}); err != nil {
		log.Errorf("创建微信媒体文件表失败: %v", err)
		return err
	}
	log.Info("✅ 微信媒体文件表创建成功")

	// 创建微信同步记录表
	if err := db.AutoMigrate(&models.WeChatSyncRecord{}); err != nil {
		log.Errorf("创建微信同步记录表失败: %v", err)
		return err
	}
	log.Info("✅ 微信同步记录表创建成功")

	// 创建索引
	if err := m.createIndexes(db); err != nil {
		log.Errorf("创建索引失败: %v", err)
		return err
	}
	log.Info("✅ 微信数据表索引创建成功")

	log.Info("✅ Migration002 执行完成")
	return nil
}

func (m *Migration002) createIndexes(db *gorm.DB) error {
	// 微信消息表索引
	messageIndexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_wechat_messages_create_time ON wechat_messages(create_time)",
		"CREATE INDEX IF NOT EXISTS idx_wechat_messages_talker_create_time ON wechat_messages(talker, create_time)",
		"CREATE INDEX IF NOT EXISTS idx_wechat_messages_type_create_time ON wechat_messages(type, create_time)",
		"CREATE INDEX IF NOT EXISTS idx_wechat_messages_is_sender_create_time ON wechat_messages(is_sender, create_time)",
		"CREATE INDEX IF NOT EXISTS idx_wechat_messages_chat_type_create_time ON wechat_messages(chat_type, create_time)",
		"CREATE INDEX IF NOT EXISTS idx_wechat_messages_processed ON wechat_messages(is_processed)",
		"CREATE INDEX IF NOT EXISTS idx_wechat_messages_content_gin ON wechat_messages_content USING gin(to_tsvector('simple', content))",
	}

	// 微信联系人表索引
	contactIndexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_wechat_contacts_type ON wechat_contacts(type)",
		"CREATE INDEX IF NOT EXISTS idx_wechat_contacts_chat_type ON wechat_contacts(chat_type)",
		"CREATE INDEX IF NOT EXISTS idx_wechat_contacts_last_active ON wechat_contacts(last_active)",
		"CREATE INDEX IF NOT EXISTS idx_wechat_contacts_nickname_gin ON wechat_contacts(nickname USING gin(to_tsvector('simple', nickname)))",
		"CREATE INDEX IF NOT EXISTS idx_wechat_contacts_remark_gin ON wechat_contacts(remark USING gin(to_tsvector('simple', remark)))",
	}

	// 微信聊天会话表索引
	chatIndexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_wechat_chats_chat_type ON wechat_chats(chat_type)",
		"CREATE INDEX IF NOT EXISTS idx_wechat_chats_last_time ON wechat_chats(last_time)",
		"CREATE INDEX IF NOT EXISTS idx_wechat_chats_is_pinned ON wechat_chats(is_pinned)",
		"CREATE INDEX IF NOT EXISTS idx_wechat_chats_name_gin ON wechat_chats(name USING gin(to_tsvector('simple', name)))",
	}

	// 微信媒体文件表索引
	mediaIndexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_wechat_media_message_id ON wechat_media(message_id)",
		"CREATE INDEX IF NOT EXISTS idx_wechat_media_type ON wechat_media(type)",
		"CREATE INDEX IF NOT EXISTS idx_wechat_media_file_size ON wechat_media(file_size)",
	}

	// 微信同步记录表索引
	syncIndexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_wechat_sync_records_sync_type ON wechat_sync_records(sync_type)",
		"CREATE INDEX IF NOT EXISTS idx_wechat_sync_records_status ON wechat_sync_records(status)",
		"CREATE INDEX IF NOT EXISTS idx_wechat_sync_records_start_time ON wechat_sync_records(start_time)",
	}

	allIndexes := append(append(append(messageIndexes, contactIndexes...), chatIndexes...), append(mediaIndexes, syncIndexes...)...)

	for _, indexSQL := range allIndexes {
		if err := db.Exec(indexSQL).Error; err != nil {
			return err
		}
	}

	return nil
}

func (m *Migration002) Rollback(db *gorm.DB) error {
	log.Info("开始回滚 Migration002: 删除微信相关数据表")

	tables := []string{
		"wechat_sync_records",
		"wechat_media",
		"wechat_chats",
		"wechat_contacts",
		"wechat_messages",
	}

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			log.Errorf("删除表 %s 失败: %v", table, err)
			return err
		}
		log.Infof("✅ 表 %s 删除成功", table)
	}

	log.Info("✅ Migration002 回滚完成")
	return nil
}