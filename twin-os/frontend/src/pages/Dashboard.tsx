import React, { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import { Brain, Users, CheckSquare, Activity, Database, AlertCircle } from 'lucide-react';
import { useAppStore } from '../store/useAppStore';

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

  useEffect(() => {
    initializeApp();
  }, [initializeApp]);

  const handleRefreshAll = async () => {
    try {
      await Promise.all([
        generateBriefing(),
        extractTodos(),
        analyzeConnections(),
      ]);
    } catch (error) {
      console.error('Failed to refresh data:', error);
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

  return (
    <div className="min-h-screen bg-black text-white">
      {/* 顶部导航栏 */}
      <motion.div
        initial={{ y: -50, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        className="border-b border-gray-800 bg-gray-900/50 backdrop-blur-sm"
      >
        <div className="px-6 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <Brain className="w-8 h-8 text-blue-500" />
              <h1 className="text-2xl font-bold bg-gradient-to-r from-blue-400 to-purple-600 bg-clip-text text-transparent">
                TWINOS
              </h1>
            </div>

            <div className="flex items-center space-x-4">
              {/* 连接状态指示器 */}
              <div className="flex items-center space-x-2">
                <div className={`w-2 h-2 rounded-full ${weChatStatus?.is_connected ? 'bg-green-500' : 'bg-red-500'} animate-pulse`} />
                <span className="text-sm text-gray-400">
                  {weChatStatus?.is_connected ? '微信已连接' : '微信未连接'}
                </span>
              </div>

              {/* 刷新按钮 */}
              <motion.button
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
                onClick={handleRefreshAll}
                disabled={isAnalyzing || isSyncing}
                className="px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-700 rounded-lg flex items-center space-x-2 transition-colors"
              >
                <Activity className="w-4 h-4" />
                <span>刷新数据</span>
              </motion.button>
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
            <motion.div variants={itemVariants} className="glass rounded-xl p-4">
              <h3 className="text-sm font-medium text-gray-400 mb-3">微信连接状态</h3>
              <div className="space-y-2">
                <div className="flex justify-between">
                  <span className="text-sm">连接状态</span>
                  <span className={`text-sm ${weChatStatus?.is_connected ? 'text-green-400' : 'text-red-400'}`}>
                    {weChatStatus?.is_connected ? '已连接' : '未连接'}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-sm">消息数量</span>
                  <span className="text-sm">{weChatStatus?.message_count || 0}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-sm">联系人数量</span>
                  <span className="text-sm">{weChatStatus?.contact_count || 0}</span>
                </div>
              </div>

              {!weChatStatus?.is_connected && (
                <button className="mt-3 w-full px-3 py-2 bg-blue-600 hover:bg-blue-700 rounded-lg text-sm transition-colors">
                  连接微信
                </button>
              )}
            </motion.div>

            {/* 快速操作 */}
            <motion.div variants={itemVariants} className="glass rounded-xl p-4">
              <h3 className="text-sm font-medium text-gray-400 mb-3">快速操作</h3>
              <div className="space-y-2">
                <button
                  onClick={syncWeChat}
                  disabled={isSyncing || !weChatStatus?.is_connected}
                  className="w-full px-3 py-2 bg-gray-800 hover:bg-gray-700 disabled:bg-gray-900 disabled:text-gray-600 rounded-lg text-sm transition-colors flex items-center justify-center space-x-2"
                >
                  <Database className="w-4 h-4" />
                  <span>{isSyncing ? '同步中...' : '同步数据'}</span>
                </button>
              </div>
            </motion.div>

            {/* 统计信息 */}
            <motion.div variants={itemVariants} className="glass rounded-xl p-4">
              <h3 className="text-sm font-medium text-gray-400 mb-3">今日统计</h3>
              <div className="space-y-2">
                <div className="flex justify-between items-center">
                  <span className="text-sm flex items-center space-x-2">
                    <Brain className="w-4 h-4 text-blue-400" />
                    <span>情报简报</span>
                  </span>
                  <span className="text-sm text-blue-400">{briefing?.message_count || 0} 条消息</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm flex items-center space-x-2">
                    <CheckSquare className="w-4 h-4 text-green-400" />
                    <span>待办事项</span>
                  </span>
                  <span className="text-sm text-green-400">{todos?.todos?.length || 0} 项</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm flex items-center space-x-2">
                    <Users className="w-4 h-4 text-purple-400" />
                    <span>人脉动态</span>
                  </span>
                  <span className="text-sm text-purple-400">{connections?.connections?.length || 0} 条</span>
                </div>
              </div>
            </motion.div>
          </div>
        </motion.div>

        {/* 主要内容区域 */}
        <div className="flex-1 flex flex-col">
          {/* 标签页导航 */}
          <div className="border-b border-gray-800">
            <div className="flex space-x-1 p-2">
              {[
                { id: 'briefing', label: '今日情报', icon: Brain },
                { id: 'todos', label: '决策待办', icon: CheckSquare },
                { id: 'connections', label: '人脉雷达', icon: Users },
              ].map((tab) => (
                <motion.button
                  key={tab.id}
                  whileHover={{ scale: 1.02 }}
                  whileTap={{ scale: 0.98 }}
                  onClick={() => setActiveTab(tab.id as any)}
                  className={`flex items-center space-x-2 px-4 py-2 rounded-lg transition-all ${
                    activeTab === tab.id
                      ? 'bg-blue-600 text-white'
                      : 'text-gray-400 hover:text-white hover:bg-gray-800'
                  }`}
                >
                  <tab.icon className="w-4 h-4" />
                  <span>{tab.label}</span>
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
                    <h2 className="text-2xl font-bold mb-6">今日情报简报</h2>
                    {briefing ? (
                      <div className="glass rounded-xl p-6">
                        <div className="prose prose-invert max-w-none">
                          <pre className="whitespace-pre-wrap font-sans text-gray-300">
                            {briefing.content}
                          </pre>
                        </div>
                        <div className="mt-4 pt-4 border-t border-gray-700 text-sm text-gray-400">
                          生成时间: {new Date(briefing.generated_at).toLocaleString()}
                          {briefing.message_count && (
                            <span className="ml-4">基于 {briefing.message_count} 条消息</span>
                          )}
                        </div>
                      </div>
                    ) : (
                      <div className="glass rounded-xl p-12 text-center">
                        <Brain className="w-16 h-16 text-gray-600 mx-auto mb-4" />
                        <h3 className="text-lg font-medium text-gray-400 mb-2">暂无情报简报</h3>
                        <p className="text-gray-500">请先同步微信数据，然后生成情报简报</p>
                      </div>
                    )}
                  </div>
                )}

                {activeTab === 'todos' && (
                  <div className="max-w-4xl">
                    <h2 className="text-2xl font-bold mb-6">决策待办</h2>
                    {todos && todos.todos.length > 0 ? (
                      <div className="space-y-4">
                        {todos.todos.map((todo, index) => (
                          <motion.div
                            key={index}
                            initial={{ opacity: 0, x: -20 }}
                            animate={{ opacity: 1, x: 0 }}
                            transition={{ delay: index * 0.1 }}
                            className="glass rounded-xl p-6"
                          >
                            <div className="flex items-start space-x-4">
                              <CheckSquare className="w-5 h-5 text-green-400 mt-1 flex-shrink-0" />
                              <div className="flex-1">
                                <h3 className="text-lg font-medium mb-2">{todo.title}</h3>
                                <p className="text-gray-400 mb-3">{todo.description}</p>
                                <div className="flex items-center space-x-4 text-sm">
                                  <span className={`px-2 py-1 rounded ${
                                    todo.priority === 'high'
                                      ? 'bg-red-500/20 text-red-400'
                                      : todo.priority === 'medium'
                                      ? 'bg-yellow-500/20 text-yellow-400'
                                      : 'bg-green-500/20 text-green-400'
                                  }`}>
                                    {todo.priority === 'high' ? '高优先级' :
                                     todo.priority === 'medium' ? '中优先级' : '低优先级'}
                                  </span>
                                  {todo.deadline && (
                                    <span className="text-gray-400">
                                      截止: {todo.deadline}
                                    </span>
                                  )}
                                </div>
                              </div>
                            </div>
                          </motion.div>
                        ))}
                      </div>
                    ) : (
                      <div className="glass rounded-xl p-12 text-center">
                        <CheckSquare className="w-16 h-16 text-gray-600 mx-auto mb-4" />
                        <h3 className="text-lg font-medium text-gray-400 mb-2">暂无待办事项</h3>
                        <p className="text-gray-500">请先同步微信数据，然后提取待办事项</p>
                      </div>
                    )}
                  </div>
                )}

                {activeTab === 'connections' && (
                  <div className="max-w-4xl">
                    <h2 className="text-2xl font-bold mb-6">人脉雷达</h2>
                    {connections && connections.connections.length > 0 ? (
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