// Package models defines the data structures used for API requests and responses.
// These structures are used to serialize and deserialize JSON data between the client and server.
package models

// RequestMessage represents a single message in a chat conversation.
// It is used to structure the input for chat completion requests.
//
// Fields:
//   - Role: The role of the message sender (e.g., "user", "assistant", "system")
//   - Content: The actual text content of the message
type RequestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// StreamRequestMessages represents a collection of messages for streaming chat completion.
// It is used as the request body for the streaming chat endpoint.
//
// Fields:
//   - Messages: A slice of RequestMessage containing the conversation history
type StreamRequestMessages struct {
	Messages []RequestMessage `json:"messages"`
}
