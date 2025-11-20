package ai

import (
	"fmt"
	"math"
	"sort"
	"time"
)

// TopicAnalyzer 主题分析器
type TopicAnalyzer struct {
	model        interface{} // *LDAModel or *SimpleTopicModel
	preprocessor *TextPreprocessor
	config       *TopicModelConfig
	useSimple    bool // 是否使用简化模型
}

// NewTopicAnalyzer 创建主题分析器
func NewTopicAnalyzer(config *TopicModelConfig) *TopicAnalyzer {
	if config == nil {
		config = DefaultTopicModelConfig()
	}

	return &TopicAnalyzer{
		model:        NewSimpleTopicModel(config),
		preprocessor: NewTextPreprocessor(),
		config:       config,
		useSimple:    true,
	}
}

// NewAdvancedTopicAnalyzer 创建高级主题分析器（暂用简化模型）
func NewAdvancedTopicAnalyzer(config *TopicModelConfig) *TopicAnalyzer {
	if config == nil {
		config = DefaultTopicModelConfig()
	}

	// 目前使用简化模型，后续可扩展为LDA
	return &TopicAnalyzer{
		model:        NewSimpleTopicModel(config),
		preprocessor: NewTextPreprocessor(),
		config:       config,
		useSimple:    true,
	}
}

// AnalyzeTopics 分析文档主题
func (ta *TopicAnalyzer) AnalyzeTopics(documents []string) ([]Topic, error) {
	if len(documents) == 0 {
		return []Topic{}, fmt.Errorf("no documents provided")
	}

	var topics []Topic
	var err error

	if ta.useSimple {
		// 使用简化模型
		simpleModel := ta.model.(*SimpleTopicModel)
		err = simpleModel.Fit(documents)
		if err != nil {
			return []Topic{}, fmt.Errorf("failed to train simple topic model: %v", err)
		}
		topics = simpleModel.GetTopics()
	} else {
		// 使用简化模型（高级模式）
		simpleModel := ta.model.(*SimpleTopicModel)
		err = simpleModel.Fit(documents)
		if err != nil {
			return []Topic{}, fmt.Errorf("failed to train topic model: %v", err)
		}
		topics = simpleModel.GetTopics()
	}

	// 计算主题质量指标
	for i := range topics {
		topics[i].Coherence = ta.calculateTopicCoherence(&topics[i])
	}

	// 按权重排序
	sort.Slice(topics, func(i, j int) bool {
		return topics[i].Weight > topics[j].Weight
	})

	return topics, nil
}

// AnalyzeDocumentTopics 分析单个文档的主题分布
func (ta *TopicAnalyzer) AnalyzeDocumentTopics(document string) (*DocumentTopicDistribution, error) {
	if ta.useSimple {
		simpleModel := ta.model.(*SimpleTopicModel)
		if !simpleModel.isTrained {
			return nil, fmt.Errorf("model not trained yet")
		}
		return simpleModel.GetDocumentTopics(document), nil
	} else {
		simpleModel := ta.model.(*SimpleTopicModel)
		if !simpleModel.isTrained {
			return nil, fmt.Errorf("model not trained yet")
		}
		return simpleModel.GetDocumentTopics(document), nil
	}
}

// TrackTopicEvolution 跟踪主题演化
func (ta *TopicAnalyzer) TrackTopicEvolution(documents []DocumentWithTime) ([]TopicEvolution, error) {
	if len(documents) < 2 {
		return []TopicEvolution{}, fmt.Errorf("need at least 2 documents for evolution analysis")
	}

	// 按时间排序文档
	sort.Slice(documents, func(i, j int) bool {
		return documents[i].Timestamp.Before(documents[j].Timestamp)
	})

	// 分时间窗口分析
	windowSize := len(documents) / 4 // 分成4个时间窗口
	if windowSize < 1 {
		windowSize = 1
	}

	evoData := make(map[string][]TopicTimePoint)

	for i := 0; i < len(documents); i += windowSize {
		end := i + windowSize
		if end > len(documents) {
			end = len(documents)
		}

		// 提取时间窗口内的文档
		windowDocs := make([]string, 0, end-i)
		for j := i; j < end; j++ {
			windowDocs = append(windowDocs, documents[j].Content)
		}

		// 分析该时间窗口的主题
		windowTopics, err := ta.AnalyzeTopics(windowDocs)
		if err != nil {
			continue // 跳过错误窗口
		}

		// 记录时间点数据
		timePoint := documents[i].Timestamp
		for _, topic := range windowTopics {
			timePointData := TopicTimePoint{
				Timestamp: timePoint,
				Weight:    topic.Weight,
				Volume:    len(topic.Documents),
				Sentiment: ta.calculateTopicSentiment(topic),
			}

			evoData[topic.ID] = append(evoData[topic.ID], timePointData)
		}

		// 这里可以添加趋势检测逻辑
		// future := make(map[string]float64)
	}

	// 构建主题演化结果
	var evolutions []TopicEvolution
	for topicID, timePoints := range evoData {
		if len(timePoints) < 2 {
			continue
		}

		evolution := TopicEvolution{
			TopicID:    topicID,
			TimePoints: timePoints,
			Trend:      ta.detectTrend(timePoints),
			ChangeRate: ta.calculateChangeRate(timePoints),
			Significance: ta.calculateSignificance(timePoints),
			Period:     "time_window",
			Metadata:   map[string]interface{}{"window_size": windowSize},
		}

		evolutions = append(evolutions, evolution)
	}

	return evolutions, nil
}

// BuildTopicNetwork 构建主题关系网络
func (ta *TopicAnalyzer) BuildTopicNetwork(topics []Topic, documents []string) (*TopicNetwork, error) {
	if len(topics) == 0 {
		return &TopicNetwork{}, fmt.Errorf("no topics provided")
	}

	// 计算主题间相关性
	relations := ta.calculateTopicRelations(topics, documents)

	// 计算网络指标
	density := ta.calculateNetworkDensity(relations)
	modularity := ta.calculateModularity(relations, topics)

	network := &TopicNetwork{
		Topics:     topics,
		Relations:  relations,
		Density:    density,
		Modularity: modularity,
		Timestamp:  time.Now(),
	}

	return network, nil
}

// ValidateTopics 验证主题质量
func (ta *TopicAnalyzer) ValidateTopics() *ValidationStatus {
	if ta.useSimple {
		simpleModel := ta.model.(*SimpleTopicModel)
		if !simpleModel.isTrained {
			return &ValidationStatus{
				IsValid:    false,
				Warnings:   []string{"Model not trained"},
				Perplexity: 0.0,
				Coherence:  0.0,
			}
		}
		return simpleModel.Validate()
	} else {
		simpleModel := ta.model.(*SimpleTopicModel)
		if !simpleModel.isTrained {
			return &ValidationStatus{
				IsValid:    false,
				Warnings:   []string{"Model not trained"},
				Perplexity: 0.0,
				Coherence:  0.0,
			}
		}
		return simpleModel.Validate()
	}
}

// OptimizeTopicNumber 优化主题数量
func (ta *TopicAnalyzer) OptimizeTopicNumber(documents []string, maxTopics int) (int, []float64, error) {
	if maxTopics < 2 {
		maxTopics = 20
	}

	if ta.useSimple {
		simpleModel := ta.model.(*SimpleTopicModel)
		return simpleModel.OptimizeTopicNumber(documents, maxTopics)
	} else {
		var perplexities []float64
		var coherences []float64
		bestScore := -math.Inf(1)
		bestTopics := 5 // 默认值

		// 测试不同主题数量
		for k := 2; k <= maxTopics && k <= len(documents)/2; k++ {
			config := *ta.config
			config.NumTopics = k

			tempModel := NewSimpleTopicModel(&config)
			err := tempModel.Fit(documents)
			if err != nil {
				continue // 跳过失败的配置
			}

			topics := tempModel.GetTopics()
			validation := tempModel.Validate()

			perplexities = append(perplexities, validation.Perplexity)

			avgCoherence := 0.0
			for _, topic := range topics {
				avgCoherence += topic.Coherence
			}
			if len(topics) > 0 {
				avgCoherence /= float64(len(topics))
			}
			coherences = append(coherences, avgCoherence)

			// 综合评分（连贯性 + 主题数量权重）
			score := avgCoherence + float64(k)/100.0
			if score > bestScore {
				bestScore = score
				bestTopics = k
			}
		}

		// 合并评分指标
		scores := make([]float64, len(perplexities))
		for i := range scores {
			scores[i] = coherences[i] + float64(i+2)/100.0
		}

		return bestTopics, scores, nil
	}
}

// calculateTopicCoherence 计算主题连贯性
func (ta *TopicAnalyzer) calculateTopicCoherence(topic *Topic) float64 {
	if len(topic.Keywords) < 2 {
		return 0.0
	}

	// 简化的连贯性计算
	coherence := 0.0
	for i := 0; i < len(topic.Keywords)-1; i++ {
		for j := i + 1; j < len(topic.Keywords); j++ {
			// 这里可以使用更复杂的共现统计
			coherence += 1.0 / float64(j-i) // 距离权重
		}
	}

	return coherence / float64(len(topic.Keywords)*(len(topic.Keywords)-1)/2)
}

// calculateTopicSentiment 计算主题情感倾向
func (ta *TopicAnalyzer) calculateTopicSentiment(topic Topic) float64 {
	// 简化的情感计算
	positiveWords := []string{"好", "棒", "优秀", "成功", "喜欢", "开心", "满意", "赞"}
	negativeWords := []string{"坏", "差", "失败", "讨厌", "糟糕", "问题", "错误", "抱怨"}

	positiveCount := 0
	negativeCount := 0

	for _, keyword := range topic.Keywords {
		for _, posWord := range positiveWords {
			if keyword == posWord {
				positiveCount++
				break
			}
		}

		for _, negWord := range negativeWords {
			if keyword == negWord {
				negativeCount++
				break
			}
		}
	}

	total := positiveCount + negativeCount
	if total == 0 {
		return 0.0 // 中性
	}

	return float64(positiveCount-negativeCount) / float64(total)
}

// calculateTopicRelations 计算主题间关系
func (ta *TopicAnalyzer) calculateTopicRelations(topics []Topic, documents []string) []TopicRelation {
	var relations []TopicRelation

	// 计算主题词汇重叠度
	for i := 0; i < len(topics); i++ {
		for j := i + 1; j < len(topics); j++ {
			overlap := ta.calculateWordOverlap(topics[i].Keywords, topics[j].Keywords)

			if overlap > 0.1 { // 只保留有意义的关系
				relation := TopicRelation{
					SourceTopicID: topics[i].ID,
					TargetTopicID: topics[j].ID,
					RelationType:  "correlation",
					Strength:      overlap,
					Confidence:    math.Min(overlap*2, 1.0),
					Description:   fmt.Sprintf("主题 %s 和 %s 在词汇上有 %f 的重叠", topics[i].Name, topics[j].Name, overlap),
				}
				relations = append(relations, relation)
			}
		}
	}

	return relations
}

// calculateWordOverlap 计算词汇重叠度
func (ta *TopicAnalyzer) calculateWordOverlap(words1, words2 []string) float64 {
	set1 := make(map[string]bool)
	set2 := make(map[string]bool)

	for _, word := range words1 {
		set1[word] = true
	}

	for _, word := range words2 {
		set2[word] = true
	}

	intersection := 0
	for word := range set1 {
		if set2[word] {
			intersection++
		}
	}

	union := len(set1) + len(set2) - intersection
	if union == 0 {
		return 0.0
	}

	return float64(intersection) / float64(union)
}

// calculateNetworkDensity 计算网络密度
func (ta *TopicAnalyzer) calculateNetworkDensity(relations []TopicRelation) float64 {
	if len(relations) == 0 {
		return 0.0
	}

	// 简化计算：关系数量 / 可能的最大关系数
	return float64(len(relations)) / 100.0 // 假设最大100个关系
}

// calculateModularity 计算网络模块度
func (ta *TopicAnalyzer) calculateModularity(relations []TopicRelation, topics []Topic) float64 {
	if len(topics) < 2 {
		return 0.0
	}

	// 简化的模块度计算
	return 0.3 + float64(len(topics))/100.0 // 简化的启发式值
}

// detectTrend 检测趋势
func (ta *TopicAnalyzer) detectTrend(timePoints []TopicTimePoint) string {
	if len(timePoints) < 2 {
		return "stable"
	}

	first := timePoints[0].Weight
	last := timePoints[len(timePoints)-1].Weight

	changeRate := (last - first) / first

	if changeRate > 0.1 {
		return "increasing"
	} else if changeRate < -0.1 {
		return "decreasing"
	} else {
		return "stable"
	}
}

// calculateChangeRate 计算变化率
func (ta *TopicAnalyzer) calculateChangeRate(timePoints []TopicTimePoint) float64 {
	if len(timePoints) < 2 {
		return 0.0
	}

	totalChange := 0.0
	for i := 1; i < len(timePoints); i++ {
		change := timePoints[i].Weight - timePoints[i-1].Weight
		totalChange += math.Abs(change)
	}

	return totalChange / float64(len(timePoints)-1)
}

// calculateSignificance 计算显著性
func (ta *TopicAnalyzer) calculateSignificance(timePoints []TopicTimePoint) float64 {
	if len(timePoints) == 0 {
		return 0.0
	}

	avgWeight := 0.0
	for _, tp := range timePoints {
		avgWeight += tp.Weight
	}
	avgWeight /= float64(len(timePoints))

	// 简化的显著性：平均权重的归一化值
	return math.Min(avgWeight*10, 1.0)
}

// DocumentWithTime 带时间戳的文档
type DocumentWithTime struct {
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	ID        string    `json:"id"`
}