package report

// ReportData 是 HTML 报告的数据结构
type ReportData struct {
	RepoInfo      *RepoInfo      `json:"repo_info"`
	GitHubStats   *GitHubStats   `json:"github_stats"`
	ScoreResult   *ScoreResult   `json:"score_result"`
	GeneratedAt   string         `json:"generated_at"`
}

// RepoInfo 是仓库基本信息
type RepoInfo struct {
	URL    string `json:"url"`
	Owner  string `json:"owner"`
	Repo   string `json:"repo"`
}

// GitHubStats 是 GitHub 统计数据
type GitHubStats struct {
	Stars      int            `json:"stars"`
	Forks      int            `json:"forks"`
	Watchers   int            `json:"watchers"`
	OpenIssues int            `json:"open_issues"`
	License    string         `json:"license"`
	Languages  map[string]int `json:"languages"`
	CreatedAt  string         `json:"created_at"`
	UpdatedAt  string         `json:"updated_at"`
}

// ScoreResult 是评分结果
type ScoreResult struct {
	OverallScore    int          `json:"overall_score"`
	Dimensions      []Dimension  `json:"dimensions"`
	Summary         string       `json:"summary"`
	Recommendations []string     `json:"recommendations"`
}

// Dimension 是单个维度的评分
type Dimension struct {
	Name    string `json:"name"`
	Score   int    `json:"score"`
	Comment string `json:"comment"`
}
