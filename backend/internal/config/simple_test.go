package config

import (
	"os"
	"testing"
)

func TestGetEnvSimple(t *testing.T) {
	// 测试环境变量获取
	os.Setenv("TEST_KEY", "test_value")
	defer os.Unsetenv("TEST_KEY")

	result := getEnv("TEST_KEY", "default")
	if result != "test_value" {
		t.Errorf("getEnv() = %v, expected test_value", result)
	}

	// 测试默认值
	result = getEnv("NON_EXISTENT_KEY", "default")
	if result != "default" {
		t.Errorf("getEnv() with non-existent key = %v, expected default", result)
	}
}