// API配置
const API_BASE_URL = 'http://localhost:1234/api/v1';

// API响应类型
export interface ApiResponse<T = any> {
  status: string;
  data?: T;
  error?: string;
}

// 微信相关API
export interface WeChatConnectRequest {
  db_path: string;
}

export interface WeChatStatus {
  is_connected: boolean;
  db_path: string;
  message_count?: number;
  contact_count: number;
  last_sync: string;
}

export interface SyncResult {
  start_time: string;
  end_time: string;
  duration: number;
  success_count: number;
  errors: number;
}

// 分析相关API
export interface BriefingResult {
  content: string;
  message_count: number;
  generated_at: string;
  analysis_id?: number;
}

export interface TodoItem {
  title: string;
  description: string;
  priority: string;
  deadline: string;
  related_people: string[];
  source_message_id: string;
  created_at: string;
}

export interface TodoResult {
  todos: TodoItem[];
  message_count: number;
  generated_at: string;
  analysis_id?: number;
}

export interface ConnectionAnalysis {
  person: string;
  action: string;
  context: string;
  importance: string;
  sentiment: string;
  message_count: number;
  last_mention: string;
}

export interface ConnectionResult {
  connections: ConnectionAnalysis[];
  message_count: number;
  generated_at: string;
  analysis_id?: number;
}

// API请求函数
class ApiClient {
  private baseURL: string;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<ApiResponse<T>> {
    const url = `${this.baseURL}${endpoint}`;

    const config: RequestInit = {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    };

    try {
      const response = await fetch(url, config);
      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || `HTTP error! status: ${response.status}`);
      }

      return data;
    } catch (error) {
      console.error('API request failed:', error);
      throw error;
    }
  }

  // 微信相关
  async connectWeChat(dbPath: string): Promise<ApiResponse> {
    return this.request('/wechat/connect', {
      method: 'POST',
      body: JSON.stringify({ db_path: dbPath }),
    });
  }

  async getWeChatStatus(): Promise<ApiResponse<WeChatStatus>> {
    return this.request('/wechat/status');
  }

  async syncWeChat(): Promise<ApiResponse<SyncResult>> {
    return this.request('/wechat/sync', {
      method: 'POST',
    });
  }

  // 分析相关
  async generateBriefing(): Promise<ApiResponse<BriefingResult>> {
    return this.request('/analysis/briefing', {
      method: 'POST',
    });
  }

  async extractTodos(): Promise<ApiResponse<TodoResult>> {
    return this.request('/analysis/todos', {
      method: 'POST',
    });
  }

  async analyzeConnections(): Promise<ApiResponse<ConnectionResult>> {
    return this.request('/analysis/connections', {
      method: 'POST',
    });
  }

  // 数据相关
  async getMessages(limit: number = 50, offset: number = 0): Promise<ApiResponse> {
    return this.request(`/data/messages?limit=${limit}&offset=${offset}`);
  }

  async getContacts(): Promise<ApiResponse> {
    return this.request('/data/contacts');
  }

  async exportData(): Promise<ApiResponse> {
    return this.request('/data/export', {
      method: 'POST',
    });
  }

  // 设置相关
  async getSettings(): Promise<ApiResponse<Record<string, string>>> {
    return this.request('/settings');
  }

  async updateSettings(settings: Record<string, string>): Promise<ApiResponse> {
    return this.request('/settings', {
      method: 'PUT',
      body: JSON.stringify(settings),
    });
  }

  // 健康检查
  async healthCheck(): Promise<ApiResponse> {
    return this.request('/health');
  }
}

export const apiClient = new ApiClient(API_BASE_URL);