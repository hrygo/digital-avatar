package ai

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"unicode"
)

// TextPreprocessor 文本预处理器
type TextPreprocessor struct {
	stopWords    map[string]bool
	tokenizer    *Tokenizer
	normalizer   *Normalizer
	featureExtractor *FeatureExtractor
}

// NewTextPreprocessor 创建文本预处理器
func NewTextPreprocessor() *TextPreprocessor {
	return &TextPreprocessor{
		stopWords:       initializeStopWords(),
		tokenizer:       NewTokenizer(),
		normalizer:      NewNormalizer(),
		featureExtractor: NewFeatureExtractor(),
	}
}

// PreprocessText 预处理文本
func (tp *TextPreprocessor) PreprocessText(text string) []string {
	// 1. 文本清理
	text = tp.cleanText(text)

	// 2. 分词
	tokens := tp.tokenizer.Tokenize(text)

	// 3. 标准化
	tokens = tp.normalizer.Normalize(tokens)

	// 4. 去除停用词
	tokens = tp.removeStopWords(tokens)

	// 5. 特征提取
	tokens = tp.featureExtractor.Extract(tokens)

	return tokens
}

// PreprocessDocuments 批量预处理文档
func (tp *TextPreprocessor) PreprocessDocuments(documents []string) [][]string {
	var results [][]string

	for _, doc := range documents {
		tokens := tp.PreprocessText(doc)
		if len(tokens) > 0 { // 只保留非空文档
			results = append(results, tokens)
		}
	}

	return results
}

// BuildVocabulary 构建词汇表
func (tp *TextPreprocessor) BuildVocabulary(documents [][]string) map[string]int {
	vocabulary := make(map[string]int)
	wordFreq := make(map[string]int)

	// 统计词频
	for _, doc := range documents {
		uniqueWords := make(map[string]bool)
		for _, word := range doc {
			if !uniqueWords[word] {
				wordFreq[word]++
				uniqueWords[word] = true
			}
		}
	}

	// 构建词汇表，过滤低频词
	minFreq := 2
	for word, freq := range wordFreq {
		if freq >= minFreq {
			vocabulary[word] = len(vocabulary)
		}
	}

	return vocabulary
}

// cleanText 清理文本
func (tp *TextPreprocessor) cleanText(text string) string {
	// 移除URL
	urlRegex := regexp.MustCompile(`http[s]?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\\(\\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+`)
	text = urlRegex.ReplaceAllString(text, "")

	// 移除邮箱
	emailRegex := regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`)
	text = emailRegex.ReplaceAllString(text, "")

	// 移除多余空白字符
	spaceRegex := regexp.MustCompile(`\s+`)
	text = spaceRegex.ReplaceAllString(text, " ")

	// 移除特殊字符，保留中文、英文、数字
	specialCharRegex := regexp.MustCompile(`[^\p{L}\p{N}\s]`)
	text = specialCharRegex.ReplaceAllString(text, "")

	return strings.TrimSpace(text)
}

// removeStopWords 去除停用词
func (tp *TextPreprocessor) removeStopWords(tokens []string) []string {
	var result []string

	for _, token := range tokens {
		if !tp.stopWords[token] && len(token) > 1 {
			result = append(result, token)
		}
	}

	return result
}

// Tokenizer 分词器
type Tokenizer struct {
	pattern *regexp.Regexp
}

// NewTokenizer 创建分词器
func NewTokenizer() *Tokenizer {
	// 匹配中文词汇、英文单词和数字
	pattern := regexp.MustCompile(`[\p{Han}]+|[a-zA-Z]+|\d+`)
	return &Tokenizer{
		pattern: pattern,
	}
}

// Tokenize 分词
func (t *Tokenizer) Tokenize(text string) []string {
	matches := t.pattern.FindAllStringSubmatch(text, -1)
	var result []string
	for _, match := range matches {
		if len(match) > 0 {
			result = append(result, match[0])
		}
	}
	return result
}

// Normalizer 标准化器
type Normalizer struct {
	minWordLength int
	maxWordLength int
}

// NewNormalizer 创建标准化器
func NewNormalizer() *Normalizer {
	return &Normalizer{
		minWordLength: 2,
		maxWordLength: 20,
	}
}

// Normalize 标准化词汇
func (n *Normalizer) Normalize(tokens []string) []string {
	var result []string

	for _, token := range tokens {
		// 转换为小写
		token = strings.ToLower(token)

		// 过滤长度
		if len(token) < n.minWordLength || len(token) > n.maxWordLength {
			continue
		}

		// 去除数字（如果纯数字）
		if n.isNumeric(token) {
			continue
		}

		result = append(result, token)
	}

	return result
}

// isNumeric 检查是否为纯数字
func (n *Normalizer) isNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// FeatureExtractor 特征提取器
type FeatureExtractor struct {
	tfidfEnabled bool
	bigramEnabled bool
}

// NewFeatureExtractor 创建特征提取器
func NewFeatureExtractor() *FeatureExtractor {
	return &FeatureExtractor{
		tfidfEnabled:  true,
		bigramEnabled: false, // 暂时禁用bigram以简化
	}
}

// Extract 提取特征
func (fe *FeatureExtractor) Extract(tokens []string) []string {
	if fe.bigramEnabled {
		return fe.extractBigrams(tokens)
	}
	return tokens
}

// extractBigrams 提取二元组
func (fe *FeatureExtractor) extractBigrams(tokens []string) []string {
	if len(tokens) < 2 {
		return tokens
	}

	var result []string
	result = append(result, tokens...) // 保留单个词

	// 添加二元组
	for i := 0; i < len(tokens)-1; i++ {
		bigram := tokens[i] + "_" + tokens[i+1]
		result = append(result, bigram)
	}

	return result
}

// calculateTFIDF 计算TF-IDF
func (fe *FeatureExtractor) calculateTFIDF(documents [][]string) map[string]map[string]float64 {
	// 计算词频
	tf := make(map[int]map[string]int)
	df := make(map[string]int)

	for docID, doc := range documents {
		tf[docID] = make(map[string]int)
		uniqueWords := make(map[string]bool)

		for _, word := range doc {
			tf[docID][word]++
			if !uniqueWords[word] {
				df[word]++
				uniqueWords[word] = true
			}
		}
	}

	// 计算TF-IDF
	tfidf := make(map[string]map[string]float64)
	totalDocs := len(documents)

	for docID, docTF := range tf {
		docKey := fmt.Sprintf("doc_%d", docID)
		tfidf[docKey] = make(map[string]float64)
		totalWords := len(documents[docID])

		for word, count := range docTF {
			tfValue := float64(count) / float64(totalWords)
			idfValue := math.Log(float64(totalDocs) / float64(df[word]))
			tfidf[docKey][word] = tfValue * idfValue
		}
	}

	return tfidf
}

// initializeStopWords 初始化停用词
func initializeStopWords() map[string]bool {
	stopWords := map[string]bool{
		// 中文停用词
		"的": true, "了": true, "在": true, "是": true, "我": true, "你": true, "他": true, "她": true, "它": true, "们": true,
		"这": true, "那": true, "有": true, "和": true, "与": true, "或": true, "但": true, "而": true, "就": true, "都": true,
		"也": true, "又": true, "还": true, "再": true, "会": true, "能": true, "可": true, "可以": true, "应该": true, "要": true,
		"想": true, "说": true, "看": true, "听": true, "去": true, "来": true, "回": true, "过": true, "到": true, "从": true,
		"为": true, "以": true, "对": true, "把": true, "被": true, "让": true, "使": true, "由": true, "给": true, "向": true,
		"关于": true, "对于": true, "根据": true, "按照": true, "除了": true, "除非": true, "只要": true, "只有": true, "不管": true, "无论": true,
		"不是": true, "没有": true, "不会": true, "不能": true, "不要": true, "不用": true, "不必": true, "不可": true, "不应": true, "不该": true,
		"一个": true, "一些": true, "所有": true, "每个": true, "任何": true, "其他": true, "另外": true, "各种": true, "整个": true, "部分": true,
		"什么": true, "怎么": true, "为什么": true, "哪里": true, "哪个": true, "多少": true, "几个": true, "第一": true, "第二": true, "第三": true,
		"非常": true, "特别": true, "很": true, "太": true, "更": true, "最": true, "比较": true, "相对": true, "绝对": true, "完全": true,
		"刚刚": true, "马上": true, "立刻": true, "现在": true, "以前": true, "以后": true, "今天": true, "昨天": true, "明天": true, "同时": true,
		"然后": true, "接着": true, "最后": true, "首先": true, "其次": true, "最终": true, "总之": true, "因此": true, "所以": true, "但是": true,

		// 英文停用词
		"the": true, "be": true, "to": true, "of": true, "and": true, "a": true, "in": true, "that": true, "have": true, "i": true,
		"it": true, "for": true, "not": true, "on": true, "with": true, "he": true, "as": true, "you": true, "do": true, "at": true,
		"this": true, "but": true, "his": true, "by": true, "from": true, "they": true, "we": true, "say": true, "her": true, "she": true,
		"or": true, "an": true, "will": true, "my": true, "one": true, "all": true, "would": true, "there": true, "their": true, "what": true,
		"so": true, "up": true, "out": true, "if": true, "about": true, "who": true, "get": true, "which": true, "go": true, "me": true,
		"when": true, "make": true, "can": true, "like": true, "time": true, "no": true, "just": true, "him": true, "know": true, "take": true,
		"people": true, "into": true, "year": true, "your": true, "good": true, "some": true, "could": true, "them": true, "see": true, "other": true,
		"than": true, "then": true, "now": true, "look": true, "only": true, "come": true, "its": true, "over": true, "think": true, "also": true,
		"back": true, "after": true, "use": true, "two": true, "how": true, "our": true, "work": true, "first": true, "well": true, "way": true,
		"even": true, "new": true, "want": true, "because": true, "any": true, "these": true, "give": true, "day": true, "most": true, "us": true,
	}

	return stopWords
}