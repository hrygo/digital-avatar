import React, { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import {
  Brain,
  CheckSquare,
  Users,
  Activity,
  Database,
  AlertCircle,
  RefreshCw,
  Settings,
  Menu,
  X,
  Bell,
  Search,
} from 'lucide-react';
import {
  useAppStore,
  useWeChatStatus,
  useBriefing,
  useTodos,
  useConnections,
  useIsLoading,
  useError,
  useActiveView,
  useRefreshStates,
  useLastRefreshTimes,
} from '../store/useAppStore';
import { BriefingCard } from '../components/BriefingCard';
import { LoadingSpinner } from '../components/LoadingSpinner';
import clsx from 'clsx';

// 导航项配置
const navigationItems = [
  {
    id: 'dashboard',
    label: '仪表盘',
    icon: Activity,
  },
  {
    id: 'briefing',
    label: '今日情报',
    icon: Brain,
  },
  {
    id: 'todos',
    label: '决策待办',
    icon: CheckSquare,
  },
  {
    id: 'connections',
    label: '人脉雷达',
    icon: Users,
  },
  {
    id: 'settings',
    label: '系统设置',
    icon: Settings,
  },
];

const Dashboard: React.FC = () => {
  const {
    isInitialized,
    error,
    initializeApp,
    setActiveView,
    toggleSidebar,
    syncWeChat,
    generateBriefing,
    extractTodos,
    analyzeConnections,
    clearError,
  } = useAppStore();

  // 状态选择器
  const weChatStatus = useWeChatStatus();
  const briefing = useBriefing();
  const todos = useTodos();
  const connections = useConnections();
  const isLoading = useIsLoading();
  const activeView = useActiveView();
  const refreshStates = useRefreshStates();
  const lastRefreshTimes = useLastRefreshTimes();

  // 本地状态
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [showNotifications, setShowNotifications] = useState(false);

  // 初始化应用
  useEffect(() => {
    if (!isInitialized) {
      initializeApp();
    }
  }, [isInitialized, initializeApp]);

  // 更新刷新时间的函数
  const updateRefreshTime = (type: keyof typeof lastRefreshTimes) => {
    useAppStore.getState().updateRefreshTime(type);
  };

  // 统一刷新函数
  const handleRefresh = async (type: 'wechat' | 'briefing' | 'todos' | 'connections') => {
    try {
      switch (type) {
        case 'wechat':
          if (weChatStatus?.is_connected) {
            await syncWeChat();
            updateRefreshTime('wechat');
          }
          break;
        case 'briefing':
          await generateBriefing();
          updateRefreshTime('briefing');
          break;
        case 'todos':
          await extractTodos();
          updateRefreshTime('todos');
          break;
        case 'connections':
          await analyzeConnections();
          updateRefreshTime('connections');
          break;
      }
    } catch (error) {
      console.error(`Failed to refresh ${type}:`, error);
    }
  };

  // 全部刷新
  const handleRefreshAll = async () => {
    try {
      await Promise.all([
        generateBriefing(),
        extractTodos(),
        analyzeConnections(),
      ]);
      updateRefreshTime('briefing');
      updateRefreshTime('todos');
      updateRefreshTime('connections');
    } catch (error) {
      console.error('Failed to refresh all:', error);
    }
  };

  // 容器动画变体
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
    hidden: { y: 20, opacity: 0 },
    visible: {
      y: 0,
      opacity: 1,
      transition: {
        type: 'spring',
        stiffness: 100,
      },
    },
  };

  // 初始化加载状态
  if (!isInitialized || isLoading) {
    return (
      <div className="min-h-screen bg-dark flex items-center justify-center">
        <LoadingSpinner size="lg" text="正在初始化 TwinOS..." />
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-dark text-white">
      {/* 顶部导航栏 */}
      <motion.div
        initial={{ y: -50, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        className="fixed top-0 left-0 right-0 z-40 glass-card border-b border-gray-700"
      >
        <div className="px-6 py-4">
          <div className="flex items-center justify-between">
            {/* 左侧：Logo和导航 */}
            <div className="flex items-center space-x-4">
              {/* 菜单按钮 */}
              <button
                onClick={toggleSidebar}
                className="p-2 text-gray-400 hover:text-white hover:bg-gray-800/50 rounded-lg transition-colors"
              >
                {sidebarCollapsed ? <Menu className="w-5 h-5" /> : <X className="w-5 h-5" />}
              </button>

              {/* Logo */}
              <div className="flex items-center space-x-3">
                <motion.div
                  animate={{ rotate: [0, 10, -10, 0] }}
                  transition={{ repeat: Infinity, duration: 4 }}
                  className="w-8 h-8 bg-gradient-to-br from-primary-500 to-accent-500 rounded-lg flex items-center justify-center glow"
                >
                  <Brain className="w-5 h-5 text-white" />
                </motion.div>
                <div>
                  <h1 className="text-xl font-bold text-white">TWINOS</h1>
                  <p className="text-xs text-gray-400">智能决策副驾</p>
                </div>
              </div>

              {/* 主导航 */}
              <nav className="hidden md:flex items-center space-x-1">
                {navigationItems.map((item, index) => (
                  <motion.button
                    key={item.id}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: index * 0.1 }}
                    whileHover={{ scale: 1.05 }}
                    whileTap={{ scale: 0.95 }}
                    onClick={() => setActiveView(item.id as any)}
                    className={clsx(
                      'flex items-center space-x-2 px-3 py-2 rounded-lg text-sm font-medium transition-all duration-200',
                      activeView === item.id
                        ? 'bg-primary-600 text-white shadow-lg'
                        : 'text-gray-400 hover:text-white hover:bg-gray-800/50'
                    )}
                  >
                    <item.icon className="w-4 h-4" />
                    <span>{item.label}</span>
                  </motion.button>
                ))}
              </nav>
            </div>

            {/* 右侧：状态、搜索、通知、操作 */}
            <div className="flex items-center space-x-4">
              {/* 微信状态 */}
              <div className="hidden sm:flex items-center space-x-2 px-3 py-2 rounded-lg border border-gray-700 bg-gray-800/50">
                <div
                  className={clsx(
                    'w-2 h-2 rounded-full',
                    weChatStatus?.is_connected
                      ? 'bg-success-500 animate-pulse'
                      : 'bg-danger-500 animate-pulse'
                  )}
                />
                <span className="text-sm text-gray-300">
                  {weChatStatus?.is_connected ? '微信已连接' : '微信未连接'}
                </span>
              </div>

              {/* 搜索框 */}
              <div className="hidden sm:block">
                <div className="relative">
                  <Search className="w-4 h-4 text-gray-400 absolute left-3 top-1/2 transform -translate-y-1/2" />
                  <input
                    type="text"
                    placeholder="搜索..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="pl-10 pr-4 py-2 bg-gray-800/50 border border-gray-700 rounded-lg text-sm text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent w-48"
                  />
                </div>
              </div>

              {/* 通知按钮 */}
              <button
                onClick={() => setShowNotifications(!showNotifications)}
                className="relative p-2 text-gray-400 hover:text-white hover:bg-gray-800/50 rounded-lg transition-colors"
              >
                <Bell className="w-5 h-5" />
                {/* 通知徽章 */}
                <span className="absolute top-1 right-1 w-2 h-2 bg-danger-500 rounded-full" />
              </button>

              {/* 刷新按钮组 */}
              <div className="flex items-center space-x-2">
                {/* 单独刷新 */}
                <div className="relative group">
                  <button className="btn btn-secondary btn-sm flex items-center space-x-2">
                    <RefreshCw className="w-4 h-4" />
                    <span className="hidden sm:inline">刷新</span>
                  </button>

                  {/* 下拉菜单 */}
                  <div className="absolute right-0 mt-2 w-48 glass-card rounded-lg shadow-xl opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 z-50">
                    <div className="p-2 space-y-1">
                      <button
                        onClick={() => handleRefresh('wechat')}
                        disabled={refreshStates.wechat || !weChatStatus?.is_connected}
                        className={clsx(
                          'w-full text-left px-3 py-2 rounded text-sm flex items-center justify-between',
                          refreshStates.wechat || !weChatStatus?.is_connected
                            ? 'opacity-50 cursor-not-allowed'
                            : 'hover:bg-gray-700/50 transition-colors'
                        )}
                      >
                        <span>同步微信</span>
                        {refreshStates.wechat && <LoadingSpinner size="sm" />}
                      </button>

                      <button
                        onClick={() => handleRefresh('briefing')}
                        disabled={refreshStates.briefing}
                        className={clsx(
                          'w-full text-left px-3 py-2 rounded text-sm flex items-center justify-between',
                          refreshStates.briefing
                            ? 'opacity-50 cursor-not-allowed'
                            : 'hover:bg-gray-700/50 transition-colors'
                        )}
                      >
                        <span>刷新简报</span>
                        {refreshStates.briefing && <LoadingSpinner size="sm" />}
                      </button>

                      <button
                        onClick={() => handleRefresh('todos')}
                        disabled={refreshStates.todos}
                        className={clsx(
                          'w-full text-left px-3 py-2 rounded text-sm flex items-center justify-between',
                          refreshStates.todos
                            ? 'opacity-50 cursor-not-allowed'
                            : 'hover:bg-gray-700/50 transition-colors'
                        )}
                      >
                        <span>刷新待办</span>
                        {refreshStates.todos && <LoadingSpinner size="sm" />}
                      </button>

                      <button
                        onClick={() => handleRefresh('connections')}
                        disabled={refreshStates.connections}
                        className={clsx(
                          'w-full text-left px-3 py-2 rounded text-sm flex items-center justify-between',
                          refreshStates.connections
                            ? 'opacity-50 cursor-not-allowed'
                            : 'hover:bg-gray-700/50 transition-colors'
                        )}
                      >
                        <span>刷新人脉</span>
                        {refreshStates.connections && <LoadingSpinner size="sm" />}
                      </button>
                    </div>
                  </div>
                </div>

                {/* 全部刷新 */}
                <button
                  onClick={handleRefreshAll}
                  disabled={Object.values(refreshStates).some(Boolean)}
                  className={clsx(
                    'btn btn-primary btn-sm flex items-center space-x-2',
                    Object.values(refreshStates).some(Boolean) && 'opacity-50 cursor-not-allowed'
                  )}
                >
                  <Activity className={clsx('w-4 h-4', Object.values(refreshStates).some(Boolean) && 'animate-spin')} />
                  <span className="hidden sm:inline">
                    {Object.values(refreshStates).some(Boolean) ? '刷新中...' : '全部刷新'}
                  </span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </motion.div>

      {/* 错误提示 */}
      {error && (
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          className="fixed top-20 left-4 right-4 z-50 glass-card border-l-4 border-danger-500 p-4"
        >
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-2">
              <AlertCircle className="w-5 h-5 text-danger-400" />
              <span className="text-danger-200">{error}</span>
            </div>
            <button
              onClick={clearError}
              className="text-gray-400 hover:text-white ml-4"
            >
              <X className="w-4 h-4" />
            </button>
          </div>
        </motion.div>
      )}

      {/* 主内容区域 */}
      <div className="flex h-screen pt-16">
        {/* 侧边栏 - 统计信息 */}
        <motion.div
          variants={containerVariants}
          initial="hidden"
          animate="visible"
          className={clsx(
            'w-80 border-r border-gray-800 p-6 overflow-y-auto transition-all duration-300',
            sidebarCollapsed ? 'w-20' : 'w-80'
          )}
        >
          <div className="space-y-6">
            {/* 微信状态卡片 */}
            <motion.div variants={itemVariants} className="card hover:scale-105 transition-transform">
              <div className="flex items-center justify-between mb-4">
                <div className="flex items-center space-x-3">
                  <div
                    className={clsx(
                      'w-3 h-3 rounded-full',
                      weChatStatus?.is_connected
                        ? 'bg-success-500 animate-pulse'
                        : 'bg-danger-500 animate-pulse'
                    )}
                  />
                  <h3 className="text-lg font-semibold text-white">微信连接</h3>
                </div>
                {lastRefreshTimes.wechat && (
                  <div className="text-xs text-gray-400">{lastRefreshTimes.wechat}</div>
                )}
              </div>

              <div className="space-y-3">
                <div className="flex justify-between items-center p-3 bg-gray-800/50 rounded">
                  <span className="text-sm text-gray-400">连接状态</span>
                  <span
                    className={clsx(
                      'px-2 py-1 rounded-full text-xs font-medium',
                      weChatStatus?.is_connected
                        ? 'bg-success-100 text-success-800'
                        : 'bg-danger-100 text-danger-800'
                    )}
                  >
                    {weChatStatus?.is_connected ? '已连接' : '未连接'}
                  </span>
                </div>

                <div className="flex justify-between items-center p-3 bg-gray-800/50 rounded">
                  <span className="text-sm text-gray-400">消息数量</span>
                  <span className="text-lg font-bold text-primary-400">
                    {weChatStatus?.message_count || 0}
                  </span>
                </div>

                <div className="flex justify-between items-center p-3 bg-gray-800/50 rounded">
                  <span className="text-sm text-gray-400">联系人</span>
                  <span className="text-lg font-bold text-accent-400">
                    {weChatStatus?.contact_count || 0}
                  </span>
                </div>
              </div>
            </motion.div>

            {/* 快速操作 */}
            <motion.div variants={itemVariants} className="card hover:scale-105 transition-transform">
              <div className="flex items-center space-x-3 mb-4">
                <div className="w-8 h-8 bg-gradient-to-br from-primary-500 to-accent-500 rounded-lg flex items-center justify-center animate-float">
                  <Database className="w-4 h-4 text-white" />
                </div>
                <h3 className="text-lg font-semibold text-white">快速操作</h3>
              </div>

              <div className="space-y-3">
                <button
                  onClick={() => handleRefresh('wechat')}
                  disabled={refreshStates.wechat || !weChatStatus?.is_connected}
                  className={clsx(
                    'w-full btn btn-secondary flex items-center justify-center space-x-2',
                    refreshStates.wechat || !weChatStatus?.is_connected
                      ? 'opacity-50 cursor-not-allowed'
                      : ''
                  )}
                >
                  {refreshStates.wechat ? (
                    <LoadingSpinner size="sm" text="同步中..." />
                  ) : (
                    <>
                      <Database className="w-4 h-4" />
                      <span>同步数据</span>
                    </>
                  )}
                </button>
              </div>
            </motion.div>

            {/* 统计信息 */}
            <motion.div variants={itemVariants} className="card hover:scale-105 transition-transform">
              <div className="flex items-center space-x-3 mb-4">
                <div className="w-8 h-8 bg-gradient-to-br from-success-500 to-emerald-600 rounded-lg flex items-center justify-center animate-float">
                  <Brain className="w-4 h-4 text-white" />
                </div>
                <h3 className="text-lg font-semibold text-white">今日统计</h3>
              </div>

              <div className="space-y-4">
                <div className="flex justify-between items-center p-3 bg-gray-800/50 rounded hover:bg-gray-800/70 transition-colors">
                  <div className="flex items-center space-x-3">
                    <Brain className="w-5 h-5 text-primary-400" />
                    <span className="text-sm text-gray-300">情报简报</span>
                  </div>
                  <span className="text-lg font-bold text-primary-400">
                    {briefing ? 1 : 0}
                  </span>
                </div>

                <div className="flex justify-between items-center p-3 bg-gray-800/50 rounded hover:bg-gray-800/70 transition-colors">
                  <div className="flex items-center space-x-3">
                    <CheckSquare className="w-5 h-5 text-success-400" />
                    <span className="text-sm text-gray-300">决策待办</span>
                  </div>
                  <span className="text-lg font-bold text-success-400">
                    {todos?.todos?.length || 0}
                  </span>
                </div>

                <div className="flex justify-between items-center p-3 bg-gray-800/50 rounded hover:bg-gray-800/70 transition-colors">
                  <div className="flex items-center space-x-3">
                    <Users className="w-5 h-5 text-accent-400" />
                    <span className="text-sm text-gray-300">人脉动态</span>
                  </div>
                  <span className="text-lg font-bold text-accent-400">
                    {connections?.connections?.length || 0}
                  </span>
                </div>
              </div>
            </motion.div>
          </div>
        </motion.div>

        {/* 主要内容区域 */}
        <div className="flex-1 overflow-hidden">
          <div className="h-full overflow-y-auto p-6">
            {/* 根据当前视图渲染内容 */}
            <motion.div
              key={activeView}
              initial="hidden"
              animate="visible"
              transition={{ duration: 0.3 }}
              variants={containerVariants}
            >
              {activeView === 'dashboard' && (
                <motion.div variants={itemVariants} className="space-y-6">
                  <h2 className="text-2xl font-bold text-white mb-6">仪表盘概览</h2>

                  <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                    <BriefingCard
                      briefing={briefing}
                      isLoading={refreshStates.briefing}
                      lastRefreshTime={lastRefreshTimes.briefing}
                      onRefresh={() => handleRefresh('briefing')}
                    />

                    {/* 其他统计卡片 */}
                    <div className="card">
                      <h3 className="text-lg font-semibold text-white mb-4">系统状态</h3>
                      <div className="space-y-4">
                        <div className="flex justify-between items-center">
                          <span className="text-gray-400">AI分析引擎</span>
                          <span className="text-success-400">运行中</span>
                        </div>
                        <div className="flex justify-between items-center">
                          <span className="text-gray-400">数据同步</span>
                          <span className="text-primary-400">正常</span>
                        </div>
                        <div className="flex justify-between items-center">
                          <span className="text-gray-400">系统性能</span>
                          <span className="text-success-400">优秀</span>
                        </div>
                      </div>
                    </div>
                  </div>
                </motion.div>
              )}

              {activeView === 'briefing' && (
                <motion.div variants={itemVariants}>
                  <BriefingCard
                    briefing={briefing}
                    isLoading={refreshStates.briefing}
                    lastRefreshTime={lastRefreshTimes.briefing}
                    onRefresh={() => handleRefresh('briefing')}
                  />
                </motion.div>
              )}

              {activeView === 'todos' && (
                <motion.div variants={itemVariants}>
                  <h2 className="text-2xl font-bold text-white mb-6">决策待办</h2>
                  {todos ? (
                    <div className="space-y-4">
                      {todos.todos.map((todo, index) => (
                        <motion.div
                          key={todo.id}
                          initial={{ opacity: 0, x: -20 }}
                          animate={{ opacity: 1, x: 0 }}
                          transition={{ delay: index * 0.1 }}
                          className="card"
                        >
                          <div className="flex items-center justify-between">
                            <div className="flex items-center space-x-3">
                              <input
                                type="checkbox"
                                checked={todo.status === 'completed'}
                                className="w-4 h-4 text-primary-500 bg-gray-800 border-gray-600 rounded focus:ring-primary-500 focus:ring-2"
                                readOnly
                              />
                              <div>
                                <p className="text-white font-medium">{todo.content}</p>
                                {todo.assignee && (
                                  <p className="text-xs text-gray-400">负责人: {todo.assignee}</p>
                                )}
                              </div>
                            </div>
                            <span
                              className={clsx(
                                'px-2 py-1 rounded text-xs font-medium',
                                todo.priority === 'high'
                                  ? 'bg-danger-100 text-danger-800'
                                  : todo.priority === 'medium'
                                  ? 'bg-warning-100 text-warning-800'
                                  : 'bg-gray-100 text-gray-800'
                              )}
                            >
                              {todo.priority === 'high' ? '高' : todo.priority === 'medium' ? '中' : '低'}
                            </span>
                          </div>
                        </motion.div>
                      ))}
                    </div>
                  ) : (
                    <div className="card text-center py-12">
                      <CheckSquare className="w-16 h-16 text-gray-600 mx-auto mb-4" />
                      <h3 className="text-xl font-semibold text-gray-400 mb-2">
                        暂无待办事项
                      </h3>
                      <button
                        onClick={() => handleRefresh('todos')}
                        className="btn btn-primary mt-4"
                      >
                        提取待办事项
                      </button>
                    </div>
                  )}
                </motion.div>
              )}

              {activeView === 'connections' && (
                <motion.div variants={itemVariants}>
                  <h2 className="text-2xl font-bold text-white mb-6">人脉雷达</h2>
                  {connections ? (
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                      {connections.connections.map((connection, index) => (
                        <motion.div
                          key={connection.id}
                          initial={{ opacity: 0, scale: 0.9 }}
                          animate={{ opacity: 1, scale: 1 }}
                          transition={{ delay: index * 0.1 }}
                          className="card hover:scale-105 transition-transform"
                        >
                          <div className="flex items-center justify-between mb-3">
                            <h4 className="font-semibold text-white">{connection.person_name}</h4>
                            <div className="flex items-center space-x-1">
                              {[...Array(5)].map((_, i) => (
                                <div
                                  key={i}
                                  className={clsx(
                                    'w-2 h-2 rounded-full',
                                    i < connection.strength / 2
                                      ? 'bg-success-400'
                                      : 'bg-gray-600'
                                  )}
                                />
                              ))}
                            </div>
                          </div>
                          <div className="space-y-2 text-sm">
                            <div className="flex justify-between">
                              <span className="text-gray-400">关系</span>
                              <span className="text-gray-300">{connection.relationship}</span>
                            </div>
                            <div className="flex justify-between">
                              <span className="text-gray-400">互动次数</span>
                              <span className="text-gray-300">{connection.interaction_count}</span>
                            </div>
                            <div className="flex justify-between">
                              <span className="text-gray-400">情感倾向</span>
                              <span
                                className={clsx(
                                  connection.sentiment === 'positive'
                                    ? 'text-success-400'
                                    : connection.sentiment === 'negative'
                                    ? 'text-danger-400'
                                    : 'text-gray-400'
                                )}
                              >
                                {connection.sentiment === 'positive' ? '积极' : connection.sentiment === 'negative' ? '消极' : '中性'}
                              </span>
                            </div>
                          </div>
                        </motion.div>
                      ))}
                    </div>
                  ) : (
                    <div className="card text-center py-12">
                      <Users className="w-16 h-16 text-gray-600 mx-auto mb-4" />
                      <h3 className="text-xl font-semibold text-gray-400 mb-2">
                        暂无人脉数据
                      </h3>
                      <button
                        onClick={() => handleRefresh('connections')}
                        className="btn btn-primary mt-4"
                      >
                        分析人脉关系
                      </button>
                    </div>
                  )}
                </motion.div>
              )}

              {activeView === 'settings' && (
                <motion.div variants={itemVariants}>
                  <h2 className="text-2xl font-bold text-white mb-6">系统设置</h2>
                  <div className="card">
                    <p className="text-gray-400">设置功能正在开发中...</p>
                  </div>
                </motion.div>
              )}
            </motion.div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;