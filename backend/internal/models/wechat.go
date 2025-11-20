package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// WeChatMessageType 微信消息类型
type WeChatMessageType int

const (
	MessageTypeText     WeChatMessageType = 1 // 文本消息
	MessageTypeImage    WeChatMessageType = 3 // 图片消息
	MessageTypeVoice    WeChatMessageType = 34 // 语音消息
	MessageTypeVideo    WeChatMessageType = 43 // 视频消息
	MessageTypeEmoji    WeChatMessageType = 47 // 表情消息
	MessageTypeLocation WeChatMessageType = 48 // 位置消息
	MessageTypeLink     WeChatMessageType = 49 // 链接消息 (包含文件消息)
	MessageTypeSystem   WeChatMessageType = 10000 // 系统消息
)

// WeChatChatType 聊天类型
type WeChatChatType string

const (
	ChatTypePrivate WeChatChatType = "private" // 私聊
	ChatTypeGroup   WeChatChatType = "group"   // 群聊
)

// WeChatMessage 微信消息模型
type WeChatMessage struct {
	ID             int64             `json:"id" gorm:"primaryKey;autoIncrement"`
	SvrID          int64             `json:"svr_id" gorm:"index"` // 微信服务器消息ID
	CreateTime     int64             `json:"create_time" gorm:"index"` // 创建时间戳
	Talker         string            `json:"talker" gorm:"index;size:255"` // 发送者ID
	Type           WeChatMessageType `json:"type" gorm:"index"` // 消息类型
	SubType        int               `json:"sub_type"` // 子类型
	IsSender       int               `json:"is_sender" gorm:"index"` // 是否为发送者 1=是 0=否
	Seq            int64             `json:"seq"` // 消息序列号
	Flag           int               `json:"flag"` // 消息标记
	Status         int               `json:"status"` // 消息状态

	// 消息内容
	Content          string `json:"content" gorm:"type:text"` // 原始消息内容
	DisplayContent   string `json:"display_content" gorm:"type:text"` // 显示内容
	CompressContent  []byte `json:"compress_content" gorm:"type:blob"` // 压缩内容
	LVBuffer         []byte `json:"lv_buffer" gorm:"type:blob"` // LV缓冲区
	BytesExtra       []byte `json:"bytes_extra" gorm:"type:blob"` // 额外字节

	// 解析后的结构化数据
	MediaPath        string `json:"media_path" gorm:"size:500"` // 媒体文件路径
	ThumbPath        string `json:"thumb_path" gorm:"size:500"` // 缩略图路径
	LinkTitle        string `json:"link_title" gorm:"size:255"` // 链接标题
	LinkDescription  string `json:"link_description" gorm:"type:text"` // 链接描述
	LinkURL          string `json:"link_url" gorm:"size:500"` // 链接URL

	// 元数据
	TalkerID         int64  `json:"talker_id" gorm:"index"` // 发送者内部ID
	ChatType         WeChatChatType `json:"chat_type" gorm:"index;size:20"` // 聊天类型
	GroupID          string `json:"group_id" gorm:"size:255;index"` // 群聊ID（群聊时有效）
	Mentions         JSON   `json:"mentions" gorm:"type:json"` // @提及的用户
	ReplyTo          string `json:"reply_to" gorm:"size:255"` // 回复的消息ID

	// 隐私和安全相关
	IsProcessed      bool   `json:"is_processed" gorm:"default:false;index"` // 是否已处理
	AnonymizedContent string `json:"anonymized_content" gorm:"type:text"` // 脱敏后的内容
	PrivacyLevel     int    `json:"privacy_level" gorm:"default:0"` // 隐私级别

	// 系统字段
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// WeChatContact 微信联系人模型
type WeChatContact struct {
	ID           int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Username     string `json:"username" gorm:"uniqueIndex;size:255"` // 微信号/用户ID
	Nickname     string `json:"nickname" gorm:"size:255"` // 昵称
	Remark       string `json:"remark" gorm:"size:255"` // 备注名
	Avatar       string `json:"avatar" gorm:"size:500"` // 头像路径

	// 联系人类型
	Type         int    `json:"type" gorm:"index"` // 联系人类型
	ChatType     WeChatChatType `json:"chat_type" gorm:"index;size:20"` // 聊天类型

	// 群聊特有字段
	MemberCount  int    `json:"member_count" gorm:"default:0"` // 群成员数量
	Owner        string `json:"owner" gorm:"size:255"` // 群主ID
	Notice       string `json:"notice" gorm:"type:text"` // 群公告

	// 统计信息
	MessageCount int    `json:"message_count" gorm:"default:0"` // 消息数量
	LastActive   int64  `json:"last_active" gorm:"index"` // 最后活跃时间

	// 隐私和安全
	IsBlocked    bool   `json:"is_blocked" gorm:"default:false"` // 是否被拉黑
	PrivacyLevel int    `json:"privacy_level" gorm:"default:0"` // 隐私级别

	// 系统字段
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// WeChatChat 微信聊天会话模型
type WeChatChat struct {
	ID           int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	ChatID       string `json:"chat_id" gorm:"uniqueIndex;size:255"` // 聊天ID
	ChatType     WeChatChatType `json:"chat_type" gorm:"index;size:20"` // 聊天类型
	Name         string `json:"name" gorm:"size:255"` // 聊天名称（群聊为群名，私聊为对方昵称）
	Avatar       string `json:"avatar" gorm:"size:500"` // 聊天头像

	// 聊天统计
	MessageCount int    `json:"message_count" gorm:"default:0"` // 总消息数
	UnreadCount  int    `json:"unread_count" gorm:"default:0"` // 未读消息数
	LastMessage  string `json:"last_message" gorm:"type:text"` // 最后一条消息
	LastTime     int64  `json:"last_time" gorm:"index"` // 最后消息时间

	// 会话设置
	IsPinned     bool   `json:"is_pinned" gorm:"default:false"` // 是否置顶
	IsMuted      bool   `json:"is_muted" gorm:"default:false"` // 是否静音
	IsArchived   bool   `json:"is_archived" gorm:"default:false"` // 是否归档

	// 系统字段
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// WeChatMedia 微信媒体文件模型
type WeChatMedia struct {
	ID           int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	MediaID      string `json:"media_id" gorm:"uniqueIndex;size:255"` // 媒体文件ID
	MessageID    int64  `json:"message_id" gorm:"index"` // 关联的消息ID
	Type         string `json:"type" gorm:"size:50"` // 媒体类型：image/voice/video/file

	// 文件信息
	FileName     string `json:"file_name" gorm:"size:500"` // 文件名
	FilePath     string `json:"file_path" gorm:"size:500"` // 文件路径
	FileSize     int64  `json:"file_size"` // 文件大小
	ThumbPath    string `json:"thumb_path" gorm:"size:500"` // 缩略图路径

	// 元数据
	Width        int    `json:"width"` // 宽度（图片/视频）
	Height       int    `json:"height"` // 高度（图片/视频）
	Duration     int    `json:"duration"` // 时长（语音/视频）

	// 系统字段
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// WeChatSyncRecord 微信数据同步记录
type WeChatSyncRecord struct {
	ID           int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	SyncType     string `json:"sync_type" gorm:"size:50;index"` // 同步类型：full/incremental
	Source       string `json:"source" gorm:"size:255"` // 数据源路径

	// 同步统计
	TotalMessages int   `json:"total_messages"` // 总消息数
	NewMessages   int   `json:"new_messages"` // 新消息数
	UpdatedMessages int `json:"updated_messages"` // 更新消息数
	TotalContacts int   `json:"total_contacts"` // 总联系人数
	NewContacts   int   `json:"new_contacts"` // 新联系人数

	// 同步状态
	Status       string `json:"status" gorm:"size:20;index"` // 状态：running/completed/failed
	Progress     float64 `json:"progress"` // 进度百分比
	ErrorMessage string `json:"error_message" gorm:"type:text"` // 错误信息

	// 时间信息
	StartTime    time.Time `json:"start_time"` // 开始时间
	EndTime      *time.Time `json:"end_time"` // 结束时间
	Duration     int64 `json:"duration"` // 耗时（秒）

	// 系统字段
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// JSON 自定义JSON类型，用于处理JSON字段
type JSON map[string]interface{}

// Value 实现driver.Valuer接口
func (j JSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan 实现sql.Scanner接口
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into JSON", value)
	}

	return json.Unmarshal(bytes, &j)
}

// TableName 指定表名
func (WeChatMessage) TableName() string {
	return "wechat_messages"
}

func (WeChatContact) TableName() string {
	return "wechat_contacts"
}

func (WeChatChat) TableName() string {
	return "wechat_chats"
}

func (WeChatMedia) TableName() string {
	return "wechat_media"
}

func (WeChatSyncRecord) TableName() string {
	return "wechat_sync_records"
}

// MessageTypeString 获取消息类型字符串
func (m WeChatMessageType) String() string {
	switch m {
	case MessageTypeText:
		return "text"
	case MessageTypeImage:
		return "image"
	case MessageTypeVoice:
		return "voice"
	case MessageTypeVideo:
		return "video"
	case MessageTypeEmoji:
		return "emoji"
	case MessageTypeLocation:
		return "location"
	case MessageTypeLink:
		return "file" // 链接和文件消息统一处理
	case MessageTypeSystem:
		return "system"
	default:
		return "unknown"
	}
}

// IsTextMessage 判断是否为文本消息
func (m WeChatMessageType) IsTextMessage() bool {
	return m == MessageTypeText
}

// IsMediaMessage 判断是否为媒体消息
func (m WeChatMessageType) IsMediaMessage() bool {
	return m == MessageTypeImage || m == MessageTypeVoice || m == MessageTypeVideo
}

// IsFileMessage 判断是否为文件消息
func (m WeChatMessageType) IsFileMessage() bool {
	return m == MessageTypeLink
}

// IsGroupChat 判断是否为群聊
func (c WeChatChatType) IsGroupChat() bool {
	return c == ChatTypeGroup
}

// IsPrivateChat 判断是否为私聊
func (c WeChatChatType) IsPrivateChat() bool {
	return c == ChatTypePrivate
}