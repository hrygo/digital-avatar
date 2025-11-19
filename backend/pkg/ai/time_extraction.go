package ai

import (
	"strings"
)

// TimePattern 时间模式
type TimePattern struct {
	Pattern   *RegexpWrapper `json:"-"`
	Type      string         `json:"type"`
	Format    string         `json:"format"`
}

// NewTimePattern 创建时间模式
func NewTimePattern(pattern, timeType, format string) *TimePattern {
	return &TimePattern{
		Pattern: NewRegexpWrapper(pattern),
		Type:    timeType,
		Format:  format,
	}
}

// TimeExtractor 时间表达提取器
type TimeExtractor struct {
	patterns []*TimePattern
}

// NewTimeExtractor 创建时间提取器
func NewTimeExtractor() *TimeExtractor {
	return &TimeExtractor{
		patterns: initializeTimePatterns(),
	}
}

// ExtractTimeExpression 提取时间表达
func (te *TimeExtractor) ExtractTimeExpression(text string) string {
	text = strings.ToLower(text)

	for _, pattern := range te.patterns {
		if pattern.Pattern == nil {
			continue
		}

		compiledPattern := pattern.Pattern.CompileRegexp()
		if compiledPattern == nil {
			continue
		}

		if compiledPattern.MatchString(text) {
			match := compiledPattern.FindString(text)
			if match != "" {
				return pattern.Type + ":" + match
			}
		}
	}

	return ""
}

// initializeTimePatterns 初始化时间模式
func initializeTimePatterns() []*TimePattern {
	return []*TimePattern{
		NewTimePattern(
			`(\d{4}-\d{1,2}-\d{1,2})`,
			"absolute_date",
			"YYYY-MM-DD",
		),
		NewTimePattern(
			`(\d{1,2}:\d{2})`,
			"absolute_time",
			"HH:MM",
		),
		NewTimePattern(
			`(今天|明天|昨天|后天)`,
			"relative_date",
			"relative",
		),
		NewTimePattern(
			`(下周|上周|本周|这周)`,
			"relative_week",
			"relative",
		),
		NewTimePattern(
			`(早上|上午|下午|晚上|中午|半夜|凌晨)`,
			"time_of_day",
			"relative",
		),
	}
}