package handlers

import (
	"strings"
	"unicode"

	"ai-chat/config"

	"github.com/kljensen/snowball"
)

const maxMessageLength = 1000 // Maximum allowed characters in a message

// normalize performs text normalization including:
// - Converting to lowercase
// - Removing punctuation
// - Removing extra whitespace
func normalize(input string) string {
	// Convert to lowercase
	input = strings.ToLower(input)

	// Remove punctuation and normalize whitespace
	input = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsSpace(r) {
			return r
		}
		return -1
	}, input)

	// Normalize whitespace
	return strings.Join(strings.Fields(input), " ")
}

// containsBannedWords checks if the input contains any banned words using NLP techniques:
// 1. Normalization
// 2. Tokenization
// 3. Stemming
func containsBannedWords(input string) bool {
	normalized := normalize(input)
	words := strings.Fields(normalized) // Tokenize

	for _, word := range words {
		// Get the stem of the word
		stemmed, err := snowball.Stem(word, "english", true)
		if err != nil {
			// If stemming fails, use the original word
			stemmed = word
		}

		// Check against banned stems
		if _, exists := config.GetBannedStems()[stemmed]; exists {
			return true
		}
	}
	return false
}

// validateMessageLength checks if a message is within the allowed length limit
func validateMessageLength(message string) bool {
	return len(strings.TrimSpace(message)) <= maxMessageLength
}

// ValidateMessageContent performs all content moderation checks on a message
func ValidateMessageContent(message string) (bool, string) {
	if !validateMessageLength(message) {
		return false, "Message exceeds maximum length"
	}

	if containsBannedWords(message) {
		return false, "Message contains inappropriate content"
	}

	return true, ""
}
