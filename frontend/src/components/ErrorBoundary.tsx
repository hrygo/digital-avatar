import React, { Component, ErrorInfo, ReactNode } from 'react';
import { motion } from 'framer-motion';
import { AlertTriangle, RefreshCw, Home } from 'lucide-react';

interface Props {
  children: ReactNode;
  fallback?: ReactNode;
  onError?: (error: Error, errorInfo: ErrorInfo) => void;
}

interface State {
  hasError: boolean;
  error: Error | null;
  errorInfo: ErrorInfo | null;
  retryCount: number;
}

export class ErrorBoundary extends Component<Props, State> {
  private maxRetries = 3;
  private retryTimeout: NodeJS.Timeout | null = null;

  constructor(props: Props) {
    super(props);
    this.state = {
      hasError: false,
      error: null,
      errorInfo: null,
      retryCount: 0
    };
  }

  static getDerivedStateFromError(error: Error): Partial<State> {
    return {
      hasError: true,
      error
    };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    this.setState({
      error,
      errorInfo
    });

    // 记录错误到控制台
    console.error('ErrorBoundary caught an error:', error, errorInfo);

    // 调用自定义错误处理函数
    if (this.props.onError) {
      this.props.onError(error, errorInfo);
    }

    // 发送错误报告到监控服务（生产环境）
    if (process.env.NODE_ENV === 'production') {
      this.reportError(error, errorInfo);
    }
  }

  componentWillUnmount() {
    if (this.retryTimeout) {
      clearTimeout(this.retryTimeout);
    }
  }

  private reportError = (error: Error, errorInfo: ErrorInfo) => {
    // 这里可以集成错误监控服务，如 Sentry
    const errorReport = {
      message: error.message,
      stack: error.stack,
      componentStack: errorInfo.componentStack,
      timestamp: new Date().toISOString(),
      userAgent: navigator.userAgent,
      url: window.location.href
    };

    // 发送到错误监控API
    fetch('/api/v1/error/report', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(errorReport)
    }).catch(console.error);
  };

  private handleRetry = () => {
    if (this.state.retryCount < this.maxRetries) {
      this.setState(prevState => ({
        hasError: false,
        error: null,
        errorInfo: null,
        retryCount: prevState.retryCount + 1
      }));

      // 延迟重试以避免立即重复错误
      this.retryTimeout = setTimeout(() => {
        window.location.reload();
      }, 1000);
    }
  };

  private handleGoHome = () => {
    window.location.href = '/';
  };

  private handleRefresh = () => {
    window.location.reload();
  };

  render() {
    if (this.state.hasError) {
      // 如果提供了自定义fallback，使用它
      if (this.props.fallback) {
        return this.props.fallback;
      }

      // 默认错误UI
      return (
        <div className="min-h-screen bg-black text-white flex items-center justify-center p-4">
          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            className="glass-card rounded-2xl p-8 max-w-md w-full text-center"
          >
            {/* 错误图标 */}
            <motion.div
              initial={{ rotate: 0 }}
              animate={{ rotate: [0, 10, -10, 0] }}
              transition={{ repeat: Infinity, duration: 4 }}
              className="w-20 h-20 bg-gradient-to-br from-red-500/20 to-pink-500/20 rounded-full flex items-center justify-center mx-auto mb-6"
            >
              <AlertTriangle className="w-10 h-10 text-red-400" />
            </motion.div>

            {/* 错误标题 */}
            <h1 className="text-2xl font-bold text-white mb-4">
              哎呀，出现了些问题
            </h1>

            {/* 错误消息 */}
            <p className="text-gray-300 mb-6 leading-relaxed">
              {process.env.NODE_ENV === 'development'
                ? this.state.error?.message || '未知错误'
                : '系统遇到了意外错误，请稍后重试'
              }
            </p>

            {/* 开发环境错误详情 */}
            {process.env.NODE_ENV === 'development' && this.state.error && (
              <details className="mb-6 text-left">
                <summary className="cursor-pointer text-blue-400 hover:text-blue-300 mb-2">
                  查看错误详情
                </summary>
                <div className="mt-2 p-4 bg-gray-900/50 rounded-lg border border-gray-700">
                  <pre className="text-xs text-gray-400 overflow-x-auto whitespace-pre-wrap">
                    {this.state.error.stack}
                    {this.state.errorInfo?.componentStack}
                  </pre>
                </div>
              </details>
            )}

            {/* 重试次数提示 */}
            {this.state.retryCount >= this.maxRetries && (
              <div className="mb-4 p-3 bg-yellow-500/20 border border-yellow-500/50 rounded-lg">
                <p className="text-sm text-yellow-300">
                  已达到最大重试次数 ({this.maxRetries})
                </p>
              </div>
            )}

            {/* 操作按钮 */}
            <div className="flex flex-col sm:flex-row gap-3 justify-center">
              {this.state.retryCount < this.maxRetries && (
                <motion.button
                  whileHover={{ scale: 1.05 }}
                  whileTap={{ scale: 0.95 }}
                  onClick={this.handleRetry}
                  className="btn-primary flex items-center justify-center space-x-2"
                >
                  <RefreshCw className="w-4 h-4" />
                  <span>重试 ({this.maxRetries - this.state.retryCount})</span>
                </motion.button>
              )}

              <motion.button
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
                onClick={this.handleRefresh}
                className="btn-secondary flex items-center justify-center space-x-2"
              >
                <RefreshCw className="w-4 h-4" />
                <span>刷新页面</span>
              </motion.button>

              <motion.button
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
                onClick={this.handleGoHome}
                className="btn-secondary flex items-center justify-center space-x-2"
              >
                <Home className="w-4 h-4" />
                <span>返回首页</span>
              </motion.button>
            </div>

            {/* 帮助信息 */}
            <div className="mt-8 pt-6 border-t border-gray-700">
              <p className="text-sm text-gray-400">
                如果问题持续存在，请联系技术支持
              </p>
            </div>
          </motion.div>
        </div>
      );
    }

    return this.props.children;
  }
}

// 简化的函数式错误边界组件
export const SimpleErrorBoundary: React.FC<{
  children: ReactNode;
  message?: string;
  onRetry?: () => void;
}> = ({ children, message = '出现错误', onRetry }) => {
  return (
    <ErrorBoundary
      fallback={
        <div className="p-4 text-center">
          <AlertTriangle className="w-8 h-8 text-red-400 mx-auto mb-2" />
          <p className="text-red-400 mb-2">{message}</p>
          {onRetry && (
            <button
              onClick={onRetry}
              className="px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-lg text-white"
            >
              重试
            </button>
          )}
        </div>
      }
    >
      {children}
    </ErrorBoundary>
  );
};

// HOC包装器
export const withErrorBoundary = <P extends object>(
  Component: React.ComponentType<P>,
  errorFallback?: ReactNode,
  onError?: (error: Error, errorInfo: ErrorInfo) => void
) => {
  const WrappedComponent = (props: P) => (
    <ErrorBoundary fallback={errorFallback} onError={onError}>
      <Component {...props} />
    </ErrorBoundary>
  );

  WrappedComponent.displayName = `withErrorBoundary(${Component.displayName || Component.name})`;

  return WrappedComponent;
};

export default ErrorBoundary;