package domain

import "testing"

func TestAgentPatchFromStateChanged(t *testing.T) {
	event := AgentStateChangedEvent{
		EventMeta: EventMeta{
			EventID:       "agent-state-game-1-match-1-owner",
			GameID:        "game-1",
			OccurredAt:    "2026-05-24T00:00:00Z",
			SchemaVersion: 1,
		},
		SimulatedDate: "2026-10-22",
		AgentID:       "owner",
		SourceEventID: "match-finished-match-1",
		SourceSubject: "partido.terminado",
		Mood:          "calm",
		State: map[string]float64{
			"sporting_trust": 0.06,
		},
		Summary: "El owner ajusta su confianza.",
	}

	patch := AgentPatchFromStateChanged(event)

	if patch.Type != SubjectAgentPatchDelta {
		t.Fatalf("expected type %q, got %q", SubjectAgentPatchDelta, patch.Type)
	}
	if patch.GameID != event.GameID {
		t.Fatalf("expected game id %q, got %q", event.GameID, patch.GameID)
	}
	if patch.AgentID != event.AgentID {
		t.Fatalf("expected agent id %q, got %q", event.AgentID, patch.AgentID)
	}
	if patch.Patch.State["sporting_trust"] != 0.06 {
		t.Fatalf("expected sporting trust in patch")
	}
	if patch.Patch.SourceSubject != SubjectMatchFinished {
		t.Fatalf("expected source subject %q, got %q", SubjectMatchFinished, patch.Patch.SourceSubject)
	}
}

func TestRelationsPatchFromRelationshipChanged(t *testing.T) {
	event := AgentRelationshipChangedEvent{
		EventMeta: EventMeta{
			EventID:       "agent-relationship-game-1-match-1",
			GameID:        "game-1",
			OccurredAt:    "2026-05-25T00:00:00Z",
			SchemaVersion: 1,
		},
		SimulatedDate: "2026-10-22",
		AgentAID:      "head_coach",
		AgentBID:      "head_analytics",
		Trust:         -0.2,
		Trend:         "deteriorating",
		LastEvent:     "La derrota reabre la tension entre datos y cancha.",
		ShortHistory:  []string{"La derrota reabre la tension entre datos y cancha."},
		SourceEventID: "match-finished-match-1",
		SourceSubject: SubjectMatchFinished,
	}

	patch := RelationsPatchFromRelationshipChanged(event)

	if patch.Type != SubjectRelationsPatchDelta {
		t.Fatalf("expected type %q, got %q", SubjectRelationsPatchDelta, patch.Type)
	}
	if patch.GameID != event.GameID {
		t.Fatalf("expected game id %q, got %q", event.GameID, patch.GameID)
	}
	if len(patch.Patch.Relationships) != 1 {
		t.Fatalf("expected one relationship patch")
	}
	if patch.Patch.Relationships[0].AgentAID != "head_coach" {
		t.Fatalf("expected agent a in relationship patch")
	}
}
