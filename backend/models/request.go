package models

// RequestMessage represents a single message in the chat
type RequestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// StreamRequestMessages represents a request with multiple messages
type StreamRequestMessages struct {
	Messages []RequestMessage `json:"messages"`
}
