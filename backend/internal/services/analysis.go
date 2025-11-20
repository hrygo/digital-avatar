package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"twin-os/backend/internal/config"
	"twin-os/backend/internal/models"
	"twin-os/backend/internal/repository"
	"twin-os/backend/pkg/ai"
	"twin-os/backend/pkg/crypto"
	"twin-os/backend/pkg/logger"
)

// AnalysisService 分析服务
type AnalysisService struct {
	config     *config.Config
	repository *repository.Repository
	aiClient   *ai.DeepSeekClient
	detector   *crypto.PIIDetector
}

// NewAnalysisService 创建分析服务
func NewAnalysisService(cfg *config.Config, repo *repository.Repository) *AnalysisService {
	aiClient := ai.NewDeepSeekClient(cfg.DeepSeekAPIKey, cfg.DeepSeekBaseURL)

	return &AnalysisService{
		config:     cfg,
		repository: repo,
		aiClient:   aiClient,
		detector:   crypto.NewPIIDetector(),
	}
}

// GenerateBriefing 生成今日情报简报
func (s *AnalysisService) GenerateBriefing() (*BriefingResult, error) {
	result := &BriefingResult{
		GeneratedAt: time.Now(),
	}

	// 获取最近24小时的消息
	messages, err := s.repository.Message.GetRecentMessages(24 * time.Hour, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent messages: %w", err)
	}

	if len(messages) == 0 {
		result.Content = "今日暂无新消息。"
		result.MessageCount = 0
		return result, nil
	}

	logger.Info(fmt.Sprintf("📊 Generating briefing from %d messages", len(messages)))

	// 准备消息内容
	var messageTexts []string
	for _, msg := range messages {
		if msg.Content != "" {
			maskedContent, _ := s.maskPII(msg.Content)
			messageTexts = append(messageTexts, maskedContent)
		}
	}

	if len(messageTexts) == 0 {
		result.Content = "今日暂无文本消息。"
		result.MessageCount = len(messages)
		return result, nil
	}

	// 使用AI生成简报
	briefing, err := s.aiClient.GenerateBriefing(messageTexts)
	if err != nil {
		logger.Warn("⚠️ Failed to generate AI briefing: " + err.Error())
		// 生成简化版简报
		briefing = s.generateSimpleBriefing(messages)
	}

	// 保存分析结果
	analysis := &models.AnalysisResult{
		Type:       "briefing",
		MessageIDs: s.extractMessageIDs(messages),
		Content:    briefing,
		Confidence: 0.8,
		CreatedAt:  time.Now(),
	}

	if err := s.repository.Analysis.Create(analysis); err != nil {
		logger.Warn("⚠️ Failed to save briefing analysis: " + err.Error())
	}

	result.Content = briefing
	result.MessageCount = len(messages)
	result.AnalysisID = analysis.ID

	logger.Info("✅ Briefing generated successfully")
	return result, nil
}

// ExtractTodos 提取待办事项
func (s *AnalysisService) ExtractTodos() (*TodoResult, error) {
	result := &TodoResult{
		GeneratedAt: time.Now(),
	}

	// 获取最近48小时的消息
	messages, err := s.repository.Message.GetRecentMessages(48*time.Hour, 200)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent messages: %w", err)
	}

	if len(messages) == 0 {
		result.Todos = []ai.TodoItem{}
		result.MessageCount = 0
		return result, nil
	}

	logger.Info(fmt.Sprintf("📋 Extracting todos from %d messages", len(messages)))

	// 准备消息内容
	var messageTexts []string
	messageMap := make(map[string]*models.Message)
	for _, msg := range messages {
		if msg.Content != "" {
			maskedContent, _ := s.maskPII(msg.Content)
			messageTexts = append(messageTexts, maskedContent)
			messageMap[msg.MessageID] = msg
		}
	}

	if len(messageTexts) == 0 {
		result.Todos = []ai.TodoItem{}
		result.MessageCount = len(messages)
		return result, nil
	}

	// 使用AI提取待办
	todos, err := s.aiClient.ExtractTodos(messageTexts)
	if err != nil {
		logger.Warn("⚠️ Failed to extract AI todos: " + err.Error())
		// 简单关键词提取
		todos = s.extractTodosByKeywords(messages)
	}

	// 过滤和增强待办项
	result.Todos = s.filterAndEnhanceTodos(todos, messageMap)
	result.MessageCount = len(messages)

	// 保存分析结果
	if len(result.Todos) > 0 {
		todosJSON, _ := json.Marshal(result.Todos)
		analysis := &models.AnalysisResult{
			Type:       "todo",
			MessageIDs: s.extractMessageIDs(messages),
			Content:    string(todosJSON),
			Confidence: 0.7,
			CreatedAt:  time.Now(),
		}

		if err := s.repository.Analysis.Create(analysis); err != nil {
			logger.Warn("⚠️ Failed to save todo analysis: " + err.Error())
		}

		result.AnalysisID = analysis.ID
	}

	logger.Info(fmt.Sprintf("✅ Extracted %d todos", len(result.Todos)))
	return result, nil
}

// AnalyzeConnections 分析人脉关系
func (s *AnalysisService) AnalyzeConnections() (*ConnectionResult, error) {
	result := &ConnectionResult{
		GeneratedAt: time.Now(),
	}

	// 获取最近一周的消息
	messages, err := s.repository.Message.GetRecentMessages(7*24*time.Hour, 300)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent messages: %w", err)
	}

	if len(messages) == 0 {
		result.Connections = []ai.ConnectionAnalysis{}
		result.MessageCount = 0
		return result, nil
	}

	logger.Info(fmt.Sprintf("👥 Analyzing connections from %d messages", len(messages)))

	// 获取联系人列表
	contacts, err := s.repository.Contact.GetAll()
	if err != nil {
		logger.Warn("⚠️ Failed to get contacts: " + err.Error())
		contacts = []*models.Contact{}
	}

	// 准备联系人姓名列表
	var contactNames []string
	for _, contact := range contacts {
		if contact.NickName != "" {
			contactNames = append(contactNames, contact.NickName)
		}
	}

	// 准备消息内容
	var messageTexts []string
	for _, msg := range messages {
		if msg.Content != "" {
			maskedContent, _ := s.maskPII(msg.Content)
			messageTexts = append(messageTexts, maskedContent)
		}
	}

	if len(messageTexts) == 0 {
		result.Connections = []ai.ConnectionAnalysis{}
		result.MessageCount = len(messages)
		return result, nil
	}

	// 使用AI分析人脉关系
	connections, err := s.aiClient.AnalyzeConnections(messageTexts, contactNames)
	if err != nil {
		logger.Warn("⚠️ Failed to analyze AI connections: " + err.Error())
		// 简单关键词分析
		connections = s.analyzeConnectionsByKeywords(messages, contacts)
	}

	result.Connections = connections
	result.MessageCount = len(messages)

	// 保存分析结果
	if len(result.Connections) > 0 {
		connectionsJSON, _ := json.Marshal(connections)
		analysis := &models.AnalysisResult{
			Type:       "connection",
			MessageIDs: s.extractMessageIDs(messages),
			Content:    string(connectionsJSON),
			Confidence: 0.6,
			CreatedAt:  time.Now(),
		}

		if err := s.repository.Analysis.Create(analysis); err != nil {
			logger.Warn("⚠️ Failed to save connection analysis: " + err.Error())
		}

		result.AnalysisID = analysis.ID
	}

	logger.Info(fmt.Sprintf("✅ Analyzed %d connections", len(result.Connections)))
	return result, nil
}

// maskPII 对文本进行PII脱敏
func (s *AnalysisService) maskPII(text string) (string, crypto.PIIReport) {
	return s.detector.DetectAndReplace(text)
}

// extractMessageIDs 提取消息ID列表
func (s *AnalysisService) extractMessageIDs(messages []*models.Message) string {
	var ids []string
	for _, msg := range messages {
		ids = append(ids, msg.MessageID)
	}
	result, _ := json.Marshal(ids)
	return string(result)
}

// generateSimpleBriefing 生成简化版简报
func (s *AnalysisService) generateSimpleBriefing(messages []*models.Message) string {
	briefing := "# 今日情报简报\n\n"

	// 统计信息
	briefing += fmt.Sprintf("- 处理消息：%d 条\n", len(messages))

	// 统计活跃联系人
	talkerCounts := make(map[string]int)
	for _, msg := range messages {
		talkerCounts[msg.TalkerID]++
	}

	if len(talkerCounts) > 0 {
		briefing += "- 活跃联系人：" + fmt.Sprintf("%d 位\n", len(talkerCounts))
	}

	briefing += "\n## 重要动态\n\n系统正在处理消息，请稍后查看详细分析。"

	return briefing
}

// extractTodosByKeywords 通过关键词提取待办
func (s *AnalysisService) extractTodosByKeywords(messages []*models.Message) []ai.TodoItem {
	var todos []ai.TodoItem

	// 待办关键词
	todoKeywords := []string{
		"请", "需要", "要求", "安排", "会议", "报告", "提交", "完成",
		"deadline", "截止", "明天", "下周", "尽快", "urgent", "紧急",
	}

	for _, msg := range messages {
		if msg.Content == "" {
			continue
		}

		content := strings.ToLower(msg.Content)
		for _, keyword := range todoKeywords {
			if strings.Contains(content, strings.ToLower(keyword)) {
				todo := ai.TodoItem{
					Title:           "待处理事项",
					Description:     s.truncateString(msg.Content, 100),
					Priority:        "medium",
					Deadline:        "",
					RelatedPeople:   []string{},
					SourceMessageID: msg.MessageID,
					CreatedAt:       msg.CreatedAt.Format(time.RFC3339),
				}

				// 简单优先级判断
				if strings.Contains(content, "紧急") || strings.Contains(content, "urgent") {
					todo.Priority = "high"
				}

				todos = append(todos, todo)
				break // 每条消息最多提取一个待办
			}
		}
	}

	return todos
}

// analyzeConnectionsByKeywords 通过关键词分析人脉
func (s *AnalysisService) analyzeConnectionsByKeywords(messages []*models.Message, contacts []*models.Contact) []ai.ConnectionAnalysis {
	connections := make(map[string]*ai.ConnectionAnalysis)

	// 分析消息中的人名提及
	for _, msg := range messages {
		if msg.Content == "" {
			continue
		}

		content := strings.ToLower(msg.Content)

		// 简化的人名检测
		for _, contact := range contacts {
			contactName := strings.ToLower(contact.NickName)
			if contactName != "" && strings.Contains(content, contactName) {
				if _, exists := connections[contact.NickName]; !exists {
					connections[contact.NickName] = &ai.ConnectionAnalysis{
						Person:       contact.NickName,
						Action:       "提及",
						Context:      s.truncateString(msg.Content, 50),
						Importance:   "medium",
						Sentiment:    "neutral",
						MessageCount: 0,
						LastInteraction:  msg.Timestamp.Format(time.RFC3339),
					}
				}

				conn := connections[contact.NickName]
				conn.MessageCount++
				conn.LastInteraction = msg.Timestamp.Format(time.RFC3339)
			}
		}
	}

	// 转换为切片
	var result []ai.ConnectionAnalysis
	for _, conn := range connections {
		result = append(result, *conn)
	}

	return result
}

// filterAndEnhanceTodos 过滤和增强待办项
func (s *AnalysisService) filterAndEnhanceTodos(todos []ai.TodoItem, messageMap map[string]*models.Message) []ai.TodoItem {
	var filtered []ai.TodoItem

	for _, todo := range todos {
		// 过滤无效待办
		if strings.TrimSpace(todo.Title) == "" {
			continue
		}

		// 增强待办项信息
		if todo.SourceMessageID != "" {
			if msg, exists := messageMap[todo.SourceMessageID]; exists {
				// 可以根据原消息增强待办项信息
				if todo.Description == "" {
					todo.Description = s.truncateString(msg.Content, 100)
				}
			}
		}

		// 设置默认值
		if todo.Priority == "" {
			todo.Priority = "medium"
		}
		if todo.CreatedAt == "" {
			todo.CreatedAt = time.Now().Format(time.RFC3339)
		}

		filtered = append(filtered, todo)
	}

	return filtered
}

// truncateString 截断字符串
func (s *AnalysisService) truncateString(str string, maxLen int) string {
	if len(str) <= maxLen {
		return str
	}
	return str[:maxLen] + "..."
}

// BriefingResult 简报结果
type BriefingResult struct {
	Content     string    `json:"content"`
	MessageCount int       `json:"message_count"`
	GeneratedAt time.Time `json:"generated_at"`
	AnalysisID  int64     `json:"analysis_id,omitempty"`
}

// TodoResult 待办结果
type TodoResult struct {
	Todos        []ai.TodoItem `json:"todos"`
	MessageCount int           `json:"message_count"`
	GeneratedAt  time.Time     `json:"generated_at"`
	AnalysisID   int64         `json:"analysis_id,omitempty"`
}

// ConnectionResult 人脉分析结果
type ConnectionResult struct {
	Connections  []ai.ConnectionAnalysis `json:"connections"`
	MessageCount int                     `json:"message_count"`
	GeneratedAt  time.Time               `json:"generated_at"`
	AnalysisID   int64                   `json:"analysis_id,omitempty"`
}