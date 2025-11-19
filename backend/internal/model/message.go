package model

import (
	"time"
)

// Message 微信消息模型
type Message struct {
	ID        int64     `json:"id" db:"id"`
	MessageID string    `json:"message_id" db:"message_id"`
	TalkerID  string    `json:"talker_id" db:"talker_id"`
	Type      int       `json:"type" db:"type"`
	Content   string    `json:"content" db:"content"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Contact 微信联系人模型
type Contact struct {
	ID         int64     `json:"id" db:"id"`
	UserName   string    `json:"user_name" db:"user_name"`
	NickName   string    `json:"nick_name" db:"nick_name"`
	Remark     string    `json:"remark" db:"remark"`
	Type       int       `json:"type" db:"type"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// AnalysisResult 分析结果模型
type AnalysisResult struct {
	ID          int64     `json:"id" db:"id"`
	Type        string    `json:"type" db:"type"` // briefing, todo, connection
	MessageIDs  string    `json:"message_ids" db:"message_ids"` // JSON array
	Content     string    `json:"content" db:"content"`
	Metadata    string    `json:"metadata" db:"metadata"` // JSON object
	Confidence  float64   `json:"confidence" db:"confidence"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	ProcessedAt time.Time `json:"processed_at" db:"processed_at"`
}

// SystemSettings 系统设置模型
type SystemSettings struct {
	ID             int64     `json:"id" db:"id"`
	Key            string    `json:"key" db:"key"`
	Value          string    `json:"value" db:"value"`
	Description    string    `json:"description" db:"description"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}