package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ollama/ollama/api"
)

func main() {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	messages := []api.Message{
		{
			Role:    "user",
			Content: "What is the weather in Florence?",
		},
	}

	tools := api.Tools{
		{Type: "function", Function: api.ToolFunction{
			Name:        "get_current_weather",
			Description: "get the current weather for a city",
			Parameters: api.ToolFunctionParameters{
				Type:     "object",
				Required: []string{"location"},
				Properties: map[string]api.ToolProperty{
					"location": {
						Type:        api.PropertyType{"string"},
						Description: "The city and state, e.g. San Francisco, CA",
					},
				},
			},
		}},
	}

	// tools := []api.Tool{
	// 	{
	// 		Type: "function",
	// 		Function: api.ToolFunction{
	// 			Name:        "get_current_weather",
	// 			Description: "Get the current weather for a city",
	// 			Parameters: api.ToolFunctionParameters{
	// 				Type: "object",
	// 				Properties: map[string]map[string]string{
	// 					"city": {
	// 						"type":        "string",
	// 						"description": "The name of the city",
	// 					},
	// 				},
	// 				Required: []string{"city"},
	// 			},
	// 		},
	// 	},
	// }

	req := &api.ChatRequest{
		Model:    "llama3.2",
		Messages: messages,
		Tools:    tools,
	}

	ctx := context.Background()

	err = client.Chat(ctx, req, func(resp api.ChatResponse) error {
		b, err := json.Marshal(resp.Message)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(b))
		// fmt.Println(resp.Message)
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
