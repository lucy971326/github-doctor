package report

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	g := NewGenerator()
	if g == nil {
		t.Fatal("NewGenerator() 返回 nil")
	}
	if g.templatePath == "" {
		t.Error("templatePath 为空")
	}
}

func TestGenerate_NoTemplate(t *testing.T) {
	g := &Generator{
		templatePath: "nonexistent.html",
	}

	data := &ReportData{
		RepoInfo: &RepoInfo{
			URL:   "https://github.com/vuejs/vue",
			Owner: "vuejs",
			Repo:  "vue",
		},
	}

	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "report-test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	outputPath := filepath.Join(tmpDir, "report.html")

	err = g.Generate(data, outputPath)
	if err == nil {
		t.Error("期望 Generate() 返回错误（模板不存在）")
	}
}

func TestReportData_Fields(t *testing.T) {
	data := &ReportData{
		RepoInfo: &RepoInfo{
			URL:   "https://github.com/vuejs/vue",
			Owner: "vuejs",
			Repo:  "vue",
		},
		GitHubStats: &GitHubStats{
			Stars: 12345,
		},
		ScoreResult: &ScoreResult{
			OverallScore: 85,
		},
		GeneratedAt: "2024-01-15 10:00:00",
	}

	if data.RepoInfo.Owner != "vuejs" {
		t.Errorf("Owner = %q", data.RepoInfo.Owner)
	}

	if data.GitHubStats.Stars != 12345 {
		t.Errorf("Stars = %d", data.GitHubStats.Stars)
	}

	if data.ScoreResult.OverallScore != 85 {
		t.Errorf("OverallScore = %d", data.ScoreResult.OverallScore)
	}
}
