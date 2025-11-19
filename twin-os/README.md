# TwinOS - 智能决策副驾系统

> **你就负责思考，剩下的交给它。**

## 项目概述

TwinOS 是面向高净值人群的"极致效率型数字分身"，通过静默分析微信数据，为用户提供结构化的决策情报。

## 核心功能

- 🔒 **数据采集**: 非侵入式读取微信数据库
- 🧠 **智能分析**: DeepSeek AI 驱动的情报分析
- 📊 **今日简报**: 300字核心信息摘要
- ✅ **决策待办**: 自动识别重要待办事项
- 👥 **人脉雷达**: 社交关系动态监控
- 🔐 **隐私保护**: 本地PII脱敏，数据永不外泄

## 技术架构

### 后端
- **语言**: Go 1.21+
- **框架**: Gin (Web Framework)
- **数据库**: SQLite (本地存储)
- **加密**: AES-256-GCM
- **AI集成**: DeepSeek API

### 前端
- **框架**: Electron + React 18
- **语言**: TypeScript
- **样式**: TailwindCSS
- **状态管理**: Zustand
- **动画**: Framer Motion

## 开发进度

- [x] 项目初始化
- [ ] 后端架构实现
- [ ] AI集成开发
- [ ] 前端界面开发
- [ ] 系统集成测试
- [ ] 产品发布准备

## 快速开始

### 环境要求
- Go 1.21+
- Node.js 18+
- npm/yarn

### 本地开发

```bash
# 后端开发
cd backend
go mod tidy
go run main.go

# 前端开发
cd frontend
npm install
npm run dev

# Electron 应用
npm run electron:dev
```

## 项目结构

```
twin-os/
├── backend/          # Go 后端服务
├── frontend/         # React 前端应用
├── docs/            # 项目文档
├── scripts/         # 构建和部署脚本
└── build/           # 构建产物
```

## 安全说明

- 原始聊天记录永不上传云端
- 本地PII脱敏处理
- 端到端加密保护
- 严格控制数据库访问频率

---

*开发团队：CEO & AI Partner*
*开发周期：12周*
*目标：打造最高效的个人决策助手*