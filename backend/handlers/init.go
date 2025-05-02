package handlers

import "chatgpt-wrapper/config"

var cfg *config.Config

// InitHandlers initializes all handlers with the given configuration
func InitHandlers(config *config.Config) {
	cfg = config
}
