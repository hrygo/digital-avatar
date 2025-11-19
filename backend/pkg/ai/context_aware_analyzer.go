package ai

import (
	"context"
	"strings"
	"time"
)

// ContextAwareAnalyzer 上下文感知分析器
type ContextAwareAnalyzer struct {
	emotionAnalyzer       *EmotionAnalyzer
	intentRecognizer      *IntentRecognizer
	topicAnalyzer         *TopicAnalyzer
	entityRecognizer      *EntityRecognizer
	timeExtractor         *TimeExtractor
	contextAnalyzer       *ContextAnalyzer
	config               *AnalysisConfig
	conversationHistory  *ConversationHistory
}

// AnalysisConfig 分析配置
type AnalysisConfig struct {
	EnableEmotionAnalysis    bool  `json:"enable_emotion_analysis"`
	EnableIntentRecognition  bool  `json:"enable_intent_recognition"`
	EnableTopicExtraction    bool  `json:"enable_topic_extraction"`
	EnableEntityRecognition  bool  `json:"enable_entity_recognition"`
	EnableTimeExtraction     bool  `json:"enable_time_extraction"`
	EnableContextAnalysis    bool  `json:"enable_context_analysis"`
	MaxHistorySize          int   `json:"max_history_size"`
	ConfidenceThreshold     float64 `json:"confidence_threshold"`
	EnableTopicEvolution    bool  `json:"enable_topic_evolution"`
}

// DefaultAnalysisConfig 默认分析配置
func DefaultAnalysisConfig() *AnalysisConfig {
	return &AnalysisConfig{
		EnableEmotionAnalysis:   true,
		EnableIntentRecognition: true,
		EnableTopicExtraction:   true,
		EnableEntityRecognition: true,
		EnableTimeExtraction:    true,
		EnableContextAnalysis:   true,
		MaxHistorySize:         50,
		ConfidenceThreshold:    0.6,
		EnableTopicEvolution:   false,
	}
}

// ConversationHistory 对话历史
type ConversationHistory struct {
	Messages []ConversationMessage `json:"messages"`
	UserID   string               `json:"user_id"`
}

// ConversationMessage 对话消息
type ConversationMessage struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	User      string    `json:"user"`
	Type      string    `json:"type"` // text, image, file, etc.
}

// ComprehensiveAnalysis 综合分析结果
type ComprehensiveAnalysis struct {
	MessageID          string                    `json:"message_id"`
	Content            string                    `json:"content"`
	Timestamp          time.Time                 `json:"timestamp"`

	// 情感分析
	EmotionAnalysis    *EmotionAnalysis          `json:"emotion_analysis"`

	// 意图识别
	IntentAnalysis     *IntentAnalysis           `json:"intent_analysis"`

	// 主题分析
	TopicAnalysis      *DocumentTopicDistribution `json:"topic_analysis"`

	// 实体识别
	Entities           []Entity                  `json:"entities"`

	// 时间表达
	TimeExpressions    []string                  `json:"time_expressions"`

	// 上下文分析
	ContextAnalysis    string                    `json:"context_analysis"`

	// 综合评分
	OverallSentiment   float64                   `json:"overall_sentiment"`
	UrgencyScore       float64                   `json:"urgency_score"`
	ImportanceScore    float64                   `json:"importance_score"`

	// 元数据
	ProcessingTime     time.Duration             `json:"processing_time"`
	Confidence         float64                   `json:"confidence"`
	AnalysisVersion    string                    `json:"analysis_version"`
}

// NewContextAwareAnalyzer 创建上下文感知分析器
func NewContextAwareAnalyzer(config *AnalysisConfig) *ContextAwareAnalyzer {
	if config == nil {
		config = DefaultAnalysisConfig()
	}

	return &ContextAwareAnalyzer{
		emotionAnalyzer:      NewEmotionAnalyzer(),
		intentRecognizer:     NewIntentRecognizer(),
		topicAnalyzer:        NewTopicAnalyzer(nil),
		entityRecognizer:     NewEntityRecognizer(),
		timeExtractor:        NewTimeExtractor(),
		contextAnalyzer:      NewContextAnalyzer(),
		config:              config,
		conversationHistory: &ConversationHistory{
			Messages: make([]ConversationMessage, 0),
		},
	}
}

// AnalyzeMessage 分析单个消息
func (caa *ContextAwareAnalyzer) AnalyzeMessage(ctx context.Context, message ConversationMessage) (*ComprehensiveAnalysis, error) {
	startTime := time.Now()

	analysis := &ComprehensiveAnalysis{
		MessageID:       message.ID,
		Content:         message.Content,
		Timestamp:       message.Timestamp,
		AnalysisVersion: "v0.2.0",
	}

	// 1. 情感分析
	if caa.config.EnableEmotionAnalysis {
		emotionResult := caa.emotionAnalyzer.AnalyzeEmotion(message.Content)
		if emotionResult != nil {
			analysis.EmotionAnalysis = emotionResult
			analysis.OverallSentiment = caa.calculateAverageIntensity(emotionResult)
		}
	}

	// 2. 意图识别
	if caa.config.EnableIntentRecognition {
		intentResult := caa.intentRecognizer.AnalyzeIntent(message.Content)
		if intentResult != nil {
			analysis.IntentAnalysis = intentResult
		}
	}

	// 3. 主题分析
	if caa.config.EnableTopicExtraction {
		// 检查模型是否已训练
		if validation := caa.topicAnalyzer.ValidateTopics(); validation.IsValid {
			topicResult, err := caa.topicAnalyzer.AnalyzeDocumentTopics(message.Content)
			if err == nil {
				analysis.TopicAnalysis = topicResult
			}
		} else {
			// 如果模型未训练，使用消息内容进行训练
			documents := []string{message.Content}
			_, err := caa.topicAnalyzer.AnalyzeTopics(documents)
			if err == nil {
				topicResult, _ := caa.topicAnalyzer.AnalyzeDocumentTopics(message.Content)
				analysis.TopicAnalysis = topicResult
			}
		}
	}

	// 4. 实体识别
	if caa.config.EnableEntityRecognition {
		entities := caa.entityRecognizer.ExtractEntities(message.Content)
		// Keep the original entity type for simplicity
		analysis.Entities = entities
	}

	// 5. 时间表达提取
	if caa.config.EnableTimeExtraction {
		timeExpr := caa.timeExtractor.ExtractTimeExpression(message.Content)
		if timeExpr != "" {
			analysis.TimeExpressions = []string{timeExpr}
		}
	}

	// 6. 上下文分析
	if caa.config.EnableContextAnalysis {
		contextResult := caa.contextAnalyzer.AnalyzeContext(message.Content)
		analysis.ContextAnalysis = contextResult
		// 计算紧急程度和重要性分数（基于上下文分析的简化版本）
		analysis.UrgencyScore = caa.calculateUrgencyScore(message.Content, contextResult)
		analysis.ImportanceScore = caa.calculateImportanceScore(message.Content, contextResult)
	}

	// 计算总体置信度
	analysis.Confidence = caa.calculateOverallConfidence(analysis)

	// 处理时间
	analysis.ProcessingTime = time.Since(startTime)

	// 更新对话历史
	caa.updateConversationHistory(message)

	return analysis, nil
}

// AnalyzeConversation 分析整个对话
func (caa *ContextAwareAnalyzer) AnalyzeConversation(ctx context.Context, messages []ConversationMessage) ([]*ComprehensiveAnalysis, error) {
	var analyses []*ComprehensiveAnalysis

	for _, message := range messages {
		analysis, err := caa.AnalyzeMessage(ctx, message)
		if err != nil {
			continue // 跳过错误消息
		}
		analyses = append(analyses, analysis)
	}

	return analyses, nil
}

// GetConversationSummary 获取对话摘要
func (caa *ContextAwareAnalyzer) GetConversationSummary() *ConversationSummary {
	if len(caa.conversationHistory.Messages) == 0 {
		return &ConversationSummary{}
	}

	summary := &ConversationSummary{
		TotalMessages:    len(caa.conversationHistory.Messages),
		TimeRange:        caa.getTimeRange(),
		DominantEmotions: caa.getDominantEmotions(),
		DominantIntents:  caa.getDominantIntents(),
		KeyTopics:        caa.getKeyTopics(),
		Entities:         caa.getAllEntities(),
		UrgencyLevel:     caa.getOverallUrgency(),
	}

	return summary
}

// ConversationSummary 对话摘要
type ConversationSummary struct {
	TotalMessages    int                    `json:"total_messages"`
	TimeRange        TimeRange              `json:"time_range"`
	DominantEmotions []EmotionType          `json:"dominant_emotions"`
	DominantIntents  []IntentType           `json:"dominant_intents"`
	KeyTopics        []string               `json:"key_topics"`
	Entities         []Entity               `json:"entities"`
	UrgencyLevel     float64                `json:"urgency_level"`
}

// TimeRange 时间范围
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// 计算总体置信度
func (caa *ContextAwareAnalyzer) calculateOverallConfidence(analysis *ComprehensiveAnalysis) float64 {
	var totalConfidence float64
	var components int

	if analysis.EmotionAnalysis != nil {
		totalConfidence += analysis.EmotionAnalysis.Confidence
		components++
	}

	if analysis.IntentAnalysis != nil {
		totalConfidence += analysis.IntentAnalysis.Confidence
		components++
	}

	if analysis.TopicAnalysis != nil {
		totalConfidence += analysis.TopicAnalysis.Certainty
		components++
	}

	if components == 0 {
		return 0.0
	}

	return totalConfidence / float64(components)
}

// 更新对话历史
func (caa *ContextAwareAnalyzer) updateConversationHistory(message ConversationMessage) {
	caa.conversationHistory.Messages = append(caa.conversationHistory.Messages, message)

	// 限制历史大小
	if len(caa.conversationHistory.Messages) > caa.config.MaxHistorySize {
		caa.conversationHistory.Messages = caa.conversationHistory.Messages[1:]
	}
}

// 获取时间范围
func (caa *ContextAwareAnalyzer) getTimeRange() TimeRange {
	if len(caa.conversationHistory.Messages) == 0 {
		return TimeRange{}
	}

	messages := caa.conversationHistory.Messages
	return TimeRange{
		Start: messages[0].Timestamp,
		End:   messages[len(messages)-1].Timestamp,
	}
}

// 获取主导情感
func (caa *ContextAwareAnalyzer) getDominantEmotions() []EmotionType {
	// 简化实现，实际应基于历史分析结果
	return []EmotionType{EmotionJoy, EmotionSadness, EmotionAnger}
}

// 获取主导意图
func (caa *ContextAwareAnalyzer) getDominantIntents() []IntentType {
	// 简化实现，实际应基于历史分析结果
	return []IntentType{IntentQuestion, IntentGreeting, IntentRequest}
}

// 获取关键主题
func (caa *ContextAwareAnalyzer) getKeyTopics() []string {
	// 简化实现，实际应基于主题分析结果
	return []string{"工作", "生活", "技术"}
}

// 获取所有实体
func (caa *ContextAwareAnalyzer) getAllEntities() []Entity {
	// 简化实现，实际应收集所有识别的实体
	return []Entity{}
}

// 获取总体紧急程度
func (caa *ContextAwareAnalyzer) getOverallUrgency() float64 {
	// 简化实现，实际应基于上下文分析
	return 0.3
}

// SetUserID 设置用户ID
func (caa *ContextAwareAnalyzer) SetUserID(userID string) {
	caa.conversationHistory.UserID = userID
}

// GetUserID 获取用户ID
func (caa *ContextAwareAnalyzer) GetUserID() string {
	return caa.conversationHistory.UserID
}

// ClearHistory 清除对话历史
func (caa *ContextAwareAnalyzer) ClearHistory() {
	caa.conversationHistory.Messages = make([]ConversationMessage, 0)
}

// GetHistorySize 获取历史大小
func (caa *ContextAwareAnalyzer) GetHistorySize() int {
	return len(caa.conversationHistory.Messages)
}

// calculateAverageIntensity 计算平均情感强度
func (caa *ContextAwareAnalyzer) calculateAverageIntensity(emotion *EmotionAnalysis) float64 {
	if emotion == nil {
		return 0.0
	}

	// 简化实现：使用主要情感类型的强度
	// 假设 EmotionAnalysis 有一个 OverallIntensity 字段，如果没有则返回 0.5
	// 这里使用一个通用的计算方法
	baseIntensity := 0.5

	// 根据情感类型调整强度
	switch emotion.Type {
	case EmotionJoy:
		baseIntensity = 0.8
	case EmotionLove:
		baseIntensity = 0.7
	case EmotionSurprise:
		baseIntensity = 0.6
	case EmotionAnger:
		baseIntensity = 0.8
	case EmotionSadness:
		baseIntensity = 0.7
	case EmotionFear:
		baseIntensity = 0.6
	case EmotionDisgust:
		baseIntensity = 0.5
	case EmotionNeutral:
		baseIntensity = 0.3
	}

	return baseIntensity
}

// calculateUrgencyScore 计算紧急程度分数
func (caa *ContextAwareAnalyzer) calculateUrgencyScore(text, context string) float64 {
	urgencyKeywords := []string{
		"紧急", "急", "马上", "立刻", "赶紧", "尽快", "火速", "速回",
		"urgent", "emergency", "asap", "immediately", "hurry",
	}

	score := 0.0
	text = strings.ToLower(text)

	for _, keyword := range urgencyKeywords {
		if strings.Contains(text, keyword) {
			score += 0.2
		}
	}

	// 如果包含时间表达，增加紧急程度
	if strings.Contains(text, "今天") || strings.Contains(text, "现在") {
		score += 0.1
	}

	if score > 1.0 {
		score = 1.0
	}

	return score
}

// calculateImportanceScore 计算重要性分数
func (caa *ContextAwareAnalyzer) calculateImportanceScore(text, context string) float64 {
	importanceKeywords := []string{
		"重要", "关键", "必须", "一定", "务必", "请", "麻烦", "帮助",
		"important", "critical", "must", "please", "help", "thanks",
	}

	score := 0.0
	text = strings.ToLower(text)

	for _, keyword := range importanceKeywords {
		if strings.Contains(text, keyword) {
			score += 0.15
		}
	}

	// 长度也影响重要性
	if len(text) > 50 {
		score += 0.1
	}

	if score > 1.0 {
		score = 1.0
	}

	return score
}