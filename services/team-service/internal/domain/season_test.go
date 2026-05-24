package domain

import "testing"

func TestGenerateInitialSeason(t *testing.T) {
	season, err := GenerateInitialSeason(GameStartedEvent{
		GameID:        "game-1",
		FranchiseName: "Lighthouses",
		Abbreviation:  "LHT",
	})
	if err != nil {
		t.Fatalf("GenerateInitialSeason() error = %v, want nil", err)
	}

	if len(season.Roster) != 15 {
		t.Fatalf("GenerateInitialSeason() roster length = %d, want 15", len(season.Roster))
	}
	if len(season.Opponents) != 30 {
		t.Fatalf("GenerateInitialSeason() opponents length = %d, want 30", len(season.Opponents))
	}
	if len(season.Schedule) != RegularSeasonGames {
		t.Fatalf("GenerateInitialSeason() schedule length = %d, want %d", len(season.Schedule), RegularSeasonGames)
	}
	if season.Schedule[0].SimulatedDate != DefaultSeasonStartDate {
		t.Fatalf("GenerateInitialSeason() first date = %q, want %q", season.Schedule[0].SimulatedDate, DefaultSeasonStartDate)
	}
	if !season.Schedule[0].HomeGame {
		t.Fatal("GenerateInitialSeason() first game should be home game")
	}
	if season.Schedule[1].HomeGame {
		t.Fatal("GenerateInitialSeason() second game should be away game")
	}
}

func TestGenerateInitialSeasonIsDeterministic(t *testing.T) {
	event := GameStartedEvent{
		GameID:        "game-1",
		FranchiseName: "Lighthouses",
		Abbreviation:  "LHT",
	}

	first, err := GenerateInitialSeason(event)
	if err != nil {
		t.Fatalf("GenerateInitialSeason(first) error = %v, want nil", err)
	}
	second, err := GenerateInitialSeason(event)
	if err != nil {
		t.Fatalf("GenerateInitialSeason(second) error = %v, want nil", err)
	}

	if first.Roster[0].PlayerID != second.Roster[0].PlayerID {
		t.Fatalf("GenerateInitialSeason() first player id = %q, want %q", first.Roster[0].PlayerID, second.Roster[0].PlayerID)
	}
	if first.Schedule[10].MatchID != second.Schedule[10].MatchID {
		t.Fatalf("GenerateInitialSeason() match id = %q, want %q", first.Schedule[10].MatchID, second.Schedule[10].MatchID)
	}
	if first.Schedule[10].Seed != second.Schedule[10].Seed {
		t.Fatalf("GenerateInitialSeason() seed = %d, want %d", first.Schedule[10].Seed, second.Schedule[10].Seed)
	}
}

func TestGenerateInitialSeasonRequiresGameID(t *testing.T) {
	if _, err := GenerateInitialSeason(GameStartedEvent{}); err == nil {
		t.Fatal("GenerateInitialSeason(empty) error = nil, want error")
	}
}

func TestNextScheduledMatch(t *testing.T) {
	season, err := GenerateInitialSeason(GameStartedEvent{GameID: "game-1"})
	if err != nil {
		t.Fatalf("GenerateInitialSeason() error = %v, want nil", err)
	}

	match, ok := NextScheduledMatch(season.Schedule, "2026-10-25")
	if !ok {
		t.Fatal("NextScheduledMatch() found = false, want true")
	}
	if match.SimulatedDate != "2026-10-26" {
		t.Fatalf("NextScheduledMatch() date = %q, want 2026-10-26", match.SimulatedDate)
	}
}

func TestBuildMatchScheduledEventIncludesFullInput(t *testing.T) {
	season, err := GenerateInitialSeason(GameStartedEvent{
		GameID:        "game-1",
		FranchiseName: "Lighthouses",
		Abbreviation:  "LHT",
	})
	if err != nil {
		t.Fatalf("GenerateInitialSeason() error = %v, want nil", err)
	}

	event := BuildMatchScheduledEvent(DayAdvancedEvent{
		EventMeta: EventMeta{
			GameID:     "game-1",
			OccurredAt: "2026-10-22T00:00:00Z",
		},
		SimulatedDate: "2026-10-22",
	}, season.Schedule[0], season.Roster)

	if event.EventID != "match-scheduled-game-1-match-001" {
		t.Fatalf("BuildMatchScheduledEvent() event_id = %q, want match-scheduled-game-1-match-001", event.EventID)
	}
	if len(event.Players) != 25 {
		t.Fatalf("BuildMatchScheduledEvent() players = %d, want 25", len(event.Players))
	}

	ownPlayers := 0
	opponentPlayers := 0
	for _, player := range event.Players {
		switch player.TeamID {
		case OwnTeamID:
			ownPlayers++
		case season.Schedule[0].OpponentTeam.TeamID:
			opponentPlayers++
		}
	}
	if ownPlayers != 15 {
		t.Fatalf("BuildMatchScheduledEvent() own players = %d, want 15", ownPlayers)
	}
	if opponentPlayers != 10 {
		t.Fatalf("BuildMatchScheduledEvent() opponent players = %d, want 10", opponentPlayers)
	}
}
