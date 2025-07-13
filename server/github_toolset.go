package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v62/github"
	"github.com/yeeaiclub/github-a2a/types"
	"golang.org/x/oauth2"
)

// GitHubToolset provides GitHub API tools for querying repositories and recent updates
type GitHubToolset struct {
	client *github.Client
}

// NewGitHubToolset creates a new GitHub toolset instance
func NewGitHubToolset() *GitHubToolset {
	toolset := &GitHubToolset{}
	toolset.initClient()
	return toolset
}

// initClient initializes the GitHub client with authentication
func (g *GitHubToolset) initClient() {
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken != "" {
		// Use authenticated client
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: githubToken},
		)
		tc := oauth2.NewClient(context.Background(), ts)
		g.client = github.NewClient(tc)
	} else {
		// Use unauthenticated client (limited rate)
		fmt.Println("Warning: No GITHUB_TOKEN found, using unauthenticated access (limited rate)")
		g.client = github.NewClient(nil)
	}
}

// GetUserRepositories gets user's repositories with recent updates
func (g *GitHubToolset) GetUserRepositories(username *string, days *int, limit *int) types.RepositoryResponse {
	// Set default values
	if days == nil {
		defaultDays := 30
		days = &defaultDays
	}
	if limit == nil {
		defaultLimit := 10
		limit = &defaultLimit
	}

	cutoffDate := time.Now().AddDate(0, 0, -*days)
	var user *github.User
	var err error

	if username != nil && *username != "" {
		// Get specific user
		user, _, err = g.client.Users.Get(context.Background(), *username)
	} else {
		// Get authenticated user
		user, _, err = g.client.Users.Get(context.Background(), "")
	}

	if err != nil {
		errorMsg := fmt.Sprintf("Failed to get user: %v", err)
		return types.RepositoryResponse{
			GitHubResponse: types.GitHubResponse{
				Status:       "error",
				Message:      errorMsg,
				ErrorMessage: &errorMsg,
			},
		}
	}

	// Get repositories
	opt := &github.RepositoryListOptions{
		Sort:      "updated",
		Direction: "desc",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var allRepos []*github.Repository
	for {
		repos, resp, err := g.client.Repositories.List(context.Background(), *user.Login, opt)
		if err != nil {
			errorMsg := fmt.Sprintf("Failed to get repositories: %v", err)
			return types.RepositoryResponse{
				GitHubResponse: types.GitHubResponse{
					Status:       "error",
					Message:      errorMsg,
					ErrorMessage: &errorMsg,
				},
			}
		}

		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	// Filter repositories by update date and limit
	var filteredRepos []types.GitHubRepository
	count := 0
	for _, repo := range allRepos {
		if count >= *limit {
			break
		}

		if repo.UpdatedAt != nil && repo.UpdatedAt.After(cutoffDate) {
			githubRepo := types.GitHubRepository{
				Name:      *repo.Name,
				FullName:  *repo.FullName,
				URL:       *repo.HTMLURL,
				UpdatedAt: *repo.UpdatedAt.GetTime(),
				Stars:     *repo.StargazersCount,
				Forks:     *repo.ForksCount,
			}

			if repo.Description != nil {
				githubRepo.Description = repo.Description
			}
			if repo.PushedAt != nil {
				githubRepo.PushedAt = repo.PushedAt.GetTime()
			}
			if repo.Language != nil {
				githubRepo.Language = repo.Language
			}

			filteredRepos = append(filteredRepos, githubRepo)
			count++
		}
	}

	message := fmt.Sprintf("Successfully retrieved %d repositories updated in the last %d days", count, *days)
	return types.RepositoryResponse{
		GitHubResponse: types.GitHubResponse{
			Status:  "success",
			Message: message,
			Count:   &count,
		},
		Data: filteredRepos,
	}
}

// GetRecentCommits gets recent commits for a repository
func (g *GitHubToolset) GetRecentCommits(repoName string, days *int, limit *int) types.CommitResponse {
	// Set default values
	if days == nil {
		defaultDays := 7
		days = &defaultDays
	}
	if limit == nil {
		defaultLimit := 10
		limit = &defaultLimit
	}

	// Parse repository name
	parts := strings.Split(repoName, "/")
	if len(parts) != 2 {
		errorMsg := "Repository name must be in format 'owner/repo'"
		return types.CommitResponse{
			GitHubResponse: types.GitHubResponse{
				Status:       "error",
				Message:      errorMsg,
				ErrorMessage: &errorMsg,
			},
		}
	}

	owner, repo := parts[0], parts[1]
	cutoffDate := time.Now().AddDate(0, 0, -*days)

	// Get commits
	opt := &github.CommitsListOptions{
		Since: cutoffDate,
		ListOptions: github.ListOptions{
			PerPage: *limit,
		},
	}

	commits, _, err := g.client.Repositories.ListCommits(context.Background(), owner, repo, opt)
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to get commits: %v", err)
		return types.CommitResponse{
			GitHubResponse: types.GitHubResponse{
				Status:       "error",
				Message:      errorMsg,
				ErrorMessage: &errorMsg,
			},
		}
	}

	// Convert to our format
	var githubCommits []types.GitHubCommit
	for _, commit := range commits {
		if len(githubCommits) >= *limit {
			break
		}

		// Get first line of commit message
		message := *commit.Commit.Message
		if idx := strings.Index(message, "\n"); idx != -1 {
			message = message[:idx]
		}

		githubCommit := types.GitHubCommit{
			SHA:     (*commit.SHA)[:8], // First 8 characters
			Message: message,
			Author:  *commit.Commit.Author.Name,
			Date:    *commit.Commit.Author.Date.GetTime(),
			URL:     *commit.HTMLURL,
		}

		githubCommits = append(githubCommits, githubCommit)
	}

	count := len(githubCommits)
	message := fmt.Sprintf("Successfully retrieved %d commits for repository %s in the last %d days", count, repoName, *days)
	return types.CommitResponse{
		GitHubResponse: types.GitHubResponse{
			Status:  "success",
			Message: message,
			Count:   &count,
		},
		Data: githubCommits,
	}
}

// SearchRepositories searches for repositories with recent activity
func (g *GitHubToolset) SearchRepositories(query string, sort *string, limit *int) types.RepositoryResponse {
	// Set default values
	if sort == nil {
		defaultSort := "updated"
		sort = &defaultSort
	}
	if limit == nil {
		defaultLimit := 10
		limit = &defaultLimit
	}

	// Add recent activity filter to query
	recentDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	searchQuery := fmt.Sprintf("%s pushed:>=%s", query, recentDate)

	// Search repositories
	opt := &github.SearchOptions{
		Sort:  *sort,
		Order: "desc",
		ListOptions: github.ListOptions{
			PerPage: *limit,
		},
	}

	result, _, err := g.client.Search.Repositories(context.Background(), searchQuery, opt)
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to search repositories: %v", err)
		return types.RepositoryResponse{
			GitHubResponse: types.GitHubResponse{
				Status:       "error",
				Message:      errorMsg,
				ErrorMessage: &errorMsg,
			},
		}
	}

	// Convert to our format
	var repos []types.GitHubRepository
	for _, repo := range result.Repositories {
		if len(repos) >= *limit {
			break
		}

		githubRepo := types.GitHubRepository{
			Name:      *repo.Name,
			FullName:  *repo.FullName,
			URL:       *repo.HTMLURL,
			UpdatedAt: *repo.UpdatedAt.GetTime(),
			Stars:     *repo.StargazersCount,
			Forks:     *repo.ForksCount,
		}

		if repo.Description != nil {
			githubRepo.Description = repo.Description
		}
		if repo.PushedAt != nil {
			githubRepo.PushedAt = repo.PushedAt.GetTime()
		}
		if repo.Language != nil {
			githubRepo.Language = repo.Language
		}

		repos = append(repos, githubRepo)
	}

	count := len(repos)
	message := fmt.Sprintf("Successfully searched for %d repositories matching \"%s\"", count, query)
	return types.RepositoryResponse{
		GitHubResponse: types.GitHubResponse{
			Status:  "success",
			Message: message,
			Count:   &count,
		},
		Data: repos,
	}
}

// GetTools returns the available tools for OpenAI function calling
func (g *GitHubToolset) GetTools() map[string]Function {
	return map[string]Function{
		"get_user_repositories": &GetUserRepositoriesTool{g},
		"get_recent_commits":    &GetRecentCommitsTool{g},
		"search_repositories":   &SearchRepositoriesTool{g},
	}
}

// Tool wrappers for go-openai compatibility
