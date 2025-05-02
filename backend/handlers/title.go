package handlers

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	"chatgpt-wrapper/models"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

// HandleTitle handles title generation requests
func HandleTitle(c *gin.Context) {
	var req models.RequestMessage
	if err := c.BindJSON(&req); err != nil {
		log.Printf("Error binding request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Content == "" {
		log.Printf("No messages provided for title generation")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No messages provided"})
		return
	}

	// Initialize OpenAI client
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	model := os.Getenv("OPENAI_MODEL")
	temperature, _ := strconv.ParseFloat(os.Getenv("OPENAI_TEMPERATURE"), 32)
	maxTokens, _ := strconv.Atoi(os.Getenv("OPENAI_MAX_TOKENS"))

	// Convert message to OpenAI format with title generation prompt
	messages := []openai.ChatCompletionMessage{
		{
			Role:    req.Role,
			Content: "Summarize the input as title of no more than 5 words. Output only the summarized title. The input is: " + req.Content,
		},
	}

	// Create chat completion for title generation
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       model,
			Messages:    messages,
			Temperature: float32(temperature),
			MaxTokens:   maxTokens,
		},
	)
	if err != nil {
		log.Printf("Error generating title: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate title"})
		return
	}

	log.Printf("Title generated successfully: %s", resp.Choices[0].Message.Content)
	c.JSON(http.StatusOK, resp)
}
