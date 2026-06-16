package github

// GitHubData 存储 GitHub API 返回的数据
type GitHubData struct {
	Stars        int            `json:"stargazers_count"`
	Forks        int            `json:"forks_count"`
	Watchers     int            `json:"subscribers_count"`
	OpenIssues   int            `json:"open_issues_count"`
	License      string         `json:"license_name"`
	Description  string         `json:"description"`
	Languages    map[string]int `json:"languages"`
	CreatedAt    string         `json:"created_at"`
	UpdatedAt    string         `json:"updated_at"`
	DefaultBranch string        `json:"default_branch"`
	Size         int            `json:"size"`
}

// RepositoryInfo 存储仓库基本信息
type RepositoryInfo struct {
	Owner string
	Repo  string
	URL   string
}
