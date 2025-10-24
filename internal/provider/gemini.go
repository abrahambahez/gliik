package provider

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// GeminiProvider implements the LLMProvider interface for Google's Gemini API.
type GeminiProvider struct {
	APIKey string
	Model  string
}

var _ LLMProvider = (*GeminiProvider)(nil)

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiStreamResponse struct {
	Candidates []struct {
		Content struct {
			Parts []geminiPart `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// NewGeminiProvider creates a new GeminiProvider instance by reading the
// GOOGLE_API_KEY environment variable and using the provided model.
func NewGeminiProvider(model string) (*GeminiProvider, error) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GOOGLE_API_KEY environment variable is not set\n\nTo use Gemini provider, set your API key:\n  export GOOGLE_API_KEY=your-api-key")
	}

	return &GeminiProvider{
		APIKey: apiKey,
		Model:  model,
	}, nil
}

// StreamCompletion sends a request to the Gemini API with the given systemPrompt
// and userMessage, then streams the response directly to stdout in real-time.
func (g *GeminiProvider) StreamCompletion(systemPrompt, userMessage string) error {
	reqBody := geminiRequest{
		Contents: []geminiContent{
			{
				Parts: []geminiPart{
					{Text: systemPrompt + "\n\n" + userMessage},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:streamGenerateContent?alt=sse&key=%s", g.Model, g.APIKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error: Cannot connect to Gemini API\n\nPlease check your network connection.\n\nError: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		switch resp.StatusCode {
		case http.StatusUnauthorized, http.StatusForbidden:
			return fmt.Errorf("Error: Invalid Gemini API key\n\nThe GOOGLE_API_KEY you provided is invalid or expired.\nPlease check your API key at: https://aistudio.google.com/app/apikey")
		case http.StatusTooManyRequests:
			return fmt.Errorf("Error: Gemini rate limit exceeded\n\nYou have exceeded your API rate limit.\nPlease wait a moment and try again.")
		case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
			return fmt.Errorf("Error: Gemini service unavailable\n\nGemini's servers are experiencing issues. Please try again later.\nStatus: %d", resp.StatusCode)
		default:
			return fmt.Errorf("Gemini API error (status %d): %s", resp.StatusCode, string(body))
		}
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" || data == "" {
			continue
		}

		var streamResp geminiStreamResponse
		if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
			continue
		}

		if len(streamResp.Candidates) > 0 && len(streamResp.Candidates[0].Content.Parts) > 0 {
			if text := streamResp.Candidates[0].Content.Parts[0].Text; text != "" {
				fmt.Print(text)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading stream: %w", err)
	}

	return nil
}
