package main

import (
	"flag"
	"log"

	"github.com/ollama/ollama/fs/gguf"
)

var (
	filename = flag.String("f", "llama3.2.gguf", "gguf filename")
)

func main() {
	f, err := gguf.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	for k, v := range f.KeyValues() {
		log.Println(k, v)
	}
	for k, v := range f.TensorInfos() {
		log.Println(k, v)
	}

	ti := f.TensorInfo("token_embd.weight")
	log.Println(ti)
	ti, r, err := f.TensorReader("token_embd.weight")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(ti)
	log.Println(r)
}
