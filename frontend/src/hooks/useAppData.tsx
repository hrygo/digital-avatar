import { useState, useEffect, useCallback } from 'react';

// 模拟数据
const mockBriefing = `# 每日情报简报

## 📊 核心数据概览
- **消息总数**: 1,247 条
- **新增对话**: 23 条
- **处理状态**: 已完成

## 🎯 关键洞察

### 1. 商业机会识别
通过分析最近的聊天记录，发现 **3个** 潜在商业机会：

- **技术合作需求**：与老王的讨论中提到需要前端开发支持
- **产品咨询项目**：李经理询问系统定制化可能性
- **培训业务扩展**：团队需要新技术栈培训

### 2. 决策待办事项
- **优先级高**: 客户演示系统准备
- **优先级中**: 技术方案文档更新
- **优先级低**: 市场调研报告

### 3. 人脉动态分析
- **新连接**: 5个新联系人
- **互动频率**: 与核心团队保持高频沟通
- **关系强度**: 商业伙伴关系稳定

## 🔍 趋势预测
基于历史数据分析，预计下周：
- 商务谈判进入关键阶段
- 技术债务需要重点关注
- 团队规模可能扩大

---
*报告生成时间: ${new Date().toLocaleString('zh-CN')}*`;

const mockTodos = [
  {
    id: '1',
    content: '准备客户演示 PPT',
    priority: 'high' as const,
    status: 'pending' as const,
    due_date: '2025-11-20',
    assignee: '张三',
    tags: ['商务', '演示']
  },
  {
    id: '2',
    content: '完成技术方案文档',
    priority: 'medium' as const,
    status: 'pending' as const,
    due_date: '2025-11-22',
    assignee: '李四',
    tags: ['技术', '文档']
  },
  {
    id: '3',
    content: '安排团队会议',
    priority: 'low' as const,
    status: 'completed' as const,
    due_date: '2025-11-19',
    assignee: '王五',
    tags: ['团队', '会议']
  }
];

const mockConnections = [
  {
    id: '1',
    person_name: '李总',
    relationship: 'client' as const,
    interaction_count: 15,
    last_interaction: '2025-11-18',
    strength_score: 0.9,
    topics: ['商业合作', '技术方案', '项目交付']
  },
  {
    id: '2',
    person_name: '王经理',
    relationship: 'partner' as const,
    interaction_count: 8,
    last_interaction: '2025-11-17',
    strength_score: 0.7,
    topics: ['技术咨询', '市场拓展']
  },
  {
    id: '3',
    person_name: '张老师',
    relationship: 'mentor' as const,
    interaction_count: 12,
    last_interaction: '2025-11-19',
    strength_score: 0.8,
    topics: ['职业发展', '行业趋势']
  }
];

export interface SystemStatus {
  isOnline: boolean;
  lastUpdate: string;
  backendStatus: 'connected' | 'disconnected' | 'loading';
  dataCount: {
    messages: number;
    contacts: number;
    briefings: number;
  };
}

export const useAppData = () => {
  const [systemStatus, setSystemStatus] = useState<SystemStatus>({
    isOnline: true,
    lastUpdate: new Date().toISOString(),
    backendStatus: 'loading',
    dataCount: {
      messages: 0,
      contacts: 0,
      briefings: 0
    }
  });

  const [briefing, setBriefing] = useState(mockBriefing);
  const [todos, setTodos] = useState(mockTodos);
  const [connections, setConnections] = useState(mockConnections);
  const [loading, setLoading] = useState({
    briefing: false,
    todos: false,
    connections: false
  });

  // 检查后端连接状态
  const checkBackendStatus = useCallback(async () => {
    try {
      const response = await fetch('http://localhost:1234/health');
      if (response.ok) {
        const data = await response.json();
        setSystemStatus(prev => ({
          ...prev,
          backendStatus: 'connected',
          lastUpdate: new Date().toISOString(),
          dataCount: {
            messages: data.message_count || 0,
            contacts: data.contact_count || 0,
            briefings: 1 // 模拟数据
          }
        }));
      }
    } catch (error) {
      setSystemStatus(prev => ({
        ...prev,
        backendStatus: 'disconnected',
        lastUpdate: new Date().toISOString()
      }));
    }
  }, []);

  // 生成简报
  const generateBriefing = useCallback(async () => {
    setLoading(prev => ({ ...prev, briefing: true }));
    await new Promise(resolve => setTimeout(resolve, 2000)); // 模拟 API 调用
    setLoading(prev => ({ ...prev, briefing: false }));
    // 在真实应用中，这里会调用实际的 API
  }, []);

  // 提取待办事项
  const extractTodos = useCallback(async () => {
    setLoading(prev => ({ ...prev, todos: true }));
    await new Promise(resolve => setTimeout(resolve, 1500)); // 模拟 API 调用
    setLoading(prev => ({ ...prev, todos: false }));
    // 在真实应用中，这里会调用实际的 API
  }, []);

  // 分析连接关系
  const analyzeConnections = useCallback(async () => {
    setLoading(prev => ({ ...prev, connections: true }));
    await new Promise(resolve => setTimeout(resolve, 1800)); // 模拟 API 调用
    setLoading(prev => ({ ...prev, connections: false }));
    // 在真实应用中，这里会调用实际的 API
  }, []);

  // 初始化时检查后端状态
  useEffect(() => {
    checkBackendStatus();
    const interval = setInterval(checkBackendStatus, 30000); // 每30秒检查一次
    return () => clearInterval(interval);
  }, [checkBackendStatus]);

  return {
    systemStatus,
    briefing,
    todos,
    connections,
    loading,
    generateBriefing,
    extractTodos,
    analyzeConnections,
    refreshStatus: checkBackendStatus
  };
};