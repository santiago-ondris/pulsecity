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
	DefaultCapBase              = 141_000_000
	DefaultLuxuryTaxLine        = 171_000_000
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

type SalaryCapSnapshot struct {
	GameID              string
	SimulatedDate       string
	CapBase             int
	LuxuryTaxLine       int
	CommittedSalary     int
	CapSpace            int
	LuxuryTaxSpace      int
	RosterCount         uint8
	Status              string
	NearLuxuryTax       bool
	ProjectedTaxPayment int
	SourceEventID       string
	SourceSubject       string
	OccurredAt          string
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

func MaterializeIncomingTradePlayer(
	gameID string,
	proposalID string,
	requestedPosition string,
	incomingSalary int,
	outgoingRating uint8,
	sortOrder int,
) RosterPlayer {
	firstNames := []string{"Jalen", "Marcus", "Tobias", "Dorian", "Caleb", "Isaiah", "Miles", "Andre"}
	lastNames := []string{"Warren", "Porter", "Holland", "Bennett", "Sullivan", "Maddox", "Foster", "Graves"}
	first := firstNames[stableHash(gameID, proposalID, "first")%uint64(len(firstNames))]
	last := lastNames[stableHash(gameID, proposalID, "last")%uint64(len(lastNames))]
	ratingDelta := int(stableHash(gameID, proposalID, "rating")%5) - 2

	return RosterPlayer{
		PlayerID:      proposalID + "-incoming",
		GameID:        gameID,
		FullName:      first + " " + last,
		Position:      requestedPosition,
		OverallRating: clampRating(int(outgoingRating) + ratingDelta),
		RosterStatus:  "active",
		ContractYears: uint8(1 + stableHash(gameID, proposalID, "years")%3),
		Salary:        incomingSalary,
		SortOrder:     sortOrder,
	}
}

func CalculateSalaryCap(
	gameID string,
	roster []RosterPlayer,
	simulatedDate string,
	occurredAt string,
	sourceEventID string,
	sourceSubject string,
) SalaryCapSnapshot {
	committedSalary := 0
	rosterCount := 0
	for _, player := range roster {
		if player.RosterStatus == "waived" || player.RosterStatus == "traded" {
			continue
		}
		committedSalary += player.Salary
		rosterCount++
	}

	capSpace := DefaultCapBase - committedSalary
	luxuryTaxSpace := DefaultLuxuryTaxLine - committedSalary
	status := "under_cap"
	projectedTaxPayment := 0
	if committedSalary > DefaultLuxuryTaxLine {
		status = "luxury_tax"
		projectedTaxPayment = (committedSalary - DefaultLuxuryTaxLine) * 2
	} else if committedSalary > DefaultCapBase {
		status = "over_cap"
	}

	return SalaryCapSnapshot{
		GameID:              gameID,
		SimulatedDate:       simulatedDate,
		CapBase:             DefaultCapBase,
		LuxuryTaxLine:       DefaultLuxuryTaxLine,
		CommittedSalary:     committedSalary,
		CapSpace:            capSpace,
		LuxuryTaxSpace:      luxuryTaxSpace,
		RosterCount:         uint8(rosterCount),
		Status:              status,
		NearLuxuryTax:       luxuryTaxSpace <= 10_000_000,
		ProjectedTaxPayment: projectedTaxPayment,
		SourceEventID:       sourceEventID,
		SourceSubject:       sourceSubject,
		OccurredAt:          occurredAt,
	}
}

func (snapshot SalaryCapSnapshot) SalaryCapCalculatedEvent() SalaryCapCalculatedEvent {
	return SalaryCapCalculatedEvent{
		EventMeta: EventMeta{
			EventID:       "salary-cap-" + snapshot.GameID + "-" + snapshot.SourceEventID,
			GameID:        snapshot.GameID,
			OccurredAt:    snapshot.OccurredAt,
			SchemaVersion: 1,
		},
		SimulatedDate:       snapshot.SimulatedDate,
		CapBase:             snapshot.CapBase,
		LuxuryTaxLine:       snapshot.LuxuryTaxLine,
		CommittedSalary:     snapshot.CommittedSalary,
		CapSpace:            snapshot.CapSpace,
		LuxuryTaxSpace:      snapshot.LuxuryTaxSpace,
		RosterCount:         snapshot.RosterCount,
		Status:              snapshot.Status,
		NearLuxuryTax:       snapshot.NearLuxuryTax,
		ProjectedTaxPayment: snapshot.ProjectedTaxPayment,
	}
}

func FinancePatchFromSalaryCap(event SalaryCapCalculatedEvent) FinancePatchEvent {
	return FinancePatchEvent{
		Type:    SubjectFinancePatch,
		Subject: SubjectFinancePatch,
		GameID:  event.GameID,
		Patch: FinanceStatePatch{
			SimulatedDate:       event.SimulatedDate,
			SourceEventID:       event.EventID,
			SourceSubject:       SubjectSalaryCap,
			CapBase:             event.CapBase,
			LuxuryTaxLine:       event.LuxuryTaxLine,
			CommittedSalary:     event.CommittedSalary,
			CapSpace:            event.CapSpace,
			LuxuryTaxSpace:      event.LuxuryTaxSpace,
			RosterCount:         event.RosterCount,
			Status:              event.Status,
			NearLuxuryTax:       event.NearLuxuryTax,
			ProjectedTaxPayment: event.ProjectedTaxPayment,
		},
	}
}

type PlayerMatchState struct {
	RecentMinutes  uint16
	EmotionalState int8
}

type InjuryAssessmentInput struct {
	GameID         string
	MatchID        string
	OccurredAt     string
	SimulatedDate  string
	PlayerID       string
	TeamID         string
	Minutes        uint8
	RecentMinutes  uint16
	EmotionalState int8
}

type MatchPreparation struct {
	PlayerStates map[string]PlayerMatchState
	Record       SeasonRecord
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
	return BuildPreparedMatchScheduledEvent(day, match, roster, MatchPreparation{})
}

func BuildPreparedMatchScheduledEvent(
	day DayAdvancedEvent,
	match ScheduledMatch,
	roster []RosterPlayer,
	preparation MatchPreparation,
) MatchScheduledEvent {
	players := make([]MatchPlayer, 0, len(roster)+10)
	for _, player := range roster {
		if player.RosterStatus != "active" {
			continue
		}
		players = append(players, PlayerForRoster(player, preparation.PlayerStates[player.PlayerID]))
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
		HomeTactics:   TacticalContextForTeam(match.HomeTeam, match.OpponentTeam, preparation.Record),
		AwayTactics:   TacticalContextForTeam(match.AwayTeam, match.OpponentTeam, preparation.Record),
		Players:       players,
		Seed:          match.Seed,
	}
}

func PlayerForRoster(player RosterPlayer, state PlayerMatchState) MatchPlayer {
	base := player.OverallRating
	return MatchPlayer{
		PlayerID:        player.PlayerID,
		TeamID:          OwnTeamID,
		ExpectedMinutes: expectedMinutesForRotationSlot(player.SortOrder - 1),
		Rating:          base,
		Scoring:         clampRating(int(base) + scoringPositionBias(player.Position)),
		Rebounding:      clampRating(int(base) + reboundingPositionBias(player.Position)),
		Playmaking:      clampRating(int(base) + playmakingPositionBias(player.Position)),
		Defense:         clampRating(int(base) - 1),
		Stamina:         82,
		Fatigue:         fatigueFromRecentMinutes(state.RecentMinutes, player.SortOrder),
		EmotionalState:  state.EmotionalState,
	}
}

func AbstractOpponentPlayers(gameID, matchID string, opponent MatchTeam) []MatchPlayer {
	positions := []string{"PG", "SG", "SF", "PF", "C", "G", "F", "F", "C", "G"}
	players := make([]MatchPlayer, 0, len(positions))
	for index, position := range positions {
		variation := int(stableHash(gameID, matchID, opponent.TeamID, position, fmt.Sprint(index))%7) - 3
		base := clampRating(int(opponent.Rating) + variation)
		players = append(players, MatchPlayer{
			PlayerID:        fmt.Sprintf("%s-%s-player-%02d", matchID, opponent.TeamID, index+1),
			TeamID:          opponent.TeamID,
			ExpectedMinutes: expectedMinutesForRotationSlot(index),
			Rating:          base,
			Scoring:         clampRating(int(base) + scoringPositionBias(position)),
			Rebounding:      clampRating(int(base) + reboundingPositionBias(position)),
			Playmaking:      clampRating(int(base) + playmakingPositionBias(position)),
			Defense:         clampRating(int(base) - 1),
			Stamina:         80,
			Fatigue:         uint8(index % 7),
			EmotionalState:  0,
		})
	}

	return players
}

func TacticalContextForTeam(team, opponent MatchTeam, record SeasonRecord) MatchTacticalContext {
	if team.TeamID != OwnTeamID {
		return opponentTacticalContext(team)
	}

	system := "balanced"
	rotationPreference := "standard"
	flexibility := uint8(58)

	if opponent.Pace >= 100 || opponent.OffenseRating >= 82 {
		system = "defensive_grind"
		flexibility = 62
	}
	if opponent.DefenseRating <= 74 {
		system = "pace_and_space"
	}
	if record.Losses > record.Wins+2 {
		rotationPreference = "top_heavy"
		flexibility += 3
	}
	if record.Wins > record.Losses+4 {
		rotationPreference = "deep"
	}

	return MatchTacticalContext{
		System:             system,
		RotationPreference: rotationPreference,
		Flexibility:        flexibility,
	}
}

func opponentTacticalContext(team MatchTeam) MatchTacticalContext {
	system := "balanced"
	if team.Pace >= 100 || team.OffenseRating >= 82 {
		system = "pace_and_space"
	}
	if team.DefenseRating >= team.OffenseRating+3 {
		system = "defensive_grind"
	}

	return MatchTacticalContext{
		System:             system,
		RotationPreference: "standard",
		Flexibility:        52,
	}
}

func expectedMinutesForRotationSlot(slot int) uint8 {
	minutes := []uint8{34, 32, 30, 29, 28, 24, 22, 18, 13, 10, 0, 0, 0, 0, 0}
	if slot < 0 || slot >= len(minutes) {
		return 0
	}

	return minutes[slot]
}

func fatigueFromRecentMinutes(recentMinutes uint16, sortOrder int) uint8 {
	baseline := uint8(sortOrder % 6)
	switch {
	case recentMinutes >= 105:
		return baseline + 18
	case recentMinutes >= 90:
		return baseline + 14
	case recentMinutes >= 72:
		return baseline + 10
	case recentMinutes >= 48:
		return baseline + 6
	default:
		return baseline
	}
}

func EmotionalStateScore(state string) int8 {
	switch strings.TrimSpace(strings.ToLower(state)) {
	case "confident", "energized", "locked_in":
		return 4
	case "satisfied", "steady", "calm":
		return 1
	case "frustrated", "anxious":
		return -3
	case "disconnected", "angry":
		return -5
	default:
		return 0
	}
}

func AssessInjuryRisk(input InjuryAssessmentInput) (PlayerInjuredEvent, bool, error) {
	if input.TeamID != OwnTeamID {
		return PlayerInjuredEvent{}, false, nil
	}

	workload := input.RecentMinutes
	if workload < 96 {
		return PlayerInjuredEvent{}, false, nil
	}

	risk := int(workload-80) / 2
	if input.Minutes >= 34 {
		risk += int(input.Minutes - 30)
	}
	if input.EmotionalState < 0 {
		risk += int(-input.EmotionalState)
	}
	if workload >= 132 {
		risk = 100
	}

	roll := int(stableHash(input.GameID, input.MatchID, input.PlayerID, "injury") % 100)
	if roll >= risk {
		return PlayerInjuredEvent{}, false, nil
	}

	injuredOn, err := time.Parse(time.DateOnly, input.SimulatedDate)
	if err != nil {
		return PlayerInjuredEvent{}, false, fmt.Errorf("assess injury risk: parse simulated date %q: %w", input.SimulatedDate, err)
	}
	severity, estimatedDays := injurySeverity(workload)

	return PlayerInjuredEvent{
		EventMeta: EventMeta{
			EventID:       fmt.Sprintf("player-injured-%s-%s", input.MatchID, input.PlayerID),
			GameID:        input.GameID,
			OccurredAt:    input.OccurredAt,
			SchemaVersion: 1,
		},
		InjuryID:             fmt.Sprintf("injury-%s-%s", input.MatchID, input.PlayerID),
		PlayerID:             input.PlayerID,
		Severity:             severity,
		EstimatedDaysOut:     estimatedDays,
		InjuredOn:            input.SimulatedDate,
		ExpectedRecoveryDate: injuredOn.AddDate(0, 0, int(estimatedDays)).Format(time.DateOnly),
		Reason:               "workload_accumulation",
		SourceMatchID:        input.MatchID,
		WorkloadScore:        workload,
	}, true, nil
}

func injurySeverity(workload uint16) (string, uint16) {
	switch {
	case workload >= 150:
		return "major", 21
	case workload >= 132:
		return "moderate", 10
	default:
		return "minor", 5
	}
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

func CoveredDayAdvancedEvents(event DayAdvancedEvent) ([]DayAdvancedEvent, error) {
	if event.DaysProcessed == 0 {
		return nil, nil
	}

	finalDate, err := time.Parse(time.DateOnly, event.SimulatedDate)
	if err != nil {
		return nil, fmt.Errorf("parse simulated date %q: %w", event.SimulatedDate, err)
	}

	covered := make([]DayAdvancedEvent, 0, event.DaysProcessed)
	startDate := finalDate.AddDate(0, 0, -int(event.DaysProcessed-1))
	for index := 0; index < int(event.DaysProcessed); index++ {
		day := event
		day.SimulatedDate = startDate.AddDate(0, 0, index).Format(time.DateOnly)
		day.DaysProcessed = 1
		covered = append(covered, day)
	}

	return covered, nil
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
		overallRating := template.base + variation
		contractYears := uint8(1 + stableHash(template.name, gameID)%4)
		players = append(players, RosterPlayer{
			PlayerID:      fmt.Sprintf("%s-player-%02d", gameID, index+1),
			GameID:        gameID,
			FullName:      template.name,
			Position:      template.position,
			OverallRating: overallRating,
			RosterStatus:  "active",
			ContractYears: contractYears,
			Salary:        calculatePlayerSalary(overallRating, template.position, contractYears),
			SortOrder:     index + 1,
		})
	}

	return players
}

func calculatePlayerSalary(overallRating uint8, position string, contractYears uint8) int {
	baseSalary := 2_000_000
	switch {
	case overallRating >= 84:
		baseSalary = 30_000_000
	case overallRating >= 81:
		baseSalary = 23_000_000
	case overallRating >= 78:
		baseSalary = 16_000_000
	case overallRating >= 75:
		baseSalary = 10_000_000
	case overallRating >= 72:
		baseSalary = 5_500_000
	}

	if position == "PG" || position == "C" {
		baseSalary += 1_500_000
	}
	if contractYears >= 3 && overallRating >= 78 {
		baseSalary += 2_000_000
	}

	return baseSalary
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
