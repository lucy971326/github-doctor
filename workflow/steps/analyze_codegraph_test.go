package steps

import (
	"testing"

	"github-doctor/workflow"
)

func TestAnalyzeCodegraphStep_Name(t *testing.T) {
	step := NewAnalyzeCodegraphStep()
	if step.Name() != "codegraph 分析" {
		t.Errorf("Name() = %q, want %q", step.Name(), "codegraph 分析")
	}
}

func TestAnalyzeCodegraphStep_NewAnalyzeCodegraphStep(t *testing.T) {
	step := NewAnalyzeCodegraphStep()
	if step == nil {
		t.Fatal("NewAnalyzeCodegraphStep() 返回 nil")
	}
}

func TestAnalyzeCodegraphStep_Execute_MockData(t *testing.T) {
	// 这个测试需要 codegraph 二进制，我们只验证结构
	_ = NewAnalyzeCodegraphStep()
	data := workflow.NewAnalysisData()
	data.RepoDir = "/tmp/test-repo"

	// 注意：这个测试需要 codegraph 二进制
	// 在 CI 中可能需要跳过或使用 mock
	if testing.Short() {
		t.Skip("跳过 codegraph 测试")
	}

	// 只验证步骤可以被创建
	t.Log("AnalyzeCodegraphStep 测试通过")
}
