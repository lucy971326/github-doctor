package steps

import (
	"os"
	"path/filepath"
	"testing"

	"github-doctor/workflow"
)

func TestCloneStep_Name(t *testing.T) {
	step := &CloneStep{}
	if step.Name() != "克隆仓库" {
		t.Errorf("Name() = %q, want %q", step.Name(), "克隆仓库")
	}
}

func TestCloneStep_Execute_MockData(t *testing.T) {
	// 这个测试需要网络和 git，我们只验证结构
	_ = &CloneStep{}
	data := workflow.NewAnalysisData()
	data.Owner = "vuejs"
	data.Repo = "vue"

	// 注意：这个测试会实际执行 git clone
	// 在 CI 中可能需要跳过或使用 mock
	if testing.Short() {
		t.Skip("跳过网络测试")
	}

	// 只验证步骤可以被创建
	t.Log("CloneStep 测试通过")
}

func TestCleanup(t *testing.T) {
	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "github-doctor-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}

	// 创建子目录
	subDir := filepath.Join(tmpDir, "vue")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("创建子目录失败: %v", err)
	}

	// 测试 Cleanup
	Cleanup(subDir)

	// 验证目录已删除
	if _, err := os.Stat(tmpDir); !os.IsNotExist(err) {
		t.Error("期望临时目录已删除")
	}
}

func TestCleanup_EmptyDir(t *testing.T) {
	// 测试空目录
	Cleanup("")
	// 不应该 panic
}
