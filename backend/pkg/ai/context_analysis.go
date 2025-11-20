package ai

import (
	"strings"
)

// ContextAnalyzer 上下文分析器
type ContextAnalyzer struct {
	businessKeywords []string
	socialKeywords   []string
	urgentKeywords   []string
	decisionKeywords []string
}

// NewContextAnalyzer 创建上下文分析器
func NewContextAnalyzer() *ContextAnalyzer {
	return &ContextAnalyzer{
		businessKeywords: []string{"项目", "任务", "工作", "会议", "报告", "客户"},
		socialKeywords:   []string{"朋友", "聚会", "聊天", "吃饭", "玩", "活动"},
		urgentKeywords:   []string{"紧急", "马上", "立即", "尽快", "急"},
		decisionKeywords: []string{"决定", "选择", "判断", "考虑", "想", "觉得"},
	}
}

// AnalyzeContext 分析文本上下文
func (ca *ContextAnalyzer) AnalyzeContext(text string) string {
	context := ""
	text = strings.ToLower(text)

	// 检测业务场景
	if ca.containsKeywords(text, ca.businessKeywords) {
		context += "工作场景 "
	}

	// 检测社交场景
	if ca.containsKeywords(text, ca.socialKeywords) {
		context += "社交场景 "
	}

	// 检测紧急程度
	if ca.containsKeywords(text, ca.urgentKeywords) {
		context += "紧急 "
	}

	// 检测决策需求
	if ca.containsKeywords(text, ca.decisionKeywords) {
		context += "决策需求 "
	}

	if context == "" {
		context = "一般对话"
	}

	return strings.TrimSpace(context)
}

// DeterminePriority 确定优先级
func (ca *ContextAnalyzer) DeterminePriority(intent IntentType, score float64, timeExpression string) string {
	// 基于意图类型的优先级
	intentPriority := map[IntentType]string{
		IntentTask:      "high",
		IntentDecision:  "high",
		IntentRequest:   "high",
		IntentComplaint: "high",
		IntentEmotional: "medium",
		IntentQuestion:  "medium",
		IntentSocial:    "medium",
		IntentOffer:     "medium",
		IntentGratitude: "low",
		IntentGreeting:  "low",
		IntentFarewell:  "low",
		IntentInfo:      "low",
		IntentNeutral:   "low",
	}

	priority, exists := intentPriority[intent]
	if !exists {
		priority = "medium"
	}

	// 基于时间表达调整优先级
	if timeExpression != "" {
		if strings.Contains(timeExpression, "今天") || strings.Contains(timeExpression, "现在") {
			priority = "high"
		} else if strings.Contains(timeExpression, "明天") || strings.Contains(timeExpression, "下周") {
			priority = "medium"
		}
	}

	// 基于得分调整
	if score >= 5.0 {
		if priority == "low" {
			priority = "medium"
		}
	} else if score >= 8.0 {
		if priority == "medium" {
			priority = "high"
		}
	}

	return priority
}

// containsKeywords 检查文本是否包含关键词
func (ca *ContextAnalyzer) containsKeywords(text string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(text, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}