import React from 'react';
import { motion } from 'framer-motion';
import {
  Brain,
  RefreshCw,
  Clock,
  FileText,
  Calendar,
  TrendingUp,
} from 'lucide-react';
import { Briefing } from '../types';
import { MarkdownRenderer } from './MarkdownRenderer';
import { LoadingSpinner } from './LoadingSpinner';
import clsx from 'clsx';

interface BriefingCardProps {
  briefing: Briefing | null;
  isLoading: boolean;
  error?: string | null;
  lastRefreshTime?: string | null;
  onRefresh: () => void;
  className?: string;
}

export const BriefingCard: React.FC<BriefingCardProps> = ({
  briefing,
  isLoading,
  error,
  lastRefreshTime,
  onRefresh,
  className = '',
}) => {
  if (error) {
    return (
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className={clsx('card border-l-4 border-danger-500', className)}
      >
        <div className="flex items-center space-x-3 text-danger-400 mb-3">
          <Brain className="w-5 h-5" />
          <span className="font-medium">简报加载失败</span>
        </div>
        <p className="text-gray-400 text-sm mb-4">{error}</p>
        <button
          onClick={onRefresh}
          className="btn btn-primary btn-sm"
        >
          重试
        </button>
      </motion.div>
    );
  }

  if (isLoading) {
    return (
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className={clsx('card', className)}
      >
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center space-x-3">
            <div className="w-10 h-10 bg-gradient-to-br from-primary-500 to-accent-500 rounded-lg flex items-center justify-center glow">
              <Brain className="w-5 h-5 text-white" />
            </div>
            <div>
              <h3 className="text-lg font-semibold text-white">今日情报简报</h3>
              <p className="text-xs text-gray-400">AI 正在分析数据...</p>
            </div>
          </div>
          <LoadingSpinner size="sm" />
        </div>

        <div className="space-y-3">
          {Array.from({ length: 5 }).map((_, index) => (
            <div
              key={index}
              className="h-4 bg-gray-800 rounded-lg animate-pulse"
              style={{
                width: `${70 + Math.random() * 30}%`,
                animationDelay: `${index * 0.1}s`,
              }}
            />
          ))}
        </div>
      </motion.div>
    );
  }

  if (!briefing) {
    return (
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className={clsx('card text-center py-12', className)}
      >
        <motion.div
          animate={{ y: [0, -10, 0] }}
          transition={{ duration: 3, repeat: Infinity }}
          className="w-16 h-16 bg-gray-800 rounded-2xl flex items-center justify-center mx-auto mb-6"
        >
          <Brain className="w-8 h-8 text-gray-600" />
        </motion.div>
        <h3 className="text-xl font-semibold text-gray-400 mb-2">
          暂无情报简报
        </h3>
        <p className="text-gray-500 mb-6 max-w-md mx-auto">
          请先同步微信数据，然后生成情报简报
        </p>
        <button
          onClick={onRefresh}
          className="btn btn-primary"
        >
          生成情报简报
        </button>
      </motion.div>
    );
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className={clsx('card', className)}
    >
      {/* 卡片头部 */}
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center space-x-3">
          <motion.div
            animate={{ rotate: [0, 5, -5, 0] }}
            transition={{ duration: 4, repeat: Infinity }}
            className="w-10 h-10 bg-gradient-to-br from-primary-500 to-accent-500 rounded-lg flex items-center justify-center glow"
          >
            <Brain className="w-5 h-5 text-white" />
          </motion.div>
          <div>
            <h3 className="text-lg font-semibold text-white">今日情报简报</h3>
            <div className="flex items-center space-x-2 text-xs text-gray-400">
              <FileText className="w-3 h-3" />
              <span>{briefing.message_count || 0} 条消息</span>
              {lastRefreshTime && (
                <>
                  <Clock className="w-3 h-3" />
                  <span>{lastRefreshTime}</span>
                </>
              )}
            </div>
          </div>
        </div>

        <button
          onClick={onRefresh}
          className="btn btn-secondary btn-sm flex items-center space-x-2"
        >
          <RefreshCw className="w-4 h-4" />
          <span>刷新</span>
        </button>
      </div>

      {/* 简报内容 */}
      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        transition={{ delay: 0.1, duration: 0.3 }}
        className="mb-6"
      >
        <MarkdownRenderer
          content={briefing.content}
          className="text-gray-300 leading-relaxed bg-gray-900/30 rounded-lg p-6 border border-gray-700/50"
        />
      </motion.div>

      {/* 统计信息 */}
      <motion.div
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.2, duration: 0.3 }}
        className="grid grid-cols-2 md:grid-cols-4 gap-4"
      >
        <div className="text-center p-3 bg-gray-800/50 rounded-lg">
          <div className="text-2xl font-bold text-primary-400">
            {briefing.message_count || 0}
          </div>
          <div className="text-xs text-gray-400">处理消息</div>
        </div>

        <div className="text-center p-3 bg-gray-800/50 rounded-lg">
          <div className="text-2xl font-bold text-success-400">
            {(() => {
              const lines = briefing.content.split('\n').filter(line => line.trim());
              return lines.length;
            })()}
          </div>
          <div className="text-xs text-gray-400">内容行数</div>
        </div>

        <div className="text-center p-3 bg-gray-800/50 rounded-lg">
          <div className="text-2xl font-bold text-warning-400">
            {(() => {
              const tasks = (briefing.content.match(/- \[ \]/g) || []).length;
              return tasks;
            })()}
          </div>
          <div className="text-xs text-gray-400">待办事项</div>
        </div>

        <div className="text-center p-3 bg-gray-800/50 rounded-lg">
          <div className="flex items-center justify-center space-x-1">
            <Calendar className="w-5 h-5 text-accent-400" />
            <span className="text-2xl font-bold text-accent-400">
              {new Date(briefing.generated_at).toLocaleDateString('zh-CN', {
                month: '2-digit',
                day: '2-digit',
              })}
            </span>
          </div>
          <div className="text-xs text-gray-400">生成日期</div>
        </div>
      </motion.div>

      {/* 卡片底部 */}
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.3, duration: 0.3 }}
        className="flex items-center justify-between pt-4 border-t border-gray-700 text-sm"
      >
        <span className="text-gray-400">由 TwinOS AI 引擎生成</span>
        <div className="flex items-center space-x-2">
          <motion.div
            animate={{ opacity: [1, 0.5, 1] }}
            transition={{ duration: 2, repeat: Infinity }}
            className="w-2 h-2 bg-success-500 rounded-full"
          />
          <span className="text-gray-400">实时分析</span>
        </div>
      </motion.div>
    </motion.div>
  );
};

// 简化版简报组件，用于侧边栏等小空间
export const BriefingMini: React.FC<{
  briefing: Briefing | null;
  isLoading: boolean;
  onClick: () => void;
}> = ({ briefing, isLoading, onClick }) => {
  return (
    <motion.div
      whileHover={{ scale: 1.02 }}
      whileTap={{ scale: 0.98 }}
      onClick={onClick}
      className="p-4 bg-gray-800/50 rounded-lg border border-gray-700 cursor-pointer hover:border-primary-500/50 transition-all"
    >
      <div className="flex items-center justify-between mb-2">
        <div className="flex items-center space-x-2">
          <Brain className="w-4 h-4 text-primary-400" />
          <span className="text-sm font-medium text-white">今日情报</span>
        </div>
        {isLoading && <LoadingSpinner size="sm" />}
      </div>

      {briefing ? (
        <div className="space-y-2">
          <div className="text-xs text-gray-400">
            {briefing.message_count} 条消息已分析
          </div>
          <div className="flex items-center justify-between">
            <span className="text-xs text-gray-500">
              {new Date(briefing.generated_at).toLocaleTimeString('zh-CN', {
                hour: '2-digit',
                minute: '2-digit',
              })}
            </span>
            <TrendingUp className="w-3 h-3 text-success-400" />
          </div>
        </div>
      ) : (
        <div className="text-xs text-gray-500">
          点击生成今日情报简报
        </div>
      )}
    </motion.div>
  );
};

export default BriefingCard;