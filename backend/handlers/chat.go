package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"chatgpt-wrapper/models"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

// HandleStream handles streaming chat completion requests
func HandleStream(c *gin.Context) {
	var req models.ChatCompletionRequest
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

	// Set headers for SSE
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
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
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	model := os.Getenv("OPENAI_MODEL")
	temperature, _ := strconv.ParseFloat(os.Getenv("OPENAI_TEMPERATURE"), 32)
	maxTokens, _ := strconv.Atoi(os.Getenv("OPENAI_MAX_TOKENS"))

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
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       model,
			Messages:    messages,
			Temperature: float32(temperature),
			MaxTokens:   maxTokens,
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
