package codegraph

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Client 是 codegraph CLI 客户端
type Client struct {
	binPath string
}

// NewClient 创建新的 codegraph 客户端
func NewClient() *Client {
	binPath := findCodegraphBinary()
	return &Client{binPath: binPath}
}

// NewClientWithPath 使用指定路径创建客户端
func NewClientWithPath(binPath string) *Client {
	return &Client{binPath: binPath}
}

// findCodegraphBinary 查找 codegraph 二进制文件
func findCodegraphBinary() string {
	// 1. 检查 PATH
	if path, err := exec.LookPath("codegraph"); err == nil {
		return path
	}

	// 2. 检查本地 bin/
	execPath, _ := os.Executable()
	binDir := filepath.Dir(execPath)

	// Windows: 检查 node.exe + codegraph.js
	if runtime.GOOS == "windows" {
		nodeExe := filepath.Join(binDir, "bin", "codegraph-win32-x64", "node.exe")
		codegraphJS := filepath.Join(binDir, "bin", "codegraph-win32-x64", "lib", "dist", "bin", "codegraph.js")
		if _, err := os.Stat(nodeExe); err == nil {
			if _, err := os.Stat(codegraphJS); err == nil {
				return nodeExe
			}
		}
	}

	// 3. 返回默认路径
	binName := "codegraph"
	if runtime.GOOS == "windows" {
		binName = "codegraph.exe"
	}
	return filepath.Join(binDir, "bin", binName)
}

// EnsureCodegraph 确保 codegraph 可用
func EnsureCodegraph() (string, error) {
	// 1. 检查 PATH
	if path, err := exec.LookPath("codegraph"); err == nil {
		return path, nil
	}

	// 2. 检查本地
	execPath, _ := os.Executable()
	binDir := filepath.Dir(execPath)

	// Windows: 检查 node.exe + codegraph.js
	if runtime.GOOS == "windows" {
		nodeExe := filepath.Join(binDir, "bin", "codegraph-win32-x64", "node.exe")
		codegraphJS := filepath.Join(binDir, "bin", "codegraph-win32-x64", "lib", "dist", "bin", "codegraph.js")
		if _, err := os.Stat(nodeExe); err == nil {
			if _, err := os.Stat(codegraphJS); err == nil {
				return nodeExe, nil
			}
		}
	}

	// 3. 下载
	fmt.Println("    📥 正在下载 codegraph...")
	if err := downloadCodegraph(binDir); err != nil {
		return "", fmt.Errorf("下载 codegraph 失败: %w", err)
	}

	// 4. 再次检查
	if runtime.GOOS == "windows" {
		nodeExe := filepath.Join(binDir, "bin", "codegraph-win32-x64", "node.exe")
		if _, err := os.Stat(nodeExe); err == nil {
			return nodeExe, nil
		}
	}

	binName := "codegraph"
	if runtime.GOOS == "windows" {
		binName = "codegraph.exe"
	}
	return filepath.Join(binDir, "bin", binName), nil
}

// downloadCodegraph 下载 codegraph
func downloadCodegraph(binDir string) error {
	binPath := filepath.Join(binDir, "bin")
	if err := os.MkdirAll(binPath, 0755); err != nil {
		return fmt.Errorf("创建 bin 目录失败: %w", err)
	}

	var downloadURL string
	switch runtime.GOOS {
	case "windows":
		if runtime.GOARCH == "arm64" {
			downloadURL = "https://github.com/colbymchenry/codegraph/releases/latest/download/codegraph-win32-arm64.zip"
		} else {
			downloadURL = "https://github.com/colbymchenry/codegraph/releases/latest/download/codegraph-win32-x64.zip"
		}
	case "darwin":
		if runtime.GOARCH == "arm64" {
			downloadURL = "https://github.com/colbymchenry/codegraph/releases/latest/download/codegraph-darwin-arm64.tar.gz"
		} else {
			downloadURL = "https://github.com/colbymchenry/codegraph/releases/latest/download/codegraph-darwin-x64.tar.gz"
		}
	default:
		if runtime.GOARCH == "arm64" {
			downloadURL = "https://github.com/colbymchenry/codegraph/releases/latest/download/codegraph-linux-arm64.tar.gz"
		} else {
			downloadURL = "https://github.com/colbymchenry/codegraph/releases/latest/download/codegraph-linux-x64.tar.gz"
		}
	}

	tmpFile := filepath.Join(binPath, "codegraph.tmp")
	if err := downloadFile(downloadURL, tmpFile); err != nil {
		return fmt.Errorf("下载失败: %w", err)
	}
	defer os.Remove(tmpFile)

	fmt.Println("    📦 正在解压...")
	if err := extractArchive(tmpFile, binPath); err != nil {
		return fmt.Errorf("解压失败: %w", err)
	}

	return nil
}

// downloadFile 下载文件
func downloadFile(url, dest string) error {
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败: HTTP %d", resp.StatusCode)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// extractArchive 解压文件
func extractArchive(archive, destDir string) error {
	if strings.HasSuffix(archive, ".zip") {
		return extractZip(archive, destDir)
	}
	return extractTarGz(archive, destDir)
}

// extractZip 解压 zip 文件
func extractZip(archive, destDir string) error {
	r, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(destDir, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("非法路径: %s", fpath)
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, 0755)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return err
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// extractTarGz 解压 tar.gz 文件
func extractTarGz(archive, destDir string) error {
	gz, err := os.Open(archive)
	if err != nil {
		return err
	}
	defer gz.Close()
	gzr, err := gzip.NewReader(gz)
	if err != nil {
		return err
	}
	defer gzr.Close()
	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fpath := filepath.Join(destDir, header.Name)
		if header.FileInfo().IsDir() {
			os.MkdirAll(fpath, 0755)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return err
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
		if err != nil {
			return err
		}
		_, err = io.Copy(outFile, tr)
		outFile.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// runCodegraph 运行 codegraph 命令
func (c *Client) runCodegraph(args ...string) ([]byte, error) {
	if strings.HasSuffix(c.binPath, "node.exe") {
		dir := filepath.Dir(c.binPath)
		codegraphJS := filepath.Join(dir, "lib", "dist", "bin", "codegraph.js")
		if _, err := os.Stat(codegraphJS); err == nil {
			allArgs := append([]string{"--liftoff-only", codegraphJS}, args...)
			cmd := exec.Command(c.binPath, allArgs...)
			return cmd.Output()
		}
	}
	cmd := exec.Command(c.binPath, args...)
	return cmd.Output()
}

func (c *Client) runCodegraphCombined(args ...string) ([]byte, error) {
	if strings.HasSuffix(c.binPath, "node.exe") {
		dir := filepath.Dir(c.binPath)
		codegraphJS := filepath.Join(dir, "lib", "dist", "bin", "codegraph.js")
		if _, err := os.Stat(codegraphJS); err == nil {
			allArgs := append([]string{"--liftoff-only", codegraphJS}, args...)
			cmd := exec.Command(c.binPath, allArgs...)
			return cmd.CombinedOutput()
		}
	}
	cmd := exec.Command(c.binPath, args...)
	return cmd.CombinedOutput()
}

// Init 初始化 codegraph 索引
func (c *Client) Init(repoDir string) error {
	output, err := c.runCodegraphCombined("init", repoDir)
	if err != nil {
		return fmt.Errorf("codegraph init 失败: %w\n输出: %s", err, output)
	}
	return nil
}

// Status 获取统计信息
func (c *Client) Status(repoDir string) (*StatusResult, error) {
	output, err := c.runCodegraph("status", repoDir, "--json")
	if err != nil {
		return nil, fmt.Errorf("codegraph status 失败: %w", err)
	}
	var result StatusResult
	if err := json.Unmarshal(output, &result); err != nil {
		return c.parseStatusText(string(output))
	}
	return &result, nil
}

// Files 获取文件结构
func (c *Client) Files(repoDir string) (string, error) {
	output, err := c.runCodegraph("files", "--path", repoDir)
	if err != nil {
		return "", fmt.Errorf("codegraph files 失败: %w", err)
	}
	return string(output), nil
}

// Query 搜索代码符号
func (c *Client) Query(repoDir, search string, limit int) (string, error) {
	output, err := c.runCodegraph("query", "--path", repoDir, search, "--limit", fmt.Sprintf("%d", limit))
	if err != nil {
		return "", fmt.Errorf("codegraph query 失败: %w", err)
	}
	return string(output), nil
}

// StatusResult 存储 codegraph 状态结果
type StatusResult struct {
	FilesIndexed int            `json:"files_indexed"`
	TotalNodes   int            `json:"total_nodes"`
	TotalEdges   int            `json:"total_edges"`
	NodesByKind  map[string]int `json:"nodes_by_kind"`
	Languages    map[string]int `json:"languages"`
}

func (c *Client) parseStatusText(text string) (*StatusResult, error) {
	result := &StatusResult{
		NodesByKind: make(map[string]int),
		Languages:   make(map[string]int),
	}
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "**Files indexed:**") {
			fmt.Sscanf(line, "**Files indexed:** %d", &result.FilesIndexed)
		}
		if strings.HasPrefix(line, "**Total nodes:**") {
			fmt.Sscanf(line, "**Total nodes:** %d", &result.TotalNodes)
		}
		if strings.HasPrefix(line, "**Total edges:**") {
			fmt.Sscanf(line, "**Total edges:** %d", &result.TotalEdges)
		}
	}
	return result, nil
}
