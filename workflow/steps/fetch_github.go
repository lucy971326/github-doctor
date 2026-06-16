package steps

import (
	"context"
	"fmt"

	"github-doctor/github"
	"github-doctor/workflow"
)

// FetchGitHubStep 是获取 GitHub 数据的步骤
type FetchGitHubStep struct {
	client *github.Client
}

// NewFetchGitHubStep 创建新的步骤
func NewFetchGitHubStep() *FetchGitHubStep {
	return &FetchGitHubStep{
		client: github.NewClient(),
	}
}

// NewFetchGitHubStepWithClient 使用指定的客户端创建步骤
func NewFetchGitHubStepWithClient(client *github.Client) *FetchGitHubStep {
	return &FetchGitHubStep{
		client: client,
	}
}

// Name 返回步骤名称
func (s *FetchGitHubStep) Name() string {
	return "获取 GitHub 数据"
}

// Execute 执行步骤
func (s *FetchGitHubStep) Execute(ctx context.Context, data *workflow.AnalysisData) error {
	githubData, err := s.client.GetRepository(data.Owner, data.Repo)
	if err != nil {
		return fmt.Errorf("获取数据失败: %w", err)
	}

	data.GitHubData = githubData
	return nil
}
