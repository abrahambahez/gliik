## Relevant Files

- `internal/config/init.go` - Contains Config struct and provider-specific config structs. Added OpenAIConfig struct and updated Config struct to include openai field. ✅ Modified
- `internal/provider/openai.go` - Implements OpenAI provider with LLMProvider interface, similar to anthropic.go structure. ✅ Created
- `internal/provider/provider.go` - Defines LLMProvider interface (no changes needed, reference only).
- `internal/provider/anthropic.go` - Reference implementation for understanding streaming patterns and error handling.
- `internal/provider/ollama.go` - Reference implementation for understanding different streaming approach.
- `cmd/run.go` - Contains provider instantiation logic in executeInstruction function. Added OpenAI provider case alongside anthropic and ollama. ✅ Modified
- `README.md` - Project documentation. Added OpenAI provider setup instructions, configuration examples, and environment variables. ✅ Modified
- `CLAUDE.md` - Developer guidance for Claude Code. Added OpenAI to technical stack and config.yaml schema. ✅ Modified

### Notes

- Follow the existing architectural pattern: Anthropic uses non-streaming API but prints immediately, Ollama uses true streaming with Server-Sent Events
- OpenAI uses Server-Sent Events for streaming, similar pattern to Ollama
- Use standard library HTTP client (no SDK) to minimize dependencies per project philosophy
- All provider error messages should be clear and actionable for end users
- Variable substitution is handled upstream, so providers only receive final prompt text
- No automated tests are required per project scope, but manual testing scenarios should be comprehensive

## Tasks

- [x] 1.0 Add OpenAI configuration structure to config package
  - [x] 1.1 Add OpenAIConfig struct to internal/config/init.go with Model and Endpoint fields
  - [x] 1.2 Add OpenAI field to Config struct with yaml tag `yaml:"openai"`
  - [x] 1.3 Update Initialize function to include default OpenAI configuration in config template (model: "gpt-4o-mini", endpoint: "https://api.openai.com/v1")

- [x] 2.0 Implement OpenAI provider with LLMProvider interface
  - [x] 2.1 [depends on: 1.0] Create new file internal/provider/openai.go with package declaration and imports
  - [x] 2.2 Define OpenAIProvider struct with APIKey, Model, and Endpoint fields
  - [x] 2.3 Implement NewOpenAIProvider constructor function that reads OPENAI_API_KEY environment variable and returns clear error if not set
  - [x] 2.4 Add compile-time interface check: var _ LLMProvider = (*OpenAIProvider)(nil)
  - [x] 2.5 Implement StreamCompletion method that constructs OpenAI API request payload (messages format with user role)
  - [x] 2.6 [depends on: 2.5] Implement streaming response handling using bufio.Scanner to read Server-Sent Events line by line
  - [x] 2.7 [depends on: 2.6] Parse SSE data lines containing JSON chunks and extract delta content to print to stdout in real-time
  - [x] 2.8 Add error handling for missing API key with clear message explaining how to set OPENAI_API_KEY
  - [x] 2.9 Add error handling for HTTP errors (rate limits, auth failures, network issues) with user-friendly messages
  - [x] 2.10 Validate and format endpoint URL to ensure proper API path construction

- [x] 3.0 Integrate OpenAI provider into command execution flow
  - [x] 3.1 [depends on: 2.0] Add OpenAI provider case to executeInstruction function in cmd/run.go after ollama case
  - [x] 3.2 [depends on: 3.1] Read cfg.OpenAI.Endpoint with fallback to default "https://api.openai.com/v1" if empty
  - [x] 3.3 [depends on: 3.1] Read cfg.OpenAI.Model with fallback to default "gpt-4o-mini" if empty
  - [x] 3.4 [depends on: 3.1] Call provider.NewOpenAIProvider with endpoint and model, handle error appropriately

- [x] 4.0 Add provider validation for OpenAI option
  - [x] 4.1 [depends on: 1.0] Update ValidateProvider function in internal/config/init.go to accept "openai" as valid provider alongside "anthropic" and "ollama"
  - [x] 4.2 Update error message in ValidateProvider to list all three valid options: 'anthropic', 'ollama', or 'openai'

- [x] 5.0 Manual testing and verification
  - [x] 5.1 [depends on: 3.0, 4.0] Test with valid OPENAI_API_KEY and default configuration - verify streaming works and response is complete
  - [x] 5.2 [depends on: 3.0, 4.0] Test without OPENAI_API_KEY environment variable - verify clear, actionable error message appears
  - [x] 5.3 [depends on: 3.0, 4.0] Test with invalid OPENAI_API_KEY - verify OpenAI API error is caught and displayed clearly
  - [x] 5.4 [depends on: 3.0, 4.0] Test with custom endpoint configuration (e.g., Azure OpenAI format) - verify endpoint is properly used
  - [x] 5.5 [depends on: 3.0, 4.0] Test with different models (gpt-4o, gpt-3.5-turbo) - verify model selection works correctly
  - [x] 5.6 [depends on: 3.0, 4.0] Test variable substitution works correctly with OpenAI provider using an existing instruction
  - [x] 5.7 [depends on: 3.0, 4.0] Test that Anthropic and Ollama providers still work correctly (no regressions)
  - [x] 5.8 [depends on: 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 5.7] Document any issues found and verify all test scenarios pass

## Test Results Summary

All test scenarios **PASSED** with no issues found:

- ✅ Valid API key with streaming works correctly
- ✅ Missing API key shows clear error message
- ✅ Invalid API key is properly caught and reported
- ✅ Custom endpoint configuration works as expected
- ✅ Model selection (gpt-4o, gpt-3.5-turbo, gpt-4o-mini) functions correctly
- ✅ Variable substitution works with OpenAI provider
- ✅ No regressions in Anthropic or Ollama providers

**Status**: OpenAI support is fully functional and ready for production use.
