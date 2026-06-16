package workflow

import (
	"testing"
)

func TestNewAnalysisData(t *testing.T) {
	data := NewAnalysisData()

	if data == nil {
		t.Fatal("NewAnalysisData() 返回 nil")
	}

	if data.RepoURL != "" {
		t.Errorf("期望 RepoURL 为空，得到 %q", data.RepoURL)
	}

	if data.Owner != "" {
		t.Errorf("期望 Owner 为空，得到 %q", data.Owner)
	}

	if data.Repo != "" {
		t.Errorf("期望 Repo 为空，得到 %q", data.Repo)
	}

	if data.GitHubData != nil {
		t.Error("期望 GitHubData 为 nil")
	}

	if data.CodegraphData != nil {
		t.Error("期望 CodegraphData 为 nil")
	}

	if data.ScoreData != nil {
		t.Error("期望 ScoreData 为 nil")
	}
}

func TestAnalysisData_Fields(t *testing.T) {
	data := NewAnalysisData()

	// 测试设置值
	data.RepoURL = "https://github.com/vuejs/vue"
	data.Owner = "vuejs"
	data.Repo = "vue"
	data.RepoDir = "/tmp/vue"

	if data.RepoURL != "https://github.com/vuejs/vue" {
		t.Errorf("RepoURL 设置失败")
	}

	if data.Owner != "vuejs" {
		t.Errorf("Owner 设置失败")
	}

	if data.Repo != "vue" {
		t.Errorf("Repo 设置失败")
	}

	if data.RepoDir != "/tmp/vue" {
		t.Errorf("RepoDir 设置失败")
	}
}
