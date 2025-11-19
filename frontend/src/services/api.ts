import axios, { AxiosResponse } from 'axios';
import { ApiResponse } from '../types';

// API基础配置
const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:1234';

// 创建axios实例
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器
apiClient.interceptors.request.use(
  (config) => {
    console.log(`🚀 API Request: ${config.method?.toUpperCase()} ${config.url}`);
    return config;
  },
  (error) => {
    console.error('❌ Request Error:', error);
    return Promise.reject(error);
  }
);

// 响应拦截器
apiClient.interceptors.response.use(
  (response: AxiosResponse) => {
    console.log(`✅ API Response: ${response.status} ${response.config.url}`);
    return response;
  },
  (error) => {
    console.error('❌ Response Error:', error.response?.data || error.message);
    return Promise.reject(error);
  }
);

// 通用请求函数
async function request<T>(
  method: 'GET' | 'POST' | 'PUT' | 'DELETE',
  url: string,
  data?: any
): Promise<ApiResponse<T>> {
  try {
    const response = await apiClient.request<ApiResponse<T>>({
      method,
      url,
      data,
    });
    return response.data;
  } catch (error: any) {
    console.error('API Error:', error);

    // 返回标准化的错误响应
    return {
      status: 'error',
      error: error.response?.data?.message || error.message || 'Unknown error',
      message: error.response?.data?.message || error.message || 'Request failed',
    };
  }
}

// API方法封装
export const api = {
  // 健康检查
  async health() {
    return request('GET', '/health');
  },

  async healthDetailed() {
    return request('GET', '/health/detailed');
  },

  // 微信相关API
  async getWeChatStatus() {
    return request('GET', '/api/v1/wechat/status');
  },

  async connectWeChat() {
    return request('POST', '/api/v1/wechat/connect');
  },

  async syncWeChat() {
    return request('POST', '/api/v1/wechat/sync');
  },

  // 分析相关API
  async generateBriefing() {
    return request('POST', '/api/v1/analysis/briefing');
  },

  async extractTodos() {
    return request('POST', '/api/v1/analysis/todos');
  },

  async analyzeConnections() {
    return request('POST', '/api/v1/analysis/connections');
  },

  async analyzeMessage(messageId: string) {
    return request('POST', '/api/v1/analysis/message', { message_id: messageId });
  },

  async analyzeConversation(options?: {
    limit?: number;
    offset?: number;
  }) {
    return request('GET', '/api/v1/analysis/conversation', options);
  },

  async getEmotionTrend(days?: number) {
    return request('GET', '/api/v1/analysis/emotion/trend', { days });
  },

  async getTopicAnalysis() {
    return request('GET', '/api/v1/analysis/topics');
  },

  async getIntentDistribution() {
    return request('GET', '/api/v1/analysis/intents');
  },

  async getConversationSummary() {
    return request('GET', '/api/v1/analysis/summary');
  },

  async batchAnalyzeMessages(messageIds: string[]) {
    return request('POST', '/api/v1/analysis/batch', { message_ids: messageIds });
  },

  async getAnalysisStatus() {
    return request('GET', '/api/v1/analysis/status');
  },

  async configureAnalysis(config: {
    model?: string;
    temperature?: number;
    max_tokens?: number;
  }) {
    return request('PUT', '/api/v1/analysis/config', config);
  },

  // 数据相关API
  async getMessages(options?: {
    limit?: number;
    offset?: number;
    search?: string;
    sender?: string;
    date_from?: string;
    date_to?: string;
  }) {
    return request('GET', '/api/v1/data/messages', options);
  },

  async getContacts(options?: {
    limit?: number;
    offset?: number;
    search?: string;
  }) {
    return request('GET', '/api/v1/data/contacts', options);
  },

  async exportData(options: {
    format: 'json' | 'csv' | 'xlsx';
    date_range?: {
      start: string;
      end: string;
    };
    include_contacts?: boolean;
    include_messages?: boolean;
  }) {
    return request('POST', '/api/v1/data/export', options);
  },

  // 设置相关API
  async getSettings() {
    return request('GET', '/api/v1/settings');
  },

  async updateSettings(settings: any) {
    return request('PUT', '/api/v1/settings', settings);
  },

  // 备份相关API
  async createBackup(name?: string, description?: string) {
    return request('POST', '/api/v1/backup/create', { name, description });
  },

  async listBackups() {
    return request('GET', '/api/v1/backup/list');
  },

  async getBackupStats() {
    return request('GET', '/api/v1/backup/stats');
  },

  async scheduleBackup(options: {
    enabled: boolean;
    frequency: 'daily' | 'weekly' | 'monthly';
    time?: string;
    retention_days?: number;
  }) {
    return request('POST', '/api/v1/backup/schedule', options);
  },

  async getBackup(id: string) {
    return request('GET', `/api/v1/backup/${id}`);
  },

  async restoreBackup(id: string, options?: {
    overwrite?: boolean;
    selected_tables?: string[];
  }) {
    return request('POST', `/api/v1/backup/${id}/restore`, options);
  },

  async exportBackup(id: string, format?: 'json' | 'sql') {
    return request('POST', `/api/v1/backup/${id}/export`, { format });
  },

  async verifyBackup(id: string) {
    return request('GET', `/api/v1/backup/${id}/verify`);
  },

  async deleteBackup(id: string) {
    return request('DELETE', `/api/v1/backup/${id}`);
  },

  // 缓存相关API
  async getCacheStats() {
    return request('GET', '/api/v1/cache/stats');
  },

  async clearCache(pattern?: string) {
    return request('POST', '/api/v1/cache/clear', { pattern });
  },

  // 错误报告API
  async reportError(error: {
    type: string;
    message: string;
    stack?: string;
    component?: string;
    action?: string;
  }) {
    return request('POST', '/api/v1/error/report', error);
  },
};

// 导出默认实例
export default apiClient;