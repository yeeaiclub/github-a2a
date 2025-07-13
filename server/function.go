package main

import (
	"fmt"

	"github.com/cohesion-org/deepseek-go"
)

type Function interface {
	OpenAIFunctionDefinition() deepseek.Function
	Call(args map[string]interface{}) interface{}
}

type GetUserRepositoriesTool struct {
	toolset *GitHubToolset
}

func (t *GetUserRepositoriesTool) OpenAIFunctionDefinition() deepseek.Function {
	return deepseek.Function{
		Name:        "get_user_repositories",
		Description: "Get user's repository list with filtering by recent update time",
		Parameters: &deepseek.FunctionParameters{
			Type: "object",
			Properties: map[string]interface{}{
				"username": map[string]interface{}{
					"type":        "string",
					"description": "GitHub username, if not provided will get repositories of the current authenticated user",
				},
				"days": map[string]interface{}{
					"type":        "integer",
					"description": "Filter repositories updated within how many days, default is 30 days",
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "Limit the number of returned results, default is 10",
				},
			},
		},
	}
}

func (t *GetUserRepositoriesTool) Call(args map[string]interface{}) interface{} {
	var username *string
	var days *int
	var limit *int

	if val, ok := args["username"].(string); ok && val != "" {
		username = &val
	}
	if val, ok := args["days"].(float64); ok {
		intVal := int(val)
		days = &intVal
	}
	if val, ok := args["limit"].(float64); ok {
		intVal := int(val)
		limit = &intVal
	}

	result := t.toolset.GetUserRepositories(username, days, limit)
	fmt.Printf("GetUserRepositories result: %+v\n", result)
	return result
}

type GetRecentCommitsTool struct {
	toolset *GitHubToolset
}

func (t *GetRecentCommitsTool) OpenAIFunctionDefinition() deepseek.Function {
	return deepseek.Function{
		Name:        "get_recent_commits",
		Description: "Get recent commit records for a specific repository",
		Parameters: &deepseek.FunctionParameters{
			Type: "object",
			Properties: map[string]interface{}{
				"repoName": map[string]interface{}{
					"type":        "string",
					"description": "Repository name in format 'owner/repo', e.g. 'microsoft/vscode'",
				},
				"days": map[string]interface{}{
					"type":        "integer",
					"description": "Get commits within how many days, default is 7 days",
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "Limit the number of returned results, default is 10",
				},
			},
			Required: []string{"repoName"},
		},
	}
}

func (t *GetRecentCommitsTool) Call(args map[string]interface{}) interface{} {
	repoName, ok := args["repoName"].(string)
	if !ok {
		return map[string]string{"error": "repoName is required"}
	}

	var days *int
	var limit *int

	if val, ok := args["days"].(float64); ok {
		intVal := int(val)
		days = &intVal
	}
	if val, ok := args["limit"].(float64); ok {
		intVal := int(val)
		limit = &intVal
	}

	result := t.toolset.GetRecentCommits(repoName, days, limit)
	return result
}

type SearchRepositoriesTool struct {
	toolset *GitHubToolset
}

func (t *SearchRepositoriesTool) OpenAIFunctionDefinition() deepseek.Function {
	return deepseek.Function{
		Name:        "search_repositories",
		Description: "Search repositories with recent activity",
		Parameters: &deepseek.FunctionParameters{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Search keywords, e.g. 'machine learning' or 'react'",
				},
				"sort": map[string]interface{}{
					"type":        "string",
					"description": "Sort method, options: 'stars', 'forks', 'updated', default is 'updated'",
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "Limit the number of returned results, default is 10",
				},
			},
			Required: []string{"query"},
		},
	}
}

func (t *SearchRepositoriesTool) Call(args map[string]interface{}) interface{} {
	query, ok := args["query"].(string)
	if !ok {
		return map[string]string{"error": "query is required"}
	}

	var sort *string
	var limit *int

	if val, ok := args["sort"].(string); ok {
		sort = &val
	}
	if val, ok := args["limit"].(float64); ok {
		intVal := int(val)
		limit = &intVal
	}

	result := t.toolset.SearchRepositories(query, sort, limit)
	fmt.Printf("SearchRepositories result: %+v\n", result)
	return result
}
