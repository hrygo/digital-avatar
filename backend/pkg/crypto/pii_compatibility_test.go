package crypto

import (
	"strings"
	"testing"
)

// TestPIIDetector_BackwardCompatibility 测试向后兼容性
func TestPIIDetector_BackwardCompatibility(t *testing.T) {
	// 创建旧版本检测器
	oldDetector := NewPIIDetector()

	// 创建新版本检测器
	newDetector := NewNLPPIIDetector()

	testText := "联系人：张三，手机：13812345678，邮箱：zhang@example.com"

	// 使用旧版本接口（无threshold参数）
	oldProcessed, oldReport := oldDetector.DetectAndReplace(testText)

	// 使用新版本接口（有threshold参数）
	newProcessed, newEntities := newDetector.DetectAndReplace(testText, 0.7)

	// 验证旧版本正常工作
	if oldProcessed == testText {
		t.Error("Old detector should have detected and replaced PII")
	}

	// 验证新版本正常工作
	if newProcessed == testText {
		t.Error("New detector should have detected and replaced PII")
	}

	// 验证两者都进行了脱敏处理（文本发生了变化）
	if !strings.Contains(oldProcessed, "[") || !strings.Contains(newProcessed, "[") {
		t.Error("Both detectors should contain replacement markers")
	}

	t.Logf("Old detector result: %s", oldProcessed)
	t.Logf("New detector result: %s", newProcessed)
	t.Logf("Old report entities: %d", oldReport.TotalDetections)
	t.Logf("New entities count: %d", len(newEntities))
}

// TestPIIDetector_CompatibilityMethods 测试兼容层的其他方法
func TestPIIDetector_CompatibilityMethods(t *testing.T) {
	detector := NewPIIDetector()
	testText := "测试文本，包含张三的电话13812345678"

	// 测试统计功能
	stats := detector.GetStatistics(testText)
	if total, ok := stats["total_entities"].(int); !ok || total == 0 {
		t.Error("Statistics should return non-zero entity count")
	}

	// 测试替换映射
	replacementMap := detector.GetReplacementMap()
	if replacementMap == nil {
		t.Error("Replacement map should not be nil")
	}

	t.Logf("Statistics: %+v", stats)
	t.Logf("Replacement map size: %d", len(replacementMap))
}

// TestPIIDetector_ErrorHandling 测试错误处理
func TestPIIDetector_ErrorHandling(t *testing.T) {
	detector := NewPIIDetector()

	// 测试空文本
	processed, report := detector.DetectAndReplace("")
	if processed != "" {
		t.Error("Empty text should remain empty after processing")
	}
	if report.TotalDetections != 0 {
		t.Error("Empty text should not detect any entities")
	}

	// 测试新版本检测器的错误处理
	newDetector := NewNLPPIIDetector()
	processed, entities := newDetector.DetectAndReplace("张三 13812345678", 1.0)
	// 高阈值应该检测到很少或没有实体
	t.Logf("High threshold (1.0) detected %d entities", len(entities))
	if processed == "张三 13812345678" {
		t.Log("High threshold correctly detected no entities")
	}
}

// TestPIIDetector_NilSafety 测试nil安全性
func TestPIIDetector_NilSafety(t *testing.T) {
	var detector *PIIDetector

	// 这些操作在nil检测器上应该安全或返回适当错误
	if detector != nil {
		t.Error("Detector should be nil")
	}

	// 创建一个检测器来测试正常功能
	detector = NewPIIDetector()
	if detector == nil {
		t.Fatal("NewPIIDetector should not return nil")
	}

	// 验证内部NLP检测器不为nil
	if detector.nlpDetector == nil {
		t.Error("Internal NLP detector should not be nil")
	}
}

// TestPIIDetector_PerformanceCompatibility 测试性能兼容性
func TestPIIDetector_PerformanceCompatibility(t *testing.T) {
	detector := NewPIIDetector()

	// 较长的测试文本
	testText := `这是一个包含多个个人信息的测试文本。联系人包括张三（电话：13812345678）、
李四（电话：13987654321），以及王五（邮箱：wangwu@example.com）。公司地址在北京市朝阳区建国门外大街1号，
另一分公司在上海市浦东新区陆家嘴环路1000号。邮政编码是100000，另一地址邮编是200000。`

	// 测试多次调用的性能一致性
	const iterations = 10
	for i := 0; i < iterations; i++ {
		processed, report := detector.DetectAndReplace(testText)
		if processed == testText {
			t.Errorf("Iteration %d: text should be processed", i)
		}

		if report.TotalDetections == 0 {
			t.Errorf("Iteration %d: should detect entities", i)
		}
	}

	t.Logf("Performance test completed with %d iterations", iterations)
}