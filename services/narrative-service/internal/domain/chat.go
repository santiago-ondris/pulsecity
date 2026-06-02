package domain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	SubjectAgentConsultationStarted = "agente.consulta_iniciada"
	SubjectAgentResponseGenerated   = "agente.respuesta_generada"
	SubjectChatMessageDelta         = "chat.message"
)

type AgentConsultationStartedEvent struct {
	EventMeta
	ConversationID  string `json:"conversation_id"`
	AgentID         string `json:"agent_id"`
	Sender          string `json:"sender"`
	Message         string `json:"message"`
	ClientMessageID string `json:"client_message_id,omitempty"`
}

type AgentChatContext struct {
	AgentID        string
	DisplayName    string
	Role           string
	Domain         string
	EmotionalState string
	Confidence     float64
	Satisfaction   float64
	Loyalty        float64
	Relationship   AgentRelationshipContext
	Decisions      []GMDecisionContext
}

type AgentRelationshipContext struct {
	Trust     float64
	Trend     string
	LastEvent string
}

type GMDecisionContext struct {
	DecisionID    string
	Kind          string
	SimulatedDate string
}

type ChatRuntimeConfig struct {
	Provider                string
	Model                   string
	OpenRouterAPIKey        string
	OpenRouterBaseURL       string
	OpenRouterAppURL        string
	OpenRouterAppTitle      string
	MaxPromptChars          int
	MaxResponseChars        int
	MaxCompletionTokens     int
	MaxTurnsPerConversation int
	RequestTimeout          time.Duration
}

type ChatGenerationRequest struct {
	Event   AgentConsultationStartedEvent
	Context AgentChatContext
	Prompt  string
	Config  ChatRuntimeConfig
}

type ChatGenerationResult struct {
	Body     string
	Metadata map[string]string
}

type ChatResponder interface {
	GenerateAgentReply(context.Context, ChatGenerationRequest) (ChatGenerationResult, error)
}

type StubChatResponder struct{}

type UnsupportedProviderResponder struct {
	Provider string
}

type openRouterChatRequest struct {
	Model               string              `json:"model"`
	Messages            []openRouterMessage `json:"messages"`
	MaxCompletionTokens int                 `json:"max_completion_tokens"`
	Temperature         float64             `json:"temperature"`
}

type openRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openRouterChatResponse struct {
	Model   string `json:"model"`
	Choices []struct {
		Message openRouterMessage `json:"message"`
	} `json:"choices"`
	Usage openRouterUsage `json:"usage"`
}

type openRouterUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type OpenRouterChatResponder struct {
	apiKey   string
	baseURL  string
	client   *http.Client
	appURL   string
	appTitle string
}

type ChatMessage struct {
	MessageID      string            `json:"message_id"`
	GameID         string            `json:"game_id"`
	ConversationID string            `json:"conversation_id"`
	AgentID        string            `json:"agent_id"`
	Sender         string            `json:"sender"`
	Body           string            `json:"body"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	CreatedAt      string            `json:"created_at"`
}

type ChatMessageEnvelope struct {
	Type           string            `json:"type"`
	Subject        string            `json:"subject"`
	GameID         string            `json:"game_id"`
	ConversationID string            `json:"conversation_id"`
	MessageID      string            `json:"message_id"`
	AgentID        string            `json:"agent_id"`
	Sender         string            `json:"sender"`
	Body           string            `json:"body"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	CreatedAt      string            `json:"created_at"`
}

func DefaultChatRuntimeConfig() ChatRuntimeConfig {
	return ChatRuntimeConfig{
		Provider:                "stub",
		Model:                   "google/gemini-2.5-flash",
		OpenRouterBaseURL:       "https://openrouter.ai/api/v1",
		OpenRouterAppTitle:      "PulseCity",
		MaxPromptChars:          6000,
		MaxResponseChars:        900,
		MaxCompletionTokens:     300,
		MaxTurnsPerConversation: 12,
		RequestTimeout:          8 * time.Second,
	}
}

func ChatRuntimeConfigFromEnv(getenv func(string) string) ChatRuntimeConfig {
	config := DefaultChatRuntimeConfig()
	config.Provider = normalizedConfigValue(getenv("LLM_PROVIDER"), config.Provider)
	config.Model = normalizedConfigValue(getenv("LLM_MODEL"), config.Model)
	config.OpenRouterAPIKey = strings.TrimSpace(getenv("OPENROUTER_API_KEY"))
	config.OpenRouterBaseURL = normalizedConfigValue(getenv("OPENROUTER_BASE_URL"), config.OpenRouterBaseURL)
	config.OpenRouterAppURL = strings.TrimSpace(getenv("OPENROUTER_APP_URL"))
	config.OpenRouterAppTitle = normalizedConfigValue(getenv("OPENROUTER_APP_TITLE"), config.OpenRouterAppTitle)
	config.MaxPromptChars = parsePositiveInt(getenv("LLM_MAX_PROMPT_CHARS"), config.MaxPromptChars)
	config.MaxResponseChars = parsePositiveInt(getenv("LLM_MAX_RESPONSE_CHARS"), config.MaxResponseChars)
	config.MaxCompletionTokens = parsePositiveInt(getenv("LLM_MAX_COMPLETION_TOKENS"), config.MaxCompletionTokens)
	config.MaxTurnsPerConversation = parsePositiveInt(getenv("LLM_MAX_TURNS_PER_CONVERSATION"), config.MaxTurnsPerConversation)
	config.RequestTimeout = time.Duration(parsePositiveInt(getenv("LLM_REQUEST_TIMEOUT_SECONDS"), int(config.RequestTimeout.Seconds()))) * time.Second
	return config
}

func NewChatResponder(config ChatRuntimeConfig) ChatResponder {
	switch strings.ToLower(strings.TrimSpace(config.Provider)) {
	case "", "stub":
		return StubChatResponder{}
	case "openrouter":
		return NewOpenRouterChatResponder(config, nil)
	default:
		return UnsupportedProviderResponder{Provider: config.Provider}
	}
}

func NewOpenRouterChatResponder(config ChatRuntimeConfig, client *http.Client) OpenRouterChatResponder {
	if client == nil {
		client = &http.Client{Timeout: config.RequestTimeout}
	}

	return OpenRouterChatResponder{
		apiKey:   strings.TrimSpace(config.OpenRouterAPIKey),
		baseURL:  strings.TrimRight(normalizedConfigValue(config.OpenRouterBaseURL, DefaultChatRuntimeConfig().OpenRouterBaseURL), "/"),
		client:   client,
		appURL:   strings.TrimSpace(config.OpenRouterAppURL),
		appTitle: normalizedConfigValue(config.OpenRouterAppTitle, "PulseCity"),
	}
}

func (StubChatResponder) GenerateAgentReply(_ context.Context, request ChatGenerationRequest) (ChatGenerationResult, error) {
	body := BuildStubAgentChatResponse(request.Event.Message, request.Context)
	return ChatGenerationResult{
		Body: truncateText(body, request.Config.MaxResponseChars),
		Metadata: map[string]string{
			"generation": "stub",
			"provider":   "stub",
		},
	}, nil
}

func (r UnsupportedProviderResponder) GenerateAgentReply(_ context.Context, _ ChatGenerationRequest) (ChatGenerationResult, error) {
	return ChatGenerationResult{}, fmt.Errorf("llm provider %q is not implemented", r.Provider)
}

func (r OpenRouterChatResponder) GenerateAgentReply(ctx context.Context, request ChatGenerationRequest) (ChatGenerationResult, error) {
	if r.apiKey == "" {
		return ChatGenerationResult{}, fmt.Errorf("openrouter api key is not configured")
	}

	body, err := json.Marshal(openRouterChatRequest{
		Model: request.Config.Model,
		Messages: []openRouterMessage{
			{Role: "system", Content: request.Prompt},
			{Role: "user", Content: strings.TrimSpace(request.Event.Message)},
		},
		MaxCompletionTokens: request.Config.MaxCompletionTokens,
		Temperature:         0.7,
	})
	if err != nil {
		return ChatGenerationResult{}, fmt.Errorf("marshal openrouter request: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, r.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return ChatGenerationResult{}, fmt.Errorf("build openrouter request: %w", err)
	}
	httpRequest.Header.Set("Authorization", "Bearer "+r.apiKey)
	httpRequest.Header.Set("Content-Type", "application/json")
	if r.appURL != "" {
		httpRequest.Header.Set("HTTP-Referer", r.appURL)
	}
	if r.appTitle != "" {
		httpRequest.Header.Set("X-OpenRouter-Title", r.appTitle)
	}

	httpResponse, err := r.client.Do(httpRequest)
	if err != nil {
		return ChatGenerationResult{}, fmt.Errorf("call openrouter: %w", err)
	}
	defer httpResponse.Body.Close()

	responseBody, err := io.ReadAll(io.LimitReader(httpResponse.Body, 1<<20))
	if err != nil {
		return ChatGenerationResult{}, fmt.Errorf("read openrouter response: %w", err)
	}
	if httpResponse.StatusCode < http.StatusOK || httpResponse.StatusCode >= http.StatusMultipleChoices {
		return ChatGenerationResult{}, fmt.Errorf("openrouter returned status %d: %s", httpResponse.StatusCode, truncateText(string(responseBody), 240))
	}

	var parsed openRouterChatResponse
	if err := json.Unmarshal(responseBody, &parsed); err != nil {
		return ChatGenerationResult{}, fmt.Errorf("decode openrouter response: %w", err)
	}
	if len(parsed.Choices) == 0 {
		return ChatGenerationResult{}, fmt.Errorf("openrouter response has no choices")
	}

	content := strings.TrimSpace(parsed.Choices[0].Message.Content)
	if content == "" {
		return ChatGenerationResult{}, fmt.Errorf("openrouter response content is empty")
	}

	return ChatGenerationResult{
		Body: truncateText(content, request.Config.MaxResponseChars),
		Metadata: map[string]string{
			"generation":        "llm",
			"provider":          "openrouter",
			"model":             emptyAs(parsed.Model, request.Config.Model),
			"prompt_tokens":     strconv.Itoa(parsed.Usage.PromptTokens),
			"completion_tokens": strconv.Itoa(parsed.Usage.CompletionTokens),
			"total_tokens":      strconv.Itoa(parsed.Usage.TotalTokens),
		},
	}, nil
}

func BuildUserChatMessage(event AgentConsultationStartedEvent, createdAt string) ChatMessage {
	return ChatMessage{
		MessageID:      event.EventID + "-gm",
		GameID:         event.GameID,
		ConversationID: event.ConversationID,
		AgentID:        event.AgentID,
		Sender:         "gm",
		Body:           strings.TrimSpace(event.Message),
		Metadata: map[string]string{
			"source_event_id": event.EventID,
			"source_subject":  SubjectAgentConsultationStarted,
		},
		CreatedAt: createdAt,
	}
}

func BuildAgentChatMessage(event AgentConsultationStartedEvent, result ChatGenerationResult, createdAt string) ChatMessage {
	metadata := result.Metadata
	if metadata == nil {
		metadata = map[string]string{}
	}
	metadata["source_event_id"] = event.EventID
	metadata["source_subject"] = SubjectAgentConsultationStarted

	return ChatMessage{
		MessageID:      "agent-response-" + uuid.NewString(),
		GameID:         event.GameID,
		ConversationID: event.ConversationID,
		AgentID:        event.AgentID,
		Sender:         "agent",
		Body:           strings.TrimSpace(result.Body),
		Metadata:       metadata,
		CreatedAt:      createdAt,
	}
}

func BuildFallbackAgentChatMessage(event AgentConsultationStartedEvent, config ChatRuntimeConfig, createdAt string) ChatMessage {
	return BuildAgentChatMessage(event, ChatGenerationResult{
		Body: truncateText(
			"No pude generar una respuesta completa ahora. Te dejo una lectura corta: mantengamos la decision dentro de mi area y volve a intentarlo en un momento.",
			config.MaxResponseChars,
		),
		Metadata: map[string]string{
			"generation": "fallback",
			"provider":   normalizedConfigValue(config.Provider, "stub"),
		},
	}, createdAt)
}

func BuildTurnLimitAgentChatMessage(event AgentConsultationStartedEvent, config ChatRuntimeConfig, createdAt string) ChatMessage {
	return BuildAgentChatMessage(event, ChatGenerationResult{
		Body: fmt.Sprintf("Esta conversacion llego al limite de %d turnos. Abrime un chat nuevo si queres seguir.", config.MaxTurnsPerConversation),
		Metadata: map[string]string{
			"generation": "turn_limit",
			"provider":   normalizedConfigValue(config.Provider, "stub"),
		},
	}, createdAt)
}

func ChatEnvelopeFromMessage(message ChatMessage) ChatMessageEnvelope {
	return ChatMessageEnvelope{
		Type:           SubjectChatMessageDelta,
		Subject:        SubjectAgentResponseGenerated,
		GameID:         message.GameID,
		ConversationID: message.ConversationID,
		MessageID:      message.MessageID,
		AgentID:        message.AgentID,
		Sender:         message.Sender,
		Body:           message.Body,
		Metadata:       message.Metadata,
		CreatedAt:      message.CreatedAt,
	}
}

func BuildAgentChatPrompt(event AgentConsultationStartedEvent, context AgentChatContext, config ChatRuntimeConfig) string {
	decisions := "sin decisiones recientes"
	if len(context.Decisions) > 0 {
		parts := make([]string, 0, len(context.Decisions))
		for _, decision := range context.Decisions {
			parts = append(parts, fmt.Sprintf("%s en %s", decision.Kind, decision.SimulatedDate))
		}
		decisions = strings.Join(parts, "; ")
	}

	relationship := "sin relacion GM directa registrada"
	if context.Relationship.Trend != "" || context.Relationship.LastEvent != "" {
		relationship = fmt.Sprintf(
			"trust %.2f, tendencia %s, ultimo evento: %s",
			context.Relationship.Trust,
			emptyAs(context.Relationship.Trend, "estable"),
			emptyAs(context.Relationship.LastEvent, "sin evento"),
		)
	}

	prompt := fmt.Sprintf(`Sistema:
Sos un agente vivo de PulseCity. Responde en espanol rioplatense, en primera persona, con tono profesional y concreto.
Tu dominio es estricto: %s.
Si el GM pregunta algo fuera de tu dominio, decilo en una linea y redirigi a tu area.
No inventes datos que no esten en el contexto. No cambies estado de juego. No prometas acciones de otros servicios.
Mantene la respuesta breve: 1 a 3 parrafos cortos.

Contexto del agente:
- agent_id: %s
- nombre: %s
- rol: %s
- estado emocional: %s
- confianza: %.2f
- satisfaccion: %.2f
- lealtad: %.2f
- relacion con GM: %s
- decisiones recientes del GM: %s

Mensaje del GM:
%s`,
		emptyAs(context.Domain, "estado general de la franquicia"),
		event.AgentID,
		emptyAs(context.DisplayName, event.AgentID),
		emptyAs(context.Role, "agente"),
		emptyAs(context.EmotionalState, "unknown"),
		context.Confidence,
		context.Satisfaction,
		context.Loyalty,
		relationship,
		decisions,
		strings.TrimSpace(event.Message),
	)

	return truncateText(prompt, config.MaxPromptChars)
}

func BuildStubAgentChatResponse(prompt string, context AgentChatContext) string {
	displayName := strings.TrimSpace(context.DisplayName)
	if displayName == "" {
		displayName = context.AgentID
	}
	role := strings.TrimSpace(context.Role)
	if role == "" {
		role = "agente"
	}
	domain := strings.TrimSpace(context.Domain)
	if domain == "" {
		domain = "mi area"
	}

	relationLine := "Todavia no tengo una lectura relacional fuerte con vos."
	if context.Relationship.Trend != "" || context.Relationship.LastEvent != "" {
		relationLine = fmt.Sprintf(
			"Nuestra relacion viene %s con confianza %.2f.",
			emptyAs(context.Relationship.Trend, "estable"),
			context.Relationship.Trust,
		)
	}

	decisionLine := "No veo decisiones recientes del GM que cambien esta respuesta."
	if len(context.Decisions) > 0 {
		latest := context.Decisions[0]
		decisionLine = fmt.Sprintf("Tengo presente tu ultima decision: %s en %s.", latest.Kind, latest.SimulatedDate)
	}

	request := strings.TrimSpace(prompt)
	if len(request) > 120 {
		request = request[:120] + "..."
	}

	return fmt.Sprintf(
		"[stub M3.6] Soy %s (%s). Puedo responder desde %s. %s %s Sobre lo que preguntaste: %q, mi lectura inicial es mantener el foco y pedirte una decision concreta antes de mover el sistema.",
		displayName,
		role,
		domain,
		relationLine,
		decisionLine,
		request,
	)
}

func truncateText(value string, limit int) string {
	if limit <= 0 || len(value) <= limit {
		return value
	}
	if limit <= 3 {
		return value[:limit]
	}

	return strings.TrimSpace(value[:limit-3]) + "..."
}

func normalizedConfigValue(value, fallback string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return fallback
	}

	return trimmed
}

func parsePositiveInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || parsed <= 0 {
		return fallback
	}

	return parsed
}

func emptyAs(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}

	return value
}
