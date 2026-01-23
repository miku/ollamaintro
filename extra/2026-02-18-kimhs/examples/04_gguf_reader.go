// 04_gguf_reader.go
//
// GGUF file reader demonstrating the model file format.
// This shows how to read metadata and tensor information.
//
// Run: go run 04_gguf_reader.go /path/to/model.gguf

package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/ollama/ollama/fs/gguf"
)

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatal("Usage: go run 04_gguf_reader.go <model.gguf>")
	}

	modelPath := flag.Arg(0)

	// Open the GGUF file
	f, err := gguf.Open(modelPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Print header information
	fmt.Println("=== GGUF File Header ===")
	fmt.Printf("Magic:    %s\n", f.Magic)
	fmt.Printf("Version:  %d\n", f.Version)
	fmt.Printf("Tensors:  %d\n", f.NumTensors())
	fmt.Println()

	// Print selected metadata
	fmt.Println("=== Key Metadata ===")
	for _, kv := range f.KeyValues() {
		// Filter to interesting keys
		if containsAny(kv.Key, []string{
			"general.architecture",
			"general.name",
			"context_length",
			"embedding_length",
			"block_count",
			"head_count",
			"vocab_size",
		}) {
			fmt.Printf("%-40s = %v\n", kv.Key, kv.Value)
		}
	}
	fmt.Println()

	// Print first few tensors
	fmt.Println("=== Tensor Information (first 10) ===")
	count := 0
	for _, t := range f.TensorInfos() {
		if count >= 10 {
			fmt.Println("... (more tensors)")
			break
		}
		fmt.Printf("%-40s  shape=%-20v  type=%v\n",
			t.Name, formatShape(t.Shape), t.Type)
		count++
	}
}

func containsAny(s string, substrs []string) bool {
	for _, sub := range substrs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

func formatShape(dims []uint64) string {
	parts := make([]string, len(dims))
	for i, d := range dims {
		parts[i] = fmt.Sprintf("%d", d)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}
