package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ollama/ollama/api"
)

func main() {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	req := &api.GenerateRequest{
		Model:  "gemma3",
		Prompt: "Why is the sky blue?",
		Stream: new(bool),
	}

	ctx := context.Background()

	var fullResponse string
	err = client.Generate(ctx, req, func(resp api.GenerateResponse) error {
		fullResponse = resp.Response
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fullResponse)
}
