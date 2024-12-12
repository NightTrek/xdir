package main

import (
	"strings"
	"unicode"
)

// TokenCounter handles token counting for LLM context
type TokenCounter struct {
	count int
}

// NewTokenCounter creates a new TokenCounter
func NewTokenCounter() *TokenCounter {
	return &TokenCounter{count: 0}
}

// CountTokens provides a simple estimation of tokens in text
// This is a basic implementation that splits on whitespace and punctuation
// For production use, you would want to use a proper tokenizer matching your LLM
func (tc *TokenCounter) CountTokens(text string) int {
	tokens := 0
	inToken := false

	// Split into words first
	words := strings.Fields(text)

	for _, word := range words {
		// Handle each character in the word
		for _, char := range word {
			if unicode.IsPunct(char) {
				// Count punctuation as separate tokens
				tokens++
				inToken = false
			} else if !inToken {
				// Start of new token
				tokens++
				inToken = true
			}
		}
		inToken = false
	}

	tc.count += tokens
	return tokens
}

// GetTotalTokens returns the total number of tokens counted
func (tc *TokenCounter) GetTotalTokens() int {
	return tc.count
}
