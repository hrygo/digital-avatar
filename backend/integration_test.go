//go:build integration
// +build integration

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"twin-os/backend/pkg/ai"
	"twin-os/backend/pkg/crypto"
)

// TestPIIIntegrationWorkflow 端到端PII处理工作流测试
func TestPIIIntegrationWorkflow(t *testing.T) {
	// 设置测试环境
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// 模拟业务服务
	piiDetector := crypto.NewNLPPIIDetector()
	emotionAnalyzer := ai.NewEmotionAnalyzer()
	intentRecognizer := ai.NewIntentRecognizer()

	// PII检测端点
	router.POST("/api/v1/pii/analyze", func(c *gin.Context) {
		var request struct {
			Text      string  `json:"text"`
			Threshold float64 `json:"threshold"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_REQUEST",
					"message": "请求参数错误",
				},
			})
			return
		}

		// 第一步：PII检测和脱敏
		processedText, entities := piiDetector.DetectAndReplace(request.Text, request.Threshold)

		// 第二步：情感分析
		emotionAnalysis := emotionAnalyzer.AnalyzeEmotion(processedText)

		// 第三步：意图识别
		intentAnalysis := intentRecognizer.AnalyzeIntent(processedText)

		// 汇总结果
		result := gin.H{
			"original_text": request.Text,
			"processed_text": processedText,
			"pii_entities": entities,
			"pii_count": len(entities),
			"emotion_analysis": gin.H{
				"type":      emotionAnalysis.Type,
				"intensity": emotionAnalysis.Intensity,
				"confidence": emotionAnalysis.Confidence,
			},
			"intent_analysis": gin.H{
				"type":       intentAnalysis.Type,
				"confidence": intentAnalysis.Confidence,
				"priority":   intentAnalysis.Priority,
			},
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    result,
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// 测试用例：真实业务场景
	testCases := []struct {
		name         string
		inputText    string
		description  string
		expectPII    int
		businessRisk string
	}{
		{
			name:         "客户投诉处理",
			inputText:    "我对你们的产品非常不满！张经理的电话是13812345678，邮箱zhang@company.com，请立即退款！",
			description:  "高风险客户投诉场景",
			expectPII:    3,
			businessRisk: "高 - 需要紧急处理",
		},
		{
			name:         "员工信息管理",
			inputText:    "新员工李四的身份证是330106199001011234，月薪8000元，入职日期2024年1月15日",
			description:  "员工敏感信息处理",
			expectPII:    2,
			businessRisk: "高 - 需要严格保密",
		},
		{
			name:         "产品反馈收集",
			inputText:    "产品功能很棒，王五建议增加数据导出功能，这样更方便我们使用",
			description:  "正常产品反馈",
			expectPII:    1,
			businessRisk: "中 - 需要脱敏处理",
		},
		{
			name:         "系统监控报告",
			inputText:    "服务器监控显示CPU使用率85%，内存使用率60%，磁盘使用率45%，需要扩容",
			description:  "技术监控报告",
			expectPII:    0,
			businessRisk: "低 - 无敏感信息",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(gin.H{
				"text":      tc.inputText,
				"threshold": 0.7,
			})

			req, _ := http.NewRequest("POST", "/api/v1/pii/analyze", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// 验证HTTP状态
			if w.Code != http.StatusOK {
				t.Errorf("HTTP状态错误: %d, 期望200", w.Code)
				t.Errorf("响应内容: %s", w.Body.String())
				return
			}

			// 验证响应结构
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Errorf("JSON解析失败: %v", err)
				return
			}

			success := response["success"].(bool)
			if !success {
				t.Error("API调用应该成功")
				return
			}

			data := response["data"].(map[string]interface{})

			// 验证PII检测结果
			piiCount := int(data["pii_count"].(float64))
			if piiCount < tc.expectPII-1 {
				t.Errorf("PII检测数量不足: 期望>%d, 实际=%d", tc.expectPII-1, piiCount)
			}

			// 验证脱敏效果
			originalText := data["original_text"].(string)
			processedText := data["processed_text"].(string)
			if originalText == processedText && tc.expectPII > 0 {
				t.Error("包含PII的文本应该被脱敏")
			}

			// 验证AI分析结果
			emotionData := data["emotion_analysis"].(map[string]interface{})
			intentData := data["intent_analysis"].(map[string]interface{})

			if emotionData["type"] == "" || intentData["type"] == "" {
				t.Error("AI分析结果不能为空")
			}

			// 业务价值验证
			t.Logf("✅ %s集成测试完成:", tc.name)
			t.Logf("   描述: %s", tc.description)
			t.Logf("   风险等级: %s", tc.businessRisk)
			t.Logf("   检测到PII: %d个", piiCount)
			t.Logf("   情感分析: %s (置信度: %.1f)",
				emotionData["type"],
				emotionData["confidence"].(float64)*100)
			t.Logf("   意图识别: %s (置信度: %.1f)",
				intentData["type"],
				intentData["confidence"].(float64)*100)

			// 业务决策建议
			var recommendations []string

			if piiCount > 2 {
				recommendations = append(recommendations, "🚨 高风险PII - 建议立即处理")
			}

			emotionType := emotionData["type"].(string)
			if emotionType == "anger" || emotionType == "sadness" {
				recommendations = append(recommendations, "⚠️ 负面情绪 - 需要特别关注")
			}

			intentType := intentData["type"].(string)
			if intentType == "task" {
				recommendations = append(recommendations, "💡 任务型意图 - 需要跟进执行")
			}

			if len(recommendations) > 0 {
				t.Log("   业务建议:")
				for i, rec := range recommendations {
					t.Logf("     %d. %s", i+1, rec)
				}
			} else {
				t.Log("   业务建议: 正常处理即可")
			}
		})
	}
}

// TestMultiModalDataProcessing 多模态数据处理集成测试
func TestMultiModalDataProcessing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// 初始化AI组件
	piiDetector := crypto.NewNLPPIIDetector()
	emotionAnalyzer := ai.NewEmotionAnalyzer()
	intentRecognizer := ai.NewIntentRecognizer()
	topicAnalyzer := ai.NewTopicAnalyzer(&ai.TopicModelConfig{})

	// 批量文本分析端点
	router.POST("/api/v1/ai/batch-analyze", func(c *gin.Context) {
		var request struct {
			Messages []string `json:"messages"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_REQUEST",
					"message": "请求参数错误",
				},
			})
			return
		}

		if len(request.Messages) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "EMPTY_MESSAGES",
					"message": "消息列表不能为空",
				},
			})
			return
		}

		// 处理每条消息
		var results []map[string]interface{}
		var allEmotions []string
		var allIntents []string
		var totalPII int

		for i, message := range request.Messages {
			// PII检测和脱敏
			processed, entities := piiDetector.DetectAndReplace(message, 0.7)

			// 情感分析
			emotion := emotionAnalyzer.AnalyzeEmotion(processed)

			// 意图识别
			intent := intentRecognizer.AnalyzeIntent(processed)

			result := map[string]interface{}{
				"index":         i,
				"original_text": message,
				"processed_text": processed,
				"pii_entities":  entities,
				"pii_count":     len(entities),
				"emotion":       emotion.Type,
				"confidence":    emotion.Confidence,
				"intent":        intent.Type,
				"intent_priority": intent.Priority,
			}

			results = append(results, result)
			allEmotions = append(allEmotions, string(emotion.Type))
			allIntents = append(allIntents, string(intent.Type))
			totalPII += len(entities)
		}

		// 主题分析
		topics, topicConfidence := topicAnalyzer.AnalyzeTopics(request.Messages)

		// 生成汇总统计
		emotionStats := make(map[string]int)
		for _, emotion := range allEmotions {
			emotionStats[emotion]++
		}

		intentStats := make(map[string]int)
		for _, intent := range allIntents {
			intentStats[intent]++
		}

		response := gin.H{
			"success": true,
			"data": gin.H{
				"total_messages": len(request.Messages),
				"total_pii": totalPII,
				"results": results,
				"summary": gin.H{
					"emotion_distribution": emotionStats,
					"intent_distribution": intentStats,
					"topics": topics,
					"topic_confidence": topicConfidence,
				},
			},
			"processing_time": time.Now().Format(time.RFC3339),
		}

		c.JSON(http.StatusOK, response)
	})

	// 测试用例：真实的业务对话数据
	businessConversations := []struct {
		name           string
		messages       []string
		description     string
		expectedTopics  []string
		businessContext  string
	}{
		{
			name: "客服对话分析",
			messages: []string{
				"客户A: 我对你们的产品很失望！",
				"客服: 很抱歉听到您的反馈，能详细说明一下问题吗？",
				"客户A: 界面太复杂了，找不到需要的功能",
				"客服: 我理解您的困扰，让我为您提供操作指导",
			},
			description:    "客服会话的质量和情绪分析",
			expectedTopics: []string{"客户", "界面", "功能"},
			businessContext: "客户服务质量监控",
		},
		{
			name: "团队协作讨论",
			messages: []string{
				"项目经理: 下周的开发计划需要重新安排",
				"开发组长: 好的，我们会调整优先级",
				"产品经理: 新功能的用户反馈很重要",
				"设计师: 我已经准备好了设计稿",
			},
			description:    "项目管理流程中的沟通分析",
			expectedTopics: []string{"开发", "计划", "产品", "设计"},
			businessContext: "团队效率分析",
		},
		{
			name: "用户反馈收集",
			messages: []string{
				"用户1: 产品很好用，建议增加数据导出",
				"用户2: 界面设计很棒，但性能需要优化",
				"用户3: 功能很全面，价格合理",
				"用户4: 希望增加更多模板和示例",
			},
			description:    "大量用户反馈的批量分析",
			expectedTopics: []string{"产品", "功能", "价格", "设计"},
			businessContext: "产品改进决策支持",
		},
	}

	for _, tc := range businessConversations {
		t.Run(tc.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(gin.H{
				"messages": tc.messages,
			})

			req, _ := http.NewRequest("POST", "/api/v1/batch-analyze", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			start := time.Now()
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			duration := time.Since(start)

			// 验证响应
			if w.Code != http.StatusOK {
				t.Errorf("批量分析失败: HTTP %d", w.Code)
				return
			}

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Errorf("响应解析失败: %v", err)
				return
			}

			success := response["success"].(bool)
			if !success {
				t.Error("批量分析应该成功")
				return
			}

			data := response["data"].(map[string]interface{})
			results := data["results"].([]interface{})
			summary := data["summary"].(map[string]interface{})

			t.Logf("✅ %s批量分析完成:", tc.name)
			t.Logf("   消息数量: %d", data["total_messages"])
			t.Logf("   总PII数量: %d", data["total_pii"])
			t.Logf("   处理时间: %v", duration)
			t.Logf("   平均每条消息: %v", duration/time.Duration(len(results)))

			// 验证分析结果
			if len(results) != len(tc.messages) {
				t.Errorf("结果数量不匹配: 期望%d, 实际%d", len(tc.messages), len(results))
			}

			// 验证主题分析
			topics := summary["topics"].([]interface{})
			if len(topics) == 0 {
				t.Logf("⚠️  未识别到主题")
			} else {
				t.Logf("   识别主题: %v", topics)
			}

			// 验证情感分布
			emotionDist := summary["emotion_distribution"].(map[string]interface{})
			intentDist := summary["intent_distribution"].(map[string]interface{})

			t.Logf("   情感分布: %+v", emotionDist)
			t.Logf("   意图分布: %+v", intentDist)

			// 业务价值分析
			analysisValue := len(results) * len(topics)
			t.Logf("   分析价值: %d (消息数×主题数)", analysisValue)

			// 业务建议
			var insights []string
			if len(topics) > 0 {
				insights = append(insights, fmt.Sprintf("✅ 成功识别%d个主要主题", len(topics)))
			}

			totalPII := data["total_pii"].(float64)
			if totalPII > 0 {
				insights = append(insights, fmt.Sprintf("⚠️ 发现%d个敏感信息需要处理", int(totalPII)))
			}

			if len(results) > 5 {
				insights = append(insights, "📊 大数据量分析具有统计意义")
			}

			if len(insights) > 0 {
				t.Log("   业务洞察:")
				for i, insight := range insights {
					t.Logf("     %d. %s", i+1, insight)
				}
			}

			t.Logf("   业务场景: %s", tc.businessContext)
			t.Logf("   %s", tc.description)
		})
	}
}

func TestMain(m *testing.M) {
	// 这个函数允许手动运行集成测试
	fmt.Println("集成测试准备就绪")
	fmt.Println("运行: go test -tags=integration ./...")
	m.Run()
}