package logger

import (
	"bytes"
	"testing"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name string
		level string
	}{
		{"debug_level", "debug"},
		{"info_level", "info"},
		{"warn_level", "warn"},
		{"error_level", "error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 测试logger初始化不会崩溃
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Init(%s) panicked: %v", tt.level, r)
				}
			}()

			Init(tt.level)
		})
	}
}

func TestLoggerMethods(t *testing.T) {
	// 初始化logger
	Init("info")

	// 重定向输出到缓冲区以测试
	oldOutput := Logger.Out
	defer func() {
		Logger.Out = oldOutput
	}()

	var buf bytes.Buffer
	Logger.Out = &buf

	tests := []struct {
		name string
		action func()
	}{
		{"info", func() { Info("test info") }},
		// {"debug", func() { Debug("test debug") }},
		{"warn", func() { Warn("test warn") }},
		{"error", func() { Error("test error") }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.action()

			// 验证日志被写入
			if buf.Len() == 0 {
				t.Errorf("Logger method %s should write output", tt.name)
			}
		})
	}
}