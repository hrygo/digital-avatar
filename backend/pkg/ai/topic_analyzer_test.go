package ai

import (
	"testing"
)

func TestTopicAnalyzer_AnalyzeTopics(t *testing.T) {
	analyzer := NewTopicAnalyzer(nil)

	// Use space-separated words to simulate segmentation
	// Ensure key terms appear >= 2 times to pass frequency filter
	documents := []string{
		"机器学习 是 人工智能 的 分支",
		"机器学习 使用 神经网络",
		"Python 是 数据科学 的 语言",
		"Python 适合 机器学习 开发",
		"人工智能 正在 改变 世界",
		"神经网络 模拟 人脑",
		"数据科学 需要 统计学",
	}

	topics, err := analyzer.AnalyzeTopics(documents)
	if err != nil {
		t.Fatalf("AnalyzeTopics() error = %v", err)
	}

	if len(topics) == 0 {
		t.Error("AnalyzeTopics() returned 0 topics")
	}

	// Check if we have valid topics
	for _, topic := range topics {
		if len(topic.Keywords) == 0 {
			t.Errorf("Topic %s has no keywords", topic.ID)
		}
	}
}

func TestTopicAnalyzer_OptimizeTopicNumber(t *testing.T) {
	analyzer := NewTopicAnalyzer(nil)

	documents := []string{
		"苹果 是 水果",
		"香蕉 是 水果",
		"橘子 是 水果",
		"猫 是 动物",
		"狗 是 动物",
		"鸟 是 动物",
		"汽车 是 交通工具",
		"飞机 是 交通工具",
	}

	// We expect roughly 3 topics (Fruits, Animals, Vehicles)
	bestK, scores, err := analyzer.OptimizeTopicNumber(documents, 5)
	if err != nil {
		t.Fatalf("OptimizeTopicNumber() error = %v", err)
	}

	if bestK < 2 || bestK > 5 {
		t.Errorf("OptimizeTopicNumber() returned unreasonable K: %d", bestK)
	}

	if len(scores) == 0 {
		t.Error("OptimizeTopicNumber() returned no scores")
	}
}
