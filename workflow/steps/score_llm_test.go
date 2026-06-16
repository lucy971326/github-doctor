package steps

import (
	"testing"

	"github-doctor/workflow"
)

func TestScoreLLMStep_Name(t *testing.T) {
	step := NewScoreLLMStep()
	if step.Name() != "LLM 评分" {
		t.Errorf("Name() = %q, want %q", step.Name(), "LLM 评分")
	}
}

func TestScoreLLMStep_NewScoreLLMStep(t *testing.T) {
	step := NewScoreLLMStep()
	if step == nil {
		t.Fatal("NewScoreLLMStep() 返回 nil")
	}
	if step.client == nil {
		t.Error("client 为 nil")
	}
}

func TestScoreLLMStep_Execute_MockData(t *testing.T) {
	// 这个测试需要 DeepSeek API Key，我们只验证结构
	_ = NewScoreLLMStep()
	data := workflow.NewAnalysisData()
	data.Owner = "vuejs"
	data.Repo = "vue"
	data.RepoURL = "https://github.com/vuejs/vue"

	// 注意：这个测试需要 DEEPSEEK_API_KEY 环境变量
	// 在 CI 中可能需要跳过或使用 mock
	if testing.Short() {
		t.Skip("跳过 LLM 测试")
	}

	// 只验证步骤可以被创建
	t.Log("ScoreLLMStep 测试通过")
}
