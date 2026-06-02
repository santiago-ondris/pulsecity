package domain

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestBuildStubAgentChatResponseUsesRealContext(t *testing.T) {
	response := BuildStubAgentChatResponse("Como ves el vestuario?", AgentChatContext{
		AgentID:     "head_coach",
		DisplayName: "Mara Ellison",
		Role:        "Head Coach",
		Domain:      "rotacion, sistema y vestuario",
		Relationship: AgentRelationshipContext{
			Trust: -0.15,
			Trend: "tensa",
		},
		Decisions: []GMDecisionContext{
			{Kind: "owner_intro_response", SimulatedDate: "2026-10-01"},
		},
	})

	for _, expected := range []string{
		"Mara Ellison",
		"Head Coach",
		"rotacion, sistema y vestuario",
		"confianza -0.15",
		"owner_intro_response",
	} {
		if !strings.Contains(response, expected) {
			t.Fatalf("expected response to contain %q, got %q", expected, response)
		}
	}
}

func TestChatEnvelopeFromMessageUsesDeltaType(t *testing.T) {
	message := ChatMessage{
		MessageID:      "message-1",
		GameID:         "game-1",
		ConversationID: "conversation-1",
		AgentID:        "owner",
		Sender:         "agent",
		Body:           "Respuesta.",
		CreatedAt:      "2026-05-29T00:00:00Z",
	}

	envelope := ChatEnvelopeFromMessage(message)

	if envelope.Type != SubjectChatMessageDelta {
		t.Fatalf("expected delta type %q, got %q", SubjectChatMessageDelta, envelope.Type)
	}
	if envelope.Subject != SubjectAgentResponseGenerated {
		t.Fatalf("expected subject %q, got %q", SubjectAgentResponseGenerated, envelope.Subject)
	}
	if envelope.MessageID != message.MessageID {
		t.Fatalf("expected message id to be preserved")
	}
}

func TestChatRuntimeConfigFromEnvOverridesCaps(t *testing.T) {
	values := map[string]string{
		"LLM_PROVIDER":                   "stub",
		"LLM_MODEL":                      "google/gemini-2.5-flash-lite",
		"OPENROUTER_API_KEY":             "test-key",
		"OPENROUTER_BASE_URL":            "https://example.test/api/v1",
		"OPENROUTER_APP_URL":             "https://pulsecity.local",
		"OPENROUTER_APP_TITLE":           "PulseCity Dev",
		"LLM_MAX_PROMPT_CHARS":           "1200",
		"LLM_MAX_RESPONSE_CHARS":         "240",
		"LLM_MAX_COMPLETION_TOKENS":      "80",
		"LLM_MAX_TURNS_PER_CONVERSATION": "4",
		"LLM_REQUEST_TIMEOUT_SECONDS":    "3",
	}

	config := ChatRuntimeConfigFromEnv(func(key string) string {
		return values[key]
	})

	if config.Provider != "stub" {
		t.Fatalf("ChatRuntimeConfigFromEnv provider = %q, want %q", config.Provider, "stub")
	}
	if config.Model != "google/gemini-2.5-flash-lite" {
		t.Fatalf("ChatRuntimeConfigFromEnv Model = %q, want %q", config.Model, "google/gemini-2.5-flash-lite")
	}
	if config.OpenRouterAPIKey != "test-key" {
		t.Fatalf("ChatRuntimeConfigFromEnv OpenRouterAPIKey = %q, want %q", config.OpenRouterAPIKey, "test-key")
	}
	if config.OpenRouterBaseURL != "https://example.test/api/v1" {
		t.Fatalf("ChatRuntimeConfigFromEnv OpenRouterBaseURL = %q, want %q", config.OpenRouterBaseURL, "https://example.test/api/v1")
	}
	if config.OpenRouterAppURL != "https://pulsecity.local" {
		t.Fatalf("ChatRuntimeConfigFromEnv OpenRouterAppURL = %q, want %q", config.OpenRouterAppURL, "https://pulsecity.local")
	}
	if config.OpenRouterAppTitle != "PulseCity Dev" {
		t.Fatalf("ChatRuntimeConfigFromEnv OpenRouterAppTitle = %q, want %q", config.OpenRouterAppTitle, "PulseCity Dev")
	}
	if config.MaxPromptChars != 1200 {
		t.Fatalf("ChatRuntimeConfigFromEnv MaxPromptChars = %d, want %d", config.MaxPromptChars, 1200)
	}
	if config.MaxResponseChars != 240 {
		t.Fatalf("ChatRuntimeConfigFromEnv MaxResponseChars = %d, want %d", config.MaxResponseChars, 240)
	}
	if config.MaxCompletionTokens != 80 {
		t.Fatalf("ChatRuntimeConfigFromEnv MaxCompletionTokens = %d, want %d", config.MaxCompletionTokens, 80)
	}
	if config.MaxTurnsPerConversation != 4 {
		t.Fatalf("ChatRuntimeConfigFromEnv MaxTurnsPerConversation = %d, want %d", config.MaxTurnsPerConversation, 4)
	}
	if config.RequestTimeout.Seconds() != 3 {
		t.Fatalf("ChatRuntimeConfigFromEnv RequestTimeout = %s, want 3s", config.RequestTimeout)
	}
}

func TestBuildAgentChatPromptIncludesDomainGuardrails(t *testing.T) {
	event := AgentConsultationStartedEvent{
		AgentID: "cfo",
		Message: "Que hacemos con la rotacion?",
	}
	context := AgentChatContext{
		AgentID:        "cfo",
		DisplayName:    "Iris Calder",
		Role:           "CFO",
		Domain:         "finanzas, presupuesto, cap y riesgo economico",
		EmotionalState: "alert",
		Relationship: AgentRelationshipContext{
			Trust:     0.22,
			Trend:     "estable",
			LastEvent: "El GM respeto una alerta financiera.",
		},
		Decisions: []GMDecisionContext{
			{Kind: "owner_intro_response", SimulatedDate: "2026-10-01"},
		},
	}

	prompt := BuildAgentChatPrompt(event, context, DefaultChatRuntimeConfig())

	for _, expected := range []string{
		"Tu dominio es estricto: finanzas, presupuesto, cap y riesgo economico.",
		"Si el GM pregunta algo fuera de tu dominio",
		"Iris Calder",
		"owner_intro_response en 2026-10-01",
		"Que hacemos con la rotacion?",
	} {
		if !strings.Contains(prompt, expected) {
			t.Fatalf("BuildAgentChatPrompt() missing %q in prompt %q", expected, prompt)
		}
	}
}

func TestNewChatResponderReturnsUnsupportedProviderResponder(t *testing.T) {
	responder := NewChatResponder(ChatRuntimeConfig{Provider: "openai"})

	_, err := responder.GenerateAgentReply(context.Background(), ChatGenerationRequest{})
	if err == nil {
		t.Fatal("NewChatResponder(openai) GenerateAgentReply error = nil, want unsupported provider error")
	}
}

func TestOpenRouterChatResponderGeneratesReply(t *testing.T) {
	var gotRequest openRouterChatRequest
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/api/v1/chat/completions" {
			t.Errorf("OpenRouterChatResponder path = %q, want %q", r.URL.Path, "/api/v1/chat/completions")
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("OpenRouterChatResponder Authorization = %q, want bearer test key", got)
		}
		if got := r.Header.Get("HTTP-Referer"); got != "https://pulsecity.local" {
			t.Errorf("OpenRouterChatResponder HTTP-Referer = %q, want %q", got, "https://pulsecity.local")
		}
		if got := r.Header.Get("X-OpenRouter-Title"); got != "PulseCity Dev" {
			t.Errorf("OpenRouterChatResponder X-OpenRouter-Title = %q, want %q", got, "PulseCity Dev")
		}
		if err := json.NewDecoder(r.Body).Decode(&gotRequest); err != nil {
			t.Fatalf("OpenRouterChatResponder request decode error = %v", err)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body: io.NopCloser(strings.NewReader(`{
			"model": "google/gemini-2.5-flash",
			"choices": [
				{"message": {"role": "assistant", "content": "Desde mi area, el riesgo financiero es manejable."}}
			],
			"usage": {"prompt_tokens": 120, "completion_tokens": 14, "total_tokens": 134}
		}`)),
		}, nil
	})}

	config := DefaultChatRuntimeConfig()
	config.Provider = "openrouter"
	config.OpenRouterAPIKey = "test-key"
	config.OpenRouterBaseURL = "https://openrouter.test/api/v1"
	config.OpenRouterAppURL = "https://pulsecity.local"
	config.OpenRouterAppTitle = "PulseCity Dev"
	config.MaxCompletionTokens = 90
	responder := NewOpenRouterChatResponder(config, client)

	result, err := responder.GenerateAgentReply(context.Background(), ChatGenerationRequest{
		Event: AgentConsultationStartedEvent{
			AgentID: "cfo",
			Message: "Podemos absorber este salario?",
		},
		Prompt: "Sistema: responde como CFO.",
		Config: config,
	})
	if err != nil {
		t.Fatalf("OpenRouterChatResponder.GenerateAgentReply() error = %v, want nil", err)
	}

	if gotRequest.Model != "google/gemini-2.5-flash" {
		t.Fatalf("OpenRouterChatResponder request model = %q, want %q", gotRequest.Model, "google/gemini-2.5-flash")
	}
	if gotRequest.MaxCompletionTokens != 90 {
		t.Fatalf("OpenRouterChatResponder request MaxCompletionTokens = %d, want %d", gotRequest.MaxCompletionTokens, 90)
	}
	if result.Body != "Desde mi area, el riesgo financiero es manejable." {
		t.Fatalf("OpenRouterChatResponder.GenerateAgentReply() body = %q, want provider content", result.Body)
	}
	if result.Metadata["provider"] != "openrouter" {
		t.Fatalf("OpenRouterChatResponder.GenerateAgentReply() provider = %q, want openrouter", result.Metadata["provider"])
	}
	if result.Metadata["total_tokens"] != "134" {
		t.Fatalf("OpenRouterChatResponder.GenerateAgentReply() total_tokens = %q, want 134", result.Metadata["total_tokens"])
	}
}

func TestOpenRouterChatResponderReturnsErrorOnProviderFailure(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusTooManyRequests,
			Body:       io.NopCloser(strings.NewReader("rate limited")),
		}, nil
	})}

	config := DefaultChatRuntimeConfig()
	config.OpenRouterAPIKey = "test-key"
	config.OpenRouterBaseURL = "https://openrouter.test/api/v1"
	responder := NewOpenRouterChatResponder(config, client)

	_, err := responder.GenerateAgentReply(context.Background(), ChatGenerationRequest{Config: config})
	if err == nil {
		t.Fatal("OpenRouterChatResponder.GenerateAgentReply() error = nil, want provider failure error")
	}
}

func TestBuildAgentChatPromptRespectsPromptCap(t *testing.T) {
	config := DefaultChatRuntimeConfig()
	config.MaxPromptChars = 80

	prompt := BuildAgentChatPrompt(AgentConsultationStartedEvent{
		AgentID: "owner",
		Message: strings.Repeat("mensaje largo ", 30),
	}, AgentChatContext{AgentID: "owner"}, config)

	if len(prompt) > config.MaxPromptChars {
		t.Fatalf("BuildAgentChatPrompt() length = %d, want <= %d", len(prompt), config.MaxPromptChars)
	}
}

func TestBuildTurnLimitAgentChatMessageMarksMetadata(t *testing.T) {
	config := DefaultChatRuntimeConfig()
	config.MaxTurnsPerConversation = 2
	message := BuildTurnLimitAgentChatMessage(AgentConsultationStartedEvent{
		EventMeta:      EventMeta{EventID: "event-1", GameID: "game-1"},
		ConversationID: "conversation-1",
		AgentID:        "owner",
	}, config, "2026-05-29T00:00:00Z")

	if message.Metadata["generation"] != "turn_limit" {
		t.Fatalf("BuildTurnLimitAgentChatMessage() generation = %q, want %q", message.Metadata["generation"], "turn_limit")
	}
	if !strings.Contains(message.Body, "2 turnos") {
		t.Fatalf("BuildTurnLimitAgentChatMessage() body = %q, want turn cap", message.Body)
	}
}
