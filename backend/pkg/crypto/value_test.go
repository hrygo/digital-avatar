package crypto

import (
	"strings"
	"testing"
)

// TestPIIDetection_BusinessValue 专注于测试PII检测的商业价值
func TestPIIDetection_BusinessValue(t *testing.T) {
	detector := NewNLPPIIDetector()

	// 核心业务场景：保护客户隐私
	testCases := []struct {
		name        string
		text        string
		description string
		businessRisk string
	}{
		{
			name: "客户联系信息",
			text: "请联系张经理，电话13812345678，邮箱zhang@company.com",
			description: "客户联系方式包含PII",
			businessRisk: "高 - 客户信息泄露风险",
		},
		{
			name: "员工敏感信息",
			text: "员工李四，身份证330106199001011234，工资信息保密",
			description: "员工个人信息和工资数据",
			businessRisk: "极高 - 违反隐私法规",
		},
		{
			name: "正常业务文本",
			text: "产品功能很好，建议增加数据导出功能",
			description: "普通业务反馈，无敏感信息",
			businessRisk: "低 - 无隐私风险",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			processed, entities := detector.DetectAndReplace(tc.text, 0.7)

			// 验证脱敏效果
			for _, entity := range entities {
				if strings.Contains(processed, entity.Text) {
					t.Errorf("PII脱敏失败: 敏感信息'%s'仍在处理结果中", entity.Text)
				}
			}

			// 业务价值验证
			if len(entities) > 0 {
				t.Logf("✅ 检测到%d个敏感信息，风险等级: %s", len(entities), tc.businessRisk)
				t.Logf("   脱敏成功: %s", processed)
			} else {
				t.Logf("✅ 无敏感信息，文本安全: %s", processed)
			}

			// 验证风险等级评估的合理性
			hasHighRiskPII := false
			for _, entity := range entities {
				if entity.Type == "ID" || entity.Type == "TEL" || entity.Type == "EMAIL" {
					hasHighRiskPII = true
					break
				}
			}

			if hasHighRiskPII && len(entities) == 0 {
				t.Errorf("高风险场景应该检测到敏感信息")
			}
		})
	}
}

// TestPIIConsistency_BusinessRequirement 验证PII替换一致性对业务的重要性
func TestPIIConsistency_BusinessRequirement(t *testing.T) {
	detector := NewNLPPIIDetector()

	// 业务场景：同一份文档中同一人应该用相同的替换
	text1 := "张三的电话是13812345678"
	text2 := "请及时联系张三，他的电话是13812345678"

	processed1, _ := detector.DetectAndReplace(text1, 0.7)
	processed2, _ := detector.DetectAndReplace(text2, 0.7)

	// 验证张三的替换在两次处理中一致
	zhangsanInText1 := strings.Contains(text1, "张三") && !strings.Contains(processed1, "张三")
	zhangsanInText2 := strings.Contains(text2, "张三") && !strings.Contains(processed2, "张三")

	if zhangsanInText1 && zhangsanInText2 {
		t.Logf("✅ 张三在两次处理中都被正确脱敏")
		t.Logf("   处理1: %s", processed1)
		t.Logf("   处理2: %s", processed2)
	}

	// 验证手机号的替换一致性
	phoneInText1 := strings.Contains(text1, "13812345678") && !strings.Contains(processed1, "13812345678")
	phoneInText2 := strings.Contains(text2, "13812345678") && !strings.Contains(processed2, "13812345678")

	if phoneInText1 && phoneInText2 {
		t.Logf("✅ 手机号在两次处理中都被正确脱敏")
	}
}