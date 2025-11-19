# Digital Avatar 前端文档

## 📚 文档导航

本文档是Digital Avatar前端系统的技术文档中心，提供组件库文档、样式指南、状态管理和部署指南。

## 🗂️ 文档分类

### 🎨 组件库
- [组件库概览](components/README.md) - UI组件库介绍和使用
- [基础组件](components/basic/) - Button、Input、Card等基础组件
- [布局组件](components/layout/) - Grid、Layout、Container等布局组件
- [业务组件](components/business/) - 与业务逻辑相关的高级组件

### 🎨 样式指南
- [设计系统](styles/design-system.md) - 色彩、字体、间距等设计规范
- [CSS架构](styles/css-architecture.md) - CSS模块化和组织架构
- [主题定制](styles/theming.md) - 主题切换和定制指南
- [响应式设计](styles/responsive.md) - 移动端适配和响应式布局

### 🔄 状态管理
- [状态管理架构](state-management/README.md) - Redux和状态管理架构
- [数据流](state-management/data-flow.md) - 数据流向和组件通信
- [异步处理](state-management/async.md) - API调用和异步状态处理
- [性能优化](state-management/performance.md) - 状态管理性能优化

### 🔧 构建部署
- [构建配置](build-deployment/build-config.md) - Webpack和构建配置
- [环境配置](build-deployment/environments.md) - 开发、测试、生产环境配置
- [部署指南](build-deployment/deployment.md) - 静态资源部署和CDN配置
- [性能优化](build-deployment/performance.md) - 前端性能优化策略

## 🚀 技术栈

### 核心框架
- **React 18**: 使用最新特性和Hooks
- **TypeScript**: 类型安全的JavaScript
- **Vite**: 快速的构建工具
- **React Router**: 客户端路由管理

### UI组件库
- **Ant Design**: 企业级UI设计语言
- **Tailwind CSS**: 实用优先的CSS框架
- **Styled Components**: CSS-in-JS解决方案
- **Framer Motion**: 动画和交互效果

### 状态管理
- **Redux Toolkit**: 现代化Redux状态管理
- **React Query**: 服务器状态管理
- **Zustand**: 轻量级状态管理
- **Context API**: 组件状态共享

### 开发工具
- **ESLint**: 代码质量检查
- **Prettier**: 代码格式化
- **Husky**: Git钩子管理
- **Jest**: 单元测试框架

## 🛠️ 开发指南

### 环境搭建
1. **克隆项目**
```bash
git clone https://github.com/digital-avatar/frontend.git
cd frontend
```

2. **安装依赖**
```bash
npm install
# 或
yarn install
```

3. **启动开发服务器**
```bash
npm run dev
# 或
yarn dev
```

4. **访问应用**
```
http://localhost:3000
```

### 项目结构
```
frontend/
├── public/                 # 静态资源
├── src/
│   ├── components/        # 通用组件
│   ├── pages/            # 页面组件
│   ├── hooks/            # 自定义Hooks
│   ├── services/         # API服务
│   ├── store/            # 状态管理
│   ├── utils/            # 工具函数
│   ├── types/            # TypeScript类型
│   └── styles/           # 全局样式
├── docs/                 # 项目文档
└── tests/                # 测试文件
```

### 开发规范
```typescript
// 组件示例
import React, { useState } from 'react';
import { Button, Card } from 'antd';
import styled from 'styled-components';

interface UserCardProps {
  user: User;
  onEdit: (user: User) => void;
}

const StyledCard = styled(Card)`
  margin-bottom: 16px;

  .ant-card-head {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }
`;

export const UserCard: React.FC<UserCardProps> = ({ user, onEdit }) => {
  const [loading, setLoading] = useState(false);

  const handleEdit = async () => {
    setLoading(true);
    try {
      await onEdit(user);
    } finally {
      setLoading(false);
    }
  };

  return (
    <StyledCard
      title={user.name}
      extra={
        <Button
          type="primary"
          loading={loading}
          onClick={handleEdit}
        >
          编辑
        </Button>
      }
    >
      <p>邮箱: {user.email}</p>
      <p>部门: {user.department}</p>
    </StyledCard>
  );
};
```

## 🎨 设计系统

### 色彩规范
```css
:root {
  /* 主色调 */
  --primary-color: #1890ff;
  --primary-hover: #40a9ff;
  --primary-active: #096dd9;

  /* 辅助色 */
  --success-color: #52c41a;
  --warning-color: #faad14;
  --error-color: #ff4d4f;
  --info-color: #1890ff;

  /* 中性色 */
  --text-primary: #262626;
  --text-secondary: #595959;
  --text-disabled: #bfbfbf;
  --border-color: #d9d9d9;
  --background-color: #fafafa;
}
```

### 字体规范
```css
:root {
  --font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  --font-size-xs: 12px;
  --font-size-sm: 14px;
  --font-size-base: 16px;
  --font-size-lg: 18px;
  --font-size-xl: 20px;
  --font-size-xxl: 24px;
}
```

### 间距规范
```css
:root {
  --spacing-xs: 4px;
  --spacing-sm: 8px;
  --spacing-md: 16px;
  --spacing-lg: 24px;
  --spacing-xl: 32px;
  --spacing-xxl: 48px;
}
```

## 🔄 状态管理

### Redux Toolkit示例
```typescript
// store/slices/userSlice.ts
import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { User } from '../types';
import { userApi } from '../services/userApi';

interface UserState {
  users: User[];
  current: User | null;
  loading: boolean;
  error: string | null;
}

const initialState: UserState = {
  users: [],
  current: null,
  loading: false,
  error: null,
};

export const fetchUsers = createAsyncThunk(
  'users/fetchUsers',
  async () => {
    const response = await userApi.getUsers();
    return response.data;
  }
);

const userSlice = createSlice({
  name: 'users',
  initialState,
  reducers: {
    clearError: (state) => {
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchUsers.pending, (state) => {
        state.loading = true;
      })
      .addCase(fetchUsers.fulfilled, (state, action) => {
        state.loading = false;
        state.users = action.payload;
      })
      .addCase(fetchUsers.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch users';
      });
  },
});

export const { clearError } = userSlice.actions;
export default userSlice.reducer;
```

## 📊 性能优化

### 代码分割
```typescript
// 路由级别的代码分割
import { lazy, Suspense } from 'react';
import { Spin } from 'antd';

const Dashboard = lazy(() => import('../pages/Dashboard'));
const Users = lazy(() => import('../pages/Users'));
const Settings = lazy(() => import('../pages/Settings'));

export const AppRoutes = () => {
  return (
    <Suspense fallback={<Spin size="large" />}>
      <Routes>
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/users" element={<Users />} />
        <Route path="/settings" element={<Settings />} />
      </Routes>
    </Suspense>
  );
};
```

### 组件优化
```typescript
import React, { memo, useMemo, useCallback } from 'react';

interface ExpensiveComponentProps {
  data: any[];
  onItemClick: (item: any) => void;
}

export const ExpensiveComponent = memo<ExpensiveComponentProps>(({
  data,
  onItemClick
}) => {
  // 使用useMemo缓存计算结果
  const processedData = useMemo(() => {
    return data.map(item => ({
      ...item,
      computed: expensiveCalculation(item),
    }));
  }, [data]);

  // 使用useCallback缓存事件处理函数
  const handleClick = useCallback((item: any) => {
    onItemClick(item);
  }, [onItemClick]);

  return (
    <div>
      {processedData.map(item => (
        <div
          key={item.id}
          onClick={() => handleClick(item)}
        >
          {item.computed}
        </div>
      ))}
    </div>
  );
});
```

## 🧪 测试

### 单元测试
```typescript
import { render, screen, fireEvent } from '@testing-library/react';
import { UserCard } from '../UserCard';

describe('UserCard', () => {
  const mockUser = {
    id: 1,
    name: 'John Doe',
    email: 'john@example.com',
    department: 'Engineering',
  };

  it('renders user information correctly', () => {
    render(<UserCard user={mockUser} onEdit={jest.fn()} />);

    expect(screen.getByText('John Doe')).toBeInTheDocument();
    expect(screen.getByText('john@example.com')).toBeInTheDocument();
    expect(screen.getByText('Engineering')).toBeInTheDocument();
  });

  it('calls onEdit when edit button is clicked', () => {
    const mockOnEdit = jest.fn();
    render(<UserCard user={mockUser} onEdit={mockOnEdit} />);

    fireEvent.click(screen.getByText('编辑'));
    expect(mockOnEdit).toHaveBeenCalledWith(mockUser);
  });
});
```

### 集成测试
```typescript
import { render, screen, waitFor } from '@testing-library/react';
import { Provider } from 'react-redux';
import { store } from '../store';
import { UserList } from '../UserList';

describe('UserList Integration', () => {
  it('loads and displays users', async () => {
    render(
      <Provider store={store}>
        <UserList />
      </Provider>
    );

    await waitFor(() => {
      expect(screen.getByText('John Doe')).toBeInTheDocument();
    });
  });
});
```

## 🚀 部署配置

### Vite配置
```typescript
// vite.config.ts
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import { resolve } from 'path';

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
      '@components': resolve(__dirname, 'src/components'),
      '@utils': resolve(__dirname, 'src/utils'),
    },
  },
  build: {
    outDir: 'dist',
    sourcemap: true,
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['react', 'react-dom'],
          antd: ['antd'],
          utils: ['lodash', 'moment'],
        },
      },
    },
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
});
```

## 📞 技术支持

### 开发团队
- **前端架构师**: 技术选型和架构设计
- **UI/UX设计师**: 组件设计和用户体验
- **前端开发**: 功能实现和性能优化
- **测试工程师**: 质量保证和自动化测试

### 问题反馈
- **GitHub Issues**: Bug报告和功能请求
- **技术讨论**: 前端技术交流
- **设计评审**: UI/UX设计讨论

---

## 📊 文档统计

- **组件数量**: 50+
- **测试覆盖率**: 目标 ≥ 85%
- **构建时间**: < 30秒
- **包体积**: < 1MB (gzipped)

*最后更新: 2024年11月*