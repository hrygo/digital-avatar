package crypto

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"
)

// BenchmarkPIIDetection PII检测性能基准测试
func BenchmarkPIIDetection(b *testing.B) {
	detector := NewNLPPIIDetector()

	// 不同规模的测试数据
	testCases := []struct {
		name  string
		text  string
		scale string
	}{
		{
			name:  "Small_Text",
			text:  "用户张三的电话是13812345678，邮箱zhang@company.com",
			scale: "Small (<100 chars)",
		},
		{
			name:  "Medium_Text",
			text: strings.Repeat("客户李四，电话13998765432，邮箱li@example.com。", 10),
			scale: "Medium (100-500 chars)",
		},
		{
			name:  "Large_Text",
			text: strings.Repeat("员工王五，身份证330106199001011234，电话13555666777，邮箱wang@enterprise.com。", 50),
			scale: "Large (500-2000 chars)",
		},
		{
			name:  "XLarge_Text",
			text: strings.Repeat("重要客户赵六，公司地址北京市朝阳区建国路88号，银行账户6222021234567890123。", 100),
			scale: "XLarge (>2000 chars)",
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_, entities := detector.DetectAndReplace(tc.text, 0.7)
				if len(entities) == 0 {
					b.Errorf("Expected to detect PII in %s", tc.scale)
				}
			}

			b.Logf("Scale: %s, Text Length: %d chars", tc.scale, len(tc.text))
		})
	}
}

// BenchmarkConcurrentPIIDetection 并发PII检测性能测试
func BenchmarkConcurrentPIIDetection(b *testing.B) {
	detector := NewNLPPIIDetector()
	testText := "用户张三的电话是13812345678，邮箱zhang@company.com，身份证330106199001011234"

	// 不同并发级别测试
	concurrencyLevels := []int{1, 5, 10, 20, 50, 100}

	for _, concurrency := range concurrencyLevels {
		b.Run(fmt.Sprintf("Concurrency_%d", concurrency), func(b *testing.B) {
			b.ResetTimer()
			b.SetParallelism(concurrency)

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, entities := detector.DetectAndReplace(testText, 0.7)
					if len(entities) == 0 {
						b.Error("Expected to detect PII")
					}
				}
			})

			b.Logf("Concurrency: %d, Completed: %d operations", concurrency, b.N)
		})
	}
}

// BenchmarkMemoryUsage 内存使用基准测试
func BenchmarkMemoryUsage(b *testing.B) {
	detector := NewNLPPIIDetector()

	// 获取初始内存状态
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// 模拟不同规模的文本处理
		text := fmt.Sprintf("用户%d张三，电话13812345678，邮箱user%d@company.com", i%100, i%100)
		_, entities := detector.DetectAndReplace(text, 0.7)
		if len(entities) == 0 {
			b.Error("Expected to detect PII")
		}
	}

	b.StopTimer()

	// 检查最终内存使用
	runtime.GC()
	runtime.ReadMemStats(&m2)

	memoryUsed := m2.Alloc - m1.Alloc
	memoryPerOp := float64(memoryUsed) / float64(b.N)

	b.ReportMetric(memoryPerOp, "bytes/op")
	b.Logf("Memory used: %.2f MB, Per operation: %.2f bytes",
		float64(memoryUsed)/1024/1024, memoryPerOp)
}

// BenchmarkReplacementMapPerformance 替换映射性能测试
func BenchmarkReplacementMapPerformance(b *testing.B) {
	detector := NewNLPPIIDetector()

	// 预先填充替换映射以测试性能影响
	testTexts := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		testTexts[i] = fmt.Sprintf("用户%d张三，电话138%08d，邮箱user%d@company.com",
			i, i, i)
	}

	// 预热阶段
	for _, text := range testTexts[:100] {
		detector.DetectAndReplace(text, 0.7)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		text := testTexts[i%len(testTexts)]
		_, entities := detector.DetectAndReplace(text, 0.7)
		if len(entities) == 0 {
			b.Error("Expected to detect PII")
		}
	}

	// 报告替换映射大小
	replacementMap := detector.GetReplacementMap()
	b.Logf("Replacement map size: %d entries", len(replacementMap))
}

// TestPerformanceRegression 性能回归测试
func TestPerformanceRegression(t *testing.T) {
	detector := NewNLPPIIDetector()

	// 性能基准（基于v0.2.0的实测数据）
	performanceBenchmarks := map[string]time.Duration{
		"small_text_processing":  1 * time.Millisecond,
		"medium_text_processing": 5 * time.Millisecond,
		"large_text_processing":  20 * time.Millisecond,
	}

	testCases := []struct {
		name           string
		text           string
		benchmarkKey   string
		maxAcceptable   time.Duration
	}{
		{
			name:         "Small_Text_Performance",
			text:         "用户张三的电话是13812345678",
			benchmarkKey: "small_text_processing",
			maxAcceptable: performanceBenchmarks["small_text_processing"] * 2, // 允许2倍基准时间
		},
		{
			name:         "Medium_Text_Performance",
			text:         strings.Repeat("客户李四，电话13998765432。", 10),
			benchmarkKey: "medium_text_processing",
			maxAcceptable: performanceBenchmarks["medium_text_processing"] * 2,
		},
		{
			name:         "Large_Text_Performance",
			text:         strings.Repeat("员工王五，身份证330106199001011234。", 50),
			benchmarkKey: "large_text_processing",
			maxAcceptable: performanceBenchmarks["large_text_processing"] * 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			iterations := 100
			start := time.Now()

			for i := 0; i < iterations; i++ {
				_, entities := detector.DetectAndReplace(tc.text, 0.7)
				if len(entities) == 0 {
					t.Errorf("Expected to detect PII in %s", tc.name)
				}
			}

			avgDuration := time.Since(start) / time.Duration(iterations)

			t.Logf("Average processing time: %v", avgDuration)
			t.Logf("Benchmark: %v", performanceBenchmarks[tc.benchmarkKey])
			t.Logf("Max acceptable: %v", tc.maxAcceptable)

			if avgDuration > tc.maxAcceptable {
				t.Errorf("Performance regression detected: %v > %v",
					avgDuration, tc.maxAcceptable)
			} else {
				t.Logf("✅ Performance within acceptable range")
			}

			// 计算性能提升倍数
			improvement := float64(performanceBenchmarks[tc.benchmarkKey]) /
				float64(avgDuration)
			t.Logf("Performance improvement: %.2fx", improvement)
		})
	}
}

// TestAccuracyVsPerformance 准确性与性能平衡测试
func TestAccuracyVsPerformance(t *testing.T) {
	detector := NewNLPPIIDetector()

	// 包含明确PII的测试文本
	testTexts := []string{
		"张经理的电话是13812345678",
		"李总的邮箱是li@company.com",
		"王女士的身份证是330106199001011234",
		"地址：北京市朝阳区建国路88号",
		"联系人：赵四，电话：13987654321",
	}

	thresholds := []float64{0.5, 0.6, 0.7, 0.8, 0.9}

	for _, threshold := range thresholds {
		t.Run(fmt.Sprintf("Threshold_%.1f", threshold), func(t *testing.T) {
			var totalDetectionTime time.Duration
			var totalDetections int
			var correctDetections int

			for _, text := range testTexts {
				start := time.Now()
				processed, entities := detector.DetectAndReplace(text, threshold)
				duration := time.Since(start)

				totalDetectionTime += duration
				totalDetections += len(entities)

				// 简单的准确性检查：如果处理后的文本包含[PHONE_REDACTED]等，认为检测正确
				if strings.Contains(processed, "[") && strings.Contains(processed, "]") {
					correctDetections++
				}
			}

			avgTime := totalDetectionTime / time.Duration(len(testTexts))
			accuracy := float64(correctDetections) / float64(len(testTexts)) * 100

			t.Logf("Threshold: %.1f", threshold)
			t.Logf("Average processing time: %v", avgTime)
			t.Logf("Accuracy: %.1f%%", accuracy)
			t.Logf("Total detections: %d", totalDetections)

			// 平衡性检查：既要保证性能，也要保证准确性
			if avgTime > 5*time.Millisecond {
				t.Logf("⚠️  Processing time is high: %v", avgTime)
			}

			if accuracy < 80 {
				t.Logf("⚠️  Accuracy is low: %.1f%%", accuracy)
			}

			// 计算综合评分
			performanceScore := 100.0 / (avgTime.Seconds() * 1000) // 性能分数（时间越短分数越高）
			compositeScore := (accuracy + performanceScore) / 2

			t.Logf("Performance score: %.1f", performanceScore)
			t.Logf("Composite score: %.1f", compositeScore)
		})
	}
}