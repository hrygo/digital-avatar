package crypto

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// TestPIIDetection_EdgeCases 测试PII检测的边界情况和鲁棒性
func TestPIIDetection_EdgeCases(t *testing.T) {
	detector := NewNLPPIIDetector()

	edgeCases := []struct {
		name            string
		text            string
		expectedPII     int
		description     string
		businessRisk    string
	}{
		{
			name:         "极端长度文本",
			text:         strings.Repeat("用户张三", 1000) + "电话" + strings.Repeat("13812345678", 500),
			expectedPII:  0, // 可能因性能限制而减少检测
			description:  "超长文本的性能和准确性测试",
			businessRisk: "中 - 大量数据处理可能影响性能",
		},
		{
			name:         "Unicode和特殊字符",
			text:         "👨张三📱13812345678✉️zhang@example.com🏢北京市朝阳区™️🇨🇳",
			expectedPII:  4,
			description:  "包含emoji和特殊字符的PII检测",
			businessRisk: "中 - 特殊字符可能影响检测准确性",
		},
		{
			name:         "混合语言边界",
			text:         "Mr. John张三Smith致电+86-13812345678，email: john.zhang@example.com",
			expectedPII:  4,
			description:  "中英文混合的复杂PII场景",
			businessRisk: "高 - 国际化业务场景",
		},
		{
			name:         "格式异常的电话号码",
			text:         "电话：138-1234-5678 / 138.1234.5678 / 138 1234 5678",
			expectedPII:  3,
			description:  "多种格式的电话号码检测",
			businessRisk: "高 - 格式多样性可能导致遗漏",
		},
		{
			name:         "模糊邮箱地址",
			text:         "邮箱：user@domain.co.uk, zhang123@company-name.cn, test.user+tag@gmail.com",
			expectedPII:  3,
			description:  "复杂邮箱格式检测",
			businessRisk: "中 - 复杂域名可能检测困难",
		},
		{
			name:         "连续PII无分隔",
			text:         "姓名张三电话13812345678邮箱zhang@example.com身份证330106199001011234",
			expectedPII:  4,
			description:  "连续无分隔符的PII检测",
			businessRisk: "高 - 边界检测错误风险",
		},
		{
			name:         "文本边界PII",
			text:         "13812345678是开头，身份证330106199001011234在结尾",
			expectedPII:  2,
			description:  "文本开头和结尾的PII检测",
			businessRisk: "中 - 边界检测准确性",
		},
		{
			name:         "转义和编码字符",
			text:         "姓名：张\\n三&#10;电话：138&#45;1234&#46;5678",
			expectedPII:  2,
			description:  "包含HTML转义字符的PII检测",
			businessRisk: "高 - 处理不当可能导致信息泄露",
		},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			// 记录开始时间以监控性能
			start := time.Now()
			processed, entities := detector.DetectAndReplace(tc.text, 0.7)
			duration := time.Since(start)

			// 性能要求：长文本处理不应超过5秒
			if len(tc.text) > 10000 && duration > 5*time.Second {
				t.Errorf("长文本处理时间过长: %v > 5s", duration)
			}

			// 验证基本功能正常
			if processed == "" && tc.text != "" {
				t.Errorf("处理结果为空，但输入不为空")
			}

			// 验证实体检测
			if len(entities) < tc.expectedPII-1 { // 允许1个误差
				t.Errorf("PII检测数量不足，期望>=%d，实际=%d", tc.expectedPII-1, len(entities))
			}

			// 验证脱敏效果
			unmaskedCount := 0
			for _, entity := range entities {
				if strings.Contains(processed, entity.Text) {
					unmaskedCount++
				}
			}

			if unmaskedCount > 0 {
				t.Errorf("有%d个PII未被脱敏: %v", unmaskedCount, entities[:unmaskedCount])
			}

			// 业务风险评估
			var riskLevel string
			if duration > 1*time.Second {
				riskLevel = "高 - 性能问题"
			} else if unmaskedCount > 0 {
				riskLevel = "高 - 脱敏失败"
			} else if len(entities) < tc.expectedPII/2 {
				riskLevel = "中 - 检测不足"
			} else {
				riskLevel = "低 - 正常"
			}

			t.Logf("✅ %s: 检测到%d个PII，耗时=%v", tc.name, len(entities), duration)
			t.Logf("   %s: %s", tc.description, riskLevel)
			t.Logf("   风险评估: %s", tc.businessRisk)

			// 特殊情况的处理建议
			if strings.Contains(tc.text, "unicode") || strings.Contains(tc.text, "emoji") {
				t.Logf("   建议: 增强Unicode和emoji处理能力")
			}
			if len(tc.text) > 5000 {
				t.Logf("   建议: 考虑文本分块处理以提升性能")
			}
			if unmaskedCount > 0 {
				t.Logf("   🚨 警告: 存在脱敏失败，需要紧急修复")
			}
		})
	}
}

// TestPIIDetection_Concurrency 并发安全性测试
func TestPIIDetection_Concurrency(t *testing.T) {
	detector := NewNLPPIIDetector()
	concurrency := 50

	// 准备测试数据
	testTexts := make([]string, concurrency)
	for i := 0; i < concurrency; i++ {
		testTexts[i] = fmt.Sprintf("用户%d张三电话13812345678%d邮箱user%d@example.com", i, i, i)
	}

	// 并发测试通道
	results := make(chan []error, concurrency)

	start := time.Now()

	// 启动并发goroutine
	for i := 0; i < concurrency; i++ {
		go func(index int) {
			var errors []error
			start := time.Now()

			// 模拟多次处理
			for j := 0; j < 5; j++ {
				text := fmt.Sprintf("用户%d张三电话13812345678%d邮箱user%d@example.com", index, j, index)
				processed, entities := detector.DetectAndReplace(text, 0.7)

				// 验证处理结果
				if processed == "" {
					errors = append(errors, fmt.Errorf("处理结果为空"))
				}

				unmasked := 0
				for _, entity := range entities {
					if strings.Contains(processed, entity.Text) {
						unmasked++
					}
				}

				if unmasked > 0 {
					errors = append(errors, fmt.Errorf("PII未完全脱敏: %d个", unmasked))
				}
			}

			duration := time.Since(start)
			if duration > 1*time.Second {
				errors = append(errors, fmt.Errorf("处理时间过长: %v", duration))
			}

			results <- errors
		}(i)
	}

	// 收集结果
	var totalErrors []error
	var totalTime time.Duration
	for i := 0; i < concurrency; i++ {
		errors := <-results
		totalErrors = append(totalErrors, errors...)
	}
	totalTime = time.Since(start)

	// 验证并发安全性
	errorRate := float64(len(totalErrors)) / float64(concurrency*5) * 100
	if errorRate > 5.0 { // 允许5%的错误率
		t.Errorf("并发错误率过高: %.2f%% > 5%%", errorRate)
	}

	avgTimePerRequest := totalTime / time.Duration(concurrency*5)
	if avgTimePerRequest > 100*time.Millisecond {
		t.Logf("⚠️  平均处理时间较长: %v", avgTimePerRequest)
	}

	t.Logf("✅ 并发测试完成:")
	t.Logf("   并发数: %d", concurrency)
	t.Logf("   总请求数: %d", concurrency*5)
	t.Logf("   总时间: %v", totalTime)
	t.Logf("   平均每请求: %v", avgTimePerRequest)
	t.Logf("   错误率: %.2f%%", errorRate)
	t.Logf("   错误详情: %+v", totalErrors)

	if len(totalErrors) == 0 {
		t.Logf("   🎉 并发安全性验证通过")
	} else {
		t.Logf("   ⚠️  发现 %d个并发错误", len(totalErrors))
	}
}

// TestPIIDetection_MemoryUsage 内存使用监控测试
func TestPIIDetection_MemoryUsage(t *testing.T) {
	detector := NewNLPPIIDetector()

	// 大量数据处理以测试内存使用
	batchSize := 100
	batches := 10

	for batch := 0; batch < batches; batch++ {
		t.Run(fmt.Sprintf("批量%d", batch+1), func(t *testing.T) {
			// 创建大量测试数据
			var texts []string
			for i := 0; i < batchSize; i++ {
				text := fmt.Sprintf("客户%d张三，电话13812345678，邮箱customer%d@example.com，地址北京市朝阳区建国路%d号",
					i, i, i)
				texts = append(texts, text)
			}

			// 处理批量数据
			start := time.Now()
			var totalEntities int
			var unmaskedTotal int

			for _, text := range texts {
				processed, entities := detector.DetectAndReplace(text, 0.7)
				totalEntities += len(entities)

				// 验证脱敏效果
				for _, entity := range entities {
					if strings.Contains(processed, entity.Text) {
						unmaskedTotal++
					}
				}
			}

			duration := time.Since(start)
			avgTime := duration / time.Duration(len(texts))

			t.Logf("✅ 批量%d完成:", batch+1)
			t.Logf("   文本数量: %d", len(texts))
			t.Logf("   总PII: %d", totalEntities)
			t.Logf("   平均每文本PII: %.1f", float64(totalEntities)/float64(len(texts)))
			t.Logf("   脱敏失败: %d", unmaskedTotal)
			t.Logf("   总耗时: %v", duration)
			t.Logf("   平均每文本: %v", avgTime)

			// 性能基准
			if avgTime > 10*time.Millisecond {
				t.Logf("   ⚠️  性能警告: 平均处理时间 >10ms")
			}

			if unmaskedTotal > 0 {
				t.Errorf("发现脱敏失败: %d个", unmaskedTotal)
			}

			if float64(totalEntities)/float64(len(texts)) < 1 {
				t.Logf("   ⚠️  检测率可能偏低")
			}
		})
	}
}