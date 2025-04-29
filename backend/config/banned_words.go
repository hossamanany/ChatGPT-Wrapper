package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/kljensen/snowball"
)

var (
	bannedStems map[string]struct{}
	once        sync.Once
)

// LoadBannedWords loads and stems the banned words from the JSON file
func LoadBannedWords() {
	once.Do(func() {
		// Get the absolute path to the banned_words.json file
		configPath := filepath.Join("config", "banned_words.json")
		data, err := os.ReadFile(configPath)
		if err != nil {
			log.Fatalf("Failed to read banned words: %v", err)
		}

		var words []string
		if err := json.Unmarshal(data, &words); err != nil {
			log.Fatalf("Failed to parse banned words: %v", err)
		}

		// Stem each word and store in memory
		bannedStems = make(map[string]struct{}, len(words))
		for _, word := range words {
			stemmed, err := snowball.Stem(word, "english", true)
			if err != nil {
				log.Printf("Warning: Failed to stem word '%s': %v", word, err)
				stemmed = word
			}
			bannedStems[stemmed] = struct{}{}
		}

		log.Printf("Loaded %d banned word stems", len(bannedStems))
	})
}

// GetBannedStems returns the cached list of banned word stems
func GetBannedStems() map[string]struct{} {
	return bannedStems
}
