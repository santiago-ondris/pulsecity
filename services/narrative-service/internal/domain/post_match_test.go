package domain

import (
	"strings"
	"testing"
)

func TestBuildPostMatchNarrativeForHomeWin(t *testing.T) {
	event := sampleMatchFinished(true, true, 118, 104)

	narrative := BuildPostMatchNarrative(event)

	if narrative.Type != "narrative.event" {
		t.Fatalf("expected narrative.event type, got %q", narrative.Type)
	}
	if narrative.Subject != SubjectNarrativeEventGenerated {
		t.Fatalf("expected subject %q, got %q", SubjectNarrativeEventGenerated, narrative.Subject)
	}
	if narrative.Kind != "post_match" {
		t.Fatalf("expected post_match kind, got %q", narrative.Kind)
	}
	if narrative.Emitter != "owner" {
		t.Fatalf("expected owner emitter for home win, got %q", narrative.Emitter)
	}
	if narrative.Urgency != "low" {
		t.Fatalf("expected low urgency for win, got %q", narrative.Urgency)
	}
	if narrative.Metadata["match_id"] != "match-1" {
		t.Fatalf("expected match metadata")
	}
	if len(narrative.Choices) != 1 || narrative.Choices[0].ID != "acknowledge" {
		t.Fatalf("expected acknowledge choice")
	}
}

func TestBuildPostMatchNarrativeForBlowoutLoss(t *testing.T) {
	event := sampleMatchFinished(true, false, 91, 119)

	narrative := BuildPostMatchNarrative(event)

	if narrative.Emitter != "head_coach" {
		t.Fatalf("expected head coach emitter, got %q", narrative.Emitter)
	}
	if narrative.Urgency != "high" {
		t.Fatalf("expected high urgency, got %q", narrative.Urgency)
	}
	if narrative.Metadata["margin"] != "28" {
		t.Fatalf("expected margin metadata, got %q", narrative.Metadata["margin"])
	}
}

func TestBuildPostMatchNarrativeUsesOwnTopScorerAndMoment(t *testing.T) {
	event := sampleMatchFinished(false, true, 107, 104)

	narrative := BuildPostMatchNarrative(event)

	if narrative.Emitter != "sports_psychologist" {
		t.Fatalf("expected sports psychologist for close game, got %q", narrative.Emitter)
	}
	if !strings.Contains(narrative.Body, "player-2 con 28 puntos") {
		t.Fatalf("expected own top scorer in body, got %q", narrative.Body)
	}
	if !strings.Contains(narrative.Body, "triple para quebrar el cierre") {
		t.Fatalf("expected own key moment in body, got %q", narrative.Body)
	}
}

func TestBuildPostMatchNarrativeIncludesStreakContext(t *testing.T) {
	event := sampleMatchFinished(true, true, 118, 104)

	narrative := BuildPostMatchNarrativeWithContext(event, NarrativeContext{WinStreak: 4})

	if narrative.Metadata["win_streak"] != "4" {
		t.Fatalf("expected win streak metadata, got %q", narrative.Metadata["win_streak"])
	}
	if !strings.Contains(narrative.Body, "racha positiva ya llega a 4 victorias") {
		t.Fatalf("expected streak in body, got %q", narrative.Body)
	}
}

func sampleMatchFinished(homeGame, won bool, ownScore, opponentScore uint16) MatchFinishedEvent {
	own := MatchTeam{TeamID: OwnTeamID, Name: "PulseCity Astrals", Abbreviation: "PCA"}
	opponent := MatchTeam{TeamID: "rival-1", Name: "Seattle Rainmakers", Abbreviation: "SEA"}
	homeTeam := opponent
	awayTeam := own
	homeScore := opponentScore
	awayScore := ownScore
	if homeGame {
		homeTeam = own
		awayTeam = opponent
		homeScore = ownScore
		awayScore = opponentScore
	}
	winnerTeamID := "rival-1"
	if won {
		winnerTeamID = OwnTeamID
	}

	return MatchFinishedEvent{
		EventMeta: EventMeta{
			EventID:       "match-finished-match-1",
			GameID:        "game-1",
			OccurredAt:    "2026-05-24T00:00:00Z",
			SchemaVersion: 1,
		},
		MatchID:       "match-1",
		SimulatedDate: "2026-10-22",
		HomeTeam:      homeTeam,
		AwayTeam:      awayTeam,
		HomeScore:     homeScore,
		AwayScore:     awayScore,
		WinnerTeamID:  winnerTeamID,
		Seed:          123,
		BoxScore: []PlayerBoxScore{
			{PlayerID: "player-1", TeamID: OwnTeamID, Points: 18},
			{PlayerID: "player-2", TeamID: OwnTeamID, Points: 28},
			{PlayerID: "rival-player", TeamID: "rival-1", Points: 42},
		},
		KeyMoments: []KeyMoment{
			{
				Quarter:     4,
				Clock:       "01:12",
				Kind:        "clutch_shot",
				Description: "player-2 mete un triple para quebrar el cierre.",
				TeamID:      OwnTeamID,
				PlayerID:    "player-2",
			},
		},
	}
}
