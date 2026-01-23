// 03_embedding.go
//
// Embedding example showing how text is converted to dense vectors.
// Embeddings can be used for semantic search, clustering, and similarity.
//
// Run: go run 03_embedding.go

package main

import (
	"context"
	"fmt"
	"log"
	"math"

	"github.com/ollama/ollama/api"
)

func main() {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	// Two sentences to compare
	texts := []string{
		"The sky is blue because of Rayleigh scattering",
		"Blue light scatters more in the atmosphere",
		"I like pizza with extra cheese",
	}

	embeddings := make([][]float64, len(texts))

	for i, text := range texts {
		req := &api.EmbedRequest{
			Model: "embeddinggemma",
			Input: text,
		}

		resp, err := client.Embed(context.TODO(), req)
		if err != nil {
			log.Fatal(err)
		}

		embeddings[i] = resp.Embeddings[0]
		fmt.Printf("Text %d: \"%s\"\n", i+1, text)
		fmt.Printf("  Dimensions: %d\n", len(embeddings[i]))
		fmt.Printf("  First 5 values: %v\n\n", embeddings[i][:5])
	}

	// Compute cosine similarity between all pairs
	fmt.Println("Cosine Similarities:")
	for i := 0; i < len(texts); i++ {
		for j := i + 1; j < len(texts); j++ {
			sim := cosineSimilarity(embeddings[i], embeddings[j])
			fmt.Printf("  [%d] vs [%d]: %.4f\n", i+1, j+1, sim)
		}
	}
}

func cosineSimilarity(a, b []float64) float64 {
	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}
