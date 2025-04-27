package handlers

import (
	"context"
	"net/http"
	"os"
	"strings"

	"ai-chat/models"
	"ai-chat/openai"

	"github.com/gin-gonic/gin"
)

type ChatRequest struct {
	Message string `json:"message"`
}

// HandleChat handles regular chat completion requests
func HandleChat(c *gin.Context) {
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

	// Initialize OpenAI client
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	// Create chat completion
	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get response from OpenAI"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

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

	// Create streaming chat completion
	stream, err := client.CreateChatCompletionStream(context.Background(), req)
	if err != nil {
		c.SSEvent("error", "Failed to create stream")
		c.Writer.Flush()
		return
	}
	defer stream.Close()

	buffer := ""
	for {
		response, err := stream.Read()
		if err != nil {
			if err.Error() == "stream closed" {
				// Send any remaining buffered content
				if buffer != "" {
					c.SSEvent("message", buffer)
					c.Writer.Flush()
				}
				break
			}
			c.SSEvent("error", "Error receiving stream")
			c.Writer.Flush()
			return
		}

		if len(response.Choices) > 0 && response.Choices[0].Delta.Content != "" {
			content := response.Choices[0].Delta.Content
			buffer += content

			// Send buffer when we encounter sentence endings or certain punctuation
			if strings.ContainsAny(content, ".!?,") || strings.Contains(content, "\n") {
				c.SSEvent("message", buffer)
				c.Writer.Flush()
				buffer = ""
			}
		}
	}
}
