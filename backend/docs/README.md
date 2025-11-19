# 项目文档

## 技术报告

### PII检测系统相关

- [测试覆盖率报告](COVERAGE_REPORT.md) - 详细的单元测试覆盖率分析
- [覆盖率摘要](COVERAGE_SUMMARY.md) - 测试覆盖率的概览和统计
- [边界情况优化报告](BOUNDARY_OPTIMIZATION_REPORT.md) - PII检测边界情况处理优化
- [代码清理报告](CLEANUP_REPORT.md) - 代码重构和清理过程记录

### API文档

- [API接口文档](api.md) - RESTful API接口详细说明
- [数据模型文档](models.md) - 数据结构和实体定义

### 开发文档

- [开发环境搭建](development.md) - 本地开发环境配置指南
- [部署指南](deployment.md) - 生产环境部署说明
- [性能优化指南](performance.md) - 系统性能调优建议

## 系统架构

### 核心模块

- **pkg/crypto** - PII检测和数据加密模块
- **pkg/ai** - AI增强分析功能
- **internal/** - 内部业务逻辑
- **cmd/** - 应用程序入口

### 设计文档

- [架构设计](architecture.md) - 系统整体架构说明
- [数据库设计](database.md) - 数据模型和关系设计
- [安全设计](security.md) - 数据安全和隐私保护机制

## 运维文档

### 监控和日志

- [监控配置](monitoring.md) - 系统监控和告警设置
- [日志管理](logging.md) - 日志收集和分析指南

### 备份和恢复

- [备份策略](backup.md) - 数据备份和恢复流程
- [灾难恢复](disaster-recovery.md) - 系统故障恢复方案

## 版本管理

- [版本发布说明](releases.md) - 各版本的详细更新记录
- [升级指南](migration.md) - 版本升级和兼容性说明