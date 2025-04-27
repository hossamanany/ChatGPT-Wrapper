package openai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"ai-chat/models"
)

// Client represents the OpenAI client
type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new OpenAI client
func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:     apiKey,
		BaseURL:    "https://api.openai.com/v1",
		HTTPClient: &http.Client{},
	}
}

// CreateChatCompletion creates a chat completion request
func (c *Client) CreateChatCompletion(ctx context.Context, req models.ChatCompletionRequest) (*models.ChatCompletionResponse, error) {
	url := fmt.Sprintf("%s/chat/completions", c.BaseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var completionResp models.ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&completionResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &completionResp, nil
}

// CreateChatCompletionStream creates a streaming chat completion request
func (c *Client) CreateChatCompletionStream(ctx context.Context, req models.ChatCompletionRequest) (*StreamReader, error) {
	req.Stream = true
	url := fmt.Sprintf("%s/chat/completions", c.BaseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return NewStreamReader(resp.Body), nil
}

// StreamReader handles reading streaming responses
type StreamReader struct {
	reader *bufio.Reader
	closer io.Closer
}

// NewStreamReader creates a new StreamReader
func NewStreamReader(reader io.ReadCloser) *StreamReader {
	return &StreamReader{
		reader: bufio.NewReader(reader),
		closer: reader,
	}
}

// Read reads the next chunk from the stream
func (s *StreamReader) Read() (*models.ChatCompletionStreamResponse, error) {
	line, err := s.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "data: ")
	if line == "[DONE]" {
		return nil, io.EOF
	}

	var resp models.ChatCompletionStreamResponse
	if err := json.Unmarshal([]byte(line), &resp); err != nil {
		return nil, fmt.Errorf("error unmarshaling stream response: %w", err)
	}

	return &resp, nil
}

// Close closes the stream reader
func (s *StreamReader) Close() error {
	return s.closer.Close()
}
