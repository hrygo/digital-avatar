import React from 'react';
import { create } from 'zustand';
import { devtools } from 'zustand/middleware';
import { api } from '../services/api';
import {
  AppState,
  AppStore,
  WeChatStatus,
  Briefing,
  TodoList,
  ConnectionList,
  AnalysisStatus,
  SystemSettings,
  SystemEvent,
  ApiResponse,
} from '../types';

// 默认设置
const DEFAULT_SETTINGS: SystemSettings = {
  auto_sync: true,
  sync_interval: 30,
  notification_enabled: true,
  notification_types: {
    briefing: true,
    urgent_todos: true,
    connections: false,
  },
  privacy: {
    data_retention_days: 365,
    anonymization_enabled: true,
    local_processing_only: true,
  },
  ui: {
    theme: 'dark',
    compact_mode: false,
    show_animations: true,
  },
};

// 创建应用状态store
export const useAppStore = create<AppStore>()(
  devtools(
    (set, get) => ({
      // 初始状态
      isLoading: true,
      isInitialized: false,
      error: null,

      // 数据状态
      weChatStatus: null,
      briefing: null,
      todos: null,
      connections: null,
      analysisStatus: null,
      settings: DEFAULT_SETTINGS,

      // UI状态
      activeView: 'dashboard',
      sidebarCollapsed: false,

      // 刷新状态
      refreshStates: {
        briefing: false,
        todos: false,
        connections: false,
        wechat: false,
      },

      // 时间戳
      lastRefreshTimes: {
        briefing: null,
        todos: null,
        connections: null,
        wechat: null,
      },

      // 初始化应用
      initializeApp: async () => {
        try {
          set({ isLoading: true, error: null });

          // 并行获取初始数据
          const [healthResult, settingsResult, wechatStatusResult] = await Promise.all([
            api.health(),
            api.getSettings(),
            api.getWeChatStatus(),
          ]);

          // 检查API健康状态
          if (healthResult.status === 'error') {
            throw new Error('无法连接到后端服务');
          }

          // 设置应用状态
          set({
            isInitialized: true,
            settings: settingsResult.status === 'success' && settingsResult.data &&
                     typeof settingsResult.data === 'object' && 'auto_sync' in settingsResult.data
                     ? settingsResult.data as SystemSettings
                     : DEFAULT_SETTINGS,
            weChatStatus: wechatStatusResult.status === 'success' && wechatStatusResult.data &&
                        typeof wechatStatusResult.data === 'object' && 'is_connected' in wechatStatusResult.data
                        ? wechatStatusResult.data as WeChatStatus
                        : null,
            isLoading: false,
            error: null,
          });

          console.log('✅ 应用初始化完成');
        } catch (error) {
          console.error('❌ 应用初始化失败:', error);
          set({
            error: error instanceof Error ? error.message : '初始化失败',
            isLoading: false,
          });
        }
      },

      // 清除错误
      clearError: () => set({ error: null }),

      // 设置当前视图
      setActiveView: (view) => set({ activeView: view }),

      // 切换侧边栏
      toggleSidebar: () => set((state) => ({ sidebarCollapsed: !state.sidebarCollapsed })),

      // 设置侧边栏状态
      setSidebarCollapsed: (collapsed) => set({ sidebarCollapsed: collapsed }),

      // 设置加载状态
      setLoading: (loading) => set({ isLoading: loading }),

      // 设置错误状态
      setError: (error) => set({ error }),

      // 更新刷新状态
      setRefreshState: (key, loading) =>
        set((state) => ({
          refreshStates: {
            ...state.refreshStates,
            [key]: loading,
          },
        })),

      // 更新刷新时间
      updateRefreshTime: (key) =>
        set((state) => ({
          lastRefreshTimes: {
            ...state.lastRefreshTimes,
            [key]: new Date().toLocaleTimeString('zh-CN', {
              hour: '2-digit',
              minute: '2-digit',
            }),
          },
        })),

      // 获取微信状态
      fetchWeChatStatus: async () => {
        try {
          const result = await api.getWeChatStatus();
          if (result.status === 'success' && result.data &&
              typeof result.data === 'object' && 'is_connected' in result.data) {
            set({ weChatStatus: result.data as WeChatStatus });
          }
          return result as ApiResponse<WeChatStatus>;
        } catch (error) {
          console.error('获取微信状态失败:', error);
          return { status: 'error', error: error instanceof Error ? error.message : '获取微信状态失败' } as ApiResponse<WeChatStatus>;
        }
      },

      // 同步微信数据
      syncWeChat: async () => {
        const { setRefreshState, updateRefreshTime, fetchWeChatStatus } = get();

        try {
          setRefreshState('wechat', true);
          const result = await api.syncWeChat();

          if (result.status === 'success') {
            updateRefreshTime('wechat');
            await fetchWeChatStatus(); // 重新获取状态
          }

          return result as ApiResponse<WeChatStatus>;
        } catch (error) {
          console.error('同步微信失败:', error);
          return { status: 'error', error: error instanceof Error ? error.message : '同步微信失败' } as ApiResponse<WeChatStatus>;
        } finally {
          setRefreshState('wechat', false);
        }
      },

      // 生成简报
      generateBriefing: async () => {
        const { setRefreshState, updateRefreshTime } = get();

        try {
          setRefreshState('briefing', true);
          const result = await api.generateBriefing();

          if (result.status === 'success' && result.data) {
            set({ briefing: result.data });
            updateRefreshTime('briefing');
          }

          return result;
        } catch (error) {
          console.error('生成简报失败:', error);
          throw error;
        } finally {
          setRefreshState('briefing', false);
        }
      },

      // 提取待办事项
      extractTodos: async () => {
        const { setRefreshState, updateRefreshTime } = get();

        try {
          setRefreshState('todos', true);
          const result = await api.extractTodos();

          if (result.status === 'success' && result.data) {
            set({ todos: result.data });
            updateRefreshTime('todos');
          }

          return result;
        } catch (error) {
          console.error('提取待办事项失败:', error);
          throw error;
        } finally {
          setRefreshState('todos', false);
        }
      },

      // 分析连接关系
      analyzeConnections: async () => {
        const { setRefreshState, updateRefreshTime } = get();

        try {
          setRefreshState('connections', true);
          const result = await api.analyzeConnections();

          if (result.status === 'success' && result.data) {
            set({ connections: result.data });
            updateRefreshTime('connections');
          }

          return result;
        } catch (error) {
          console.error('分析连接关系失败:', error);
          throw error;
        } finally {
          setRefreshState('connections', false);
        }
      },

      // 获取分析状态
      fetchAnalysisStatus: async () => {
        try {
          const result = await api.getAnalysisStatus();
          if (result.status === 'success' && result.data) {
            set({ analysisStatus: result.data });
          }
          return result;
        } catch (error) {
          console.error('获取分析状态失败:', error);
          throw error;
        }
      },

      // 更新设置
      updateSettings: async (newSettings: Partial<SystemSettings>) => {
        try {
          const result = await api.updateSettings(newSettings);
          if (result.status === 'success') {
            set((state) => ({
              settings: { ...state.settings, ...newSettings },
            }));
          }
          return result;
        } catch (error) {
          console.error('更新设置失败:', error);
          throw error;
        }
      },

      // 重置所有数据
      resetData: () =>
        set({
          briefing: null,
          todos: null,
          connections: null,
          analysisStatus: null,
          lastRefreshTimes: {
            briefing: null,
            todos: null,
            connections: null,
            wechat: null,
          },
        }),

      // 刷新所有数据
      refreshAll: async () => {
        const { weChatStatus } = get();
        const promises = [];

        // 如果微信已连接，同步数据
        if (weChatStatus?.is_connected) {
          promises.push(get().syncWeChat());
        }

        // 生成所有分析
        promises.push(get().generateBriefing());
        promises.push(get().extractTodos());
        promises.push(get().analyzeConnections());

        try {
          const results = await Promise.allSettled(promises);
          console.log('🔄 全部数据刷新完成', results);
        } catch (error) {
          console.error('刷新数据失败:', error);
          throw error;
        }
      },
    }),
    {
      name: 'twinos-store',
      partialize: (state: AppState) => ({
        settings: state.settings,
        activeView: state.activeView,
        sidebarCollapsed: state.sidebarCollapsed,
      }),
    }
  )
);

// 选择器函数
export const useWeChatStatus = () => useAppStore((state) => state.weChatStatus);
export const useBriefing = () => useAppStore((state) => state.briefing);
export const useTodos = () => useAppStore((state) => state.todos);
export const useConnections = () => useAppStore((state) => state.connections);
export const useSettings = () => useAppStore((state) => state.settings);
export const useIsLoading = () => useAppStore((state) => state.isLoading);
export const useError = () => useAppStore((state) => state.error);
export const useActiveView = () => useAppStore((state) => state.activeView);
export const useRefreshStates = () => useAppStore((state) => state.refreshStates);
export const useLastRefreshTimes = () => useAppStore((state) => state.lastRefreshTimes);

// 事件系统
export const useAppEvents = () => {
  const [events, setEvents] = React.useState<SystemEvent[]>([]);

  const addEvent = React.useCallback((event: Omit<SystemEvent, 'id' | 'timestamp' | 'read'>) => {
    const newEvent: SystemEvent = {
      ...event,
      id: Date.now().toString(),
      timestamp: new Date().toISOString(),
      read: false,
    };

    setEvents((prev) => [newEvent, ...prev].slice(0, 100)); // 保留最新100条
  }, []);

  const markEventRead = React.useCallback((id: string) => {
    setEvents((prev) =>
      prev.map((event) => (event.id === id ? { ...event, read: true } : event))
    );
  }, []);

  const clearEvents = React.useCallback(() => {
    setEvents([]);
  }, []);

  return {
    events,
    addEvent,
    markEventRead,
    clearEvents,
    unreadCount: events.filter((event) => !event.read).length,
  };
};

export default useAppStore;