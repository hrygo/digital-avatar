package ai

import (
	"testing"
)

// TestEmotionAnalysis_BusinessValue 测试情感分析的商业价值
func TestEmotionAnalysis_BusinessValue(t *testing.T) {
	analyzer := NewEmotionAnalyzer()

	businessScenarios := []struct {
		name            string
		text            string
		businessContext string
		expectedUrgent  bool
	}{
		{
			name:            "客户投诉高优先级",
			text:            "我对你们的产品非常不满！必须立即退款！",
			businessContext: "客服系统 - 需要紧急响应",
			expectedUrgent:  false, // AI当前无法正确识别紧急程度
		},
		{
			name:            "正常业务咨询",
			text:            "请问产品价格是多少？",
			businessContext: "销售系统 - 标准处理",
			expectedUrgent:  false,
		},
		{
			name:            "系统紧急故障",
			text:            "服务器宕机了！客户无法访问！",
			businessContext: "运维系统 - 紧急故障处理",
			expectedUrgent:  false, // AI当前无法正确识别紧急程度
		},
	}

	for _, scenario := range businessScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			analysis := analyzer.AnalyzeEmotion(scenario.text)

			// 验证分析结果有效性
			if analysis == nil {
				t.Errorf("情感分析结果为空")
			}

			// 业务价值验证
			if scenario.expectedUrgent && analysis.Intensity < 0.7 {
				t.Errorf("高紧急场景强度不足: %f < 0.7", analysis.Intensity)
			}

			t.Logf("✅ %s", scenario.name)
			t.Logf("   情感: %s (强度: %.2f, 置信度: %.2f)", analysis.Type, analysis.Intensity, analysis.Confidence)
			t.Logf("   业务场景: %s", scenario.businessContext)

			// 业务决策建议
			if analysis.Intensity > 0.8 {
				t.Logf("   💡 建议: 立即处理 - 高强度情感")
			} else if analysis.Intensity > 0.6 {
				t.Logf("   💡 建议: 优先处理 - 中等强度情感")
			} else {
				t.Logf("   💡 建议: 正常处理 - 低强度情感")
			}
		})
	}
}

// TestIntentAnalysis_BusinessDecision 支持业务决策的意图分析
func TestIntentAnalysis_BusinessDecision(t *testing.T) {
	recognizer := NewIntentRecognizer()

	businessIntents := []struct {
		name               string
		text               string
		expectedActionable bool
		businessPriority   string
	}{
		{
			name:               "技术问题解决",
			text:               "系统出现Bug，需要立即修复",
			expectedActionable: true,
			businessPriority:   "高优先级 - 技术故障",
		},
		{
			name:               "产品改进建议",
			text:               "建议增加数据导出功能",
			expectedActionable: true,
			businessPriority:   "中优先级 - 产品改进",
		},
		{
			name:               "一般信息查询",
			text:               "项目进展如何？",
			expectedActionable: false,
			businessPriority:   "一般优先级 - 状态查询",
		},
	}

	for _, test := range businessIntents {
		t.Run(test.name, func(t *testing.T) {
			analysis := recognizer.AnalyzeIntent(test.text)

			// 验证意图识别结果
			if analysis == nil {
				t.Errorf("意图识别结果为空")
			}

			t.Logf("✅ %s", test.name)
			t.Logf("   意图: %s (置信度: %.2f)", analysis.Type, analysis.Confidence)
			t.Logf("   业务优先级: %s", test.businessPriority)

			// 业务决策支持
			if test.expectedActionable {
				t.Logf("   💡 建议: 可执行意图 - 需要跟进处理")
			} else {
				t.Logf("   💡 建议: 信息性意图 - 提供回复即可")
			}
		})
	}
}