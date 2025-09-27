package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ollama/ollama/llama"
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatal("model path required")
	}
	path := flag.Arg(0)
	arch, err := llama.GetModelArch(path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(arch)
}
