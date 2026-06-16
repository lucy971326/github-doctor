package workflow

import (
	"context"
	"testing"
)

// MockStep 是一个用于测试的 mock 步骤
type MockStep struct {
	name    string
	execute func(ctx context.Context, data *AnalysisData) error
}

func (s *MockStep) Name() string {
	return s.name
}

func (s *MockStep) Execute(ctx context.Context, data *AnalysisData) error {
	if s.execute != nil {
		return s.execute(ctx, data)
	}
	return nil
}

func TestStepInterface(t *testing.T) {
	// 验证 MockStep 实现了 Step 接口
	var step Step = &MockStep{name: "test"}

	if step.Name() != "test" {
		t.Errorf("期望 Name() 返回 'test'，得到 %q", step.Name())
	}

	// 测试 Execute
	err := step.Execute(context.Background(), NewAnalysisData())
	if err != nil {
		t.Errorf("Execute() 返回错误: %v", err)
	}
}

func TestStepExecute_WithError(t *testing.T) {
	expectedErr := context.DeadlineExceeded
	step := &MockStep{
		name: "error-step",
		execute: func(ctx context.Context, data *AnalysisData) error {
			return expectedErr
		},
	}

	err := step.Execute(context.Background(), NewAnalysisData())
	if err != expectedErr {
		t.Errorf("期望错误 %v，得到 %v", expectedErr, err)
	}
}
