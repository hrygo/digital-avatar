import React from 'react';
import Dashboard from './pages/Dashboard';
import { ErrorBoundary } from './components/ErrorBoundary';

const App: React.FC = () => {
  return (
    <ErrorBoundary
      onError={(error, errorInfo) => {
        // 记录应用级错误
        console.error('App level error:', error);
        console.error('Component stack:', errorInfo.componentStack);
      }}
    >
      <div className="App">
        <Dashboard />
      </div>
    </ErrorBoundary>
  );
};

export default App;