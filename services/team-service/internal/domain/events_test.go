package domain

import "testing"

func TestEventSubjects(t *testing.T) {
	tests := map[string]string{
		"day advanced":    SubjectTimeDayAdvanced,
		"match scheduled": SubjectMatchScheduled,
		"match finished":  SubjectMatchFinished,
		"trade proposed":  SubjectTradeProposed,
		"trade rejected":  SubjectTradeRejected,
		"trade countered": SubjectTradeCountered,
		"trade accepted":  SubjectTradeAccepted,
	}

	if tests["day advanced"] != "tiempo.dia_avanzado" {
		t.Fatalf("SubjectTimeDayAdvanced = %q", tests["day advanced"])
	}
	if tests["match scheduled"] != "partido.programado" {
		t.Fatalf("SubjectMatchScheduled = %q", tests["match scheduled"])
	}
	if tests["match finished"] != "partido.terminado" {
		t.Fatalf("SubjectMatchFinished = %q", tests["match finished"])
	}
	if tests["trade proposed"] != "trade.propuesta_enviada" {
		t.Fatalf("SubjectTradeProposed = %q", tests["trade proposed"])
	}
	if tests["trade rejected"] != "trade.rechazada" {
		t.Fatalf("SubjectTradeRejected = %q", tests["trade rejected"])
	}
	if tests["trade countered"] != "trade.contraoferta" {
		t.Fatalf("SubjectTradeCountered = %q", tests["trade countered"])
	}
	if tests["trade accepted"] != "trade.aceptada" {
		t.Fatalf("SubjectTradeAccepted = %q", tests["trade accepted"])
	}
}

func TestSeasonPatchFromRecord(t *testing.T) {
	record := SeasonRecord{
		GameID:        "game-1",
		Wins:          3,
		Losses:        2,
		PointsFor:     550,
		PointsAgainst: 540,
		LastResult: &SeasonMatchSummary{
			MatchID:       "match-1",
			SimulatedDate: "2026-10-22",
			HomeScore:     110,
			AwayScore:     101,
			WinnerTeamID:  OwnTeamID,
		},
	}

	event := SeasonPatchFromRecord(record)

	if event.Type != SubjectSeasonPatchDelta {
		t.Fatalf("SeasonPatchFromRecord() Type = %q, want %q", event.Type, SubjectSeasonPatchDelta)
	}
	if event.GameID != "game-1" {
		t.Fatalf("SeasonPatchFromRecord() GameID = %q, want game-1", event.GameID)
	}
	if event.Patch.Wins != 3 {
		t.Fatalf("SeasonPatchFromRecord() wins = %d, want 3", event.Patch.Wins)
	}
	if event.Patch.LastResult == nil {
		t.Fatal("SeasonPatchFromRecord() LastResult = nil, want result")
	}
}
