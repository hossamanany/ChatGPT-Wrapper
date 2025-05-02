package handlers

import (
	"log"
	"net/http"

	"chatgpt-wrapper/models"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

// HandleTitle handles title generation requests
func HandleTitle(c *gin.Context) {
	var req models.RequestMessage
	if err := c.BindJSON(&req); err != nil {
		log.Printf("Error binding title request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if req.Content == "" {
		log.Printf("Empty content received for title generation")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content cannot be empty"})
		return
	}

	// Convert message to OpenAI format with title generation prompt
	messages := []openai.ChatCompletionMessage{
		{
			Role:    req.Role,
			Content: "Summarize the input as title of no more than 5 words. Output only the summarized title. The input is: " + req.Content,
		},
	}

	// Create chat completion for title generation
	resp, err := openAIService.CreateChatCompletion(c.Request.Context(), messages)
	if err != nil {
		log.Printf("Error generating title for content '%s': %v", req.Content, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate title"})
		return
	}

	if len(resp.Choices) == 0 {
		log.Printf("No choices returned from title generation for content '%s'", req.Content)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No title generated"})
		return
	}

	title := resp.Choices[0].Message.Content
	log.Printf("Title generated successfully for content '%s': %s", req.Content, title)
	c.JSON(http.StatusOK, resp)
}
