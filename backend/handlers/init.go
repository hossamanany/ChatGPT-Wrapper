package handlers

import (
	"chatgpt-wrapper/config"
	"chatgpt-wrapper/services"
)

var cfg *config.Config
var openAIService *services.OpenAIService

// InitHandlers initializes all handlers with the given configuration
func InitHandlers(config *config.Config) {
	cfg = config
	openAIService = services.NewOpenAIService(config)
}
