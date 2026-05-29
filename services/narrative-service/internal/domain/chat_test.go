package domain

import (
	"strings"
	"testing"
)

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
