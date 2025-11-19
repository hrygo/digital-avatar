package ai

import (
	"fmt"
	"strings"
)

// EntityPattern 实体模式
type EntityPattern struct {
	Type    string         `json:"type"`
	Pattern *RegexpWrapper `json:"-"`
	Example string         `json:"example"`
}

// NewEntityPattern 创建实体模式
func NewEntityPattern(entityType, pattern, example string) *EntityPattern {
	return &EntityPattern{
		Type:    entityType,
		Pattern: NewRegexpWrapper(pattern),
		Example: example,
	}
}

// EntityRecognizer 实体识别器
type EntityRecognizer struct {
	patterns map[string][]*EntityPattern
}

// NewEntityRecognizer 创建实体识别器
func NewEntityRecognizer() *EntityRecognizer {
	return &EntityRecognizer{
		patterns: initializeEntityPatterns(),
	}
}

// ExtractEntities 识别文本中的实体
func (er *EntityRecognizer) ExtractEntities(text string) []Entity {
	var entities []Entity
	textLower := strings.ToLower(text)

	for entityType, patterns := range er.patterns {
		for _, pattern := range patterns {
			if pattern.Pattern == nil {
				continue
			}

			compiledPattern := pattern.Pattern.CompileRegexp()
			if compiledPattern == nil {
				continue
			}

			matches := compiledPattern.FindAllStringSubmatch(textLower, -1)
			for _, match := range matches {
				if len(match) > 1 {
					entity := Entity{
						Type:      entityType,
						Value:     match[1],
						Position: strings.Index(textLower, match[0]),
						Confidence: 0.8,
					}
					entities = append(entities, entity)
				}
			}
		}
	}

	// 去重
	uniqueEntities := make(map[string]Entity)
	for _, entity := range entities {
		key := fmt.Sprintf("%s:%s", entity.Type, entity.Value)
		uniqueEntities[key] = entity
	}

	var result []Entity
	for _, entity := range uniqueEntities {
		result = append(result, entity)
	}

	return result
}

// initializeEntityPatterns 初始化实体模式
func initializeEntityPatterns() map[string][]*EntityPattern {
	patterns := make(map[string][]*EntityPattern)

	// 人物实体
	patterns["person"] = []*EntityPattern{
		NewEntityPattern("person",
			`([a-zA-Z\x{4e00}-\x{9fa5}]{2,10}(?:先生|女士|老师|经理|总监|总|哥|姐|弟|妹)|张|王|李|刘|陈|杨|黄|赵|吴|周|徐|孙|马|朱|胡|郭|何|高|林|罗|郑|梁|谢|宋|唐|许|韩|冯|邓|曹|彭|曾|萧|田|董|袁|潘|于|蒋|蔡|余|杜|叶|程|魏|苏|吕|丁|任|沈|姚|卢|姜|崔|钟|谭|陆|汪|范|金|石|廖|贾|夏|韦|付|方|白|邹|孟|熊|秦|邱|江|尹|薛|闫|段|雷|侯|龙|史|陶|黎|贺|顾|毛|郝|龚|邵|万|钱|严|覃|武|戴|莫|孔|向|汤)`,
			"张三,李四"),
	}

	// 时间实体
	patterns["time"] = []*EntityPattern{
		NewEntityPattern("time",
			`(\d{4}-\d{1,2}-\d{1,2}|\d{1,2}月\d{1,2}日|\d{1,2}:\d{2}|今天|明天|昨天|后天|下周|上周|今年|明年|早上|晚上|中午|下午)`,
			"2024-01-01,明天"),
	}

	// 地点实体
	patterns["location"] = []*EntityPattern{
		NewEntityPattern("location",
			`([a-zA-Z\x{4e00}-\x{9fa5}]+(?:市|区|县|镇|村|街|路|号|楼|层|室|厅|馆|店|场|站|机场|车站|学校|医院|银行|公园|广场|市场|餐厅|咖啡厅|酒店|宾馆|办公楼|写字楼))`,
			"北京市,办公室"),
	}

	// 组织实体
	patterns["organization"] = []*EntityPattern{
		NewEntityPattern("organization",
			`([a-zA-Z\x{4e00}-\x{9fa5}]+(?:公司|企业|集团|机构|部门|团队|工作室|实验室|研究所|大学|学院|学校|医院|银行|商店|超市|餐厅|酒店|宾馆))`,
			"腾讯公司,阿里巴巴"),
	}

	return patterns
}