package ai

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// IntentRecognizer 意图识别器
type IntentRecognizer struct {
	patterns      map[IntentType][]*IntentPattern
	keywords      map[IntentType][]string
	actionVerbs   []string
	entityRec     *EntityRecognizer
	timeExtract   *TimeExtractor
	contextAnal   *ContextAnalyzer
}

// NewIntentRecognizer 创建意图识别器
func NewIntentRecognizer() *IntentRecognizer {
	return &IntentRecognizer{
		patterns:    initializeIntentPatterns(),
		keywords:    initializeIntentKeywords(),
		actionVerbs: initializeActionVerbs(),
		entityRec:   NewEntityRecognizer(),
		timeExtract: NewTimeExtractor(),
		contextAnal: NewContextAnalyzer(),
	}
}

// AnalyzeIntent 分析文本意图
func (ir *IntentRecognizer) AnalyzeIntent(text string) *IntentAnalysis {
	text = strings.TrimSpace(text)
	if text == "" {
		return &IntentAnalysis{
			Type:      IntentNeutral,
			Confidence: 0.5,
			Timestamp: time.Now(),
		}
	}

	// 计算各意图类型的得分
	scores := make(map[IntentType]float64)
	for intent, patterns := range ir.patterns {
		scores[intent] = ir.calculateIntentScore(text, patterns, intent)
	}

	// 找出主要意图
	dominantIntent := IntentNeutral
	maxScore := 0.0
	for intent, score := range scores {
		if score > maxScore {
			maxScore = score
			dominantIntent = intent
		}
	}

	// 识别实体
	entities := ir.entityRec.ExtractEntities(text)

	// 提取行动词
	actionVerbs := ir.extractActionVerbs(text)

	// 提取时间表达
	timeExpression := ir.timeExtract.ExtractTimeExpression(text)

	// 确定优先级
	priority := ir.contextAnal.DeterminePriority(dominantIntent, maxScore, timeExpression)

	// 分析上下文
	context := ir.contextAnal.AnalyzeContext(text)

	// 计算置信度
	confidence := ir.calculateConfidence(scores, maxScore)

	// 分析含义和影响
	implications := ir.analyzeImplications(dominantIntent, entities, context)

	return &IntentAnalysis{
		Type:           dominantIntent,
		SubType:        ir.getSubType(dominantIntent, text),
		Confidence:     confidence,
		Entities:       entities,
		ActionVerbs:    actionVerbs,
		TimeExpression: timeExpression,
		Priority:       priority,
		Context:        context,
		Implications:   implications,
		Timestamp:      time.Now(),
	}
}

// AnalyzeImplicitIntent 分析隐式意图
func (ir *IntentRecognizer) AnalyzeImplicitIntent(context []string) *IntentAnalysis {
	if len(context) == 0 {
		return &IntentAnalysis{
			Type:      IntentNeutral,
			Confidence: 0.0,
			Timestamp: time.Now(),
		}
	}

	// 合并上下文文本
	combinedText := strings.Join(context, " ")
	analysis := ir.AnalyzeIntent(combinedText)

	// 降低隐式意图的置信度
	analysis.Confidence *= 0.8

	// 标记为隐式分析
	analysis.Context += " (隐式意图)"

	return analysis
}

// TrackIntentPattern 跟踪意图模式
func (ir *IntentRecognizer) TrackIntentPattern(analyses []*IntentAnalysis) []IntentPattern {
	if len(analyses) < 3 {
		return []IntentPattern{}
	}

	// 分析意图序列
	intentSequence := make([]IntentType, len(analyses))
	for i, analysis := range analyses {
		intentSequence[i] = analysis.Type
	}

	// 识别重复模式
	patterns := make(map[string]int)
	for i := 0; i < len(intentSequence)-1; i++ {
		pattern := fmt.Sprintf("%s->%s", intentSequence[i], intentSequence[i+1])
		patterns[pattern]++
	}

	// 提取有意义的模式
	var result []IntentPattern
	for pattern, count := range patterns {
		if count >= 2 { // 至少重复2次
			result = append(result, IntentPattern{
				Pattern:  NewRegexpWrapper(""),
				Weight:   0.0,
				Keywords: []string{},
				Entities: []string{},
				Context:  fmt.Sprintf("意图转换模式：%s", pattern),
			})
		}
	}

	return result
}

// calculateIntentScore 计算意图得分
func (ir *IntentRecognizer) calculateIntentScore(text string, patterns []*IntentPattern, intentType IntentType) float64 {
	score := 0.0
	text = strings.ToLower(text)

	for _, pattern := range patterns {
		if pattern.Pattern == nil {
			continue
		}

		compiledPattern := pattern.Pattern.CompileRegexp()
		if compiledPattern == nil {
			continue
		}

		if compiledPattern.MatchString(text) {
			score += pattern.Weight

			// 关键词匹配加分
			for _, keyword := range pattern.Keywords {
				if strings.Contains(text, keyword) {
					score += 0.5
				}
			}
		}
	}

	// 基于关键词的直接匹配
	if intentKeywords, exists := ir.keywords[intentType]; exists {
		for _, keyword := range intentKeywords {
			if strings.Contains(text, keyword) {
				score += 1.0
			}
		}
	}

	return math.Min(score, 10.0)
}

// extractActionVerbs 提取行动词
func (ir *IntentRecognizer) extractActionVerbs(text string) []string {
	text = strings.ToLower(text)
	var verbs []string

	for _, verb := range ir.actionVerbs {
		if strings.Contains(text, verb) {
			verbs = append(verbs, verb)
		}
	}

	// 去重
	unique := make(map[string]bool)
	for _, verb := range verbs {
		if !unique[verb] {
			unique[verb] = true
		}
	}

	var result []string
	for verb := range unique {
		result = append(result, verb)
	}

	return result
}

// calculateConfidence 计算置信度
func (ir *IntentRecognizer) calculateConfidence(scores map[IntentType]float64, maxScore float64) float64 {
	if len(scores) == 0 {
		return 0.0
	}

	// 计算得分分布
	total := 0.0
	for _, score := range scores {
		total += score
	}

	if total == 0 {
		return 0.0
	}

	// 计算最高得分占比
	confidence := maxScore / total

	// 确保置信度在合理范围内
	return math.Max(0.0, math.Min(1.0, confidence))
}

// analyzeImplications 分析含义和影响
func (ir *IntentRecognizer) analyzeImplications(intent IntentType, entities []Entity, context string) []string {
	var implications []string

	switch intent {
	case IntentTask:
		implications = []string{"需要执行具体动作", "可能需要资源分配", "可能涉及时间管理"}
	case IntentDecision:
		implications = []string{"需要做出选择", "可能产生后果", "需要信息支持"}
	case IntentRequest:
		implications = []string{"需要响应或协助", "可能涉及资源协调", "需要时间安排"}
	case IntentQuestion:
		implications = []string{"需要信息提供", "可能涉及知识查询", "需要专业支持"}
	case IntentComplaint:
		implications = []string{"需要问题解决", "可能涉及情绪安抚", "需要改进措施"}
	case IntentSocial:
		implications = []string{"需要人际互动", "可能涉及关系维护", "需要情感交流"}
	case IntentInfo:
		implications = []string{"需要信息共享", "可能涉及知识传递", "需要理解确认"}
	case IntentEmotional:
		implications = []string{"需要情感支持", "可能涉及情绪调节", "需要关怀关注"}
	}

	// 基于实体调整含义
	for _, entity := range entities {
		switch entity.Type {
		case "person":
			implications = append(implications, "涉及人员互动")
		case "time":
			implications = append(implications, "涉及时间安排")
		case "location":
			implications = append(implications, "涉及地点相关")
		case "organization":
			implications = append(implications, "涉及组织协调")
		}
	}

	return implications
}

// getSubType 获取子类型
func (ir *IntentRecognizer) getSubType(intent IntentType, text string) string {
	text = strings.ToLower(text)

	switch intent {
	case IntentTask:
		if strings.Contains(text, "会议") {
			return "meeting"
		} else if strings.Contains(text, "报告") {
			return "report"
		} else if strings.Contains(text, "邮件") {
			return "email"
		} else if strings.Contains(text, "电话") {
			return "call"
		}
	case IntentQuestion:
		if strings.Contains(text, "什么") || strings.Contains(text, "如何") {
			return "information"
		} else if strings.Contains(text, "为什么") || strings.Contains(text, "怎么") {
			return "explanation"
		} else if strings.Contains(text, "是否") || strings.Contains(text, "能不能") {
			return "confirmation"
		}
	}

	return ""
}