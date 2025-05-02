package services

import (
	"context"

	"chatgpt-wrapper/config"

	"github.com/sashabaranov/go-openai"
)

// OpenAIService handles OpenAI client operations
type OpenAIService struct {
	client *openai.Client
	cfg    *config.Config
}

// NewOpenAIService creates a new OpenAIService instance
func NewOpenAIService(cfg *config.Config) *OpenAIService {
	return &OpenAIService{
		client: openai.NewClient(cfg.OpenAIAPIKey),
		cfg:    cfg,
	}
}

// CreateChatCompletion creates a chat completion request
func (s *OpenAIService) CreateChatCompletion(ctx context.Context, messages []openai.ChatCompletionMessage) (openai.ChatCompletionResponse, error) {
	return s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       s.cfg.OpenAIModel,
			Messages:    messages,
			Temperature: s.cfg.OpenAITemperature,
			MaxTokens:   s.cfg.OpenAIMaxTokens,
		},
	)
}

// CreateChatCompletionStream creates a streaming chat completion request
func (s *OpenAIService) CreateChatCompletionStream(ctx context.Context, messages []openai.ChatCompletionMessage) (*openai.ChatCompletionStream, error) {
	return s.client.CreateChatCompletionStream(
		ctx,
		openai.ChatCompletionRequest{
			Model:       s.cfg.OpenAIModel,
			Messages:    messages,
			Temperature: s.cfg.OpenAITemperature,
			MaxTokens:   s.cfg.OpenAIMaxTokens,
			Stream:      true,
		},
	)
}
