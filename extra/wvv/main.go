package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

// Word2Vec implements Skip-gram with Negative Sampling
type Word2Vec struct {
	EmbeddingDim    int
	WindowSize      int
	MinCount        int
	NegativeSamples int
	LearningRate    float64
	Epochs          int

	vocab       map[string]int
	indexToWord map[int]string
	wordFreq    []float64
	wInput      [][]float64 // Input embeddings
	wOutput     [][]float64 // Output embeddings
}

// NewWord2Vec creates a new Word2Vec model
func NewWord2Vec(embeddingDim, windowSize, minCount, negativeSamples int,
	learningRate float64, epochs int) *Word2Vec {
	return &Word2Vec{
		EmbeddingDim:    embeddingDim,
		WindowSize:      windowSize,
		MinCount:        minCount,
		NegativeSamples: negativeSamples,
		LearningRate:    learningRate,
		Epochs:          epochs,
		vocab:           make(map[string]int),
		indexToWord:     make(map[int]string),
	}
}

// buildVocab builds vocabulary from sentences
func (w *Word2Vec) buildVocab(sentences [][]string) int {
	// Count word frequencies
	wordCounts := make(map[string]int)
	for _, sentence := range sentences {
		for _, word := range sentence {
			wordCounts[word]++
		}
	}

	// Filter by minimum count
	vocabWords := make([]string, 0)
	for word, count := range wordCounts {
		if count >= w.MinCount {
			vocabWords = append(vocabWords, word)
		}
	}

	// Create mappings
	for idx, word := range vocabWords {
		w.vocab[word] = idx
		w.indexToWord[idx] = word
	}

	// Store word frequencies for negative sampling
	w.wordFreq = make([]float64, len(vocabWords))
	sum := 0.0
	for idx, word := range vocabWords {
		freq := math.Pow(float64(wordCounts[word]), 0.75)
		w.wordFreq[idx] = freq
		sum += freq
	}

	// Normalize to probabilities
	for i := range w.wordFreq {
		w.wordFreq[i] /= sum
	}

	return len(w.vocab)
}

// getNegativeSamples samples negative examples
func (w *Word2Vec) getNegativeSamples(targetIdx, nSamples int) []int {
	samples := make([]int, 0, nSamples)
	for len(samples) < nSamples {
		sample := w.sampleFromDist(w.wordFreq)
		if sample != targetIdx {
			samples = append(samples, sample)
		}
	}
	return samples
}

// sampleFromDist samples an index from a probability distribution
func (w *Word2Vec) sampleFromDist(probs []float64) int {
	r := rand.Float64()
	cumSum := 0.0
	for i, p := range probs {
		cumSum += p
		if r <= cumSum {
			return i
		}
	}
	return len(probs) - 1
}

// sigmoid activation function
func sigmoid(x float64) float64 {
	if x < -500 {
		x = -500
	} else if x > 500 {
		x = 500
	}
	return 1.0 / (1.0 + math.Exp(-x))
}

// dot computes dot product of two vectors
func dot(a, b []float64) float64 {
	sum := 0.0
	for i := range a {
		sum += a[i] * b[i]
	}
	return sum
}

// trainPair trains on a single (center, context) pair with negative sampling
func (w *Word2Vec) trainPair(centerIdx, contextIdx int) {
	// Get embeddings
	h := w.wInput[centerIdx]
	u := w.wOutput[contextIdx]

	// Positive sample
	score := dot(h, u)
	pred := sigmoid(score)
	grad := (pred - 1) * w.LearningRate

	// Update output embedding for context word
	for i := range u {
		w.wOutput[contextIdx][i] -= grad * h[i]
	}

	// Accumulate gradient for center word
	gradH := make([]float64, w.EmbeddingDim)
	for i := range gradH {
		gradH[i] = grad * u[i]
	}

	// Negative samples
	negativeIndices := w.getNegativeSamples(contextIdx, w.NegativeSamples)
	for _, negIdx := range negativeIndices {
		uNeg := w.wOutput[negIdx]
		scoreNeg := dot(h, uNeg)
		predNeg := sigmoid(scoreNeg)
		gradNeg := predNeg * w.LearningRate

		// Update negative sample embedding
		for i := range uNeg {
			w.wOutput[negIdx][i] -= gradNeg * h[i]
		}

		// Accumulate gradient
		for i := range gradH {
			gradH[i] += gradNeg * uNeg[i]
		}
	}

	// Update center word embedding
	for i := range h {
		w.wInput[centerIdx][i] -= gradH[i]
	}
}

// Train trains the Word2Vec model
func (w *Word2Vec) Train(sentences [][]string) {
	// Build vocabulary
	vocabSize := w.buildVocab(sentences)
	fmt.Printf("Vocabulary size: %d\n", vocabSize)

	// Initialize embeddings
	w.wInput = make([][]float64, vocabSize)
	w.wOutput = make([][]float64, vocabSize)
	for i := 0; i < vocabSize; i++ {
		w.wInput[i] = make([]float64, w.EmbeddingDim)
		w.wOutput[i] = make([]float64, w.EmbeddingDim)
		for j := 0; j < w.EmbeddingDim; j++ {
			w.wInput[i][j] = (rand.Float64() - 0.5) / float64(w.EmbeddingDim)
			w.wOutput[i][j] = (rand.Float64() - 0.5) / float64(w.EmbeddingDim)
		}
	}

	// Training loop
	for epoch := 0; epoch < w.Epochs; epoch++ {
		totalPairs := 0

		for _, sentence := range sentences {
			// Convert words to indices
			wordIndices := make([]int, 0, len(sentence))
			for _, word := range sentence {
				if idx, ok := w.vocab[word]; ok {
					wordIndices = append(wordIndices, idx)
				}
			}

			// Generate training pairs
			for i, centerIdx := range wordIndices {
				start := max(0, i-w.WindowSize)
				end := min(len(wordIndices), i+w.WindowSize+1)

				for j := start; j < end; j++ {
					if i != j {
						contextIdx := wordIndices[j]
						w.trainPair(centerIdx, contextIdx)
						totalPairs++
					}
				}
			}
		}

		fmt.Printf("Epoch %d/%d - Trained on %d pairs\n", epoch+1, w.Epochs, totalPairs)
	}
}

// GetVector returns the embedding vector for a word
func (w *Word2Vec) GetVector(word string) ([]float64, error) {
	idx, ok := w.vocab[word]
	if !ok {
		return nil, fmt.Errorf("word '%s' not in vocabulary", word)
	}
	return w.wInput[idx], nil
}

// WordSimilarity represents a word and its similarity score
type WordSimilarity struct {
	Word       string
	Similarity float64
}

// MostSimilar finds most similar words using cosine similarity
func (w *Word2Vec) MostSimilar(word string, topN int) ([]WordSimilarity, error) {
	idx, ok := w.vocab[word]
	if !ok {
		return nil, fmt.Errorf("word '%s' not in vocabulary", word)
	}

	wordVec := w.wInput[idx]
	wordNorm := norm(wordVec)

	// Compute similarities
	similarities := make([]WordSimilarity, 0, len(w.vocab))
	for otherWord, otherIdx := range w.vocab {
		if otherIdx == idx {
			continue // Skip the word itself
		}

		otherVec := w.wInput[otherIdx]
		cosineSim := dot(wordVec, otherVec) / (wordNorm * norm(otherVec))
		similarities = append(similarities, WordSimilarity{
			Word:       otherWord,
			Similarity: cosineSim,
		})
	}

	// Sort by similarity descending
	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].Similarity > similarities[j].Similarity
	})

	// Return top N
	if topN > len(similarities) {
		topN = len(similarities)
	}
	return similarities[:topN], nil
}

// norm computes the L2 norm of a vector
func norm(vec []float64) float64 {
	sum := 0.0
	for _, v := range vec {
		sum += v * v
	}
	return math.Sqrt(sum)
}

// Helper functions for Go < 1.21 compatibility
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	// Sample corpus
	sentences := [][]string{
		{"the", "quick", "brown", "fox", "jumps", "over", "the", "lazy", "dog"},
		{"the", "dog", "is", "lazy"},
		{"the", "cat", "is", "quick"},
		{"the", "fox", "is", "brown"},
		{"a", "quick", "brown", "dog", "jumps"},
		{"the", "lazy", "cat", "sleeps"},
		{"quick", "animals", "run", "fast"},
		{"brown", "animals", "are", "common"},
	}

	// Replicate sentences for more training data
	var allSentences [][]string
	for i := 0; i < 100; i++ {
		allSentences = append(allSentences, sentences...)
	}

	// Train model
	model := NewWord2Vec(50, 2, 1, 5, 0.025, 10)
	model.Train(allSentences)

	// Test similarity
	fmt.Println("\nMost similar to 'quick':")
	if similar, err := model.MostSimilar("quick", 5); err == nil {
		for _, ws := range similar {
			fmt.Printf("  %s: %.4f\n", ws.Word, ws.Similarity)
		}
	}

	fmt.Println("\nMost similar to 'dog':")
	if similar, err := model.MostSimilar("dog", 5); err == nil {
		for _, ws := range similar {
			fmt.Printf("  %s: %.4f\n", ws.Word, ws.Similarity)
		}
	}
}
