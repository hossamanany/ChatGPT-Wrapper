package handlers

import (
	"encoding/json"
	"net/http"

	"chatgpt-wrapper/models"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

// HandleStream handles streaming chat completion requests
func HandleStream(c *gin.Context) {
	var req models.StreamRequestMessages
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Check if messages array is empty
	if len(req.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Messages array cannot be empty"})
		return
	}

	// Validate message content
	lastMessage := req.Messages[len(req.Messages)-1].Content
	isValid, _ := ValidateMessageContent(lastMessage)

	// Set headers for streaming response
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Flush()

	// If content is invalid, send warning message and return
	if !isValid {
		response := openai.ChatCompletionStreamResponse{
			Choices: []openai.ChatCompletionStreamChoice{
				{
					Delta: openai.ChatCompletionStreamChoiceDelta{
						Content: "I apologize, but I cannot process messages with inappropriate content. Please rephrase your message appropriately.",
					},
				},
			},
		}
		jsonData, _ := json.Marshal(response)
		c.Writer.Write(jsonData)
		c.Writer.Write([]byte("\n"))
		c.Writer.Flush()
		return
	}

	// Initialize OpenAI client
	client := openai.NewClient(cfg.OpenAIAPIKey)

	// Convert messages to OpenAI format
	messages := make([]openai.ChatCompletionMessage, len(req.Messages))
	for i, msg := range req.Messages {
		messages[i] = openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// Create streaming chat completion
	stream, err := client.CreateChatCompletionStream(
		c.Request.Context(),
		openai.ChatCompletionRequest{
			Model:       cfg.OpenAIModel,
			Messages:    messages,
			Temperature: cfg.OpenAITemperature,
			MaxTokens:   cfg.OpenAIMaxTokens,
			Stream:      true,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create stream"})
		return
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if err != nil {
			break
		}

		// Marshal the response into JSON
		jsonData, err := json.Marshal(response)
		if err != nil {
			continue
		}

		// Write it directly without "data:" prefix
		c.Writer.Write(jsonData)
		c.Writer.Write([]byte("\n"))
		c.Writer.Flush()
	}
}
