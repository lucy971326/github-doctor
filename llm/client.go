package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	// DeepSeek API 地址
	deepseekAPIURL = "https://api.deepseek.com/chat/completions"

	// 系统提示词
	systemPrompt = `你是一个专业的 GitHub 仓库代码质量分析专家。你的任务是分析给定的仓库数据，并给出专业的评分和评价。

## 重要规则（必须遵守）：

1. **绝对不要提及任何内部工具或实现细节**，如 "codegraph"、"索引"、"分析工具" 等。你只是一个分析专家，直接基于数据说话。

2. **不要猜测或假设文件是否存在**。如果数据中没有显示某个文件，不要说"缺少"该文件，而是说"数据中未体现"或直接忽略该项。

3. **只基于提供的数据进行评价**，不要添加数据中没有的信息。

4. **评价要正面、建设性**，避免负面猜测。

## 输出格式（严格 JSON）：

{
  "overall_score": 0-100的整数,
  "dimensions": [
    {
      "name": "维度名称",
      "score": 0-100的整数,
      "comment": "该维度的详细评价（基于数据，不猜测）"
    }
  ],
  "summary": "整体评价总结",
  "recommendations": ["建议1", "建议2"]
}

## 评分维度：
1. 代码质量 - 基于代码结构、函数数量、类数量等
2. 社区活跃度 - 基于 Star、Fork、Issues 等
3. 文档完善度 - 基于 README、License 等（只评价数据中明确有的）
4. 维护状态 - 基于最近更新时间、项目大小等
5. 项目规模 - 基于文件数量、代码行数等

## 示例（正确）：
"项目使用 TypeScript，代码结构清晰，函数和类的组织合理。"

## 示例（错误 - 不要这样做）：
"注意到分析工具数据可能不完整..."（暴露了实现细节）
"缺少 LICENSE 文件"（数据中没有就不要说缺少）`
)

// Client 是 DeepSeek API 客户端
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// NewClient 创建新的 DeepSeek API 客户端
func NewClient() *Client {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("DEEPSEEK_API_KEY")
	}

	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

// NewClientWithKey 使用指定的 API Key 创建客户端
func NewClientWithKey(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

// Analyze 分析仓库数据并返回评分
func (c *Client) Analyze(input *AnalysisInput) (*ScoreData, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("DEEPSEEK_API_KEY 环境变量未设置")
	}

	// 构建用户消息
	userMessage, err := c.buildUserMessage(input)
	if err != nil {
		return nil, fmt.Errorf("构建消息失败: %w", err)
	}

	// 构建请求
	req := &ChatRequest{
		Model: "deepseek-v4-pro",
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userMessage},
		},
		ResponseFormat: &ResponseFormat{
			Type: "json_object",
		},
	}

	// 发送请求
	resp, err := c.sendRequest(req)
	if err != nil {
		return nil, fmt.Errorf("API 请求失败: %w", err)
	}

	// 解析响应
	var scoreData ScoreData
	if err := json.Unmarshal([]byte(resp), &scoreData); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 后处理：过滤掉泄露的实现细节
	filterImplementationDetails(&scoreData)

	return &scoreData, nil
}

// filterImplementationDetails 过滤掉泄露的实现细节
func filterImplementationDetails(data *ScoreData) {
	// 需要过滤的关键词
	bannedWords := []string{"codegraph", "Codegraph", "code graph", "Code Graph", "索引工具", "分析工具"}

	// 过滤 summary
	data.Summary = filterText(data.Summary, bannedWords)

	// 过滤 dimensions
	for i := range data.Dimensions {
		data.Dimensions[i].Comment = filterText(data.Dimensions[i].Comment, bannedWords)
	}

	// 过滤 recommendations
	for i := range data.Recommendations {
		data.Recommendations[i] = filterText(data.Recommendations[i], bannedWords)
	}
}

// filterText 过滤文本中的敏感词
func filterText(text string, bannedWords []string) string {
	for _, word := range bannedWords {
		// 简单替换，实际可以更复杂
		text = strings.ReplaceAll(text, word, "***")
	}
	return text
}

// buildUserMessage 构建用户消息
func (c *Client) buildUserMessage(input *AnalysisInput) (string, error) {
	// 基础信息
	message := fmt.Sprintf(`请分析以下 GitHub 仓库数据：

仓库信息：
- URL: %s
- Owner: %s
- Repo: %s

GitHub 数据：
- Stars: %d
- Forks: %d
- Open Issues: %d
- License: %s
- 创建时间: %s
- 最后更新: %s
- 语言分布: %s`,
		input.RepoURL,
		input.Owner,
		input.Repo,
		input.GitHubData.Stars,
		input.GitHubData.Forks,
		input.GitHubData.OpenIssues,
		input.GitHubData.License,
		input.GitHubData.CreatedAt,
		input.GitHubData.UpdatedAt,
		formatLanguages(input.GitHubData.Languages),
	)

	// 如果有 codegraph 数据，添加到消息
	if input.CodegraphData != nil && input.CodegraphData.Statistics != nil {
		message += fmt.Sprintf(`

代码分析数据：
- 文件数量: %d
- 函数数量: %d
- 类数量: %d
- 方法数量: %d
- 调用关系数: %d

文件结构：
%s

核心代码：
%s`,
			input.CodegraphData.Statistics.Files,
			input.CodegraphData.Statistics.Functions,
			input.CodegraphData.Statistics.Classes,
			input.CodegraphData.Statistics.Methods,
			input.CodegraphData.Statistics.Edges,
			input.CodegraphData.FileTree,
			input.CodegraphData.CoreCode,
		)
	} else {
		message += `

（注意：本次分析未包含代码深度分析数据，仅基于 GitHub API 数据进行评价）`
	}

	message += `

请根据以上数据给出专业的评分和评价。`

	return message, nil
}

// sendRequest 发送 API 请求
func (c *Client) sendRequest(req *ChatRequest) (string, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequest("POST", deepseekAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API 返回错误: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("API 返回空选择")
	}

	return chatResp.Choices[0].Message.Content, nil
}

// formatLanguages 格式化语言分布
func formatLanguages(languages map[string]int) string {
	if len(languages) == 0 {
		return "无"
	}

	result := ""
	for lang, bytes := range languages {
		if result != "" {
			result += ", "
		}
		result += fmt.Sprintf("%s: %d bytes", lang, bytes)
	}
	return result
}
