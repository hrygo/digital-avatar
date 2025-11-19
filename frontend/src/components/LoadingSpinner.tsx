import React from 'react';
import { motion } from 'framer-motion';
import { Loader2, Brain, Activity } from 'lucide-react';
import clsx from 'clsx';

interface LoadingSpinnerProps {
  size?: 'sm' | 'md' | 'lg';
  color?: 'primary' | 'secondary' | 'success' | 'warning' | 'danger';
  text?: string;
  className?: string;
}

export const LoadingSpinner: React.FC<LoadingSpinnerProps> = ({
  size = 'md',
  color = 'primary',
  text,
  className = '',
}) => {
  const sizeClasses = {
    sm: 'w-4 h-4',
    md: 'w-6 h-6',
    lg: 'w-8 h-8',
  };

  const colorClasses = {
    primary: 'text-primary-500',
    secondary: 'text-gray-500',
    success: 'text-success-500',
    warning: 'text-warning-500',
    danger: 'text-danger-500',
  };

  return (
    <div className={clsx('flex items-center justify-center', className)}>
      <motion.div
        animate={{ rotate: 360 }}
        transition={{ duration: 1, repeat: Infinity, ease: 'linear' }}
      >
        <Loader2 className={clsx(sizeClasses[size], colorClasses[color])} />
      </motion.div>
      {text && (
        <span className={clsx('ml-2 text-sm text-gray-400')}>
          {text}
        </span>
      )}
    </div>
  );
};

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
    <div className={clsx('card', className)}>
      <div className="flex items-center justify-between mb-4">
        <h3 className="card-title">{title}</h3>
        <LoadingSpinner size="sm" />
      </div>
      {description && (
        <p className="text-gray-400 text-sm mb-4">{description}</p>
      )}
      <div className="space-y-2">
        {Array.from({ length: lines }).map((_, index) => (
          <motion.div
            key={index}
            className="h-4 bg-gray-800 rounded-lg"
            initial={{ width: '60%' }}
            animate={{ width: ['60%', '90%', '60%'] }}
            transition={{
              duration: 1.5,
              repeat: Infinity,
              delay: index * 0.1,
            }}
            style={{
              width: `${60 + Math.random() * 30}%`,
            }}
          />
        ))}
      </div>
    </div>
  );
};

interface FullScreenLoaderProps {
  text?: string;
  subtext?: string;
  icon?: 'brain' | 'activity' | 'loader';
  className?: string;
}

export const FullScreenLoader: React.FC<FullScreenLoaderProps> = ({
  text = '系统初始化中',
  subtext = '请稍候...',
  icon = 'brain',
  className = '',
}) => {
  const iconComponents = {
    brain: Brain,
    activity: Activity,
    loader: Loader2,
  };

  const IconComponent = iconComponents[icon];

  return (
    <div className={clsx(
      'fixed inset-0 bg-dark flex flex-col items-center justify-center z-50',
      className
    )}>
      <motion.div
        initial={{ scale: 0.8, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        transition={{ duration: 0.5 }}
        className="text-center"
      >
        <motion.div
          animate={{
            rotate: [0, 10, -10, 0],
            scale: [1, 1.1, 1],
          }}
          transition={{
            duration: 4,
            repeat: Infinity,
            ease: 'easeInOut',
          }}
          className="w-16 h-16 bg-gradient-to-br from-primary-500 to-accent-500 rounded-2xl flex items-center justify-center mb-8 glow-lg"
        >
          <IconComponent className="w-8 h-8 text-white" />
        </motion.div>

        <motion.h1
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.2, duration: 0.5 }}
          className="text-2xl font-bold text-white mb-2"
        >
          {text}
        </motion.h1>

        <motion.p
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.4, duration: 0.5 }}
          className="text-gray-400 mb-8"
        >
          {subtext}
        </motion.p>

        <motion.div
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.6, duration: 0.5 }}
          className="flex space-x-2"
        >
          {[0, 1, 2].map((index) => (
            <motion.div
              key={index}
              className="w-2 h-2 bg-primary-500 rounded-full"
              animate={{
                scale: [1, 1.5, 1],
                opacity: [0.5, 1, 0.5],
              }}
              transition={{
                duration: 1.5,
                repeat: Infinity,
                delay: index * 0.2,
              }}
            />
          ))}
        </motion.div>
      </motion.div>

      {/* 背景动画 */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <motion.div
          className="absolute top-0 left-0 w-full h-full"
          animate={{
            background: [
              'radial-gradient(circle at 20% 50%, rgba(59, 130, 246, 0.1) 0%, transparent 50%)',
              'radial-gradient(circle at 80% 50%, rgba(147, 51, 234, 0.1) 0%, transparent 50%)',
              'radial-gradient(circle at 50% 20%, rgba(34, 197, 94, 0.1) 0%, transparent 50%)',
              'radial-gradient(circle at 50% 80%, rgba(245, 158, 11, 0.1) 0%, transparent 50%)',
              'radial-gradient(circle at 20% 50%, rgba(59, 130, 246, 0.1) 0%, transparent 50%)',
            ],
          }}
          transition={{ duration: 8, repeat: Infinity }}
        />
      </div>
    </div>
  );
};

interface LoadingProgressProps {
  progress: number;
  text?: string;
  showPercentage?: boolean;
  className?: string;
}

export const LoadingProgress: React.FC<LoadingProgressProps> = ({
  progress,
  text,
  showPercentage = true,
  className = '',
}) => {
  return (
    <div className={clsx('w-full', className)}>
      {text && (
        <div className="flex justify-between items-center mb-2">
          <span className="text-sm text-gray-400">{text}</span>
          {showPercentage && (
            <span className="text-sm text-gray-400">{progress}%</span>
          )}
        </div>
      )}
      <div className="w-full bg-gray-800 rounded-full h-2 overflow-hidden">
        <motion.div
          className="h-full bg-gradient-to-r from-primary-500 to-accent-500 rounded-full"
          initial={{ width: 0 }}
          animate={{ width: `${progress}%` }}
          transition={{ duration: 0.5, ease: 'easeOut' }}
        />
      </div>
    </div>
  );
};

interface SkeletonProps {
  className?: string;
  children?: React.ReactNode;
}

export const Skeleton: React.FC<SkeletonProps> = ({ className = '', children }) => {
  return (
    <div
      className={clsx(
        'animate-pulse bg-gray-800 rounded',
        className
      )}
    >
      {children}
    </div>
  );
};

// 预设的加载骨架组件
export const CardSkeleton: React.FC<{ lines?: number; showAvatar?: boolean }> = ({
  lines = 3,
  showAvatar = false,
}) => {
  return (
    <div className="card">
      <div className="flex items-start space-x-4">
        {showAvatar && (
          <Skeleton className="w-12 h-12 rounded-full flex-shrink-0" />
        )}
        <div className="flex-1 space-y-2">
          <Skeleton className="h-4 w-3/4" />
          {Array.from({ length: lines - 1 }).map((_, index) => (
            <Skeleton
              key={index}
              className={`h-4 w-[${80 + Math.random() * 20}%]`}
            />
          ))}
        </div>
      </div>
    </div>
  );
};

export const TableSkeleton: React.FC<{ rows?: number; columns?: number }> = ({
  rows = 5,
  columns = 4,
}) => {
  return (
    <div className="overflow-x-auto">
      <table className="table">
        <thead>
          <tr>
            {Array.from({ length: columns }).map((_, index) => (
              <th key={index}>
                <Skeleton className="h-4 w-20" />
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {Array.from({ length: rows }).map((_, rowIndex) => (
            <tr key={rowIndex}>
              {Array.from({ length: columns }).map((_, colIndex) => (
                <td key={colIndex}>
                  <Skeleton
                    className={`h-4 w-[${60 + Math.random() * 40}%]`}
                  />
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default LoadingSpinner;