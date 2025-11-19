package crypto

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
)

// NLPPIIDetector 基于NLP技术的PII检测器
type NLPPIIDetector struct {
	replacementMap map[string]string
	mutex          sync.RWMutex

	// 姓名词典
	nameDictionary map[string]bool
	// 常见姓氏
	commonSurnames map[string]bool
	// 公司名后缀
	companySuffixes []string
	// 地理位置后缀
	locationSuffixes []string

	// 模式匹配器（作为后备）
	patternMatchers map[string]*regexp.Regexp
}

// PIIType PII类型枚举
type PIIType string

const (
	PIITypePerson      PIIType = "PER"      // 人名
	PIITypeOrganization PIIType = "ORG"      // 组织机构
	PIITypeLocation    PIIType = "LOC"      // 地理位置
	PIITypePhone       PIIType = "TEL"      // 电话号码
	PIITypeEmail       PIIType = "EMAIL"    // 邮箱地址
	PIITypeIDCard      PIIType = "ID"       // 身份证
	PIITypeAddress     PIIType = "ADDR"     // 地址
)

// PIIEntity PII实体
type PIIEntity struct {
	Text     string  `json:"text"`
	Type     PIIType `json:"type"`
	StartPos int     `json:"start_pos"`
	EndPos   int     `json:"end_pos"`
	Score    float64 `json:"score"`
	Replaced string  `json:"replaced,omitempty"`
}

// NewNLPPIIDetector 创建基于NLP的PII检测器
func NewNLPPIIDetector() *NLPPIIDetector {
	detector := &NLPPIIDetector{
		replacementMap:  make(map[string]string),
		nameDictionary:  make(map[string]bool),
		commonSurnames:  make(map[string]bool),
		patternMatchers: make(map[string]*regexp.Regexp),
	}

	detector.initDictionaries()
	detector.initPatterns()

	return detector
}

// initDictionaries 初始化词典
func (d *NLPPIIDetector) initDictionaries() {
	// 常见中文姓氏
	d.commonSurnames = map[string]bool{
		"李": true, "王": true, "张": true, "刘": true, "陈": true, "杨": true, "赵": true,
		"黄": true, "周": true, "吴": true, "徐": true, "孙": true, "马": true, "朱": true,
		"胡": true, "郭": true, "何": true, "高": true, "林": true, "罗": true, "郑": true,
		"梁": true, "谢": true, "宋": true, "唐": true, "许": true, "韩": true, "冯": true,
		"邓": true, "曹": true, "彭": true, "曾": true, "萧": true, "田": true, "董": true,
		"袁": true, "潘": true, "于": true, "蒋": true, "蔡": true, "余": true, "杜": true,
		"叶": true, "程": true, "魏": true, "吕": true, "丁": true, "任": true, "沈": true,
		"姚": true, "卢": true, "姜": true, "崔": true, "钟": true, "谭": true, "陆": true,
		"汪": true, "范": true, "金": true, "石": true, "廖": true, "贾": true, "夏": true,
		"韦": true, "付": true, "方": true, "白": true, "邹": true, "孟": true, "熊": true,
		"秦": true, "邱": true, "江": true, "尹": true, "薛": true, "闫": true, "段": true,
		"雷": true, "侯": true, "龙": true, "史": true, "陶": true, "黎": true, "贺": true,
		"顾": true, "毛": true, "郝": true, "龚": true, "邵": true, "万": true, "钱": true,
		"严": true, "覃": true, "武": true, "戴": true, "莫": true, "孔": true, "向": true,
	}

	// 常见名字
	commonNames := []string{
		"伟", "芳", "娜", "秀英", "敏", "静", "丽", "强", "磊", "军", "洋", "勇", "艳", "杰", "娟", "涛",
		"超", "明", "霞", "平", "刚", "桂英", "建华", "文", "华", "金", "春梅", "玉兰", "萍", "飞",
		"浩", "建国", "小红", "晓东", "小军", "建华", "小红", "张伟", "王芳", "李娜", "刘洋", "陈静",
		"杨敏", "赵强", "黄丽", "周磊", "吴军", "徐艳", "孙杰", "马娟", "朱涛", "胡超", "郭明",
	}

	for _, name := range commonNames {
		d.nameDictionary[name] = true
	}

	// 公司名后缀
	d.companySuffixes = []string{
		"有限公司", "股份有限公司", "集团", "科技公司", "网络公司", "信息公司",
		"文化公司", "教育公司", "医疗公司", "金融公司", "投资公司",
		"公司", "企业", "机构", "中心", "集团", "控股",
	}

	// 地理位置后缀
	d.locationSuffixes = []string{
		"省", "市", "区", "县", "镇", "乡", "街道", "路", "道", "巷", "弄", "号", "楼", "室",
		"村", "庄", "屯", "堡", "店", "铺", "馆", "所", "院", "校", "场", "站", "港",
		"自治区", "自治州", "自治县", "特别行政区",
	}
}

// initPatterns 初始化模式匹配器
func (d *NLPPIIDetector) initPatterns() {
	// 注意：检测顺序很重要，更具体的模式应该优先检测
	d.patternMatchers = map[string]*regexp.Regexp{
		// 身份证模式 - 最优先，避免被电话号码模式误匹配
		"idcard": regexp.MustCompile(`\d{17}[\dXx]|\d{15}`),
		// 手机号模式 - 限制为11位数字且以1开头
		"phone":  regexp.MustCompile(`(?:\+?86)?1[3-9]\d{9}|(?:\d{3,4}-\d{7,8})`),
		// 邮箱模式
		"email":  regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`),
	}
}

// DetectPII 检测文本中的PII信息
func (d *NLPPIIDetector) DetectPII(text string, threshold float64) ([]PIIEntity, error) {
	var entities []PIIEntity

	// 1. 使用模式匹配检测明确的PII类型
	entities = append(entities, d.detectByPatterns(text)...)

	// 2. 使用词典和规则检测姓名
	entities = append(entities, d.detectNames(text)...)

	// 3. 检测组织机构
	entities = append(entities, d.detectOrganizations(text)...)

	// 4. 检测地址（更具体，优先检测）
	entities = append(entities, d.detectAddresses(text)...)

	// 5. 检测地理位置
	entities = append(entities, d.detectLocations(text)...)

	// 6. 去重和排序
	entities = d.deduplicateAndSort(entities)

	// 7. 过滤低置信度结果
	var filtered []PIIEntity
	for _, entity := range entities {
		if entity.Score >= threshold {
			filtered = append(filtered, entity)
		}
	}

	return filtered, nil
}

// detectByPatterns 使用模式匹配检测PII
func (d *NLPPIIDetector) detectByPatterns(text string) []PIIEntity {
	var entities []PIIEntity

	// 检测顺序很重要：更具体和更重要的PII类型优先检测

	// 1. 检测身份证 - 最高优先级，避免被电话号码误匹配
	for _, match := range d.patternMatchers["idcard"].FindAllStringSubmatchIndex(text, -1) {
		start, end := match[0], match[1]
		original := text[start:end]

		score := 0.9
		if len(original) == 18 {
			score = 0.95 // 18位身份证更可信
		}

		// 验证身份证格式的基本合理性
		if d.isValidIDCard(original) {
			entities = append(entities, PIIEntity{
				Text:     original,
				Type:     PIITypeIDCard,
				StartPos: start,
				EndPos:   end,
				Score:    score,
			})
		}
	}

	// 2. 检测邮箱 - 高优先级，格式明确
	for _, match := range d.patternMatchers["email"].FindAllStringSubmatchIndex(text, -1) {
		start, end := match[0], match[1]
		original := text[start:end]

		entities = append(entities, PIIEntity{
			Text:     original,
			Type:     PIITypeEmail,
			StartPos: start,
			EndPos:   end,
			Score:    0.98,
		})
	}

	// 3. 检测电话号码 - 最后检测，优先级最低
	for _, match := range d.patternMatchers["phone"].FindAllStringSubmatchIndex(text, -1) {
		start, end := match[0], match[1]
		original := text[start:end]

		// 避免与已检测的身份证重叠
		overlap := false
		for _, entity := range entities {
			if entity.Type == PIITypeIDCard && d.isOverlapping(PIIEntity{StartPos: start, EndPos: end}, entity) {
				overlap = true
				break
			}
		}

		if !overlap {
			entities = append(entities, PIIEntity{
				Text:     original,
				Type:     PIITypePhone,
				StartPos: start,
				EndPos:   end,
				Score:    0.95,
			})
		}
	}

	return entities
}

// isValidIDCard 验证身份证号码格式的合理性
func (d *NLPPIIDetector) isValidIDCard(idCard string) bool {
	// 基本长度检查
	if len(idCard) != 15 && len(idCard) != 18 {
		return false
	}

	// 检查是否全为数字（除了最后一位可能是X）
	if len(idCard) == 15 {
		for i := 0; i < 15; i++ {
			if idCard[i] < '0' || idCard[i] > '9' {
				return false
			}
		}
	} else if len(idCard) == 18 {
		// 检查前17位
		for i := 0; i < 17; i++ {
			if idCard[i] < '0' || idCard[i] > '9' {
				return false
			}
		}
		// 检查最后一位
		lastChar := idCard[17]
		if !(lastChar >= '0' && lastChar <= '9' || lastChar == 'X' || lastChar == 'x') {
			return false
		}
	}

	return true
}

// detectNames 检测姓名
func (d *NLPPIIDetector) detectNames(text string) []PIIEntity {
	var entities []PIIEntity
	runes := []rune(text)

	for i := 0; i < len(runes); i++ {
		// 检查是否是常见姓氏
		if i < len(runes) && d.commonSurnames[string(runes[i])] {
			// 尝试匹配2-3个字符的姓名
			for nameLen := 2; nameLen <= 3 && i+nameLen <= len(runes); nameLen++ {
				name := string(runes[i : i+nameLen])

				// 计算置信度
				score := d.calculateNameScore(name, text, i)
				if score > 0.6 {
					byteStart := len(string(runes[:i]))
					byteEnd := len(string(runes[:i+nameLen]))

					entities = append(entities, PIIEntity{
						Text:     name,
						Type:     PIITypePerson,
						StartPos: byteStart,
						EndPos:   byteEnd,
						Score:    score,
					})
					break // 找到最合适的长度就停止
				}
			}
		}
	}

	return entities
}

// calculateNameScore 计算姓名置信度
func (d *NLPPIIDetector) calculateNameScore(name, context string, pos int) float64 {
	score := 0.0
	runes := []rune(name)

	// 基础分数：姓氏匹配
	if d.commonSurnames[string(runes[0])] {
		score += 0.4
	}

	// 长度检查
	if len(runes) == 2 || len(runes) == 3 {
		score += 0.2
	}

	// 词典检查
	if d.nameDictionary[name] {
		score += 0.3
	}

	// 上下文分析
	contextRunes := []rune(context)

	// 前后检查是否是姓名的合理位置
	if pos > 0 {
		prevChar := string(contextRunes[pos-1])
		if prevChar == "的" || prevChar == "是" || prevChar == "叫" || prevChar == "姓" {
			score += 0.1
		}
	}

	if pos+len(runes) < len(contextRunes) {
		nextChar := string(contextRunes[pos+len(runes)])
		if nextChar == "的" || nextChar == "是" || nextChar == "在" || nextChar == "，" || nextChar == "。" {
			score += 0.1
		}
	}

	return score
}

// detectOrganizations 检测组织机构
func (d *NLPPIIDetector) detectOrganizations(text string) []PIIEntity {
	var entities []PIIEntity

	// 使用常见组织名前缀列表来提高检测准确性
	orgPrefixes := []string{
		"腾讯", "阿里巴巴", "百度", "京东", "小米", "华为", "网易", "搜狐", "新浪",
		"字节跳动", "美团", "滴滴", "今日头条", "抖音", "快手", "拼多多",
		"平安", "工商", "建设", "农业", "中国银行", "招商", "中信", "光大",
	}

	// 首先检查已知的组织名
	for _, prefix := range orgPrefixes {
		for _, suffix := range d.companySuffixes {
			fullName := prefix + suffix
			if strings.Contains(text, fullName) {
				// 找到组织名的位置
				start := strings.Index(text, fullName)
				if start != -1 {
					score := d.calculateOrganizationScore(fullName)
					if score > 0.3 {
						entities = append(entities, PIIEntity{
							Text:     fullName,
							Type:     PIITypeOrganization,
							StartPos: start,
							EndPos:   start + len(fullName),
							Score:    score,
						})
					}
				}
			}
		}
	}

	// 然后使用正则表达式进行通用检测
	for _, suffix := range d.companySuffixes {
		// 改进的正则表达式，避免匹配过长的不相关内容
		pattern := regexp.MustCompile(`([\p{Han}]{2,8}` + regexp.QuoteMeta(suffix) + `)`)
		matches := pattern.FindAllStringSubmatchIndex(text, -1)

		for _, match := range matches {
			if len(match) >= 4 {
				start, end := match[2], match[3]
				original := text[start:end]

				score := d.calculateOrganizationScore(original)
				if score > 0.3 {
					entities = append(entities, PIIEntity{
						Text:     original,
						Type:     PIITypeOrganization,
						StartPos: start,
						EndPos:   end,
						Score:    score,
					})
				}
			}
		}
	}

	return entities
}

// calculateOrganizationScore 计算组织机构置信度
func (d *NLPPIIDetector) calculateOrganizationScore(name string) float64 {
	score := 0.0

	// 包含明确的组织后缀
	hasSuffix := false
	for _, suffix := range d.companySuffixes {
		if strings.HasSuffix(name, suffix) {
			score += 0.7 // 提高后缀分数
			hasSuffix = true
			break
		}
	}

	// 长度检查
	runes := []rune(name)
	if len(runes) >= 3 && len(runes) <= 15 { // 放宽长度要求
		score += 0.2
	}

	// 智能排除词检查 - 只在开头或结尾的明显非组织词才扣分
	excludeWords := []string{"我家", "在", "我", "的", "是", "一家", "大", "是", "上班", "工作", "的地方"}

	// 检查开头是否有排除词
	for _, word := range excludeWords {
		if strings.HasPrefix(name, word) && len(name) > len(word) {
			score -= 0.6 // 如果以排除词开头，大幅扣分
			break
		}
	}

	// 检查结尾是否有明显非组织词（但不包含在后缀中的词）
	if !hasSuffix {
		excludeEndingWords := []string{"了", "的", "是", "啊", "呢", "吗"}
		for _, word := range excludeEndingWords {
			if strings.HasSuffix(name, word) && len(name) > len(word) {
				score -= 0.3
				break
			}
		}
	}

	return score
}

// detectLocations 检测地理位置
func (d *NLPPIIDetector) detectLocations(text string) []PIIEntity {
	var entities []PIIEntity

	// 排除的词汇，避免误检测
	excludeWords := []string{"身份证号", "身份证", "银行卡号", "手机号", "电话号码"}

	for _, suffix := range d.locationSuffixes {
		// 查找包含地理位置后缀的匹配
		pattern := regexp.MustCompile(`([\p{Han}]{1,8}` + regexp.QuoteMeta(suffix) + `)`)
		matches := pattern.FindAllStringSubmatchIndex(text, -1)

		for _, match := range matches {
			start, end := match[2], match[3]
			original := text[start:end]

			// 检查是否应该排除
			shouldExclude := false
			for _, word := range excludeWords {
				if strings.Contains(original, word) {
					shouldExclude = true
					break
				}
			}
			if shouldExclude {
				continue
			}

			score := d.calculateLocationScore(original)
			if score > 0.5 {
				entities = append(entities, PIIEntity{
					Text:     original,
					Type:     PIITypeLocation,
					StartPos: start,
					EndPos:   end,
					Score:    score,
				})
			}
		}
	}

	return entities
}

// calculateLocationScore 计算地理位置置信度
func (d *NLPPIIDetector) calculateLocationScore(location string) float64 {
	score := 0.0

	// 包含地理后缀
	for _, suffix := range d.locationSuffixes {
		if strings.HasSuffix(location, suffix) {
			score += 0.5
			break
		}
	}

	// 长度检查
	runes := []rune(location)
	if len(runes) >= 2 && len(runes) <= 8 {
		score += 0.2
	}

	return score
}

// detectAddresses 检测地址（复合地理位置）
func (d *NLPPIIDetector) detectAddresses(text string) []PIIEntity {
	var entities []PIIEntity

	// 地址模式：更灵活的地址检测
	// 1. 完整地址模式：省市区+街道门牌号
	pattern1 := regexp.MustCompile(`([\p{Han}]{8,50}(?:省|市|区|县)(?:[\p{Han}]{1,20}(?:路|街道|大街|大道)(?:[\p{Han}0-9]{1,10}(?:号|弄|巷|室)?)?)?)`)

	// 2. 简单地址模式：包含号、室等后缀
	pattern2 := regexp.MustCompile(`([\p{Han}]{6,50}(?:号|室|弄|巷|楼|单元)(?:[\p{Han}0-9\-]*)?)`)

	patterns := []*regexp.Regexp{pattern1, pattern2}

	for _, pattern := range patterns {
		matches := pattern.FindAllStringSubmatchIndex(text, -1)

		for _, match := range matches {
			start, end := match[0], match[1]
			original := text[start:end]

			// 避免匹配过短或过长的文本
			runes := []rune(original)
			if len(runes) < 6 || len(runes) > 50 {
				continue
			}

			// 避免包含明确的排除词
			excludeWords := []string{"身份证号", "身份证", "银行卡号", "手机号", "电话号码", "邮箱", "邮件"}
			shouldExclude := false
			for _, word := range excludeWords {
				if strings.Contains(original, word) {
					shouldExclude = true
					break
				}
			}
			if shouldExclude {
				continue
			}

			entities = append(entities, PIIEntity{
				Text:     original,
				Type:     PIITypeAddress,
				StartPos: start,
				EndPos:   end,
				Score:    0.85, // 比地理位置更高的置信度
			})
		}
	}

	return entities
}

// deduplicateAndSort 去重并排序实体
func (d *NLPPIIDetector) deduplicateAndSort(entities []PIIEntity) []PIIEntity {
	// 去重：移除重叠的实体，保留置信度更高的
	var result []PIIEntity
	for _, entity := range entities {
		overlap := false
		for i, existing := range result {
			if d.isOverlapping(entity, existing) {
				overlap = true
				// 保留置信度更高的
				if entity.Score > existing.Score {
					result[i] = entity
				}
				break
			}
		}
		if !overlap {
			result = append(result, entity)
		}
	}

	// 按位置排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].StartPos < result[j].StartPos
	})

	return result
}

// isOverlapping 检查两个实体是否重叠
func (d *NLPPIIDetector) isOverlapping(e1, e2 PIIEntity) bool {
	return !(e1.EndPos <= e2.StartPos || e2.EndPos <= e1.StartPos)
}

// DetectAndReplace 检测并替换PII信息
func (d *NLPPIIDetector) DetectAndReplace(text string, threshold float64) (string, []PIIEntity) {
	// 检测PII
	entities, err := d.DetectPII(text, threshold)
	if err != nil {
		return text, []PIIEntity{}
	}

	// 按位置倒序排序（从后往前替换）
	sort.Slice(entities, func(i, j int) bool {
		return entities[i].StartPos > entities[j].StartPos
	})

	result := text

	// 从后往前替换避免位置偏移
	for _, entity := range entities {
		replacement := d.getReplacement(string(entity.Type), entity.Text)
		result = result[:entity.StartPos] + replacement + result[entity.EndPos:]

		// 更新实体信息
		entity.Replaced = replacement
	}

	// 按正序排列返回
	sort.Slice(entities, func(i, j int) bool {
		return entities[i].StartPos < entities[j].StartPos
	})

	return result, entities
}

// getReplacement 获取替换文本
func (d *NLPPIIDetector) getReplacement(piiType, original string) string {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	key := fmt.Sprintf("%s_%s", piiType, original)

	// 如果已经生成过替换，使用相同的替换
	if replacement, exists := d.replacementMap[key]; exists {
		return replacement
	}

	// 生成基于哈希的一致性替换
	hash := sha256.Sum256([]byte(original))

	var replacement string
	switch piiType {
	case "PER":
		replacement = fmt.Sprintf("[姓名%x]", hash[:4])
	case "ORG":
		replacement = fmt.Sprintf("[公司%x]", hash[:4])
	case "TEL":
		replacement = "[手机号码]"
	case "EMAIL":
		replacement = "[邮箱地址]"
	case "ID":
		replacement = "[身份证号]"
	case "ADDR", "LOC":
		replacement = "[地址]"
	default:
		replacement = "[敏感信息]"
	}

	d.replacementMap[key] = replacement
	return replacement
}

// GetStatistics 获取检测统计信息
func (d *NLPPIIDetector) GetStatistics(text string) map[string]interface{} {
	entities, _ := d.DetectPII(text, 0.0)

	stats := make(map[string]interface{})
	typeCount := make(map[PIIType]int)

	for _, entity := range entities {
		typeCount[entity.Type]++
	}

	stats["total_entities"] = len(entities)
	stats["type_distribution"] = typeCount
	stats["replacement_map_size"] = len(d.replacementMap)

	return stats
}

// GetReplacementMap 获取替换映射
func (d *NLPPIIDetector) GetReplacementMap() map[string]string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	result := make(map[string]string)
	for k, v := range d.replacementMap {
		result[k] = v
	}
	return result
}

// ToJSON 将检测结果转换为JSON
func (d *NLPPIIDetector) ToJSON(entities []PIIEntity) (string, error) {
	data, err := json.MarshalIndent(entities, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}