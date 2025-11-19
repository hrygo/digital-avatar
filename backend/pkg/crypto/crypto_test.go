package crypto

import (
	"bytes"
	"encoding/base64"
	"testing"
)

// TestEncryptDecryptData 测试基本加密解密功能
func TestEncryptDecryptData(t *testing.T) {
	plaintext := []byte("这是一个需要加密的敏感信息")
	key := "test-key-12345"

	// 加密
	ciphertext, err := EncryptData(plaintext, key)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// 验证密文与原文不同
	if bytes.Equal(plaintext, ciphertext) {
		t.Error("Ciphertext should be different from plaintext")
	}

	// 解密
	decrypted, err := DecryptData(ciphertext, key)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	// 验证解密结果与原文相同
	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("Decrypted text mismatch. Original: %s, Decrypted: %s", string(plaintext), string(decrypted))
	}
}

// TestEncryptDecryptBase64 测试Base64编码的加密解密
func TestEncryptDecryptBase64(t *testing.T) {
	plaintext := "这是一个需要加密的敏感信息"
	key := "test-key-12345"

	// 加密到Base64
	ciphertextBase64, err := EncryptToBase64(plaintext, key)
	if err != nil {
		t.Fatalf("Encryption to Base64 failed: %v", err)
	}

	// 验证结果不是原文
	if ciphertextBase64 == plaintext {
		t.Error("Ciphertext should be different from plaintext")
	}

	// 验证是有效的Base64
	_, err = base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		t.Errorf("Result should be valid Base64: %v", err)
	}

	// 从Base64解密
	decrypted, err := DecryptFromBase64(ciphertextBase64, key)
	if err != nil {
		t.Fatalf("Decryption from Base64 failed: %v", err)
	}

	// 验证解密结果与原文相同
	if decrypted != plaintext {
		t.Errorf("Decrypted text mismatch. Original: %s, Decrypted: %s", plaintext, decrypted)
	}
}

// TestEncryptEmptyData 测试空数据加密
func TestEncryptEmptyData(t *testing.T) {
	plaintext := []byte{}
	key := "test-key-12345"

	ciphertext, err := EncryptData(plaintext, key)
	if err != nil {
		t.Fatalf("Empty data encryption failed: %v", err)
	}

	decrypted, err := DecryptData(ciphertext, key)
	if err != nil {
		t.Fatalf("Empty data decryption failed: %v", err)
	}

	if len(decrypted) != 0 {
		t.Error("Decrypted empty data should be empty")
	}
}

// TestEncryptWithInvalidInputs 测试无效输入处理
func TestEncryptWithInvalidInputs(t *testing.T) {
	testCases := []struct {
		name     string
		plaintext []byte
		key      string
		valid    bool
	}{
		{"Empty key", []byte("data"), "", false},
		{"Valid key", []byte("data"), "valid-key", true},
		{"Empty data", []byte{}, "key", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := EncryptData(tc.plaintext, tc.key)

			if tc.valid && err != nil {
				t.Errorf("Valid input should not cause error: %v", err)
			}

			if !tc.valid && err == nil {
				t.Error("Invalid input should cause error")
			}
		})
	}
}

// TestDecryptWithWrongKey 测试错误密钥解密
func TestDecryptWithWrongKey(t *testing.T) {
	plaintext := []byte("敏感信息")
	correctKey := "correct-key-123"
	wrongKey := "wrong-key-456"

	// 用正确密钥加密
	ciphertext, err := EncryptData(plaintext, correctKey)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// 用错误密钥解密
	_, err = DecryptData(ciphertext, wrongKey)
	if err == nil {
		t.Error("Decryption with wrong key should fail")
	}
}

// TestDecryptInvalidCiphertext 测试无效密文解密
func TestDecryptInvalidCiphertext(t *testing.T) {
	key := "test-key-12345"

	testCases := []struct {
		name        string
		ciphertext  []byte
		shouldFail  bool
	}{
		{"Empty ciphertext", []byte{}, true},
		{"Nil ciphertext", nil, true},
		{"Too short ciphertext", []byte("short"), true},
		{"Invalid data", []byte("invalid_ciphertext_data"), true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := DecryptData(tc.ciphertext, key)

			if tc.shouldFail && err == nil {
				t.Error("Invalid ciphertext should cause decryption error")
			}
		})
	}
}

// TestMultipleEncryptions 测试多次加密的结果不同
func TestMultipleEncryptions(t *testing.T) {
	plaintext := []byte("需要加密的数据")
	key := "test-key-12345"

	// 加密多次
	var ciphertexts [][]byte
	for i := 0; i < 3; i++ {
		ciphertext, err := EncryptData(plaintext, key)
		if err != nil {
			t.Fatalf("Encryption %d failed: %v", i, err)
		}
		ciphertexts = append(ciphertexts, ciphertext)
	}

	// 验证所有密文都不相同（由于随机nonce）
	for i := 0; i < len(ciphertexts); i++ {
		for j := i + 1; j < len(ciphertexts); j++ {
			if bytes.Equal(ciphertexts[i], ciphertexts[j]) {
				t.Errorf("Ciphertext %d and %d should be different", i, j)
			}
		}
	}

	// 验证所有密文都能正确解密
	for i, ciphertext := range ciphertexts {
		decrypted, err := DecryptData(ciphertext, key)
		if err != nil {
			t.Errorf("Decryption %d failed: %v", i, err)
		}
		if !bytes.Equal(plaintext, decrypted) {
			t.Errorf("Decryption %d result mismatch", i)
		}
	}
}

// TestLargeDataEncryption 测试大量数据加密
func TestLargeDataEncryption(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large data test in short mode")
	}

	// 创建10KB的测试数据
	data := make([]byte, 10*1024)
	for i := range data {
		data[i] = byte(i % 256)
	}

	key := "test-key-12345"

	// 加密
	ciphertext, err := EncryptData(data, key)
	if err != nil {
		t.Fatalf("Large data encryption failed: %v", err)
	}

	// 解密
	decrypted, err := DecryptData(ciphertext, key)
	if err != nil {
		t.Fatalf("Large data decryption failed: %v", err)
	}

	// 验证数据完整性
	if !bytes.Equal(data, decrypted) {
		t.Error("Large data encryption/decryption failed")
	}

	t.Logf("Successfully encrypted/decrypted %d bytes", len(data))
}

// TestPIIConsistency 测试PII替换的一致性
func TestPIIConsistency(t *testing.T) {
	detector := NewNLPPIIDetector()

	// 测试文本
	text := "联系人：张三，手机：13812345678，邮箱：zhang@example.com"

	// 第一次检测替换
	processed1, entities1 := detector.DetectAndReplace(text, 0.7)

	// 第二次检测替换
	processed2, entities2 := detector.DetectAndReplace(text, 0.7)

	// 验证处理结果相同
	if processed1 != processed2 {
		t.Error("Same input should produce same output")
	}

	// 验证检测到的实体数量相同
	if len(entities1) != len(entities2) {
		t.Error("Same input should detect same number of entities")
	}

	t.Logf("Consistency test passed with %d entities", len(entities1))
}