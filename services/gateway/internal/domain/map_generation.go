package domain

const (
	SubjectMapGenerationStarted  = "mapa.generacion_iniciada"
	SubjectMapTerrainReady       = "mapa.terreno_listo"
	SubjectMapZonesCalculated    = "mapa.zonas_calculadas"
	SubjectMapStadiumLocated     = "mapa.estadio_ubicado"
	SubjectMapGenerationComplete = "mapa.generacion_completa"
	SubjectMapWildcard           = "mapa.*"
)

type StartGameRequest struct {
	CityName           string `json:"city_name"`
	FranchiseName      string `json:"franchise_name"`
	Abbreviation       string `json:"abbreviation"`
	PrimaryColor       string `json:"primary_color"`
	SecondaryColor     string `json:"secondary_color"`
	AccentColor        string `json:"accent_color"`
	InitialScenario    string `json:"initial_scenario"`
	CityManagementMode string `json:"city_management_mode"`
}

type GuestSession struct {
	GuestToken string `json:"guest_token"`
	CreatedAt  string `json:"created_at,omitempty"`
	LastSeenAt string `json:"last_seen_at,omitempty"`
}

type MapGenerationRequest struct {
	GameID        string `json:"game_id"`
	CityName      string `json:"city_name,omitempty"`
	FranchiseName string `json:"franchise_name,omitempty"`
	Abbreviation  string `json:"abbreviation,omitempty"`
}

type GameSetup struct {
	GameID             string           `json:"game_id"`
	GuestToken         string           `json:"-"`
	UserID             string           `json:"-"`
	OwnerKind          string           `json:"owner_kind,omitempty"`
	CityName           string           `json:"city_name"`
	FranchiseName      string           `json:"franchise_name"`
	Abbreviation       string           `json:"abbreviation"`
	PrimaryColor       string           `json:"primary_color"`
	SecondaryColor     string           `json:"secondary_color"`
	AccentColor        string           `json:"accent_color"`
	InitialScenario    string           `json:"initial_scenario"`
	CityManagementMode string           `json:"city_management_mode"`
	OwnerIntroEvent    *NarrativeEvent  `json:"owner_intro_event,omitempty"`
	OwnerIntroResponse *NarrativeChoice `json:"owner_intro_response,omitempty"`
	Status             string           `json:"status"`
	CreatedAt          string           `json:"created_at,omitempty"`
	UpdatedAt          string           `json:"updated_at,omitempty"`
}

type GameSummary struct {
	GameID             string `json:"game_id"`
	CityName           string `json:"city_name"`
	FranchiseName      string `json:"franchise_name"`
	OwnerKind          string `json:"owner_kind"`
	InitialScenario    string `json:"initial_scenario"`
	CityManagementMode string `json:"city_management_mode"`
	Status             string `json:"status"`
	UpdatedAt          string `json:"updated_at"`
}

type OwnerIntroResponseRequest struct {
	ChoiceID string `json:"choice_id"`
}

type NarrativeChoice struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type NarrativeEvent struct {
	EventID  string            `json:"event_id"`
	GameID   string            `json:"game_id"`
	Type     string            `json:"type"`
	Subject  string            `json:"subject"`
	Emitter  string            `json:"emitter"`
	Kind     string            `json:"kind"`
	Urgency  string            `json:"urgency"`
	Title    string            `json:"title"`
	Body     string            `json:"body"`
	Metadata map[string]string `json:"metadata,omitempty"`
	Choices  []NarrativeChoice `json:"choices,omitempty"`
}

type NarrativeResponseEvent struct {
	Type      string            `json:"type"`
	Subject   string            `json:"subject"`
	GameID    string            `json:"game_id"`
	EventID   string            `json:"event_id"`
	Choice    NarrativeChoice   `json:"choice"`
	Emitter   string            `json:"emitter"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	Timestamp string            `json:"timestamp"`
}

type MapGenerationProgress struct {
	GameID   string     `json:"game_id"`
	Stage    string     `json:"stage"`
	Progress int        `json:"progress"`
	Message  string     `json:"message"`
	MapData  *MapData   `json:"map_data,omitempty"`
	Stadium  *GridPoint `json:"stadium,omitempty"`
}

type MapData struct {
	Width  int         `json:"width"`
	Height int         `json:"height"`
	Cells  [][]MapCell `json:"cells"`
}

type MapCell struct {
	Terrain string `json:"terrain"`
	Zone    string `json:"zone,omitempty"`
}

type GridPoint struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type MapClientState struct {
	GameID   string     `json:"game_id"`
	Stage    string     `json:"stage"`
	Progress int        `json:"progress"`
	Message  string     `json:"message"`
	MapData  *MapData   `json:"map_data,omitempty"`
	Stadium  *GridPoint `json:"stadium,omitempty"`
}

type MapSnapshotEnvelope struct {
	Type    string         `json:"type"`
	Subject string         `json:"subject"`
	State   MapClientState `json:"state"`
}

type MapPatchEnvelope struct {
	Type    string        `json:"type"`
	Subject string        `json:"subject"`
	GameID  string        `json:"game_id"`
	Patch   MapStatePatch `json:"patch"`
}

type MapStatePatch struct {
	Stage    *string    `json:"stage,omitempty"`
	Progress *int       `json:"progress,omitempty"`
	Message  *string    `json:"message,omitempty"`
	MapData  *MapData   `json:"map_data,omitempty"`
	Stadium  *GridPoint `json:"stadium,omitempty"`
}

const (
	SubjectTimeSessionStarted = "tiempo.sesion_iniciada"
	SubjectTimeSessionEnded   = "tiempo.sesion_terminada"
	SubjectTimeSpeedChanged   = "tiempo.velocidad_cambiada"
	SubjectTimePauseChanged   = "tiempo.pausa_activada"
	SubjectTimeDayAdvanced    = "tiempo.dia_avanzado"
	SubjectMatchFinished      = "partido.terminado"
	SubjectSeasonPatchDelta   = "season.patch"
	SubjectCityPatchDelta     = "city.patch"
	SubjectAgentStateChanged  = "agente.estado_cambio"
	SubjectAgentPatchDelta    = "agent.patch"
)

type EventMeta struct {
	EventID       string `json:"event_id"`
	GameID        string `json:"game_id"`
	OccurredAt    string `json:"occurred_at"`
	SchemaVersion uint16 `json:"schema_version"`
}

type TimeSessionStartedEvent struct {
	EventMeta
	SessionID string `json:"session_id"`
	ClientID  string `json:"client_id"`
}

type TimeSessionEndedEvent struct {
	EventMeta
	SessionID string `json:"session_id"`
	Reason    string `json:"reason"`
}

type TimeSpeedChangedEvent struct {
	EventMeta
	Speed uint8 `json:"speed"`
}

type TimePauseChangedEvent struct {
	EventMeta
	Paused bool `json:"paused"`
}

type TimeDayAdvancedEvent struct {
	EventMeta
	SimulatedDate string `json:"simulated_date"`
	Speed         uint8  `json:"speed"`
	DaysProcessed uint16 `json:"days_processed"`
}

type TimeControlRequest struct {
	Speed  *uint8 `json:"speed,omitempty"`
	Paused *bool  `json:"paused,omitempty"`
}

type TimePatchEnvelope struct {
	Type    string         `json:"type"`
	Subject string         `json:"subject"`
	GameID  string         `json:"game_id"`
	Patch   TimeStatePatch `json:"patch"`
}

type TimeStatePatch struct {
	SimulatedDate *string `json:"simulated_date,omitempty"`
	Speed         *uint8  `json:"speed,omitempty"`
	Paused        *bool   `json:"paused,omitempty"`
	DaysProcessed *uint16 `json:"days_processed,omitempty"`
}

type SeasonPatchEnvelope struct {
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

type CityPatchEnvelope struct {
	Type    string         `json:"type"`
	Subject string         `json:"subject"`
	GameID  string         `json:"game_id"`
	Patch   CityStatePatch `json:"patch"`
}

type CityStatePatch struct {
	FanSentiment             float64 `json:"fan_sentiment"`
	TicketSalesIndex         float64 `json:"ticket_sales_index"`
	LocalEconomyIndex        float64 `json:"local_economy_index"`
	StadiumDistrictLandValue float64 `json:"stadium_district_land_value"`
	WinStreak                uint16  `json:"win_streak"`
	LossStreak               uint16  `json:"loss_streak"`
	LastMatchID              string  `json:"last_match_id"`
	Reason                   string  `json:"reason"`
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

type AgentPatchEnvelope struct {
	Type    string          `json:"type"`
	Subject string          `json:"subject"`
	GameID  string          `json:"game_id"`
	AgentID string          `json:"agent_id"`
	Patch   AgentStatePatch `json:"patch"`
}

type AgentStatePatch struct {
	Mood          string             `json:"mood"`
	State         map[string]float64 `json:"state"`
	Summary       string             `json:"summary"`
	SimulatedDate string             `json:"simulated_date"`
	SourceEventID string             `json:"source_event_id"`
	SourceSubject string             `json:"source_subject"`
}

func AgentPatchFromStateChanged(event AgentStateChangedEvent) AgentPatchEnvelope {
	return AgentPatchEnvelope{
		Type:    SubjectAgentPatchDelta,
		Subject: SubjectAgentPatchDelta,
		GameID:  event.GameID,
		AgentID: event.AgentID,
		Patch: AgentStatePatch{
			Mood:          event.Mood,
			State:         event.State,
			Summary:       event.Summary,
			SimulatedDate: event.SimulatedDate,
			SourceEventID: event.SourceEventID,
			SourceSubject: event.SourceSubject,
		},
	}
}
