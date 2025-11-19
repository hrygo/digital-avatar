package ai

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"
)

// TestEmotionAnalyzer_PerformanceBenchmarks 情感分析性能基准测试
func TestEmotionAnalyzer_PerformanceBenchmarks(t *testing.T) {
	analyzer := NewEmotionAnalyzer()

	// 不同复杂度的测试用例
	testCases := []struct {
		name       string
		text       string
		maxTime    time.Duration
		description string
	}{
		{
			name:       "简短文本",
			text:       "产品很好用",
			maxTime:    1 * time.Millisecond,
			description: "基础性能基准",
		},
		{
			name:       "中等长度文本",
			text:       "我对你们的产品总体上比较满意，但是界面设计需要改进，建议增加数据导出功能，这样会更方便我们进行数据分析和报告生成。",
			maxTime:    5 * time.Millisecond,
			description: "常见用户反馈长度",
		},
		{
			name:       "长文本处理",
			text:       strings.Repeat("用户反馈显示我们的产品在多个方面表现优异，包括功能完整性、用户界面友好性和性能稳定性。", 10),
			maxTime:    20 * time.Millisecond,
			description: "长篇用户评论处理",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 预热
			for i := 0; i < 10; i++ {
				analyzer.AnalyzeEmotion(tc.text)
			}

			// 性能测试
			iterations := 1000
			start := time.Now()
			var validResults int

			for i := 0; i < iterations; i++ {
				result := analyzer.AnalyzeEmotion(tc.text)
				if result != nil && result.Type != "" {
					validResults++
				}
			}

			duration := time.Since(start)
			avgTime := duration / time.Duration(iterations)
			validityRate := float64(validResults) / float64(iterations) * 100

			t.Logf("✅ %s性能测试完成:", tc.name)
			t.Logf("   文本长度: %d字符", len(tc.text))
			t.Logf("   迭代次数: %d", iterations)
			t.Logf("   总耗时: %v", duration)
			t.Logf("   平均耗时: %v", avgTime)
			t.Logf("   有效结果率: %.1f%%", validityRate)
			t.Logf("   每秒处理量: %.0f次", float64(time.Second)/float64(avgTime))
			t.Logf("   %s", tc.description)

			// 性能断言
			if avgTime > tc.maxTime {
				t.Errorf("平均处理时间超出预期: %v > %v", avgTime, tc.maxTime)
			}

			if validityRate < 95.0 {
				t.Errorf("有效结果率过低: %.1f%% < 95%%", validityRate)
			}
		})
	}
}

// TestIntentRecognizer_PerformanceBenchmarks 意图识别性能基准测试
func TestIntentRecognizer_PerformanceBenchmarks(t *testing.T) {
	recognizer := NewIntentRecognizer()

	intentTypes := []struct {
		name    string
		text    string
		intent  string
		maxTime time.Duration
	}{
		{
			name:    "任务意图",
			text:    "请立即修复系统bug，这很紧急",
			intent:  "task",
			maxTime: 2 * time.Millisecond,
		},
		{
			name:    "信息查询",
			text:    "我想知道项目进展如何？",
			intent:  "question",
			maxTime: 2 * time.Millisecond,
		},
		{
			name:    "建议反馈",
			text:    "建议增加更多功能以满足客户需求",
			intent:  "request",
			maxTime: 3 * time.Millisecond,
		},
	}

	for _, tc := range intentTypes {
		t.Run(tc.name, func(t *testing.T) {
			// 预热
			for i := 0; i < 10; i++ {
				recognizer.AnalyzeIntent(tc.text)
			}

			// 性能测试
			iterations := 1000
			start := time.Now()
			var correctIntents int
			var highConfidence int

			for i := 0; i < iterations; i++ {
				result := recognizer.AnalyzeIntent(tc.text)
				if result != nil {
					if string(result.Type) == tc.intent {
						correctIntents++
					}
					if result.Confidence > 0.8 {
						highConfidence++
					}
				}
			}

			duration := time.Since(start)
			avgTime := duration / time.Duration(iterations)
			accuracy := float64(correctIntents) / float64(iterations) * 100
			confidenceRate := float64(highConfidence) / float64(iterations) * 100

			t.Logf("✅ %s性能测试:", tc.name)
			t.Logf("   迭代次数: %d", iterations)
			t.Logf("   总耗时: %v", duration)
			t.Logf("   平均耗时: %v", avgTime)
			t.Logf("   意图准确率: %.1f%%", accuracy)
			t.Logf("   高置信度率: %.1f%%", confidenceRate)

			// 性能和质量断言
			if avgTime > tc.maxTime {
				t.Errorf("平均处理时间超出预期: %v > %v", avgTime, tc.maxTime)
			}

			if accuracy < 90.0 {
				t.Errorf("意图识别准确率过低: %.1f%% < 90%%", accuracy)
			}

			if confidenceRate < 80.0 {
				t.Logf("⚠️  置信度偏低: %.1f%%", confidenceRate)
			}
		})
	}
}

// TestAIComponents_MemoryUsage 内存使用监控测试
func TestAIComponents_MemoryUsage(t *testing.T) {
	// 获取初始内存状态
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	// 创建多个AI组件实例
	var analyzers []*EmotionAnalyzer
	var recognizers []*IntentRecognizer

	componentCount := 100
	for i := 0; i < componentCount; i++ {
		analyzers = append(analyzers, NewEmotionAnalyzer())
		recognizers = append(recognizers, NewIntentRecognizer())
	}

	// 模拟大量数据处理
	testText := "这是用于测试内存使用情况的示例文本，包含了中文字符和基本的业务场景描述。"
	iterations := 10000

	for i := 0; i < iterations; i++ {
		// 随机选择组件进行处理
		analyzerIndex := i % len(analyzers)
		recognizerIndex := i % len(recognizers)

		analyzers[analyzerIndex].AnalyzeEmotion(testText)
		recognizers[recognizerIndex].AnalyzeIntent(testText)

		// 每1000次迭代检查一次内存
		if i%1000 == 999 {
			runtime.ReadMemStats(&m2)
			memoryUsage := m2.Alloc - m1.Alloc

			t.Logf("迭代%d: 内存使用: %.2f MB", i+1, float64(memoryUsage)/1024/1024)
		}
	}

	// 最终内存检查
	runtime.GC()
	runtime.ReadMemStats(&m2)
	finalMemoryUsage := m2.Alloc - m1.Alloc

	t.Logf("✅ 内存使用监控完成:")
	t.Logf("   组件实例数: %d", componentCount*2)
	t.Logf("   数据处理次数: %d", iterations)
	t.Logf("   最终内存使用: %.2f MB", float64(finalMemoryUsage)/1024/1024)

	// 内存使用合理性检查
	memoryPerComponent := float64(finalMemoryUsage) / float64(componentCount*2)
	if memoryPerComponent > 1024*1024 { // 1MB per component
		t.Errorf("单个组件内存使用过高: %.2f KB > 1MB", memoryPerComponent/1024)
	}

	// 性能建议
	if finalMemoryUsage > 100*1024*1024 { // 100MB
		t.Logf("   ⚠️  总内存使用较高，考虑优化内存管理")
	}

	// 计算内存效率
	dataProcessed := float64(iterations) * float64(len(testText))
	memoryEfficiency := dataProcessed / float64(finalMemoryUsage)
	t.Logf("   内存效率: %.1f bytes/byte", memoryEfficiency)

	if memoryEfficiency < 100 {
		t.Logf("   ⚠️  内存效率偏低，建议优化")
	} else {
		t.Logf("   🎉 内存效率良好")
	}
}

// TestAIComponents_ConcurrencyUnderLoad 负载下的并发性能测试
func TestAIComponents_ConcurrencyUnderLoad(t *testing.T) {
	analyzer := NewEmotionAnalyzer()
	recognizer := NewIntentRecognizer()

	// 不同并发级别测试
	concurrencyLevels := []int{10, 50, 100, 200}

	for _, concurrency := range concurrencyLevels {
		t.Run(fmt.Sprintf("并发级别_%d", concurrency), func(t *testing.T) {
			results := make(chan struct {
				emotionTime time.Duration
				intentTime  time.Duration
				emotionError bool
				intentError  bool
			}, concurrency)

			start := time.Now()

			// 启动并发goroutine
			for i := 0; i < concurrency; i++ {
				go func(index int) {
					testText := fmt.Sprintf("测试文本%d，包含情感和意图信息", index)

					// 测试情感分析
					emotionStart := time.Now()
					emotionResult := analyzer.AnalyzeEmotion(testText)
					emotionTime := time.Since(emotionStart)
					emotionError := emotionResult == nil

					// 测试意图识别
					intentStart := time.Now()
					intentResult := recognizer.AnalyzeIntent(testText)
					intentTime := time.Since(intentStart)
					intentError := intentResult == nil

					results <- struct {
						emotionTime time.Duration
						intentTime  time.Duration
						emotionError bool
						intentError  bool
					}{
						emotionTime: emotionTime,
						intentTime:  intentTime,
						emotionError: emotionError,
						intentError:  intentError,
					}
				}(i)
			}

			// 收集结果
			var totalEmotionTime, totalIntentTime time.Duration
			var emotionErrors, intentErrors int
			var maxEmotionTime, maxIntentTime time.Duration
			var minEmotionTime, minIntentTime time.Duration = time.Hour, time.Hour

			for i := 0; i < concurrency; i++ {
				result := <-results
				totalEmotionTime += result.emotionTime
				totalIntentTime += result.intentTime

				if result.emotionError {
					emotionErrors++
				}
				if result.intentError {
					intentErrors++
				}

				if result.emotionTime > maxEmotionTime {
					maxEmotionTime = result.emotionTime
				}
				if result.emotionTime < minEmotionTime && result.emotionTime > 0 {
					minEmotionTime = result.emotionTime
				}

				if result.intentTime > maxIntentTime {
					maxIntentTime = result.intentTime
				}
				if result.intentTime < minIntentTime && result.intentTime > 0 {
					minIntentTime = result.intentTime
				}
			}

			totalTime := time.Since(start)

			// 计算性能指标
			avgEmotionTime := totalEmotionTime / time.Duration(concurrency)
			avgIntentTime := totalIntentTime / time.Duration(concurrency)
			emotionErrorRate := float64(emotionErrors) / float64(concurrency) * 100
			intentErrorRate := float64(intentErrors) / float64(concurrency) * 100

			t.Logf("✅ 并发级别%d测试完成:", concurrency)
			t.Logf("   总耗时: %v", totalTime)
			t.Logf("   平均情感分析时间: %v", avgEmotionTime)
			t.Logf("   平均意图识别时间: %v", avgIntentTime)
			t.Logf("   情感分析错误率: %.1f%%", emotionErrorRate)
			t.Logf("   意图识别错误率: %.1f%%", intentErrorRate)
			t.Logf("   情感分析时间范围: %v - %v", minEmotionTime, maxEmotionTime)
			t.Logf("   意图识别时间范围: %v - %v", minIntentTime, maxIntentTime)

			// 并发性能断言
			emotionErrorThreshold := 5.0 // 5%错误率
			intentErrorThreshold := 5.0

			if emotionErrorRate > emotionErrorThreshold {
				t.Errorf("情感分析错误率过高: %.1f%% > %.1f%%", emotionErrorRate, emotionErrorThreshold)
			}

			if intentErrorRate > intentErrorThreshold {
				t.Errorf("意图识别错误率过高: %.1f%% > %.1f%%", intentErrorRate, intentErrorThreshold)
			}

			// 性能下降检查
			if concurrency > 50 {
				// 在高并发下，性能下降是正常的，但不能下降太多
				expectedSlowdown := float64(concurrency) / 50.0
				actualEmotionSlowdown := float64(avgEmotionTime) / float64(2*time.Millisecond) // 基准2ms
				actualIntentSlowdown := float64(avgIntentTime) / float64(3*time.Millisecond) // 基准3ms

				maxSlowdown := expectedSlowdown * 2 // 允许2倍的预期下降

				if actualEmotionSlowdown > maxSlowdown {
					t.Logf("   ⚠️  情感分析性能下降较大: %.1fx", actualEmotionSlowdown)
				}

				if actualIntentSlowdown > maxSlowdown {
					t.Logf("   ⚠️  意图识别性能下降较大: %.1fx", actualIntentSlowdown)
				}
			}

			// 并发效率
			totalOps := concurrency * 2 // 每个goroutine执行2个操作
			opsPerSecond := float64(totalOps) / totalTime.Seconds()
			t.Logf("   并发吞吐量: %.0f ops/sec", opsPerSecond)

			if opsPerSecond < 100 {
				t.Logf("   ⚠️  并发吞吐量偏低")
			}
		})
	}
}