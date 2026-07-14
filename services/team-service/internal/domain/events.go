package domain

const (
	SubjectTimeDayAdvanced  = "tiempo.dia_avanzado"
	SubjectMatchScheduled   = "partido.programado"
	SubjectMatchStarting    = "partido.iniciando"
	SubjectMatchFinished    = "partido.terminado"
	SubjectRosterPatchDelta = "roster.patch"
	SubjectSeasonPatchDelta = "season.patch"
	SubjectMatchResultDelta = "match.result"
)

type EventMeta struct {
	EventID       string `json:"event_id"`
	GameID        string `json:"game_id"`
	OccurredAt    string `json:"occurred_at"`
	SchemaVersion uint16 `json:"schema_version"`
}

type DayAdvancedEvent struct {
	EventMeta
	SimulatedDate string `json:"simulated_date"`
	Speed         uint8  `json:"speed"`
	DaysProcessed uint16 `json:"days_processed"`
}

type MatchScheduledEvent struct {
	EventMeta
	MatchID       string               `json:"match_id"`
	SimulatedDate string               `json:"simulated_date"`
	HomeTeam      MatchTeam            `json:"home_team"`
	AwayTeam      MatchTeam            `json:"away_team"`
	HomeTactics   MatchTacticalContext `json:"home_tactics"`
	AwayTactics   MatchTacticalContext `json:"away_tactics"`
	Players       []MatchPlayer        `json:"players"`
	Seed          uint64               `json:"seed"`
}

type MatchFinishedEvent struct {
	EventMeta
	MatchID       string           `json:"match_id"`
	SimulatedDate string           `json:"simulated_date"`
	HomeTeam      MatchTeam        `json:"home_team"`
	AwayTeam      MatchTeam        `json:"away_team"`
	HomeScore     uint16           `json:"home_score"`
	AwayScore     uint16           `json:"away_score"`
	WinnerTeamID  string           `json:"winner_team_id"`
	Seed          uint64           `json:"seed"`
	BoxScore      []PlayerBoxScore `json:"box_score"`
	KeyMoments    []KeyMoment      `json:"key_moments"`
}

type RosterPatchEnvelope struct {
	Type    string           `json:"type"`
	Subject string           `json:"subject"`
	GameID  string           `json:"game_id"`
	Patch   RosterStatePatch `json:"patch"`
}

type RosterStatePatch struct {
	SimulatedDate string                 `json:"simulated_date"`
	SourceEventID string                 `json:"source_event_id"`
	SourceSubject string                 `json:"source_subject"`
	Players       []PlayerEmotionalPatch `json:"players"`
}

type PlayerEmotionalPatch struct {
	PlayerID         string  `json:"player_id"`
	EmotionalState   string  `json:"emotional_state"`
	Satisfaction     float64 `json:"satisfaction"`
	Loyalty          float64 `json:"loyalty"`
	Ego              float64 `json:"ego"`
	CompetitiveDrive float64 `json:"competitive_drive"`
	CityConnection   float64 `json:"city_connection"`
	Summary          string  `json:"summary"`
}

type SeasonPatchEvent struct {
	Type    string           `json:"type"`
	Subject string           `json:"subject"`
	GameID  string           `json:"game_id"`
	Patch   SeasonStatePatch `json:"patch"`
}

type SeasonStatePatch struct {
	Wins          uint16              `json:"wins"`
	Losses        uint16              `json:"losses"`
	PointsFor     uint16              `json:"points_for"`
	PointsAgainst uint16              `json:"points_against"`
	LastResult    *SeasonMatchSummary `json:"last_result,omitempty"`
}

type SeasonMatchSummary struct {
	MatchID       string `json:"match_id"`
	SimulatedDate string `json:"simulated_date"`
	HomeTeamID    string `json:"home_team_id"`
	AwayTeamID    string `json:"away_team_id"`
	HomeScore     uint16 `json:"home_score"`
	AwayScore     uint16 `json:"away_score"`
	WinnerTeamID  string `json:"winner_team_id"`
}

type SeasonRecord struct {
	GameID        string
	Wins          uint16
	Losses        uint16
	PointsFor     uint16
	PointsAgainst uint16
	LastResult    *SeasonMatchSummary
}

func SeasonPatchFromRecord(record SeasonRecord) SeasonPatchEvent {
	return SeasonPatchEvent{
		Type:    SubjectSeasonPatchDelta,
		Subject: SubjectSeasonPatchDelta,
		GameID:  record.GameID,
		Patch: SeasonStatePatch{
			Wins:          record.Wins,
			Losses:        record.Losses,
			PointsFor:     record.PointsFor,
			PointsAgainst: record.PointsAgainst,
			LastResult:    record.LastResult,
		},
	}
}

type MatchTeam struct {
	TeamID             string `json:"team_id"`
	Name               string `json:"name"`
	Abbreviation       string `json:"abbreviation"`
	Rating             uint8  `json:"rating"`
	OffenseRating      uint8  `json:"offense_rating"`
	DefenseRating      uint8  `json:"defense_rating"`
	Pace               uint8  `json:"pace"`
	HomeCourtAdvantage int8   `json:"home_court_advantage"`
}

type MatchTacticalContext struct {
	System             string `json:"system"`
	RotationPreference string `json:"rotation_preference"`
	Flexibility        uint8  `json:"flexibility"`
}

type MatchPlayer struct {
	PlayerID        string `json:"player_id"`
	TeamID          string `json:"team_id"`
	ExpectedMinutes uint8  `json:"expected_minutes"`
	Rating          uint8  `json:"rating"`
	Scoring         uint8  `json:"scoring"`
	Rebounding      uint8  `json:"rebounding"`
	Playmaking      uint8  `json:"playmaking"`
	Defense         uint8  `json:"defense"`
	Stamina         uint8  `json:"stamina"`
	Fatigue         uint8  `json:"fatigue"`
	EmotionalState  int8   `json:"emotional_state"`
}

type PlayerBoxScore struct {
	PlayerID  string `json:"player_id"`
	TeamID    string `json:"team_id"`
	Minutes   uint8  `json:"minutes"`
	Points    uint16 `json:"points"`
	Rebounds  uint16 `json:"rebounds"`
	Assists   uint16 `json:"assists"`
	Steals    uint16 `json:"steals"`
	Blocks    uint16 `json:"blocks"`
	Turnovers uint16 `json:"turnovers"`
}

type KeyMoment struct {
	Quarter     uint8  `json:"quarter"`
	Clock       string `json:"clock"`
	Kind        string `json:"kind"`
	Description string `json:"description"`
	TeamID      string `json:"team_id"`
	PlayerID    string `json:"player_id,omitempty"`
}
