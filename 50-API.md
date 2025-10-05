# API and SDK


## Go SDK

```go
package main

import (
	"context"
	"log"

	"github.com/ollama/ollama/api"
)

func main() {
	// api.NewClient or api.ClientFromEnvironment
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.List(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	for _, m := range resp.Models {
		log.Printf("%v %v %v", m.Digest, m.Name, m.Details.ParameterSize)
	}
}
```

## Text Completion

```go
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	req := &api.GenerateRequest{
		Model:  "gemma3",
		Prompt: "Why is the sky blue?",
	}

	ctx := context.Background()

	err = client.Generate(ctx, req, func(resp api.GenerateResponse) error {
		fmt.Print(resp.Response)
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
```

## Embedding model

```go

func main() {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	req := &api.EmbedRequest{
		Model: "embeddinggemma",
		Input: "The sky is blue because of Rayleigh scattering",
	}

	resp, err := client.Embed(context.TODO(), req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("dimensions: %d\n", len(resp.Embeddings[0]))
	fmt.Printf("%v ...\n", resp.Embeddings[0][:10])
}

```

## Image processing

```go
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
		Model:  "moondream",
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
``` 