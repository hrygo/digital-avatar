# API文档中心

## 📚 API文档导航

Digital Avatar提供RESTful API和GraphQL接口，支持AI分析、PII处理、用户管理等核心功能。

## 🔗 基础信息

### API基础URL
- **开发环境**: `http://localhost:8080/api/v1`
- **测试环境**: `https://api-test.digital-avatar.com/api/v1`
- **生产环境**: `https://api.digital-avatar.com/api/v1`

### 认证方式
使用JWT Bearer Token进行认证：

```http
Authorization: Bearer <your_jwt_token>
```

### 响应格式
所有API响应使用统一的JSON格式：

```json
{
  "success": true,
  "data": {},
  "message": "操作成功",
  "timestamp": "2024-11-19T20:00:00Z",
  "requestId": "uuid-string"
}
```

## 📋 API分类

### 🔐 认证授权 ([auth/](endpoints/auth/))
用户认证、授权和权限管理。

- [用户登录](endpoints/auth/login.md) - 用户登录和获取token
- [用户注册](endpoints/auth/register.md) - 新用户注册
- [刷新Token](endpoints/auth/refresh.md) - JWT token刷新
- [密码重置](endpoints/auth/reset-password.md) - 密码重置流程

### 👥 用户管理 ([users/](endpoints/users/))
用户信息的CRUD操作。

- [获取用户列表](endpoints/users/list.md) - 分页获取用户列表
- [获取用户详情](endpoints/users/detail.md) - 获取指定用户信息
- [创建用户](endpoints/users/create.md) - 创建新用户
- [更新用户](endpoints/users/update.md) - 更新用户信息
- [删除用户](endpoints/users/delete.md) - 删除用户账号

### 🤖 AI功能 ([ai/](endpoints/ai/))
AI分析和智能处理功能。

- [文本分析](endpoints/ai/text-analysis.md) - 文本智能分析
- [情感分析](endpoints/ai/sentiment.md) - 情感倾向分析
- [实体识别](endpoints/ai/entity-recognition.md) - 命名实体识别
- [智能问答](endpoints/ai/qa.md) - 智能问答系统
- [决策建议](endpoints/ai/recommendation.md) - 智能决策建议

### 🔒 PII处理 ([pii/](endpoints/pii/))
个人敏感信息检测和处理。

- [PII检测](endpoints/pii/detect.md) - 检测文本中的敏感信息
- [PII替换](endpoints/pii/replace.md) - 敏感信息脱敏替换
- [批量处理](endpoints/pii/batch.md) - 批量PII处理
- [配置管理](endpoints/pii/config.md) - PII检测配置
- [统计报告](endpoints/pii/statistics.md) - PII处理统计

### 📊 数据分析 ([analytics/](endpoints/analytics/))
数据统计和分析接口。

- [用户统计](endpoints/analytics/users.md) - 用户数据分析
- [使用统计](endpoints/analytics/usage.md) - 功能使用统计
- [性能指标](endpoints/analytics/performance.md) - 系统性能指标
- [报表导出](endpoints/analytics/export.md) - 数据报表导出

## 🛠️ 快速开始

### 1. 获取访问Token

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

**响应**:
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expiresIn": 3600,
    "user": {
      "id": "uuid-string",
      "email": "user@example.com",
      "name": "张三",
      "role": "user"
    }
  }
}
```

### 2. 调用API

使用获取到的token调用其他API：

```http
GET /api/v1/users/profile
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 3. PII检测示例

```http
POST /api/v1/pii/detect
Authorization: Bearer <your-token>
Content-Type: application/json

{
  "text": "联系人：张三，手机：13812345678，邮箱：zhang@example.com",
  "threshold": 0.7
}
```

**响应**:
```json
{
  "success": true,
  "data": {
    "originalText": "联系人：张三，手机：13812345678，邮箱：zhang@example.com",
    "processedText": "联系人：[姓名1d841bc0]，手机：[手机号码]，邮箱：[邮箱地址]",
    "entities": [
      {
        "text": "张三",
        "type": "PER",
        "startPos": 3,
        "endPos": 5,
        "score": 0.8
      },
      {
        "text": "13812345678",
        "type": "TEL",
        "startPos": 8,
        "endPos": 17,
        "score": 0.95
      },
      {
        "text": "zhang@example.com",
        "type": "EMAIL",
        "startPos": 20,
        "endPos": 37,
        "score": 0.98
      }
    ],
    "statistics": {
      "totalEntities": 3,
      "typeDistribution": {
        "PER": 1,
        "TEL": 1,
        "EMAIL": 1
      }
    }
  }
}
```

## 📝 错误处理

### HTTP状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 401 | 未授权访问 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 429 | 请求频率限制 |
| 500 | 服务器内部错误 |

### 错误响应格式

```json
{
  "success": false,
  "error": {
    "code": "INVALID_PARAMETER",
    "message": "请求参数不合法",
    "details": {
      "field": "email",
      "reason": "邮箱格式不正确"
    }
  },
  "timestamp": "2024-11-19T20:00:00Z",
  "requestId": "uuid-string"
}
```

### 常见错误码

| 错误码 | 说明 | 解决方案 |
|--------|------|----------|
| `INVALID_TOKEN` | Token无效或过期 | 使用refresh token刷新 |
| `PERMISSION_DENIED` | 权限不足 | 联系管理员分配权限 |
| `RATE_LIMIT_EXCEEDED` | 请求频率超限 | 降低请求频率 |
| `PII_PROCESSING_FAILED` | PII处理失败 | 检查输入文本格式 |

## ⚡ 请求限制

### 频率限制

- **普通用户**: 1000次/小时
- **高级用户**: 5000次/小时
- **企业用户**: 10000次/小时

### 请求大小限制

- **JSON请求**: 最大10MB
- **文件上传**: 最大100MB
- **批量处理**: 最大1000条记录

## 🔧 开发工具

### Postman集合

提供完整的Postman API测试集合：

```json
{
  "info": {
    "name": "Digital Avatar API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "variable": [
    {
      "key": "baseUrl",
      "value": "http://localhost:8080/api/v1"
    },
    {
      "key": "token",
      "value": ""
    }
  ]
}
```

### SDK支持

#### JavaScript/TypeScript

```typescript
import { DigitalAvatarAPI } from '@digital-avatar/sdk';

const api = new DigitalAvatarAPI({
  baseURL: 'http://localhost:8080/api/v1',
  token: 'your-jwt-token'
});

// PII检测
const result = await api.pii.detect({
  text: '用户文本内容',
  threshold: 0.7
});

// 用户管理
const users = await api.users.list({
  page: 1,
  pageSize: 20
});
```

#### Go

```go
import (
    "github.com/hrygo/digital-avatar/go-sdk"
    "context"
)

client := digitalavatar.NewClient("http://localhost:8080/api/v1")
client.SetToken("your-jwt-token")

// PII检测
result, err := client.PIIDetect(context.Background(), &digitalavatar.PIIDetectRequest{
    Text:      "用户文本内容",
    Threshold: 0.7,
})
```

#### Python

```python
from digital_avatar_sdk import DigitalAvatarClient

client = DigitalAvatarClient(
    base_url='http://localhost:8080/api/v1',
    token='your-jwt-token'
)

# PII检测
result = client.pii.detect(
    text='用户文本内容',
    threshold=0.7
)
```

## 🔄 版本管理

### API版本策略

- **URL版本**: `/api/v1/`, `/api/v2/`
- **Header版本**: `API-Version: v1`
- **向后兼容**: 保持旧版本至少6个月

### 版本变更通知

- **重大变更**: 提前3个月通知
- **功能废弃**: 提前1个月通知
- **安全更新**: 立即生效

## 📊 监控和日志

### API监控

- **响应时间**: 平均响应时间 < 200ms
- **可用性**: 99.9% SLA
- **错误率**: < 0.1%

### 日志记录

所有API请求都会记录详细日志：

```json
{
  "timestamp": "2024-11-19T20:00:00Z",
  "method": "POST",
  "path": "/api/v1/pii/detect",
  "status": 200,
  "duration": 150,
  "requestId": "uuid-string",
  "userId": "uuid-string",
  "userAgent": "DigitalAvatar-SDK/1.0.0"
}
```

## 📞 技术支持

### API支持

- **文档问题**: 提交GitHub Issue
- **技术问题**: 发送邮件到 api-support@digital-avatar.com
- **紧急问题**: 联系技术支持热线

### 开发者社区

- **开发者论坛**: https://community.digital-avatar.com
- **Slack频道**: #digital-avatar-developers
- **技术博客**: https://blog.digital-avatar.com

---

*API文档最后更新: 2024年11月*