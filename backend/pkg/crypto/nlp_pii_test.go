package crypto

import (
	"regexp"
	"strings"
	"testing"
)

func TestNLPPIIDetector_DetectAndReplace(t *testing.T) {
	detector := NewNLPPIIDetector()

	tests := []struct {
		name          string
		inputText     string
		expectedCount int
		shouldContain []string
	}{
		{
			name:          "Basic Name and Phone",
			inputText:     "我的名字是张三，电话是13812345678。",
			expectedCount: 2,
			shouldContain: []string{"[姓名", "[手机号码]"},
		},
		{
			name:          "Email and ID Card",
			inputText:     "我的邮箱是test@example.com，身份证号是11010119900307123X。",
			expectedCount: 2,
			shouldContain: []string{"[邮箱地址]", "[身份证号]"},
		},
		{
			name:          "Organization",
			inputText:     "腾讯有限公司是一家大公司。",
			expectedCount: 1,
			shouldContain: []string{"[公司"},
		},
		{
			name:          "Address",
			inputText:     "我住在北京朝阳区建国门外大街1号。",
			expectedCount: 1,
			shouldContain: []string{"[地址]"},
		},
		{
			name:          "Mixed PII",
			inputText:     "联系人：李四，手机：13900001111，邮箱：lisi@company.com，在阿里巴巴集团工作。",
			expectedCount: 4,
			shouldContain: []string{"[姓名", "[手机号码]", "[邮箱地址]", "[公司"},
		},
		{
			name:          "No PII",
			inputText:     "今天天气很好，适合出门散步。",
			expectedCount: 0,
			shouldContain: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processedText, entities := detector.DetectAndReplace(tt.inputText, 0.6)

			// 检查检测数量
			if len(entities) != tt.expectedCount {
				t.Errorf("Test %s: Expected %d entities, got %d", tt.name, tt.expectedCount, len(entities))
			}

			// 检查是否包含期望的替换标签
			for _, expected := range tt.shouldContain {
				if !strings.Contains(processedText, expected) {
					t.Errorf("Test %s: Expected processed text to contain '%s', but got: %s", tt.name, expected, processedText)
				}
			}

			// 打印详细信息用于调试
			if len(entities) > 0 {
				t.Logf("Test %s - Detected entities:", tt.name)
				for i, entity := range entities {
					t.Logf("  %d: %s (%s) -> %s (score: %.2f)",
						i+1, entity.Text, entity.Type, entity.Replaced, entity.Score)
				}
				t.Logf("Processed: %s", processedText)
			}
		})
	}
}

func TestNLPPIIDetector_ConsistentReplacement(t *testing.T) {
	detector := NewNLPPIIDetector()

	// 测试姓名一致性
	nameTests := []string{
		"我的名字是张三。",
		"你好，张三。",
		"再次见到张三。",
		"张三是我的朋友。",
	}

	var nameReplacements []string
	for _, text := range nameTests {
		processed, _ := detector.DetectAndReplace(text, 0.6)

		// 查找姓名替换
		nameRe := regexp.MustCompile(`\[姓名[a-f0-9]{8}\]`)
		if match := nameRe.FindString(processed); match != "" {
			nameReplacements = append(nameReplacements, match)
		}
	}

	// 检查一致性
	if len(nameReplacements) > 1 {
		for i := 1; i < len(nameReplacements); i++ {
			if nameReplacements[i] != nameReplacements[0] {
				t.Errorf("Expected consistent name replacement, but got different values: %s vs %s",
					nameReplacements[0], nameReplacements[i])
			}
		}
	} else if len(nameReplacements) == 0 {
		t.Error("Expected to find name replacements, but found none")
	}

	// 测试不同姓名的不同替换
	diffTests := []string{
		"这是李四。",
		"这是王五。",
	}

	var diffReplacements []string
	for _, text := range diffTests {
		processed, _ := detector.DetectAndReplace(text, 0.6)
		nameRe := regexp.MustCompile(`\[姓名[a-f0-9]{8}\]`)
		if match := nameRe.FindString(processed); match != "" {
			diffReplacements = append(diffReplacements, match)
		}
	}

	if len(diffReplacements) == 2 && diffReplacements[0] == diffReplacements[1] {
		t.Errorf("Expected different replacements for different names, but got same: %s", diffReplacements[0])
	}
}

func TestNLPPIIDetector_Performance(t *testing.T) {
	detector := NewNLPPIIDetector()

	// 构造较长的测试文本
	text := strings.Repeat("我的名字是张三，电话是13812345678，邮箱是zhang@example.com。", 10)

	// 测试性能
	t.Run("Performance Test", func(t *testing.T) {
		processed, entities := detector.DetectAndReplace(text, 0.6)

		t.Logf("Processed text length: %d characters", len(processed))
		t.Logf("Detected %d entities", len(entities))
		t.Logf("Replacement map size: %d", len(detector.GetReplacementMap()))

		// 验证结果合理性
		if len(entities) == 0 {
			t.Error("Expected to detect entities in repeated text, but found none")
		}

		if !strings.Contains(processed, "[姓名") || !strings.Contains(processed, "[手机号码]") {
			t.Error("Expected to find replacement tags in processed text")
		}
	})
}

func TestNLPPIIDetector_EdgeCases(t *testing.T) {
	detector := NewNLPPIIDetector()

	tests := []struct {
		name     string
		text     string
		validate func(t *testing.T, text string, entities []PIIEntity)
	}{
		{
			name: "Empty Text",
			text: "",
			validate: func(t *testing.T, text string, entities []PIIEntity) {
				if len(entities) != 0 {
					t.Errorf("Expected no entities in empty text, got %d", len(entities))
				}
			},
		},
		{
			name: "Only Spaces",
			text: "   \t\n   ",
			validate: func(t *testing.T, text string, entities []PIIEntity) {
				if len(entities) != 0 {
					t.Errorf("Expected no entities in whitespace-only text, got %d", len(entities))
				}
			},
		},
		{
			name: "Mixed Languages",
			text: "My name is 张三 and my email is john@example.com.",
			validate: func(t *testing.T, text string, entities []PIIEntity) {
				// 应该检测到中文姓名和英文邮箱
				hasName := false
				hasEmail := false

				for _, entity := range entities {
					if entity.Type == PIITypePerson {
						hasName = true
					}
					if entity.Type == PIITypeEmail {
						hasEmail = true
					}
				}

				if !hasName {
					t.Error("Expected to detect Chinese name in mixed language text")
				}
				if !hasEmail {
					t.Error("Expected to detect email in mixed language text")
				}
			},
		},
		{
			name: "Border Cases",
			text: "张三",
			validate: func(t *testing.T, text string, entities []PIIEntity) {
				if len(entities) == 0 || entities[0].Type != PIITypePerson {
					t.Error("Expected to detect name in minimal text")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, entities := detector.DetectAndReplace(tt.text, 0.6)
			tt.validate(t, tt.text, entities)
		})
	}
}

func TestNLPPIIDetector_Statistics(t *testing.T) {
	detector := NewNLPPIIDetector()

	text := "联系人：张三（13812345678），李四（lisi@example.com），在阿里巴巴集团工作，地址是北京市朝阳区建国门外大街1号。"

	stats := detector.GetStatistics(text)

	// 检查统计信息
	if total, ok := stats["total_entities"].(int); !ok || total == 0 {
		t.Error("Expected non-zero total entities count")
	}

	if typeDist, ok := stats["type_distribution"].(map[PIIType]int); ok {
		expectedTypes := []PIIType{PIITypePerson, PIITypePhone, PIITypeEmail, PIITypeOrganization, PIITypeAddress}
		for _, expectedType := range expectedTypes {
			if count, exists := typeDist[expectedType]; !exists || count == 0 {
				t.Errorf("Expected to detect entities of type %s", expectedType)
			}
		}
	} else {
		t.Error("Expected type distribution in statistics")
	}

	t.Logf("Statistics: %+v", stats)
}

func TestNLPPIIDetector_JSONOutput(t *testing.T) {
	detector := NewNLPPIIDetector()

	text := "我的名字是张三，电话是13812345678。"
	_, entities := detector.DetectAndReplace(text, 0.6)

	jsonStr, err := detector.ToJSON(entities)
	if err != nil {
		t.Fatalf("Failed to convert entities to JSON: %v", err)
	}

	// 验证JSON包含基本信息
	if !strings.Contains(jsonStr, "张三") || !strings.Contains(jsonStr, "13812345678") {
		t.Errorf("JSON output missing expected content: %s", jsonStr)
	}

	if !strings.Contains(jsonStr, "\"type\"") || !strings.Contains(jsonStr, "\"text\"") {
		t.Errorf("JSON output missing expected fields: %s", jsonStr)
	}

	t.Logf("JSON output:\n%s", jsonStr)
}

// BenchmarkPIIDetection 基准测试
func BenchmarkNLPPIIDetector_DetectAndReplace(b *testing.B) {
	detector := NewNLPPIIDetector()
	text := "我的名字是张三，电话是13812345678，邮箱是zhang@example.com，在腾讯有限公司工作。"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		detector.DetectAndReplace(text, 0.6)
	}
}

// 基准测试不同长度的文本
func BenchmarkNLPPIIDetector_DifferentLengths(b *testing.B) {
	detector := NewNLPPIIDetector()

	shortText := "张三 13812345678"
	mediumText := strings.Repeat(shortText+" ", 10)   // 10个重复
	longText := strings.Repeat(shortText+" ", 100)    // 100个重复

	tests := []struct {
		name string
		text string
	}{
		{"Short", shortText},
		{"Medium", mediumText},
		{"Long", longText},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				detector.DetectAndReplace(tt.text, 0.6)
			}
		})
	}
}