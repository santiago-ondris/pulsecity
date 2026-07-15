package domain

import "testing"

func TestBuildPostTradeNarrative(t *testing.T) {
	event := TradeAcceptedEvent{
		EventMeta: EventMeta{
			EventID:       "trade-accepted-trade-1",
			GameID:        "game-1",
			OccurredAt:    "2026-11-01T00:00:00Z",
			SchemaVersion: 1,
		},
		ProposalID:              "trade-1",
		SimulatedDate:           "2026-11-01",
		RivalTeamID:             "bos",
		OutgoingPlayerID:        "player-1",
		OutgoingPlayerName:      "Adrian Vale",
		IncomingPlayerID:        "trade-1-incoming",
		IncomingPlayerName:      "Jalen Warren",
		IncomingPosition:        "PG",
		IncomingRating:          76,
		IncomingSalary:          12_000_000,
		AcceptedAdditionalAsset: "second_round_pick",
	}

	narrative := BuildPostTradeNarrative(event)

	if narrative.EventID != "post-trade-trade-1" {
		t.Fatalf("BuildPostTradeNarrative() event id = %q, want post-trade-trade-1", narrative.EventID)
	}
	if narrative.Kind != "post_trade" {
		t.Fatalf("BuildPostTradeNarrative() kind = %q, want post_trade", narrative.Kind)
	}
	if narrative.Emitter != "director_player_personnel" {
		t.Fatalf("BuildPostTradeNarrative() emitter = %q, want director_player_personnel", narrative.Emitter)
	}
	if narrative.Metadata["source_subject"] != SubjectTradeAccepted {
		t.Fatalf("BuildPostTradeNarrative() source subject = %q, want %q", narrative.Metadata["source_subject"], SubjectTradeAccepted)
	}
}
