import React, { useState, useEffect } from 'react';
import {
  Database,
  Upload,
  RefreshCw,
  Settings,
  MessageSquare,
  Users,
  Calendar,
  AlertCircle,
  CheckCircle,
  Plus,
  BarChart3,
  Eye,
  Download,
  Trash2
} from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import { wechatService, WeChatStatus, SyncRecord } from '../services/wechatService';
import WeChatImportWizard from './WeChatImportWizard';
import WeChatImportProgress from './WeChatImportProgress';

interface WeChatDataManagerProps {
  onStatusChange?: (status: WeChatStatus) => void;
}

const WeChatDataManager: React.FC<WeChatDataManagerProps> = ({ onStatusChange }) => {
  const [showImportWizard, setShowImportWizard] = useState(false);
  const [showProgress, setShowProgress] = useState(false);
  const [currentSyncId, setCurrentSyncId] = useState<number | null>(null);
  const [wechatStatus, setWeChatStatus] = useState<WeChatStatus | null>(null);
  const [syncRecords, setSyncRecords] = useState<SyncRecord[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string>('');

  // 获取微信状态
  const fetchWeChatStatus = async () => {
    try {
      const status = await wechatService.getStatus();
      setWeChatStatus(status);
      onStatusChange?.(status);
    } catch (error) {
      console.error('获取微信状态失败:', error);
      setError('无法连接到微信数据服务');
    }
  };

  // 获取同步记录
  const fetchSyncRecords = async () => {
    try {
      const response = await wechatService.getSyncRecords(1, 10);
      setSyncRecords(response.data);
    } catch (error) {
      console.error('获取同步记录失败:', error);
    }
  };

  // 初始化数据
  useEffect(() => {
    fetchWeChatStatus();
    fetchSyncRecords();

    // 定期刷新状态
    const interval = setInterval(() => {
      if (!showImportWizard && !showProgress) {
        fetchWeChatStatus();
      }
    }, 30000); // 30秒刷新一次

    return () => clearInterval(interval);
  }, [showImportWizard, showProgress]);

  // 处理导入完成
  const handleImportComplete = (response: { sync_id: number; status: string }) => {
    setCurrentSyncId(response.sync_id);
    setShowImportWizard(false);
    setShowProgress(true);
    fetchSyncRecords();
  };

  // 处理导入进度完成
  const handleProgressComplete = (record: SyncRecord) => {
    setShowProgress(false);
    setCurrentSyncId(null);
    fetchWeChatStatus();
    fetchSyncRecords();
  };

  // 手动同步
  const handleManualSync = async () => {
    if (!wechatStatus || !wechatStatus.connected) {
      setError('请先连接微信数据源');
      return;
    }

    setIsLoading(true);
    setError('');

    try {
      const response = await wechatService.sync({
        type: 'database',
        path: '', // 使用默认路径
        mode: 'incremental'
      });

      setCurrentSyncId(response.sync_id);
      setShowProgress(true);
    } catch (error) {
      console.error('手动同步失败:', error);
      setError(error instanceof Error ? error.message : '同步失败');
    } finally {
      setIsLoading(false);
    }
  };

  // 格式化时间
  const formatTime = (timestamp: number) => {
    return new Date(timestamp).toLocaleString('zh-CN');
  };

  // 格式化持续时间
  const formatDuration = (ms: number) => {
    if (ms < 1000) return `${ms}ms`;
    if (ms < 60000) return `${Math.round(ms / 1000)}s`;
    return `${Math.round(ms / 60000)}min`;
  };

  const StatusIndicator = ({ status }: { status: string }) => {
    const config = {
      connected: { color: 'text-green-400', bg: 'bg-green-600/20', border: 'border-green-600/30', icon: CheckCircle, text: '已连接' },
      disconnected: { color: 'text-red-400', bg: 'bg-red-600/20', border: 'border-red-600/30', icon: AlertCircle, text: '未连接' },
      syncing: { color: 'text-blue-400', bg: 'bg-blue-600/20', border: 'border-blue-600/30', icon: RefreshCw, text: '同步中' }
    };

    const configItem = config[status as keyof typeof config] || config.disconnected;
    const Icon = configItem.icon;

    return (
      <div className={`flex items-center gap-2 px-3 py-1 rounded-full border ${configItem.bg} ${configItem.border} ${configItem.color}`}>
        <Icon className="w-4 h-4" />
        <span className="text-sm font-medium">{configItem.text}</span>
      </div>
    );
  };

  if (showImportWizard) {
    return (
      <WeChatImportWizard
        onComplete={handleImportComplete}
        onCancel={() => setShowImportWizard(false)}
      />
    );
  }

  return (
    <div className="space-y-6">
      {/* 微信数据管理组件 */}
      <div className="bg-gray-800/50 backdrop-blur-sm rounded-xl p-6 border border-gray-700">
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-green-600 rounded-lg flex items-center justify-center">
              <MessageSquare className="w-5 h-5 text-white" />
            </div>
            <div>
              <h2 className="text-xl font-bold text-white">微信数据管理</h2>
              <p className="text-gray-400 text-sm">管理和同步微信聊天数据</p>
            </div>
          </div>

          <div className="flex items-center gap-3">
            {wechatStatus && <StatusIndicator status={wechatStatus.connected ? 'connected' : 'disconnected'} />}
            <button
              onClick={() => setShowImportWizard(true)}
              className="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
            >
              <Plus className="w-4 h-4" />
              导入数据
            </button>
            <button
              onClick={handleManualSync}
              disabled={isLoading || !wechatStatus?.connected}
              className="flex items-center gap-2 px-4 py-2 bg-green-600 hover:bg-green-700 disabled:bg-gray-700 disabled:text-gray-400 text-white rounded-lg transition-colors"
            >
              <RefreshCw className={`w-4 h-4 ${isLoading ? 'animate-spin' : ''}`} />
              {isLoading ? '同步中...' : '手动同步'}
            </button>
          </div>
        </div>

        {/* 错误提示 */}
        <AnimatePresence>
          {error && (
            <motion.div
              initial={{ opacity: 0, y: -10 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -10 }}
              className="mb-4 bg-red-900/20 border border-red-700 rounded-lg p-4"
            >
              <div className="flex items-start gap-3">
                <AlertCircle className="w-5 h-5 text-red-400 flex-shrink-0 mt-0.5" />
                <div>
                  <h4 className="font-semibold text-red-400 mb-1">错误</h4>
                  <p className="text-red-300 text-sm">{error}</p>
                </div>
                <button
                  onClick={() => setError('')}
                  className="ml-auto text-red-400 hover:text-red-300"
                >
                  ×
                </button>
              </div>
            </motion.div>
          )}
        </AnimatePresence>

        {/* 数据统计 */}
        {wechatStatus && wechatStatus.connected && (
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
            <div className="bg-gray-900/50 rounded-lg p-4 border border-gray-600/50">
              <div className="flex items-center gap-3 mb-2">
                <MessageSquare className="w-5 h-5 text-blue-400" />
                <span className="text-gray-400 text-sm">总消息数</span>
              </div>
              <div className="text-2xl font-bold text-white">
                {wechatStatus.messageCount.toLocaleString()}
              </div>
            </div>

            <div className="bg-gray-900/50 rounded-lg p-4 border border-gray-600/50">
              <div className="flex items-center gap-3 mb-2">
                <Users className="w-5 h-5 text-green-400" />
                <span className="text-gray-400 text-sm">联系人</span>
              </div>
              <div className="text-2xl font-bold text-white">
                {wechatStatus.contactCount.toLocaleString()}
              </div>
            </div>

            <div className="bg-gray-900/50 rounded-lg p-4 border border-gray-600/50">
              <div className="flex items-center gap-3 mb-2">
                <Database className="w-5 h-5 text-purple-400" />
                <span className="text-gray-400 text-sm">聊天会话</span>
              </div>
              <div className="text-2xl font-bold text-white">
                {wechatStatus.chatCount.toLocaleString()}
              </div>
            </div>

            <div className="bg-gray-900/50 rounded-lg p-4 border border-gray-600/50">
              <div className="flex items-center gap-3 mb-2">
                <Calendar className="w-5 h-5 text-yellow-400" />
                <span className="text-gray-400 text-sm">最后同步</span>
              </div>
              <div className="text-sm font-medium text-white">
                {formatTime(wechatStatus.lastSyncTime)}
              </div>
            </div>
          </div>
        )}

        {/* 同步记录 */}
        <div>
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold text-white">同步历史</h3>
            <button
              onClick={fetchSyncRecords}
              className="text-gray-400 hover:text-white transition-colors"
            >
              <RefreshCw className="w-4 h-4" />
            </button>
          </div>

          {syncRecords.length > 0 ? (
            <div className="space-y-3">
              {syncRecords.map((record, index) => (
                <motion.div
                  key={record.id}
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: index * 0.05 }}
                  className="bg-gray-900/50 rounded-lg p-4 border border-gray-700/50"
                >
                  <div className="flex items-center justify-between">
                    <div className="flex-1">
                      <div className="flex items-center gap-3 mb-2">
                        <span className="text-white font-medium">
                          {record.sync_type === 'full' ? '全量同步' : '增量同步'}
                        </span>
                        <span className="text-gray-400 text-sm">
                          来源: {record.source}
                        </span>
                        {record.status === 'completed' ? (
                          <CheckCircle className="w-4 h-4 text-green-400" />
                        ) : record.status === 'failed' ? (
                          <AlertCircle className="w-4 h-4 text-red-400" />
                        ) : (
                          <RefreshCw className="w-4 h-4 text-blue-400 animate-spin" />
                        )}
                      </div>

                      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                        <div>
                          <span className="text-gray-400">消息: </span>
                          <span className="text-white">
                            {record.new_messages > 0 && `+${record.new_messages}/`}
                            {record.total_messages}
                          </span>
                        </div>
                        <div>
                          <span className="text-gray-400">联系人: </span>
                          <span className="text-white">
                            {record.new_contacts > 0 && `+${record.new_contacts}/`}
                            {record.total_contacts}
                          </span>
                        </div>
                        <div>
                          <span className="text-gray-400">进度: </span>
                          <span className="text-white">{record.progress}%</span>
                        </div>
                        <div>
                          <span className="text-gray-400">用时: </span>
                          <span className="text-white">
                            {record.duration ? formatDuration(record.duration) : '进行中...'}
                          </span>
                        </div>
                      </div>

                      {record.error_message && (
                        <div className="mt-2 text-xs text-red-400 bg-red-900/20 rounded p-2">
                          {record.error_message}
                        </div>
                      )}
                    </div>

                    <div className="ml-4 text-right">
                      <div className="text-xs text-gray-500 mb-1">
                        {formatTime(new Date(record.created_at).getTime())}
                      </div>
                      {record.status === 'completed' && (
                        <button className="text-blue-400 hover:text-blue-300 text-sm">
                          <Eye className="w-4 h-4" />
                        </button>
                      )}
                    </div>
                  </div>
                </motion.div>
              ))}
            </div>
          ) : (
            <div className="text-center py-8 text-gray-400">
              <Database className="w-12 h-12 mx-auto mb-3 opacity-50" />
              <p>暂无同步记录</p>
              <button
                onClick={() => setShowImportWizard(true)}
                className="mt-4 text-blue-400 hover:text-blue-300 text-sm"
              >
                开始导入数据
              </button>
            </div>
          )}
        </div>
      </div>

      {/* 进度监控 */}
      <WeChatImportProgress
        syncId={currentSyncId}
        isVisible={showProgress}
        onClose={() => {
          setShowProgress(false);
          setCurrentSyncId(null);
        }}
        onComplete={handleProgressComplete}
        onError={(error) => {
          setError(error);
          setShowProgress(false);
          setCurrentSyncId(null);
        }}
      />
    </div>
  );
};

export default WeChatDataManager;