package main

import (
	"log"

	"github.com/ollama/ollama/runner/llamarunner"
)

func main() {
	s := llamarunner.Server{}
	prompt := "hello world"
	params := llamarunner.NewSequenceParams{}

	seq, err := s.NewSequence(prompt, nil, params)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(s)
	log.Println(seq)
}
