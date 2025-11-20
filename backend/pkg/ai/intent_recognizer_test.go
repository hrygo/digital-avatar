package ai

import (
	"testing"
)

func TestIntentRecognizer_AnalyzeIntent(t *testing.T) {
	recognizer := NewIntentRecognizer()

	tests := []struct {
		name           string
		text           string
		expectedType   IntentType
		expectedSub    string
		minConfidence  float64
	}{
		{
			name:          "Task Intent - Meeting",
			text:          "我需要安排一个会议",
			expectedType:  IntentTask,
			expectedSub:   "meeting",
			minConfidence: 0.3,
		},
		{
			name:          "Question Intent - Information",
			text:          "请问这是什么？",
			expectedType:  IntentQuestion,
			expectedSub:   "information",
			minConfidence: 0.3,
		},
		{
			name:          "Emotional Intent - Happy",
			text:          "我今天很开心",
			expectedType:  IntentEmotional,
			expectedSub:   "",
			minConfidence: 0.3,
		},
		{
			name:          "Decision Intent",
			text:          "我正在考虑选择哪个方案",
			expectedType:  IntentDecision,
			expectedSub:   "",
			minConfidence: 0.3,
		},
		{
			name:          "Neutral Intent",
			text:          "",
			expectedType:  IntentNeutral,
			expectedSub:   "",
			minConfidence: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := recognizer.AnalyzeIntent(tt.text)

			if result.Type != tt.expectedType {
				t.Errorf("AnalyzeIntent() Type = %v, want %v", result.Type, tt.expectedType)
			}

			if tt.expectedSub != "" && result.SubType != tt.expectedSub {
				t.Errorf("AnalyzeIntent() SubType = %v, want %v", result.SubType, tt.expectedSub)
			}

			if result.Confidence < tt.minConfidence {
				t.Errorf("AnalyzeIntent() Confidence = %v, want >= %v", result.Confidence, tt.minConfidence)
			}
		})
	}
}

func TestIntentRecognizer_AnalyzeImplicitIntent(t *testing.T) {
	recognizer := NewIntentRecognizer()

	context := []string{
		"项目截止日期快到了",
		"我们还有很多工作没做",
	}

	result := recognizer.AnalyzeImplicitIntent(context)

	// 隐式意图通常会被识别为 Task 或 Emotional（焦虑）
	// 但主要是检查是否返回了结果且置信度经过了调整
	if result == nil {
		t.Error("AnalyzeImplicitIntent() returned nil")
	}

	if result.Confidence > 0.8 {
		t.Errorf("AnalyzeImplicitIntent() Confidence = %v, expected it to be dampened (<= 0.8)", result.Confidence)
	}
}

func TestIntentRecognizer_TrackIntentPattern(t *testing.T) {
	recognizer := NewIntentRecognizer()

	// 模拟一个重复的模式: Question -> Task
	analyses := []*IntentAnalysis{
		{Type: IntentQuestion},
		{Type: IntentTask},
		{Type: IntentQuestion},
		{Type: IntentTask},
		{Type: IntentQuestion},
		{Type: IntentTask},
	}

	patterns := recognizer.TrackIntentPattern(analyses)

	if len(patterns) == 0 {
		t.Error("TrackIntentPattern() returned empty patterns")
	}

	found := false
	for _, p := range patterns {
		if p.Context == "意图转换模式：question->task" {
			found = true
			break
		}
	}

	if !found {
		t.Error("TrackIntentPattern() did not find expected pattern 'question->task'")
	}
}
