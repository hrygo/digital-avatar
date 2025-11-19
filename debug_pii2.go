package main

import (
	"fmt"
	"twin-os/backend/pkg/crypto"
)

func main() {
	detector := crypto.NewPIIDetector()

	tests := []string{
		"我的邮箱是test@example.com，身份证号是11010119900307123X。",
		"我住在北京市朝阳区某某路123号，在腾讯有限公司工作。",
		"这是一段不包含敏感信息的文本。",
	}

	for i, text := range tests {
		fmt.Printf("\n=== Test %d ===\n", i+1)
		fmt.Printf("Original: %s\n", text)

		processedText, report := detector.DetectAndReplace(text)
		fmt.Printf("Processed: %s\n", processedText)
		fmt.Printf("Total detections: %d\n", report.TotalDetections)

		for j, detection := range report.Detections {
			fmt.Printf("  %d: Type=%s, Original='%s', Replaced='%s'\n",
				j+1, detection.Type, detection.Original, detection.Replaced)
		}
	}
}