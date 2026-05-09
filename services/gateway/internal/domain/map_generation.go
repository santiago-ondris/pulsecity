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
	CityName        string `json:"city_name"`
	FranchiseName   string `json:"franchise_name"`
	Abbreviation    string `json:"abbreviation"`
	PrimaryColor    string `json:"primary_color"`
	SecondaryColor  string `json:"secondary_color"`
	AccentColor     string `json:"accent_color"`
	InitialScenario string `json:"initial_scenario"`
}

type MapGenerationRequest struct {
	GameID   string `json:"game_id"`
	CityName string `json:"city_name,omitempty"`
}

type GameSetup struct {
	GameID          string `json:"game_id"`
	CityName        string `json:"city_name"`
	FranchiseName   string `json:"franchise_name"`
	Abbreviation    string `json:"abbreviation"`
	PrimaryColor    string `json:"primary_color"`
	SecondaryColor  string `json:"secondary_color"`
	AccentColor     string `json:"accent_color"`
	InitialScenario string `json:"initial_scenario"`
	Status          string `json:"status"`
	CreatedAt       string `json:"created_at,omitempty"`
	UpdatedAt       string `json:"updated_at,omitempty"`
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
