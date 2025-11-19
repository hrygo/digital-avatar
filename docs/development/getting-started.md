# 快速开始指南

欢迎参与Digital Avatar项目！本指南将帮助您在5分钟内搭建开发环境并开始贡献代码。

## 🎯 项目概览

Digital Avatar是一个智能决策副驾系统，提供：
- 🤖 **AI智能分析**: 基于NLP的文本分析和决策支持
- 🔒 **隐私保护**: 先进的PII检测和数据脱敏技术
- 📊 **数据洞察**: 智能数据分析和可视化
- 🌐 **全栈应用**: 现代化的前后端分离架构

## 🛠️ 技术栈

### 后端技术栈
- **Go 1.21+**: 高性能后端服务
- **PostgreSQL**: 主数据库
- **Redis**: 缓存和会话存储
- **Docker**: 容器化部署
- **gRPC**: 服务间通信

### 前端技术栈
- **React 18**: 用户界面框架
- **TypeScript**: 类型安全的JavaScript
- **Vite**: 快速构建工具
- **Tailwind CSS**: 样式框架
- **Ant Design**: UI组件库

## 🚀 快速启动

### 1. 环境准备

**系统要求**:
- Git 2.30+
- Go 1.21+
- Node.js 18+
- Docker 20.10+
- PostgreSQL 14+ (可选)

### 2. 克隆项目

```bash
git clone https://github.com/digital-avatar/digital-avatar.git
cd digital-avatar
```

### 3. 后端环境搭建

```bash
# 进入后端目录
cd backend

# 安装Go依赖
go mod tidy

# 复制环境配置
cp .env.example .env

# 编辑配置文件
vim .env
```

**配置示例** (.env):
```env
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=digital_avatar
DB_PASSWORD=your_password
DB_NAME=digital_avatar

# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your_redis_password

# 服务配置
SERVER_PORT=8080
JWT_SECRET=your_jwt_secret
```

### 4. 前端环境搭建

```bash
# 进入前端目录
cd ../frontend

# 安装依赖
npm install
# 或使用 yarn
yarn install

# 复制环境配置
cp .env.example .env.local

# 编辑配置文件
vim .env.local
```

**配置示例** (.env.local):
```env
# API配置
VITE_API_BASE_URL=http://localhost:8080
VITE_API_TIMEOUT=30000

# 应用配置
VITE_APP_TITLE=Digital Avatar
VITE_APP_VERSION=1.0.0
```

### 5. 启动开发服务

**启动后端服务**:
```bash
cd backend
go run cmd/server/main.go
```
服务将在 http://localhost:8080 启动

**启动前端服务**:
```bash
cd frontend
npm run dev
# 或
yarn dev
```
应用将在 http://localhost:3000 启动

### 6. 验证安装

打开浏览器访问 http://localhost:3000，您应该能看到Digital Avatar的欢迎界面。

## 🧪 运行测试

### 后端测试

```bash
# 运行所有测试
go test ./...

# 运行PII模块测试（我们的核心模块）
go test ./pkg/crypto -v -cover

# 生成覆盖率报告
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### 前端测试

```bash
# 运行单元测试
npm run test
# 或
yarn test

# 运行测试并生成覆盖率
npm run test:coverage
# 或
yarn test:coverage

# 运行E2E测试
npm run test:e2e
# 或
yarn test:e2e
```

## 📁 项目结构

```
digital-avatar/
├── docs/                      # 项目文档中心
│   ├── product/              # 产品文档
│   ├── project/              # 项目管理
│   ├── development/          # 开发指南
│   ├── api/                  # API文档
│   └── deployment/           # 部署运维
├── backend/                  # 后端服务
│   ├── cmd/                  # 应用入口
│   ├── internal/             # 内部代码
│   ├── pkg/                  # 公共包
│   │   ├── crypto/           # PII检测模块
│   │   ├── ai/               # AI分析模块
│   │   └── database/         # 数据库模块
│   └── docs/                 # 后端文档
├── frontend/                 # 前端应用
│   ├── src/                  # 源代码
│   │   ├── components/       # 组件
│   │   ├── pages/            # 页面
│   │   ├── hooks/            # 自定义Hooks
│   │   ├── services/         # API服务
│   │   └── store/            # 状态管理
│   └── docs/                 # 前端文档
├── scripts/                  # 构建和部署脚本
├── docker-compose.yml        # 开发环境容器编排
└── README.md                # 项目说明
```

## 🔧 开发工具配置

### VS Code配置

推荐安装以下VS Code扩展：

```json
{
  "recommendations": [
    "golang.go",
    "ms-vscode.vscode-typescript-next",
    "bradlc.vscode-tailwindcss",
    "esbenp.prettier-vscode",
    "ms-vscode.vscode-eslint",
    "ms-vscode-remote.remote-containers"
  ]
}
```

**设置** (.vscode/settings.json):
```json
{
  "go.useLanguageServer": true,
  "go.formatTool": "goimports",
  "go.lintTool": "golangci-lint",
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.fixAll.eslint": true
  },
  "typescript.preferences.importModuleSpecifier": "relative"
}
```

### Git配置

**提交前钩子** (.husky/pre-commit):
```bash
#!/bin/sh
# 运行Go格式化
go fmt ./...
go vet ./...

# 运行前端检查
cd frontend && npm run lint

# 运行测试
go test ./...
cd frontend && npm run test:unit
```

## 🎯 核心功能开发

### PII检测系统

我们最核心的功能是PII（个人敏感信息）检测：

```go
package main

import (
    "fmt"
    "github.com/digital-avatar/backend/pkg/crypto"
)

func main() {
    // 创建PII检测器
    detector := crypto.NewNLPPIIDetector()

    // 检测文本中的敏感信息
    text := "联系人：张三，手机：13812345678，邮箱：zhang@example.com"
    processed, entities, err := detector.DetectAndReplace(text, 0.7)

    if err != nil {
        panic(err)
    }

    fmt.Printf("原文: %s\n", text)
    fmt.Printf("处理后: %s\n", processed)
    fmt.Printf("检测到 %d 个敏感实体\n", len(entities))
}
```

### 前端组件开发

```typescript
import React from 'react';
import { Card, Button } from 'antd';
import { usePIIDetection } from '@/hooks/usePIIDetection';

interface PIIProcessorProps {
  text: string;
  onProcessed: (result: string) => void;
}

export const PIIProcessor: React.FC<PIIProcessorProps> = ({
  text,
  onProcessed
}) => {
  const { process, loading, result } = usePIIDetection();

  const handleProcess = async () => {
    const processedText = await process(text);
    onProcessed(processedText);
  };

  return (
    <Card title="PII处理器">
      <p>原文: {text}</p>
      <Button
        type="primary"
        onClick={handleProcess}
        loading={loading}
      >
        处理敏感信息
      </Button>
      {result && <p>处理后: {result}</p>}
    </Card>
  );
};
```

## 🤝 贡献指南

### 分支策略

- **main**: 生产环境分支
- **develop**: 开发环境分支
- **feature/***: 功能开发分支
- **hotfix/***: 紧急修复分支

### 提交规范

使用[约定式提交](https://www.conventionalcommits.org/)：

```bash
# 功能开发
git commit -m "feat(ai): add new sentiment analysis model"

# Bug修复
git commit -m "fix(pii): resolve phone number detection issue"

# 文档更新
git commit -m "docs(api): update authentication endpoints"

# 性能优化
git commit -m "perf(backend): optimize database query performance"
```

### 代码审查

1. **创建Pull Request**
2. **通过所有测试**
3. **代码覆盖率不能降低**
4. **通过至少一人审查**
5. **合并到develop分支**

## 📞 获取帮助

### 联系方式

- **技术问题**: [GitHub Issues](https://github.com/digital-avatar/digital-avatar/issues)
- **功能讨论**: [GitHub Discussions](https://github.com/digital-avatar/digital-avatar/discussions)
- **团队沟通**: Slack频道 #digital-avatar

### 学习资源

- [项目文档](../README.md)
- [API文档](api/README.md)
- [PII模块文档](../../backend/pkg/crypto/README.md)
- [开发规范](coding-standards.md)

### 常见问题

**Q: 后端服务启动失败？**
A: 检查数据库连接配置，确保PostgreSQL和Redis服务正常运行。

**Q: 前端API调用失败？**
A: 确认后端服务在8080端口启动，检查网络代理配置。

**Q: PII检测不准确？**
A: 查看测试覆盖率报告，考虑调整检测阈值或添加更多训练数据。

---

🎉 **恭喜！** 您已经成功搭建了Digital Avatar开发环境。现在可以开始贡献代码了！

如果遇到任何问题，请随时通过上述渠道联系我们。