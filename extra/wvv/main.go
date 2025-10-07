package main

import (
	"bufio"
	"encoding/gob"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strings"
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

	// Clip gradient to prevent exploding gradients
	grad = clipValue(grad, -10.0, 10.0)

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

		// Clip gradient
		gradNeg = clipValue(gradNeg, -10.0, 10.0)

		// Update negative sample embedding
		for i := range uNeg {
			w.wOutput[negIdx][i] -= gradNeg * h[i]
		}

		// Accumulate gradient
		for i := range gradH {
			gradH[i] += gradNeg * uNeg[i]
		}
	}

	// Clip accumulated gradient
	for i := range gradH {
		gradH[i] = clipValue(gradH[i], -10.0, 10.0)
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

	// Check for zero norm
	if wordNorm < 1e-10 {
		return nil, fmt.Errorf("word '%s' has zero or near-zero embedding", word)
	}

	// Compute similarities
	similarities := make([]WordSimilarity, 0, len(w.vocab))
	for otherWord, otherIdx := range w.vocab {
		if otherIdx == idx {
			continue // Skip the word itself
		}

		otherVec := w.wInput[otherIdx]
		otherNorm := norm(otherVec)

		// Skip vectors with zero or near-zero norm
		if otherNorm < 1e-10 {
			continue
		}

		cosineSim := dot(wordVec, otherVec) / (wordNorm * otherNorm)

		// Skip NaN or Inf values
		if math.IsNaN(cosineSim) || math.IsInf(cosineSim, 0) {
			continue
		}

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

// Save saves the model to a file
func (w *Word2Vec) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)

	// Encode all model data
	if err := encoder.Encode(w.EmbeddingDim); err != nil {
		return fmt.Errorf("failed to encode EmbeddingDim: %w", err)
	}
	if err := encoder.Encode(w.WindowSize); err != nil {
		return fmt.Errorf("failed to encode WindowSize: %w", err)
	}
	if err := encoder.Encode(w.MinCount); err != nil {
		return fmt.Errorf("failed to encode MinCount: %w", err)
	}
	if err := encoder.Encode(w.NegativeSamples); err != nil {
		return fmt.Errorf("failed to encode NegativeSamples: %w", err)
	}
	if err := encoder.Encode(w.LearningRate); err != nil {
		return fmt.Errorf("failed to encode LearningRate: %w", err)
	}
	if err := encoder.Encode(w.Epochs); err != nil {
		return fmt.Errorf("failed to encode Epochs: %w", err)
	}
	if err := encoder.Encode(w.vocab); err != nil {
		return fmt.Errorf("failed to encode vocab: %w", err)
	}
	if err := encoder.Encode(w.indexToWord); err != nil {
		return fmt.Errorf("failed to encode indexToWord: %w", err)
	}
	if err := encoder.Encode(w.wordFreq); err != nil {
		return fmt.Errorf("failed to encode wordFreq: %w", err)
	}
	if err := encoder.Encode(w.wInput); err != nil {
		return fmt.Errorf("failed to encode wInput: %w", err)
	}
	if err := encoder.Encode(w.wOutput); err != nil {
		return fmt.Errorf("failed to encode wOutput: %w", err)
	}

	return nil
}

// Load loads the model from a file
func (w *Word2Vec) Load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)

	// Decode all model data
	if err := decoder.Decode(&w.EmbeddingDim); err != nil {
		return fmt.Errorf("failed to decode EmbeddingDim: %w", err)
	}
	if err := decoder.Decode(&w.WindowSize); err != nil {
		return fmt.Errorf("failed to decode WindowSize: %w", err)
	}
	if err := decoder.Decode(&w.MinCount); err != nil {
		return fmt.Errorf("failed to decode MinCount: %w", err)
	}
	if err := decoder.Decode(&w.NegativeSamples); err != nil {
		return fmt.Errorf("failed to decode NegativeSamples: %w", err)
	}
	if err := decoder.Decode(&w.LearningRate); err != nil {
		return fmt.Errorf("failed to decode LearningRate: %w", err)
	}
	if err := decoder.Decode(&w.Epochs); err != nil {
		return fmt.Errorf("failed to decode Epochs: %w", err)
	}
	if err := decoder.Decode(&w.vocab); err != nil {
		return fmt.Errorf("failed to decode vocab: %w", err)
	}
	if err := decoder.Decode(&w.indexToWord); err != nil {
		return fmt.Errorf("failed to decode indexToWord: %w", err)
	}
	if err := decoder.Decode(&w.wordFreq); err != nil {
		return fmt.Errorf("failed to decode wordFreq: %w", err)
	}
	if err := decoder.Decode(&w.wInput); err != nil {
		return fmt.Errorf("failed to decode wInput: %w", err)
	}
	if err := decoder.Decode(&w.wOutput); err != nil {
		return fmt.Errorf("failed to decode wOutput: %w", err)
	}

	return nil
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

// clipValue clips a value to be within [minVal, maxVal]
func clipValue(val, minVal, maxVal float64) float64 {
	if val < minVal {
		return minVal
	}
	if val > maxVal {
		return maxVal
	}
	return val
}

// filterAlpha keeps only lowercase alphabetic characters
func filterAlpha(s string) string {
	var result strings.Builder
	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func main() {
	// Command line flags
	useStdin := flag.Bool("stdin", false, "Read sentences from stdin (one sentence per line)")
	saveFile := flag.String("save", "", "Save trained model to file")
	loadFile := flag.String("load", "", "Load model from file (skips training)")
	embeddingDim := flag.Int("dim", 50, "Embedding dimension")
	windowSize := flag.Int("window", 2, "Context window size")
	minCount := flag.Int("min-count", 1, "Minimum word frequency")
	negativeSamples := flag.Int("negative", 5, "Number of negative samples")
	learningRate := flag.Float64("lr", 0.025, "Learning rate")
	epochs := flag.Int("epochs", 10, "Number of training epochs")
	flag.Parse()

	model := NewWord2Vec(*embeddingDim, *windowSize, *minCount, *negativeSamples, *learningRate, *epochs)

	// Load existing model if specified
	if *loadFile != "" {
		fmt.Fprintf(os.Stderr, "Loading model from %s...\n", *loadFile)
		if err := model.Load(*loadFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading model: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Model loaded. Vocabulary size: %d\n", len(model.vocab))
	} else {
		// Train a new model
		var sentences [][]string

		if *useStdin {
			// Read from stdin
			fmt.Fprintln(os.Stderr, "Reading sentences from stdin...")
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line == "" {
					continue
				}
				// Split line into words and convert to lowercase
				fields := strings.Fields(strings.ToLower(line))
				words := make([]string, 0, len(fields))
				for _, field := range fields {
					// Keep only alphabetic characters
					word := filterAlpha(field)
					if word != "" {
						words = append(words, word)
					}
				}
				if len(words) > 0 {
					sentences = append(sentences, words)
				}
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stderr, "Read %d sentences\n", len(sentences))
		} else {
			// Use demo corpus
			demoSentences := [][]string{
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
			for i := 0; i < 100; i++ {
				sentences = append(sentences, demoSentences...)
			}
		}

		// Train model
		model.Train(sentences)

		// Save model if specified
		if *saveFile != "" {
			fmt.Fprintf(os.Stderr, "Saving model to %s...\n", *saveFile)
			if err := model.Save(*saveFile); err != nil {
				fmt.Fprintf(os.Stderr, "Error saving model: %v\n", err)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stderr, "Model saved successfully.\n")
		}
	}

	// Test similarity
	fmt.Println("\nMost similar to 'quick':")
	if similar, err := model.MostSimilar("quick", 5); err == nil {
		for _, ws := range similar {
			fmt.Printf("  %s: %.4f\n", ws.Word, ws.Similarity)
		}
	} else {
		fmt.Printf("  Error: %v\n", err)
	}

	fmt.Println("\nMost similar to 'dog':")
	if similar, err := model.MostSimilar("dog", 5); err == nil {
		for _, ws := range similar {
			fmt.Printf("  %s: %.4f\n", ws.Word, ws.Similarity)
		}
	} else {
		fmt.Printf("  Error: %v\n", err)
	}

	fmt.Println("\nMost similar to 'sky':")
	if similar, err := model.MostSimilar("sky", 5); err == nil {
		for _, ws := range similar {
			fmt.Printf("  %s: %.4f\n", ws.Word, ws.Similarity)
		}
	} else {
		fmt.Printf("  Error: %v\n", err)
	}
}
