// 微信数据服务
export interface WeChatStatus {
  connected: boolean;
  messageCount: number;
  contactCount: number;
  chatCount: number;
  lastSyncTime: number;
  lastSyncStatus: string;
}

export interface ImportRequest {
  type: 'database' | 'backup' | 'file';
  path: string;
  mode?: 'full' | 'incremental';
}

export interface ImportResponse {
  sync_id: number;
  status: string;
}

export interface SyncRecord {
  id: number;
  sync_type: string;
  source: string;
  total_messages: number;
  new_messages: number;
  updated_messages: number;
  total_contacts: number;
  new_contacts: number;
  status: string;
  progress: number;
  error_message: string;
  start_time: string;
  end_time: string | null;
  duration: number;
  created_at: string;
  updated_at: string;
}

export interface WeChatMessage {
  id: number;
  svr_id: number;
  create_time: number;
  talker: string;
  type: number;
  sub_type: number;
  is_sender: number;
  seq: number;
  flag: number;
  status: number;
  content: string;
  display_content: string;
  chat_type: string;
  group_id: string;
  mentions: any;
  reply_to: string;
  created_at: string;
  updated_at: string;
}

export interface WeChatContact {
  id: number;
  username: string;
  nickname: string;
  remark: string;
  avatar: string;
  type: number;
  chat_type: string;
  member_count: number;
  owner: string;
  notice: string;
  message_count: number;
  last_active: number;
  is_blocked: boolean;
  privacy_level: number;
  created_at: string;
  updated_at: string;
}

export interface WeChatChat {
  id: number;
  chat_id: string;
  chat_type: string;
  name: string;
  avatar: string;
  message_count: number;
  unread_count: number;
  last_message: string;
  last_time: number;
  is_pinned: boolean;
  is_muted: boolean;
  is_archived: boolean;
  created_at: string;
  updated_at: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface MessageStatistics {
  total_messages: number;
  by_type: Record<string, number>;
  by_day: Record<string, number>;
}

class WeChatService {
  private baseURL = '/api/v1/wechat';

  // 连接微信数据源
  async connect(request: ImportRequest): Promise<ImportResponse> {
    const response = await fetch(`${this.baseURL}/connect`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || '连接失败');
    }

    const result = await response.json();
    return result.data;
  }

  // 同步微信数据
  async sync(request: ImportRequest): Promise<ImportResponse> {
    const response = await fetch(`${this.baseURL}/sync`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || '同步失败');
    }

    const result = await response.json();
    return result.data;
  }

  // 获取微信状态
  async getStatus(): Promise<WeChatStatus> {
    const response = await fetch(`${this.baseURL}/status`);

    if (!response.ok) {
      throw new Error('获取状态失败');
    }

    const result = await response.json();
    return result.data;
  }

  // 获取消息列表
  async getMessages(
    page: number = 1,
    pageSize: number = 20,
    chatType?: string,
    talker?: string
  ): Promise<PaginatedResponse<WeChatMessage>> {
    const params = new URLSearchParams({
      page: page.toString(),
      page_size: pageSize.toString(),
    });

    if (chatType) {
      params.append('chat_type', chatType);
    }
    if (talker) {
      params.append('talker', talker);
    }

    const response = await fetch(`${this.baseURL}/messages?${params}`);

    if (!response.ok) {
      throw new Error('获取消息列表失败');
    }

    const result = await response.json();
    return result.data;
  }

  // 获取联系人列表
  async getContacts(
    page: number = 1,
    pageSize: number = 20,
    chatType?: string
  ): Promise<PaginatedResponse<WeChatContact>> {
    const params = new URLSearchParams({
      page: page.toString(),
      page_size: pageSize.toString(),
    });

    if (chatType) {
      params.append('chat_type', chatType);
    }

    const response = await fetch(`${this.baseURL}/contacts?${params}`);

    if (!response.ok) {
      throw new Error('获取联系人列表失败');
    }

    const result = await response.json();
    return result.data;
  }

  // 获取聊天会话列表
  async getChats(
    page: number = 1,
    pageSize: number = 20,
    chatType?: string
  ): Promise<PaginatedResponse<WeChatChat>> {
    const params = new URLSearchParams({
      page: page.toString(),
      page_size: pageSize.toString(),
    });

    if (chatType) {
      params.append('chat_type', chatType);
    }

    const response = await fetch(`${this.baseURL}/chats?${params}`);

    if (!response.ok) {
      throw new Error('获取聊天会话列表失败');
    }

    const result = await response.json();
    return result.data;
  }

  // 获取同步记录
  async getSyncRecords(
    page: number = 1,
    pageSize: number = 20
  ): Promise<PaginatedResponse<SyncRecord>> {
    const params = new URLSearchParams({
      page: page.toString(),
      page_size: pageSize.toString(),
    });

    const response = await fetch(`${this.baseURL}/sync-records?${params}`);

    if (!response.ok) {
      throw new Error('获取同步记录失败');
    }

    const result = await response.json();
    return result.data;
  }

  // 搜索消息
  async searchMessages(
    keyword: string,
    page: number = 1,
    pageSize: number = 20
  ): Promise<PaginatedResponse<WeChatMessage>> {
    const params = new URLSearchParams({
      keyword,
      page: page.toString(),
      page_size: pageSize.toString(),
    });

    const response = await fetch(`${this.baseURL}/search?${params}`);

    if (!response.ok) {
      throw new Error('搜索消息失败');
    }

    const result = await response.json();
    return result.data;
  }

  // 获取聊天消息
  async getChatMessages(
    chatId: string,
    page: number = 1,
    pageSize: number = 20
  ): Promise<PaginatedResponse<WeChatMessage>> {
    const params = new URLSearchParams({
      page: page.toString(),
      page_size: pageSize.toString(),
    });

    const response = await fetch(`${this.baseURL}/chats/${chatId}/messages?${params}`);

    if (!response.ok) {
      throw new Error('获取聊天消息失败');
    }

    const result = await response.json();
    return result.data;
  }

  // 获取统计信息
  async getStatistics(): Promise<MessageStatistics> {
    const response = await fetch(`${this.baseURL}/statistics`);

    if (!response.ok) {
      throw new Error('获取统计信息失败');
    }

    const result = await response.json();
    return result.data;
  }
}

export const wechatService = new WeChatService();
export default wechatService;