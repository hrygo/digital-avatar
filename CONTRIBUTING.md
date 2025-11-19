# 贡献指南

感谢您对 TwinOS 项目的关注！我们欢迎所有形式的贡献，包括但不限于代码、文档、测试、反馈和建议。

## 🤝 贡献方式

### 报告问题
- 使用 [GitHub Issues](https://github.com/hrygo/digital-avatar/issues) 报告 bug
- 在报告中提供详细的重现步骤和环境信息
- 使用已有的问题模板来确保信息完整

### 功能建议
- 在 [GitHub Discussions](https://github.com/your-username/twin-os/discussions) 中讨论新功能想法
- 提交 Issue 时详细描述功能需求和使用场景
- 考虑功能的向后兼容性和影响范围

### 代码贡献
- Fork 项目到您的 GitHub 账户
- 创建功能分支进行开发
- 遵循项目的代码规范和测试要求
- 提交 Pull Request 并详细描述更改内容

### 文档改进
- 改进现有文档的准确性和可读性
- 添加缺失的 API 文档或使用示例
- 翻译文档到其他语言

## 🚀 开始贡献

### 环境设置

1. **Fork 项目**
   ```bash
   # 在 GitHub 上 Fork 项目
   # 然后克隆您的 Fork
   git clone https://github.com/YOUR_USERNAME/twin-os.git
   cd twin-os
   ```

2. **添加上游仓库**
   ```bash
   git remote add upstream https://github.com/original-owner/twin-os.git
   git fetch upstream
   ```

3. **创建开发分支**
   ```bash
   git checkout -b feature/your-feature-name
   ```

4. **设置开发环境**
   ```bash
   # 初始化开发环境
   make init

   # 启动开发服务器
   make dev
   ```

### 开发流程

1. **代码规范**
   ```bash
   # 格式化代码
   make format

   # 代码检查
   make lint

   # 运行测试
   make test
   ```

2. **提交代码**
   ```bash
   # 添加更改
   git add .

   # 提交（使用有意义的提交信息）
   git commit -m "feat: add new analysis feature"
   ```

3. **保持同步**
   ```bash
   # 同步上游更改
   git fetch upstream
   git rebase upstream/main
   ```

4. **提交 Pull Request**
   - 推送到您的 Fork
   - 在 GitHub 上创建 Pull Request
   - 填写 PR 模板
   - 等待代码审查

## 📝 提交规范

我们使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

### 提交类型

- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更改
- `style`: 代码格式化（不影响功能）
- `refactor`: 代码重构
- `test`: 添加或修改测试
- `chore`: 构建过程或辅助工具的变动

### 提交格式

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### 示例

```bash
# 新功能
git commit -m "feat(ai): add sentiment analysis capability"

# Bug 修复
git commit -m "fix(backup): resolve file permission issue on Linux"

# 文档更新
git commit -m "docs: update API documentation for v1.1"
```

## 🧪 测试要求

### 测试覆盖率

- 新功能必须包含测试
- 保持测试覆盖率在 80% 以上
- 所有测试必须通过 CI 检查

### 运行测试

```bash
# 运行所有测试
make test

# 运行特定测试
make test-server  # 后端测试
make test-web      # 前端测试

# 生成覆盖率报告
make test coverage
```

### 测试编写指南

#### 后端测试

```go
// 示例：backend/pkg/analysis/analysis_test.go
func TestAnalyzeMessage(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "simple message",
            input:    "Hello world",
            expected: "analyzed: Hello world",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := AnalyzeMessage(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

#### 前端测试

```typescript
// 示例：frontend/src/components/__tests__/Dashboard.test.tsx
import { render, screen } from '@testing-library/react';
import Dashboard from '../Dashboard';

describe('Dashboard Component', () => {
  it('renders title correctly', () => {
    render(<Dashboard />);
    expect(screen.getByText('TwinOS Dashboard')).toBeInTheDocument();
  });
});
```

## 🎨 代码规范

### Go 代码规范

- 遵循 [Go 官方代码规范](https://golang.org/doc/effective_go.html)
- 使用 `gofmt` 和 `goimports` 格式化代码
- 包名使用小写字母
- 接口名以 -er 结尾
- 错误处理要完整

### TypeScript/React 代码规范

- 使用 TypeScript 进行类型检查
- 组件使用函数式组件和 Hooks
- 遵循 [React 最佳实践](https://react.dev/learn)
- 使用 Prettier 格式化代码

### 命名规范

#### Go
```go
// 常量
const MaxRetries = 3

// 变量
var userService *UserService

// 函数
func GetUserByID(id string) (*User, error)

// 结构体
type User struct {
    ID   string
    Name string
}
```

#### TypeScript
```typescript
// 常量
const MAX_RETRIES = 3;

// 变量
const userService = new UserService();

// 函数
function getUserById(id: string): Promise<User> {
  // implementation
}

// 接口
interface User {
  id: string;
  name: string;
}
```

## 📚 文档贡献

### API 文档

- 为所有公共 API 添加文档注释
- 使用 Swagger/OpenAPI 规范
- 提供请求/响应示例

### 代码注释

```go
// GetUserByID 根据用户ID获取用户信息
//
// 参数:
//   - id: 用户唯一标识符
//
// 返回:
//   - *User: 用户信息指针
//   - error: 错误信息，如果用户不存在则返回 ErrUserNotFound
//
// 示例:
//   user, err := GetUserByID("123")
//   if err != nil {
//       log.Fatal(err)
//   }
//   fmt.Printf("User: %+v\n", user)
func GetUserByID(id string) (*User, error) {
    // implementation
}
```

## 🔍 代码审查

### 审查清单

- [ ] 代码符合项目规范
- [ ] 包含适当的测试
- [ ] 文档已更新
- [ ] 性能影响可接受
- [ ] 安全性已考虑
- [ ] 向后兼容性

### 审查流程

1. **自动化检查**: CI/CD 管道运行测试和代码检查
2. **人工审查**: 维护者进行代码审查
3. **反馈处理**: 根据反馈修改代码
4. **合并代码**: 审查通过后合并到主分支

## 🏷️ 发布流程

### 版本号规范

我们使用 [语义化版本](https://semver.org/lang/zh-CN/)：

- `MAJOR.MINOR.PATCH`
- `MAJOR`: 不兼容的 API 更改
- `MINOR`: 向后兼容的新功能
- `PATCH`: 向后兼容的 Bug 修复

### 发布步骤

1. 更新版本号
2. 更新 CHANGELOG.md
3. 创建 Git 标签
4. 构建 release 包
5. 发布到 GitHub Releases

## 💬 社区

### 沟通渠道

- **GitHub Issues**: Bug 报告和功能请求
- **GitHub Discussions**: 一般讨论和问答
- **邮件列表**: 重要公告和讨论
- **社交媒体**: 项目动态和社区互动

### 行为准则

- 尊重所有参与者
- 保持友好和专业
- 接受建设性反馈
- 关注问题本身而非个人

## 🙏 致谢

感谢所有为 TwinOS 项目做出贡献的开发者！

### 贡献者列表

- [@your-username](https://github.com/your-username) - 项目创始人
- 添加更多贡献者...

---

## 📞 联系我们

如果您有任何问题或建议，欢迎通过以下方式联系我们：

- 邮箱: dev@twin-os.com
- GitHub: [TwinOS Organization](https://github.com/twin-os)
- 官网: https://twin-os.com

---

再次感谢您的贡献！🎉