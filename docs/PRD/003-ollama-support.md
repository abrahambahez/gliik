# PRD: Ollama Support

## Introduction/Overview

Currently, Gliik only supports Anthropic's Claude API for executing AI instructions. This feature adds support for Ollama, a local LLM runtime, allowing users to run instructions using locally-hosted open-source models instead of cloud-based APIs.

**Problem:** Users who want to use local models for privacy, cost, or offline usage cannot use Gliik.

**Goal:** Enable Gliik to work with Ollama as an alternative LLM provider alongside Anthropic.

## Goals

1. Allow users to choose between Anthropic and Ollama via configuration
2. Maintain existing Anthropic functionality without breaking changes
3. Support streaming responses from Ollama (matching current UX)
4. Provide clear error messages when Ollama is unavailable
5. Keep implementation minimal and aligned with MVP philosophy

## User Stories

1. **As a privacy-conscious developer**, I want to use Gliik with local Ollama models so that my prompts and data never leave my machine.

2. **As a cost-conscious user**, I want to use free local models via Ollama so that I don't incur API costs for every instruction execution.

3. **As an offline developer**, I want to use Gliik without internet connectivity so that I can work in environments with restricted network access.

4. **As a current Anthropic user**, I want to switch providers without changing my existing instructions so that I can experiment with Ollama without disruption.

## Functional Requirements

### Configuration (config.yaml)

1. The system must add a `provider` field to `config.yaml` with valid values: `"anthropic"` or `"ollama"`
2. The system must default to `"anthropic"` if no provider is specified (backward compatibility)
3. When `provider` is set to `"ollama"`, the system must use Ollama endpoint instead of Anthropic API
4. The system must validate the provider value and error if it's not `"anthropic"` or `"ollama"`

### Ollama Integration

5. The system must send requests to Ollama's `/api/generate` endpoint
6. The system must use the endpoint `http://localhost:11434` for Ollama (hardcoded default)
7. The system must use a default Ollama model `llama3.2` (hardcoded default)
8. The system must format prompts for Ollama as a single combined string (system + user content)
9. The system must handle Ollama's newline-delimited JSON streaming responses
10. The system must stream Ollama responses to stdout in real-time (matching current Anthropic behavior)

### Error Handling

11. If Ollama provider is selected but the endpoint is unreachable, the system must display:
    ```
    Error: Cannot connect to Ollama

    Make sure Ollama is running:
      ollama serve

    Default endpoint: http://localhost:11434
    ```
12. The system must not fall back to Anthropic automatically (fail fast)
13. The system must preserve existing Anthropic error handling without changes

### Code Architecture

14. The system must create an abstraction layer for LLM providers (interface or similar pattern)
15. The system must refactor existing Anthropic code into a provider implementation
16. The system must implement a new Ollama provider following the same interface
17. The system must use only Go standard library for HTTP calls to Ollama (no new dependencies)

## Non-Goals (Out of Scope)

1. **Custom Ollama endpoints:** Hardcode `http://localhost:11434`, no user configuration
2. **Custom models per instruction:** Use single default model for all instructions
3. **Model validation:** No checks for whether specified model exists in Ollama
4. **Ollama installation detection:** Assume user has installed Ollama if they configure it
5. **Multiple provider support per execution:** Cannot mix providers in single run
6. **Per-instruction provider override:** Cannot specify provider in meta.yaml
7. **Anthropic and Ollama automatic fallback:** No automatic switching between providers
8. **Model management:** No `gliik` commands to list/pull/manage Ollama models

## Technical Considerations

### Provider Interface

Create an interface similar to:
```go
type LLMProvider interface {
    StreamCompletion(systemPrompt string, userMessage string) error
}
```

### Ollama API Request Format

Use Ollama's `/api/generate` endpoint with:
```json
{
  "model": "llama3.2",
  "prompt": "system: <system_prompt>\n\nuser: <user_message>",
  "stream": true
}
```

### Ollama Streaming Response

Parse newline-delimited JSON, each line contains:
```json
{"model":"llama3.2","response":"text chunk","done":false}
```

### Provider Selection Logic

In execute command:
1. Load config.yaml
2. Read `provider` field (default: "anthropic")
3. Instantiate appropriate provider
4. Call `StreamCompletion()` method

### Files to Modify

- `internal/config/config.go` - Add provider field
- `cmd/execute.go` or similar - Provider selection logic
- New file: `internal/provider/provider.go` - Interface definition
- New file: `internal/provider/anthropic.go` - Refactored Anthropic code
- New file: `internal/provider/ollama.go` - New Ollama implementation

## Success Metrics

1. Users can successfully execute instructions using Ollama by setting `provider: ollama` in config
2. Existing Anthropic users experience zero breaking changes
3. Ollama streaming output appears in real-time (no perceptible delay vs Anthropic)
4. Error messages for unreachable Ollama are clear and actionable
5. Implementation adds <200 lines of code (excluding refactoring)

## Open Questions

1. ~~Should we support custom Ollama endpoints?~~ → No (out of scope for MVP)
2. ~~Should we support per-instruction model selection?~~ → No (global default only)
3. ~~What should happen if Ollama is unreachable?~~ → Show error, don't fallback
4. Should we validate that the configured model exists in Ollama? → No (out of scope)
5. Should we add a `gliik config set-provider <name>` helper command? → Optional nice-to-have

## Implementation Notes for Developer

- Follow existing patterns in the codebase (see CLAUDE.md and CODE_RULES.md)
- Use long, descriptive function names
- No inline comments, only Godoc
- Prefer standard library over dependencies
- Test manually with: `ollama serve` + `ollama pull llama3.2` before running Gliik
