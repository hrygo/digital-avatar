package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

// PIIDetector PII检测器（兼容接口）
// 注意：这是一个兼容层，内部使用新的NLPPIIDetector
// 建议新代码直接使用NLPPIIDetector
type PIIDetector struct {
	nlpDetector *NLPPIIDetector
}

// PIIReport PII检测报告（兼容结构）
type PIIReport struct {
	OriginalText    string        `json:"original_text"`
	ProcessedText   string        `json:"processed_text"`
	TotalDetections int           `json:"total_detections"`
	Detections      []PIIDetection `json:"detections"`
}

// PIIDetection PII检测结果（兼容结构）
type PIIDetection struct {
	Type      string  `json:"type"`
	Original  string  `json:"original"`
	Replaced  string  `json:"replaced"`
	StartPos  int     `json:"start_pos"`
	EndPos    int     `json:"end_pos"`
	Confidence float64 `json:"confidence"`
}

// NewPIIDetector 创建PII检测器（兼容接口）
// 注意：此函数为了向后兼容而保留，建议使用NewNLPPIIDetector
func NewPIIDetector() *PIIDetector {
	return &PIIDetector{
		nlpDetector: NewNLPPIIDetector(),
	}
}

// DetectAndReplace 检测并替换PII信息（兼容接口）
func (d *PIIDetector) DetectAndReplace(text string) (string, PIIReport) {
	// 使用新的NLP检测器，置信度阈值设为0.6
	processedText, entities := d.nlpDetector.DetectAndReplace(text, 0.6)

	// 转换新的实体格式到旧的检测格式
	detections := make([]PIIDetection, 0, len(entities))
	for _, entity := range entities {
		detection := PIIDetection{
			Type:      string(entity.Type),
			Original:  entity.Text,
			Replaced:  entity.Replaced,
			StartPos:  entity.StartPos,
			EndPos:    entity.EndPos,
			Confidence: entity.Score,
		}
		detections = append(detections, detection)
	}

	report := PIIReport{
		OriginalText:    text,
		ProcessedText:   processedText,
		TotalDetections: len(detections),
		Detections:      detections,
	}

	return processedText, report
}

// GetReplacementMap 获取替换映射（兼容接口）
func (d *PIIDetector) GetReplacementMap() map[string]string {
	return d.nlpDetector.GetReplacementMap()
}

// ClearReplacementMap 清空替换映射（兼容接口）
func (d *PIIDetector) ClearReplacementMap() {
	// 新的NLP检测器没有这个方法，我们可以重新创建
	d.nlpDetector = NewNLPPIIDetector()
}

// GetStatistics 获取统计信息（兼容接口）
func (d *PIIDetector) GetStatistics(text string) map[string]interface{} {
	stats := d.nlpDetector.GetStatistics(text)

	// 为向后兼容，确保返回的统计信息包含旧的字段
	if _, exists := stats["total_detections"]; !exists {
		if total, ok := stats["total_entities"].(int); ok {
			stats["total_detections"] = total
		}
	}

	return stats
}

// SetConfidenceThreshold 设置置信度阈值（新增方法，用于兼容）
func (d *PIIDetector) SetConfidenceThreshold(threshold float64) {
	// 新的NLP检测器允许在调用时指定阈值，这里为了兼容提供一个方法
	// 实际使用时，阈值会在DetectAndReplace中指定
}

// GetConfidenceThreshold 获取当前置信度阈值（新增方法，用于兼容）
func (d *PIIDetector) GetConfidenceThreshold() float64 {
	// 返回默认阈值
	return 0.6
}

// ==================== 加密相关功能 ====================
// 以下是与PII检测无关的加密功能，保留在原文件中

// EncryptData 使用AES加密数据
func EncryptData(plaintext []byte, key string) ([]byte, error) {
	// 生成32字节的密钥
	keyBytes := sha256.Sum256([]byte(key))

	block, err := aes.NewCipher(keyBytes[:])
	if err != nil {
		return nil, err
	}

	// 创建GCM模式的加密器
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 生成随机nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// 加密数据
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// DecryptData 使用AES解密数据
func DecryptData(ciphertext []byte, key string) ([]byte, error) {
	// 生成32字节的密钥
	keyBytes := sha256.Sum256([]byte(key))

	block, err := aes.NewCipher(keyBytes[:])
	if err != nil {
		return nil, err
	}

	// 创建GCM模式的解密器
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext_bytes := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext_bytes, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// EncryptToBase64 加密数据并返回Base64编码
func EncryptToBase64(plaintext string, key string) (string, error) {
	ciphertext, err := EncryptData([]byte(plaintext), key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptFromBase64 从Base64编码解密数据
func DecryptFromBase64(ciphertext string, key string) (string, error) {
	ciphertext_bytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	plaintext, err := DecryptData(ciphertext_bytes, key)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}