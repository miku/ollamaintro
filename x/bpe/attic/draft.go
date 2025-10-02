// Example: Byte-Pair Encoding.
package main

import (
	"io"
	"log"
	"os"
	"unicode/utf8"
)

type Pair struct {
	A, B int
}

func getStats(tokens []int) map[Pair]int {
	counts := make(map[Pair]int)
	for i := 0; i < len(tokens)-1; i++ {
		counts[Pair{tokens[i], tokens[i+1]}] += 1
	}
	return counts
}

func findMostCommonPair(tokens []int) Pair {
	counts := getStats(tokens)
	mostCommon := Pair{}
	maxCount := -1
	for k, v := range counts {
		if v > maxCount {
			maxCount = v
			mostCommon = k
		}
	}
	return mostCommon
}

func merge(ids []int, pair Pair, idx int) []int {
	newIds := make([]int, 0)
	i := 0
	for i < len(ids) {
		if i < len(ids)-1 && ids[i] == pair.A && ids[i+1] == pair.B {
			newIds = append(newIds, idx)
			i += 2
		} else {
			newIds = append(newIds, ids[i])
			i += 1
		}
	}
	return newIds
}

func main() {
	var (
		bs, _  = io.ReadAll(os.Stdin)
		tokens = make([]int, len(bs))
	)
	for i, b := range bs {
		tokens[i] = int(b)
	}
	log.Printf("text length: %d, byte length: %d \n",
		utf8.RuneCountInString(string(bs)), len(tokens))
	log.Printf("num tokens: %d", len(tokens))

	mostCommon := findMostCommonPair(tokens)
	log.Println(mostCommon)
	ids := merge(tokens, mostCommon, 256)
	log.Println(ids)
	log.Printf("num ids: %d", len(ids))
}
