# Task List: Ollama Support Implementation

## Relevant Files

- `internal/config/init.go` - Add `provider` field to Config struct
- `internal/provider/provider.go` - Define LLMProvider interface (NEW FILE)
- `internal/provider/anthropic.go` - Refactored Anthropic implementation (NEW FILE)
- `internal/provider/ollama.go` - New Ollama provider implementation (NEW FILE)
- `internal/ai/client.go` - Current Anthropic client (will be refactored into provider)
- `cmd/run.go` - Update to use provider abstraction instead of direct ai.Client

### Notes

- No test files required for MVP (manual testing only)
- Use `go run main.go run <instruction>` to test
- Manual testing requires: `ollama serve` + `ollama pull llama3.2`
- Follow CODE_RULES.md: long function names, Godoc only, no inline comments

## Tasks

- [x] 1.0 Create provider abstraction layer
  - [x] 1.1 Create `internal/provider/` directory
  - [x] 1.2 Create `provider.go` with `LLMProvider` interface definition
  - [x] 1.3 Define `StreamCompletion(systemPrompt, userMessage string) error` method signature
  - [x] 1.4 Add Godoc documentation explaining the provider abstraction pattern

- [x] 2.0 Update configuration to support provider selection
  - [x] 2.1 Add `Provider string` field to `Config` struct in `internal/config/init.go`
  - [x] 2.2 Add `yaml:"provider"` tag to Provider field
  - [x] 2.3 Create `ValidateProvider()` method that checks value is "anthropic" or "ollama"
  - [x] 2.4 Update `Initialize()` to set default `provider: "anthropic"` in config.yaml
  - [x] 2.5 Add Godoc for Provider field explaining valid values
  - [x] 2.6 Add `AnthropicConfig` struct with `Model` field to `Config` struct
  - [x] 2.7 Add `OllamaConfig` struct with `Endpoint` and `Model` fields to `Config` struct
  - [x] 2.8 Add `yaml:"anthropic"` and `yaml:"ollama"` tags to provider config fields
  - [x] 2.9 Update `Initialize()` to set default anthropic and ollama config sections
  - [x] 2.10 Add Godoc for provider config structs explaining their purpose

- [x] 3.0 Refactor Anthropic client into provider implementation
  - [x] 3.1 [depends on: 1.0] Create `internal/provider/anthropic.go` file
  - [x] 3.2 [depends on: 3.1] Copy `Client` struct from `ai/client.go` and rename to `AnthropicProvider`
  - [x] 3.3 [depends on: 3.2] Copy `NewClient()` function and rename to `NewAnthropicProvider()`
  - [x] 3.4 [depends on: 3.3] Copy `Complete()` method and rename to `StreamCompletion()`
  - [x] 3.5 [depends on: 3.4] Update `StreamCompletion()` signature to match interface (systemPrompt, userMessage)
  - [x] 3.6 [depends on: 3.5] Update method to combine systemPrompt and userMessage into messages array
  - [x] 3.7 [depends on: 3.6] Add Godoc documentation for `AnthropicProvider` struct and methods
  - [x] 3.8 [depends on: 3.7] Verify AnthropicProvider implements LLMProvider interface
  - [x] 3.9 [depends on: 2.6] Update `NewAnthropicProvider()` to accept model parameter
  - [x] 3.10 [depends on: 3.9] Remove hardcoded model, use parameter instead

- [x] 4.0 Implement Ollama provider
  - [x] 4.1 [depends on: 1.0] Create `internal/provider/ollama.go` file
  - [x] 4.2 [depends on: 4.1] Define `OllamaProvider` struct with `Endpoint` and `Model` fields
  - [x] 4.3 [depends on: 2.7] Update `NewOllamaProvider()` to accept endpoint and model parameters
  - [x] 4.4 [depends on: 4.3] Remove hardcoded defaults, use parameters instead
  - [x] 4.5 [depends on: 4.4] Implement `StreamCompletion(systemPrompt, userMessage string) error` method
  - [x] 4.6 [depends on: 4.5] Build JSON request body with model, prompt (combined system+user), and stream:true
  - [x] 4.7 [depends on: 4.6] Create HTTP POST request to `{endpoint}/api/generate`
  - [x] 4.8 [depends on: 4.7] Set Content-Type header to `application/json`
  - [x] 4.9 [depends on: 4.8] Execute request and handle connection errors with specific error message from PRD
  - [x] 4.10 [depends on: 4.9] Implement newline-delimited JSON streaming response parser
  - [x] 4.11 [depends on: 4.10] For each line, parse JSON and extract `response` field, write to stdout
  - [x] 4.12 [depends on: 4.11] Stop parsing when `done: true` is received
  - [x] 4.13 [depends on: 4.12] Add Godoc documentation for OllamaProvider struct and methods
  - [x] 4.14 [depends on: 4.13] Verify OllamaProvider implements LLMProvider interface

- [x] 5.0 Update run command to use provider abstraction
  - [x] 5.1 [depends on: 3.0, 4.0] Import `internal/provider` package in `cmd/run.go`
  - [x] 5.2 [depends on: 5.1] Import `internal/config` package in `cmd/run.go`
  - [x] 5.3 [depends on: 5.2] Update `executeInstruction()` to call `config.Load()` to get config
  - [x] 5.4 [depends on: 5.3] Add provider selection logic: if config.Provider == "ollama" use Ollama, else Anthropic
  - [x] 5.5 [depends on: 5.4] Instantiate Anthropic provider with model from config.Anthropic.Model
  - [x] 5.6 [depends on: 5.5] Instantiate Ollama provider with endpoint and model from config.Ollama
  - [x] 5.7 [depends on: 5.6] Add fallback defaults if config sections are missing
  - [x] 5.8 [depends on: 5.7] Replace `aiClient.Complete()` call with `provider.StreamCompletion()`
  - [x] 5.9 [depends on: 5.8] Pass empty string for systemPrompt and finalPrompt for userMessage
  - [x] 5.10 [depends on: 5.9] Remove direct dependency on `internal/ai` package
  - [x] 5.11 [depends on: 5.10] Test with Anthropic provider (default config)
  - [x] 5.12 [depends on: 5.11] Test with Ollama provider (set `provider: ollama` in config.yaml)
  - [x] 5.13 [depends on: 5.12] Test error handling when Ollama is not running
