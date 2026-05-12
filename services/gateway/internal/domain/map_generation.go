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

type MapGenerationRequest struct {
	GameID   string `json:"game_id"`
	CityName string `json:"city_name,omitempty"`
}

type GameSetup struct {
	GameID             string           `json:"game_id"`
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
