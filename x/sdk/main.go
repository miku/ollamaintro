package main

import (
	"context"
	"log"

	"github.com/ollama/ollama/api"
)

func main() {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.List(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	for _, m := range resp.Models {
		log.Printf("%v %v %v", m.Digest, m.Name, m.Details.ParameterSize)
	}
}
