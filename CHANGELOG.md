# 📝 变更日志 (CHANGELOG)

所有重要的项目变更都会记录在此文件中。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
项目遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [未发布]

### 新增
- GitHub开源项目发布准备
- 项目目录结构扁平化优化
- 文档重组和分类管理

---

## [0.2.0-dev] - 2025-11-19

### 🌟 新增功能 (Enhanced Features)
- **🧠 增强型AI分析引擎**
  - **意图识别系统**: 支持任务(Task)、决策(Decision)、社交(Social)等12种意图类型的自动识别，包含置信度评分和上下文感知。
  - **情感计算模块**: 实现Joy, Anger, Sadness等8种细粒度情感分析，支持强度计算和情感趋势(Trend)追踪。
  - **主题建模引擎**: 基于简化的TF-IDF和共现分析算法，自动提取对话主题和关键词。
  - **本地化处理**: 引入正则表达式引擎和本地词典，减少对云端LLM的依赖，提高响应速度和隐私安全性。

### 🛠 技术改进
- **单元测试**: 新增 `pkg/ai` 模块的完整单元测试覆盖 (`intent_test.go`, `emotion_test.go`, `topic_test.go`)。
- **API扩展**: 新增 `/api/v1/analysis/*` 系列接口，支持批量分析、趋势查询和主题统计。
- **代码重构**: 优化了Handler层与Service层的依赖注入结构。

---

## [0.1.0] - 2025-11-19

### 🎉 首次发布 - Digital Avatar 智能决策副驾系统 MVP

#### 🎯 MVP核心功能
- **🧠 基础智能分析**
  - 集成DeepSeek API进行聊天记录分析
  - 自动提取待办事项和行动项
  - 基础人脉关系识别
  - 生成简单的个人简报

- **💾 数据同步管理**
  - 微信聊天记录导入功能
  - 基础数据存储和管理
  - PII敏感信息基础脱敏
  - 数据备份和恢复功能

- **⚡ 系统基础架构**
  - Go后端 + React前端技术栈
  - SQLite数据库存储
  - 基础RESTful API设计
  - Docker容器化支持

#### 🔧 MVP技术实现
- **后端**: Go 1.21+, Gin框架, SQLite
- **前端**: React 18, TypeScript, Tailwind CSS
- **AI引擎**: DeepSeek API基础集成
- **部署**: Docker基础支持, Makefile构建

#### 📋 MVP项目特色
- **最小可用产品**: 核心功能完整可用
- **技术验证**: 验证核心技术可行性
- **用户体验**: 基础界面和交互体验
- **扩展性**: 为后续版本奠定基础

#### 📋 项目结构
```
digital-avatar/
├── backend/              # Go后端服务
├── frontend/             # React前端应用
├── docs/                 # 项目文档
│   ├── product/         # 产品文档
│   └── project/         # 项目管理文档
├── build/               # 构建输出
├── scripts/             # 构建脚本
├── .github/             # GitHub配置
├── Dockerfile           # 容器配置
├── docker-compose.yml   # 容器编排
├── Makefile            # 构建系统
├── README.md           # 项目说明
└── CONTRIBUTING.md     # 贡献指南
```

#### 🚀 快速开始
```bash
# 克隆项目
git clone https://github.com/your-username/digital-avatar.git
cd digital-avatar

# 初始化开发环境
make init

# 启动开发服务器
make dev
```

#### 📚 文档
- [产品需求文档](docs/product/prd.md)
- [产品设计文档](docs/product/product-design.md)
- [项目启动文档](docs/project/kickoff.md)
- [贡献指南](CONTRIBUTING.md)

#### 🎯 路线图
- **v1.1**: 移动端适配、更多AI模型支持
- **v2.0**: 插件生态系统、企业级功能
- **v3.0**: 云端部署版本、团队协作功能

---

## 📄 版本说明

- **主版本号**: 不兼容的API修改
- **次版本号**: 向下兼容的功能性新增
- **修订号**: 向下兼容的问题修正

---

## 🤝 贡献

欢迎提交Issue和Pull Request！请查看[贡献指南](CONTRIBUTING.md)了解详细信息。

---

*Digital Avatar - 你就负责思考，剩下的交给它*

**项目组成员: 大鸿 & 宫蕴*