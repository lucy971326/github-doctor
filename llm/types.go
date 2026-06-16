package llm

// ScoreData 存储 LLM 评分结果
type ScoreData struct {
	OverallScore    int          `json:"overall_score"`
	Dimensions      []Dimension  `json:"dimensions"`
	Summary         string       `json:"summary"`
	Recommendations []string     `json:"recommendations"`
}

// Dimension 存储单个维度的评分
type Dimension struct {
	Name    string `json:"name"`
	Score   int    `json:"score"`
	Comment string `json:"comment"`
}

// ChatRequest 是 DeepSeek API 的请求格式
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"`
}

// Message 是聊天消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ResponseFormat 指定响应格式
type ResponseFormat struct {
	Type string `json:"type"`
}

// ChatResponse 是 DeepSeek API 的响应格式
type ChatResponse struct {
	Choices []Choice `json:"choices"`
}

// Choice 是单个选择
type Choice struct {
	Message Message `json:"message"`
}

// AnalysisInput 是发送给 LLM 的分析输入
type AnalysisInput struct {
	RepoURL     string            `json:"repo_url"`
	Owner       string            `json:"owner"`
	Repo        string            `json:"repo"`
	GitHubData  *GitHubInput      `json:"github_data"`
	CodegraphData *CodegraphInput `json:"codegraph_data"`
}

// GitHubInput 是 GitHub 数据的输入格式
type GitHubInput struct {
	Stars      int            `json:"stars"`
	Forks      int            `json:"forks"`
	OpenIssues int            `json:"open_issues"`
	License    string         `json:"license"`
	Languages  map[string]int `json:"languages"`
	CreatedAt  string         `json:"created_at"`
	UpdatedAt  string         `json:"updated_at"`
}

// CodegraphInput 是 codegraph 数据的输入格式
type CodegraphInput struct {
	Statistics *StatisticsInput `json:"statistics"`
	FileTree   string           `json:"file_tree"`
	CoreCode   string           `json:"core_code"`
}

// StatisticsInput 是统计数据的输入格式
type StatisticsInput struct {
	Files     int `json:"files"`
	Functions int `json:"functions"`
	Classes   int `json:"classes"`
	Methods   int `json:"methods"`
	Edges     int `json:"edges"`
}
