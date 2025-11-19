import { create } from 'zustand';
import { devtools } from 'zustand/middleware';
import { apiClient, WeChatStatus, BriefingResult, TodoResult, ConnectionResult } from '../services/api';

// 应用状态类型
interface AppState {
  // 初始化状态
  isInitialized: boolean;

  // 微信连接状态
  weChatStatus: WeChatStatus | null;
  isConnecting: boolean;

  // 数据同步状态
  isSyncing: boolean;
  lastSyncTime: string | null;

  // 分析结果
  briefing: BriefingResult | null;
  todos: TodoResult | null;
  connections: ConnectionResult | null;

  // 分析状态
  isAnalyzing: boolean;

  // 错误状态
  error: string | null;

  // 操作方法
  initializeApp: () => void;
  connectWeChat: (dbPath: string) => Promise<void>;
  checkWeChatStatus: () => Promise<void>;
  syncWeChat: () => Promise<void>;
  generateBriefing: () => Promise<void>;
  extractTodos: () => Promise<void>;
  analyzeConnections: () => Promise<void>;
  clearError: () => void;
}

export const useAppStore = create<AppState>()(
  devtools(
    (set, get) => ({
      // 初始状态
      isInitialized: false,
      weChatStatus: null,
      isConnecting: false,
      isSyncing: false,
      lastSyncTime: null,
      briefing: null,
      todos: null,
      connections: null,
      isAnalyzing: false,
      error: null,

      // 初始化应用
      initializeApp: async () => {
        try {
          // 检查微信状态
          await get().checkWeChatStatus();

          set({
            isInitialized: true,
            error: null
          });
        } catch (error) {
          console.error('Failed to initialize app:', error);
          set({
            error: error instanceof Error ? error.message : '初始化失败',
            isInitialized: true
          });
        }
      },

      // 连接微信数据库
      connectWeChat: async (dbPath: string) => {
        try {
          set({
            isConnecting: true,
            error: null
          });

          await apiClient.connectWeChat(dbPath);
          await get().checkWeChatStatus();

          set({
            isConnecting: false
          });
        } catch (error) {
          set({
            isConnecting: false,
            error: error instanceof Error ? error.message : '连接微信失败'
          });
          throw error;
        }
      },

      // 检查微信状态
      checkWeChatStatus: async () => {
        try {
          const response = await apiClient.getWeChatStatus();
          if (response.data) {
            set({ weChatStatus: response.data });
          }
        } catch (error) {
          console.error('Failed to check WeChat status:', error);
          set({
            weChatStatus: {
              is_connected: false,
              db_path: '',
              contact_count: 0,
              last_sync: new Date().toISOString(),
            }
          });
        }
      },

      // 同步微信数据
      syncWeChat: async () => {
        try {
          set({
            isSyncing: true,
            error: null
          });

          await apiClient.syncWeChat();

          set({
            isSyncing: false,
            lastSyncTime: new Date().toISOString()
          });
        } catch (error) {
          set({
            isSyncing: false,
            error: error instanceof Error ? error.message : '同步失败'
          });
          throw error;
        }
      },

      // 生成情报简报
      generateBriefing: async () => {
        try {
          set({
            isAnalyzing: true,
            error: null
          });

          const response = await apiClient.generateBriefing();

          set({
            briefing: response.data || null,
            isAnalyzing: false
          });
        } catch (error) {
          set({
            isAnalyzing: false,
            error: error instanceof Error ? error.message : '生成简报失败'
          });
          throw error;
        }
      },

      // 提取待办事项
      extractTodos: async () => {
        try {
          set({
            isAnalyzing: true,
            error: null
          });

          const response = await apiClient.extractTodos();

          set({
            todos: response.data || null,
            isAnalyzing: false
          });
        } catch (error) {
          set({
            isAnalyzing: false,
            error: error instanceof Error ? error.message : '提取待办失败'
          });
          throw error;
        }
      },

      // 分析人脉关系
      analyzeConnections: async () => {
        try {
          set({
            isAnalyzing: true,
            error: null
          });

          const response = await apiClient.analyzeConnections();

          set({
            connections: response.data || null,
            isAnalyzing: false
          });
        } catch (error) {
          set({
            isAnalyzing: false,
            error: error instanceof Error ? error.message : '分析人脉失败'
          });
          throw error;
        }
      },

      // 清除错误
      clearError: () => {
        set({ error: null });
      },
    }),
    {
      name: 'twin-os-store',
    }
  )
);