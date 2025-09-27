// e.g. go run main.go /usr/share/ollama/.ollama/models/blobs/sha256-dde5aa3fc5ffc17176b5e8bdc82f587b24b2678c6c66101bf7da77af9f7ccdff
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
