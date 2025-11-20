package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"twin-os/backend/pkg/logger"
)

// DeepSeekClient DeepSeek API客户端
type DeepSeekClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewDeepSeekClient 创建DeepSeek客户端
func NewDeepSeekClient(apiKey, baseURL string) *DeepSeekClient {
	if baseURL == "" {
		baseURL = "https://api.deepseek.com/v1"
	}

	return &DeepSeekClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
}

// ChatMessage 聊天消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	ID      string        `json:"id"`
	Object  string        `json:"object"`
	Created int64         `json:"created"`
	Model   string        `json:"model"`
	Choices []ChatChoice  `json:"choices"`
	Usage   ChatUsage     `json:"usage"`
}

// ChatChoice 聊天选择
type ChatChoice struct {
	Index        int           `json:"index"`
	Message      ChatMessage   `json:"message"`
	FinishReason string        `json:"finish_reason"`
}

// ChatUsage 使用统计
type ChatUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// GenerateBriefing 生成今日情报简报
func (c *DeepSeekClient) GenerateBriefing(messages []string) (string, error) {
	if c.apiKey == "" {
		return "", fmt.Errorf("DeepSeek API key not configured")
	}

	// 构建简报生成Prompt
	prompt := c.buildBriefingPrompt(messages)

	request := ChatRequest{
		Model: "deepseek-chat",
		Messages: []ChatMessage{
			{
				Role:    "system",
				Content: "你是TwinOS智能分析引擎的核心情报分析助手。你的使命是从海量聊天信息中提炼出高价值的战略情报。",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.4,
		MaxTokens:   1500,
	}

	response, err := c.chat(request)
	if err != nil {
		return "", fmt.Errorf("failed to generate briefing: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	return response.Choices[0].Message.Content, nil
}

// ExtractTodos 提取待办事项
func (c *DeepSeekClient) ExtractTodos(messages []string) ([]TodoItem, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("DeepSeek API key not configured")
	}

	prompt := c.buildTodoPrompt(messages)

	request := ChatRequest{
		Model: "deepseek-chat",
		Messages: []ChatMessage{
			{
				Role:    "system",
				Content: "你是TwinOS智能任务管理系统的核心AI助手。你擅长从非结构化的聊天记录中识别、提取并结构化待办事项。你的分析准确、实用且具有前瞻性。",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.2,
		MaxTokens:   2500,
	}

	response, err := c.chat(request)
	if err != nil {
		return nil, fmt.Errorf("failed to extract todos: %w", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no response generated")
	}

	// 提取并解析JSON响应
	jsonContent := c.extractJSON(response.Choices[0].Message.Content)
	
	// 定义响应结构（包含todos和summary）
	type TodoResponse struct {
		Todos []TodoItem `json:"todos"`
	}
	var todoResp TodoResponse
	
	err = json.Unmarshal([]byte(jsonContent), &todoResp)
	if err != nil {
		// 尝试直接解析为数组（兼容性）
		var todos []TodoItem
		err2 := json.Unmarshal([]byte(jsonContent), &todos)
		if err2 == nil {
			return todos, nil
		}
		
		// 如果JSON解析失败，尝试手动解析
		logger.Warn("⚠️ JSON unmarshal failed, trying manual parsing: " + err.Error())
		return c.parseTodosManually(response.Choices[0].Message.Content)
	}

	return todoResp.Todos, nil
}

// AnalyzeConnections 分析人际关系
func (c *DeepSeekClient) AnalyzeConnections(messages []string, contacts []string) ([]ConnectionAnalysis, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("DeepSeek API key not configured")
	}

	prompt := c.buildConnectionPrompt(messages, contacts)

	request := ChatRequest{
		Model: "deepseek-chat",
		Messages: []ChatMessage{
			{
				Role:    "system",
				Content: "你是TwinOS智能人脉分析系统的专业AI分析师。你精通社交网络分析、关系挖掘和商业智能，能够从复杂的聊天交互中识别有价值的人脉动态和潜在机会。",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.3,
		MaxTokens:   2000,
	}

	response, err := c.chat(request)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze connections: %w", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no response generated")
	}

	// 提取并解析JSON响应
	jsonContent := c.extractJSON(response.Choices[0].Message.Content)

	type ConnectionResponse struct {
		Connections []ConnectionAnalysis `json:"connections"`
	}
	var connResp ConnectionResponse

	err = json.Unmarshal([]byte(jsonContent), &connResp)
	if err != nil {
		// 兼容直接返回数组的情况
		var conns []ConnectionAnalysis
		if err2 := json.Unmarshal([]byte(jsonContent), &conns); err2 == nil {
			return conns, nil
		}
		
		logger.Warn("⚠️ Failed to parse connection analysis JSON: " + err.Error())
		return []ConnectionAnalysis{}, nil
	}

	return connResp.Connections, nil
}

// chat 执行聊天请求
func (c *DeepSeekClient) chat(request ChatRequest) (*ChatResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var chatResponse ChatResponse
	err = json.Unmarshal(body, &chatResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &chatResponse, nil
}

// buildBriefingPrompt 构建简报生成Prompt
func (c *DeepSeekClient) buildBriefingPrompt(messages []string) string {
	// 预处理消息：去重
	processedMessages := c.preprocessMessages(messages, 100)
	// 截断消息：防止超过 Context Window (假设 4096 或 8192，这里保守设为 3000 字符)
	processedMessages = c.truncateMessages(processedMessages, 3000)

	currentTime := time.Now().Format("2006-01-02 15:04")
	weekday := time.Now().Weekday().String()

	prompt := fmt.Sprintf(`你是TwinOS智能分析引擎的首席情报官。今天是 %s (%s)。
你的任务是阅读用户的加密聊天记录，剔除噪音，提炼出对用户决策有价值的战略情报。

## 核心原则
1. **结论先行**：不要复述聊天记录，直接给出分析结论。
2. **识别风险**：敏锐地发现潜在的冲突、延期、误解风险。
3. **行动导向**：所有的洞察都必须转化为可执行的建议。
4. **隐私意识**：消息中的人名已被脱敏为 [姓名X]、[公司Y]，请直接使用代号引用，不要尝试猜测真实姓名。

## 分析维度
- **🔴 紧急 (Urgent)**: 需要立即处理的危机、即将到期的Deadline。
- **🟡 重要 (Important)**: 关键项目的进展、重要合作机会、人脉维系。
- **🔵 关注 (Watch)**: 值得注意的市场动态、潜在趋势。

## 输出格式 (Markdown)

# 每日情报简报

### 🚨 紧急警报 (P0)
> *无紧急事项则显示"今日无紧急风险"*
- [截止时间] **事项标题**: 核心内容与风险点。

### 📋 决策待办 (Action Items)
- [ ] **待办1**: 对应上下文及建议行动。
- [ ] **待办2**: ...

### 💡 深度洞察 (Insights)
- **项目/业务**: 分析关键项目的健康度。
- **人脉/情绪**: [姓名X] 似乎对某事有顾虑... / [公司Y] 表现出合作意向...

### 📊 数据概览
- 消息量: %d 条 (仅供参考)

## 聊天记录片段 (最新的在最后):
`, currentTime, weekday, len(processedMessages))

	for i, msg := range processedMessages {
		prompt += fmt.Sprintf("[%d] %s\n", i+1, msg)
	}

	return prompt
}

// buildTodoPrompt 构建待办提取Prompt
func (c *DeepSeekClient) buildTodoPrompt(messages []string) string {
	// 预处理消息
	processedMessages := c.preprocessMessages(messages, 100)

	currentDate := time.Now().Format("2006-01-02")
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	prompt := fmt.Sprintf(`你是TwinOS智能任务管理助手，专门从聊天记录中识别和提取结构化的待办事项。今天是%s。

## 任务识别标准

### 🎯 任务类型
1. **行动类任务**：需要用户执行的具体动作
2. **决策类任务**：需要用户做出判断或选择
3. **响应类任务**：需要回复或确认的事项
4. **准备类任务**：为会议或活动做准备
5. **跟进类任务**：需要后续跟踪的事项

### 🔥 优先级判定
- **HIGH**：包含"紧急"、"立即"、"今天"、"明天"等时间敏感词
- **MEDIUM**：包含"下周"、"尽快"、"需要"等中等时间词
- **LOW**：包含"有时间"、"不急"、"后续"等低时间敏感词

### ⏰ 时间表达式识别
- **绝对时间**："2024-01-15"、"1月15号"、"下周一"
- **相对时间**："明天"、"后天"、"下个月"、"3天后"
- **模糊时间**："近期"、"不久"、"尽快"
- **截止时间**："截止"、"之前"、"务必在"

### 👥 人员识别
- **直接提及**：@用户名、真实姓名
- **间接提及**：职位、角色（"经理"、"客户"、"技术"）
- **协作需求**："配合"、"协调"、"沟通"

## 聊天消息数据：
`, currentDate)

	for i, msg := range processedMessages {
		// 标记消息中的时间敏感词
		markedMsg := c.highlightTimeKeywords(msg)
		prompt += fmt.Sprintf("[%d] %s\n", i+1, markedMsg)
	}

	prompt += fmt.Sprintf(`

## 输出要求

请严格按照以下JSON格式返回，每个字段都必须填写：

` + "```json" + `
{
  "todos": [
    {
      "title": "简洁的任务标题（不超过20字）",
      "description": "详细描述（包含上下文，不超过100字）",
      "priority": "high|medium|low",
      "deadline": "YYYY-MM-DD格式或具体时间描述",
      "related_people": ["人员1", "人员2"],
      "action_type": "execute|decide|respond|prepare|follow",
      "estimated_duration": "预估耗时（如：30分钟、2小时）",
      "source_message_index": 消息序号,
      "confidence": 0.8,
      "tags": ["标签1", "标签2"]
    }
  ],
  "summary": {
    "total_tasks": 数量,
    "urgent_count": 紧急任务数,
    "upcoming_deadlines": ["最近截止时间1", "最近截止时间2"],
    "key_people": ["最频繁提及的人员"]
  }
}
` + "```" + `

## 注意事项：
1. 今天是%s，明天是%s
2. 如果没有明确截止日期，根据紧急程度推断合理时间
3. 每条消息最多提取1个主要任务
4. 优先提取需要用户主动行动的任务
5. 忽略纯信息性、已完成的任务
`, currentDate, tomorrow)

	return prompt
}

// buildConnectionPrompt 构建人脉分析Prompt
func (c *DeepSeekClient) buildConnectionPrompt(messages []string, contacts []string) string {
	// 预处理消息
	processedMessages := c.preprocessMessages(messages, 50)

	prompt := fmt.Sprintf(`你是TwinOS智能人脉分析助手，专门从聊天记录中识别人际关系动态和互动模式。

## 分析维度

### 🔍 关系类型识别
1. **上下级关系**：工作汇报、任务分配、决策指示
2. **协作关系**：项目合作、信息共享、资源协调
3. **客户关系**：商务洽谈、需求沟通、服务提供
4. **社交关系**：私人交流、情感互动、兴趣分享
5. **竞争关系**：资源争夺、立场对立、意见分歧

### 📈 互动强度分析
- **高频互动**：频繁提及、持续关注
- **深度互动**：复杂话题、长时间讨论
- **正向互动**：积极评价、支持鼓励
- **负向互动**：批评建议、质疑反对
- **中性互动**：信息传递、客观陈述

### 💡 机会与风险识别
- **合作机会**：潜在项目、资源互补
- **关系风险**：误解冲突、关系紧张
- **维护需求**：需要主动联系的关系
- **拓展机会**：新的人脉连接可能

## 已知联系人列表：
%s

## 聊天消息数据：
`, fmt.Sprintf("%v", contacts))

	for i, msg := range processedMessages {
		// 标记消息中的人名关键词
		markedMsg := c.highlightPersonKeywords(msg, contacts)
		prompt += fmt.Sprintf("[%d] %s\n", i+1, markedMsg)
	}

	prompt += `

## 输出要求

请严格按照以下JSON格式返回，提供深度的人脉分析：

` + "```json" + `
{
  "connections": [
    {
      "person": "人员姓名",
      "relationship_type": "superior|subordinate|peer|client|supplier|friend|competitor|other",
      "interaction_pattern": "frequent|regular|occasional|rare",
      "action": "提及|讨论|评价|求助|指导|合作|邀请|批评",
      "context": "具体的互动内容和背景",
      "importance": "high|medium|low",
      "sentiment": "positive|neutral|negative",
      "message_count": 交互消息数量,
      "last_interaction": "最后互动时间描述",
      "opportunity_score": 0.8,
      "risk_level": "low|medium|high",
      "relationship_strength": "strong|moderate|weak",
      "suggested_action": "建议的跟进行动",
      "topics": ["话题1", "话题2"]
    }
  ],
  "insights": {
    "total_connections": 关系总数,
    "key_relationships": ["最重要的关系人"],
    "emerging_opportunities": ["新出现的合作机会"],
    "attention_required": ["需要关注的关系"],
    "network_health": "网络健康度评分",
    "growth_areas": ["可以拓展的关系领域"]
  }
}
` + "```" + `

## 分析重点：
1. 识别对人名、代词的引用关系
2. 分析情绪倾向和互动质量
3. 评估关系的重要性和发展潜力
4. 发现潜在的商业或合作机会
5. 识别需要维护或改善的关系

请提供深度、准确的分析结果。
`

	return prompt
}

// extractJSON 从响应中提取JSON字符串
func (c *DeepSeekClient) extractJSON(content string) string {
	// 移除Markdown代码块标记
	content = strings.TrimSpace(content)
	if strings.HasPrefix(content, "```json") {
		content = content[7:]
	} else if strings.HasPrefix(content, "```") {
		content = content[3:]
	}
	if strings.HasSuffix(content, "```") {
		content = content[:len(content)-3]
	}
	return strings.TrimSpace(content)
}

// parseTodosManually 手动解析待办事项
func (c *DeepSeekClient) parseTodosManually(content string) ([]TodoItem, error) {
	// 简单的手动解析逻辑
	// 这里可以添加更复杂的解析规则
	return []TodoItem{}, nil
}

// preprocessMessages 预处理消息
func (c *DeepSeekClient) preprocessMessages(messages []string, maxCount int) []string {
	if len(messages) == 0 {
		return []string{}
	}

	// 去重和清理
	seen := make(map[string]bool)
	var processed []string

	for i := len(messages) - 1; i >= 0 && len(processed) < maxCount; i-- {
		msg := strings.TrimSpace(messages[i])
		if msg != "" && !seen[msg] && len(msg) > 5 {
			processed = append([]string{msg}, processed...)
			seen[msg] = true
		}
	}

	return processed
}

// truncateMessages 截断消息以适应Token限制
func (c *DeepSeekClient) truncateMessages(messages []string, maxTokens int) []string {
	var truncated []string
	currentTokens := 0
	// 预估Token：平均1个汉字=2个Token，1个英文单词=1.3个Token。
	// 保守估计，按字符数直接作为Token数的参考（中文偏多）。
	// System Prompt 和 Template 大约占 500 Tokens。
	limit := maxTokens - 500

	// 优先保留最新的消息
	for i := 0; i < len(messages); i++ {
		msg := messages[i]
		estimatedTokens := len(msg) // 粗略估算
		if currentTokens+estimatedTokens > limit {
			break
		}
		truncated = append(truncated, msg)
		currentTokens += estimatedTokens
	}

	return truncated
}

// highlightTimeKeywords 高亮时间关键词
func (c *DeepSeekClient) highlightTimeKeywords(message string) string {
	timeKeywords := []string{
		"紧急", "立即", "马上", "今天", "明天", "后天", "本周", "下周",
		"截止", "之前", "务必", "尽快", "尽快", "不急", "有时间",
		"2024-", "2025-", "月", "号", "点", "分",
	}

	highlighted := message
	for _, keyword := range timeKeywords {
		highlighted = strings.ReplaceAll(highlighted, keyword, "**"+keyword+"**")
	}

	return highlighted
}

// highlightPersonKeywords 高亮人名关键词
func (c *DeepSeekClient) highlightPersonKeywords(message string, contacts []string) string {
	if len(contacts) == 0 {
		return message
	}

	highlighted := message
	for _, contact := range contacts {
		if contact != "" {
			highlighted = strings.ReplaceAll(highlighted, contact, "**"+contact+"**")
		}
	}

	// 也标记常见的关系词
	relationKeywords := []string{
		"张总", "李总", "王经理", "刘经理", "老板", "领导",
		"同事", "朋友", "客户", "供应商", "合作伙伴",
	}

	for _, keyword := range relationKeywords {
		highlighted = strings.ReplaceAll(highlighted, keyword, "**"+keyword+"**")
	}

	return highlighted
}

// TodoItem 待办事项
type TodoItem struct {
	Title             string   `json:"title"`
	Description       string   `json:"description"`
	Priority          string   `json:"priority"`          // high, medium, low
	Deadline          string   `json:"deadline"`
	RelatedPeople     []string `json:"related_people"`
	SourceMessageID   string   `json:"source_message_id"`
	CreatedAt         string   `json:"created_at"`
	ActionType        string   `json:"action_type"`        // execute, decide, respond, prepare, follow
	EstimatedDuration string   `json:"estimated_duration"` // 预估耗时
	Confidence        float64  `json:"confidence"`         // 置信度 0-1
	Tags              []string `json:"tags"`               // 标签
}

// ConnectionAnalysis 人脉分析结果
type ConnectionAnalysis struct {
	Person            string   `json:"person"`
	RelationshipType  string   `json:"relationship_type"`  // superior, subordinate, peer, client, supplier, friend, competitor, other
	InteractionPattern string  `json:"interaction_pattern"` // frequent, regular, occasional, rare
	Action            string   `json:"action"`            // 提及, 讨论, 评价, 求助, 指导, 合作, 邀请, 批评
	Context           string   `json:"context"`           // 上下文
	Importance        string   `json:"importance"`        // high, medium, low
	Sentiment         string   `json:"sentiment"`         // positive, neutral, negative
	MessageCount      int      `json:"message_count"`
	LastInteraction   string   `json:"last_interaction"`  // 最后互动时间描述
	OpportunityScore  float64  `json:"opportunity_score"`  // 机会评分 0-1
	RiskLevel         string   `json:"risk_level"`         // low, medium, high
	RelationshipStrength string `json:"relationship_strength"` // strong, moderate, weak
	SuggestedAction   string   `json:"suggested_action"`   // 建议的跟进行动
	Topics            []string `json:"topics"`            // 话题
}