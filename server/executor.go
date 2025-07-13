package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/cohesion-org/deepseek-go"

	"github.com/yeeaiclub/a2a-go/sdk/server/event"
	"github.com/yeeaiclub/a2a-go/sdk/server/execution"
	"github.com/yeeaiclub/a2a-go/sdk/server/tasks"
	"github.com/yeeaiclub/a2a-go/sdk/server/tasks/updater"
	"github.com/yeeaiclub/a2a-go/sdk/types"
)

type DeepSeekExecutor struct {
	store        tasks.TaskStore
	card         *types.AgentCard
	tools        map[string]Function
	apiKey       string
	systemPrompt string
	client       *deepseek.Client
	model        string
}

func NewExecutor(store tasks.TaskStore, card *types.AgentCard, tools map[string]Function, apiKey string, systemPrompt string) *DeepSeekExecutor {
	client := deepseek.NewClient(apiKey)
	log.Printf("Initializing DeepSeekExecutor")

	return &DeepSeekExecutor{
		store:        store,
		card:         card,
		tools:        tools,
		apiKey:       apiKey,
		systemPrompt: systemPrompt,
		client:       client,
		model:        "deepseek-chat",
	}
}

func (e *DeepSeekExecutor) Execute(ctx context.Context, requestContext *execution.RequestContext, queue *event.Queue) error {
	u := updater.NewTaskUpdater(queue, requestContext.TaskId, requestContext.ContextId)
	if requestContext.Task == nil {
		u.Submit()
	}
	u.StartWork()

	messageText := ""
	for _, msg := range requestContext.Params.Message.Parts {
		if msg.GetKind() == "text" {
			part := msg.(*types.TextPart)
			messageText += part.Text
		}
	}

	return e.processRequest(ctx, messageText, u)
}

func (e *DeepSeekExecutor) Cancel(ctx context.Context, requestContext *execution.RequestContext, queue *event.Queue) error {
	return fmt.Errorf("cancel operation not supported")
}

func (e *DeepSeekExecutor) processRequest(ctx context.Context, messageText string, taskUpdater *updater.TaskUpdater) error {
	if e.client == nil {
		log.Printf("ERROR: DeepSeekExecutor client is nil!")
		return fmt.Errorf("DeepSeekExecutor client is nil")
	}

	log.Printf("Processing request with message: %s", messageText)

	messages := []deepseek.ChatCompletionMessage{
		{Role: deepseek.ChatMessageRoleSystem, Content: e.systemPrompt},
		{Role: deepseek.ChatMessageRoleUser, Content: messageText},
	}

	var tools []deepseek.Tool
	for _, function := range e.tools {
		tools = append(tools, deepseek.Tool{
			Type:     "function",
			Function: function.OpenAIFunctionDefinition(),
		})
	}

	log.Printf("Created %d tools for the request", len(tools))

	maxIterations := 10
	iteration := 0

	for iteration < maxIterations {
		iteration += 1
		log.Printf("Making API call iteration %d/%d", iteration, maxIterations)

		response, err := e.client.CreateChatCompletion(ctx, &deepseek.ChatCompletionRequest{
			Model:    e.model,
			Messages: messages,
			Tools:    tools,
		})

		if err != nil {
			log.Printf("Error in API call: %v", err)
			return err
		}

		log.Printf("API call successful, got %d choices", len(response.Choices))

		message := response.Choices[0].Message
		messages = append(messages, deepseek.ChatCompletionMessage{
			Role:      deepseek.ChatMessageRoleAssistant,
			Content:   message.Content,
			ToolCalls: message.ToolCalls,
		})

		if len(message.ToolCalls) == 0 {
			if message.Content != "" {
				log.Println("Assistant response:", message.Content)
				agentParts := []types.Part{
					&types.TextPart{Kind: "text", Text: message.Content},
				}
				taskUpdater.AddArtifact(agentParts)
				taskUpdater.Complete(updater.WithFinal(true))
			}
			break
		}

		log.Printf("Processing %d tool calls", len(message.ToolCalls))

		for _, tool := range message.ToolCalls {
			name := tool.Function.Name
			args := tool.Function.Arguments
			if function, ok := e.tools[name]; ok {
				var arg map[string]interface{}
				err = json.Unmarshal([]byte(args), &arg)
				if err != nil {
					log.Printf("Error parsing function arguments: %v", err)
					continue
				}

				res := function.Call(arg)
				// Serialize the result to JSON string
				resultJSON, err := json.Marshal(res)
				log.Printf("Result JSON: %s", string(resultJSON))
				if err != nil {
					resultJSON = []byte(fmt.Sprintf(`{"error": "Failed to serialize result: %v"}`, err))
				}

				messages = append(messages, deepseek.ChatCompletionMessage{
					Role:       deepseek.ChatMessageRoleTool,
					ToolCallID: tool.ID,
					Content:    string(resultJSON),
				})
			} else {
				messages = append(messages, deepseek.ChatCompletionMessage{
					Role:       deepseek.ChatMessageRoleTool,
					ToolCallID: tool.ID,
					Content:    fmt.Sprintf(`{"error": "Method %s not found on tool instance"}`, name),
				})
			}
		}

		agentMessage := taskUpdater.NewAgentMessage([]types.Part{
			&types.TextPart{Kind: "text", Text: "Processing tool calls..."},
		})
		taskUpdater.UpdateStatus(types.WORKING, updater.WithMessage(&agentMessage))
	}

	if iteration >= maxIterations {
		parts := []types.Part{&types.TextPart{Kind: "text", Text: "Sorry, the request has exceeded the maximum number of iterations."}}
		taskUpdater.Complete(updater.WithFinal(true), updater.WithMessage(&types.Message{
			Parts: parts,
		}))
	}
	return nil
}
