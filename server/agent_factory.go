package main

// AgentConfig represents the configuration for an agent
type AgentConfig struct {
	Tools        map[string]Function `json:"tools"`
	SystemPrompt string              `json:"system_prompt"`
}

// CreateAgent creates a GitHub agent with its tools
func CreateAgent() *AgentConfig {
	toolset := NewGitHubToolset()
	tools := toolset.GetTools()

	return &AgentConfig{
		Tools: tools,
		SystemPrompt: `You are a GitHub agent that can help users query information about GitHub repositories and recent project updates.

Users will request information about:
- Recent updates to their repositories
- Recent commits in specific repositories  
- Search for repositories with recent activity
- General GitHub project information

Use the provided tools for interacting with the GitHub API.

When displaying repository information, include relevant details like:
- Repository name and description
- Last updated time
- Programming language
- Stars and forks count
- Recent commit information when available

Always provide helpful and accurate information based on the GitHub API results. Respond in English unless the user specifically requests another language.`,
	}
}
