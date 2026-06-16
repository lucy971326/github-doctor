package llm

import (
	"encoding/json"
	"testing"
)

func TestScoreData_JSON(t *testing.T) {
	data := &ScoreData{
		OverallScore: 85,
		Dimensions: []Dimension{
			{Name: "代码质量", Score: 88, Comment: "代码结构清晰"},
			{Name: "社区活跃度", Score: 75, Comment: "Star 数较多"},
		},
		Summary:         "这是一个成熟的开源项目",
		Recommendations: []string{"建议增加测试覆盖率", "建议添加 CI/CD"},
	}

	// 序列化
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("json.Marshal() 错误: %v", err)
	}

	// 反序列化
	var decoded ScoreData
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("json.Unmarshal() 错误: %v", err)
	}

	// 验证
	if decoded.OverallScore != 85 {
		t.Errorf("OverallScore = %d, want 85", decoded.OverallScore)
	}

	if len(decoded.Dimensions) != 2 {
		t.Errorf("Dimensions 长度 = %d, want 2", len(decoded.Dimensions))
	}

	if decoded.Summary != "这是一个成熟的开源项目" {
		t.Errorf("Summary = %q, want %q", decoded.Summary, "这是一个成熟的开源项目")
	}
}

func TestChatRequest(t *testing.T) {
	req := &ChatRequest{
		Model: "deepseek-chat",
		Messages: []Message{
			{Role: "system", Content: "你是一个代码质量分析专家"},
			{Role: "user", Content: "请分析以下 GitHub 仓库数据"},
		},
		ResponseFormat: &ResponseFormat{
			Type: "json_object",
		},
	}

	if req.Model != "deepseek-chat" {
		t.Errorf("Model = %q, want %q", req.Model, "deepseek-chat")
	}

	if len(req.Messages) != 2 {
		t.Errorf("Messages 长度 = %d, want 2", len(req.Messages))
	}

	if req.ResponseFormat.Type != "json_object" {
		t.Errorf("ResponseFormat.Type = %q, want %q", req.ResponseFormat.Type, "json_object")
	}
}

func TestDimension(t *testing.T) {
	d := &Dimension{
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
}
