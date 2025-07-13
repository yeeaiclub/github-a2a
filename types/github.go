package types

import "time"

type GitHubUser struct {
	Login string  `json:"login"`
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}

// GitHubRepository represents GitHub repository information
type GitHubRepository struct {
	Name        string     `json:"name"`
	FullName    string     `json:"full_name"`
	Description *string    `json:"description,omitempty"`
	URL         string     `json:"url"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PushedAt    *time.Time `json:"pushed_at,omitempty"`
	Language    *string    `json:"language,omitempty"`
	Stars       int        `json:"stars"`
	Forks       int        `json:"forks"`
}

// GitHubCommit represents GitHub commit information
type GitHubCommit struct {
	SHA     string    `json:"sha"`
	Message string    `json:"message"`
	Author  string    `json:"author"`
	Date    time.Time `json:"date"`
	URL     string    `json:"url"`
}

// GitHubResponse represents base response model for GitHub API operations
type GitHubResponse struct {
	Status       string  `json:"status"`
	Message      string  `json:"message"`
	Count        *int    `json:"count,omitempty"`
	ErrorMessage *string `json:"error_message,omitempty"`
}

// RepositoryResponse represents response model for repository operations
type RepositoryResponse struct {
	GitHubResponse
	Data []GitHubRepository `json:"data,omitempty"`
}

// CommitResponse represents response model for commit operations
type CommitResponse struct {
	GitHubResponse
	Data []GitHubCommit `json:"data,omitempty"`
}

// ToolFunction represents a tool function for OpenAI function calling
type ToolFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// Tool represents a tool for OpenAI function calling
type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}
