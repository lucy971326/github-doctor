package steps

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github-doctor/workflow"
)

// ValidateStep 是验证 URL 的步骤
type ValidateStep struct{}

// Name 返回步骤名称
func (s *ValidateStep) Name() string {
	return "验证 URL"
}

// Execute 执行步骤
func (s *ValidateStep) Execute(ctx context.Context, data *workflow.AnalysisData) error {
	// 验证 URL 格式
	if data.RepoURL == "" {
		return fmt.Errorf("URL 不能为空")
	}

	// 解析 URL
	owner, repo, err := parseGitHubURL(data.RepoURL)
	if err != nil {
		return err
	}

	data.Owner = owner
	data.Repo = repo

	return nil
}

// parseGitHubURL 解析 GitHub URL
func parseGitHubURL(rawURL string) (owner, repo string, err error) {
	// 处理不带协议的 URL
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", "", fmt.Errorf("URL 解析失败: %w", err)
	}

	// 验证是 GitHub 域名
	host := parsed.Hostname()
	if host != "github.com" && host != "www.github.com" {
		return "", "", fmt.Errorf("不是 GitHub URL: %s", host)
	}

	// 解析路径
	path := strings.TrimPrefix(parsed.Path, "/")
	path = strings.TrimSuffix(path, ".git")
	parts := strings.Split(path, "/")

	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("URL 格式错误，应为: https://github.com/owner/repo")
	}

	return parts[0], parts[1], nil
}
