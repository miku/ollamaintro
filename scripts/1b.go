package main

import (
	"bufio"
	"io"
	"os"
)

func main() {
	bw := bufio.NewWriter(os.Stdout)
	defer bw.Flush()
	for _ = range 1_000_000_000 {
		_, _ = io.WriteString(bw, "word\n")
	}
}
