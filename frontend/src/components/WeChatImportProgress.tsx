import React, { useState, useEffect, useCallback } from 'react';
import {
  Activity,
  Users,
  MessageSquare,
  Clock,
  CheckCircle,
  AlertCircle,
  Loader2,
  X,
  RefreshCw,
  Eye,
  FileText,
  Database
} from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import { wechatService, SyncRecord } from '../services/wechatService';

interface WeChatImportProgressProps {
  syncId: number | null;
  isVisible: boolean;
  onClose: () => void;
  onComplete?: (record: SyncRecord) => void;
  onError?: (error: string) => void;
}

interface ProgressStep {
  id: string;
  title: string;
  description: string;
  status: 'pending' | 'running' | 'completed' | 'error';
  progress: number;
  icon: React.ReactNode;
  details?: string;
}

const WeChatImportProgress: React.FC<WeChatImportProgressProps> = ({
  syncId,
  isVisible,
  onClose,
  onComplete,
  onError
}) => {
  const [currentRecord, setCurrentRecord] = useState<SyncRecord | null>(null);
  const [steps, setSteps] = useState<ProgressStep[]>([]);
  const [isPolling, setIsPolling] = useState(false);
  const [showDetails, setShowDetails] = useState(false);
  const [error, setError] = useState<string>('');

  // 初始化进度步骤
  const initializeSteps = useCallback(() => {
    const initialSteps: ProgressStep[] = [
      {
        id: 'connect',
        title: '连接数据源',
        description: '建立与微信数据库的连接',
        status: 'pending',
        progress: 0,
        icon: <Database className="w-4 h-4" />
      },
      {
        id: 'validate',
        title: '验证数据完整性',
        description: '检查数据文件的有效性和完整性',
        status: 'pending',
        progress: 0,
        icon: <CheckCircle className="w-4 h-4" />
      },
      {
        id: 'parse-contacts',
        title: '解析联系人信息',
        description: '提取和处理联系人数据',
        status: 'pending',
        progress: 0,
        icon: <Users className="w-4 h-4" />
      },
      {
        id: 'parse-messages',
        title: '解析消息记录',
        description: '处理聊天消息和媒体文件',
        status: 'pending',
        progress: 0,
        icon: <MessageSquare className="w-4 h-4" />
      },
      {
        id: 'analyze-content',
        title: '分析内容结构',
        description: '分析聊天内容和关系网络',
        status: 'pending',
        progress: 0,
        icon: <Activity className="w-4 h-4" />
      },
      {
        id: 'save-data',
        title: '保存到数据库',
        description: '将处理后的数据保存到系统',
        status: 'pending',
        progress: 0,
        icon: <FileText className="w-4 h-4" />
      }
    ];

    setSteps(initialSteps);
  }, []);

  // 更新步骤状态
  const updateStepsFromProgress = useCallback((record: SyncRecord) => {
    const progress = record.progress || 0;
    const status = record.status.toLowerCase();

    const newSteps = [...steps];

    // 根据进度和状态更新步骤
    if (status === 'running') {
      const completedSteps = Math.floor((progress / 100) * newSteps.length);

      newSteps.forEach((step, index) => {
        if (index < completedSteps) {
          step.status = 'completed';
          step.progress = 100;
        } else if (index === completedSteps) {
          step.status = 'running';
          step.progress = (progress % (100 / newSteps.length)) * newSteps.length;
        } else {
          step.status = 'pending';
          step.progress = 0;
        }
      });
    } else if (status === 'completed') {
      newSteps.forEach(step => {
        step.status = 'completed';
        step.progress = 100;
      });
    } else if (status === 'failed') {
      // 找到失败的步骤
      const failedStepIndex = Math.floor((progress / 100) * newSteps.length);
      if (failedStepIndex < newSteps.length) {
        newSteps[failedStepIndex].status = 'error';
        newSteps[failedStepIndex].details = record.error_message;
      }
    }

    setSteps(newSteps);
  }, [steps]);

  // 获取同步记录
  const fetchSyncRecord = useCallback(async () => {
    if (!syncId) return;

    try {
      const response = await wechatService.getSyncRecords(1, 1);
      const record = response.data.find(r => r.id === syncId);

      if (record) {
        setCurrentRecord(record);
        updateStepsFromProgress(record);

        if (record.status === 'completed') {
          setIsPolling(false);
          onComplete?.(record);
        } else if (record.status === 'failed') {
          setIsPolling(false);
          setError(record.error_message || '导入失败');
          onError?.(record.error_message || '导入失败');
        }
      }
    } catch (error) {
      console.error('获取同步记录失败:', error);
      setError('获取进度信息失败');
      setIsPolling(false);
    }
  }, [syncId, updateStepsFromProgress, onComplete, onError]);

  // 开始轮询
  useEffect(() => {
    if (isVisible && syncId && !isPolling) {
      setIsPolling(true);
      initializeSteps();
      fetchSyncRecord();
    }
  }, [isVisible, syncId, isPolling, initializeSteps, fetchSyncRecord]);

  // 定时轮询进度
  useEffect(() => {
    if (!isPolling || !syncId) return;

    const interval = setInterval(() => {
      fetchSyncRecord();
    }, 2000);

    return () => clearInterval(interval);
  }, [isPolling, syncId, fetchSyncRecord]);

  // 重试
  const handleRetry = () => {
    setError('');
    setIsPolling(true);
    fetchSyncRecord();
  };

  // 计算总体进度
  const overallProgress = steps.reduce((acc, step) => acc + step.progress, 0) / steps.length;

  if (!isVisible) return null;

  return (
    <AnimatePresence>
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        exit={{ opacity: 0 }}
        className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center p-6 z-50"
      >
        <motion.div
          initial={{ scale: 0.9, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          exit={{ scale: 0.9, opacity: 0 }}
          transition={{ duration: 0.2 }}
          className="bg-gray-900 rounded-2xl border border-gray-700 w-full max-w-2xl max-h-[90vh] overflow-hidden"
        >
          {/* Header */}
          <div className="flex items-center justify-between p-6 border-b border-gray-800">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 bg-blue-600 rounded-lg flex items-center justify-center">
                <Activity className="w-5 h-5 text-white" />
              </div>
              <div>
                <h2 className="text-xl font-bold text-white">微信数据导入</h2>
                <p className="text-gray-400 text-sm">
                  {currentRecord ? `任务ID: ${currentRecord.id}` : '正在初始化...'}
                </p>
              </div>
            </div>

            <div className="flex items-center gap-2">
              <button
                onClick={() => setShowDetails(!showDetails)}
                className="p-2 text-gray-400 hover:text-white hover:bg-gray-800 rounded-lg transition-colors"
              >
                <Eye className="w-4 h-4" />
              </button>
              <button
                onClick={onClose}
                className="p-2 text-gray-400 hover:text-white hover:bg-gray-800 rounded-lg transition-colors"
              >
                <X className="w-4 h-4" />
              </button>
            </div>
          </div>

          {/* Content */}
          <div className="p-6 space-y-6 max-h-[60vh] overflow-y-auto">
            {/* 总体进度 */}
            <div className="bg-gray-800/50 rounded-xl p-6 border border-gray-700">
              <div className="flex items-center justify-between mb-4">
                <div className="flex items-center gap-3">
                  <div className="text-sm font-medium text-gray-300">总体进度</div>
                  {currentRecord && (
                    <div className={`
                      px-2 py-1 rounded-full text-xs font-medium
                      ${currentRecord.status === 'completed' ? 'bg-green-600/20 text-green-400' :
                        currentRecord.status === 'running' ? 'bg-blue-600/20 text-blue-400' :
                        currentRecord.status === 'failed' ? 'bg-red-600/20 text-red-400' :
                        'bg-gray-600/20 text-gray-400'}
                    `}>
                      {currentRecord.status === 'completed' ? '已完成' :
                       currentRecord.status === 'running' ? '进行中' :
                       currentRecord.status === 'failed' ? '失败' : '等待中'}
                    </div>
                  )}
                </div>
                <div className="text-2xl font-bold text-white">
                  {Math.round(overallProgress)}%
                </div>
              </div>

              <div className="w-full bg-gray-700 rounded-full h-2 mb-4">
                <motion.div
                  initial={{ width: 0 }}
                  animate={{ width: `${overallProgress}%` }}
                  transition={{ duration: 0.3 }}
                  className="bg-gradient-to-r from-blue-500 to-blue-600 h-2 rounded-full"
                />
              </div>

              {showDetails && currentRecord && (
                <div className="space-y-2 text-sm">
                  <div className="flex justify-between text-gray-400">
                    <span>开始时间:</span>
                    <span className="text-white">
                      {new Date(currentRecord.start_time).toLocaleString('zh-CN')}
                    </span>
                  </div>
                  <div className="flex justify-between text-gray-400">
                    <span>已用时间:</span>
                    <span className="text-white">
                      {currentRecord.duration ? `${Math.round(currentRecord.duration / 1000)}秒` : '计算中...'}
                    </span>
                  </div>
                  {currentRecord.total_messages > 0 && (
                    <div className="flex justify-between text-gray-400">
                      <span>消息数量:</span>
                      <span className="text-white">{currentRecord.total_messages.toLocaleString()} 条</span>
                    </div>
                  )}
                  {currentRecord.total_contacts > 0 && (
                    <div className="flex justify-between text-gray-400">
                      <span>联系人:</span>
                      <span className="text-white">{currentRecord.total_contacts.toLocaleString()} 人</span>
                    </div>
                  )}
                </div>
              )}
            </div>

            {/* 步骤详情 */}
            <div className="space-y-4">
              {steps.map((step, index) => (
                <motion.div
                  key={step.id}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: index * 0.1 }}
                  className={`
                    bg-gray-800/30 rounded-xl p-4 border transition-all
                    ${step.status === 'running' ? 'border-blue-500/50 bg-blue-600/5' :
                      step.status === 'completed' ? 'border-green-500/50 bg-green-600/5' :
                      step.status === 'error' ? 'border-red-500/50 bg-red-600/5' :
                      'border-gray-700'}
                  `}
                >
                  <div className="flex items-start gap-4">
                    <div className={`
                      w-10 h-10 rounded-lg flex items-center justify-center flex-shrink-0
                      ${step.status === 'running' ? 'bg-blue-600 text-white animate-pulse' :
                        step.status === 'completed' ? 'bg-green-600 text-white' :
                        step.status === 'error' ? 'bg-red-600 text-white' :
                        'bg-gray-700 text-gray-400'}
                    `}>
                      {step.status === 'running' ? (
                        <Loader2 className="w-5 h-5 animate-spin" />
                      ) : step.status === 'error' ? (
                        <AlertCircle className="w-5 h-5" />
                      ) : (
                        step.icon
                      )}
                    </div>

                    <div className="flex-1 min-w-0">
                      <div className="flex items-center justify-between mb-2">
                        <div>
                          <h3 className="font-semibold text-white">{step.title}</h3>
                          <p className="text-sm text-gray-400">{step.description}</p>
                        </div>

                        <div className="text-sm font-medium text-white ml-4">
                          {step.status === 'completed' ? '100%' :
                           step.status === 'running' ? `${Math.round(step.progress)}%` :
                           step.status === 'error' ? '失败' : '等待'}
                        </div>
                      </div>

                      {step.status === 'running' && (
                        <div className="w-full bg-gray-700 rounded-full h-1.5">
                          <motion.div
                            initial={{ width: 0 }}
                            animate={{ width: `${step.progress}%` }}
                            transition={{ duration: 0.3 }}
                            className="bg-blue-500 h-1.5 rounded-full"
                          />
                        </div>
                      )}

                      {step.status === 'error' && step.details && (
                        <div className="mt-2 text-xs text-red-400 bg-red-900/20 rounded p-2">
                          {step.details}
                        </div>
                      )}
                    </div>
                  </div>
                </motion.div>
              ))}
            </div>

            {/* 错误信息 */}
            {error && (
              <motion.div
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                className="bg-red-900/20 border border-red-700 rounded-xl p-4"
              >
                <div className="flex items-start gap-3">
                  <AlertCircle className="w-5 h-5 text-red-400 flex-shrink-0 mt-0.5" />
                  <div className="flex-1">
                    <h4 className="font-semibold text-red-400 mb-1">导入失败</h4>
                    <p className="text-red-300 text-sm">{error}</p>
                  </div>
                </div>

                <div className="mt-4 flex gap-3">
                  <button
                    onClick={handleRetry}
                    className="flex items-center gap-2 px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors text-sm"
                  >
                    <RefreshCw className="w-4 h-4" />
                    重试
                  </button>
                  <button
                    onClick={onClose}
                    className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors text-sm"
                  >
                    关闭
                  </button>
                </div>
              </motion.div>
            )}
          </div>
        </motion.div>
      </motion.div>
    </AnimatePresence>
  );
};

export default WeChatImportProgress;