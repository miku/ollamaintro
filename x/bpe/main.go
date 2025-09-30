// Example: Byte-Pair Encoding.
package main

import (
	"fmt"
	"io"
	"os"
)

type Pair struct {
	A, B int
}

func main() {
	var (
		text, _ = io.ReadAll(os.Stdin)
		bs      = []byte(text)
		tokens  = make([]int, len(bs))
	)
	for i, b := range bs {
		tokens[i] = int(b)
	}
	fmt.Println(tokens)
}
