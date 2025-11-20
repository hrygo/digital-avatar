package handler

import (
	"testing"
)

func TestWeChatHandlerCreation(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("WeChatHandler creation panicked: %v", r)
		}
	}()

	// 测试创建函数不会崩溃
	_ = NewWeChatHandler(nil)
}

func TestAnalysisHandlerCreation(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("AnalysisHandler creation panicked: %v", r)
		}
	}()

	_ = NewAnalysisHandler(nil)
}

func TestDataHandlerCreation(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("DataHandler creation panicked: %v", r)
		}
	}()

	_ = NewDataHandler(nil)
}

func TestSettingsHandlerCreation(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("SettingsHandler creation panicked: %v", r)
		}
	}()

	_ = NewSettingsHandler(nil)
}

func TestBackupHandlerCreation(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("BackupHandler creation panicked: %v", r)
		}
	}()

	_ = NewBackupHandler(nil)
}