package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
				Content: "你是TwinOS的情报分析助手，专门负责从聊天记录中提取并生成简洁的情报简报。",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.3,
		MaxTokens:   1000,
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
				Content: "你是TwinOS的任务分析助手，专门负责从聊天记录中识别和提取待办事项。请以JSON格式返回结果。",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.1,
		MaxTokens:   2000,
	}

	response, err := c.chat(request)
	if err != nil {
		return nil, fmt.Errorf("failed to extract todos: %w", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no response generated")
	}

	// 解析JSON响应
	var todos []TodoItem
	err = json.Unmarshal([]byte(response.Choices[0].Message.Content), &todos)
	if err != nil {
		// 如果JSON解析失败，尝试手动解析
		return c.parseTodosManually(response.Choices[0].Message.Content)
	}

	return todos, nil
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
				Content: "你是TwinOS的人脉分析助手，专门负责分析聊天记录中的人际关系和提及情况。请以JSON格式返回分析结果。",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.2,
		MaxTokens:   1500,
	}

	response, err := c.chat(request)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze connections: %w", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no response generated")
	}

	// 解析JSON响应
	var connections []ConnectionAnalysis
	err = json.Unmarshal([]byte(response.Choices[0].Message.Content), &connections)
	if err != nil {
		// 如果JSON解析失败，返回空结果
		logger.Warn("⚠️ Failed to parse connection analysis JSON: " + err.Error())
		return []ConnectionAnalysis{}, nil
	}

	return connections, nil
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
	prompt := `请分析以下聊天消息，生成一份300字左右的今日情报简报。要求：

1. 提取最重要的信息和事件
2. 识别关键决策点和待办事项
3. 总结重要的人脉动态
4. 使用简洁、结构化的语言
5. 按重要性排序

聊天消息：
`

	for i, msg := range messages {
		if i >= 50 { // 限制消息数量避免Token超限
			break
		}
		prompt += fmt.Sprintf("- %s\n", msg)
	}

	prompt += `
请生成简报，格式如下：
# 今日情报简报

## 重要事件
[事件1]
[事件2]

## 关键决策
[决策1]
[决策2]

## 人脉动态
[动态1]
[动态2]
`

	return prompt
}

// buildTodoPrompt 构建待办提取Prompt
func (c *DeepSeekClient) buildTodoPrompt(messages []string) string {
	prompt := `请分析以下聊天消息，提取所有需要处理的待办事项。要求：

1. 识别明确需要用户处理的事项
2. 提取时间节点和截止日期
3. 确定优先级（高/中/低）
4. 识别相关人员和协作需求
5. 以JSON格式返回

JSON格式示例：
[
  {
    "title": "完成项目报告",
    "description": "需要提交Q4季度报告给管理层",
    "priority": "high",
    "deadline": "2024-01-15",
    "related_people": ["张总", "李经理"],
    "source_message_index": 5
  }
]

聊天消息：
`

	for i, msg := range messages {
		if i >= 100 { // 限制消息数量
			break
		}
		prompt += fmt.Sprintf("[%d] %s\n", i, msg)
	}

	return prompt
}

// buildConnectionPrompt 构建人脉分析Prompt
func (c *DeepSeekClient) buildConnectionPrompt(messages []string, contacts []string) string {
	prompt := `请分析以下聊天消息，识别人际关系动态。要求：

1. 识别谁在谈论你或提到与你相关的事情
2. 分析重要合作机会或潜在冲突
3. 识别需要维护的重要关系
4. 以JSON格式返回

JSON格式示例：
[
  {
    "person": "张总",
    "action": "提及",
    "context": "讨论项目进展",
    "importance": "high",
    "sentiment": "positive",
    "message_count": 3
  }
]

联系人列表：` + fmt.Sprintf("%v", contacts) + `

聊天消息：
`

	for i, msg := range messages {
		if i >= 50 { // 限制消息数量
			break
		}
		prompt += fmt.Sprintf("- %s\n", msg)
	}

	return prompt
}

// parseTodosManually 手动解析待办事项
func (c *DeepSeekClient) parseTodosManually(content string) ([]TodoItem, error) {
	// 简单的手动解析逻辑
	// 这里可以添加更复杂的解析规则
	return []TodoItem{}, nil
}

// TodoItem 待办事项
type TodoItem struct {
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Priority        string   `json:"priority"` // high, medium, low
	Deadline        string   `json:"deadline"`
	RelatedPeople   []string `json:"related_people"`
	SourceMessageID string   `json:"source_message_id"`
	CreatedAt       string   `json:"created_at"`
}

// ConnectionAnalysis 人脉分析结果
type ConnectionAnalysis struct {
	Person       string  `json:"person"`
	Action       string  `json:"action"`       // 提及, 讨论, 评价
	Context      string  `json:"context"`      // 上下文
	Importance   string  `json:"importance"`   // high, medium, low
	Sentiment    string  `json:"sentiment"`    // positive, neutral, negative
	MessageCount int     `json:"message_count"`
	LastMention  string  `json:"last_mention"` // 最后提及时间
}