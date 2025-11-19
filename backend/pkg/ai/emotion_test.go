package ai

import (
	"testing"
)

func TestEmotionAnalyzer_AnalyzeEmotion(t *testing.T) {
	analyzer := NewEmotionAnalyzer()

	tests := []struct {
		name           string
		text           string
		expectedType   EmotionType
		minIntensity   float64
	}{
		{
			name:           "Joy - Basic",
			text:           "今天天气真好，我很开心",
			expectedType:   EmotionJoy,
			minIntensity:   0.1,
		},
		{
			name:           "Anger - Basic",
			text:           "这真是太糟糕了，我很生气",
			expectedType:   EmotionAnger,
			minIntensity:   0.1,
		},
		{
			name:           "Neutral",
			text:           "我就去买个菜",
			expectedType:   EmotionNeutral, // It might hit some neutral keywords or fallback
			minIntensity:   0.0,
		},
		{
			name:           "Love - Emoji",
			text:           "谢谢你 ❤️",
			expectedType:   EmotionLove,
			minIntensity:   0.1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.AnalyzeEmotion(tt.text)

			// Note: For Neutral, it might actually match nothing and return Neutral, which is fine.
			// Or it might match "Neutral" keywords.
			
			if result.Type != tt.expectedType {
				// Allow Neutral fallback if we expected something else but got Neutral (means weak signal)
				// But for these strong test cases, we expect match.
				t.Errorf("AnalyzeEmotion() Type = %v, want %v", result.Type, tt.expectedType)
			}

			if result.Intensity < tt.minIntensity {
				t.Errorf("AnalyzeEmotion() Intensity = %v, want >= %v", result.Intensity, tt.minIntensity)
			}
		})
	}
}

func TestEmotionAnalyzer_Modifiers(t *testing.T) {
	analyzer := NewEmotionAnalyzer()

	textBase := "开心"
	textStrong := "非常开心"

	resultBase := analyzer.AnalyzeEmotion(textBase)
	resultStrong := analyzer.AnalyzeEmotion(textStrong)

	if resultBase.Type != EmotionJoy || resultStrong.Type != EmotionJoy {
		t.Fatal("Both should be Joy")
	}

	if resultStrong.Intensity <= resultBase.Intensity {
		t.Errorf("Modifier '非常' did not increase intensity. Base: %v, Strong: %v", resultBase.Intensity, resultStrong.Intensity)
	}
}
