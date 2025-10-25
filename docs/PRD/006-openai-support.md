# PRD: OpenAI Provider Support

## Introduction/Overview

This feature adds OpenAI as a third LLM provider option for Gliik, alongside the existing Anthropic Claude and Ollama integrations. Users will be able to configure Gliik to use OpenAI's GPT models for executing instructions, maintaining the same user experience and workflow as other providers.

The problem this solves: Users who prefer or have access to OpenAI's API (GPT-4o, GPT-4o-mini, etc.) currently cannot use Gliik with these models. By adding OpenAI support, we expand Gliik's flexibility and allow users to choose the LLM provider that best fits their needs, budget, or organizational requirements.

## Goals

1. Enable users to execute Gliik instructions using OpenAI's GPT models
2. Maintain architectural consistency with existing provider integrations (Anthropic, Ollama)
3. Support configurable OpenAI API endpoints for compatibility with Azure OpenAI and OpenAI-compatible APIs
4. Provide clear error messages when API credentials are missing or invalid
5. Deliver feature parity with existing providers (streaming responses, variable substitution, etc.)

## User Stories

1. **As a Gliik user with an OpenAI API key**, I want to configure Gliik to use OpenAI as my LLM provider so that I can execute instructions using GPT models.

2. **As a developer using Azure OpenAI**, I want to configure a custom OpenAI API endpoint so that I can use Gliik with my organization's Azure OpenAI deployment.

3. **As a Gliik user**, I want to see streamed responses from OpenAI in real-time (just like Anthropic and Ollama) so that I get immediate feedback while my instruction executes.

4. **As a new OpenAI user**, I want to receive a clear, helpful error message if I forget to set my OPENAI_API_KEY environment variable so that I understand what I need to do to fix the issue.

5. **As a Gliik user**, I want to easily switch between providers (Anthropic, Ollama, OpenAI) by changing my config.yaml so that I can use the best model for each task.

## Functional Requirements

1. The system must support `provider: "openai"` as a valid option in `config.yaml`

2. The system must read the `OPENAI_API_KEY` environment variable for authentication when OpenAI provider is selected

3. The system must implement an OpenAI provider that satisfies the existing `LLMProvider` interface defined in `internal/provider/provider.go`

4. The system must support the following config.yaml structure for OpenAI configuration:
   ```yaml
   provider: "openai"
   openai:
     model: "gpt-4o-mini"  # default model
     endpoint: "https://api.openai.com/v1"  # configurable for Azure/compatible APIs
   ```

5. The system must use `gpt-4o-mini` as the default OpenAI model when not explicitly specified

6. The system must support configurable API endpoints in the `openai.endpoint` config field to enable Azure OpenAI and OpenAI-compatible API usage

7. The system must stream responses from OpenAI to stdout in real-time, matching the behavior of Anthropic and Ollama providers

8. The system must display a clear, actionable error message when `OPENAI_API_KEY` is not set, explaining:
   - That the environment variable is required
   - How to set it (e.g., `export OPENAI_API_KEY=sk-...`)

9. The system must handle OpenAI API errors gracefully (rate limits, invalid keys, network issues) with user-friendly error messages

10. The system must support variable substitution in instruction templates when using OpenAI provider (same behavior as other providers)

11. The system must preserve markdown formatting in prompts sent to OpenAI (frontmatter excluded, body only)

12. The system must send only the markdown body (not frontmatter) of instruction.md files to OpenAI's API

## Non-Goals (Out of Scope)

1. **Per-instruction provider override** - OpenAI provider selection remains global via config.yaml only (consistent with project scope)

2. **OpenAI-specific features** - Advanced OpenAI features like function calling, vision, or DALL-E integration are not included in this MVP

3. **Cost tracking** - No token usage tracking or cost estimation features

4. **Model auto-selection** - No automatic model selection based on prompt complexity or length

5. **Fine-tuned models** - Initial implementation targets base OpenAI models only; custom fine-tuned model support is out of scope

6. **Conversation history** - OpenAI integration will send single-turn prompts only (consistent with current architecture)

7. **Temperature/parameter configuration** - Beyond model and endpoint, no additional OpenAI API parameters (temperature, top_p, etc.) in MVP

## Technical Considerations

1. **Provider Interface Compliance**: The OpenAI provider must implement the existing `LLMProvider` interface without modifying it, ensuring drop-in compatibility

2. **Dependencies**: Use official OpenAI Go SDK or a minimal HTTP client approach (prefer standard library per project philosophy)

3. **Streaming Implementation**: OpenAI's streaming API uses Server-Sent Events (SSE); implementation should mirror the streaming patterns used in Anthropic provider

4. **Error Handling**: Leverage OpenAI SDK's error types if using SDK, or parse API error responses carefully if using HTTP client

5. **Configuration Validation**: Validate OpenAI configuration on initialization and fail fast with clear errors if misconfigured

6. **Environment Variable**: Follow existing pattern - check `OPENAI_API_KEY` at provider initialization time, not at config load time

7. **Testing**: Manual testing required (no automated tests per project scope), verify with actual OpenAI API calls

8. **Endpoint Validation**: Ensure custom endpoint URLs are properly validated and formatted (trailing slashes, HTTPS enforcement, etc.)

## Design Considerations

No UI changes required. Configuration follows existing pattern in `config.yaml`:

```yaml
# Example configuration
default_model: "gpt-4o-mini"
provider: "openai"

openai:
  model: "gpt-4o-mini"
  endpoint: "https://api.openai.com/v1"

# Alternative: Azure OpenAI
# openai:
#   model: "gpt-4"
#   endpoint: "https://your-resource.openai.azure.com"
```

## Success Metrics

1. **Functional Parity**: Users can execute any existing Gliik instruction with OpenAI provider and receive correct, streamed responses

2. **Error Clarity**: 100% of authentication/configuration errors provide actionable error messages that users can resolve without consulting documentation

3. **Performance**: Response streaming latency comparable to Anthropic provider (no noticeable buffering delays)

4. **Adoption**: Feature is documented in README with clear setup instructions, enabling users to switch providers easily

5. **Stability**: No regressions in existing Anthropic or Ollama provider functionality after OpenAI integration

## Open Questions

1. Should we add OpenAI provider documentation to an existing README section or create a dedicated providers guide?

2. Do we need to document recommended models for different use cases (e.g., gpt-4o for complex reasoning, gpt-4o-mini for simple tasks)?

3. Should the error message for missing OPENAI_API_KEY include a link to OpenAI's API key generation page?

4. Are there specific OpenAI API rate limits or quotas we should warn users about in documentation?

5. Should we validate the API key format (starts with sk-) before making API calls, or rely on OpenAI's error responses?

## Implementation Notes for Developer

- Start by reviewing `internal/provider/anthropic.go` as a reference implementation
- The `LLMProvider` interface is defined in `internal/provider/provider.go`
- Config parsing logic is in the config package - add OpenAI section there
- Follow existing error message patterns for consistency
- Remember: no inline comments, code must be self-documenting through clear naming
- Test with multiple scenarios: valid key, missing key, invalid key, custom endpoint, different models
- Ensure streaming writes to stdout progressively, not in chunks
