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

func TestCoveredDayAdvancedEventsExpandsSkippedDates(t *testing.T) {
	days, err := CoveredDayAdvancedEvents(DayAdvancedEvent{
		EventMeta:     EventMeta{GameID: "game-1"},
		SimulatedDate: "2026-10-24",
		Speed:         20,
		DaysProcessed: 3,
	})
	if err != nil {
		t.Fatalf("CoveredDayAdvancedEvents() error = %v, want nil", err)
	}

	if len(days) != 3 {
		t.Fatalf("CoveredDayAdvancedEvents() len = %d, want 3", len(days))
	}
	if days[0].SimulatedDate != "2026-10-22" {
		t.Fatalf("first covered date = %q, want 2026-10-22", days[0].SimulatedDate)
	}
	if days[1].SimulatedDate != "2026-10-23" {
		t.Fatalf("second covered date = %q, want 2026-10-23", days[1].SimulatedDate)
	}
	if days[2].SimulatedDate != "2026-10-24" {
		t.Fatalf("third covered date = %q, want 2026-10-24", days[2].SimulatedDate)
	}
	for _, day := range days {
		if day.DaysProcessed != 1 {
			t.Fatalf("covered day DaysProcessed = %d, want 1", day.DaysProcessed)
		}
	}
}

func TestCoveredDayAdvancedEventsRejectsInvalidDate(t *testing.T) {
	_, err := CoveredDayAdvancedEvents(DayAdvancedEvent{
		SimulatedDate: "invalid",
		DaysProcessed: 1,
	})
	if err == nil {
		t.Fatal("CoveredDayAdvancedEvents() error = nil, want error")
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
	if event.HomeTactics.System == "" || event.AwayTactics.System == "" {
		t.Fatal("BuildMatchScheduledEvent() tactics missing system")
	}
	if event.Players[0].ExpectedMinutes != 34 {
		t.Fatalf("BuildMatchScheduledEvent() first player minutes = %d, want 34", event.Players[0].ExpectedMinutes)
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

func TestBuildPreparedMatchScheduledEventUsesPlayerMatchState(t *testing.T) {
	season, err := GenerateInitialSeason(GameStartedEvent{
		GameID: "game-1",
	})
	if err != nil {
		t.Fatalf("GenerateInitialSeason() error = %v, want nil", err)
	}
	player := season.Roster[0]

	event := BuildPreparedMatchScheduledEvent(DayAdvancedEvent{
		EventMeta: EventMeta{
			GameID:     "game-1",
			OccurredAt: "2026-10-22T00:00:00Z",
		},
		SimulatedDate: "2026-10-22",
	}, season.Schedule[0], season.Roster, MatchPreparation{
		PlayerStates: map[string]PlayerMatchState{
			player.PlayerID: {
				RecentMinutes:  96,
				EmotionalState: 4,
			},
		},
	})

	line := matchPlayerByID(t, event.Players, player.PlayerID)
	if line.Fatigue <= uint8(player.SortOrder%6) {
		t.Fatalf("prepared player fatigue = %d, want above baseline", line.Fatigue)
	}
	if line.EmotionalState != 4 {
		t.Fatalf("prepared player emotional state = %d, want 4", line.EmotionalState)
	}
}

func TestBuildPreparedMatchScheduledEventDerivesCoachTactics(t *testing.T) {
	season, err := GenerateInitialSeason(GameStartedEvent{
		GameID: "game-1",
	})
	if err != nil {
		t.Fatalf("GenerateInitialSeason() error = %v, want nil", err)
	}
	match := season.Schedule[0]
	match.OpponentTeam.OffenseRating = 85
	match.OpponentTeam.Pace = 103
	match.AwayTeam = match.OpponentTeam

	event := BuildPreparedMatchScheduledEvent(DayAdvancedEvent{
		EventMeta: EventMeta{
			GameID:     "game-1",
			OccurredAt: "2026-10-22T00:00:00Z",
		},
		SimulatedDate: "2026-10-22",
	}, match, season.Roster, MatchPreparation{
		Record: SeasonRecord{
			GameID: "game-1",
			Wins:   2,
			Losses: 6,
		},
	})

	if event.HomeTactics.System != "defensive_grind" {
		t.Fatalf("home tactics system = %q, want defensive_grind", event.HomeTactics.System)
	}
	if event.HomeTactics.RotationPreference != "top_heavy" {
		t.Fatalf("home rotation = %q, want top_heavy", event.HomeTactics.RotationPreference)
	}
	if event.HomeTactics.Flexibility <= 58 {
		t.Fatalf("home flexibility = %d, want above 58", event.HomeTactics.Flexibility)
	}
	if event.AwayTactics.System != "pace_and_space" {
		t.Fatalf("away tactics system = %q, want pace_and_space", event.AwayTactics.System)
	}
}

func matchPlayerByID(t *testing.T, players []MatchPlayer, playerID string) MatchPlayer {
	t.Helper()
	for _, player := range players {
		if player.PlayerID == playerID {
			return player
		}
	}

	t.Fatalf("player %q not found", playerID)
	return MatchPlayer{}
}
