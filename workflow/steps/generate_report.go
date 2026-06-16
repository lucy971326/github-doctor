package steps

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github-doctor/llm"
	"github-doctor/report"
	"github-doctor/workflow"
)

// GenerateReportStep 是生成 HTML 报告的步骤
type GenerateReportStep struct {
	generator *report.Generator
}

// NewGenerateReportStep 创建新的步骤
func NewGenerateReportStep() *GenerateReportStep {
	return &GenerateReportStep{
		generator: report.NewGenerator(),
	}
}

// Name 返回步骤名称
func (s *GenerateReportStep) Name() string {
	return "生成报告"
}

// Execute 执行步骤
func (s *GenerateReportStep) Execute(ctx context.Context, data *workflow.AnalysisData) error {
	// 构建报告数据
	reportData := &report.ReportData{
		RepoInfo: &report.RepoInfo{
			URL:   data.RepoURL,
			Owner: data.Owner,
			Repo:  data.Repo,
		},
		GitHubStats: &report.GitHubStats{
			Stars:      data.GitHubData.Stars,
			Forks:      data.GitHubData.Forks,
			Watchers:   data.GitHubData.Watchers,
			OpenIssues: data.GitHubData.OpenIssues,
			License:    data.GitHubData.License,
			Languages:  data.GitHubData.Languages,
			CreatedAt:  data.GitHubData.CreatedAt,
			UpdatedAt:  data.GitHubData.UpdatedAt,
		},
	}

	// 如果有 AI 评分数据
	if data.ScoreData != nil {
		reportData.ScoreResult = &report.ScoreResult{
			OverallScore:    data.ScoreData.OverallScore,
			Dimensions:      convertDimensions(data.ScoreData.Dimensions),
			Summary:         data.ScoreData.Summary,
			Recommendations: data.ScoreData.Recommendations,
		}
	}

	// 生成报告文件
	outputPath := filepath.Join(os.TempDir(), fmt.Sprintf("github-doctor-%s-%s.html", data.Owner, data.Repo))
	if err := s.generator.Generate(reportData, outputPath); err != nil {
		return fmt.Errorf("生成报告失败: %w", err)
	}

	// 打开浏览器
	report.OpenInBrowser(outputPath)

	return nil
}

// convertDimensions 转换维度格式
func convertDimensions(dimensions []llm.Dimension) []report.Dimension {
	result := make([]report.Dimension, len(dimensions))
	for i, d := range dimensions {
		result[i] = report.Dimension{
			Name:    d.Name,
			Score:   d.Score,
			Comment: d.Comment,
		}
	}
	return result
}
