package ai

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

// EmotionType 细粒度情感类型
type EmotionType string

const (
	EmotionJoy       EmotionType = "joy"        // 喜悦
	EmotionLove      EmotionType = "love"       // 爱
	EmotionSurprise  EmotionType = "surprise"   // 惊讶
	EmotionAnger     EmotionType = "anger"      // 愤怒
	EmotionSadness   EmotionType = "sadness"    // 悲伤
	EmotionFear      EmotionType = "fear"       // 恐惧
	EmotionDisgust   EmotionType = "disgust"    // 厌恶
	EmotionNeutral   EmotionType = "neutral"    // 中性
)

// EmotionAnalysis 细粒度情感分析结果
type EmotionAnalysis struct {
	Type        EmotionType `json:"type"`        // 主要情感类型
	Intensity   float64    `json:"intensity"`   // 情感强度 0-1
	Confidence  float64    `json:"confidence"`  // 置信度 0-1
	Trend       string     `json:"trend"`       // 趋势：increasing, decreasing, stable
	KeyFactors  []string   `json:"key_factors"` // 影响因素
	Context     string     `json:"context"`     // 上下文
	Timestamp   time.Time  `json:"timestamp"`   // 时间戳
}

// EmotionTrendData 情感趋势数据点
type EmotionTrendData struct {
	Timestamp   time.Time              `json:"timestamp"`
	Emotions    map[EmotionType]float64 `json:"emotions"`    // 各情感类型的强度
	Dominant    EmotionType             `json:"dominant"`    // 主导情感
	Intensity   float64                 `json:"intensity"`   // 整体情感强度
	Valence     float64                 `json:"valence"`     // 情感极性 -1到1
	Arousal     float64                 `json:"arousal"`     // 激活度 0到1
}

// EmotionPattern 情感模式
type EmotionPattern struct {
	Type         string    `json:"type"`         // 模式类型：daily, weekly, conversation
	Description  string    `json:"description"`  // 模式描述
	Frequency    float64   `json:"frequency"`    // 出现频率
	Confidence   float64   `json:"confidence"`   // 置信度
	Triggers     []string  `json:"triggers"`     // 触发因素
	Implications  []string  `json:"implications"` // 含义和影响
}

// EmotionAnalyzer 情感分析器
type EmotionAnalyzer struct {
	emotionKeywords map[EmotionType][]string
	modifiers       map[string]float64
	contextAnalyzer *ContextAnalyzer
}

// NewEmotionAnalyzer 创建情感分析器
func NewEmotionAnalyzer() *EmotionAnalyzer {
	analyzer := &EmotionAnalyzer{
		emotionKeywords: initializeEmotionKeywords(),
		modifiers:       initializeModifiers(),
		contextAnalyzer: NewContextAnalyzer(),
	}
	return analyzer
}

// AnalyzeEmotion 分析单条消息的情感
func (ea *EmotionAnalyzer) AnalyzeEmotion(text string) *EmotionAnalysis {
	text = strings.ToLower(strings.TrimSpace(text))
	if text == "" {
		return &EmotionAnalysis{
			Type:       EmotionNeutral,
			Intensity:  0.0,
			Confidence: 1.0,
			Trend:      "stable",
			Timestamp:  time.Now(),
		}
	}

	// 计算各情感类型的得分
	scores := make(map[EmotionType]float64)
	for emotion, keywords := range ea.emotionKeywords {
		scores[emotion] = ea.calculateEmotionScore(text, keywords)
	}

	// 找出主要情感
	dominantEmotion := EmotionNeutral
	maxScore := 0.0
	for emotion, score := range scores {
		if score > maxScore {
			maxScore = score
			dominantEmotion = emotion
		}
	}

	// 计算强度和置信度
	intensity := math.Min(maxScore, 1.0)
	confidence := ea.calculateConfidence(scores, maxScore)

	// 识别影响因素
	keyFactors := ea.identifyKeyFactors(text, dominantEmotion)

	// 分析上下文
	context := ea.contextAnalyzer.AnalyzeContext(text)

	return &EmotionAnalysis{
		Type:       dominantEmotion,
		Intensity:   intensity,
		Confidence:  confidence,
		Trend:      "stable", // 需要历史数据才能分析趋势
		KeyFactors: keyFactors,
		Context:    context,
		Timestamp:  time.Now(),
	}
}

// AnalyzeEmotionTrend 分析情感趋势
func (ea *EmotionAnalyzer) AnalyzeEmotionTrend(messages []MessageData) ([]EmotionTrendData, error) {
	if len(messages) == 0 {
		return []EmotionTrendData{}, nil
	}

	// 按时间排序
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Timestamp.Before(messages[j].Timestamp)
	})

	var trendData []EmotionTrendData

	// 滑动窗口分析（每10条消息为一个数据点）
	windowSize := 10
	step := 5

	for i := 0; i < len(messages); i += step {
		end := i + windowSize
		if end > len(messages) {
			end = len(messages)
		}

		windowMessages := messages[i:end]
		if len(windowMessages) == 0 {
			continue
		}

		// 分析窗口内消息的情感
		emotions := make(map[EmotionType]float64)
		textContent := ""

		for _, msg := range windowMessages {
			analysis := ea.AnalyzeEmotion(msg.Content)
			emotions[analysis.Type] += analysis.Intensity
			textContent += msg.Content + " "
		}

		// 标准化情感得分
		for emotion := range emotions {
			emotions[emotion] /= float64(len(windowMessages))
		}

		// 计算主导情感和整体强度
		dominant, intensity := ea.getDominantEmotion(emotions)
		valence, arousal := ea.calculateValenceArousal(emotions)

		dataPoint := EmotionTrendData{
			Timestamp: windowMessages[len(windowMessages)-1].Timestamp,
			Emotions:  emotions,
			Dominant: dominant,
			Intensity: intensity,
			Valence:   valence,
			Arousal:   arousal,
		}

		trendData = append(trendData, dataPoint)
	}

	return trendData, nil
}

// DetectEmotionPatterns 检测情感模式
func (ea *EmotionAnalyzer) DetectEmotionPatterns(trendData []EmotionTrendData) ([]EmotionPattern, error) {
	if len(trendData) < 3 {
		return []EmotionPattern{}, nil
	}

	var patterns []EmotionPattern

	// 检测日常模式
	dailyPattern := ea.detectDailyPattern(trendData)
	if dailyPattern != nil {
		patterns = append(patterns, *dailyPattern)
	}

	// 检测周期性模式
	periodicPattern := ea.detectPeriodicPattern(trendData)
	if periodicPattern != nil {
		patterns = append(patterns, *periodicPattern)
	}

	// 检测情感转换模式
	transitionPattern := ea.detectTransitionPattern(trendData)
	if transitionPattern != nil {
		patterns = append(patterns, *transitionPattern)
	}

	return patterns, nil
}

// calculateEmotionScore 计算特定情感的得分
func (ea *EmotionAnalyzer) calculateEmotionScore(text string, keywords []string) float64 {
	score := 0.0
	textWords := strings.Fields(text)

	// 基于关键词匹配
	for _, keyword := range keywords {
		keywordWords := strings.Fields(keyword)
		for _, kwWord := range keywordWords {
			for _, textWord := range textWords {
				if strings.Contains(textWord, kwWord) {
					score += 1.0
				}
			}
		}
	}

	// 应用修饰词调整
	modifiedScore := score
	for word, modifier := range ea.modifiers {
		if strings.Contains(text, word) {
			modifiedScore *= modifier
		}
	}

	// 标准化得分
	// 移除对 wordCount 的依赖，因为中文文本可能被识别为一个词
	// 使用简单的缩放：假设5.0分代表最高强度
	return math.Min(modifiedScore/5.0, 1.0)
}

// calculateConfidence 计算置信度
func (ea *EmotionAnalyzer) calculateConfidence(scores map[EmotionType]float64, maxScore float64) float64 {
	// 计算得分的分布熵
	total := 0.0
	for _, score := range scores {
		total += score
	}

	if total == 0 {
		return 0.0
	}

	entropy := 0.0
	for _, score := range scores {
		if score > 0 {
			prob := score / total
			entropy -= prob * math.Log(prob)
		}
	}

	// 转换为置信度（熵越小，置信度越高）
	maxEntropy := math.Log(float64(len(scores)))
	confidence := 1.0 - (entropy / maxEntropy)

	// 考虑最高得分的影响
	scoreConfidence := maxScore
	confidence = (confidence + scoreConfidence) / 2.0

	return math.Max(0.0, math.Min(1.0, confidence))
}

// identifyKeyFactors 识别情感影响因素
func (ea *EmotionAnalyzer) identifyKeyFactors(text string, emotion EmotionType) []string {
	var factors []string

	// 基于情感类型的特定因素
	switch emotion {
	case EmotionJoy:
		joyKeywords := []string{"成功", "完成", "好", "棒", "太好了", "满意", "开心", "高兴"}
		factors = ea.extractKeywordsFromText(text, joyKeywords)
	case EmotionAnger:
		angerKeywords := []string{"糟糕", "烦", "气", "生气", "不满", "问题", "错误", "失败"}
		factors = ea.extractKeywordsFromText(text, angerKeywords)
	case EmotionSadness:
		sadKeywords := []string{"难过", "伤心", "失望", "遗憾", "可惜", "不容易", "辛苦"}
		factors = ea.extractKeywordsFromText(text, sadKeywords)
	case EmotionFear:
		fearKeywords := []string{"担心", "害怕", "紧张", "焦虑", "风险", "危险", "不确定"}
		factors = ea.extractKeywordsFromText(text, fearKeywords)
	case EmotionSurprise:
		surpriseKeywords := []string{"哇", "天啊", "没想到", "意外", "突然", "居然"}
		factors = ea.extractKeywordsFromText(text, surpriseKeywords)
	}

	// 限制因素数量
	if len(factors) > 5 {
		factors = factors[:5]
	}

	return factors
}

// extractKeywordsFromText 从文本中提取关键词
func (ea *EmotionAnalyzer) extractKeywordsFromText(text string, keywords []string) []string {
	text = strings.ToLower(text)
	var found []string

	for _, keyword := range keywords {
		if strings.Contains(text, strings.ToLower(keyword)) {
			found = append(found, keyword)
		}
	}

	return found
}

// getDominantEmotion 获取主导情感
func (ea *EmotionAnalyzer) getDominantEmotion(emotions map[EmotionType]float64) (EmotionType, float64) {
	dominant := EmotionNeutral
	maxIntensity := 0.0

	for emotion, intensity := range emotions {
		if intensity > maxIntensity {
			maxIntensity = intensity
			dominant = emotion
		}
	}

	return dominant, maxIntensity
}

// calculateValenceArousal 计算情感极性和激活度
func (ea *EmotionAnalyzer) calculateValenceArousal(emotions map[EmotionType]float64) (valence, arousal float64) {
	// 情感极性映射
	valenceMap := map[EmotionType]float64{
		EmotionJoy:      1.0,
		EmotionLove:     0.8,
		EmotionSurprise: 0.2,
		EmotionNeutral:   0.0,
		EmotionFear:     -0.3,
		EmotionSadness:  -0.5,
		EmotionDisgust:  -0.7,
		EmotionAnger:    -0.8,
	}

	// 激活度映射
	arousalMap := map[EmotionType]float64{
		EmotionAnger:    0.9,
		EmotionFear:     0.8,
		EmotionSurprise: 0.7,
		EmotionJoy:      0.6,
		EmotionSadness:  0.5,
		EmotionDisgust:  0.4,
		EmotionLove:     0.3,
		EmotionNeutral:   0.1,
	}

	// 计算加权平均
	totalIntensity := 0.0
	for _, intensity := range emotions {
		totalIntensity += intensity
	}

	if totalIntensity == 0 {
		return 0.0, 0.0
	}

	weightedValence := 0.0
	weightedArousal := 0.0

	for emotion, intensity := range emotions {
		if valence, exists := valenceMap[emotion]; exists {
			weightedValence += valence * intensity
		}
		if arousal, exists := arousalMap[emotion]; exists {
			weightedArousal += arousal * intensity
		}
	}

	return weightedValence / totalIntensity, weightedArousal / totalIntensity
}

// detectDailyPattern 检测日常情感模式
func (ea *EmotionAnalyzer) detectDailyPattern(trendData []EmotionTrendData) *EmotionPattern {
	if len(trendData) < 24 { // 需要至少一天的数据
		return nil
	}

	// 分析每天的情感分布
	hourlyEmotions := make(map[int]map[EmotionType]float64)
	for _, data := range trendData {
		hour := data.Timestamp.Hour()
		if hourlyEmotions[hour] == nil {
			hourlyEmotions[hour] = make(map[EmotionType]float64)
		}
		for emotion, intensity := range data.Emotions {
			hourlyEmotions[hour][emotion] += intensity
		}
	}

	// 寻找规律（例如：早上积极，晚上平静）
	patterns := []string{}
	for hour, emotions := range hourlyEmotions {
		dominant, _ := ea.getDominantEmotion(emotions)
		patterns = append(patterns, fmt.Sprintf("%d时: %s", hour, dominant))
	}

	return &EmotionPattern{
		Type:        "daily",
		Description: "日常情感波动模式",
		Frequency:   1.0,
		Confidence:  0.7,
		Triggers:    []string{"时间", "工作节奏", "生理周期"},
		Implications: []string{"建议在高峰时段安排重要任务", "注意低谷期的情绪调节"},
	}
}

// detectPeriodicPattern 检测周期性情感模式
func (ea *EmotionAnalyzer) detectPeriodicPattern(trendData []EmotionTrendData) *EmotionPattern {
	if len(trendData) < 7 {
		return nil
	}

	// 简化的周期性检测
	weeklyPattern := ea.detectWeeklyPattern(trendData)
	if weeklyPattern {
		return &EmotionPattern{
			Type:        "weekly",
			Description: "周度情感周期",
			Frequency:   0.8,
			Confidence:  0.6,
			Triggers:    []string{"工作日vs周末", "项目周期", "社交活动"},
			Implications: []string{"利用高峰期提高效率", "在低谷期安排休息"},
		}
	}

	return nil
}

// detectWeeklyPattern 检测周度模式
func (ea *EmotionAnalyzer) detectWeeklyPattern(trendData []EmotionTrendData) bool {
	// 简化实现：检查是否存在明显的周度规律
	weekdays := make(map[int][]float64)

	for _, data := range trendData {
		weekday := int(data.Timestamp.Weekday())
		weekdays[weekday] = append(weekdays[weekday], data.Valence)
	}

	// 检查工作日和周末的差异
	workdayAvg := 0.0
	weekendAvg := 0.0
	workdayCount := 0
	weekendCount := 0

	for day, values := range weekdays {
		avg := 0.0
		for _, v := range values {
			avg += v
		}
		if len(values) > 0 {
			avg /= float64(len(values))
		}

		if day == 0 || day == 6 { // 周末
			weekendAvg += avg
			weekendCount++
		} else { // 工作日
			workdayAvg += avg
			workdayCount++
		}
	}

	if workdayCount > 0 {
		workdayAvg /= float64(workdayCount)
	}
	if weekendCount > 0 {
		weekendAvg /= float64(weekendCount)
	}

	// 如果工作日和周末的情感差异显著，则认为存在周度模式
	return math.Abs(workdayAvg-weekendAvg) > 0.3
}

// detectTransitionPattern 检测情感转换模式
func (ea *EmotionAnalyzer) detectTransitionPattern(trendData []EmotionTrendData) *EmotionPattern {
	if len(trendData) < 5 {
		return nil
	}

	transitions := make(map[string]int)
	for i := 1; i < len(trendData); i++ {
		from := trendData[i-1].Dominant
		to := trendData[i].Dominant
		if from != to {
			transition := fmt.Sprintf("%s->%s", from, to)
			transitions[transition]++
		}
	}

	// 找出最常见的转换
	maxCount := 0
	mostCommonTransition := ""
	for transition, count := range transitions {
		if count > maxCount {
			maxCount = count
			mostCommonTransition = transition
		}
	}

	if maxCount >= 3 { // 至少出现3次才认为有模式
		return &EmotionPattern{
			Type:        "transition",
			Description: fmt.Sprintf("情感转换模式：%s", mostCommonTransition),
			Frequency:   float64(maxCount) / float64(len(trendData)),
			Confidence:  0.5,
			Triggers:    []string{"外部事件", "话题变化", "社交互动"},
			Implications: []string{"注意情感转换的触发因素", "主动调节情感状态"},
		}
	}

	return nil
}

// MessageData 消息数据结构
type MessageData struct {
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	Sender    string    `json:"sender"`
}


// initializeEmotionKeywords 初始化情感关键词
func initializeEmotionKeywords() map[EmotionType][]string {
	return map[EmotionType][]string{
		EmotionJoy: {
			"开心", "高兴", "快乐", "愉快", "满意", "棒", "好", "太好了", "赞", "厉害", "成功",
			"完成", "解决了", "搞定了", "完美", "优秀", "美好", "精彩", "喜欢", "爱",
			"哈哈", "嘿嘿", "嘻嘻", "😊", "😄", "🎉", "👍", "💪",
		},
		EmotionLove: {
			"爱", "喜欢", "宝贝", "亲爱的", "想你", "关心", "在乎", "温暖", "感动", "感激",
			"谢谢", "感谢", "珍惜", "重要", "离不开", "甜蜜", "浪漫", "温柔", "体贴",
			"❤️", "💕", "💖", "💗", "💝",
		},
		EmotionSurprise: {
			"哇", "天啊", "天哪", "我天", "我的天", "没想到", "意外", "突然", "居然",
			"竟然", "真的吗", "不会吧", "太意外了", "吓我一跳", "震惊", "惊奇", "惊讶",
			"😮", "😲", "🤯", "😱",
		},
		EmotionAnger: {
			"生气", "愤怒", "气死", "烦死", "讨厌", "可恶", "该死", "混蛋", "愚蠢", "笨蛋",
			"搞什么", "什么鬼", "神经病", "有病", "无聊", "差劲", "糟糕", "失败", "错误",
			"问题", "麻烦", "头疼", "😡", "😤", "😠", "🤬",
		},
		EmotionSadness: {
			"难过", "伤心", "悲伤", "沮丧", "失望", "遗憾", "可惜", "心疼", "难过", "痛苦",
			"难受", "不舒服", "想哭", "哭", "叹气", "唉", "太惨了", "不幸", "糟糕",
			"😢", "😭", "😔", "💔",
		},
		EmotionFear: {
			"害怕", "恐惧", "担心", "焦虑", "紧张", "不安", "害怕", "恐怖", "可怕", "危险",
			"担心", "忧虑", "害怕", "恐惧", "紧张", "不确定", "没把握", "没信心", "怀疑",
			"😨", "😰", "😱", "🙀",
		},
		EmotionDisgust: {
			"恶心", "厌恶", "反感", "讨厌", "恶心", "令人作呕", "受不了", "难以忍受",
			"脏", "恶心", "反胃", "作呕", "厌恶", "鄙视", "鄙视", "看不起", "轻视",
			"🤢", "🤮", "🤧",
		},
		EmotionNeutral: {
			"正常", "一般", "还行", "可以", "好的", "嗯", "哦", "嗯嗯", "好吧", "知道了",
			"明白", "了解", "清楚", "没问题", "行", "可以", "好的", "收到", "谢谢",
		},
	}
}

// initializeModifiers 初始化修饰词
func initializeModifiers() map[string]float64 {
	return map[string]float64{
		// 加强词
		"很": 1.5,
		"非常": 2.0,
		"特别": 1.8,
		"极其": 2.5,
		"超级": 2.2,
		"太": 2.0,
		"真的": 1.7,
		"确实": 1.6,

		// 减弱词
		"有点": 0.7,
		"稍微": 0.6,
		"略微": 0.5,
		"比较": 0.8,
		"还算": 0.7,
		"一般": 0.5,
		"不太": 0.4,
		"不怎么": 0.3,

		// 否定词
		"不": -1.0,
		"没": -1.0,
		"无": -1.0,
		"不是": -1.0,
		"没有": -1.0,

		// 程度词
		"最": 2.5,
		"第一": 2.2,
		"首位": 2.0,
		"无比": 2.3,
	}
}