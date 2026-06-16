package github

import (
	"encoding/json"
	"testing"
)

func TestGitHubData_JSON(t *testing.T) {
	data := &GitHubData{
		Stars:      12345,
		Forks:      678,
		Watchers:   100,
		OpenIssues: 42,
		License:    "MIT",
		Description: "Vue.js is a progressive JavaScript framework",
		Languages: map[string]int{
			"TypeScript": 2958283,
			"JavaScript": 137661,
		},
		CreatedAt:    "2020-01-01T00:00:00Z",
		UpdatedAt:    "2024-01-15T00:00:00Z",
		DefaultBranch: "main",
		Size:         1024,
	}

	// 序列化
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("json.Marshal() 错误: %v", err)
	}

	// 反序列化
	var decoded GitHubData
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("json.Unmarshal() 错误: %v", err)
	}

	// 验证
	if decoded.Stars != 12345 {
		t.Errorf("Stars = %d, want 12345", decoded.Stars)
	}

	if decoded.Forks != 678 {
		t.Errorf("Forks = %d, want 678", decoded.Forks)
	}

	if decoded.License != "MIT" {
		t.Errorf("License = %q, want %q", decoded.License, "MIT")
	}

	if len(decoded.Languages) != 2 {
		t.Errorf("Languages 长度 = %d, want 2", len(decoded.Languages))
	}
}

func TestRepositoryInfo(t *testing.T) {
	info := &RepositoryInfo{
		Owner: "vuejs",
		Repo:  "vue",
		URL:   "https://github.com/vuejs/vue",
	}

	if info.Owner != "vuejs" {
		t.Errorf("Owner = %q, want %q", info.Owner, "vuejs")
	}

	if info.Repo != "vue" {
		t.Errorf("Repo = %q, want %q", info.Repo, "vue")
	}

	if info.URL != "https://github.com/vuejs/vue" {
		t.Errorf("URL = %q, want %q", info.URL, "https://github.com/vuejs/vue")
	}
}
