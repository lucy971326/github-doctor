package report

import (
	"encoding/json"
	"testing"
)

func TestReportData_JSON(t *testing.T) {
	data := &ReportData{
		RepoInfo: &RepoInfo{
			URL:   "https://github.com/vuejs/vue",
			Owner: "vuejs",
			Repo:  "vue",
		},
		GitHubStats: &GitHubStats{
			Stars:      12345,
			Forks:      678,
			OpenIssues: 42,
			License:    "MIT",
			Languages: map[string]int{
				"TypeScript": 2958283,
				"JavaScript": 137661,
			},
		},
		ScoreResult: &ScoreResult{
			OverallScore: 85,
			Dimensions: []Dimension{
				{Name: "代码质量", Score: 88, Comment: "代码结构清晰"},
			},
			Summary:         "这是一个成熟的开源项目",
			Recommendations: []string{"建议增加测试覆盖率"},
		},
		GeneratedAt: "2024-01-15 10:00:00",
	}

	// 序列化
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("json.Marshal() 错误: %v", err)
	}

	// 反序列化
	var decoded ReportData
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("json.Unmarshal() 错误: %v", err)
	}

	// 验证
	if decoded.RepoInfo.Owner != "vuejs" {
		t.Errorf("Owner = %q, want %q", decoded.RepoInfo.Owner, "vuejs")
	}

	if decoded.GitHubStats.Stars != 12345 {
		t.Errorf("Stars = %d, want 12345", decoded.GitHubStats.Stars)
	}

	if decoded.ScoreResult.OverallScore != 85 {
		t.Errorf("OverallScore = %d, want 85", decoded.ScoreResult.OverallScore)
	}
}

func TestDimension(t *testing.T) {
	d := Dimension{
		Name:    "代码质量",
		Score:   88,
		Comment: "代码结构清晰",
	}

	if d.Name != "代码质量" {
		t.Errorf("Name = %q, want %q", d.Name, "代码质量")
	}

	if d.Score != 88 {
		t.Errorf("Score = %d, want 88", d.Score)
	}

	if d.Comment != "代码结构清晰" {
		t.Errorf("Comment = %q, want %q", d.Comment, "代码结构清晰")
	}
}
