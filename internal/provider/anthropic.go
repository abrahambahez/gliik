package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// AnthropicProvider implements the LLMProvider interface for Anthropic's Claude API.
// It handles authentication and communication with the Anthropic API endpoint.
type AnthropicProvider struct {
	APIKey string
	Model  string
}

var _ LLMProvider = (*AnthropicProvider)(nil)

type messageRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []message `json:"messages"`
	System    string    `json:"system,omitempty"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type messageResponse struct {
	Content []contentBlock `json:"content"`
}

type contentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// NewAnthropicProvider creates a new AnthropicProvider instance by reading the
// ANTHROPIC_API_KEY environment variable and using the provided model.
// Returns an error if the API key is not set.
func NewAnthropicProvider(model string) (*AnthropicProvider, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY environment variable is not set")
	}

	return &AnthropicProvider{
		APIKey: apiKey,
		Model:  model,
	}, nil
}

// StreamCompletion sends a request to the Anthropic API with the given systemPrompt
// and userMessage, then prints the response directly to stdout in real-time.
// The systemPrompt provides context and instructions, while userMessage contains
// the actual user input or request.
func (a *AnthropicProvider) StreamCompletion(systemPrompt, userMessage string) error {
	reqBody := messageRequest{
		Model:     a.Model,
		MaxTokens: 4096,
		System:    systemPrompt,
		Messages: []message{
			{
				Role:    "user",
				Content: userMessage,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", a.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response messageResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Content) == 0 {
		return fmt.Errorf("no content in response")
	}

	fmt.Print(response.Content[0].Text)
	return nil
}
