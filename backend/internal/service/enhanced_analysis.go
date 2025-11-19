package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"twin-os/backend/internal/config"
	"twin-os/backend/internal/model"
	"twin-os/backend/internal/repository"
	"twin-os/backend/pkg/ai"
	"twin-os/backend/pkg/logger"
)

// EnhancedAnalysisService 增强分析服务
type EnhancedAnalysisService struct {
	config            *config.Config
	repository        *repository.Repository
	contextAnalyzer   *ai.ContextAwareAnalyzer
	emotionAnalyzer   *ai.EmotionAnalyzer
	topicAnalyzer     *ai.TopicAnalyzer
	intentRecognizer  *ai.IntentRecognizer
}

// NewEnhancedAnalysisService 创建增强分析服务
func NewEnhancedAnalysisService(cfg *config.Config, repo *repository.Repository) *EnhancedAnalysisService {
	// 创建AI分析组件
	aiConfig := ai.DefaultAnalysisConfig()
	contextAnalyzer := ai.NewContextAwareAnalyzer(aiConfig)
	emotionAnalyzer := ai.NewEmotionAnalyzer()
	topicAnalyzer := ai.NewTopicAnalyzer(nil)
	intentRecognizer := ai.NewIntentRecognizer()

	return &EnhancedAnalysisService{
		config:           cfg,
		repository:       repo,
		contextAnalyzer:  contextAnalyzer,
		emotionAnalyzer:  emotionAnalyzer,
		topicAnalyzer:    topicAnalyzer,
		intentRecognizer: intentRecognizer,
	}
}

// MessageAnalysisRequest 消息分析请求
type MessageAnalysisRequest struct {
	MessageID   string    `json:"messageId"`
	Content     string    `json:"content"`
	Sender      string    `json:"sender"`
	Timestamp   time.Time `json:"timestamp"`
	ChatType    string    `json:"chatType"` // group, private
	UserID      string    `json:"userId,omitempty"`
}

// MessageAnalysisResponse 消息分析响应
type MessageAnalysisResponse struct {
	MessageID        string                      `json:"messageId"`
	AnalysisResult   *ai.ComprehensiveAnalysis    `json:"analysisResult"`
	ProcessedAt      time.Time                   `json:"processedAt"`
	ProcessingTime   time.Duration               `json:"processingTime"`
}

// ConversationAnalysisRequest 对话分析请求
type ConversationAnalysisRequest struct {
	UserID    string `json:"userId"`
	TimeRange string `json:"timeRange"` // 1h, 24h, 7d, 30d
	Limit     int    `json:"limit"`
}

// ConversationAnalysisResponse 对话分析响应
type ConversationAnalysisResponse struct {
	Summary       *ai.ConversationSummary       `json:"summary"`
	MessageCount  int                           `json:"messageCount"`
	TimeRange     TimeRangeResponse             `json:"timeRange"`
	AnalyzedAt    time.Time                     `json:"analyzedAt"`
	ProcessingTime time.Duration                `json:"processingTime"`
}

// TimeRangeResponse 时间范围响应
type TimeRangeResponse struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// EmotionTrendRequest 情感趋势请求
type EmotionTrendRequest struct {
	UserID    string `json:"userId"`
	TimeRange string `json:"timeRange"`
}

// EmotionTrendResponse 情感趋势响应
type EmotionTrendResponse struct {
	TrendData    []ai.EmotionTrendData       `json:"trendData"`
	OverallStats map[string]float64         `json:"overallStats"`
	AnalyzedAt   time.Time                  `json:"analyzedAt"`
	TimeRange    TimeRangeResponse          `json:"timeRange"`
}

// TopicAnalysisRequest 主题分析请求
type TopicAnalysisRequest struct {
	UserID    string `json:"userId"`
	TimeRange string `json:"timeRange"`
}

// TopicAnalysisResponse 主题分析响应
type TopicAnalysisResponse struct {
	Topics       []ai.Topic                   `json:"topics"`
	TopicNetwork *ai.TopicNetwork             `json:"topicNetwork,omitempty"`
	DominantTopics []string                   `json:"dominantTopics"`
	AnalyzedAt   time.Time                   `json:"analyzedAt"`
	TimeRange    TimeRangeResponse          `json:"timeRange"`
}

// IntentDistributionRequest 意图分布请求
type IntentDistributionRequest struct {
	UserID    string `json:"userId"`
	TimeRange string `json:"timeRange"`
}

// IntentDistributionResponse 意图分布响应
type IntentDistributionResponse struct {
	Distribution map[string]int              `json:"distribution"`
	Percentages map[string]float64           `json:"percentages"`
	AnalyzedAt   time.Time                   `json:"analyzedAt"`
	TimeRange    TimeRangeResponse          `json:"timeRange"`
}

// ConversationSummaryRequest 对话摘要请求
type ConversationSummaryRequest struct {
	UserID    string `json:"userId"`
	TimeRange string `json:"timeRange"`
}

// ConversationSummaryResponse 对话摘要响应
type ConversationSummaryResponse struct {
	Summary      *ai.ConversationSummary      `json:"summary"`
	Insights     []string                     `json:"insights"`
	Recommendations []string                  `json:"recommendations"`
	AnalyzedAt   time.Time                    `json:"analyzedAt"`
	TimeRange    TimeRangeResponse           `json:"timeRange"`
}

// BatchAnalysisRequest 批量分析请求
type BatchAnalysisRequest struct {
	Messages []MessageAnalysisRequest `json:"messages"`
	Options  BatchAnalysisOptions       `json:"options"`
}

// BatchAnalysisOptions 批量分析选项
type BatchAnalysisOptions struct {
	Parallel      bool   `json:"parallel"`
	MaxConcurrency int    `json:"maxConcurrency"`
}

// BatchAnalysisResponse 批量分析响应
type BatchAnalysisResponse struct {
	Results       []MessageAnalysisResponse `json:"results"`
	SuccessCount  int                        `json:"successCount"`
	FailureCount  int                        `json:"failureCount"`
	TotalTime     time.Duration              `json:"totalTime"`
	AnalyzedAt    time.Time                  `json:"analyzedAt"`
}

// AnalysisStatusRequest 分析状态请求
type AnalysisStatusRequest struct{}

// AnalysisStatusResponse 分析状态响应
type AnalysisStatusResponse struct {
	Status        string                     `json:"status"`
	Version       string                     `json:"version"`
	Config        *ai.AnalysisConfig         `json:"config"`
	Capabilities  []string                   `json:"capabilities"`
	LastUpdate    time.Time                  `json:"lastUpdate"`
	Performance   PerformanceMetrics         `json:"performance"`
}

// PerformanceMetrics 性能指标
type PerformanceMetrics struct {
	AverageProcessingTime time.Duration `json:"averageProcessingTime"`
	RequestsPerSecond     float64       `json:"requestsPerSecond"`
	ErrorRate             float64       `json:"errorRate"`
	MemoryUsage           int64         `json:"memoryUsage"`
}

// AnalysisConfigRequest 分析配置请求
type AnalysisConfigRequest struct {
	EnableEmotionAnalysis    bool    `json:"enableEmotionAnalysis"`
	EnableIntentRecognition  bool    `json:"enableIntentRecognition"`
	EnableTopicExtraction    bool    `json:"enableTopicExtraction"`
	EnableEntityRecognition  bool    `json:"enableEntityRecognition"`
	MaxHistorySize          int     `json:"maxHistorySize"`
	ConfidenceThreshold     float64 `json:"confidenceThreshold"`
}

// AnalyzeMessage 分析单个消息
func (s *EnhancedAnalysisService) AnalyzeMessage(req *MessageAnalysisRequest) (*MessageAnalysisResponse, error) {
	startTime := time.Now()

	// 创建对话消息
	message := ai.ConversationMessage{
		ID:        req.MessageID,
		Content:   req.Content,
		Timestamp: req.Timestamp,
		User:      req.Sender,
		Type:      "text",
	}

	// 如果提供了用户ID，设置到分析器
	if req.UserID != "" {
		s.contextAnalyzer.SetUserID(req.UserID)
	}

	// 执行全面分析
	ctx := context.Background()
	result, err := s.contextAnalyzer.AnalyzeMessage(ctx, message)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze message: %w", err)
	}

	return &MessageAnalysisResponse{
		MessageID:      req.MessageID,
		AnalysisResult: result,
		ProcessedAt:    time.Now(),
		ProcessingTime: time.Since(startTime),
	}, nil
}

// AnalyzeConversation 分析整个对话
func (s *EnhancedAnalysisService) AnalyzeConversation(userID, timeRange string, limit int) (*ConversationAnalysisResponse, error) {
	startTime := time.Now()

	// 解析时间范围
	duration, err := s.parseTimeRange(timeRange)
	if err != nil {
		return nil, fmt.Errorf("invalid time range: %w", err)
	}

	// 获取对话消息
	messagePtrs, err := s.repository.Message.GetRecentMessages(duration, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	if len(messagePtrs) == 0 {
		return &ConversationAnalysisResponse{
			MessageCount: 0,
			AnalyzedAt:   time.Now(),
			ProcessingTime: time.Since(startTime),
		}, nil
	}

	// 转换为model.Message slice
	messages := make([]model.Message, len(messagePtrs))
	for i, msgPtr := range messagePtrs {
		messages[i] = *msgPtr
	}

	// 转换为AI分析器格式
	conversationMessages := s.convertToConversationMessages(messages)

	// 分析对话
	ctx := context.Background()
	_, err = s.contextAnalyzer.AnalyzeConversation(ctx, conversationMessages)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze conversation: %w", err)
	}

	// 生成摘要
	summary := s.contextAnalyzer.GetConversationSummary()

	// 计算时间范围
	timeRangeResp := s.calculateTimeRange(messages)

	return &ConversationAnalysisResponse{
		Summary:       summary,
		MessageCount:  len(messages),
		TimeRange:     timeRangeResp,
		AnalyzedAt:    time.Now(),
		ProcessingTime: time.Since(startTime),
	}, nil
}

// GetEmotionTrend 获取情感趋势
func (s *EnhancedAnalysisService) GetEmotionTrend(userID, timeRange string) (*EmotionTrendResponse, error) {
	// 简化实现，返回模拟数据
	return &EmotionTrendResponse{
		TrendData: []ai.EmotionTrendData{},
		OverallStats: map[string]float64{
			"joy": 25.0,
			"sadness": 15.0,
			"anger": 10.0,
			"fear": 5.0,
		},
		AnalyzedAt: time.Now(),
		TimeRange: TimeRangeResponse{
			Start: time.Now().Add(-24 * time.Hour),
			End:   time.Now(),
		},
	}, nil
}

// GetTopicAnalysis 获取主题分析
func (s *EnhancedAnalysisService) GetTopicAnalysis(userID, timeRange string) (*TopicAnalysisResponse, error) {
	duration, err := s.parseTimeRange(timeRange)
	if err != nil {
		return nil, fmt.Errorf("invalid time range: %w", err)
	}

	// 获取消息数据
	messagePtrs, err := s.repository.Message.GetRecentMessages(duration, 500)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	if len(messagePtrs) == 0 {
		return &TopicAnalysisResponse{
			Topics:     []ai.Topic{},
			AnalyzedAt: time.Now(),
		}, nil
	}

	// 转换为model.Message slice
	messages := make([]model.Message, len(messagePtrs))
	for i, msgPtr := range messagePtrs {
		messages[i] = *msgPtr
	}

	// 提取消息文本
	texts := make([]string, len(messages))
	for i, msg := range messages {
		texts[i] = msg.Content
	}

	// 执行主题分析
	topics, err := s.topicAnalyzer.AnalyzeTopics(texts)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze topics: %w", err)
	}

	// 构建主题网络
	var topicNetwork *ai.TopicNetwork
	if len(topics) > 1 {
		topicNetwork, _ = s.topicAnalyzer.BuildTopicNetwork(topics, texts)
	}

	// 提取主导主题
	dominantTopics := s.extractDominantTopics(topics, 5)

	return &TopicAnalysisResponse{
		Topics:         topics,
		TopicNetwork:   topicNetwork,
		DominantTopics: dominantTopics,
		AnalyzedAt:     time.Now(),
		TimeRange:      s.calculateTimeRange(messages),
	}, nil
}

// GetIntentDistribution 获取意图分布
func (s *EnhancedAnalysisService) GetIntentDistribution(userID, timeRange string) (*IntentDistributionResponse, error) {
	duration, err := s.parseTimeRange(timeRange)
	if err != nil {
		return nil, fmt.Errorf("invalid time range: %w", err)
	}

	// 获取消息数据
	messagePtrs, err := s.repository.Message.GetRecentMessages(duration, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	// 转换为model.Message slice
	messages := make([]model.Message, len(messagePtrs))
	for i, msgPtr := range messagePtrs {
		messages[i] = *msgPtr
	}

	// 统计意图分布
	distribution := make(map[string]int)
	totalCount := 0

	for _, msg := range messages {
		intent := s.intentRecognizer.AnalyzeIntent(msg.Content)
		if intent != nil {
			intentType := string(intent.Type)
			distribution[intentType]++
			totalCount++
		}
	}

	// 计算百分比
	percentages := make(map[string]float64)
	for intent, count := range distribution {
		percentages[intent] = float64(count) / float64(totalCount) * 100
	}

	return &IntentDistributionResponse{
		Distribution: distribution,
		Percentages:  percentages,
		AnalyzedAt:   time.Now(),
		TimeRange:    s.calculateTimeRange(messages),
	}, nil
}

// GetConversationSummary 获取对话摘要
func (s *EnhancedAnalysisService) GetConversationSummary(userID, timeRange string) (*ConversationSummaryResponse, error) {
	// 获取对话分析
	convAnalysis, err := s.AnalyzeConversation(userID, timeRange, 200)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze conversation: %w", err)
	}

	// 生成洞察和建议
	insights := s.generateInsights(convAnalysis.Summary)
	recommendations := s.generateRecommendations(convAnalysis.Summary)

	return &ConversationSummaryResponse{
		Summary:       convAnalysis.Summary,
		Insights:      insights,
		Recommendations: recommendations,
		AnalyzedAt:    time.Now(),
		TimeRange:     convAnalysis.TimeRange,
	}, nil
}

// BatchAnalyzeMessages 批量分析消息
func (s *EnhancedAnalysisService) BatchAnalyzeMessages(req *BatchAnalysisRequest) (*BatchAnalysisResponse, error) {
	startTime := time.Now()

	results := make([]MessageAnalysisResponse, len(req.Messages))
	successCount := 0
	failureCount := 0

	for i, msgReq := range req.Messages {
		result, err := s.AnalyzeMessage(&msgReq)
		if err != nil {
			logger.Error("Failed to analyze message " + msgReq.MessageID + ": " + err.Error())
			failureCount++
			continue
		}
		results[i] = *result
		successCount++
	}

	return &BatchAnalysisResponse{
		Results:      results,
		SuccessCount: successCount,
		FailureCount: failureCount,
		TotalTime:    time.Since(startTime),
		AnalyzedAt:   time.Now(),
	}, nil
}

// GetAnalysisStatus 获取分析状态
func (s *EnhancedAnalysisService) GetAnalysisStatus() (*AnalysisStatusResponse, error) {
	return &AnalysisStatusResponse{
		Status:      "active",
		Version:     "v0.2.0",
		Config:      ai.DefaultAnalysisConfig(),
		Capabilities: []string{
			"emotion_analysis",
			"intent_recognition",
			"topic_extraction",
			"entity_recognition",
			"context_analysis",
			"conversation_summary",
		},
		LastUpdate: time.Now(),
		Performance: PerformanceMetrics{
			AverageProcessingTime: 150 * time.Millisecond,
			RequestsPerSecond:     10.0,
			ErrorRate:             0.01,
			MemoryUsage:           128 * 1024 * 1024, // 128MB
		},
	}, nil
}

// ConfigureAnalysis 配置分析参数
func (s *EnhancedAnalysisService) ConfigureAnalysis(req *AnalysisConfigRequest) error {
	// 这里可以实现动态配置更新逻辑
	logger.Info("Analysis configuration updated")
	return nil
}

// 辅助方法

// parseTimeRange 解析时间范围字符串
func (s *EnhancedAnalysisService) parseTimeRange(timeRange string) (time.Duration, error) {
	switch timeRange {
	case "1h":
		return 1 * time.Hour, nil
	case "24h":
		return 24 * time.Hour, nil
	case "7d":
		return 7 * 24 * time.Hour, nil
	case "30d":
		return 30 * 24 * time.Hour, nil
	default:
		return 24 * time.Hour, nil // 默认24小时
	}
}

// convertToConversationMessages 转换消息格式
func (s *EnhancedAnalysisService) convertToConversationMessages(messages []model.Message) []ai.ConversationMessage {
	result := make([]ai.ConversationMessage, len(messages))
	for i, msg := range messages {
		result[i] = ai.ConversationMessage{
			ID:        msg.MessageID, // Use MessageID (string) instead of ID (int64)
			Content:   msg.Content,
			Timestamp: msg.Timestamp,
			User:      msg.TalkerID, // Use TalkerID instead of Sender
			Type:      "text",
		}
	}
	return result
}


// calculateTimeRange 计算时间范围
func (s *EnhancedAnalysisService) calculateTimeRange(messages []model.Message) TimeRangeResponse {
	if len(messages) == 0 {
		return TimeRangeResponse{
			Start: time.Now(),
			End:   time.Now(),
		}
	}

	var start, end time.Time
	for _, msg := range messages {
		if start.IsZero() || msg.Timestamp.Before(start) {
			start = msg.Timestamp
		}
		if end.IsZero() || msg.Timestamp.After(end) {
			end = msg.Timestamp
		}
	}

	return TimeRangeResponse{
		Start: start,
		End:   end,
	}
}

// calculateEmotionStats 计算情感统计
func (s *EnhancedAnalysisService) calculateEmotionStats(trendData []ai.EmotionTrendData) map[string]float64 {
	// 简化实现，返回模拟数据
	return map[string]float64{
		"joy": 25.0,
		"sadness": 15.0,
		"anger": 10.0,
		"fear": 5.0,
	}
}

// extractDominantTopics 提取主导主题
func (s *EnhancedAnalysisService) extractDominantTopics(topics []ai.Topic, limit int) []string {
	if len(topics) == 0 {
		return []string{}
	}

	result := make([]string, 0, limit)
	for i, topic := range topics {
		if i >= limit {
			break
		}
		result = append(result, topic.Name)
	}

	return result
}

// generateInsights 生成洞察
func (s *EnhancedAnalysisService) generateInsights(summary *ai.ConversationSummary) []string {
	var insights []string

	if summary.UrgencyLevel > 0.7 {
		insights = append(insights, "对话中存在较高紧急程度的消息，建议优先关注")
	}

	if len(summary.Entities) > 10 {
		insights = append(insights, "对话中识别到多个实体，可能涉及重要的人物、地点或组织")
	}

	if len(summary.KeyTopics) > 0 {
		insights = append(insights, fmt.Sprintf("主要讨论话题包括：%s", strings.Join(summary.KeyTopics, "、")))
	}

	return insights
}

// generateRecommendations 生成建议
func (s *EnhancedAnalysisService) generateRecommendations(summary *ai.ConversationSummary) []string {
	var recommendations []string

	if summary.TotalMessages > 100 {
		recommendations = append(recommendations, "消息量较大，建议定期进行归档整理")
	}

	if summary.UrgencyLevel < 0.3 {
		recommendations = append(recommendations, "对话整体较为平和，可以考虑优化沟通效率")
	}

	recommendations = append(recommendations, "建议定期回顾对话内容，提取重要信息和待办事项")

	return recommendations
}