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
