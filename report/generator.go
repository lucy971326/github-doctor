package report

import (
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// Generator 是报告生成器
type Generator struct {
	templatePath string
}

// NewGenerator 创建新的报告生成器
func NewGenerator() *Generator {
	// 查找模板文件
	execPath, _ := os.Executable()
	binDir := filepath.Dir(execPath)
	templatePath := filepath.Join(binDir, "templates", "report.html")

	// 如果模板不存在，使用当前目录
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		templatePath = "templates/report.html"
	}

	return &Generator{
		templatePath: templatePath,
	}
}

// Generate 生成 HTML 报告
func (g *Generator) Generate(data *ReportData, outputPath string) error {
	// 设置生成时间
	data.GeneratedAt = time.Now().Format("2006-01-02 15:04:05")

	// 自定义模板函数
	funcMap := template.FuncMap{
		"mul": func(a, b float64) float64 {
			return a * b
		},
		"sub": func(a, b float64) float64 {
			return a - b
		},
		"float64": func(v int) float64 {
			return float64(v)
		},
		"scoreOffset": func(score int) float64 {
			// Circumference = 2 * PI * r = 2 * 3.14159 * 65 ≈ 408
			// offset = circumference * (1 - score/100)
			return 408 * (1 - float64(score)/100)
		},
		"langColor": func(lang string) string {
			colors := map[string]string{
				"JavaScript":  "#f1e05a",
				"TypeScript":  "#3178c6",
				"Python":      "#3572A5",
				"Java":        "#b07219",
				"Go":          "#00ADD8",
				"Rust":        "#dea584",
				"C++":         "#f34b7d",
				"C":           "#555555",
				"C#":          "#178600",
				"Ruby":        "#701516",
				"PHP":         "#4F5D95",
				"Swift":       "#F05138",
				"Kotlin":      "#A97BFF",
				"Dart":        "#00B4AB",
				"Lua":         "#000080",
				"Shell":       "#89e051",
				"HTML":        "#e34c26",
				"CSS":         "#563d7c",
				"Vue":         "#41b883",
				"Svelte":      "#ff3e00",
			}
			if c, ok := colors[lang]; ok {
				return c
			}
			return "#6b7280"
		},
	}

	// 读取模板
	tmpl, err := template.New("report.html").Funcs(funcMap).ParseFiles(g.templatePath)
	if err != nil {
		return fmt.Errorf("解析模板失败: %w", err)
	}

	// 创建输出文件
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	// 执行模板
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("执行模板失败: %w", err)
	}

	return nil
}

// OpenInBrowser 在浏览器中打开报告
func OpenInBrowser(path string) error {
	// 获取绝对路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("获取路径失败: %w", err)
	}

	// 转换为 file:// URL
	url := "file:///" + filepath.ToSlash(absPath)

	// 根据操作系统打开浏览器
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	return cmd.Start()
}
