package workflow

import (
	"context"
	"testing"
)

func TestNewWorkflow(t *testing.T) {
	w := NewWorkflow()

	if w == nil {
		t.Fatal("NewWorkflow() 返回 nil")
	}

	if len(w.Steps) != 0 {
		t.Errorf("期望 Steps 为空，得到 %d 个步骤", len(w.Steps))
	}

	if w.Data == nil {
		t.Error("期望 Data 不为 nil")
	}
}

func TestWorkflow_AddStep(t *testing.T) {
	w := NewWorkflow()
	step := &MockStep{name: "test-step"}

	w.AddStep(step)

	if len(w.Steps) != 1 {
		t.Errorf("期望 1 个步骤，得到 %d 个", len(w.Steps))
	}

	if w.Steps[0].Name() != "test-step" {
		t.Errorf("步骤名称错误")
	}
}

func TestWorkflow_Run_Success(t *testing.T) {
	w := NewWorkflow()

	// 添加一个成功的步骤
	executed := false
	w.AddStep(&MockStep{
		name: "success-step",
		execute: func(ctx context.Context, data *AnalysisData) error {
			executed = true
			return nil
		},
	})

	err := w.Run("https://github.com/vuejs/vue")
	if err != nil {
		t.Errorf("Run() 返回错误: %v", err)
	}

	if !executed {
		t.Error("步骤未被执行")
	}

	if w.Data.Owner != "vuejs" {
		t.Errorf("期望 Owner 为 'vuejs'，得到 %q", w.Data.Owner)
	}

	if w.Data.Repo != "vue" {
		t.Errorf("期望 Repo 为 'vue'，得到 %q", w.Data.Repo)
	}
}

func TestWorkflow_Run_StepError(t *testing.T) {
	w := NewWorkflow()

	// 添加一个失败的步骤
	w.AddStep(&MockStep{
		name: "fail-step",
		execute: func(ctx context.Context, data *AnalysisData) error {
			return context.DeadlineExceeded
		},
	})

	err := w.Run("https://github.com/vuejs/vue")
	if err == nil {
		t.Error("期望 Run() 返回错误")
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
			name:      "带 www",
			url:       "https://www.github.com/vuejs/vue",
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
		{
			name:    "空路径",
			url:     "https://github.com/",
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
