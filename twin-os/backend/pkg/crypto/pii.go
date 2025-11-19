package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// PIIDetector PII检测器
type PIIDetector struct {
	namePatterns      []*regexp.Regexp
	phonePatterns     []*regexp.Regexp
	emailPatterns     []*regexp.Regexp
	addressPatterns   []*regexp.Regexp
	idCardPatterns    []*regexp.Regexp
	companyPatterns   []*regexp.Regexp
	replacementMap    map[string]string
}

// NewPIIDetector 创建PII检测器
func NewPIIDetector() *PIIDetector {
	detector := &PIIDetector{
		replacementMap: make(map[string]string),
	}

	// 初始化正则表达式模式
	detector.initPatterns()

	return detector
}

// initPatterns 初始化检测模式
func (d *PIIDetector) initPatterns() {
	// 中文姓名模式（2-4个中文字符）
	d.namePatterns = append(d.namePatterns,
		regexp.MustCompile(`[\p{Han}]{2,4}`),
	)

	// 电话号码模式
	d.phonePatterns = append(d.phonePatterns,
		regexp.MustCompile(`1[3-9]\d{9}`),                     // 手机号
		regexp.MustCompile(`\d{3,4}-\d{7,8}`),               // 固定电话
		regexp.MustCompile(`\+\d{2,3}\s?\d{3,4}\s?\d{7,8}`), // 国际号码
	)

	// 邮箱模式
	d.emailPatterns = append(d.emailPatterns,
		regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`),
	)

	// 身份证模式
	d.idCardPatterns = append(d.idCardPatterns,
		regexp.MustCompile(`\d{17}[\dXx]`), // 18位身份证
		regexp.MustCompile(`\d{15}`),       // 15位身份证
	)

	// 公司名模式（简化版）
	d.companyPatterns = append(d.companyPatterns,
		regexp.MustCompile(`[\p{Han}]+(?:有限公司|股份有限公司|集团|科技|网络|信息|文化|教育|医疗|金融|投资)`),
		regexp.MustCompile(`[\p{Han}]+(?:公司|企业|机构|中心)`),
	)
}

// DetectAndReplace 检测并替换PII信息
func (d *PIIDetector) DetectAndReplace(text string) (string, PIIReport) {
	report := PIIReport{
		OriginalText: text,
		Detections:   make([]PIIDetection, 0),
	}

	result := text

	// 检测和替换姓名
	result = d.detectAndReplaceNames(result, &report)

	// 检测和替换电话号码
	result = d.detectAndReplacePhones(result, &report)

	// 检测和替换邮箱
	result = d.detectAndReplaceEmails(result, &report)

	// 检测和替换身份证
	result = d.detectAndReplaceIDCards(result, &report)

	// 检测和替换公司名
	result = d.detectAndReplaceCompanies(result, &report)

	report.ProcessedText = result
	report.TotalDetections = len(report.Detections)

	return result, report
}

// detectAndReplaceNames 检测和替换姓名
func (d *PIIDetector) detectAndReplaceNames(text string, report *PIIReport) string {
	for _, pattern := range d.namePatterns {
		matches := pattern.FindAllStringSubmatchIndex(text, -1)
		for _, match := range matches {
			if len(match) >= 2 {
				start, end := match[0], match[1]
				original := text[start:end]

				// 过滤掉一些明显不是姓名的词
				if d.isLikelyName(original) {
					replacement := d.getReplacement("name", original)
					text = text[:start] + replacement + text[end:]

					detection := PIIDetection{
						Type:      "name",
						Original:  original,
						Replaced:  replacement,
						StartPos:  start,
						EndPos:    end,
						Confidence: d.calculateConfidence("name", original),
					}
					report.Detections = append(report.Detections, detection)
				}
			}
		}
	}
	return text
}

// detectAndReplacePhones 检测和替换电话号码
func (d *PIIDetector) detectAndReplacePhones(text string, report *PIIReport) string {
	for _, pattern := range d.phonePatterns {
		matches := pattern.FindAllStringSubmatchIndex(text, -1)
		for _, match := range matches {
			if len(match) >= 2 {
				start, end := match[0], match[1]
				original := text[start:end]

				replacement := d.getReplacement("phone", original)
				text = text[:start] + replacement + text[end:]

				detection := PIIDetection{
					Type:      "phone",
					Original:  original,
					Replaced:  replacement,
					StartPos:  start,
					EndPos:    end,
					Confidence: 0.95, // 电话号码匹配置信度高
				}
				report.Detections = append(report.Detections, detection)
			}
		}
	}
	return text
}

// detectAndReplaceEmails 检测和替换邮箱
func (d *PIIDetector) detectAndReplaceEmails(text string, report *PIIReport) string {
	for _, pattern := range d.emailPatterns {
		matches := pattern.FindAllStringSubmatchIndex(text, -1)
		for _, match := range matches {
			if len(match) >= 2 {
				start, end := match[0], match[1]
				original := text[start:end]

				replacement := d.getReplacement("email", original)
				text = text[:start] + replacement + text[end:]

				detection := PIIDetection{
					Type:      "email",
					Original:  original,
					Replaced:  replacement,
					StartPos:  start,
					EndPos:    end,
					Confidence: 0.95, // 邮箱匹配置信度高
				}
				report.Detections = append(report.Detections, detection)
			}
		}
	}
	return text
}

// detectAndReplaceIDCards 检测和替换身份证
func (d *PIIDetector) detectAndReplaceIDCards(text string, report *PIIReport) string {
	for _, pattern := range d.idCardPatterns {
		matches := pattern.FindAllStringSubmatchIndex(text, -1)
		for _, match := range matches {
			if len(match) >= 2 {
				start, end := match[0], match[1]
				original := text[start:end]

				// 验证身份证格式
				if d.isValidIDCard(original) {
					replacement := d.getReplacement("idcard", original)
					text = text[:start] + replacement + text[end:]

					detection := PIIDetection{
						Type:      "idcard",
						Original:  original,
						Replaced:  replacement,
						StartPos:  start,
						EndPos:    end,
						Confidence: 0.9,
					}
					report.Detections = append(report.Detections, detection)
				}
			}
		}
	}
	return text
}

// detectAndReplaceCompanies 检测和替换公司名
func (d *PIIDetector) detectAndReplaceCompanies(text string, report *PIIReport) string {
	for _, pattern := range d.companyPatterns {
		matches := pattern.FindAllStringSubmatchIndex(text, -1)
		for _, match := range matches {
			if len(match) >= 2 {
				start, end := match[0], match[1]
				original := text[start:end]

				replacement := d.getReplacement("company", original)
				text = text[:start] + replacement + text[end:]

				detection := PIIDetection{
					Type:      "company",
					Original:  original,
					Replaced:  replacement,
					StartPos:  start,
					EndPos:    end,
					Confidence: 0.7, // 公司名匹配置信度中等
				}
				report.Detections = append(report.Detections, detection)
			}
		}
	}
	return text
}

// isLikelyName 判断是否可能是姓名
func (d *PIIDetector) isLikelyName(text string) bool {
	// 排除一些明显不是姓名的词
	excludeWords := []string{
		"微信", "公司", "有限公司", "集团", "科技", "网络", "信息",
		"文化", "教育", "医疗", "金融", "投资", "产品", "项目",
		"系统", "平台", "服务", "客户", "用户", "经理", "总监",
		"老师", "医生", "律师", "工程师", "设计师", "开发",
	}

	for _, word := range excludeWords {
		if strings.Contains(text, word) {
			return false
		}
	}

	// 检查字符长度
	runes := []rune(text)
	if len(runes) < 2 || len(runes) > 4 {
		return false
	}

	return true
}

// isValidIDCard 验证身份证格式
func (d *PIIDetector) isValidIDCard(idCard string) bool {
	if len(idCard) == 18 {
		// 18位身份证验证
		return d.validateIDCard18(idCard)
	} else if len(idCard) == 15 {
		// 15位身份证验证
		return d.validateIDCard15(idCard)
	}
	return false
}

// validateIDCard18 验证18位身份证
func (d *PIIDetector) validateIDCard18(idCard string) bool {
	if len(idCard) != 18 {
		return false
	}

	// 检查前17位是否为数字
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

	return true
}

// validateIDCard15 验证15位身份证
func (d *PIIDetector) validateIDCard15(idCard string) bool {
	if len(idCard) != 15 {
		return false
	}

	// 检查是否全为数字
	for _, char := range idCard {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}

// getReplacement 获取替换文本
func (d *PIIDetector) getReplacement(piiType, original string) string {
	key := fmt.Sprintf("%s_%s", piiType, original)

	// 如果已经生成过替换，使用相同的替换
	if replacement, exists := d.replacementMap[key]; exists {
		return replacement
	}

	// 生成新的替换
	var replacement string
	switch piiType {
	case "name":
		replacement = fmt.Sprintf("[姓名%d]", len(d.replacementMap)+1)
	case "phone":
		replacement = "[手机号码]"
	case "email":
		replacement = "[邮箱地址]"
	case "idcard":
		replacement = "[身份证号]"
	case "company":
		replacement = fmt.Sprintf("[公司%d]", len(d.replacementMap)+1)
	default:
		replacement = "[敏感信息]"
	}

	d.replacementMap[key] = replacement
	return replacement
}

// calculateConfidence 计算检测置信度
func (d *PIIDetector) calculateConfidence(piiType, original string) float64 {
	switch piiType {
	case "phone":
		return 0.95
	case "email":
		return 0.95
	case "idcard":
		return 0.9
	case "name":
		// 根据长度和复杂度计算置信度
		runes := []rune(original)
		if len(runes) == 2 || len(runes) == 3 {
			return 0.8
		}
		return 0.6
	case "company":
		return 0.7
	default:
		return 0.5
	}
}

// PIIReport PII检测报告
type PIIReport struct {
	OriginalText   string        `json:"original_text"`
	ProcessedText  string        `json:"processed_text"`
	TotalDetections int          `json:"total_detections"`
	Detections     []PIIDetection `json:"detections"`
}

// PIIDetection PII检测结果
type PIIDetection struct {
	Type      string  `json:"type"`
	Original  string  `json:"original"`
	Replaced  string  `json:"replaced"`
	StartPos  int     `json:"start_pos"`
	EndPos    int     `json:"end_pos"`
	Confidence float64 `json:"confidence"`
}

// EncryptData 加密数据
func EncryptData(plaintext []byte, key string) ([]byte, error) {
	// 使用SHA256生成32字节的密钥
	keyBytes := sha256.Sum256([]byte(key))

	block, err := aes.NewCipher(keyBytes[:])
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// 创建GCM模式的AEAD
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// 生成随机nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// 加密数据
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return ciphertext, nil
}

// DecryptData 解密数据
func DecryptData(ciphertext []byte, key string) ([]byte, error) {
	// 使用SHA256生成32字节的密钥
	keyBytes := sha256.Sum256([]byte(key))

	block, err := aes.NewCipher(keyBytes[:])
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// 创建GCM模式的AEAD
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	// 分离nonce和实际密文
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// 解密数据
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// EncryptToBase64 加密并转换为Base64
func EncryptToBase64(plaintext string, key string) (string, error) {
	ciphertext, err := EncryptData([]byte(plaintext), key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptFromBase64 从Base64解密
func DecryptFromBase64(ciphertext string, key string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	plaintext, err := DecryptData(data, key)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}