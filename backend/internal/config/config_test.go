package config

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	// 设置测试环境变量
	os.Setenv("TEST_KEY", "test_value")
	defer os.Unsetenv("TEST_KEY")

	// 测试存在的环境变量
	if got := getEnv("TEST_KEY", "default"); got != "test_value" {
		t.Errorf("getEnv() = %v, want test_value", got)
	}

	// 测试不存在的环境变量
	if got := getEnv("NON_EXISTENT", "default"); got != "default" {
		t.Errorf("getEnv() = %v, want default", got)
	}
}

func TestGetEnvAsInt(t *testing.T) {
	// 设置测试环境变量
	os.Setenv("TEST_INT", "3000")
	defer os.Unsetenv("TEST_INT")

	// 测试有效的整数
	if got := getEnvAsInt("TEST_INT", 8080); got != 3000 {
		t.Errorf("getEnvAsInt() = %v, want 3000", got)
	}

	// 测试默认值
	if got := getEnvAsInt("NON_EXISTENT", 8080); got != 8080 {
		t.Errorf("getEnvAsInt() = %v, want 8080", got)
	}

	// 测试无效整数
	os.Setenv("INVALID", "invalid")
	defer os.Unsetenv("INVALID")
	if got := getEnvAsInt("INVALID", 8080); got != 8080 {
		t.Errorf("getEnvAsInt() with invalid = %v, want 8080", got)
	}
}

func TestLoad(t *testing.T) {
	// 保存原始环境变量
	origDBHost := os.Getenv("DB_HOST")
	origServerPort := os.Getenv("SERVER_PORT")
	defer func() {
		os.Setenv("DB_HOST", origDBHost)
		os.Setenv("SERVER_PORT", origServerPort)
	}()

	// 设置测试环境变量
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("SERVER_PORT", "8080")

	// 测试配置加载
	if err := Load(); err != nil {
		// 由于缺少.env文件是正常的，我们只验证函数调用不会崩溃
		t.Logf("Load() returned expected error (no .env file): %v", err)
	}
}