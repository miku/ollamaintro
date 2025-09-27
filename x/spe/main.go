package main

import (
	"container/heap"
	"log"
	"unicode/utf8"
)

// Token represents a vocabulary token
type Token int

const (
	NullToken Token = -1
)

// Vocabulary interface - simplified version of what the original code expects
type Vocabulary interface {
	TextToToken(text string) Token
	ByteToToken(b byte) Token
	GetTokenScore(token Token) float64
	IsValidToken(token Token) bool
}

// Symbol represents a text segment during tokenization
type Symbol struct {
	Text string // the actual text content
	N    int    // length of the text segment
	Prev int    // index of previous symbol (-1 if none)
	Next int    // index of next symbol (-1 if none)
}

// Bigram represents a potential merge between two adjacent symbols
type Bigram struct {
	Left  int     // index of left symbol
	Right int     // index of right symbol
	Score float64 // score from vocabulary
	Size  int     // total size of merged text
}

// BigramQueue implements a max-heap for bigrams (highest score first)
type BigramQueue []*Bigram

func (pq BigramQueue) Len() int { return len(pq) }

func (pq BigramQueue) Less(i, j int) bool {
	// Max heap - higher scores have priority
	return pq[i].Score > pq[j].Score
}

func (pq BigramQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *BigramQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Bigram))
}

func (pq *BigramQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
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
	// Clear previous state
	ts.symbols = ts.symbols[:0]
	ts.workQueue = ts.workQueue[:0]
	for k := range ts.revMerge {
		delete(ts.revMerge, k)
	}

	// Step 1: Split string into UTF-8 characters
	ts.splitIntoCharacters(text)

	// Step 2: Seed work queue with all possible 2-character tokens
	ts.seedBigrams()

	// Step 3: Keep substituting highest frequency pairs
	ts.mergeBigrams()

	// Step 4: Convert remaining symbols to final tokens
	return ts.convertToTokens()
}

// splitIntoCharacters splits text into UTF-8 characters and creates symbol chain
func (ts *TokenizerSession) splitIntoCharacters(text string) {
	index := 0
	offset := 0

	for offset < len(text) {
		// Get the length of the current UTF-8 character
		_, runeLen := utf8.DecodeRuneInString(text[offset:])
		if runeLen == 0 {
			runeLen = 1 // Handle invalid UTF-8
		}

		// Ensure we don't go beyond the string
		if offset+runeLen > len(text) {
			runeLen = len(text) - offset
		}

		symbol := Symbol{
			Text: text[offset : offset+runeLen],
			N:    runeLen,
			Prev: index - 1,
		}

		// Set next pointer
		if offset+runeLen == len(text) {
			symbol.Next = -1 // Last symbol
		} else {
			symbol.Next = index + 1
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

		// Skip if one of the symbols was already merged
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

		// Find new potential merges involving the newly merged symbol
		ts.tryAddBigram(leftSym.Prev, bigram.Left)
		ts.tryAddBigram(bigram.Left, leftSym.Next)
	}
}

// convertToTokens converts remaining symbols to final token sequence
func (ts *TokenizerSession) convertToTokens() []Token {
	var output []Token

	// Walk through the linked list of remaining symbols
	for i := 0; i != -1; i = ts.symbols[i].Next {
		if ts.symbols[i].N > 0 { // Skip deleted symbols
			ts.resegment(&ts.symbols[i], &output)
		}
	}

	return output
}

// tryAddBigram attempts to add a bigram to the work queue if it forms a valid token
func (ts *TokenizerSession) tryAddBigram(left, right int) {
	if left == -1 || right == -1 {
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
		Size:  len(combinedText),
	}

	heap.Push(&ts.workQueue, bigram)

	// Record this merge for potential later use in resegment
	ts.revMerge[combinedText] = [2]int{left, right}
}

// resegment recursively converts a symbol to tokens
func (ts *TokenizerSession) resegment(symbol *Symbol, output *[]Token) {
	// Try to find the symbol directly in vocabulary
	token := ts.vocab.TextToToken(symbol.Text)
	if token != NullToken {
		*output = append(*output, token)
		return
	}

	// Check if this text was created by a previous merge
	if mergeInfo, exists := ts.revMerge[symbol.Text]; exists {
		leftIdx, rightIdx := mergeInfo[0], mergeInfo[1]
		ts.resegment(&ts.symbols[leftIdx], output)
		ts.resegment(&ts.symbols[rightIdx], output)
		return
	}

	// Fallback: output individual bytes as tokens
	for i := 0; i < len(symbol.Text); i++ {
		byteToken := ts.vocab.ByteToToken(symbol.Text[i])
		*output = append(*output, byteToken)
	}
}

// Example usage and simple vocabulary implementation for testing
type SimpleVocab struct {
	textToToken map[string]Token
	tokenScores map[Token]float64
}

func NewSimpleVocab() *SimpleVocab {
	return &SimpleVocab{
		textToToken: make(map[string]Token),
		tokenScores: make(map[Token]float64),
	}
}

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
	return Token(b) // Simple mapping for example
}

func (v *SimpleVocab) GetTokenScore(token Token) float64 {
	if score, exists := v.tokenScores[token]; exists {
		return score
	}
	return 0.0
}

func (v *SimpleVocab) IsValidToken(token Token) bool {
	_, exists := v.tokenScores[token]
	return exists
}

// Example usage
func main() {
	// Create a simple vocabulary
	vocab := NewSimpleVocab()
	vocab.AddToken("th", 100, 10.0)
	vocab.AddToken("he", 101, 9.0)
	vocab.AddToken("the", 102, 15.0)
	vocab.AddToken("hello", 103, 12.0)

	// Create tokenizer session
	session := NewTokenizerSession(vocab)

	// Tokenize text
	tokens := session.Tokenize("hello")

	// tokens will contain the result of tokenization
	log.Println(tokens)

}
