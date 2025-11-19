import React, { useState } from 'react';
import {
  Activity,
  Brain,
  CheckSquare,
  Users,
  RefreshCw,
  Settings,
  TrendingUp,
  AlertCircle,
  Wifi,
  WifiOff,
  Loader2
} from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import EnhancedMarkdownRenderer from '../components/EnhancedMarkdownRenderer';
import { useAppData } from '../hooks/useAppData';
import clsx from 'clsx';

interface StatCardProps {
  title: string;
  value: string | number;
  change?: {
    value: number;
    type: 'increase' | 'decrease' | 'neutral';
  };
  icon: React.ReactNode;
  loading?: boolean;
}

const StatCard: React.FC<StatCardProps> = ({ title, value, change, icon, loading }) => (
  <motion.div
    initial={{ opacity: 0, y: 20 }}
    animate={{ opacity: 1, y: 0 }}
    transition={{ duration: 0.3 }}
    className="bg-gray-800/50 backdrop-blur-sm rounded-xl p-6 border border-gray-700"
  >
    <div className="flex items-center justify-between mb-4">
      <div className="text-gray-400">{title}</div>
      {loading ? (
        <Loader2 className="w-5 h-5 text-gray-400 animate-spin" />
      ) : (
        <div className="text-blue-400">{icon}</div>
      )}
    </div>
    <div className="text-2xl font-bold text-white mb-2">{value}</div>
    {change && (
      <div className={clsx(
        'flex items-center text-sm',
        change.type === 'increase' && 'text-green-400',
        change.type === 'decrease' && 'text-red-400',
        change.type === 'neutral' && 'text-gray-400'
      )}>
        <TrendingUp className="w-4 h-4 mr-1" />
        {change.value > 0 ? '+' : ''}{change.value}%
      </div>
    )}
  </motion.div>
);

const NavigationItem: React.FC<{
  icon: React.ReactNode;
  label: string;
  active: boolean;
  count?: number;
  onClick: () => void;
}> = ({ icon, label, active, count, onClick }) => (
  <motion.button
    whileHover={{ scale: 1.02 }}
    whileTap={{ scale: 0.98 }}
    onClick={onClick}
    className={clsx(
      'w-full flex items-center justify-between px-4 py-3 rounded-xl transition-all duration-200 text-left',
      active
        ? 'bg-gradient-to-r from-blue-600 to-blue-500 text-white shadow-lg shadow-blue-600/25 border border-blue-500/30'
        : 'text-gray-300 hover:bg-gray-700/50 hover:text-white border border-transparent'
    )}
  >
    <div className="flex items-center gap-3 flex-1">
      <div className="w-5 h-5 flex items-center justify-center flex-shrink-0">
        {icon}
      </div>
      <span className="font-medium text-sm truncate">{label}</span>
    </div>
    {count !== undefined && count > 0 && (
      <div className="bg-red-500/20 text-red-400 border border-red-500/30 text-xs px-2 py-1 rounded-full font-medium min-w-[20px] text-center flex-shrink-0">
        {count}
      </div>
    )}
  </motion.button>
);

const EnhancedDashboard: React.FC = () => {
  const [activeView, setActiveView] = useState<'dashboard' | 'briefing' | 'todos' | 'connections' | 'settings'>('dashboard');
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);

  const {
    systemStatus,
    briefing,
    todos,
    connections,
    loading,
    generateBriefing,
    extractTodos,
    analyzeConnections,
    refreshStatus
  } = useAppData();

  const pendingTodosCount = todos.filter(todo => todo.status === 'pending').length;
  const strongConnectionsCount = connections.filter(conn => conn.strength_score > 0.7).length;

  const navigationItems = [
    {
      id: 'dashboard',
      label: '仪表盘',
      icon: <Activity />,
    },
    {
      id: 'briefing',
      label: '今日情报',
      icon: <Brain />,
      count: 1,
    },
    {
      id: 'todos',
      label: '决策待办',
      icon: <CheckSquare />,
      count: pendingTodosCount,
    },
    {
      id: 'connections',
      label: '人脉雷达',
      icon: <Users />,
      count: strongConnectionsCount,
    },
    {
      id: 'settings',
      label: '系统设置',
      icon: <Settings />,
    },
  ];

  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1,
      },
    },
  };

  const itemVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: { opacity: 1, y: 0 },
  };

  return (
    <div className="min-h-screen bg-black text-white flex">
      {/* 侧边栏 */}
      <motion.aside
        initial={false}
        animate={{
          width: sidebarCollapsed ? 80 : 280,
        }}
        transition={{ duration: 0.3 }}
        className="bg-gray-900/50 backdrop-blur-sm border-r border-gray-800 flex flex-col"
      >
        {/* Logo */}
        <div className="p-6 border-b border-gray-800">
          <div className={clsx('flex items-center gap-3', sidebarCollapsed && 'justify-center')}>
            <div className="w-10 h-10 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center font-bold text-lg">
              T
            </div>
            {!sidebarCollapsed && (
              <div>
                <h1 className="text-xl font-bold">TwinOS</h1>
                <p className="text-sm text-gray-400">智能决策副驾系统</p>
              </div>
            )}
          </div>
        </div>

        {/* 导航 */}
        <nav className="flex-1 px-3 py-4 overflow-y-auto">
          <div className="mb-6">
            <div className="space-y-1">
              {navigationItems.map((item) => (
                <NavigationItem
                  key={item.id}
                  icon={item.icon}
                  label={item.label}
                  active={activeView === item.id}
                  count={item.count}
                  onClick={() => setActiveView(item.id as any)}
                />
              ))}
            </div>
          </div>
        </nav>

        {/* 系统状态 */}
        <div className="px-3 py-4 border-t border-gray-700/50">
          <div className="mb-4">
            <h2 className="text-xs font-semibold text-gray-500 uppercase tracking-wider px-4 mb-3">
              系统状态
            </h2>
            <div className="bg-gray-800/50 rounded-xl p-3 border border-gray-700/50 space-y-3">
              {/* 服务器连接状态 */}
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-300 font-medium">服务器连接</span>
                <div className="flex items-center gap-2">
                  {systemStatus.backendStatus === 'connected' ? (
                    <>
                      <Wifi className="w-4 h-4 text-green-400" />
                      <span className="text-xs text-green-400">正常</span>
                    </>
                  ) : systemStatus.backendStatus === 'loading' ? (
                    <>
                      <Loader2 className="w-4 h-4 text-yellow-400 animate-spin" />
                      <span className="text-xs text-yellow-400">连接中</span>
                    </>
                  ) : (
                    <>
                      <WifiOff className="w-4 h-4 text-red-400" />
                      <span className="text-xs text-red-400">断开</span>
                    </>
                  )}
                </div>
              </div>

              {/* 微信数据库连接状态 */}
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-300 font-medium">微信数据库</span>
                <div className="flex items-center gap-2">
                  {systemStatus.isOnline && systemStatus.backendStatus === 'connected' ? (
                    <>
                      <Activity className="w-4 h-4 text-green-400" />
                      <span className="text-xs text-green-400">已同步</span>
                    </>
                  ) : systemStatus.backendStatus === 'loading' ? (
                    <>
                      <Loader2 className="w-4 h-4 text-yellow-400 animate-spin" />
                      <span className="text-xs text-yellow-400">检测中</span>
                    </>
                  ) : (
                    <>
                      <AlertCircle className="w-4 h-4 text-red-400" />
                      <span className="text-xs text-red-400">未连接</span>
                    </>
                  )}
                </div>
              </div>

              {/* 数据统计 */}
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-300 font-medium">消息数量</span>
                <span className="text-xs text-gray-400">
                  {systemStatus.dataCount.messages.toLocaleString()} 条
                </span>
              </div>

              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-300 font-medium">联系人数量</span>
                <span className="text-xs text-gray-400">
                  {systemStatus.dataCount.contacts.toLocaleString()} 人
                </span>
              </div>

              <div className="text-xs text-gray-500 pt-2 border-t border-gray-700/50">
                最后更新: {new Date(systemStatus.lastUpdate).toLocaleTimeString('zh-CN')}
              </div>

              <button
                onClick={refreshStatus}
                className="w-full text-xs bg-gray-700/50 hover:bg-gray-600/50 text-gray-300 px-3 py-2 rounded-lg transition-colors border border-gray-600/30"
              >
                刷新状态
              </button>
            </div>
          </div>
        </div>
      </motion.aside>

      {/* 主内容区 */}
      <main className="flex-1 overflow-hidden">
        <div className="h-full overflow-y-auto p-6">
          <AnimatePresence mode="wait">
            <motion.div
              key={activeView}
              variants={containerVariants}
              initial="hidden"
              animate="visible"
              exit="hidden"
            >
              {activeView === 'dashboard' && (
                <motion.div variants={itemVariants} className="space-y-6">
                  <h1 className="text-3xl font-bold text-white mb-2">仪表盘概览</h1>

                  {/* 统计卡片 */}
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                    <StatCard
                      title="总消息数"
                      value={systemStatus.dataCount.messages}
                      change={{ value: 12.5, type: 'increase' }}
                      icon={<Activity />}
                    />
                    <StatCard
                      title="联系人"
                      value={systemStatus.dataCount.contacts}
                      change={{ value: 8.3, type: 'increase' }}
                      icon={<Users />}
                    />
                    <StatCard
                      title="今日简报"
                      value="1"
                      icon={<Brain />}
                    />
                    <StatCard
                      title="待办事项"
                      value={pendingTodosCount}
                      icon={<CheckSquare />}
                    />
                  </div>

                  {/* 快速操作 */}
                  <div className="bg-gray-800/50 backdrop-blur-sm rounded-xl p-6 border border-gray-700">
                    <h2 className="text-xl font-semibold text-white mb-4">快速操作</h2>
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                      <button
                        onClick={generateBriefing}
                        disabled={loading.briefing}
                        className={clsx(
                          'flex items-center justify-center gap-2 px-4 py-3 rounded-lg transition-colors',
                          loading.briefing
                            ? 'bg-gray-700 text-gray-400'
                            : 'bg-blue-600 hover:bg-blue-700 text-white'
                        )}
                      >
                        <Brain className="w-4 h-4" />
                        {loading.briefing ? '生成中...' : '生成简报'}
                      </button>
                      <button
                        onClick={extractTodos}
                        disabled={loading.todos}
                        className={clsx(
                          'flex items-center justify-center gap-2 px-4 py-3 rounded-lg transition-colors',
                          loading.todos
                            ? 'bg-gray-700 text-gray-400'
                            : 'bg-green-600 hover:bg-green-700 text-white'
                        )}
                      >
                        <CheckSquare className="w-4 h-4" />
                        {loading.todos ? '提取中...' : '提取待办'}
                      </button>
                      <button
                        onClick={analyzeConnections}
                        disabled={loading.connections}
                        className={clsx(
                          'flex items-center justify-center gap-2 px-4 py-3 rounded-lg transition-colors',
                          loading.connections
                            ? 'bg-gray-700 text-gray-400'
                            : 'bg-purple-600 hover:bg-purple-700 text-white'
                        )}
                      >
                        <Users className="w-4 h-4" />
                        {loading.connections ? '分析中...' : '分析人脉'}
                      </button>
                    </div>
                  </div>
                </motion.div>
              )}

              {activeView === 'briefing' && (
                <motion.div variants={itemVariants} className="space-y-6">
                  <div className="flex items-center justify-between">
                    <h1 className="text-3xl font-bold text-white">今日情报</h1>
                    <button
                      onClick={generateBriefing}
                      disabled={loading.briefing}
                      className="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors disabled:bg-gray-700 disabled:text-gray-400"
                    >
                      <RefreshCw className={clsx('w-4 h-4', loading.briefing && 'animate-spin')} />
                      {loading.briefing ? '生成中...' : '刷新简报'}
                    </button>
                  </div>
                  <div className="bg-gray-800/50 backdrop-blur-sm rounded-xl p-6 border border-gray-700">
                    <EnhancedMarkdownRenderer
                      content={briefing}
                      enableCopy={true}
                      maxHeight="70vh"
                    />
                  </div>
                </motion.div>
              )}

              {activeView === 'todos' && (
                <motion.div variants={itemVariants} className="space-y-6">
                  <div className="flex items-center justify-between">
                    <h1 className="text-3xl font-bold text-white">决策待办</h1>
                    <button
                      onClick={extractTodos}
                      disabled={loading.todos}
                      className="flex items-center gap-2 px-4 py-2 bg-green-600 hover:bg-green-700 text-white rounded-lg transition-colors disabled:bg-gray-700 disabled:text-gray-400"
                    >
                      <RefreshCw className={clsx('w-4 h-4', loading.todos && 'animate-spin')} />
                      {loading.todos ? '刷新中...' : '刷新待办'}
                    </button>
                  </div>

                  <div className="space-y-4">
                    {todos.map((todo, index) => (
                      <motion.div
                        key={todo.id}
                        initial={{ opacity: 0, x: -20 }}
                        animate={{ opacity: 1, x: 0 }}
                        transition={{ delay: index * 0.1 }}
                        className={clsx(
                          'bg-gray-800/50 backdrop-blur-sm rounded-xl p-6 border border-gray-700',
                          todo.status === 'completed' && 'opacity-60'
                        )}
                      >
                        <div className="flex items-start justify-between">
                          <div className="flex-1">
                            <div className="flex items-center gap-3 mb-2">
                              <div
                                className={clsx(
                                  'w-2 h-2 rounded-full',
                                  todo.priority === 'high' && 'bg-red-500',
                                  todo.priority === 'medium' && 'bg-yellow-500',
                                  todo.priority === 'low' && 'bg-green-500'
                                )}
                              />
                              <h3 className={clsx(
                                'font-semibold',
                                todo.status === 'completed' && 'line-through text-gray-500'
                              )}>
                                {todo.content}
                              </h3>
                            </div>
                            <div className="flex items-center gap-4 text-sm text-gray-400 mb-2">
                              <span>负责人: {todo.assignee}</span>
                              <span>截止: {todo.due_date}</span>
                            </div>
                            {todo.tags.length > 0 && (
                              <div className="flex gap-2 flex-wrap">
                                {todo.tags.map((tag) => (
                                  <span
                                    key={tag}
                                    className="px-2 py-1 bg-blue-600/20 text-blue-400 text-xs rounded"
                                  >
                                    {tag}
                                  </span>
                                ))}
                              </div>
                            )}
                          </div>
                          <div className="ml-4">
                            <span
                              className={clsx(
                                'px-3 py-1 rounded-lg text-sm font-medium',
                                todo.status === 'completed'
                                  ? 'bg-green-600/20 text-green-400'
                                  : todo.status === 'pending'
                                  ? 'bg-yellow-600/20 text-yellow-400'
                                  : 'bg-gray-600/20 text-gray-400'
                              )}
                            >
                              {todo.status === 'completed' ? '已完成' : todo.status === 'pending' ? '待处理' : '已取消'}
                            </span>
                          </div>
                        </div>
                      </motion.div>
                    ))}
                  </div>
                </motion.div>
              )}

              {activeView === 'connections' && (
                <motion.div variants={itemVariants} className="space-y-6">
                  <div className="flex items-center justify-between">
                    <h1 className="text-3xl font-bold text-white">人脉雷达</h1>
                    <button
                      onClick={analyzeConnections}
                      disabled={loading.connections}
                      className="flex items-center gap-2 px-4 py-2 bg-purple-600 hover:bg-purple-700 text-white rounded-lg transition-colors disabled:bg-gray-700 disabled:text-gray-400"
                    >
                      <RefreshCw className={clsx('w-4 h-4', loading.connections && 'animate-spin')} />
                      {loading.connections ? '分析中...' : '刷新分析'}
                    </button>
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {connections.map((connection, index) => (
                      <motion.div
                        key={connection.id}
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: index * 0.1 }}
                        className="bg-gray-800/50 backdrop-blur-sm rounded-xl p-6 border border-gray-700"
                      >
                        <div className="flex items-center justify-between mb-4">
                          <h3 className="text-lg font-semibold text-white">
                            {connection.person_name}
                          </h3>
                          <div className="px-2 py-1 bg-blue-600/20 text-blue-400 text-xs rounded">
                            {connection.relationship}
                          </div>
                        </div>
                        <div className="space-y-2 text-sm">
                          <div className="flex justify-between">
                            <span className="text-gray-400">互动次数:</span>
                            <span className="text-white">{connection.interaction_count}</span>
                          </div>
                          <div className="flex justify-between">
                            <span className="text-gray-400">关系强度:</span>
                            <span className="text-white">
                              {(connection.strength_score * 100).toFixed(0)}%
                            </span>
                          </div>
                          <div className="flex justify-between">
                            <span className="text-gray-400">最后互动:</span>
                            <span className="text-white">{connection.last_interaction}</span>
                          </div>
                        </div>
                        <div className="mt-4">
                          <div className="flex flex-wrap gap-2">
                            {connection.topics.map((topic) => (
                              <span
                                key={topic}
                                className="px-2 py-1 bg-purple-600/20 text-purple-400 text-xs rounded"
                              >
                                {topic}
                              </span>
                            ))}
                          </div>
                        </div>
                      </motion.div>
                    ))}
                  </div>
                </motion.div>
              )}

              {activeView === 'settings' && (
                <motion.div variants={itemVariants} className="space-y-6">
                  <h1 className="text-3xl font-bold text-white mb-6">系统设置</h1>

                  <div className="bg-gray-800/50 backdrop-blur-sm rounded-xl p-6 border border-gray-700">
                    <h2 className="text-xl font-semibold text-white mb-4">系统信息</h2>
                    <div className="space-y-4">
                      <div className="flex justify-between">
                        <span className="text-gray-400">应用版本:</span>
                        <span className="text-white">v0.3.0</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-400">后端状态:</span>
                        <span className={clsx(
                          'text-sm font-medium',
                          systemStatus.backendStatus === 'connected' && 'text-green-400',
                          systemStatus.backendStatus === 'disconnected' && 'text-red-400',
                          systemStatus.backendStatus === 'loading' && 'text-yellow-400'
                        )}>
                          {systemStatus.backendStatus === 'connected' && '连接正常'}
                          {systemStatus.backendStatus === 'disconnected' && '连接失败'}
                          {systemStatus.backendStatus === 'loading' && '检查中...'}
                        </span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-400">最后同步:</span>
                        <span className="text-white">
                          {new Date(systemStatus.lastUpdate).toLocaleString('zh-CN')}
                        </span>
                      </div>
                    </div>
                  </div>
                </motion.div>
              )}
            </motion.div>
          </AnimatePresence>
        </div>
      </main>
    </div>
  );
};

export default EnhancedDashboard;