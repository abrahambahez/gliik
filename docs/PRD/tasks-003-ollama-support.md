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

- [ ] 1.0 Create provider abstraction layer
  - [ ] 1.1 Create `internal/provider/` directory
  - [ ] 1.2 Create `provider.go` with `LLMProvider` interface definition
  - [ ] 1.3 Define `StreamCompletion(systemPrompt, userMessage string) error` method signature
  - [ ] 1.4 Add Godoc documentation explaining the provider abstraction pattern

- [ ] 2.0 Update configuration to support provider selection
  - [ ] 2.1 Add `Provider string` field to `Config` struct in `internal/config/init.go`
  - [ ] 2.2 Add `yaml:"provider"` tag to Provider field
  - [ ] 2.3 Create `ValidateProvider()` method that checks value is "anthropic" or "ollama"
  - [ ] 2.4 Update `Initialize()` to set default `provider: "anthropic"` in config.yaml
  - [ ] 2.5 Add Godoc for Provider field explaining valid values

- [ ] 3.0 Refactor Anthropic client into provider implementation
  - [ ] 3.1 [depends on: 1.0] Create `internal/provider/anthropic.go` file
  - [ ] 3.2 [depends on: 3.1] Copy `Client` struct from `ai/client.go` and rename to `AnthropicProvider`
  - [ ] 3.3 [depends on: 3.2] Copy `NewClient()` function and rename to `NewAnthropicProvider()`
  - [ ] 3.4 [depends on: 3.3] Copy `Complete()` method and rename to `StreamCompletion()`
  - [ ] 3.5 [depends on: 3.4] Update `StreamCompletion()` signature to match interface (systemPrompt, userMessage)
  - [ ] 3.6 [depends on: 3.5] Update method to combine systemPrompt and userMessage into messages array
  - [ ] 3.7 [depends on: 3.6] Add Godoc documentation for `AnthropicProvider` struct and methods
  - [ ] 3.8 [depends on: 3.7] Verify AnthropicProvider implements LLMProvider interface

- [ ] 4.0 Implement Ollama provider
  - [ ] 4.1 [depends on: 1.0] Create `internal/provider/ollama.go` file
  - [ ] 4.2 [depends on: 4.1] Define `OllamaProvider` struct with `Endpoint` and `Model` fields
  - [ ] 4.3 [depends on: 4.2] Implement `NewOllamaProvider()` constructor with hardcoded defaults
  - [ ] 4.4 [depends on: 4.3] Set default endpoint to `http://localhost:11434` and model to `llama3.2`
  - [ ] 4.5 [depends on: 4.4] Implement `StreamCompletion(systemPrompt, userMessage string) error` method
  - [ ] 4.6 [depends on: 4.5] Build JSON request body with model, prompt (combined system+user), and stream:true
  - [ ] 4.7 [depends on: 4.6] Create HTTP POST request to `{endpoint}/api/generate`
  - [ ] 4.8 [depends on: 4.7] Set Content-Type header to `application/json`
  - [ ] 4.9 [depends on: 4.8] Execute request and handle connection errors with specific error message from PRD
  - [ ] 4.10 [depends on: 4.9] Implement newline-delimited JSON streaming response parser
  - [ ] 4.11 [depends on: 4.10] For each line, parse JSON and extract `response` field, write to stdout
  - [ ] 4.12 [depends on: 4.11] Stop parsing when `done: true` is received
  - [ ] 4.13 [depends on: 4.12] Add Godoc documentation for OllamaProvider struct and methods
  - [ ] 4.14 [depends on: 4.13] Verify OllamaProvider implements LLMProvider interface

- [ ] 5.0 Update run command to use provider abstraction
  - [ ] 5.1 [depends on: 3.0, 4.0] Import `internal/provider` package in `cmd/run.go`
  - [ ] 5.2 [depends on: 5.1] Import `internal/config` package in `cmd/run.go`
  - [ ] 5.3 [depends on: 5.2] Update `executeInstruction()` to call `config.Load()` to get config
  - [ ] 5.4 [depends on: 5.3] Add provider selection logic: if config.Provider == "ollama" use Ollama, else Anthropic
  - [ ] 5.5 [depends on: 5.4] Instantiate correct provider based on selection
  - [ ] 5.6 [depends on: 5.5] Replace `aiClient.Complete()` call with `provider.StreamCompletion()`
  - [ ] 5.7 [depends on: 5.6] Pass empty string for systemPrompt and finalPrompt for userMessage
  - [ ] 5.8 [depends on: 5.7] Remove direct dependency on `internal/ai` package
  - [ ] 5.9 [depends on: 5.8] Test with Anthropic provider (default config)
  - [ ] 5.10 [depends on: 5.9] Test with Ollama provider (set `provider: ollama` in config.yaml)
  - [ ] 5.11 [depends on: 5.10] Test error handling when Ollama is not running
