// 02_multimodal.go
//
// Multimodal example showing how vision models process images.
// The image bytes are passed in the Images field and processed
// by the vision encoder into embeddings that join the text tokens.
//
// Run: go run 02_multimodal.go image.png

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ollama/ollama/api"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run 02_multimodal.go <image_path>")
	}

	imagePath := os.Args[1]

	// Read image file
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		log.Fatalf("Failed to read image: %v", err)
	}
	fmt.Printf("Loaded image: %s (%d bytes)\n", imagePath, len(imageData))

	// Create client
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	// Build multimodal request
	// The Images field contains raw image bytes
	// Internally, Ollama:
	// 1. Passes image through vision encoder to get embeddings
	// 2. Replaces [img-0] placeholder with those embeddings
	// 3. Processes combined image + text through the transformer
	req := &api.GenerateRequest{
		Model:  "qwen2.5vl", // Vision-capable model
		Prompt: "Describe this image in detail. What do you see?",
		Images: []api.ImageData{imageData},
	}

	fmt.Println("Prompt:", req.Prompt)
	fmt.Println("---")

	ctx := context.Background()
	err = client.Generate(ctx, req, func(resp api.GenerateResponse) error {
		fmt.Print(resp.Response)
		if resp.Done {
			fmt.Printf("\n---\n")
			fmt.Printf("Generated tokens: %d\n", resp.EvalCount)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
