package domain

import (
	"fmt"
	"hash/fnv"
	"strings"
	"time"
)

const (
	SubjectMapGenerationStarted = "mapa.generacion_iniciada"
	DefaultSeasonStartDate      = "2026-10-22"
	RegularSeasonGames          = 82
	OwnTeamID                   = "pulsecity"
)

type GameStartedEvent struct {
	GameID        string `json:"game_id"`
	CityName      string `json:"city_name,omitempty"`
	FranchiseName string `json:"franchise_name,omitempty"`
	Abbreviation  string `json:"abbreviation,omitempty"`
}

type InitialSeason struct {
	GameID    string
	Team      MatchTeam
	Roster    []RosterPlayer
	Opponents []MatchTeam
	Schedule  []ScheduledMatch
}

type RosterPlayer struct {
	PlayerID      string
	GameID        string
	FullName      string
	Position      string
	OverallRating uint8
	RosterStatus  string
	ContractYears uint8
	Salary        int
	SortOrder     int
}

type ScheduledMatch struct {
	MatchID       string
	GameID        string
	SimulatedDate string
	HomeTeam      MatchTeam
	AwayTeam      MatchTeam
	OpponentTeam  MatchTeam
	HomeGame      bool
	Seed          uint64
	Status        string
}

func BuildMatchScheduledEvent(day DayAdvancedEvent, match ScheduledMatch, roster []RosterPlayer) MatchScheduledEvent {
	players := make([]MatchPlayer, 0, len(roster)+10)
	for _, player := range roster {
		if player.RosterStatus != "active" {
			continue
		}
		players = append(players, PlayerForRoster(player))
	}
	players = append(players, AbstractOpponentPlayers(match.GameID, match.MatchID, match.OpponentTeam)...)

	return MatchScheduledEvent{
		EventMeta: EventMeta{
			EventID:       fmt.Sprintf("match-scheduled-%s", match.MatchID),
			GameID:        match.GameID,
			OccurredAt:    day.OccurredAt,
			SchemaVersion: 1,
		},
		MatchID:       match.MatchID,
		SimulatedDate: match.SimulatedDate,
		HomeTeam:      match.HomeTeam,
		AwayTeam:      match.AwayTeam,
		Players:       players,
		Seed:          match.Seed,
	}
}

func PlayerForRoster(player RosterPlayer) MatchPlayer {
	base := player.OverallRating
	return MatchPlayer{
		PlayerID:       player.PlayerID,
		TeamID:         OwnTeamID,
		Rating:         base,
		Scoring:        clampRating(int(base) + scoringPositionBias(player.Position)),
		Rebounding:     clampRating(int(base) + reboundingPositionBias(player.Position)),
		Playmaking:     clampRating(int(base) + playmakingPositionBias(player.Position)),
		Defense:        clampRating(int(base) - 1),
		Stamina:        82,
		Fatigue:        uint8(player.SortOrder % 9),
		EmotionalState: 0,
	}
}

func AbstractOpponentPlayers(gameID, matchID string, opponent MatchTeam) []MatchPlayer {
	positions := []string{"PG", "SG", "SF", "PF", "C", "G", "F", "F", "C", "G"}
	players := make([]MatchPlayer, 0, len(positions))
	for index, position := range positions {
		variation := int(stableHash(gameID, matchID, opponent.TeamID, position, fmt.Sprint(index))%7) - 3
		base := clampRating(int(opponent.Rating) + variation)
		players = append(players, MatchPlayer{
			PlayerID:       fmt.Sprintf("%s-%s-player-%02d", matchID, opponent.TeamID, index+1),
			TeamID:         opponent.TeamID,
			Rating:         base,
			Scoring:        clampRating(int(base) + scoringPositionBias(position)),
			Rebounding:     clampRating(int(base) + reboundingPositionBias(position)),
			Playmaking:     clampRating(int(base) + playmakingPositionBias(position)),
			Defense:        clampRating(int(base) - 1),
			Stamina:        80,
			Fatigue:        uint8(index % 7),
			EmotionalState: 0,
		})
	}

	return players
}

func GenerateInitialSeason(event GameStartedEvent) (InitialSeason, error) {
	gameID := strings.TrimSpace(event.GameID)
	if gameID == "" {
		return InitialSeason{}, fmt.Errorf("generate initial season: missing game id")
	}

	team := ownTeam(event)
	roster := generateRoster(gameID)
	opponents := generateOpponents(gameID)
	schedule, err := generateSchedule(gameID, team, opponents)
	if err != nil {
		return InitialSeason{}, err
	}

	return InitialSeason{
		GameID:    gameID,
		Team:      team,
		Roster:    roster,
		Opponents: opponents,
		Schedule:  schedule,
	}, nil
}

func NextScheduledMatch(schedule []ScheduledMatch, simulatedDate string) (ScheduledMatch, bool) {
	for _, match := range schedule {
		if match.Status == "scheduled" && match.SimulatedDate >= simulatedDate {
			return match, true
		}
	}

	return ScheduledMatch{}, false
}

func ownTeam(event GameStartedEvent) MatchTeam {
	name := strings.TrimSpace(event.FranchiseName)
	if name == "" {
		name = "PulseCity Basketball Club"
	}
	abbreviation := strings.ToUpper(strings.TrimSpace(event.Abbreviation))
	if abbreviation == "" {
		abbreviation = "PUL"
	}

	return MatchTeam{
		TeamID:             OwnTeamID,
		Name:               name,
		Abbreviation:       abbreviation,
		Rating:             76,
		OffenseRating:      75,
		DefenseRating:      74,
		Pace:               98,
		HomeCourtAdvantage: 3,
	}
}

func generateRoster(gameID string) []RosterPlayer {
	templates := []struct {
		name     string
		position string
		base     uint8
	}{
		{"Mateo Cross", "PG", 80},
		{"Elias Monroe", "SG", 78},
		{"Noah Sterling", "SF", 79},
		{"Julian Mercer", "PF", 77},
		{"Theo Banks", "C", 81},
		{"Adrian Vale", "PG", 74},
		{"Camden Price", "SG", 73},
		{"Lucian Reed", "SF", 75},
		{"Nico Hayes", "PF", 72},
		{"Dante Brooks", "C", 74},
		{"Silas Hart", "G", 71},
		{"Roman Ellis", "F", 70},
		{"Kieran Fox", "F", 69},
		{"Milo Grant", "C", 68},
		{"Jonas Pike", "G", 67},
	}

	players := make([]RosterPlayer, 0, len(templates))
	for index, template := range templates {
		variation := uint8(stableHash(gameID, template.name) % 3)
		players = append(players, RosterPlayer{
			PlayerID:      fmt.Sprintf("%s-player-%02d", gameID, index+1),
			GameID:        gameID,
			FullName:      template.name,
			Position:      template.position,
			OverallRating: template.base + variation,
			RosterStatus:  "active",
			ContractYears: uint8(1 + stableHash(template.name, gameID)%4),
			Salary:        1_800_000 + int(stableHash(gameID, template.position, template.name)%18)*250_000,
			SortOrder:     index + 1,
		})
	}

	return players
}

func generateOpponents(gameID string) []MatchTeam {
	templates := []struct {
		name         string
		abbreviation string
		rating       uint8
		offense      uint8
		defense      uint8
		pace         uint8
	}{
		{"Baltimore Foundry", "BAL", 77, 78, 75, 97},
		{"Cincinnati Monarchs", "CIN", 74, 73, 75, 95},
		{"Las Vegas Neon", "LVN", 82, 84, 78, 101},
		{"Seattle Rainmakers", "SEA", 80, 79, 81, 96},
		{"Kansas City Rail", "KCR", 73, 72, 74, 94},
		{"Vancouver Tides", "VAN", 78, 77, 79, 98},
		{"Austin Comets", "AUS", 81, 83, 77, 102},
		{"Pittsburgh Iron", "PIT", 76, 74, 78, 93},
		{"Nashville Sound", "NSH", 75, 76, 73, 99},
		{"San Diego Breakers", "SDB", 79, 80, 77, 100},
		{"Louisville Colonels", "LOU", 72, 71, 73, 92},
		{"St. Louis Arches", "STL", 78, 76, 80, 96},
		{"Tampa Bay Suns", "TBS", 77, 79, 74, 101},
		{"Montreal Saints", "MTL", 80, 78, 82, 95},
		{"Mexico City Volcanes", "MCV", 83, 85, 79, 103},
		{"Columbus Forge", "CLB", 74, 75, 73, 97},
		{"Raleigh Flight", "RAL", 76, 77, 75, 99},
		{"San Jose Circuit", "SJC", 79, 81, 76, 100},
		{"Omaha Plains", "OMA", 71, 70, 72, 94},
		{"Virginia Beach Admirals", "VBA", 75, 74, 76, 96},
		{"Providence Lanterns", "PRO", 73, 72, 74, 95},
		{"Albuquerque Mesa", "ABQ", 76, 78, 73, 101},
		{"Birmingham Steel", "BHM", 72, 71, 74, 93},
		{"Portland Pines", "POR", 81, 80, 82, 97},
		{"Buffalo Blizzard", "BUF", 74, 73, 76, 94},
		{"Honolulu Waves", "HNL", 78, 80, 75, 102},
		{"Richmond Union", "RIC", 73, 72, 75, 95},
		{"Boise Summit", "BOI", 75, 74, 76, 97},
		{"Jacksonville Current", "JAX", 77, 78, 76, 99},
		{"Milwaukee Northstars", "MNS", 82, 81, 83, 96},
	}

	opponents := make([]MatchTeam, 0, len(templates))
	for index, template := range templates {
		delta := int8(stableHash(gameID, template.abbreviation)%3) - 1
		opponents = append(opponents, MatchTeam{
			TeamID:             fmt.Sprintf("rival-%02d", index+1),
			Name:               template.name,
			Abbreviation:       template.abbreviation,
			Rating:             applyRatingDelta(template.rating, delta),
			OffenseRating:      applyRatingDelta(template.offense, delta),
			DefenseRating:      applyRatingDelta(template.defense, delta),
			Pace:               applyRatingDelta(template.pace, delta),
			HomeCourtAdvantage: 2,
		})
	}

	return opponents
}

func generateSchedule(gameID string, team MatchTeam, opponents []MatchTeam) ([]ScheduledMatch, error) {
	startDate, err := time.Parse(time.DateOnly, DefaultSeasonStartDate)
	if err != nil {
		return nil, fmt.Errorf("parse season start date: %w", err)
	}

	schedule := make([]ScheduledMatch, 0, RegularSeasonGames)
	for index := 0; index < RegularSeasonGames; index++ {
		opponent := opponents[index%len(opponents)]
		matchDate := startDate.AddDate(0, 0, index*2)
		homeGame := index%2 == 0
		homeTeam := team
		awayTeam := opponent
		if !homeGame {
			homeTeam = opponent
			awayTeam = team
		}

		schedule = append(schedule, ScheduledMatch{
			MatchID:       fmt.Sprintf("%s-match-%03d", gameID, index+1),
			GameID:        gameID,
			SimulatedDate: matchDate.Format(time.DateOnly),
			HomeTeam:      homeTeam,
			AwayTeam:      awayTeam,
			OpponentTeam:  opponent,
			HomeGame:      homeGame,
			Seed:          stablePositiveSeed(gameID, opponent.TeamID, matchDate.Format(time.DateOnly)),
			Status:        "scheduled",
		})
	}

	return schedule, nil
}

func applyRatingDelta(base uint8, delta int8) uint8 {
	if delta < 0 {
		return base - uint8(-delta)
	}

	return base + uint8(delta)
}

func clampRating(value int) uint8 {
	if value < 40 {
		return 40
	}
	if value > 99 {
		return 99
	}

	return uint8(value)
}

func scoringPositionBias(position string) int {
	switch position {
	case "PG", "SG", "G":
		return 3
	case "SF":
		return 1
	default:
		return -1
	}
}

func reboundingPositionBias(position string) int {
	switch position {
	case "C":
		return 5
	case "PF", "F":
		return 3
	case "SF":
		return 1
	default:
		return -3
	}
}

func playmakingPositionBias(position string) int {
	switch position {
	case "PG":
		return 6
	case "SG", "G":
		return 3
	case "SF":
		return 1
	default:
		return -3
	}
}

func stableHash(parts ...string) uint64 {
	hasher := fnv.New64a()
	for _, part := range parts {
		_, _ = hasher.Write([]byte(part))
		_, _ = hasher.Write([]byte{0})
	}

	return hasher.Sum64()
}

func stablePositiveSeed(parts ...string) uint64 {
	const maxInt64 = uint64(1<<63 - 1)

	return stableHash(parts...) % maxInt64
}
