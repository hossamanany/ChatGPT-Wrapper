package handlers

import (
	"context"
	"net/http"
	"os"

	"ai-chat/models"

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

	// Validate message content
	lastMessage := req.Messages[len(req.Messages)-1].Content
	if isValid, errMsg := ValidateMessageContent(lastMessage); !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	// Set headers for SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.Flush()

	// Initialize OpenAI client
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

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
			Model:       req.Model,
			Messages:    messages,
			Temperature: req.Temperature,
			MaxTokens:   req.MaxTokens,
			Stream:      true,
		},
	)
	if err != nil {
		c.SSEvent("error", "Failed to create stream")
		c.Writer.Flush()
		return
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if err != nil {
			if err.Error() == "stream closed" {
				break
			}
			c.SSEvent("error", "Error receiving stream")
			c.Writer.Flush()
			return
		}

		if len(response.Choices) > 0 && response.Choices[0].Delta.Content != "" {
			c.SSEvent("message", response.Choices[0].Delta.Content)
			c.Writer.Flush()
		}
	}
}
