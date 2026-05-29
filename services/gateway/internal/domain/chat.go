package domain

const (
	SubjectAgentConsultationStarted = "agente.consulta_iniciada"
	SubjectAgentResponseGenerated   = "agente.respuesta_generada"
	SubjectChatMessageDelta         = "chat.message"
)

type AgentChatRequest struct {
	AgentID        string `json:"agent_id"`
	Message        string `json:"message"`
	ConversationID string `json:"conversation_id,omitempty"`
}

type AgentConsultationStartedEvent struct {
	EventMeta
	ConversationID  string `json:"conversation_id"`
	AgentID         string `json:"agent_id"`
	Sender          string `json:"sender"`
	Message         string `json:"message"`
	ClientMessageID string `json:"client_message_id,omitempty"`
}

type AgentChatAcceptedResponse struct {
	ConversationID string `json:"conversation_id"`
	RequestEventID string `json:"request_event_id"`
	Status         string `json:"status"`
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
