package config

import (
	"log"
	"os"
	"strconv"
)

// Config holds all configuration values
type Config struct {
	OpenAIAPIKey      string
	OpenAIModel       string
	OpenAITemperature float32
	OpenAIMaxTokens   int
	Port              string
}

// NewConfig creates a new Config instance with values from environment variables
func NewConfig() *Config {
	cfg := &Config{
		OpenAIAPIKey:      os.Getenv("OPENAI_API_KEY"),
		OpenAIModel:       os.Getenv("OPENAI_MODEL"),
		OpenAITemperature: 0.7,    // Default value
		OpenAIMaxTokens:   1000,   // Default value
		Port:              "8080", // Default value
	}

	// Parse temperature if provided
	if tempStr := os.Getenv("OPENAI_TEMPERATURE"); tempStr != "" {
		if temp, err := strconv.ParseFloat(tempStr, 32); err == nil {
			cfg.OpenAITemperature = float32(temp)
		} else {
			log.Printf("Warning: Invalid OPENAI_TEMPERATURE value, using default: %v", err)
		}
	}

	// Parse max tokens if provided
	if maxTokensStr := os.Getenv("OPENAI_MAX_TOKENS"); maxTokensStr != "" {
		if maxTokens, err := strconv.Atoi(maxTokensStr); err == nil {
			cfg.OpenAIMaxTokens = maxTokens
		} else {
			log.Printf("Warning: Invalid OPENAI_MAX_TOKENS value, using default: %v", err)
		}
	}

	// Get port if provided
	if port := os.Getenv("PORT"); port != "" {
		cfg.Port = port
	}

	return cfg
}
