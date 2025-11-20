-- 创建微信相关数据表
-- 执行脚本: sqlite3 data/twin-os.db < create_wechat_tables.sql

-- 创建微信消息表
CREATE TABLE IF NOT EXISTS wechat_messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    message_id TEXT UNIQUE NOT NULL,
    talker TEXT NOT NULL,
    content TEXT,
    type INTEGER NOT NULL DEFAULT 1,
    create_time DATETIME NOT NULL,
    send_time DATETIME,
    chat_type INTEGER DEFAULT 1,
    is_sender BOOLEAN DEFAULT FALSE,
    is_processed BOOLEAN DEFAULT FALSE,
    metadata TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 创建微信联系人表
CREATE TABLE IF NOT EXISTS wechat_contacts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_name TEXT UNIQUE NOT NULL,
    nickname TEXT,
    remark TEXT,
    avatar TEXT,
    type INTEGER NOT NULL DEFAULT 1,
    chat_type INTEGER DEFAULT 1,
    last_active DATETIME,
    contact_count INTEGER DEFAULT 0,
    is_blocked BOOLEAN DEFAULT FALSE,
    is_pinned BOOLEAN DEFAULT FALSE,
    metadata TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 创建微信聊天会话表
CREATE TABLE IF NOT EXISTS wechat_chats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    chat_id TEXT UNIQUE NOT NULL,
    name TEXT,
    avatar TEXT,
    chat_type INTEGER NOT NULL DEFAULT 1,
    last_message TEXT,
    last_time DATETIME,
    message_count INTEGER DEFAULT 0,
    unread_count INTEGER DEFAULT 0,
    is_pinned BOOLEAN DEFAULT FALSE,
    is_muted BOOLEAN DEFAULT FALSE,
    metadata TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 创建微信媒体文件表
CREATE TABLE IF NOT EXISTS wechat_media (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    message_id TEXT NOT NULL,
    file_name TEXT,
    file_path TEXT,
    file_size INTEGER DEFAULT 0,
    file_type TEXT,
    md5_hash TEXT,
    thumbnail_path TEXT,
    download_status INTEGER DEFAULT 0,
    metadata TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (message_id) REFERENCES wechat_messages(message_id)
);

-- 创建微信同步记录表
CREATE TABLE IF NOT EXISTS wechat_sync_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sync_id TEXT UNIQUE NOT NULL,
    sync_type TEXT NOT NULL,
    status TEXT NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME,
    message_count INTEGER DEFAULT 0,
    contact_count INTEGER DEFAULT 0,
    error_message TEXT,
    metadata TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
-- 微信消息表索引
CREATE INDEX IF NOT EXISTS idx_wechat_messages_create_time ON wechat_messages(create_time);
CREATE INDEX IF NOT EXISTS idx_wechat_messages_talker_create_time ON wechat_messages(talker, create_time);
CREATE INDEX IF NOT EXISTS idx_wechat_messages_type_create_time ON wechat_messages(type, create_time);
CREATE INDEX IF NOT EXISTS idx_wechat_messages_is_sender_create_time ON wechat_messages(is_sender, create_time);
CREATE INDEX IF NOT EXISTS idx_wechat_messages_chat_type_create_time ON wechat_messages(chat_type, create_time);
CREATE INDEX IF NOT EXISTS idx_wechat_messages_processed ON wechat_messages(is_processed);

-- 微信联系人表索引
CREATE INDEX IF NOT EXISTS idx_wechat_contacts_type ON wechat_contacts(type);
CREATE INDEX IF NOT EXISTS idx_wechat_contacts_chat_type ON wechat_contacts(chat_type);
CREATE INDEX IF NOT EXISTS idx_wechat_contacts_last_active ON wechat_contacts(last_active);

-- 微信聊天会话表索引
CREATE INDEX IF NOT EXISTS idx_wechat_chats_chat_type ON wechat_chats(chat_type);
CREATE INDEX IF NOT EXISTS idx_wechat_chats_last_time ON wechat_chats(last_time);
CREATE INDEX IF NOT EXISTS idx_wechat_chats_is_pinned ON wechat_chats(is_pinned);

-- 微信媒体文件表索引
CREATE INDEX IF NOT EXISTS idx_wechat_media_message_id ON wechat_media(message_id);
CREATE INDEX IF NOT EXISTS idx_wechat_media_file_type ON wechat_media(file_type);

-- 微信同步记录表索引
CREATE INDEX IF NOT EXISTS idx_wechat_sync_records_sync_type ON wechat_sync_records(sync_type);
CREATE INDEX IF NOT EXISTS idx_wechat_sync_records_status ON wechat_sync_records(status);
CREATE INDEX IF NOT EXISTS idx_wechat_sync_records_start_time ON wechat_sync_records(start_time);