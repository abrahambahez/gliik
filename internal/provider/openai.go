package provider

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// OpenAIProvider implements the LLMProvider interface for OpenAI's API.
// It handles authentication and communication with OpenAI or OpenAI-compatible endpoints.
type OpenAIProvider struct {
	APIKey   string
	Model    string
	Endpoint string
}

var _ LLMProvider = (*OpenAIProvider)(nil)

// NewOpenAIProvider creates a new OpenAIProvider instance by reading the
// OPENAI_API_KEY environment variable and using the provided endpoint and model.
// Returns an error with clear instructions if the API key is not set.
func NewOpenAIProvider(endpoint, model string) (*OpenAIProvider, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is not set\n\nTo use OpenAI provider, set your API key:\n  export OPENAI_API_KEY=sk-...")
	}

	if endpoint == "" {
		return nil, fmt.Errorf("OpenAI endpoint cannot be empty")
	}

	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		return nil, fmt.Errorf("OpenAI endpoint must start with http:// or https://\nProvided: %s", endpoint)
	}

	normalizedEndpoint := strings.TrimSuffix(endpoint, "/")

	return &OpenAIProvider{
		APIKey:   apiKey,
		Model:    model,
		Endpoint: normalizedEndpoint,
	}, nil
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIRequest struct {
	Model    string          `json:"model"`
	Messages []openAIMessage `json:"messages"`
	Stream   bool            `json:"stream"`
}

type openAIStreamResponse struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
}

// StreamCompletion sends a request to the OpenAI API with the given systemPrompt
// and userMessage, then streams the response directly to stdout in real-time.
// The systemPrompt and userMessage are sent as separate messages in the conversation.
func (o *OpenAIProvider) StreamCompletion(systemPrompt, userMessage string) error {
	messages := []openAIMessage{}

	if systemPrompt != "" {
		messages = append(messages, openAIMessage{
			Role:    "system",
			Content: systemPrompt,
		})
	}

	messages = append(messages, openAIMessage{
		Role:    "user",
		Content: userMessage,
	})

	reqBody := openAIRequest{
		Model:    o.Model,
		Messages: messages,
		Stream:   true,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := o.Endpoint + "/chat/completions"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error: Cannot connect to OpenAI API\n\nPlease check your network connection and endpoint configuration.\n\nConfigured endpoint: %s\nError: %w", o.Endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := bufio.NewReader(resp.Body).ReadString('\n')

		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return fmt.Errorf("Error: Invalid OpenAI API key\n\nThe OPENAI_API_KEY you provided is invalid or expired.\nPlease check your API key at: https://platform.openai.com/api-keys")
		case http.StatusTooManyRequests:
			return fmt.Errorf("Error: OpenAI rate limit exceeded\n\nYou have exceeded your API rate limit.\nPlease wait a moment and try again, or check your usage at: https://platform.openai.com/usage")
		case http.StatusForbidden:
			return fmt.Errorf("Error: Access forbidden\n\nYour API key does not have permission to access this resource.\nPlease check your API key permissions.")
		case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
			return fmt.Errorf("Error: OpenAI service unavailable\n\nOpenAI's servers are experiencing issues. Please try again later.\nStatus: %d", resp.StatusCode)
		default:
			return fmt.Errorf("OpenAI API error (status %d): %s", resp.StatusCode, body)
		}
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")

		if data == "[DONE]" {
			break
		}

		var streamResp openAIStreamResponse
		if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
			continue
		}

		if len(streamResp.Choices) > 0 {
			content := streamResp.Choices[0].Delta.Content
			if content != "" {
				fmt.Print(content)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading stream: %w", err)
	}

	return nil
}
