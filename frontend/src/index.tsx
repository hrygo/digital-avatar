import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import './styles/index.css';

// 错误边界组件
class ErrorBoundary extends React.Component<
  { children: React.ReactNode },
  { hasError: boolean; error: Error | null }
> {
  constructor(props: { children: React.ReactNode }) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error: Error) {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error('React Error Boundary caught an error:', error, errorInfo);

    // 这里可以添加错误报告逻辑
    // reportError(error, errorInfo);
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="min-h-screen bg-dark flex items-center justify-center p-6">
          <div className="glass-card max-w-md w-full p-6 text-center">
            <h1 className="text-2xl font-bold text-danger-400 mb-4">
              应用程序遇到错误
            </h1>
            <div className="text-left mb-6">
              <p className="text-gray-300 mb-2">错误信息:</p>
              <pre className="text-xs text-gray-500 bg-gray-900 p-3 rounded overflow-auto max-h-40">
                {this.state.error?.stack || this.state.error?.message || '未知错误'}
              </pre>
            </div>
            <button
              onClick={() => window.location.reload()}
              className="btn btn-primary"
            >
              重新加载应用
            </button>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}

// 获取根元素
const container = document.getElementById('root');
if (!container) {
  throw new Error('Root element not found');
}

// 创建React根实例
const root = ReactDOM.createRoot(container);

// 渲染应用
root.render(
  <React.StrictMode>
    <ErrorBoundary>
      <App />
    </ErrorBoundary>
  </React.StrictMode>
);