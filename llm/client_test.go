package llm

import (
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("NewClient() 返回 nil")
	}
	if client.httpClient == nil {
		t.Error("httpClient 为 nil")
	}
}

func TestFormatLanguages(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]int
		expected string
	}{
		{
			name:     "空语言",
			input:    map[string]int{},
			expected: "无",
		},
		{
			name:     "单个语言",
			input:    map[string]int{"Go": 1000},
			expected: "Go: 1000 bytes",
		},
		{
			name: "多个语言",
			input: map[string]int{
				"TypeScript": 2958283,
				"JavaScript": 137661,
			},
			// 顺序不确定，只检查包含关系
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatLanguages(tt.input)

			if tt.name == "多个语言" {
				// 检查是否包含两种语言
				if !contains(result, "TypeScript") || !contains(result, "JavaScript") {
					t.Errorf("formatLanguages() = %q, 应包含 TypeScript 和 JavaScript", result)
				}
			} else {
				if result != tt.expected {
					t.Errorf("formatLanguages() = %q, want %q", result, tt.expected)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestAnalyze_NoAPIKey(t *testing.T) {
	// 保存原始环境变量
	origKey := os.Getenv("DEEPSEEK_API_KEY")
	defer os.Setenv("DEEPSEEK_API_KEY", origKey)

	// 清除环境变量
	os.Unsetenv("DEEPSEEK_API_KEY")

	client := NewClient()
	input := &AnalysisInput{
		RepoURL: "https://github.com/vuejs/vue",
		Owner:   "vuejs",
		Repo:    "vue",
	}

	_, err := client.Analyze(input)
	if err == nil {
		t.Error("期望 Analyze() 返回错误（无 API Key）")
	}
}

func TestAnalysisInput(t *testing.T) {
	input := &AnalysisInput{
		RepoURL: "https://github.com/vuejs/vue",
		Owner:   "vuejs",
		Repo:    "vue",
		GitHubData: &GitHubInput{
			Stars:      12345,
			Forks:      678,
			OpenIssues: 42,
			License:    "MIT",
		},
		CodegraphData: &CodegraphInput{
			Statistics: &StatisticsInput{
				Files:     245,
				Functions: 1058,
				Classes:   55,
				Methods:   838,
				Edges:     11639,
			},
			FileTree: "📁 文件树...",
			CoreCode: "🔑 核心代码...",
		},
	}

	if input.RepoURL != "https://github.com/vuejs/vue" {
		t.Errorf("RepoURL = %q", input.RepoURL)
	}

	if input.GitHubData.Stars != 12345 {
		t.Errorf("Stars = %d", input.GitHubData.Stars)
	}
}
