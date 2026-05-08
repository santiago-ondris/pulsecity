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
	CityName string `json:"city_name"`
}

type MapGenerationRequest struct {
	GameID   string `json:"game_id"`
	CityName string `json:"city_name,omitempty"`
}

type MapGenerationProgress struct {
	GameID   string `json:"game_id"`
	Stage    string `json:"stage"`
	Progress int    `json:"progress"`
	Message  string `json:"message"`
}

type MapDeltaEnvelope struct {
	Type    string                `json:"type"`
	Subject string                `json:"subject"`
	Payload MapGenerationProgress `json:"payload"`
}
