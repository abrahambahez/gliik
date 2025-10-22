package provider

// LLMProvider defines the interface for language model providers that can execute
// AI completions with streaming output. Implementations include Anthropic's Claude
// API and Ollama for local model execution.
type LLMProvider interface {
	// StreamCompletion sends a prompt to the language model and streams the response
	// to stdout in real-time. The systemPrompt provides context and instructions,
	// while userMessage contains the actual user input or request.
	StreamCompletion(systemPrompt, userMessage string) error
}
