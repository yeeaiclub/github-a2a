package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/yeeaiclub/a2a-go/sdk/server/handler"
	"github.com/yeeaiclub/a2a-go/sdk/server/tasks"
	"github.com/yeeaiclub/a2a-go/sdk/types"
)

var AgentCard = types.AgentCard{
	Description: "An A2A-compliant agent that provides GitHub capabilities",
	Version:     "1.0.0",
	DefaultInputModes: []string{
		"text",
	},
	DefaultOutputModes: []string{
		"text",
	},
}

var store = tasks.NewInMemoryTaskStore()

func main() {
	// Create agent configuration
	agentConfig := CreateAgent()

	// Get API key from environment
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		log.Fatal("DEEPSEEK_API_KEY environment variable is required")
	}

	store.Save(context.Background(), &types.Task{Id: "1"})

	defaultHandler := handler.NewDefaultHandler(
		store,
		NewExecutor(store, &AgentCard, agentConfig.Tools, apiKey, agentConfig.SystemPrompt),
		handler.WithQueueManger(NewQueueManager()),
	)

	server := handler.NewServer(
		"/agent_card",
		"/api",
		AgentCard,
		defaultHandler,
		handler.WithWriteTimeout(1*time.Minute),
		handler.WithIdleTimeout(1*time.Minute),
		handler.WithReadTimeout(1*time.Minute),
	)

	server.Start(8080)
}
