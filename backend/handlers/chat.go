// Package handlers contains HTTP request handlers for the ChatGPT Wrapper API.
// It implements the core chat functionality including streaming responses and message validation.
package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"chatgpt-wrapper/models"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

// HandleStream processes streaming chat completion requests.
// It:
// - Validates the incoming request format and content
// - Sets up streaming response headers
// - Checks message content for inappropriate content
// - Creates and manages the OpenAI chat completion stream
// - Streams the responses back to the client in real-time
//
// Parameters:
//   - c: Gin context containing the HTTP request and response
//
// The function returns a streaming response with chat completions or an error if the request is invalid.
func HandleStream(c *gin.Context) {
	var req models.StreamRequestMessages
	if err := c.BindJSON(&req); err != nil {
		log.Printf("Error binding request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Check if messages array is empty
	if len(req.Messages) == 0 {
		log.Printf("Empty messages array received")
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
		jsonData, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshaling invalid content response: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
			return
		}
		if _, err := c.Writer.Write(jsonData); err != nil {
			log.Printf("Error writing invalid content response: %v", err)
			return
		}
		if _, err := c.Writer.Write([]byte("\n")); err != nil {
			log.Printf("Error writing newline: %v", err)
			return
		}
		c.Writer.Flush()
		return
	}

	// Convert messages to OpenAI format
	messages := make([]openai.ChatCompletionMessage, len(req.Messages))
	for i, msg := range req.Messages {
		messages[i] = openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// Create streaming chat completion
	stream, err := openAIService.CreateChatCompletionStream(c.Request.Context(), messages)
	if err != nil {
		log.Printf("Error creating chat completion stream: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create stream"})
		return
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if err != nil {
			log.Printf("Error receiving stream response: %v", err)
			break
		}

		// Marshal the response into JSON
		jsonData, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshaling stream response: %v", err)
			continue
		}

		// Write it directly without "data:" prefix
		if _, err := c.Writer.Write(jsonData); err != nil {
			log.Printf("Error writing stream response: %v", err)
			return
		}
		if _, err := c.Writer.Write([]byte("\n")); err != nil {
			log.Printf("Error writing newline: %v", err)
			return
		}
		c.Writer.Flush()
	}
}
