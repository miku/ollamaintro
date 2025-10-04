package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand/v2"
	"os"

	"github.com/ollama/ollama/api"
	"golang.org/x/exp/constraints"
)

var (
	ollamaHost         = flag.String("s", os.Getenv("OLLAMA_HOST"), "ollama host")
	modelName          = flag.String("m", "nomic-embed-text", "model name")
	filename           = flag.String("f", "data/pg1661.json", "data file to read text snippets from")
	store              = flag.String("w", "embed.json", "file to save embeddings to")
	limit              = flag.Int("l", 10, "limit requests")
	doCreateEmbeddings = flag.Bool("e", false, "create embeddings for data")
	doRandomTriplet    = flag.Bool("3", false, "random three")

	// (1) take three snippets and compare / which one is closer to the other?
	// (2) enter some text and find similar snippets
	// (3) cluster snippets
)

// Chunk is a small piece of text.
type Chunk struct {
	ID        int64     `json:"id"`
	Text      string    `json:"text"`
	Embedding []float32 `json:"embedding,omitempty"`
}

type Set struct {
	Chunks []Chunk
}

func main() {
	flag.Parse()

	switch {
	case *doCreateEmbeddings:
		f, err := os.Open(*filename)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		br := bufio.NewReader(f)
		set := Set{}
		for {
			b, err := br.ReadBytes('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			var chunk Chunk
			if err := json.Unmarshal(b, &chunk); err != nil {
				log.Fatal(err)
			}
			set.Chunks = append(set.Chunks, chunk)
		}
		if _, err := os.Stat(*store); err == nil {
			log.Printf("output %s exists, nothing to do", *store)
			return
		}
		client, err := api.ClientFromEnvironment()
		if err != nil {
			log.Fatal(err)
		}
		version, err := client.Version(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("ollama client %v", version)
		log.Printf("generating embeddings...")
		// TODO: batching
		for i := range len(set.Chunks) {
			if i%50 == 0 {
				log.Printf("@%v", i)
			}
			if *limit > 0 && i == *limit {
				break
			}
			req := &api.EmbedRequest{
				Model: *modelName,
				Input: set.Chunks[i].Text,
			}
			resp, err := client.Embed(context.TODO(), req)
			if err != nil {
				log.Fatal(err)
			}
			if len(resp.Embeddings) < 1 {
				log.Println("warn: empty reply")
				continue
			}
			set.Chunks[i].Embedding = resp.Embeddings[0]
		}
		// save updated chunks
		var buf bytes.Buffer
		var enc = json.NewEncoder(&buf)
		for i, chunk := range set.Chunks {
			if *limit > 0 && i == *limit {
				break
			}
			err := enc.Encode(chunk)
			if err != nil {
				log.Fatal(err)
			}
		}
		if err := os.WriteFile(*store, buf.Bytes(), 0644); err != nil {
			fmt.Println(buf.String())
		}
	case *doRandomTriplet:
		if _, err := os.Stat(*store); os.IsNotExist(err) {
			log.Fatal("create embedding first")
		}
		set, err := loadChunksFromFile(*store)
		if err != nil {
			log.Fatal(err)
		}
		if len(set.Chunks) < 3 {
			log.Printf("not enough chunks, need 3, found %v", len(set.Chunks))
			return
		}
		perm := rand.Perm(len(set.Chunks))
		var selected = make([]Chunk, 3)
		for i := range 3 {
			idx := perm[i]
			selected[i] = set.Chunks[idx]
		}
		for _, v := range Combinations(selected, 2) {
			sim := CosineSimilarity(v[0].Embedding, v[1].Embedding)
			log.Println(v[0].ID, v[1].ID, sim)
		}
		fmt.Println()
		for i := range selected {
			fmt.Printf("%d: %v\n\n", selected[i].ID, selected[i].Text)
		}
	}
}

func loadChunksFromFile(filename string) (*Set, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	br := bufio.NewReader(f)
	set := &Set{}
	for {
		b, err := br.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		var chunk Chunk
		if err := json.Unmarshal(b, &chunk); err != nil {
			log.Fatal(err)
		}
		set.Chunks = append(set.Chunks, chunk)
	}
	return set, nil
}

// CosineSimilarity computes the cosine similarity between two vectors.
// Works with any numeric type but returns float64 for precision.
func CosineSimilarity[T constraints.Float | constraints.Integer](a, b []T) float64 {
	if len(a) != len(b) {
		return 0
	}
	var dotProduct, normA, normB float64
	for i := 0; i < len(a); i++ {
		aVal := float64(a[i])
		bVal := float64(b[i])
		dotProduct += aVal * bVal
		normA += aVal * aVal
		normB += bVal * bVal
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// Combinations returns all r-length combinations of elements from the input slice.
// The combinations are generated in lexicographic order.
func Combinations[T any](items []T, r int) [][]T {
	if r > len(items) || r <= 0 {
		return [][]T{}
	}

	var result [][]T
	indices := make([]int, r)

	// Initialize indices to 0, 1, 2, ..., r-1
	for i := range indices {
		indices[i] = i
	}

	for {
		// Copy current combination
		temp := make([]T, r)
		for i, idx := range indices {
			temp[i] = items[idx]
		}
		result = append(result, temp)

		// Find the rightmost index that can be incremented
		i := r - 1
		for i >= 0 && indices[i] == len(items)-r+i {
			i--
		}

		// If no such index exists, we're done
		if i < 0 {
			break
		}

		// Increment this index and reset all indices to its right
		indices[i]++
		for j := i + 1; j < r; j++ {
			indices[j] = indices[j-1] + 1
		}
	}

	return result
}
