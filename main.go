package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github-doctor/codegraph"
	"github-doctor/github"
	"github-doctor/workflow"
	"github-doctor/workflow/steps"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println()
	fmt.Println("  ╔══════════════════════════════════════╗")
	fmt.Println("  ║     GitHub Doctor - 仓库体检工具     ║")
	fmt.Println("  ╚══════════════════════════════════════╝")
	fmt.Println()

	// 1. 输入 GitHub URL
	fmt.Print("  仓库 URL: ")
	scanner.Scan()
	url := strings.TrimSpace(scanner.Text())

	if url == "" {
		fmt.Println("  ✗ URL 不能为空")
		os.Exit(1)
	}

	// 2. 询问 GitHub Token
	fmt.Print("  GitHub Token (回车跳过): ")
	scanner.Scan()
	githubToken := strings.TrimSpace(scanner.Text())

	// 3. 询问代理
	fmt.Print("  代理端口 (回车跳过):   ")
	scanner.Scan()
	proxyPort := strings.TrimSpace(scanner.Text())

	// 设置代理
	if proxyPort != "" {
		proxyURL := fmt.Sprintf("http://127.0.0.1:%s", proxyPort)
		os.Setenv("HTTPS_PROXY", proxyURL)
		os.Setenv("HTTP_PROXY", proxyURL)
	}

	// 4. 询问是否需要代码质量分析
	fmt.Print("  深度分析? (y/n):      ")
	scanner.Scan()
	choice := strings.TrimSpace(strings.ToLower(scanner.Text()))

	needCodeAnalysis := choice == "y" || choice == "yes"

	// 5. DeepSeek API Key
	var deepseekKey string
	if needCodeAnalysis {
		fmt.Print("  DeepSeek API Key:     ")
		scanner.Scan()
		deepseekKey = strings.TrimSpace(scanner.Text())

		if deepseekKey == "" {
			fmt.Println("  └─ Key 为空，跳过 AI 评分")
			needCodeAnalysis = false
		}
	}

	fmt.Println()
	fmt.Println("  ─────────────────────────────────────")

	// 创建工作流
	w := workflow.NewWorkflow()

	// 创建 GitHub 客户端
	var githubClient *github.Client
	if githubToken != "" {
		githubClient = github.NewClientWithToken(githubToken)
		fmt.Println("  [✓] GitHub Token 已配置")
	}

	// 添加基础步骤
	w.AddStep(&steps.ValidateStep{})
	w.AddStep(steps.NewFetchGitHubStepWithClient(githubClient))
	w.AddStep(&steps.CloneStep{})

	// 代码分析
	if needCodeAnalysis {
		binPath, err := codegraph.EnsureCodegraph()
		if err != nil {
			fmt.Printf("  [✗] codegraph 不可用: %v\n", err)
			os.Exit(1)
		}
		w.AddStep(steps.NewAnalyzeCodegraphStepWithBinPath(binPath))
		w.AddStep(steps.NewScoreLLMStepWithKey(deepseekKey))
	}

	w.AddStep(steps.NewGenerateReportStep())
	w.Data.NeedCodeAnalysis = needCodeAnalysis

	// 执行工作流
	if err := w.Run(url); err != nil {
		fmt.Printf("  [✗] 失败: %v\n", err)
		os.Exit(1)
	}

	defer steps.Cleanup(w.Data.RepoDir)

	fmt.Println()
	fmt.Println("  ─────────────────────────────────────")
	fmt.Println("  [✓] 完成！报告已在浏览器中打开")
	fmt.Println()
}
