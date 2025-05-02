// Package config manages application configuration and environment variables.
// It provides functionality to load and validate configuration values from environment variables
// with sensible defaults for missing values.
package config

import (
	"log"
	"os"
	"strconv"
)

// Config holds the application configuration values.
// It contains settings for OpenAI API integration and server configuration.
//
// Fields:
//   - OpenAIAPIKey: The API key for OpenAI services
//   - OpenAIModel: The OpenAI model to use (e.g., "gpt-3.5-turbo")
//   - OpenAITemperature: Controls randomness in responses (0.0 to 1.0)
//   - OpenAIMaxTokens: Maximum number of tokens to generate in responses
//   - Port: The port number the server will listen on
type Config struct {
	OpenAIAPIKey      string
	OpenAIModel       string
	OpenAITemperature float32
	OpenAIMaxTokens   int
	Port              string
}

// NewConfig creates a new Config instance with values from environment variables.
// It loads configuration values from environment variables with the following precedence:
// 1. Environment variable value if set
// 2. Default value if environment variable is not set or invalid
//
// Returns:
//   - *Config: A pointer to a new Config instance with loaded values
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
