package ai

import (
	"regexp"
	"time"
)

// IntentType 意图类型
type IntentType string

const (
	IntentTask       IntentType = "task"       // 任务意图
	IntentQuestion   IntentType = "question"   // 问题意图
	IntentDecision   IntentType = "decision"   // 决策意图
	IntentInfo       IntentType = "info"       // 信息意图
	IntentSocial     IntentType = "social"     // 社交意图
	IntentEmotional  IntentType = "emotional"  // 情感表达
	IntentRequest    IntentType = "request"    // 请求意图
	IntentOffer      IntentType = "offer"      // 提供意图
	IntentComplaint  IntentType = "complaint"  // 抱怨意图
	IntentGratitude  IntentType = "gratitude"  // 感谢意图
	IntentGreeting   IntentType = "greeting"   // 问候意图
	IntentFarewell   IntentType = "farewell"   // 告别意图
	IntentNeutral    IntentType = "neutral"    // 中性意图
)

// IntentAnalysis 意图分析结果
type IntentAnalysis struct {
	Type          IntentType   `json:"type"`          // 主要意图类型
	SubType       string       `json:"sub_type"`       // 子类型
	Confidence     float64      `json:"confidence"`     // 置信度 0-1
	Entities      []Entity     `json:"entities"`      // 识别的实体
	ActionVerbs    []string     `json:"action_verbs"`  // 行动词
	TimeExpression string       `json:"time_expression"` // 时间表达
	Priority      string       `json:"priority"`      // 优先级
	Context       string       `json:"context"`       // 上下文
	Implications  []string     `json:"implications"`  // 含义和影响
	Timestamp     time.Time    `json:"timestamp"`     // 时间戳
}

// Entity 实体识别结果
type Entity struct {
	Type     string  `json:"type"`     // 实体类型：person, location, organization, time, amount
	Value    string  `json:"value"`    // 实体值
	Position int     `json:"position"` // 在文本中的位置
	Confidence float64 `json:"confidence"` // 置信度
}

// IntentPattern 意图模式
type IntentPattern struct {
	Pattern   *RegexpWrapper `json:"-"`
	Weight    float64        `json:"weight"`
	Keywords  []string       `json:"keywords"`
	Entities  []string       `json:"entities"`
	Context   string         `json:"context"`
}

// RegexpWrapper 包装正则表达式以支持JSON序列化
type RegexpWrapper struct {
	Pattern string `json:"pattern"`
}

// CompileRegexp 编译正则表达式
func (rw *RegexpWrapper) CompileRegexp() *regexp.Regexp {
	if rw.Pattern == "" {
		return nil
	}
	return regexp.MustCompile(rw.Pattern)
}

// NewRegexpWrapper 创建正则表达式包装器
func NewRegexpWrapper(pattern string) *RegexpWrapper {
	return &RegexpWrapper{Pattern: pattern}
}