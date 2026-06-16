package steps

import (
	"context"
	"fmt"

	"github-doctor/llm"
	"github-doctor/workflow"
)

// ScoreLLMStep 是 LLM 评分的步骤
type ScoreLLMStep struct {
	client *llm.Client
}

// NewScoreLLMStep 创建新的步骤
func NewScoreLLMStep() *ScoreLLMStep {
	return &ScoreLLMStep{
		client: llm.NewClient(),
	}
}

// NewScoreLLMStepWithKey 使用指定的 API Key 创建步骤
func NewScoreLLMStepWithKey(apiKey string) *ScoreLLMStep {
	return &ScoreLLMStep{
		client: llm.NewClientWithKey(apiKey),
	}
}

// Name 返回步骤名称
func (s *ScoreLLMStep) Name() string {
	return "LLM 评分"
}

// Execute 执行步骤
func (s *ScoreLLMStep) Execute(ctx context.Context, data *workflow.AnalysisData) error {
	// 构建输入
	input := &llm.AnalysisInput{
		RepoURL: data.RepoURL,
		Owner:   data.Owner,
		Repo:    data.Repo,
		GitHubData: &llm.GitHubInput{
			Stars:      data.GitHubData.Stars,
			Forks:      data.GitHubData.Forks,
			OpenIssues: data.GitHubData.OpenIssues,
			License:    data.GitHubData.License,
			Languages:  data.GitHubData.Languages,
			CreatedAt:  data.GitHubData.CreatedAt,
			UpdatedAt:  data.GitHubData.UpdatedAt,
		},
	}

	// 如果有 codegraph 数据，添加到输入
	if data.CodegraphData != nil && data.CodegraphData.Statistics != nil {
		input.CodegraphData = &llm.CodegraphInput{
			Statistics: &llm.StatisticsInput{
				Files:     data.CodegraphData.Statistics.Files,
				Functions: data.CodegraphData.Statistics.Functions,
				Classes:   data.CodegraphData.Statistics.Classes,
				Methods:   data.CodegraphData.Statistics.Methods,
				Edges:     data.CodegraphData.Statistics.Edges,
			},
			FileTree: data.CodegraphData.FileTree,
			CoreCode: data.CodegraphData.CoreCode,
		}
	}

	// 调用 LLM
	scoreData, err := s.client.Analyze(input)
	if err != nil {
		return fmt.Errorf("LLM 分析失败: %w", err)
	}

	data.ScoreData = scoreData
	return nil
}
