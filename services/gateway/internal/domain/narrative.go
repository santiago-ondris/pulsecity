package domain

const (
	SubjectNarrativeOwnerIntroRequested = "narrativa.owner_intro_solicitada"
	SubjectNarrativeEventGenerated      = "narrativa.evento_generado"
)

type OwnerIntroRequestedEvent struct {
	GameID             string `json:"game_id"`
	CityName           string `json:"city_name"`
	FranchiseName      string `json:"franchise_name"`
	InitialScenario    string `json:"initial_scenario"`
	CityManagementMode string `json:"city_management_mode"`
}
