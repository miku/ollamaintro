package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ollama/ollama/api"
)

func main() {
	// Check if image path was provided
	if len(os.Args) < 2 {
		log.Fatal("usage: go run main.go <path-to-image>")
	}

	imagePath := os.Args[1]

	// Read the image file
	imgData, err := os.ReadFile(imagePath)
	if err != nil {
		log.Fatalf("failed to read image: %v", err)
	}

	// Create Ollama client
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	// Create the request for text recognition
	req := &api.GenerateRequest{
		Model:  "qwen2.5vl", // or "llava:13b" for better accuracy
		Prompt: "Extract and transcribe all text visible in this image. List the text exactly as it appears.",
		Images: []api.ImageData{imgData},
	}

	ctx := context.Background()

	fmt.Println("Recognizing text from image...")
	fmt.Println("---")

	// Handle the streaming response
	respFunc := func(resp api.GenerateResponse) error {
		fmt.Print(resp.Response)
		return nil
	}

	err = client.Generate(ctx, req, respFunc)
	if err != nil {
		log.Fatalf("failed to generate: %v", err)
	}

	fmt.Println()
}
