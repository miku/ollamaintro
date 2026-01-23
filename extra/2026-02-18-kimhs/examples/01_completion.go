// 01_completion.go
//
// Basic text completion example showing what happens when you type a prompt.
// This demonstrates the simplest Ollama API interaction.
//
// Run: go run 01_completion.go

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ollama/ollama/api"
)

func main() {
	// Create client from environment (uses OLLAMA_HOST or defaults to localhost:11434)
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	// Build the request
	// This maps directly to the GenerateRequest struct in the slides
	req := &api.GenerateRequest{
		Model:  "gemma3:270m", // Use a small model for quick demo
		Prompt: "Why is the sky blue? Answer in one sentence.",
	}

	fmt.Println("Prompt:", req.Prompt)
	fmt.Println("---")

	// Generate with streaming callback
	// Each token is passed to this function as it's generated
	ctx := context.Background()
	err = client.Generate(ctx, req, func(resp api.GenerateResponse) error {
		// Print each token as it arrives
		fmt.Print(resp.Response)

		// When done, print metadata
		if resp.Done {
			fmt.Printf("\n---\n")
			fmt.Printf("Total duration: %d ms\n", resp.TotalDuration/1_000_000)
			fmt.Printf("Prompt tokens: %d\n", resp.PromptEvalCount)
			fmt.Printf("Generated tokens: %d\n", resp.EvalCount)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
