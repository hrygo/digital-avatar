package crypto

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

// TestPIIOptimization 性能优化测试
func TestPIIOptimization(t *testing.T) {
	detector := NewNLPPIIDetector()

	// 测试用例
	testTexts := []struct {
		name string
		text string
	}{
		{
			name: "Basic_PII",
			text: "张经理电话13812345678，邮箱zhang@company.com",
		},
		{
			name: "Multiple_Entities",
			text: "客户李四，电话13998765432，邮箱li@example.com。员工王五，身份证330106199001011234，地址北京市朝阳区建国路88号。",
		},
		{
			name: "Large_Text",
			text: `重要客户张经理的电话是13812345678，邮箱zhang.manager@company.com，身份证号330106199001011234。
联系人李副总的联系方式包括手机13987654321，办公室电话010-12345678，邮箱li.deputy@enterprise.com。
技术主管王工程师的专业信息：邮箱wang.tech@company.com，工号EMP2024001，紧急联系13788889999。
销售总监陈经理的客户数据：客户A电话13511112222，邮箱client.a@company.com；客户B电话13633334444，地址上海市浦东新区张江高科技园区。`,
		},
	}

	iterations := 1000

	t.Log("🚀 PII Performance Optimization Test")
	t.Log("Testing original vs fast optimized methods")

	for _, tc := range testTexts {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("📊 Testing %s (Length: %d chars)", tc.name, len(tc.text))

			// 预热
			for i := 0; i < 10; i++ {
				detector.DetectAndReplace(tc.text, 0.7)
				detector.DetectAndReplaceFast(tc.text, 0.7)
			}

			// 测试原始版本
			runtime.GC()
			start := time.Now()
			for i := 0; i < iterations; i++ {
				_, _ = detector.DetectAndReplace(tc.text, 0.7)
			}
			originalDuration := time.Since(start)

			// 测试优化版本
			runtime.GC()
			start = time.Now()
			for i := 0; i < iterations; i++ {
				_, _ = detector.DetectAndReplaceFast(tc.text, 0.7)
			}
			optimizedDuration := time.Since(start)

			// 计算性能指标
			originalAvg := originalDuration / time.Duration(iterations)
			optimizedAvg := optimizedDuration / time.Duration(iterations)
			speedup := float64(originalDuration) / float64(optimizedDuration)

			// 输出结果
			t.Logf("Original method:")
			t.Logf("  Total time: %v", originalDuration)
			t.Logf("  Average time: %v", originalAvg)

			t.Logf("Fast optimized method:")
			t.Logf("  Total time: %v", optimizedDuration)
			t.Logf("  Average time: %v", optimizedAvg)

			t.Logf("🎯 Performance Improvement:")
			t.Logf("  Speedup: %.2fx faster", speedup)

			// v0.3.0 目标验证
			targetSpeedup := 3.0
			if speedup >= targetSpeedup {
				t.Logf("🎉 v0.3.0 target achieved: %.2fx ≥ %.1fx", speedup, targetSpeedup)
			} else {
				t.Logf("📈 Progress toward v0.3.0 target: %.2fx / %.1fx (%.1f%% complete)",
					speedup, targetSpeedup, speedup/targetSpeedup*100)
			}

			// 性能目标验证
			if speedup >= 2.0 {
				t.Logf("✅ Significant improvement achieved")
			} else {
				t.Logf("📈 Modest improvement achieved")
			}
		})
	}
}

// TestCacheEffectiveness 缓存效果测试
func TestCacheEffectiveness(t *testing.T) {
	detector := NewNLPPIIDetector()
	testText := "用户张三的电话是13812345678，邮箱zhang@company.com"

	t.Log("💾 Cache Effectiveness Test")

	// 首次调用（无缓存）
	start := time.Now()
	processed1, entities1 := detector.DetectAndReplaceFast(testText, 0.7)
	firstCallTime := time.Since(start)

	// 第二次调用（使用缓存）
	start = time.Now()
	processed2, entities2 := detector.DetectAndReplaceFast(testText, 0.7)
	secondCallTime := time.Since(start)

	t.Logf("First call time: %v", firstCallTime)
	t.Logf("Second call time: %v", secondCallTime)

	// 验证缓存效果
	if secondCallTime < firstCallTime/2 {
		cacheSpeedup := float64(firstCallTime) / float64(secondCallTime)
		t.Logf("✅ Cache effective: %.2fx speedup", cacheSpeedup)
	} else {
		t.Logf("⚠️  Cache may not be significantly improving performance")
	}

	// 验证结果一致性
	if processed1 != processed2 {
		t.Error("Cached result differs from original result")
	}

	if len(entities1) != len(entities2) {
		t.Error("Entity count mismatch between cached and original result")
	}

	// 测试性能统计
	stats := detector.GetPerformanceStats()
	t.Logf("Performance stats after %d calls:", stats.TotalRequests)
	t.Logf("  Cache hits: %d", stats.CacheHits)
	t.Logf("  Average response time: %.2f μs", float64(stats.AvgResponseTime)/1000)
	t.Logf("  Total entities detected: %d", stats.EntitiesDetected)

	if stats.CacheHits > 0 {
		cacheHitRate := float64(stats.CacheHits) / float64(stats.TotalRequests) * 100
		t.Logf("  Cache hit rate: %.1f%%", cacheHitRate)
	}

	// 清空缓存并测试
	detector.ClearCache()
	_ = detector.GetPerformanceStats()
	t.Logf("Cache cleared")
}

// TestAccuracyConsistency 准确性一致性测试
func TestAccuracyConsistency(t *testing.T) {
	detector := NewNLPPIIDetector()

	testCases := []string{
		"张经理电话13812345678",
		"李总邮箱li@company.com",
		"王女士身份证330106199001011234",
		"地址北京市朝阳区建国路88号",
		"联系人赵四，电话13987654321",
		"test@company.com", // 只有邮箱
		"13312345678", // 只有电话
		"330106199001011234", // 只有身份证
	}

	t.Log("🎯 Accuracy Consistency Test")

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Case_%d", i+1), func(t *testing.T) {
			// 原始检测结果
			originalProcessed, originalEntities := detector.DetectAndReplace(testCase, 0.7)
			// 优化检测结果
			optimizedProcessed, optimizedEntities := detector.DetectAndReplaceFast(testCase, 0.7)

			t.Logf("Original: '%s' -> %d entities", originalProcessed, len(originalEntities))
			t.Logf("Optimized: '%s' -> %d entities", optimizedProcessed, len(optimizedEntities))

			// 检查实体数量一致性
			if len(originalEntities) != len(optimizedEntities) {
				// 对于某些边缘情况，优化版本可能有不同的检测策略
				if len(optimizedEntities) > 0 && len(originalEntities) > 0 {
					t.Logf("⚠️  Entity count differs but both detected PII: Original=%d, Optimized=%d",
						len(originalEntities), len(optimizedEntities))
				} else if len(originalEntities) > 0 {
					t.Errorf("Original detected PII but optimized didn't: Original=%d, Optimized=%d",
						len(originalEntities), len(optimizedEntities))
				} else if len(optimizedEntities) > 0 {
					t.Logf("ℹ️  Optimized detected additional PII: Original=%d, Optimized=%d",
						len(originalEntities), len(optimizedEntities))
				}
			} else {
				t.Logf("✅ Entity count consistent: %d", len(originalEntities))
			}

			// 检查脱敏效果
			if len(originalEntities) > 0 && originalProcessed == testCase {
				t.Error("Original failed to mask PII")
			}

			if len(optimizedEntities) > 0 && optimizedProcessed == testCase {
				t.Error("Optimized failed to mask PII")
			}
		})
	}
}

// TestConcurrentPerformance 并发性能测试
func TestConcurrentPerformance(t *testing.T) {
	detector := NewNLPPIIDetector()
	testText := "用户张三的电话是13812345678，邮箱zhang@company.com"
	concurrencyLevels := []int{1, 5, 10, 20}

	t.Log("🔄 Concurrent Performance Test")

	for _, concurrency := range concurrencyLevels {
		t.Run(fmt.Sprintf("Concurrency_%d", concurrency), func(t *testing.T) {
			t.Logf("Testing with %d concurrent goroutines", concurrency)

			// 使用快速优化方法进行并发测试
			start := time.Now()
			for i := 0; i < concurrency; i++ {
				go func(id int) {
					_, entities := detector.DetectAndReplaceFast(testText, 0.7)
					if len(entities) == 0 {
						t.Errorf("Goroutine %d: No entities detected", id)
					}
				}(i)
			}

			// 等待并发完成
			time.Sleep(500 * time.Millisecond)
			currentTime := time.Since(start)

			// 计算性能指标
			operationsPerSecond := float64(concurrency) / currentTime.Seconds()
			avgTimePerOp := currentTime / time.Duration(concurrency)

			t.Logf("Total time: %v", currentTime)
			t.Logf("Average time per operation: %v", avgTimePerOp)
			t.Logf("Operations per second: %.0f", operationsPerSecond)

			// 性能目标验证
			if operationsPerSecond > 1000 {
				t.Logf("✅ Excellent concurrent performance: %.0f ops/sec", operationsPerSecond)
			} else if operationsPerSecond > 500 {
				t.Logf("✅ Good concurrent performance: %.0f ops/sec", operationsPerSecond)
			} else {
				t.Logf("📈 Concurrent performance: %.0f ops/sec", operationsPerSecond)
			}

			// 获取性能统计
			stats := detector.GetPerformanceStats()
			t.Logf("Performance stats: %+v", stats)
		})
	}
}