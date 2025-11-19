# Digital Avatar 后端文档

## 📚 文档导航

本文档是Digital Avatar后端系统的技术文档中心，提供架构设计、API文档、部署指南和技术报告。

## 🗂️ 文档分类

### 🏗️ 系统架构
- [系统架构概览](architecture/README.md) - 整体架构和设计决策
- [数据库设计](database/README.md) - 数据模型和关系设计
- [微服务架构](architecture/microservices.md) - 服务拆分和通信机制

### 🔌 API文档
- [API导航](../docs/api/README.md) - 完整的API文档索引
- [认证授权](../docs/api/authentication.md) - JWT认证和权限管理
- [AI功能接口](../docs/api/endpoints/ai/) - AI分析和处理API
- [PII处理接口](../docs/api/endpoints/pii/) - 隐私信息处理API
- [用户管理接口](../docs/api/endpoints/users/) - 用户CRUD操作

### 🧪 测试报告
- [测试覆盖率报告](reports/test-coverage/) - 代码测试覆盖率分析
- [性能测试报告](reports/performance/) - 系统性能测试结果
- [安全测试报告](reports/security/) - 安全漏洞和渗透测试

### 🔧 技术实现
- [开发指南](guides/) - 后端开发最佳实践
- [核心技术模块](../pkg/crypto/README.md) - PII检测和加密模块
- [中间件使用](guides/middleware.md) - 认证、日志、监控中间件

### 📊 技术报告
- [边界情况优化](reports/technical/BOUNDARY_OPTIMIZATION_REPORT.md) - PII检测优化报告
- [代码清理报告](reports/technical/CLEANUP_REPORT.md) - 重构和清理记录
- [架构演进](reports/technical/architecture-evolution.md) - 架构变更历史

## 🚀 核心模块

### PII检测系统 ([pkg/crypto/](../pkg/crypto/))
基于NLP技术的个人敏感信息检测和脱敏系统。

**核心特性**:
- 🎯 **检测精度**: 90.5%测试覆盖率，95%+检测准确率
- ⚡ **处理性能**: <1ms/100字符，毫秒级响应
- 🔒 **数据安全**: AES-GCM加密 + SHA256哈希脱敏
- 🔄 **向后兼容**: 100%API兼容性

**快速使用**:
```go
import "github.com/digital-avatar/backend/pkg/crypto"

// 创建PII检测器
detector := crypto.NewNLPPIIDetector()

// 检测和替换敏感信息
processed, entities, err := detector.DetectAndReplace(text, 0.7)
```

### AI分析引擎
智能文本分析和决策支持引擎。

**功能模块**:
- 🧠 **智能分析**: 基于预训练模型
- 📊 **数据洞察**: 趋势分析和预测
- 🔍 **内容理解**: 语义分析和实体识别
- ⚙️ **配置管理**: 灵活的模型配置

### 数据访问层
数据库操作和缓存管理系统。

**技术栈**:
- 🗄️ **数据库**: PostgreSQL/MySQL支持
- 🚀 **ORM**: GORM集成
- 💾 **缓存**: Redis多级缓存
- 📊 **监控**: 查询性能监控

## 🛠️ 开发指南

### 环境搭建
1. **克隆项目**
```bash
git clone https://github.com/digital-avatar/backend.git
cd backend
```

2. **安装依赖**
```bash
go mod tidy
```

3. **配置环境**
```bash
cp .env.example .env
# 编辑.env配置数据库连接等
```

4. **启动服务**
```bash
go run cmd/server/main.go
```

### 测试运行
```bash
# 运行所有测试
go test ./...

# 运行PII模块测试
go test ./pkg/crypto -v -cover

# 生成覆盖率报告
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### 代码质量
- **测试覆盖率**: 目标 ≥ 80%
- **代码审查**: 所有PR必须经过审查
- **静态分析**: 使用golangci-lint
- **文档更新**: 代码变更同步更新文档

## 📈 性能指标

### 系统性能
- **响应时间**: < 100ms (95%ile)
- **吞吐量**: > 1000 RPS
- **并发用户**: > 1000
- **可用性**: 99.9%

### PII检测性能
- **检测速度**: < 1ms/100字符
- **内存效率**: 线性内存使用
- **准确率**: 95%+ 真实场景
- **一致性**: 100% 相同实体替换

## 🔒 安全保障

### 数据安全
- **传输加密**: TLS 1.3
- **存储加密**: AES-256
- **密钥管理**: 安全的密钥轮换
- **访问控制**: RBAC权限模型

### PII处理安全
- **不可逆脱敏**: SHA256哈希替换
- **一致性保证**: 相同信息相同替换
- **审计日志**: 完整的操作记录
- **合规支持**: GDPR/CCPA兼容

## 🚀 部署运维

### 容器化部署
```dockerfile
# Dockerfile示例
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

### Kubernetes部署
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: digital-avatar-backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: digital-avatar-backend
  template:
    metadata:
      labels:
        app: digital-avatar-backend
    spec:
      containers:
      - name: backend
        image: digital-avatar/backend:latest
        ports:
        - containerPort: 8080
```

## 📞 技术支持

### 问题反馈
- **GitHub Issues**: 技术问题和Bug报告
- **技术讨论**: GitHub Discussions
- **紧急问题**: 联系技术负责人

### 开发团队
- **架构师**: 系统架构和技术决策
- **后端开发**: 功能开发和性能优化
- **DevOps**: 部署和运维支持

---

## 📊 文档统计

- **API文档**: 15个接口
- **测试覆盖率**: 90.5%
- **核心模块**: 8个
- **技术报告**: 5份

*最后更新: 2024年11月*