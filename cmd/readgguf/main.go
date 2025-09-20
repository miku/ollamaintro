package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ollama/ollama/fs/gguf"
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		log.Fatal("path to model required")
	}
	modelPath := flag.Arg(0)
	f, err := gguf.Open(modelPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	fmt.Println(f.Magic)
	fmt.Println(f.Version)
	for k, v := range f.KeyValues() {
		fmt.Printf("  %v => %v\n", k, v)
	}
	fmt.Printf("%v\n", f.NumTensors())
	for k, v := range f.TensorInfos() {
		fmt.Printf("  %v => %v\n", k, v)
	}
}
