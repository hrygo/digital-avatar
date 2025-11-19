package ai

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"
)

// SimpleTopicModel 简化的主题模型
type SimpleTopicModel struct {
	config        *TopicModelConfig
	vocabulary    map[string]int
	invVocabulary []string
	wordsPerTopic map[string][]string // 主题 -> 词汇列表
	topicWeights  map[string]float64  // 主题 -> 权重
	isTrained     bool
	rng           *rand.Rand
}

// NewSimpleTopicModel 创建简化主题模型
func NewSimpleTopicModel(config *TopicModelConfig) *SimpleTopicModel {
	if config == nil {
		config = DefaultTopicModelConfig()
	}

	return &SimpleTopicModel{
		config:        config,
		wordsPerTopic: make(map[string][]string),
		topicWeights:  make(map[string]float64),
		rng:           rand.New(rand.NewSource(config.Seed)),
	}
}

// Fit 训练简化的主题模型
func (stm *SimpleTopicModel) Fit(documents []string) error {
	if len(documents) == 0 {
		return fmt.Errorf("no documents provided")
	}

	// 1. 文本预处理
	preprocessor := NewTextPreprocessor()
	tokenizedDocs := preprocessor.PreprocessDocuments(documents)

	if len(tokenizedDocs) == 0 {
		return fmt.Errorf("no valid documents after preprocessing")
	}

	// 2. 构建词汇表
	stm.vocabulary = preprocessor.BuildVocabulary(tokenizedDocs)

	if len(stm.vocabulary) == 0 {
		return fmt.Errorf("empty vocabulary after preprocessing")
	}

	// 3. 构建逆向词汇表
	var words []string
	for word := range stm.vocabulary {
		words = append(words, word)
	}
	sort.Strings(words)

	stm.invVocabulary = words

	// 4. 简化的主题提取
	err := stm.extractTopics(tokenizedDocs)
	if err != nil {
		return fmt.Errorf("topic extraction failed: %v", err)
	}

	stm.isTrained = true
	return nil
}

// extractTopics 提取主题
func (stm *SimpleTopicModel) extractTopics(documents [][]string) error {
	numTopics := stm.config.NumTopics

	// 统计所有词频
	wordFreq := make(map[string]int)
	totalWords := 0

	for _, doc := range documents {
		for _, word := range doc {
			if _, exists := stm.vocabulary[word]; exists {
				wordFreq[word]++
				totalWords++
			}
		}
	}

	// 将词汇按频率排序
	type wordFreqPair struct {
		word  string
		freq  int
		score float64
	}

	var wordPairs []wordFreqPair
	for word, freq := range wordFreq {
		// 计算TF-IDF风格的分数
		score := float64(freq) / float64(totalWords)
		if score > 0.001 { // 过滤低频词
			wordPairs = append(wordPairs, wordFreqPair{word, freq, score})
		}
	}

	sort.Slice(wordPairs, func(i, j int) bool {
		return wordPairs[i].score > wordPairs[j].score
	})

	// 分配词汇到主题
	wordsPerTopic := numTopics * 10 // 每个主题10个词
	if len(wordPairs) < wordsPerTopic {
		wordsPerTopic = len(wordPairs) / numTopics
		if wordsPerTopic < 5 {
			wordsPerTopic = 5
		}
	}

	for k := 0; k < numTopics && k*wordsPerTopic < len(wordPairs); k++ {
		topicID := fmt.Sprintf("topic_%d", k)
		var topicWords []string

		for i := k * wordsPerTopic; i < (k+1)*wordsPerTopic && i < len(wordPairs); i++ {
			topicWords = append(topicWords, wordPairs[i].word)
		}

		// 计算主题权重
		topicWeight := 0.0
		for _, word := range topicWords {
			topicWeight += float64(wordFreq[word])
		}
		topicWeight /= float64(totalWords)

		stm.wordsPerTopic[topicID] = topicWords
		stm.topicWeights[topicID] = topicWeight
	}

	return nil
}

// GetTopics 获取主题
func (stm *SimpleTopicModel) GetTopics() []Topic {
	if !stm.isTrained {
		return []Topic{}
	}

	var topics []Topic

	for topicID, words := range stm.wordsPerTopic {
		weight := stm.topicWeights[topicID]
		coherence := stm.calculateTopicCoherence(words)

		topic := Topic{
			ID:          topicID,
			Name:        stm.generateTopicName(words),
			Description: stm.generateTopicDescription(words),
			Keywords:    words,
			Weight:      weight,
			Coherence:   coherence,
			Timestamp:   time.Now(),
			Metadata:    map[string]string{"algorithm": "SimpleTopicModel"},
		}

		topics = append(topics, topic)
	}

	// 按权重排序
	sort.Slice(topics, func(i, j int) bool {
		return topics[i].Weight > topics[j].Weight
	})

	return topics
}

// GetDocumentTopics 获取文档主题分布
func (stm *SimpleTopicModel) GetDocumentTopics(document string) *DocumentTopicDistribution {
	if !stm.isTrained {
		return nil
	}

	// 预处理文档
	preprocessor := NewTextPreprocessor()
	tokens := preprocessor.PreprocessText(document)

	if len(tokens) == 0 {
		return nil
	}

	// 计算与各主题的匹配度
	var distributions []TopicDistribution
	totalScore := 0.0

	for topicID, topicWords := range stm.wordsPerTopic {
		matchScore := stm.calculateDocumentTopicMatch(tokens, topicWords)
		totalScore += matchScore

		distributions = append(distributions, TopicDistribution{
			TopicID:    topicID,
			Probability: matchScore,
			Relevance:  matchScore,
		})
	}

	// 归一化概率
	if totalScore > 0 {
		for i := range distributions {
			distributions[i].Probability /= totalScore
			distributions[i].Relevance = distributions[i].Probability
		}
	}

	// 找到主导主题
	dominantTopic := ""
	maxProb := 0.0
	for _, dist := range distributions {
		if dist.Probability > maxProb {
			maxProb = dist.Probability
			dominantTopic = dist.TopicID
		}
	}

	// 计算确定性
	certainty := maxProb

	return &DocumentTopicDistribution{
		DocumentID:    fmt.Sprintf("doc_%d", time.Now().Unix()),
		Distributions: distributions,
		DominantTopic: dominantTopic,
		Certainty:     certainty,
		Timestamp:     time.Now(),
	}
}

// calculateDocumentTopicMatch 计算文档与主题的匹配度
func (stm *SimpleTopicModel) calculateDocumentTopicMatch(docTokens []string, topicWords []string) float64 {
	if len(topicWords) == 0 {
		return 0.0
	}

	matchCount := 0
	docWordSet := make(map[string]bool)

	for _, word := range docTokens {
		docWordSet[word] = true
	}

	for _, topicWord := range topicWords {
		if docWordSet[topicWord] {
			matchCount++
		}
	}

	return float64(matchCount) / math.Sqrt(float64(len(docTokens)*len(topicWords)))
}

// calculateTopicCoherence 计算主题连贯性
func (stm *SimpleTopicModel) calculateTopicCoherence(words []string) float64 {
	if len(words) < 2 {
		return 0.0
	}

	coherence := 0.0
	for i := 0; i < len(words)-1; i++ {
		for j := i + 1; j < len(words); j++ {
			coherence += 1.0 / float64(j-i) // 距离权重
		}
	}

	return coherence / float64(len(words)*(len(words)-1)/2)
}

// generateTopicName 生成主题名称
func (stm *SimpleTopicModel) generateTopicName(words []string) string {
	if len(words) == 0 {
		return "未知主题"
	}

	if len(words) >= 2 {
		return fmt.Sprintf("%s & %s", words[0], words[1])
	}

	return words[0]
}

// generateTopicDescription 生成主题描述
func (stm *SimpleTopicModel) generateTopicDescription(words []string) string {
	if len(words) == 0 {
		return "该主题没有明确的词汇特征"
	}

	limit := len(words)
	if limit > 5 {
		limit = 5
	}

	return fmt.Sprintf("该主题主要涉及：%s", strings.Join(words[:limit], "、"))
}

// Validate 验证模型
func (stm *SimpleTopicModel) Validate() *ValidationStatus {
	status := &ValidationStatus{
		IsValid:      true,
		Perplexity:   100.0, // 简化模型使用固定值
		Coherence:    0.5,
		Stability:    1.0,
		Warnings:     []string{},
		Recommendations: []string{},
	}

	if !stm.isTrained {
		status.IsValid = false
		status.Warnings = append(status.Warnings, "模型未训练")
		return status
	}

// 计算平均连贯性
	if len(stm.wordsPerTopic) > 0 {
		totalCoherence := 0.0
		for _, words := range stm.wordsPerTopic {
			totalCoherence += stm.calculateTopicCoherence(words)
		}
		status.Coherence = totalCoherence / float64(len(stm.wordsPerTopic))
	}

	// 检查问题
	if status.Coherence < 0.3 {
		status.Warnings = append(status.Warnings, "主题连贯性较低")
	}

	if len(stm.vocabulary) < 50 {
		status.Warnings = append(status.Warnings, "词汇表过小")
	}

	if len(stm.wordsPerTopic) < 3 {
		status.Warnings = append(status.Warnings, "提取的主题数量过少")
	}

	// 生成建议
	if len(status.Warnings) > 0 {
		status.Recommendations = append(status.Recommendations, "建议增加更多文档或调整预处理参数")
	}

	return status
}

// OptimizeTopicNumber 优化主题数量
func (stm *SimpleTopicModel) OptimizeTopicNumber(documents []string, maxTopics int) (int, []float64, error) {
	if maxTopics < 2 {
		maxTopics = 20
	}

	var scores []float64
	bestScore := -math.Inf(1)
	bestTopics := 5 // 默认值

	// 测试不同主题数量
	for k := 2; k <= maxTopics && k <= len(documents)/2; k++ {
		config := *stm.config
		config.NumTopics = k

		tempModel := NewSimpleTopicModel(&config)
		err := tempModel.Fit(documents)
		if err != nil {
			continue // 跳过失败的配置
		}

		validation := tempModel.Validate()

		// 简化的评分：连贯性 + 主题数量权重
		score := validation.Coherence + float64(k)/100.0

		scores = append(scores, score)

		if score > bestScore {
			bestScore = score
			bestTopics = k
		}
	}

	return bestTopics, scores, nil
}