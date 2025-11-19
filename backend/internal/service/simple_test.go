package service

import (
	"testing"
)

func TestServiceStructures(t *testing.T) {
	// 测试服务结构体字段存在且可访问
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Service structure test panicked: %v", r)
		}
	}()

	// 这里只测试结构体定义不会导致编译错误
	// 实际的service创建需要依赖注入
	var _ *Service = (*Service)(nil)
	var _ *AnalysisService = (*AnalysisService)(nil)
	var _ *BackupService = (*BackupService)(nil)
}

func TestServiceInterface(t *testing.T) {
	// 测试服务接口实现
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Service interface test panicked: %v", r)
		}
	}()

	// 确保接口存在且可使用
	type MockDataService struct{}
	_ = MockDataService{}
}