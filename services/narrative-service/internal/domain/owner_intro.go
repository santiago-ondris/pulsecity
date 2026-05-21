package domain

import (
	"strings"

	"github.com/google/uuid"
)

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

func BuildOwnerIntroEvent(request OwnerIntroRequestedEvent) NarrativeEvent {
	return NarrativeEvent{
		EventID: "owner-intro-" + uuid.NewString(),
		GameID:  request.GameID,
		Type:    "narrative.event",
		Subject: "narrativa.owner_intro_generado",
		Emitter: "owner",
		Kind:    "owner_intro",
		Urgency: "critical",
		Title:   "Llamada del Owner",
		Body:    ownerIntroBody(request),
		Metadata: map[string]string{
			"city_name":            request.CityName,
			"franchise_name":       request.FranchiseName,
			"initial_scenario":     request.InitialScenario,
			"city_management_mode": request.CityManagementMode,
		},
		Choices: []NarrativeChoice{
			{ID: "build_culture", Label: "Empezá por identidad y cultura"},
			{ID: "win_now", Label: "Acelerá para competir rapido"},
			{ID: "city_first", Label: "Usá la franquicia para activar la ciudad"},
		},
	}
}

func ownerIntroBody(request OwnerIntroRequestedEvent) string {
	franchise := strings.TrimSpace(request.FranchiseName)
	city := strings.TrimSpace(request.CityName)
	modeLine := "Quiero que entiendas algo desde el primer dia: esta franquicia y esta ciudad van a empujarse entre si."
	if request.CityManagementMode == "dual_figure" {
		modeLine = "Vas a llevar dos sombreros desde el primer dia: la franquicia y la ciudad. No quiero que uses ese poder a medias."
	}

	switch request.InitialScenario {
	case "rebuild":
		return "Te traje a " + city + " para construir con paciencia. " + franchise + " no necesita humo, necesita direccion. " + modeLine + " Si hacemos bien las bases, el resto llega."
	case "contention":
		return "No te contraté para aprender en el cargo. " + franchise + " tiene talento, gasto y expectativas encima desde el primer dia. " + modeLine + " Quiero resultados rapido y no pienso disfrazarlo."
	case "decline":
		return "Acá todavía pesa demasiado el pasado. " + franchise + " vive bajo la sombra de lo que fue, y la ciudad lo siente. " + modeLine + " Tu trabajo es devolver autoridad antes de que esto se vuelva costumbre."
	default:
		return city + " no tiene historia previa que la sostenga; eso es una ventaja y una responsabilidad. " + franchise + " empieza desde cero, y con eso viene libertad real. " + modeLine + " Quiero vision, criterio y una identidad que se note desde el primer movimiento."
	}
}
