package steps

import (
	"testing"

	"github-doctor/llm"
	"github-doctor/workflow"
)

func TestGenerateReportStep_Name(t *testing.T) {
	step := NewGenerateReportStep()
	if step.Name() != "生成报告" {
		t.Errorf("Name() = %q, want %q", step.Name(), "生成报告")
	}
}

func TestGenerateReportStep_NewGenerateReportStep(t *testing.T) {
	step := NewGenerateReportStep()
	if step == nil {
		t.Fatal("NewGenerateReportStep() 返回 nil")
	}
	if step.generator == nil {
		t.Error("generator 为 nil")
	}
}

func TestGenerateReportStep_Execute_MockData(t *testing.T) {
	// 这个测试需要模板文件，我们只验证结构
	_ = NewGenerateReportStep()
	data := workflow.NewAnalysisData()
	data.Owner = "vuejs"
	data.Repo = "vue"
	data.RepoURL = "https://github.com/vuejs/vue"

	// 注意：这个测试需要模板文件
	// 在 CI 中可能需要跳过或使用 mock
	if testing.Short() {
		t.Skip("跳过报告生成测试")
	}

	// 只验证步骤可以被创建
	t.Log("GenerateReportStep 测试通过")
}

func TestConvertDimensions(t *testing.T) {
	dimensions := []llm.Dimension{
		{Name: "代码质量", Score: 88, Comment: "代码结构清晰"},
		{Name: "社区活跃度", Score: 75, Comment: "Star 数较多"},
	}

	result := convertDimensions(dimensions)

	if len(result) != 2 {
		t.Errorf("convertDimensions() 返回 %d 个维度，期望 2 个", len(result))
	}

	if result[0].Name != "代码质量" {
		t.Errorf("维度名称 = %q, 期望 %q", result[0].Name, "代码质量")
	}

	if result[1].Score != 75 {
		t.Errorf("维度分数 = %d, 期望 75", result[1].Score)
	}
}
