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
	"os"

	"github.com/ollama/ollama/api"
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
	}
}
