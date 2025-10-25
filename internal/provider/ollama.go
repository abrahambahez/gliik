package provider

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// OllamaProvider implements the LLMProvider interface for Ollama's local LLM runtime.
// It handles communication with a local Ollama server for executing AI completions.
type OllamaProvider struct {
	Endpoint string
	Model    string
}

var _ LLMProvider = (*OllamaProvider)(nil)

// NewOllamaProvider creates a new OllamaProvider instance with the specified endpoint
// and model. The endpoint should be the full URL to the Ollama server (e.g.,
// "http://localhost:11434"), and the model should be a valid Ollama model name.
func NewOllamaProvider(endpoint, model string) *OllamaProvider {
	return &OllamaProvider{
		Endpoint: endpoint,
		Model:    model,
	}
}

// StreamCompletion sends a request to the Ollama server with the given systemPrompt
// and userMessage, then streams the response directly to stdout in real-time.
// The systemPrompt and userMessage are combined into a single prompt for Ollama.
// Returns an error if the connection fails or if there's an issue with the response.
func (o *OllamaProvider) StreamCompletion(systemPrompt, userMessage string) error {
	combinedPrompt := systemPrompt + "\n\n" + userMessage

	reqBody := map[string]interface{}{
		"model":  o.Model,
		"prompt": combinedPrompt,
		"stream": true,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/generate", o.Endpoint)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error: Cannot connect to Ollama\n\nMake sure Ollama is running:\n  ollama serve\n\nConfigured endpoint: %s", o.Endpoint)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ollama API error (status %d)", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Bytes()

		var response map[string]interface{}
		if err := json.Unmarshal(line, &response); err != nil {
			continue
		}

		if done, ok := response["done"].(bool); ok && done {
			break
		}

		if responseText, ok := response["response"].(string); ok {
			fmt.Print(responseText)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	return nil
}
