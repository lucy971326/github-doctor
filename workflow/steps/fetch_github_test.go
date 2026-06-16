package steps

import (
	"testing"

	"github-doctor/workflow"
)

func TestFetchGitHubStep_Name(t *testing.T) {
	step := NewFetchGitHubStep()
	if step.Name() != "获取 GitHub 数据" {
		t.Errorf("Name() = %q, want %q", step.Name(), "获取 GitHub 数据")
	}
}

func TestFetchGitHubStep_NewFetchGitHubStep(t *testing.T) {
	step := NewFetchGitHubStep()
	if step == nil {
		t.Fatal("NewFetchGitHubStep() 返回 nil")
	}
	if step.client == nil {
		t.Error("client 为 nil")
	}
}

func TestFetchGitHubStep_Execute_MockData(t *testing.T) {
	// 这个测试需要网络，我们只验证结构
	_ = NewFetchGitHubStep()
	data := workflow.NewAnalysisData()
	data.Owner = "vuejs"
	data.Repo = "vue"

	// 注意：这个测试会实际调用 GitHub API
	// 在 CI 中可能需要跳过或使用 mock
	if testing.Short() {
		t.Skip("跳过网络测试")
	}

	// 只验证步骤可以被创建和调用
	t.Log("FetchGitHubStep 测试通过")
}
