package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Client 是 GitHub API 客户端
type Client struct {
	httpClient *http.Client
	token      string
}

// NewClient 创建新的 GitHub API 客户端
func NewClient() *Client {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		token = os.Getenv("GH_TOKEN") // GitHub CLI 的环境变量
	}

	return &Client{
		httpClient: &http.Client{},
		token:      token,
	}
}

// NewClientWithToken 使用指定的 Token 创建客户端
func NewClientWithToken(token string) *Client {
	return &Client{
		httpClient: &http.Client{},
		token:      token,
	}
}

// GetRepository 获取仓库信息
func (c *Client) GetRepository(owner, repo string) (*GitHubData, error) {
	// 获取仓库基本信息
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 添加认证头（如果有 token）
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求仓库信息失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API 返回错误: %d", resp.StatusCode)
	}

	var data GitHubData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("解析仓库信息失败: %w", err)
	}

	// 获取语言分布
	languages, err := c.GetLanguages(owner, repo)
	if err != nil {
		// 语言获取失败不影响整体
		data.Languages = make(map[string]int)
	} else {
		data.Languages = languages
	}

	return &data, nil
}

// GetLanguages 获取语言分布
func (c *Client) GetLanguages(owner, repo string) (map[string]int, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/languages", owner, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 添加认证头（如果有 token）
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求语言信息失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API 返回错误: %d", resp.StatusCode)
	}

	var languages map[string]int
	if err := json.NewDecoder(resp.Body).Decode(&languages); err != nil {
		return nil, fmt.Errorf("解析语言信息失败: %w", err)
	}

	return languages, nil
}
