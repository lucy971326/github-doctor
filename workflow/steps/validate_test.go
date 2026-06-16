package steps

import (
	"context"
	"testing"

	"github-doctor/workflow"
)

func TestValidateStep_Name(t *testing.T) {
	step := &ValidateStep{}
	if step.Name() != "验证 URL" {
		t.Errorf("Name() = %q, want %q", step.Name(), "验证 URL")
	}
}

func TestValidateStep_Execute_Success(t *testing.T) {
	step := &ValidateStep{}
	data := workflow.NewAnalysisData()
	data.RepoURL = "https://github.com/vuejs/vue"

	err := step.Execute(context.Background(), data)
	if err != nil {
		t.Errorf("Execute() 返回错误: %v", err)
	}

	if data.Owner != "vuejs" {
		t.Errorf("Owner = %q, want %q", data.Owner, "vuejs")
	}

	if data.Repo != "vue" {
		t.Errorf("Repo = %q, want %q", data.Repo, "vue")
	}
}

func TestValidateStep_Execute_EmptyURL(t *testing.T) {
	step := &ValidateStep{}
	data := workflow.NewAnalysisData()
	data.RepoURL = ""

	err := step.Execute(context.Background(), data)
	if err == nil {
		t.Error("期望 Execute() 返回错误（空 URL）")
	}
}

func TestValidateStep_Execute_InvalidURL(t *testing.T) {
	step := &ValidateStep{}
	data := workflow.NewAnalysisData()
	data.RepoURL = "https://gitlab.com/vuejs/vue"

	err := step.Execute(context.Background(), data)
	if err == nil {
		t.Error("期望 Execute() 返回错误（非 GitHub URL）")
	}
}

func TestParseGitHubURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantOwner string
		wantRepo  string
		wantErr   bool
	}{
		{
			name:      "标准 URL",
			url:       "https://github.com/vuejs/vue",
			wantOwner: "vuejs",
			wantRepo:  "vue",
			wantErr:   false,
		},
		{
			name:      "不带协议",
			url:       "github.com/vuejs/vue",
			wantOwner: "vuejs",
			wantRepo:  "vue",
			wantErr:   false,
		},
		{
			name:      "带 .git 后缀",
			url:       "https://github.com/vuejs/vue.git",
			wantOwner: "vuejs",
			wantRepo:  "vue",
			wantErr:   false,
		},
		{
			name:    "非 GitHub URL",
			url:     "https://gitlab.com/vuejs/vue",
			wantErr: true,
		},
		{
			name:    "格式错误",
			url:     "https://github.com/vuejs",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo, err := parseGitHubURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseGitHubURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if owner != tt.wantOwner {
					t.Errorf("parseGitHubURL() owner = %v, want %v", owner, tt.wantOwner)
				}
				if repo != tt.wantRepo {
					t.Errorf("parseGitHubURL() repo = %v, want %v", repo, tt.wantRepo)
				}
			}
		})
	}
}
