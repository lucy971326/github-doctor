package codegraph

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("NewClient() 返回 nil")
	}
	if client.binPath == "" {
		t.Error("binPath 为空")
	}
}

func TestStatusResult(t *testing.T) {
	result := &StatusResult{
		FilesIndexed: 245,
		TotalNodes:   4005,
		TotalEdges:   11639,
		NodesByKind: map[string]int{
			"function": 1058,
			"method":   838,
			"class":    55,
		},
		Languages: map[string]int{
			"typescript": 225,
			"javascript": 18,
		},
	}

	if result.FilesIndexed != 245 {
		t.Errorf("FilesIndexed = %d, want 245", result.FilesIndexed)
	}

	if result.TotalNodes != 4005 {
		t.Errorf("TotalNodes = %d, want 4005", result.TotalNodes)
	}

	if len(result.NodesByKind) != 3 {
		t.Errorf("NodesByKind 长度 = %d, want 3", len(result.NodesByKind))
	}
}

func TestParseStatusText(t *testing.T) {
	client := &Client{}
	text := `## CodeGraph Status

**Files indexed:** 245
**Total nodes:** 4005
**Total edges:** 11639`

	result, err := client.parseStatusText(text)
	if err != nil {
		t.Fatalf("parseStatusText() 错误: %v", err)
	}

	if result.FilesIndexed != 245 {
		t.Errorf("FilesIndexed = %d, want 245", result.FilesIndexed)
	}

	if result.TotalNodes != 4005 {
		t.Errorf("TotalNodes = %d, want 4005", result.TotalNodes)
	}

	if result.TotalEdges != 11639 {
		t.Errorf("TotalEdges = %d, want 11639", result.TotalEdges)
	}
}
