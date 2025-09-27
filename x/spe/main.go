package main

import (
	"container/heap"
	"fmt"
	"unicode/utf8"
)

// Token represents a vocabulary token
type Token int

const NullToken Token = -1

// Vocabulary interface for token lookups
type Vocabulary interface {
	TextToToken(text string) Token
	ByteToToken(b byte) Token
	GetTokenScore(token Token) float64
	IsValidToken(token Token) bool
}

// Symbol represents a text segment during tokenization
type Symbol struct {
	Text string // the actual text content
	N    int    // length in bytes (for compatibility checks)
	Prev int    // index of previous symbol (-1 if none)
	Next int    // index of next symbol (-1 if none)
}

// Bigram represents a potential merge between two adjacent symbols
type Bigram struct {
	Left  int     // index of left symbol
	Right int     // index of right symbol
	Score float64 // score from vocabulary
	Size  int     // total size of merged text in bytes
}

// BigramQueue implements a max-heap for bigrams
type BigramQueue []*Bigram

func (pq BigramQueue) Len() int           { return len(pq) }
func (pq BigramQueue) Less(i, j int) bool { return pq[i].Score > pq[j].Score }
func (pq BigramQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }

func (pq *BigramQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Bigram))
}

func (pq *BigramQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

// TokenizerSession represents a single tokenization session
type TokenizerSession struct {
	vocab     Vocabulary
	symbols   []Symbol
	workQueue BigramQueue
	revMerge  map[string][2]int // maps merged text to [left_idx, right_idx]
}

// NewTokenizerSession creates a new tokenization session
func NewTokenizerSession(vocab Vocabulary) *TokenizerSession {
	return &TokenizerSession{
		vocab:    vocab,
		revMerge: make(map[string][2]int),
	}
}

// Tokenize converts input text to tokens using SentencePiece algorithm
func (ts *TokenizerSession) Tokenize(text string) []Token {
	if len(text) == 0 {
		return []Token{}
	}

	// Reset state for new tokenization
	ts.reset()

	// Step 1: Split string into UTF-8 characters
	ts.splitIntoCharacters(text)

	// Step 2: Seed work queue with all possible 2-character tokens
	ts.seedBigrams()

	// Step 3: Keep substituting highest frequency pairs
	ts.mergeBigrams()

	// Step 4: Convert remaining symbols to final tokens
	return ts.convertToTokens()
}

// reset clears the session state
func (ts *TokenizerSession) reset() {
	ts.symbols = ts.symbols[:0]
	ts.workQueue = ts.workQueue[:0]
	// Clear map more efficiently
	if len(ts.revMerge) > 0 {
		ts.revMerge = make(map[string][2]int)
	}
}

// splitIntoCharacters splits text into UTF-8 characters
func (ts *TokenizerSession) splitIntoCharacters(text string) {
	offset := 0
	index := 0

	for offset < len(text) {
		r, runeLen := utf8.DecodeRuneInString(text[offset:])
		if r == utf8.RuneError && runeLen == 1 {
			// Handle invalid UTF-8 as single byte
			runeLen = 1
		}

		// Ensure we don't exceed string bounds
		if offset+runeLen > len(text) {
			runeLen = len(text) - offset
		}

		symbol := Symbol{
			Text: text[offset : offset+runeLen],
			N:    runeLen,
			Prev: -1,
			Next: -1,
		}

		// Set prev/next pointers
		if index > 0 {
			symbol.Prev = index - 1
			ts.symbols[index-1].Next = index
		}

		ts.symbols = append(ts.symbols, symbol)
		offset += runeLen
		index++
	}
}

// seedBigrams adds all possible 2-character combinations to work queue
func (ts *TokenizerSession) seedBigrams() {
	for i := 1; i < len(ts.symbols); i++ {
		ts.tryAddBigram(i-1, i)
	}
}

// mergeBigrams repeatedly merges the highest scoring bigram pairs
func (ts *TokenizerSession) mergeBigrams() {
	for ts.workQueue.Len() > 0 {
		bigram := heap.Pop(&ts.workQueue).(*Bigram)

		leftSym := &ts.symbols[bigram.Left]
		rightSym := &ts.symbols[bigram.Right]

		// Skip if symbols were already merged or size changed
		if leftSym.N == 0 || rightSym.N == 0 || leftSym.N+rightSym.N != bigram.Size {
			continue
		}

		// Merge right symbol into left symbol
		leftSym.Text = leftSym.Text + rightSym.Text
		leftSym.N += rightSym.N
		rightSym.N = 0 // Mark as deleted

		// Update linked list: remove right symbol from chain
		leftSym.Next = rightSym.Next
		if rightSym.Next >= 0 {
			ts.symbols[rightSym.Next].Prev = bigram.Left
		}

		// Find new potential merges with neighbors
		if leftSym.Prev >= 0 {
			ts.tryAddBigram(leftSym.Prev, bigram.Left)
		}
		if leftSym.Next >= 0 {
			ts.tryAddBigram(bigram.Left, leftSym.Next)
		}
	}
}

// convertToTokens converts remaining symbols to final token sequence
func (ts *TokenizerSession) convertToTokens() []Token {
	output := make([]Token, 0, len(ts.symbols))

	// Walk through the linked list of remaining symbols
	for i := 0; i >= 0 && i < len(ts.symbols); i = ts.symbols[i].Next {
		if ts.symbols[i].N > 0 { // Skip deleted symbols
			ts.resegment(&ts.symbols[i], &output)
		}
	}

	return output
}

// tryAddBigram attempts to add a bigram to the work queue
func (ts *TokenizerSession) tryAddBigram(left, right int) {
	if left < 0 || right < 0 || left >= len(ts.symbols) || right >= len(ts.symbols) {
		return
	}

	leftSym := &ts.symbols[left]
	rightSym := &ts.symbols[right]

	// Skip deleted symbols
	if leftSym.N == 0 || rightSym.N == 0 {
		return
	}

	// Combine text from both symbols
	combinedText := leftSym.Text + rightSym.Text
	token := ts.vocab.TextToToken(combinedText)

	if token == NullToken || !ts.vocab.IsValidToken(token) {
		return
	}

	// Create and add the bigram
	bigram := &Bigram{
		Left:  left,
		Right: right,
		Score: ts.vocab.GetTokenScore(token),
		Size:  len(combinedText), // Store size in bytes
	}

	heap.Push(&ts.workQueue, bigram)

	// Record this merge for potential use in resegment
	ts.revMerge[combinedText] = [2]int{left, right}
}

// resegment recursively converts a symbol to tokens
func (ts *TokenizerSession) resegment(symbol *Symbol, output *[]Token) {
	// Try to find the symbol directly in vocabulary
	token := ts.vocab.TextToToken(symbol.Text)
	if token != NullToken && ts.vocab.IsValidToken(token) {
		*output = append(*output, token)
		return
	}

	// Check if this text was created by a previous merge
	if mergeInfo, exists := ts.revMerge[symbol.Text]; exists {
		leftIdx, rightIdx := mergeInfo[0], mergeInfo[1]
		if leftIdx < len(ts.symbols) && rightIdx < len(ts.symbols) {
			ts.resegment(&ts.symbols[leftIdx], output)
			ts.resegment(&ts.symbols[rightIdx], output)
			return
		}
	}

	// Fallback: output individual bytes as tokens
	for i := 0; i < len(symbol.Text); i++ {
		byteToken := ts.vocab.ByteToToken(symbol.Text[i])
		*output = append(*output, byteToken)
	}
}

// SimpleVocab is a basic vocabulary implementation for testing
type SimpleVocab struct {
	textToToken map[string]Token
	tokenScores map[Token]float64
	byteTokens  [256]Token // Direct byte to token mapping
}

// NewSimpleVocab creates a new simple vocabulary
func NewSimpleVocab() *SimpleVocab {
	v := &SimpleVocab{
		textToToken: make(map[string]Token),
		tokenScores: make(map[Token]float64),
	}
	// Initialize byte tokens (256-511 range for bytes)
	for i := 0; i < 256; i++ {
		v.byteTokens[i] = Token(256 + i)
		v.tokenScores[Token(256+i)] = -float64(i) // Lower scores for byte fallback
	}
	return v
}

// AddToken adds a token to the vocabulary
func (v *SimpleVocab) AddToken(text string, token Token, score float64) {
	v.textToToken[text] = token
	v.tokenScores[token] = score
}

func (v *SimpleVocab) TextToToken(text string) Token {
	if token, exists := v.textToToken[text]; exists {
		return token
	}
	return NullToken
}

func (v *SimpleVocab) ByteToToken(b byte) Token {
	return v.byteTokens[b]
}

func (v *SimpleVocab) GetTokenScore(token Token) float64 {
	if score, exists := v.tokenScores[token]; exists {
		return score
	}
	return -1000.0 // Very low default score
}

func (v *SimpleVocab) IsValidToken(token Token) bool {
	_, exists := v.tokenScores[token]
	return exists
}

// Helper function to display tokens with their text representation
func displayTokens(tokens []Token, vocab *SimpleVocab) {
	fmt.Printf("Tokens: ")
	for _, tok := range tokens {
		// Find the text for this token (reverse lookup for display)
		text := "<byte>"
		for t, tokID := range vocab.textToToken {
			if tokID == tok {
				text = t
				break
			}
		}
		if text == "<byte>" && tok >= 256 && tok < 512 {
			text = fmt.Sprintf("<0x%02X>", int(tok-256))
		}
		fmt.Printf("[%d:%s] ", tok, text)
	}
	fmt.Println()
}

func main() {
	// Create a vocabulary with common English bigrams and words
	vocab := NewSimpleVocab()

	// Add some common tokens with scores (higher score = higher priority)
	// Common bigrams
	vocab.AddToken("th", 100, 15.0)
	vocab.AddToken("he", 101, 14.0)
	vocab.AddToken("in", 102, 13.0)
	vocab.AddToken("er", 103, 12.0)
	vocab.AddToken("an", 104, 11.0)

	// Common words
	vocab.AddToken("the", 200, 20.0)
	vocab.AddToken("hello", 201, 18.0)
	vocab.AddToken("world", 202, 17.0)
	vocab.AddToken("and", 203, 16.0)

	// Longer sequences
	vocab.AddToken("ing", 300, 10.0)
	vocab.AddToken("tion", 301, 9.0)

	// Create tokenizer session
	session := NewTokenizerSession(vocab)

	// Example 1: Simple word tokenization
	fmt.Println("Example 1: Tokenizing 'hello'")
	tokens := session.Tokenize("hello")
	displayTokens(tokens, vocab)

	// Example 2: Text with known bigrams
	fmt.Println("\nExample 2: Tokenizing 'the'")
	tokens = session.Tokenize("the")
	displayTokens(tokens, vocab)

	// Example 3: Mixed known and unknown text
	fmt.Println("\nExample 3: Tokenizing 'hello world'")
	tokens = session.Tokenize("hello world")
	displayTokens(tokens, vocab)

	// Example 4: Text that will be broken into bigrams
	fmt.Println("\nExample 4: Tokenizing 'there'")
	tokens = session.Tokenize("there")
	displayTokens(tokens, vocab)

	// Example 5: Completely unknown text (falls back to bytes)
	fmt.Println("\nExample 5: Tokenizing 'xyz'")
	tokens = session.Tokenize("xyz")
	displayTokens(tokens, vocab)

	// Example 6: UTF-8 handling
	fmt.Println("\nExample 6: Tokenizing '世界' (UTF-8)")
	tokens = session.Tokenize("世界")
	displayTokens(tokens, vocab)

	// Example 7: Demonstrating greedy merging
	fmt.Println("\nExample 7: Tokenizing 'mother' (demonstrates bigram merging)")
	// Add specific tokens to show the merging process
	vocab.AddToken("mo", 400, 8.0)
	vocab.AddToken("ther", 401, 7.0)
	vocab.AddToken("mother", 402, 25.0) // High score for full word
	tokens = session.Tokenize("mother")
	displayTokens(tokens, vocab)
}
