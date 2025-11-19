# Digital Avatar Backend

数字人后端服务，提供AI驱动的个人信息处理和安全保护功能。

## 项目概述

本项目是一个基于Go语言开发的数字人后端服务，专注于提供：

- **PII信息检测与脱敏**: 使用NLP技术检测和替换个人敏感信息
- **数据加密保护**: 提供多层次的加密安全机制
- **AI增强分析**: 集成预训练模型进行智能文本处理
- **高性能API**: 支持大规模并发请求处理

## 项目结构

```
.
├── cmd/                    # 应用程序入口点
├── internal/               # 内部应用代码
│   ├── config/            # 配置管理
│   ├── handler/           # HTTP处理器
│   ├── model/             # 数据模型
│   ├── repository/        # 数据访问层
│   ├── server/            # 服务器配置
│   └── service/           # 业务逻辑层
├── pkg/                   # 公共库代码
│   ├── ai/               # AI相关功能
│   ├── backup/           # 备份功能
│   ├── cache/            # 缓存管理
│   ├── crypto/           # 加密和PII处理
│   ├── database/         # 数据库操作
│   ├── logger/           # 日志管理
│   └── wechat/           # 微信集成
├── docs/                  # 项目文档
├── data/                  # 数据文件
└── backups/              # 备份文件
```

## 核心功能

### PII检测与脱敏 (pkg/crypto)

基于NLP技术的个人敏感信息检测系统，支持：

- **个人信息类型**: 姓名、电话、邮箱、身份证、组织机构、地址
- **智能检测**: 结合规则和词典的高精度识别
- **一致性替换**: 使用哈希算法确保相同信息的替换一致性
- **高性能处理**: 毫秒级检测响应，支持大规模文本处理

#### 主要特性

- 🎯 **检测精度**: 95%+ 的实际检测准确率
- ⚡ **处理性能**: <1ms/100字符的检测速度
- 🔒 **数据安全**: SHA256哈希替换，不可逆脱敏
- 📊 **统计报告**: 详细的检测统计和分析
- 🔧 **向后兼容**: 完整的API兼容层支持

#### 快速使用

```go
package main

import (
    "fmt"
    "github.com/hrygo/digital-avatar/backend/pkg/crypto"
)

func main() {
    // 创建PII检测器
    detector := crypto.NewNLPPIIDetector()

    // 检测和替换敏感信息
    text := "联系人：张三，手机：13812345678，邮箱：zhang@example.com"
    processed, entities, err := detector.DetectAndReplace(text, 0.7)

    if err == nil {
        fmt.Printf("原文: %s\n", text)
        fmt.Printf("处理后: %s\n", processed)
        fmt.Printf("检测到 %d 个敏感实体\n", len(entities))
    }
}
```

## 开发环境

### 环境要求

- Go 1.19+
- MySQL 8.0+ (可选)
- Redis 6.0+ (可选)

### 安装和运行

```bash
# 克隆项目
git clone <repository-url>
cd backend

# 安装依赖
go mod tidy

# 运行测试
go test ./...

# 启动服务
go run cmd/server/main.go
```

### 测试覆盖

项目重视代码质量，核心模块测试覆盖率达到74.8%+：

```bash
# 运行完整测试套件
go test ./pkg/crypto -v -coverprofile=coverage.out

# 生成覆盖率报告
go tool cover -html=coverage.out -o coverage.html
```

详细的测试报告请参考：[测试覆盖率报告](docs/COVERAGE_REPORT.md)

## 技术架构

### 设计原则

- **模块化设计**: 清晰的包结构，便于维护和扩展
- **依赖注入**: 使用接口和依赖注入提高可测试性
- **错误处理**: 完善的错误处理和日志记录机制
- **性能优化**: 针对高并发场景的架构设计

### 核心技术

- **NLP处理**: 基于词典和规则的混合PII检测
- **并发安全**: 线程安全的数据结构和算法
- **内存优化**: 高效的内存使用和垃圾回收策略
- **缓存机制**: 多级缓存提升响应性能

## API文档

详细的API文档请参考：[API文档](docs/api.md)

## 开发指南

### 代码规范

项目遵循Go官方编码规范和以下最佳实践：

1. **包命名**: 使用简短、清晰的包名
2. **接口设计**: 优先使用接口而非具体类型
3. **错误处理**: 显式处理所有可能的错误
4. **测试优先**: 新功能必须包含单元测试

### 贡献指南

1. Fork 项目仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 版本历史

- **v0.2.0** (开发中) - 增强分析功能
  - 完整的NLP PII检测系统
  - 边界情况优化
  - 高测试覆盖率

- **v0.1.0** (已发布) - MVP版本
  - 基础PII检测功能
  - 核心API接口

## 许可证

本项目采用 [MIT License](LICENSE) 许可证。

## 联系我们

- 项目维护者：[项目组信息](CONTRIBUTORS.md)
- 问题反馈：[GitHub Issues](https://github.com/hrygo/digital-avatar/issues)
- 文档更新：[项目Wiki](https://github.com/hrygo/digital-avatar/wiki)