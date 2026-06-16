package steps

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github-doctor/workflow"
)

// CloneStep 是克隆仓库的步骤
type CloneStep struct{}

// Name 返回步骤名称
func (s *CloneStep) Name() string {
	return "克隆仓库"
}

// Execute 执行步骤
func (s *CloneStep) Execute(ctx context.Context, data *workflow.AnalysisData) error {
	// 检查仓库大小
	if data.GitHubData != nil && data.GitHubData.Size > 0 {
		sizeMB := float64(data.GitHubData.Size) / 1024

		// 大于 300MB 跳过代码分析和克隆
		if sizeMB > 300 {
			fmt.Printf("    仓库 %.0f MB，跳过克隆\n", sizeMB)
			data.SkipCodeAnalysis = true
			return nil
		}
	}

	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "github-doctor-*")
	if err != nil {
		return fmt.Errorf("创建临时目录失败: %w", err)
	}

	// 构建 clone URL
	cloneURL := fmt.Sprintf("https://github.com/%s/%s.git", data.Owner, data.Repo)

	// 执行 git clone (shallow clone)
	cmd := exec.CommandContext(ctx, "git", "clone", "--depth", "1", cloneURL, filepath.Join(tmpDir, data.Repo))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return fmt.Errorf("git clone 失败: %w", err)
	}

	data.RepoDir = filepath.Join(tmpDir, data.Repo)
	return nil
}

// Cleanup 清理临时目录
func Cleanup(repoDir string) {
	if repoDir != "" {
		// 获取父目录（临时目录）
		parentDir := filepath.Dir(repoDir)
		if strings.HasPrefix(filepath.Base(parentDir), "github-doctor-") {
			os.RemoveAll(parentDir)
		}
	}
}
