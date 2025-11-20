package ai

import "time"

// Topic 主题
type Topic struct {
	ID           string            `json:"id"`           // 主题ID
	Name         string            `json:"name"`         // 主题名称
	Description  string            `json:"description"`  // 主题描述
	Keywords     []string          `json:"keywords"`     // 关键词列表
	Weight       float64           `json:"weight"`       // 主题权重
	Coherence    float64           `json:"coherence"`    // 主题连贯性
	Documents    []string          `json:"documents"`    // 相关文档ID
	Timestamp    time.Time         `json:"timestamp"`    // 创建时间
	Metadata     map[string]string `json:"metadata"`     // 元数据
}

// TopicWord 主题词汇
type TopicWord struct {
	Word     string  `json:"word"`     // 词汇
	Probability float64 `json:"probability"` // 概率
	Weight   float64 `json:"weight"`   // 权重
}

// TopicDistribution 主题分布
type TopicDistribution struct {
	TopicID    string  `json:"topic_id"`    // 主题ID
	Probability float64 `json:"probability"` // 概率
	Relevance  float64 `json:"relevance"`  // 相关性
}

// DocumentTopicDistribution 文档主题分布
type DocumentTopicDistribution struct {
	DocumentID   string             `json:"document_id"`   // 文档ID
	Distributions []TopicDistribution `json:"distributions"` // 主题分布
	DominantTopic string            `json:"dominant_topic"` // 主导主题
	Certainty     float64           `json:"certainty"`     // 确定性
	Timestamp     time.Time         `json:"timestamp"`     // 时间戳
}

// TopicEvolution 主题演化
type TopicEvolution struct {
	TopicID      string                 `json:"topic_id"`      // 主题ID
	TimePoints   []TopicTimePoint       `json:"time_points"`   // 时间点数据
	Trend        string                 `json:"trend"`         // 趋势：increasing, decreasing, stable
	ChangeRate   float64               `json:"change_rate"`   // 变化率
	Significance float64               `json:"significance"`  // 显著性
	Period       string                `json:"period"`        // 周期：daily, weekly, monthly
	Metadata     map[string]interface{} `json:"metadata"`     // 元数据
}

// TopicTimePoint 主题时间点
type TopicTimePoint struct {
	Timestamp time.Time `json:"timestamp"` // 时间戳
	Weight    float64   `json:"weight"`    // 权重
	Volume    int       `json:"volume"`    // 讨论量
	Sentiment float64   `json:"sentiment"` // 情感倾向
}

// TopicRelation 主题关系
type TopicRelation struct {
	SourceTopicID string  `json:"source_topic_id"` // 源主题ID
	TargetTopicID string  `json:"target_topic_id"` // 目标主题ID
	RelationType  string  `json:"relation_type"`  // 关系类型：correlation, causal, hierarchical
	Strength      float64 `json:"strength"`       // 关系强度
	Confidence    float64 `json:"confidence"`     // 置信度
	Description   string  `json:"description"`    // 关系描述
}

// TopicNetwork 主题网络
type TopicNetwork struct {
	Topics     []Topic         `json:"topics"`     // 主题列表
	Relations  []TopicRelation `json:"relations"`  // 关系列表
	Density    float64         `json:"density"`    // 网络密度
	Modularity float64         `json:"modularity"` // 模块度
	Timestamp  time.Time       `json:"timestamp"`  // 时间戳
}

// TopicModelConfig 主题模型配置
type TopicModelConfig struct {
	NumTopics        int     `json:"num_topics"`         // 主题数量
	Alpha            float64 `json:"alpha"`             // 文档-主题分布参数
	Beta             float64 `json:"beta"`              // 主题-词汇分布参数
	Iterations       int     `json:"iterations"`        // 迭代次数
	Burnin           int     `json:"burnin"`            // 预烧期
	ThinInterval     int     `json:"thin_interval"`     // 采样间隔
	MinWordFreq      int     `json:"min_word_freq"`     // 最小词频
	MaxWordFreq      int     `json:"max_word_freq"`     // 最大词频
	VocabularySize   int     `json:"vocabulary_size"`   // 词汇表大小
	StopWords        []string `json:"stop_words"`       // 停用词
	Seed             int64   `json:"seed"`              // 随机种子
	OptimizeInterval int     `json:"optimize_interval"` // 优化间隔
}

// DefaultTopicModelConfig 默认主题模型配置
func DefaultTopicModelConfig() *TopicModelConfig {
	return &TopicModelConfig{
		NumTopics:        10,
		Alpha:            0.1,
		Beta:             0.01,
		Iterations:       1000,
		Burnin:           200,
		ThinInterval:     10,
		MinWordFreq:      2,
		MaxWordFreq:      100,
		VocabularySize:   1000,
		StopWords:        []string{"的", "了", "在", "是", "我", "你", "他", "她", "它", "们", "这", "那", "有", "和", "与", "或"},
		Seed:             42,
		OptimizeInterval: 50,
	}
}

// ValidationStatus 验证状态
type ValidationStatus struct {
	IsValid      bool     `json:"is_valid"`       // 是否有效
	Perplexity   float64  `json:"perplexity"`     // 困惑度
	Coherence    float64  `json:"coherence"`     // 连贯性
	Stability    float64  `json:"stability"`     // 稳定性
	Warnings     []string `json:"warnings"`       // 警告信息
	Recommendations []string `json:"recommendations"` // 建议
}