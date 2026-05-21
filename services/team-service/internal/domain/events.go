package domain

const (
	SubjectTimeDayAdvanced  = "tiempo.dia_avanzado"
	SubjectMatchScheduled   = "partido.programado"
	SubjectMatchStarting    = "partido.iniciando"
	SubjectMatchFinished    = "partido.terminado"
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
	MatchID       string    `json:"match_id"`
	SimulatedDate string    `json:"simulated_date"`
	HomeTeam      MatchTeam `json:"home_team"`
	AwayTeam      MatchTeam `json:"away_team"`
	Seed          uint64    `json:"seed"`
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
