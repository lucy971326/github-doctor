package workflow

import (
	"github-doctor/github"
	"github-doctor/llm"
)

// AnalysisData 是各步骤共享的数据上下文
type AnalysisData struct {
	// 输入
	RepoURL          string
	Owner            string
	Repo             string
	RepoDir          string // clone 的临时目录
	NeedCodeAnalysis bool   // 是否需要代码质量分析
	SkipCodeAnalysis bool   // 是否跳过代码分析（大仓库时）

	// GitHub API 数据
	GitHubData *github.GitHubData

	// codegraph 数据
	CodegraphData *CodegraphData

	// LLM 评分
	ScoreData *llm.ScoreData
}

// CodegraphData 存储 codegraph 分析结果
type CodegraphData struct {
	Statistics *Statistics
	FileTree   string
	CoreCode   string
}

// Statistics 存储代码统计信息
type Statistics struct {
	Files     int
	Functions int
	Classes   int
	Methods   int
	Edges     int
}

// NewAnalysisData 创建新的分析数据上下文
func NewAnalysisData() *AnalysisData {
	return &AnalysisData{}
}
