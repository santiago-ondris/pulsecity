package domain

import (
	"fmt"
	"strings"

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

func BuildStubAgentChatMessage(event AgentConsultationStartedEvent, context AgentChatContext, createdAt string) ChatMessage {
	return ChatMessage{
		MessageID:      "agent-response-" + uuid.NewString(),
		GameID:         event.GameID,
		ConversationID: event.ConversationID,
		AgentID:        event.AgentID,
		Sender:         "agent",
		Body:           BuildStubAgentChatResponse(event.Message, context),
		Metadata: map[string]string{
			"generation":      "stub",
			"source_event_id": event.EventID,
			"source_subject":  SubjectAgentConsultationStarted,
		},
		CreatedAt: createdAt,
	}
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
		"[stub M3.5] Soy %s (%s). Puedo responder desde %s. %s %s Sobre lo que preguntaste: %q, mi lectura inicial es mantener el foco y pedirte una decision concreta antes de mover el sistema.",
		displayName,
		role,
		domain,
		relationLine,
		decisionLine,
		request,
	)
}

func emptyAs(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}

	return value
}
