package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ollama/ollama/api"
)

var (
	ollamaHost = flag.String("s", os.Getenv("OLLAMA_HOST"), "ollama host")
	modelName  = flag.String("m", "nomic-embed-text", "model name")
	filename   = flag.String("f", "data/pg1661.json", "data file to read text snippets from")
	store      = flag.String("w", "embed.json", "file to save embeddings to")
	limit      = flag.Int("l", 100, "limit requests")

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
		if i == *limit {
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
	for i, chunk := range set.Chunks {
		if i == *limit {
			break
		}
		b, err := json.Marshal(chunk)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(b))
	}

}
