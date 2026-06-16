package workflow

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// Workflow 是工作流引擎
type Workflow struct {
	Steps []Step
	Data  *AnalysisData
}

// NewWorkflow 创建新的工作流
func NewWorkflow() *Workflow {
	return &Workflow{
		Steps: []Step{},
		Data:  NewAnalysisData(),
	}
}

// AddStep 添加步骤到工作流
func (w *Workflow) AddStep(step Step) {
	w.Steps = append(w.Steps, step)
}

// Run 执行工作流
func (w *Workflow) Run(repoURL string) error {
	// 验证并解析 URL
	owner, repo, err := parseGitHubURL(repoURL)
	if err != nil {
		return fmt.Errorf("无效的 GitHub URL: %w", err)
	}

	w.Data.RepoURL = repoURL
	w.Data.Owner = owner
	w.Data.Repo = repo

	// 执行所有步骤
	for _, step := range w.Steps {
		// 检查是否应该跳过
		if w.shouldSkipStep(step) {
			continue
		}

		fmt.Printf("  [%s] %s\n", "...", step.Name())
		if err := step.Execute(context.Background(), w.Data); err != nil {
			return fmt.Errorf("%s: %w", step.Name(), err)
		}
		fmt.Printf("  [%s] %s\n", "✓", step.Name())
	}

	return nil
}

// shouldSkipStep 判断是否应该跳过该步骤
func (w *Workflow) shouldSkipStep(step Step) bool {
	// 如果跳过代码分析，跳过 codegraph 和 LLM 步骤
	if w.Data.SkipCodeAnalysis {
		name := step.Name()
		if name == "codegraph 分析" || name == "LLM 评分" {
			return true
		}
	}
	return false
}

// parseGitHubURL 解析 GitHub URL，返回 owner 和 repo
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
