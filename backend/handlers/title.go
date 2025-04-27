package handlers

import (
	"context"
	"net/http"
	"os"

	"ai-chat/models"
	"ai-chat/openai"

	"github.com/gin-gonic/gin"
)

// HandleTitle handles title generation requests
func HandleTitle(c *gin.Context) {
	var req models.ChatCompletionRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Initialize OpenAI client
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	// Create chat completion for title generation
	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate title"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
