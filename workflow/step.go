package workflow

import "context"

// Step 定义工作流步骤的接口
type Step interface {
	// Name 返回步骤名称
	Name() string

	// Execute 执行步骤
	Execute(ctx context.Context, data *AnalysisData) error
}
