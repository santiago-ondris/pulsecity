package domain

const (
	SubjectMatchFinished     = "partido.terminado"
	SubjectCityEconomyChange = "ciudad.economia_cambio"
	SubjectCityLandUpdated   = "ciudad.suelo_actualizado"
	SubjectAgentStateChanged = "agente.estado_cambio"
)

type EventMeta struct {
	EventID       string `json:"event_id"`
	GameID        string `json:"game_id"`
	OccurredAt    string `json:"occurred_at"`
	SchemaVersion uint16 `json:"schema_version"`
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
}

type MatchTeam struct {
	TeamID       string `json:"team_id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
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

type CityEconomyChangeEvent struct {
	EventMeta
	SimulatedDate     string  `json:"simulated_date"`
	SourceEventID     string  `json:"source_event_id"`
	SourceSubject     string  `json:"source_subject"`
	FanSentimentDelta float64 `json:"fan_sentiment_delta"`
	TicketSalesDelta  float64 `json:"ticket_sales_delta"`
	LocalEconomyDelta float64 `json:"local_economy_delta"`
	FanSentiment      float64 `json:"fan_sentiment"`
	TicketSalesIndex  float64 `json:"ticket_sales_index"`
	LocalEconomyIndex float64 `json:"local_economy_index"`
	WinStreak         uint16  `json:"win_streak"`
	LossStreak        uint16  `json:"loss_streak"`
	Reason            string  `json:"reason"`
}

type CityLandUpdatedEvent struct {
	EventMeta
	SimulatedDate  string  `json:"simulated_date"`
	ZoneID         string  `json:"zone_id"`
	LandValueDelta float64 `json:"land_value_delta"`
	NewLandValue   float64 `json:"new_land_value"`
	SourceEventID  string  `json:"source_event_id"`
	Reason         string  `json:"reason"`
}

type AgentStateChangedEvent struct {
	EventMeta
	SimulatedDate string             `json:"simulated_date"`
	AgentID       string             `json:"agent_id"`
	SourceEventID string             `json:"source_event_id"`
	SourceSubject string             `json:"source_subject"`
	Mood          string             `json:"mood"`
	State         map[string]float64 `json:"state"`
	Summary       string             `json:"summary"`
}
