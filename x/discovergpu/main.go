package main

import (
	"log"

	"github.com/ollama/ollama/discover"
)

func main() {

	gi := discover.GetGPUInfo()
	log.Println(gi)
	// Machine w/ dedicated GPU
	// 2025/09/27 22:18:29 INFO looking for compatible GPUs
	// 2025/09/27 22:18:29 [{{20989804544 18706530304 0} cuda v13 479199232 [/lib/ollama/cuda_v13] [] false GPU-ddce76c8-78ec-75d8-4d6c-007c17a668a8 0 NVIDIA RTX 4000 SFF Ada Generation 8.9 13 0}]

	// Machine w/o dedicated GPU
	// 2025/09/27 22:18:57 INFO looking for compatible GPUs
	// 2025/09/27 22:18:57 INFO no compatible GPUs were discovered
	// 2025/09/27 22:18:57 [{{33291141120 19012984832 1024454656} cpu  0 [] [] false 0 0   0 0}]

}
