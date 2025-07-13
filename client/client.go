package main

import (
	"log"
	"net/http"

	"github.com/yeeaiclub/a2a-go/sdk/client"
	"github.com/yeeaiclub/a2a-go/sdk/types"
)

func main() {
	httpClient := http.Client{}
	newClient := client.NewClient(&httpClient, "http://localhost:8080/api")

	message, err := newClient.SendMessage(types.MessageSendParam{
		Message: &types.Message{
			ContextID: "context-1",
			Role:      types.User,
			Parts: []types.Part{
				&types.TextPart{Text: "Show recent commits for repository 'facebook/react", Kind: "text"},
			},
		},
	})
	if err != nil {
		log.Println(err)
		return
	}

	if message.Error != nil {
		log.Printf("Server error: %s (code: %d)", message.Error.Message, message.Error.Code)
		return
	}

	task, err := types.MapTo[types.Task](message.Result)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("Task ID: %s", task.Id)
	log.Printf("Task State: %s", task.Status.State)
	log.Printf("History length: %d", len(task.History))

	for _, art := range task.Artifacts {
		for j, part := range art.Parts {
			if textPart, ok := part.(*types.TextPart); ok {
				log.Printf("  Part[%d]: %s", j, textPart.Text)
			}
		}
	}
}
