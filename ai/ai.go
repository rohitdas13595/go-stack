package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client is an OpenAI-compatible chat client.
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
	Model      string
}

// Message is a chat message.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// CompletionRequest is a chat completion request.
type CompletionRequest struct {
	Model string    `json:"model"`
	Messages []Message `json:"messages"`
}

// CompletionResponse minimal parse.
type CompletionResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// Complete calls POST /v1/chat/completions.
func (c *Client) Complete(ctx context.Context, req *CompletionRequest) (string, error) {
	if c.HTTPClient == nil {
		c.HTTPClient = http.DefaultClient
	}
	if req.Model == "" {
		req.Model = c.Model
	}
	b, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	url := c.BaseURL + "/v1/chat/completions"
	if c.BaseURL == "" {
		url = "https://api.openai.com/v1/chat/completions"
	}
	hreq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	hreq.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		hreq.Header.Set("Authorization", "Bearer "+c.APIKey)
	}
	resp, err := c.HTTPClient.Do(hreq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("ai: %s: %s", resp.Status, string(body))
	}
	var out CompletionResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return "", err
	}
	if len(out.Choices) == 0 {
		return "", fmt.Errorf("ai: empty choices")
	}
	return out.Choices[0].Message.Content, nil
}

// StreamChunk is one streamed delta (simplified).
type StreamChunk struct {
	Delta string
}

// Stream calls chat completions; MVP uses non-streaming Complete and emits one chunk.
func (c *Client) Stream(ctx context.Context, req *CompletionRequest, onChunk func(StreamChunk) error) error {
	if c.HTTPClient == nil {
		c.HTTPClient = http.DefaultClient
	}
	s, err := c.Complete(ctx, req)
	if err != nil {
		return err
	}
	return onChunk(StreamChunk{Delta: s})
}

// Embed is a placeholder for embedding APIs (pgvector / external).
func Embed(ctx context.Context, client *Client, text string) ([]float32, error) {
	_ = ctx
	_ = client
	_ = text
	return nil, fmt.Errorf("ai: Embed not implemented; use provider-specific client")
}

// Agent runs a multi-step tool loop (stub).
type Agent struct {
	Model string
	MaxSteps int
}

// Run executes the agent (stub: returns prompt echo).
func (a *Agent) Run(ctx context.Context, prompt string) (string, error) {
	_ = ctx
	_ = a
	return prompt, nil
}
