import React from 'react';
import { motion } from 'framer-motion';
import { Loader2, Brain } from 'lucide-react';

interface LoadingSpinnerProps {
  size?: 'sm' | 'md' | 'lg';
  color?: 'primary' | 'secondary' | 'white';
  text?: string;
  overlay?: boolean;
  className?: string;
}

export const LoadingSpinner: React.FC<LoadingSpinnerProps> = ({
  size = 'md',
  color = 'primary',
  text,
  overlay = false,
  className = '',
}) => {
  const sizeClasses = {
    sm: 'w-4 h-4',
    md: 'w-8 h-8',
    lg: 'w-12 h-12'
  };

  const colorClasses = {
    primary: 'text-blue-500',
    secondary: 'text-gray-400',
    white: 'text-white'
  };

  const textSizeClasses = {
    sm: 'text-sm',
    md: 'text-base',
    lg: 'text-lg'
  };

  const content = (
    <div className={`flex flex-col items-center justify-center space-y-3 ${className}`}>
      <motion.div
        animate={{ rotate: 360 }}
        transition={{
          duration: 1.5,
          repeat: Infinity,
          ease: "linear"
        }}
        className={`${sizeClasses[size]} ${colorClasses[color]}`}
      >
        <Loader2 className="w-full h-full" />
      </motion.div>

      {text && (
        <motion.p
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.2 }}
          className={`${textSizeClasses[size]} ${colorClasses[color]} text-center font-medium`}
        >
          {text}
        </motion.p>
      )}
    </div>
  );

  if (overlay) {
    return (
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        exit={{ opacity: 0 }}
        className="absolute inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50"
      >
        {content}
      </motion.div>
    );
  }

  return content;
};

// 加载状态卡片组件
interface LoadingCardProps {
  title: string;
  description?: string;
  lines?: number;
  className?: string;
}

export const LoadingCard: React.FC<LoadingCardProps> = ({
  title,
  description,
  lines = 3,
  className = '',
}) => {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className={`glass-card rounded-xl p-6 ${className}`}
    >
      {/* 标题骨架 */}
      <div className="flex items-center space-x-3 mb-4">
        <div className="w-8 h-8 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center">
          <motion.div
            animate={{ rotate: 360 }}
            transition={{ duration: 2, repeat: Infinity, ease: "linear" }}
          >
            <Brain className="w-4 h-4 text-white" />
          </motion.div>
        </div>
        <div className="flex-1">
          <div className="h-6 bg-gray-700 rounded-lg w-32 animate-pulse" />
          {description && (
            <div className="h-4 bg-gray-700 rounded-lg w-48 mt-2 animate-pulse opacity-70" />
          )}
        </div>
      </div>

      {/* 内容骨架 */}
      <div className="space-y-3">
        {Array.from({ length: lines }).map((_, index) => (
          <div key={index} className="flex items-start space-x-3">
            <div className="w-2 h-2 bg-blue-500 rounded-full mt-2 flex-shrink-0 animate-pulse" />
            <div className="flex-1 space-y-2">
              <div
                className="h-4 bg-gray-700 rounded-lg animate-pulse"
                style={{
                  width: `${Math.random() * 30 + 70}%`,
                  animationDelay: `${index * 0.1}s`
                }}
              />
              {index % 2 === 0 && (
                <div
                  className="h-3 bg-gray-700 rounded-lg animate-pulse opacity-60"
                  style={{
                    width: `${Math.random() * 20 + 60}%`,
                    animationDelay: `${index * 0.1 + 0.1}s`
                  }}
                />
              )}
            </div>
          </div>
        ))}
      </div>
    </motion.div>
  );
};

// 全屏加载组件
interface FullScreenLoaderProps {
  text?: string;
  subtext?: string;
}

export const FullScreenLoader: React.FC<FullScreenLoaderProps> = ({
  text = '正在加载...',
  subtext
}) => {
  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className="fixed inset-0 bg-black text-white flex items-center justify-center z-50"
    >
      <div className="text-center space-y-6">
        {/* Logo 动画 */}
        <motion.div
          animate={{
            rotate: [0, 10, -10, 0],
            scale: [1, 1.1, 1]
          }}
          transition={{
            rotate: { repeat: Infinity, duration: 4 },
            scale: { repeat: Infinity, duration: 2 }
          }}
          className="w-20 h-20 bg-gradient-to-br from-blue-500 to-purple-600 rounded-2xl flex items-center justify-center mx-auto glow"
        >
          <Brain className="w-10 h-10 text-white" />
        </motion.div>

        {/* 加载文本 */}
        <div className="space-y-2">
          <motion.h1
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 }}
            className="text-2xl font-bold"
          >
            {text}
          </motion.h1>

          {subtext && (
            <motion.p
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.3 }}
              className="text-gray-400"
            >
              {subtext}
            </motion.p>
          )}
        </div>

        {/* 进度指示器 */}
        <motion.div
          initial={{ width: 0 }}
          animate={{ width: '200px' }}
          transition={{ delay: 0.5, duration: 1.5 }}
          className="h-1 bg-gradient-to-r from-blue-500 to-purple-600 rounded-full mx-auto overflow-hidden"
        >
          <motion.div
            animate={{ x: ['-100%', '100%'] }}
            transition={{ repeat: Infinity, duration: 1 }}
            className="h-full bg-white/30"
          />
        </motion.div>
      </div>
    </motion.div>
  );
};

// 进度条组件
interface ProgressBarProps {
  progress: number; // 0-100
  text?: string;
  showPercentage?: boolean;
  color?: 'blue' | 'green' | 'purple';
  className?: string;
}

export const ProgressBar: React.FC<ProgressBarProps> = ({
  progress,
  text,
  showPercentage = true,
  color = 'blue',
  className = '',
}) => {
  const colorClasses = {
    blue: 'bg-blue-500',
    green: 'bg-green-500',
    purple: 'bg-purple-500'
  };

  return (
    <div className={`space-y-2 ${className}`}>
      {(text || showPercentage) && (
        <div className="flex items-center justify-between text-sm text-gray-400">
          <span>{text}</span>
          {showPercentage && <span>{Math.round(progress)}%</span>}
        </div>
      )}
      <div className="w-full h-2 bg-gray-700 rounded-full overflow-hidden">
        <motion.div
          initial={{ width: 0 }}
          animate={{ width: `${progress}%` }}
          transition={{ duration: 0.3, ease: 'easeOut' }}
          className={`h-full ${colorClasses[color]} relative overflow-hidden`}
        >
          <motion.div
            animate={{ x: ['-100%', '100%'] }}
            transition={{ repeat: Infinity, duration: 1 }}
            className="absolute inset-0 bg-white/20"
          />
        </motion.div>
      </div>
    </div>
  );
};

export default LoadingSpinner;