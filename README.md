# TwinOS - 智能决策副驾系统

<div align="center">

![TwinOS Logo](https://via.placeholder.com/200x80/000000/FFFFFF?text=TwinOS)

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![React Version](https://img.shields.io/badge/React-18+-61DAFB?style=flat-square&logo=react)](https://reactjs.org)
[![License](https://img.shields.io/badge/License-MIT-blue?style=flat-square)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen?style=flat-square)](https://github.com/your-username/twin-os/actions)

**你就负责思考，剩下的交给它**

[快速开始](#快速开始) • [功能特性](#功能特性) • [文档](#📚-文档) • [贡献](#贡献)

</div>

## 📖 项目简介

TwinOS是一个基于人工智能的智能决策副驾系统，专为提升个人和团队决策效率而设计。系统结合了深度学习、自然语言处理和知识图谱技术，为用户提供智能化的数据分析、决策支持和知识管理服务。

### 🎯 核心理念

- **智能辅助**: 通过AI技术增强人类决策能力
- **知识驱动**: 构建个人和组织知识图谱
- **数据安全**: 端到端加密，保护用户隐私
- **开放生态**: 模块化设计，支持插件扩展

## ✨ 功能特性

### 🧠 智能分析引擎
- **深度洞察**: 基于聊天记录的智能分析
- **待办提取**: 自动识别和整理任务清单
- **关系图谱**: 构建人脉关系和社交网络
- **情感分析**: 理解对话情绪和语调

### 💾 数据管理系统
- **微信集成**: 安全连接和同步微信数据
- **多源支持**: 支持多种数据源导入
- **实时同步**: 增量同步，保持数据最新
- **隐私保护**: PII自动脱敏，保护敏感信息

### 🛡️ 企业级安全
- **端到端加密**: 军用级数据加密
- **备份恢复**: 自动备份，一键恢复
- **访问控制**: 细粒度权限管理
- **审计日志**: 完整的操作记录追踪

### 🚀 高性能架构
- **微服务设计**: 松耦合，高可扩展性
- **智能缓存**: 多层缓存优化性能
- **容器化部署**: Docker/Kubernetes支持
- **监控告警**: 实时性能监控

## 🏗️ 技术架构

### 后端技术栈
- **语言**: Go 1.21+
- **框架**: Gin (高性能HTTP框架)
- **数据库**: SQLite (轻量级) + Redis (缓存)
- **AI引擎**: DeepSeek API
- **消息队列**: 内置任务调度器

### 前端技术栈
- **框架**: React 18 + TypeScript
- **构建工具**: Create React App
- **状态管理**: Zustand
- **UI组件**: Tailwind CSS + Headless UI
- **动画**: Framer Motion

### 部署架构
- **容器化**: Docker + Docker Compose
- **反向代理**: Nginx
- **进程管理**: Systemd
- **监控**: 内置性能监控

## 🚀 快速开始

### 环境要求

- **Go**: 1.21 或更高版本
- **Node.js**: 18.0 或更高版本
- **Git**: 最新版本
- **操作系统**: Linux/macOS/Windows

### 一键安装

```bash
# 克隆项目
git clone https://github.com/your-username/twin-os.git
cd twin-os

# 初始化开发环境
make init

# 启动开发服务器
make dev
```

### 手动安装

<details>
<summary>点击展开详细步骤</summary>

#### 1. 克隆项目
```bash
git clone https://github.com/your-username/twin-os.git
cd twin-os
```

#### 2. 安装依赖
```bash
# 安装Go依赖
cd backend
go mod download
go mod tidy

# 安装Node.js依赖
cd ../frontend
npm install
```

#### 3. 配置环境
```bash
# 复制环境配置文件
cp backend/.env.example backend/.env

# 编辑配置文件，设置API密钥等
vim backend/.env
```

#### 4. 启动服务
```bash
# 启动后端服务 (终端1)
cd backend
go run main.go

# 启动前端服务 (终端2)
cd frontend
npm start
```
</details>

### 访问应用

- **Web界面**: http://localhost:3000
- **API文档**: http://localhost:1234/health
- **后端服务**: http://localhost:1234

## 📚 文档

我们提供了详细的文档来帮助你更好地使用和参与项目：

### 🎯 产品文档
- **[产品需求文档 (PRD)](docs/product/prd.md)** - 了解产品定位、用户画像和核心功能
- **[产品设计文档](docs/product/product-design.md)** - 详细的UI/UX设计规范
- **[产品路线图](docs/product/product-roadmap.md)** - 版本规划和未来发展方向

### 🚀 项目管理
- **[项目启动文档](docs/project/kickoff.md)** - 项目背景、目标和团队介绍
- **[冲刺计划](docs/project/sprint-plan.md)** - 开发里程碑和迭代计划
- **[项目状态](docs/project/project-status.md)** - 当前进度和风险跟踪
- **[MVP范围调整](docs/project/mvp-scope-adjustment.md)** - 功能优先级和变更记录

### 📖 完整文档中心
👉 **[访问完整文档中心](docs/README.md)** 查看所有文档的详细导航

## 📖 使用指南

### 基础配置

1. **微信连接**: 在设置页面配置微信连接参数
2. **数据同步**: 首次使用需要同步历史数据
3. **AI配置**: 配置DeepSeek API密钥
4. **备份设置**: 设置自动备份策略

### 核心功能

#### 🔍 智能分析
- 自动分析聊天记录，提取关键信息
- 生成个人简报和洞察报告
- 识别重要对话和决策节点

#### 📋 待办管理
- 从对话中自动提取任务
- 智能分类和优先级排序
- 进度跟踪和提醒功能

#### 👥 关系网络
- 构建人脉关系图谱
- 分析互动频率和关系强度
- 发现潜在合作机会

## 🛠️ 开发指南

### 项目结构

```
twin-os/
├── backend/              # Go后端服务
│   ├── internal/        # 内部包
│   ├── pkg/            # 公共包
│   ├── main.go         # 入口文件
│   └── go.mod          # Go模块定义
├── frontend/           # React前端应用
│   ├── src/           # 源代码
│   ├── public/        # 静态资源
│   └── package.json   # Node.js依赖
├── docs/              # 项目文档
├── scripts/           # 构建脚本
├── build/             # 构建输出
├── dist/              # 发布文件
├── Makefile           # 构建命令
└── README.md          # 项目说明
```

### 开发命令

```bash
# 查看所有可用命令
make help

# 开发环境
make dev              # 启动开发服务器
make dev-server       # 仅启动后端
make dev-web          # 仅启动前端

# 构建
make build            # 构建所有组件
make build-server     # 构建后端
make build-web        # 构建前端

# 测试
make test             # 运行所有测试
make test-server      # 测试后端
make test-web         # 测试前端

# 代码质量
make lint             # 代码检查
make format           # 代码格式化

# 部署
make docker-build     # 构建Docker镜像
make docker-run       # 运行Docker容器
```

### 贡献代码

我们欢迎所有形式的贡献！请查看 [贡献指南](CONTRIBUTING.md) 了解详细信息。

#### 开发流程
1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 📦 部署指南

### Docker部署

```bash
# 构建镜像
make docker-build

# 运行容器
make docker-run

# 查看运行状态
docker ps
```

### 生产部署

<details>
<summary>生产环境部署指南</summary>

#### 1. 服务器要求
- CPU: 2核心以上
- 内存: 4GB以上
- 存储: 20GB以上
- 操作系统: Linux (推荐Ubuntu 20.04+)

#### 2. 环境准备
```bash
# 安装Docker
curl -fsSL https://get.docker.com | sh

# 安装Docker Compose
pip install docker-compose

# 克隆项目
git clone https://github.com/your-username/twin-os.git
cd twin-os
```

#### 3. 配置文件
```bash
# 复制生产配置
cp docker-compose.prod.yml docker-compose.yml

# 编辑环境变量
vim .env.production
```

#### 4. 启动服务
```bash
# 构建和启动
docker-compose up -d

# 查看服务状态
docker-compose ps
```
</details>

## 🔧 配置说明

### 环境变量

| 变量名 | 描述 | 默认值 | 必需 |
|--------|------|--------|------|
| `PORT` | 服务端口 | 1234 | ❌ |
| `GIN_MODE` | 运行模式 | debug | ❌ |
| `DB_PATH` | 数据库路径 | ./data/twin-os.db | ❌ |
| `DEEPSEEK_API_KEY` | AI服务密钥 | - | ✅ |
| `JWT_SECRET` | JWT密钥 | random | ✅ |
| `BACKUP_DIR` | 备份目录 | ./backups | ❌ |

### 完整配置示例

```bash
# backend/.env
PORT=8080
GIN_MODE=release
DB_PATH=/app/data/twin-os.db
DEEPSEEK_API_KEY=your_api_key_here
JWT_SECRET=your_jwt_secret_here
BACKUP_DIR=/app/backups
LOG_LEVEL=info
MAX_BACKUP_FILES=10
```

## 🤝 社区支持

- **GitHub Issues**: [报告问题](https://github.com/your-username/twin-os/issues)
- **GitHub Discussions**: [社区讨论](https://github.com/your-username/twin-os/discussions)
- **Wiki**: [详细文档](https://github.com/your-username/twin-os/wiki)

## 📊 项目状态

### 开发进度

- [x] 核心架构设计
- [x] 用户认证系统
- [x] 数据同步引擎
- [x] AI分析模块
- [x] 备份恢复系统
- [x] Web用户界面
- [ ] 移动端应用
- [ ] 插件系统
- [ ] 企业版功能

### 路线图

#### v1.0 (当前版本)
- ✅ 基础功能实现
- ✅ Web应用发布
- ✅ Docker支持

#### v1.1 (计划中)
- 🔄 移动端适配
- 🔄 更多AI模型支持
- 🔄 团队协作功能

#### v2.0 (未来版本)
- 📋 插件生态系统
- 📋 企业级功能
- 📋 云端部署版本

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

感谢以下开源项目的支持：

- [Gin](https://github.com/gin-gonic/gin) - Go Web框架
- [React](https://github.com/facebook/react) - 用户界面库
- [DeepSeek](https://github.com/deepseek-ai) - AI模型支持
- [Tailwind CSS](https://github.com/tailwindlabs/tailwindcss) - CSS框架

---

<div align="center">

**[⬆ 回到顶部](#twinos---智能决策副驾系统)**

Made with ❤️ by the TwinOS Team

</div>