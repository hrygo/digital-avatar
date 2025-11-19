import React, { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import { Brain, Users, CheckSquare, Activity, Database, AlertCircle, RefreshCw, Clock } from 'lucide-react';
import { useAppStore } from '../store/useAppStore';
import { LoadingSpinner, LoadingCard, FullScreenLoader } from '../components/LoadingSpinner';

const Dashboard: React.FC = () => {
  const {
    weChatStatus,
    isConnecting,
    isSyncing,
    isAnalyzing,
    briefing,
    todos,
    connections,
    error,
    initializeApp,
    syncWeChat,
    generateBriefing,
    extractTodos,
    analyzeConnections,
    clearError,
  } = useAppStore();

  const [activeTab, setActiveTab] = useState<'briefing' | 'todos' | 'connections'>('briefing');
  const [isInitialLoading, setIsInitialLoading] = useState(true);

  // Individual refresh states
  const [refreshStates, setRefreshStates] = useState({
    briefing: false,
    todos: false,
    connections: false,
    wechat: false,
    all: false
  });

  // Last refresh times
  const [lastRefreshTimes, setLastRefreshTimes] = useState({
    briefing: null as string | null,
    todos: null as string | null,
    connections: null as string | null,
    wechat: null as string | null,
    all: null as string | null
  });

  useEffect(() => {
    const initialize = async () => {
      setIsInitialLoading(true);
      try {
        await initializeApp();
      } finally {
        // 延迟隐藏加载状态，让用户看到加载效果
        setTimeout(() => {
          setIsInitialLoading(false);
        }, 1000);
      }
    };

    initialize();
  }, [initializeApp]);

  const updateRefreshTime = (type: keyof typeof lastRefreshTimes) => {
    setLastRefreshTimes(prev => ({
      ...prev,
      [type]: new Date().toLocaleTimeString('zh-CN', {
        hour: '2-digit',
        minute: '2-digit'
      })
    }));
  };

  const setRefreshState = (type: keyof typeof refreshStates, loading: boolean) => {
    setRefreshStates(prev => ({
      ...prev,
      [type]: loading
    }));
  };

  const handleRefreshBriefing = async () => {
    if (refreshStates.briefing) return;

    setRefreshState('briefing', true);
    try {
      await generateBriefing();
      updateRefreshTime('briefing');
    } catch (error) {
      console.error('Failed to refresh briefing:', error);
    } finally {
      setRefreshState('briefing', false);
    }
  };

  const handleRefreshTodos = async () => {
    if (refreshStates.todos) return;

    setRefreshState('todos', true);
    try {
      await extractTodos();
      updateRefreshTime('todos');
    } catch (error) {
      console.error('Failed to refresh todos:', error);
    } finally {
      setRefreshState('todos', false);
    }
  };

  const handleRefreshConnections = async () => {
    if (refreshStates.connections) return;

    setRefreshState('connections', true);
    try {
      await analyzeConnections();
      updateRefreshTime('connections');
    } catch (error) {
      console.error('Failed to refresh connections:', error);
    } finally {
      setRefreshState('connections', false);
    }
  };

  const handleRefreshWeChat = async () => {
    if (refreshStates.wechat || !weChatStatus?.is_connected) return;

    setRefreshState('wechat', true);
    try {
      await syncWeChat();
      updateRefreshTime('wechat');
    } catch (error) {
      console.error('Failed to refresh WeChat data:', error);
    } finally {
      setRefreshState('wechat', false);
    }
  };

  const handleRefreshAll = async () => {
    if (refreshStates.all) return;

    setRefreshState('all', true);
    try {
      await Promise.all([
        generateBriefing(),
        extractTodos(),
        analyzeConnections(),
      ]);
      updateRefreshTime('all');
    } catch (error) {
      console.error('Failed to refresh all data:', error);
    } finally {
      setRefreshState('all', false);
    }
  };

  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1
      }
    }
  };

  const itemVariants = {
    hidden: { y: 20, opacity: 0 },
    visible: {
      y: 0,
      opacity: 1,
      transition: {
        type: "spring",
        stiffness: 100
      }
    }
  };

  // 显示初始加载全屏加载器
  if (isInitialLoading) {
    return (
      <FullScreenLoader
        text="正在初始化 TwinOS"
        subtext="请稍候，正在连接智能分析引擎..."
      />
    );
  }

  return (
    <div className="min-h-screen bg-black text-white">
      {/* 顶部导航栏 */}
      <motion.div
        initial={{ y: -50, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        className="border-b border-gray-800 bg-gradient-to-r from-gray-900/80 to-gray-800/60 backdrop-blur-xl shadow-2xl"
      >
        <div className="px-6 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <motion.div
                animate={{ rotate: [0, 10, -10, 0] }}
                transition={{ repeat: Infinity, duration: 4 }}
                className="w-10 h-10 bg-gradient-to-br from-blue-500 to-purple-600 rounded-xl flex items-center justify-center glow"
              >
                <Brain className="w-6 h-6 text-white" />
              </motion.div>
              <div>
                <h1 className="text-3xl font-bold text-gradient">TWINOS</h1>
                <p className="text-xs text-gray-400">智能决策副驾</p>
              </div>
            </div>

            <div className="flex items-center space-x-6">
              {/* 连接状态指示器 */}
              <motion.div
                whileHover={{ scale: 1.05 }}
                className="flex items-center space-x-3 px-4 py-2 rounded-full border border-gray-700 bg-gray-800/50"
              >
                <div className={`w-3 h-3 rounded-full ${
                  weChatStatus?.is_connected ? 'bg-green-500 glow-green animate-pulse' : 'bg-red-500 glow-red animate-pulse'
                }`} />
                <span className="text-sm font-medium text-gray-300">
                  {weChatStatus?.is_connected ? '微信已连接' : '微信未连接'}
                </span>
              </motion.div>

              {/* 刷新按钮组 */}
              <div className="flex items-center space-x-3">
                {/* 单独刷新下拉菜单 */}
                <div className="relative group">
                  <motion.button
                    whileHover={{ scale: 1.05 }}
                    whileTap={{ scale: 0.95 }}
                    className="btn-secondary flex items-center space-x-2"
                  >
                    <RefreshCw className="w-4 h-4" />
                    <span>单独刷新</span>
                  </motion.button>

                  {/* 下拉菜单 */}
                  <div className="absolute right-0 mt-2 w-48 glass-card rounded-lg shadow-xl opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 z-50">
                    <div className="p-2 space-y-1">
                      <motion.button
                        whileHover={{ scale: 1.02 }}
                        whileTap={{ scale: 0.98 }}
                        onClick={handleRefreshBriefing}
                        disabled={refreshStates.briefing}
                        className={`w-full text-left px-3 py-2 rounded-lg flex items-center justify-between text-sm ${
                          refreshStates.briefing
                            ? 'opacity-50 cursor-not-allowed'
                            : 'hover:bg-gray-700/50 transition-colors'
                        }`}
                      >
                        <span>刷新简报</span>
                        <div className="flex items-center space-x-2">
                          {refreshStates.briefing && <RefreshCw className="w-3 h-3 animate-spin" />}
                          {lastRefreshTimes.briefing && (
                            <span className="text-xs text-gray-400">{lastRefreshTimes.briefing}</span>
                          )}
                        </div>
                      </motion.button>

                      <motion.button
                        whileHover={{ scale: 1.02 }}
                        whileTap={{ scale: 0.98 }}
                        onClick={handleRefreshTodos}
                        disabled={refreshStates.todos}
                        className={`w-full text-left px-3 py-2 rounded-lg flex items-center justify-between text-sm ${
                          refreshStates.todos
                            ? 'opacity-50 cursor-not-allowed'
                            : 'hover:bg-gray-700/50 transition-colors'
                        }`}
                      >
                        <span>刷新任务</span>
                        <div className="flex items-center space-x-2">
                          {refreshStates.todos && <RefreshCw className="w-3 h-3 animate-spin" />}
                          {lastRefreshTimes.todos && (
                            <span className="text-xs text-gray-400">{lastRefreshTimes.todos}</span>
                          )}
                        </div>
                      </motion.button>

                      <motion.button
                        whileHover={{ scale: 1.02 }}
                        whileTap={{ scale: 0.98 }}
                        onClick={handleRefreshConnections}
                        disabled={refreshStates.connections}
                        className={`w-full text-left px-3 py-2 rounded-lg flex items-center justify-between text-sm ${
                          refreshStates.connections
                            ? 'opacity-50 cursor-not-allowed'
                            : 'hover:bg-gray-700/50 transition-colors'
                        }`}
                      >
                        <span>刷新关系</span>
                        <div className="flex items-center space-x-2">
                          {refreshStates.connections && <RefreshCw className="w-3 h-3 animate-spin" />}
                          {lastRefreshTimes.connections && (
                            <span className="text-xs text-gray-400">{lastRefreshTimes.connections}</span>
                          )}
                        </div>
                      </motion.button>

                      <motion.button
                        whileHover={{ scale: 1.02 }}
                        whileTap={{ scale: 0.98 }}
                        onClick={handleRefreshWeChat}
                        disabled={refreshStates.wechat || !weChatStatus?.is_connected}
                        className={`w-full text-left px-3 py-2 rounded-lg flex items-center justify-between text-sm ${
                          refreshStates.wechat || !weChatStatus?.is_connected
                            ? 'opacity-50 cursor-not-allowed'
                            : 'hover:bg-gray-700/50 transition-colors'
                        }`}
                      >
                        <span>同步微信</span>
                        <div className="flex items-center space-x-2">
                          {refreshStates.wechat && <RefreshCw className="w-3 h-3 animate-spin" />}
                          {lastRefreshTimes.wechat && (
                            <span className="text-xs text-gray-400">{lastRefreshTimes.wechat}</span>
                          )}
                        </div>
                      </motion.button>
                    </div>
                  </div>
                </div>

                {/* 全部刷新按钮 */}
                <motion.button
                  whileHover={{ scale: 1.05 }}
                  whileTap={{ scale: 0.95 }}
                  onClick={handleRefreshAll}
                  disabled={refreshStates.all}
                  className={`btn-primary flex items-center space-x-2 ${
                    refreshStates.all ? 'opacity-50 cursor-not-allowed' : ''
                  }`}
                >
                  <Activity className={`w-4 h-4 ${refreshStates.all ? 'animate-spin' : ''}`} />
                  <span>{refreshStates.all ? '刷新中...' : '全部刷新'}</span>
                  {lastRefreshTimes.all && (
                    <span className="text-xs opacity-75">{lastRefreshTimes.all}</span>
                  )}
                </motion.button>
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
          className="bg-red-500/20 border border-red-500/50 text-red-200 px-4 py-3 m-4 rounded-lg flex items-center space-x-2"
        >
          <AlertCircle className="w-5 h-5" />
          <span>{error}</span>
          <button onClick={clearError} className="ml-auto hover:text-white">×</button>
        </motion.div>
      )}

      {/* 主内容区域 */}
      <div className="flex h-[calc(100vh-80px)]">
        {/* 侧边栏统计信息 */}
        <motion.div
          variants={containerVariants}
          initial="hidden"
          animate="visible"
          className="w-80 border-r border-gray-800 p-6 overflow-y-auto"
        >
          <div className="space-y-6">
            {/* 微信状态卡片 */}
            <motion.div variants={itemVariants} className="glass-card rounded-xl p-6 hover:scale-105 transition-transform duration-200">
              <div className="flex items-center justify-between mb-4">
                <div className="flex items-center space-x-3">
                  <div className={`w-3 h-3 rounded-full ${weChatStatus?.is_connected ? 'bg-green-500 glow-green animate-pulse' : 'bg-red-500 glow-red'} `} />
                  <h3 className="text-lg font-semibold text-white">微信连接状态</h3>
                  {lastRefreshTimes.wechat && (
                    <div className="flex items-center space-x-1 text-xs text-gray-400">
                      <Clock className="w-3 h-3" />
                      <span>{lastRefreshTimes.wechat}</span>
                    </div>
                  )}
                </div>

                {weChatStatus?.is_connected && (
                  <motion.button
                    whileHover={{ scale: 1.05 }}
                    whileTap={{ scale: 0.95 }}
                    onClick={handleRefreshWeChat}
                    disabled={refreshStates.wechat}
                    className={`btn-secondary flex items-center space-x-1 px-3 py-1 text-sm ${
                      refreshStates.wechat ? 'opacity-50 cursor-not-allowed' : ''
                    }`}
                  >
                    <RefreshCw className={`w-3 h-3 ${refreshStates.wechat ? 'animate-spin' : ''}`} />
                    <span>{refreshStates.wechat ? '同步中' : '同步'}</span>
                  </motion.button>
                )}
              </div>

              <div className="space-y-3">
                <div className="flex justify-between items-center p-3 bg-gray-800/50 rounded-lg">
                  <span className="text-sm text-gray-400">连接状态</span>
                  <span className={`px-2 py-1 rounded-full text-xs font-medium ${
                    weChatStatus?.is_connected
                      ? 'status-online'
                      : 'status-offline'
                  }`}>
                    {weChatStatus?.is_connected ? '已连接' : '未连接'}
                  </span>
                </div>

                <div className="flex justify-between items-center p-3 bg-gray-800/50 rounded-lg">
                  <span className="text-sm text-gray-400">消息数量</span>
                  <span className="text-lg font-bold text-blue-400">{weChatStatus?.message_count || 0}</span>
                </div>

                <div className="flex justify-between items-center p-3 bg-gray-800/50 rounded-lg">
                  <span className="text-sm text-gray-400">联系人数量</span>
                  <span className="text-lg font-bold text-purple-400">{weChatStatus?.contact_count || 0}</span>
                </div>
              </div>

              {!weChatStatus?.is_connected && (
                <motion.button
                  whileHover={{ scale: 1.05 }}
                  whileTap={{ scale: 0.95 }}
                  className="btn-primary mt-4 w-full"
                >
                  连接微信
                </motion.button>
              )}
            </motion.div>

            {/* 快速操作 */}
            <motion.div variants={itemVariants} className="glass-card rounded-xl p-6 hover:scale-105 transition-transform duration-200">
              <div className="flex items-center space-x-3 mb-4">
                <div className="w-8 h-8 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center animate-float">
                  <Activity className="w-4 h-4 text-white" />
                </div>
                <h3 className="text-lg font-semibold text-white">快速操作</h3>
              </div>

              <div className="space-y-3">
                <motion.button
                  whileHover={{ scale: 1.02 }}
                  whileTap={{ scale: 0.98 }}
                  onClick={syncWeChat}
                  disabled={isSyncing || !weChatStatus?.is_connected}
                  className={`w-full px-4 py-3 rounded-lg flex items-center justify-center space-x-2 font-medium transition-all duration-200 ${
                    isSyncing || !weChatStatus?.is_connected
                      ? 'bg-gray-800/50 text-gray-500 cursor-not-allowed'
                      : 'btn-secondary'
                  }`}
                >
                  {isSyncing ? (
                    <LoadingSpinner size="sm" color="secondary" text="同步中..." />
                  ) : (
                    <>
                      <Database className="w-4 h-4" />
                      <span>同步数据</span>
                    </>
                  )}
                </motion.button>
              </div>
            </motion.div>

            {/* 统计信息 */}
            <motion.div variants={itemVariants} className="glass-card rounded-xl p-6 hover:scale-105 transition-transform duration-200">
              <div className="flex items-center space-x-3 mb-4">
                <div className="w-8 h-8 bg-gradient-to-br from-green-500 to-emerald-600 rounded-lg flex items-center justify-center animate-float">
                  <Brain className="w-4 h-4 text-white" />
                </div>
                <h3 className="text-lg font-semibold text-white">今日统计</h3>
              </div>

              <div className="space-y-4">
                <motion.div
                  whileHover={{ x: 5 }}
                  className="flex justify-between items-center p-3 bg-gray-800/50 rounded-lg hover:bg-gray-800/70 transition-colors"
                >
                  <div className="flex items-center space-x-3">
                    <Brain className="w-5 h-5 text-blue-400" />
                    <span className="text-sm text-gray-300">情报简报</span>
                  </div>
                  <span className="text-lg font-bold text-blue-400">{briefing?.message_count || 0}</span>
                </motion.div>

                <motion.div
                  whileHover={{ x: 5 }}
                  className="flex justify-between items-center p-3 bg-gray-800/50 rounded-lg hover:bg-gray-800/70 transition-colors"
                >
                  <div className="flex items-center space-x-3">
                    <CheckSquare className="w-5 h-5 text-green-400" />
                    <span className="text-sm text-gray-300">待办事项</span>
                  </div>
                  <span className="text-lg font-bold text-green-400">{todos?.todos?.length || 0}</span>
                </motion.div>

                <motion.div
                  whileHover={{ x: 5 }}
                  className="flex justify-between items-center p-3 bg-gray-800/50 rounded-lg hover:bg-gray-800/70 transition-colors"
                >
                  <div className="flex items-center space-x-3">
                    <Users className="w-5 h-5 text-purple-400" />
                    <span className="text-sm text-gray-300">人脉动态</span>
                  </div>
                  <span className="text-lg font-bold text-purple-400">{connections?.connections?.length || 0}</span>
                </motion.div>
              </div>
            </motion.div>
          </div>
        </motion.div>

        {/* 主要内容区域 */}
        <div className="flex-1 flex flex-col">
          {/* 标签页导航 */}
          <div className="border-b border-gray-800 bg-gray-900/50 backdrop-blur-sm">
            <div className="flex space-x-2 p-4">
              {[
                { id: 'briefing', label: '今日情报', icon: Brain, color: 'blue' },
                { id: 'todos', label: '决策待办', icon: CheckSquare, color: 'green' },
                { id: 'connections', label: '人脉雷达', icon: Users, color: 'purple' },
              ].map((tab, index) => (
                <motion.button
                  key={tab.id}
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: index * 0.1 }}
                  whileHover={{ scale: 1.05, y: -2 }}
                  whileTap={{ scale: 0.98 }}
                  onClick={() => setActiveTab(tab.id as any)}
                  className={`flex items-center space-x-3 px-6 py-3 rounded-xl transition-all duration-200 font-medium ${
                    activeTab === tab.id
                      ? `bg-gradient-to-r from-${tab.color}-600 to-${tab.color}-700 text-white shadow-lg glow`
                      : 'text-gray-400 hover:text-white hover:bg-gray-800/50'
                  }`}
                >
                  <motion.div
                    animate={{
                      rotate: activeTab === tab.id ? 360 : 0,
                      scale: activeTab === tab.id ? 1.1 : 1
                    }}
                    transition={{ duration: 0.3 }}
                  >
                    <tab.icon className={`w-5 h-5 ${
                      activeTab === tab.id ? 'text-white' : `text-${tab.color}-400`
                    }`} />
                  </motion.div>
                  <span>{tab.label}</span>
                  {activeTab === tab.id && (
                    <motion.div
                      layoutId="activeTab"
                      className="absolute bottom-0 left-0 right-0 h-1 bg-gradient-to-r from-blue-400 to-purple-600 rounded-t-full"
                    />
                  )}
                </motion.button>
              ))}
            </div>
          </div>

          {/* 内容面板 */}
          <div className="flex-1 overflow-y-auto p-6">
            {isAnalyzing && (
              <div className="flex items-center justify-center h-64">
                <div className="flex flex-col items-center space-y-4">
                  <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
                  <p className="text-gray-400">正在分析数据...</p>
                </div>
              </div>
            )}

            {!isAnalyzing && (
              <motion.div
                key={activeTab}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.3 }}
              >
                {activeTab === 'briefing' && (
                  <div className="max-w-4xl">
                    <motion.div
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      className="flex items-center justify-between mb-6"
                    >
                      <div className="flex items-center space-x-3">
                        <div className="w-8 h-8 bg-gradient-to-br from-blue-500 to-blue-600 rounded-lg flex items-center justify-center glow">
                          <Brain className="w-4 h-4 text-white" />
                        </div>
                        <h2 className="text-2xl font-bold text-white">今日情报简报</h2>
                        {lastRefreshTimes.briefing && (
                          <div className="flex items-center space-x-1 text-xs text-gray-400">
                            <Clock className="w-3 h-3" />
                            <span>{lastRefreshTimes.briefing}</span>
                          </div>
                        )}
                      </div>

                      <motion.button
                        whileHover={{ scale: 1.05 }}
                        whileTap={{ scale: 0.95 }}
                        onClick={handleRefreshBriefing}
                        disabled={refreshStates.briefing}
                        className={`btn-secondary flex items-center space-x-2 ${
                          refreshStates.briefing ? 'opacity-50 cursor-not-allowed' : ''
                        }`}
                      >
                        <RefreshCw className={`w-4 h-4 ${refreshStates.briefing ? 'animate-spin' : ''}`} />
                        <span>{refreshStates.briefing ? '刷新中...' : '刷新'}</span>
                      </motion.button>
                    </motion.div>

                    {refreshStates.briefing ? (
                      <LoadingCard
                        title="今日情报简报"
                        description="正在分析微信消息，生成 AI 情报简报..."
                        lines={5}
                      />
                    ) : briefing ? (
                      <motion.div
                        initial={{ opacity: 0, scale: 0.95 }}
                        animate={{ opacity: 1, scale: 1 }}
                        className="glass-card rounded-2xl p-8 hover:scale-[1.01] transition-transform duration-300"
                      >
                        <div className="mb-6">
                          <div className="flex items-center justify-between mb-4">
                            <div className="flex items-center space-x-2">
                              <div className="w-3 h-3 bg-blue-500 rounded-full animate-pulse" />
                              <span className="text-sm text-blue-400 font-medium">AI 分析结果</span>
                            </div>
                            <div className="flex items-center space-x-4 text-sm text-gray-400">
                              <span>📊 {briefing.message_count || 0} 条消息</span>
                              <span>🕐 {new Date(briefing.generated_at).toLocaleString()}</span>
                            </div>
                          </div>
                        </div>

                        <div className="prose prose-invert max-w-none">
                          <pre className="whitespace-pre-wrap font-sans text-gray-300 leading-relaxed bg-gray-900/30 rounded-lg p-6 border border-gray-700/50">
                            {briefing.content}
                          </pre>
                        </div>

                        <div className="mt-6 pt-6 border-t border-gray-700">
                          <div className="flex items-center justify-between text-sm">
                            <span className="text-gray-400">由 TwinOS AI 引擎生成</span>
                            <motion.div
                              animate={{ opacity: [1, 0.5, 1] }}
                              transition={{ repeat: Infinity, duration: 2 }}
                              className="w-2 h-2 bg-green-500 rounded-full"
                            />
                          </div>
                        </div>
                      </motion.div>
                    ) : (
                      <motion.div
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        className="glass-card rounded-2xl p-16 text-center"
                      >
                        <motion.div
                          animate={{ y: [0, -10, 0] }}
                          transition={{ repeat: Infinity, duration: 3 }}
                          className="w-20 h-20 bg-gradient-to-br from-gray-700 to-gray-800 rounded-2xl flex items-center justify-center mx-auto mb-6"
                        >
                          <Brain className="w-10 h-10 text-gray-500" />
                        </motion.div>
                        <h3 className="text-xl font-semibold text-gray-400 mb-3">暂无情报简报</h3>
                        <p className="text-gray-500 mb-6">请先同步微信数据，然后生成情报简报</p>
                        <motion.button
                          whileHover={{ scale: 1.05 }}
                          whileTap={{ scale: 0.95 }}
                          onClick={generateBriefing}
                          className="btn-primary"
                        >
                          生成情报简报
                        </motion.button>
                      </motion.div>
                    )}
                  </div>
                )}

                {activeTab === 'todos' && (
                  <div className="max-w-4xl">
                    <motion.div
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      className="flex items-center justify-between mb-6"
                    >
                      <div className="flex items-center space-x-3">
                        <div className="w-8 h-8 bg-gradient-to-br from-green-500 to-emerald-600 rounded-lg flex items-center justify-center glow-green">
                          <CheckSquare className="w-4 h-4 text-white" />
                        </div>
                        <h2 className="text-2xl font-bold text-white">决策待办</h2>
                        {lastRefreshTimes.todos && (
                          <div className="flex items-center space-x-1 text-xs text-gray-400">
                            <Clock className="w-3 h-3" />
                            <span>{lastRefreshTimes.todos}</span>
                          </div>
                        )}
                      </div>

                      <motion.button
                        whileHover={{ scale: 1.05 }}
                        whileTap={{ scale: 0.95 }}
                        onClick={handleRefreshTodos}
                        disabled={refreshStates.todos}
                        className={`btn-secondary flex items-center space-x-2 ${
                          refreshStates.todos ? 'opacity-50 cursor-not-allowed' : ''
                        }`}
                      >
                        <RefreshCw className={`w-4 h-4 ${refreshStates.todos ? 'animate-spin' : ''}`} />
                        <span>{refreshStates.todos ? '刷新中...' : '刷新'}</span>
                      </motion.button>
                    </motion.div>

                    {refreshStates.todos ? (
                      <LoadingCard
                        title="决策待办"
                        description="正在从消息中提取待办事项和决策要点..."
                        lines={4}
                      />
                    ) : todos && todos.todos.length > 0 ? (
                      <div className="space-y-4">
                        {todos.todos.map((todo, index) => (
                          <motion.div
                            key={index}
                            initial={{ opacity: 0, x: -20 }}
                            animate={{ opacity: 1, x: 0 }}
                            whileHover={{ scale: 1.02, x: 10 }}
                            transition={{ delay: index * 0.1 }}
                            className="glass-card rounded-2xl p-6 border-l-4 hover:border-l-8 transition-all duration-300"
                            style={{
                              borderLeftColor: todo.priority === 'high' ? '#ef4444' :
                                             todo.priority === 'medium' ? '#eab308' : '#22c55e'
                            }}
                          >
                            <div className="flex items-start space-x-4">
                              <motion.div
                                whileHover={{ rotate: 15, scale: 1.2 }}
                                className="w-12 h-12 bg-gradient-to-br from-green-500/20 to-emerald-600/20 rounded-xl flex items-center justify-center mt-1"
                              >
                                <CheckSquare className="w-6 h-6 text-green-400" />
                              </motion.div>
                              <div className="flex-1">
                                <div className="flex items-center justify-between mb-3">
                                  <h3 className="text-xl font-semibold text-white">{todo.title}</h3>
                                  <motion.span
                                    whileHover={{ scale: 1.1 }}
                                    className={`px-3 py-1 rounded-full text-xs font-bold shadow-lg ${
                                      todo.priority === 'high'
                                        ? 'bg-gradient-to-r from-red-500 to-pink-500 text-white glow-red'
                                        : todo.priority === 'medium'
                                        ? 'bg-gradient-to-r from-yellow-500 to-orange-500 text-white'
                                        : 'bg-gradient-to-r from-green-500 to-emerald-500 text-white glow-green'
                                    }`}
                                  >
                                    {todo.priority === 'high' ? '🔥 高优先级' :
                                     todo.priority === 'medium' ? '⚡ 中优先级' : '✓ 低优先级'}
                                  </motion.span>
                                </div>
                                <p className="text-gray-300 mb-4 leading-relaxed">{todo.description}</p>
                                <div className="flex items-center justify-between">
                                  {todo.deadline && (
                                    <div className="flex items-center space-x-2 text-sm text-gray-400">
                                      <span>⏰</span>
                                      <span>截止: {todo.deadline}</span>
                                    </div>
                                  )}
                                  <motion.button
                                    whileHover={{ scale: 1.05 }}
                                    whileTap={{ scale: 0.95 }}
                                    className="text-sm text-blue-400 hover:text-blue-300 font-medium"
                                  >
                                    查看详情 →
                                  </motion.button>
                                </div>
                              </div>
                            </div>
                          </motion.div>
                        ))}
                      </div>
                    ) : (
                      <motion.div
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        className="glass-card rounded-2xl p-16 text-center"
                      >
                        <motion.div
                          animate={{ rotate: [0, 5, -5, 0] }}
                          transition={{ repeat: Infinity, duration: 4 }}
                          className="w-20 h-20 bg-gradient-to-br from-gray-700 to-gray-800 rounded-2xl flex items-center justify-center mx-auto mb-6"
                        >
                          <CheckSquare className="w-10 h-10 text-gray-500" />
                        </motion.div>
                        <h3 className="text-xl font-semibold text-gray-400 mb-3">暂无待办事项</h3>
                        <p className="text-gray-500 mb-6">请先同步微信数据，然后提取待办事项</p>
                        <motion.button
                          whileHover={{ scale: 1.05 }}
                          whileTap={{ scale: 0.95 }}
                          onClick={extractTodos}
                          className="btn-primary"
                        >
                          提取待办事项
                        </motion.button>
                      </motion.div>
                    )}
                  </div>
                )}

                {activeTab === 'connections' && (
                  <div className="max-w-4xl">
                    <motion.div
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      className="flex items-center justify-between mb-6"
                    >
                      <div className="flex items-center space-x-3">
                        <div className="w-8 h-8 bg-gradient-to-br from-purple-500 to-pink-600 rounded-lg flex items-center justify-center glow-purple">
                          <Users className="w-4 h-4 text-white" />
                        </div>
                        <h2 className="text-2xl font-bold text-white">人脉雷达</h2>
                        {lastRefreshTimes.connections && (
                          <div className="flex items-center space-x-1 text-xs text-gray-400">
                            <Clock className="w-3 h-3" />
                            <span>{lastRefreshTimes.connections}</span>
                          </div>
                        )}
                      </div>

                      <motion.button
                        whileHover={{ scale: 1.05 }}
                        whileTap={{ scale: 0.95 }}
                        onClick={handleRefreshConnections}
                        disabled={refreshStates.connections}
                        className={`btn-secondary flex items-center space-x-2 ${
                          refreshStates.connections ? 'opacity-50 cursor-not-allowed' : ''
                        }`}
                      >
                        <RefreshCw className={`w-4 h-4 ${refreshStates.connections ? 'animate-spin' : ''}`} />
                        <span>{refreshStates.connections ? '刷新中...' : '刷新'}</span>
                      </motion.button>
                    </motion.div>
                    {refreshStates.connections ? (
                      <LoadingCard
                        title="人脉雷达"
                        description="正在分析消息中的社交网络和人际关系..."
                        lines={3}
                      />
                    ) : connections && connections.connections.length > 0 ? (
                      <div className="grid gap-4 md:grid-cols-2">
                        {connections.connections.map((connection, index) => (
                          <motion.div
                            key={index}
                            initial={{ opacity: 0, scale: 0.9 }}
                            animate={{ opacity: 1, scale: 1 }}
                            transition={{ delay: index * 0.1 }}
                            className="glass rounded-xl p-6"
                          >
                            <div className="flex items-center space-x-3 mb-4">
                              <div className="w-12 h-12 bg-purple-500/20 rounded-full flex items-center justify-center">
                                <Users className="w-6 h-6 text-purple-400" />
                              </div>
                              <div>
                                <h3 className="text-lg font-medium">{connection.person}</h3>
                                <p className="text-sm text-gray-400">{connection.action}</p>
                              </div>
                            </div>
                            <p className="text-gray-300 mb-3">{connection.context}</p>
                            <div className="flex items-center justify-between text-sm">
                              <span className={`px-2 py-1 rounded ${
                                connection.importance === 'high'
                                  ? 'bg-red-500/20 text-red-400'
                                  : connection.importance === 'medium'
                                  ? 'bg-yellow-500/20 text-yellow-400'
                                  : 'bg-green-500/20 text-green-400'
                              }`}>
                                {connection.importance === 'high' ? '重要' :
                                 connection.importance === 'medium' ? '一般' : '普通'}
                              </span>
                              <span className="text-gray-400">
                                {connection.message_count} 次提及
                              </span>
                            </div>
                          </motion.div>
                        ))}
                      </div>
                    ) : (
                      <div className="glass rounded-xl p-12 text-center">
                        <Users className="w-16 h-16 text-gray-600 mx-auto mb-4" />
                        <h3 className="text-lg font-medium text-gray-400 mb-2">暂无人脉动态</h3>
                        <p className="text-gray-500">请先同步微信数据，然后分析人脉关系</p>
                      </div>
                    )}
                  </div>
                )}
              </motion.div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;