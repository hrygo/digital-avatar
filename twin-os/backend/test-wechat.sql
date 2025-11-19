-- 模拟微信数据库结构 (SQLite)
-- 用于测试TwinOS的数据处理功能

-- 创建消息表
CREATE TABLE IF NOT EXISTS message (
    msgId INTEGER PRIMARY KEY,
    talkerId TEXT NOT NULL,
    type INTEGER NOT NULL,
    content TEXT,
    createTime INTEGER NOT NULL
);

-- 创建联系人表
CREATE TABLE IF NOT EXISTS rcontact (
    username TEXT PRIMARY KEY,
    nickname TEXT,
    remark TEXT,
    type INTEGER NOT NULL
);

-- 插入模拟数据 (使用2025年11月19日的时间戳)
INSERT OR IGNORE INTO message (msgId, talkerId, type, content, createTime) VALUES
(1, 'filehelper', 1, '欢迎使用TwinOS系统！', 1763532000),
(2, 'testuser1', 1, '明天下午2点开会讨论项目进度', 1763532100),
(3, 'testuser1', 1, '请记得准备项目文档和进度报告', 1763532200),
(4, 'testuser2', 1, '产品已经完成，可以安排发布了', 1763532300),
(5, 'testuser3', 1, '下周三有个重要的客户会议，需要提前准备材料', 1763532400),
(6, 'testuser2', 1, '张总对项目很满意，说我们的团队效率很高', 1763532500),
(7, 'testuser1', 1, '收到，我会在周五前完成所有准备工作', 1763532600),
(8, 'testuser3', 1, '李经理提到你最近的决策很准确', 1763532700),
(9, 'testuser4', 1, '新版本的系统架构设计很棒', 1763532800),
(10, 'testuser1', 1, '感谢认可，我们继续努力改进', 1763532900);

INSERT OR IGNORE INTO rcontact (username, nickname, remark, type) VALUES
('testuser1', '张明', '项目经理', 1),
('testuser2', '李华', '产品经理', 1),
('testuser3', '王经理', '市场总监', 1),
('testuser4', '赵工', '技术专家', 1),
('filehelper', '文件传输助手', '系统', 33);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_message_time ON message(createTime);
CREATE INDEX IF NOT EXISTS idx_message_talker ON message(talkerId);