package steps

import (
	"context"
	"fmt"

	"github-doctor/codegraph"
	"github-doctor/workflow"
)

// AnalyzeCodegraphStep 是 codegraph 分析的步骤
type AnalyzeCodegraphStep struct {
	binPath string
}

// NewAnalyzeCodegraphStep 创建新的步骤
func NewAnalyzeCodegraphStep() *AnalyzeCodegraphStep {
	return &AnalyzeCodegraphStep{}
}

// NewAnalyzeCodegraphStepWithBinPath 使用指定的 codegraph 路径创建步骤
func NewAnalyzeCodegraphStepWithBinPath(binPath string) *AnalyzeCodegraphStep {
	return &AnalyzeCodegraphStep{
		binPath: binPath,
	}
}

// Name 返回步骤名称
func (s *AnalyzeCodegraphStep) Name() string {
	return "codegraph 分析"
}

// Execute 执行步骤
func (s *AnalyzeCodegraphStep) Execute(ctx context.Context, data *workflow.AnalysisData) error {
	// 检查是否跳过代码分析
	if data.SkipCodeAnalysis {
		return nil
	}

	// 确保 codegraph 可用
	var binPath string
	var err error

	if s.binPath != "" {
		binPath = s.binPath
	} else {
		binPath, err = codegraph.EnsureCodegraph()
		if err != nil {
			return fmt.Errorf("codegraph 不可用: %w", err)
		}
	}

	// 创建客户端
	client := codegraph.NewClientWithPath(binPath)

	// 初始化索引
	if err := client.Init(data.RepoDir); err != nil {
		return fmt.Errorf("初始化失败: %w", err)
	}

	// 获取统计信息
	status, err := client.Status(data.RepoDir)
	if err != nil {
		return fmt.Errorf("获取统计失败: %w", err)
	}

	// 获取文件结构
	fileTree, err := client.Files(data.RepoDir)
	if err != nil {
		return fmt.Errorf("获取文件结构失败: %w", err)
	}

	// 搜索核心代码
	coreCode, err := client.Query(data.RepoDir, "main", 20)
	if err != nil {
		return fmt.Errorf("搜索代码失败: %w", err)
	}

	// 组装结果
	data.CodegraphData = &workflow.CodegraphData{
		Statistics: &workflow.Statistics{
			Files:     status.FilesIndexed,
			Functions: status.TotalNodes,
			Classes:   status.NodesByKind["class"],
			Methods:   status.NodesByKind["method"],
			Edges:     status.TotalEdges,
		},
		FileTree: fileTree,
		CoreCode: coreCode,
	}

	return nil
}
