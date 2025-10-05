package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ollama/ollama/api"
)

func main() {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	b, err := os.ReadFile("image.png")
	if err != nil {
		log.Fatal(err)
	}

	req := &api.GenerateRequest{
		Model:  "qwen2.5vl",
		Prompt: "Describe this image in detail",
		Images: []api.ImageData{b},
	}

	err = client.Generate(context.TODO(), req, func(resp api.GenerateResponse) error {
		fmt.Print(resp.Response)
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
}
