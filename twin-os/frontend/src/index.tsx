import React, { useState, useEffect } from 'react';
import { createRoot } from 'react-dom/client';
import './index.css';

interface WeChatStatus {
  is_connected: boolean;
  db_path: string;
  contact_count: number;
  last_sync: string;
}

interface BriefingResult {
  content: string;
  message_count: number;
  generated_at: string;
}

const App: React.FC = () => {
  const [weChatStatus, setWeChatStatus] = useState<WeChatStatus | null>(null);
  const [briefing, setBriefing] = useState<BriefingResult | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    checkWeChatStatus();
    generateBriefing();
  }, []);

  const checkWeChatStatus = async () => {
    try {
      const response = await fetch('http://localhost:1234/api/v1/wechat/status');
      const data = await response.json();
      if (data.status === 'success') {
        setWeChatStatus(data.data);
      }
    } catch (err) {
      setError('连接后端服务失败');
    }
  };

  const generateBriefing = async () => {
    setIsLoading(true);
    try {
      const response = await fetch('http://localhost:1234/api/v1/analysis/briefing', {
        method: 'POST',
      });
      const data = await response.json();
      if (data.status === 'success') {
        setBriefing(data.data);
      }
    } catch (err) {
      setError('生成情报简报失败');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-black text-white">
      {/* 顶部导航栏 */}
      <div className="bg-gray-900 border-b border-gray-800 px-6 py-4">
        <div className="flex items-center justify-between">
          <h1 className="text-2xl font-bold text-blue-500">TWINOS</h1>
          <div className="flex items-center space-x-4">
            <div className="flex items-center space-x-2">
              <div className={`w-2 h-2 rounded-full ${weChatStatus?.is_connected ? 'bg-green-500' : 'bg-red-500'}`}></div>
              <span className="text-sm">
                {weChatStatus?.is_connected ? '微信已连接' : '微信未连接'}
              </span>
            </div>
          </div>
        </div>
      </div>

      {/* 主内容区域 */}
      <div className="p-6">
        {/* 系统状态 */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
          <div className="bg-gray-900 rounded-lg p-6 border border-gray-800">
            <h2 className="text-lg font-semibold mb-4">微信连接状态</h2>
            <div className="space-y-2">
              <div className="flex justify-between">
                <span className="text-gray-400">连接状态</span>
                <span className={weChatStatus?.is_connected ? 'text-green-400' : 'text-red-400'}>
                  {weChatStatus?.is_connected ? '已连接' : '未连接'}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-400">联系人数量</span>
                <span>{weChatStatus?.contact_count || 0}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-400">最后同步</span>
                <span className="text-sm">
                  {weChatStatus?.last_sync ? new Date(weChatStatus.last_sync).toLocaleString() : '未同步'}
                </span>
              </div>
            </div>
          </div>

          <div className="bg-gray-900 rounded-lg p-6 border border-gray-800">
            <h2 className="text-lg font-semibold mb-4">今日情报简报</h2>
            <div className="space-y-4">
              {isLoading ? (
                <div className="flex items-center justify-center py-8">
                  <div className="w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full animate-spin mr-3"></div>
                  <span>正在生成情报简报...</span>
                </div>
              ) : briefing ? (
                <div>
                  <div className="bg-gray-800 rounded p-4 mb-4">
                    <pre className="whitespace-pre-wrap text-sm text-gray-300">
                      {briefing.content}
                    </pre>
                  </div>
                  <div className="flex justify-between text-sm text-gray-400">
                    <span>消息数量: {briefing.message_count}</span>
                    <span>生成时间: {new Date(briefing.generated_at).toLocaleString()}</span>
                  </div>
                </div>
              ) : (
                <div className="text-center py-8 text-gray-400">
                  暂无情报简报
                </div>
              )}
            </div>
            <button
              onClick={generateBriefing}
              disabled={isLoading}
              className="mt-4 w-full px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-700 rounded-lg transition-colors"
            >
              {isLoading ? '生成中...' : '刷新简报'}
            </button>
          </div>
        </div>

        {/* 错误提示 */}
        {error && (
          <div className="bg-red-900 border border-red-800 text-red-200 px-4 py-3 rounded-lg mb-6">
            <div className="flex items-center justify-between">
              <span>{error}</span>
              <button
                onClick={() => setError(null)}
                className="text-red-200 hover:text-white"
              >
                ×
              </button>
            </div>
          </div>
        )}

        {/* 系统信息 */}
        <div className="bg-gray-900 rounded-lg p-6 border border-gray-800">
          <h2 className="text-lg font-semibold mb-4">系统信息</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
            <div>
              <span className="text-gray-400">前端服务</span>
              <span className="text-green-400">http://localhost:3000</span>
            </div>
            <div>
              <span className="text-gray-400">后端服务</span>
              <span className="text-green-400">http://localhost:1234</span>
            </div>
            <div>
              <span className="text-gray-400">团队模式</span>
              <span className="text-blue-400">两人敏捷团队</span>
            </div>
          </div>
        </div>

        {/* 快速操作 */}
        <div className="mt-6 flex space-x-4">
          <button
            onClick={checkWeChatStatus}
            className="px-6 py-2 bg-gray-800 hover:bg-gray-700 rounded-lg transition-colors"
          >
            检查状态
          </button>
          <button
            onClick={generateBriefing}
            disabled={isLoading}
            className="px-6 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-700 rounded-lg transition-colors"
          >
            {isLoading ? '分析中...' : '生成简报'}
          </button>
        </div>
      </div>
    </div>
  );
};

const container = document.getElementById('root');
if (!container) {
  throw new Error('Root container not found');
}

const root = createRoot(container);
root.render(<App />);