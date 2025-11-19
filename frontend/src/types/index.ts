// API响应类型
export interface ApiResponse<T = any> {
  status: 'success' | 'error';
  data?: T;
  message?: string;
  error?: string;
}

// 微信状态类型
export interface WeChatStatus {
  is_connected: boolean;
  message_count: number;
  contact_count: number;
  last_sync?: string;
  error?: string;
}

// 简报类型
export interface Briefing {
  id: number;
  content: string;
  message_count: number;
  generated_at: string;
  analysis_id: number;
}

// 待办事项类型
export interface Todo {
  id: string;
  content: string;
  priority: 'high' | 'medium' | 'low';
  status: 'pending' | 'completed' | 'cancelled';
  created_at: string;
  due_date?: string;
  assignee?: string;
  tags?: string[];
}

export interface TodoList {
  todos: Todo[];
  total_count: number;
  completed_count: number;
  pending_count: number;
}

// 连接关系类型
export interface Connection {
  id: string;
  person_name: string;
  relationship: 'colleague' | 'friend' | 'client' | 'partner' | 'other';
  strength: number; // 1-10
  last_interaction: string;
  interaction_count: number;
  topics: string[];
  sentiment: 'positive' | 'neutral' | 'negative';
}

export interface ConnectionList {
  connections: Connection[];
  total_count: number;
  strong_connections: number; // strength >= 7
  recent_interactions: number; // 最近7天内的互动
}

// 分析状态类型
export interface AnalysisStatus {
  is_processing: boolean;
  last_analysis: string;
  queue_size: number;
  processing_stage?: string;
  estimated_completion?: string;
}

// 系统设置类型
export interface SystemSettings {
  auto_sync: boolean;
  sync_interval: number; // 分钟
  notification_enabled: boolean;
  notification_types: {
    briefing: boolean;
    urgent_todos: boolean;
    connections: boolean;
  };
  privacy: {
    data_retention_days: number;
    anonymization_enabled: boolean;
    local_processing_only: boolean;
  };
  ui: {
    theme: 'dark' | 'light';
    compact_mode: boolean;
    show_animations: boolean;
  };
}

// 应用状态类型
export interface AppState {
  // 系统状态
  isLoading: boolean;
  isInitialized: boolean;
  error: string | null;

  // 数据状态
  weChatStatus: WeChatStatus | null;
  briefing: Briefing | null;
  todos: TodoList | null;
  connections: ConnectionList | null;
  analysisStatus: AnalysisStatus | null;
  settings: SystemSettings;

  // UI状态
  activeView: 'dashboard' | 'briefing' | 'todos' | 'connections' | 'settings';
  sidebarCollapsed: boolean;

  // 刷新状态
  refreshStates: {
    briefing: boolean;
    todos: boolean;
    connections: boolean;
    wechat: boolean;
  };

  // 时间戳
  lastRefreshTimes: {
    briefing: string | null;
    todos: string | null;
    connections: string | null;
    wechat: string | null;
  };
}

// 应用状态和动作类型
export interface AppStore extends AppState {
  // 初始化方法
  initializeApp: () => Promise<void>;

  // UI状态方法
  setActiveView: (view: AppState['activeView']) => void;
  toggleSidebar: () => void;
  setSidebarCollapsed: (collapsed: boolean) => void;

  // 加载状态方法
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  clearError: () => void;

  // 刷新状态方法
  setRefreshState: (key: keyof AppState['refreshStates'], loading: boolean) => void;
  updateRefreshTime: (key: keyof AppState['lastRefreshTimes']) => void;

  // 数据操作方法
  syncWeChat: () => Promise<ApiResponse<WeChatStatus>>;
  generateBriefing: () => Promise<ApiResponse<Briefing>>;
  extractTodos: () => Promise<ApiResponse<TodoList>>;
  analyzeConnections: () => Promise<ApiResponse<ConnectionList>>;
  fetchWeChatStatus: () => Promise<ApiResponse<WeChatStatus>>;
  fetchAnalysisStatus: () => Promise<ApiResponse<AnalysisStatus>>;

  // 批量操作
  refreshAll: () => Promise<void>;
}

// 事件类型
export interface SystemEvent {
  id: string;
  type: 'info' | 'warning' | 'error' | 'success';
  title: string;
  message: string;
  timestamp: string;
  read: boolean;
  actions?: Array<{
    label: string;
    action: string;
  }>;
}

// 导航菜单类型
export interface NavigationItem {
  id: string;
  label: string;
  icon: string;
  path: string;
  badge?: number;
  active?: boolean;
}

// 卡片组件Props类型
export interface CardProps {
  children: React.ReactNode;
  className?: string;
  title?: string;
  subtitle?: string;
  actions?: React.ReactNode;
  loading?: boolean;
  error?: string;
}

// 统计卡片Props类型
export interface StatCardProps {
  title: string;
  value: string | number;
  change?: {
    value: number;
    type: 'increase' | 'decrease' | 'neutral';
  };
  icon?: React.ReactNode;
  className?: string;
}

// 列表项Props类型
export interface ListItemProps {
  title: string;
  subtitle?: string;
  badge?: string | number;
  icon?: React.ReactNode;
  actions?: React.ReactNode;
  onClick?: () => void;
  className?: string;
}

// 模态框Props类型
export interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  title?: string;
  children: React.ReactNode;
  size?: 'sm' | 'md' | 'lg' | 'xl';
  className?: string;
}

// 通知Props类型
export interface NotificationProps {
  type: 'success' | 'error' | 'warning' | 'info';
  title: string;
  message?: string;
  duration?: number;
  onClose?: () => void;
  actions?: Array<{
    label: string;
    onClick: () => void;
  }>;
}

// API错误类型
export interface ApiError {
  code: string;
  message: string;
  details?: any;
  timestamp: string;
}

// 请求配置类型
export interface RequestConfig {
  timeout?: number;
  retries?: number;
  headers?: Record<string, string>;
}

// 分页类型
export interface Pagination {
  page: number;
  limit: number;
  total: number;
  totalPages: number;
  hasNext: boolean;
  hasPrev: boolean;
}

// 搜索类型
export interface SearchOptions {
  query: string;
  filters?: Record<string, any>;
  sort?: {
    field: string;
    order: 'asc' | 'desc';
  };
  pagination?: {
    page: number;
    limit: number;
  };
}

// 导出类型
export interface ExportOptions {
  format: 'json' | 'csv' | 'pdf';
  dateRange?: {
    start: string;
    end: string;
  };
  filters?: Record<string, any>;
}