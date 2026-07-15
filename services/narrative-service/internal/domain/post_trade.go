package domain

import "fmt"

const SubjectTradeAccepted = "trade.aceptada"

type TradeAcceptedEvent struct {
	EventMeta
	ProposalID              string `json:"proposal_id"`
	SimulatedDate           string `json:"simulated_date"`
	RivalTeamID             string `json:"rival_team_id"`
	OutgoingPlayerID        string `json:"outgoing_player_id"`
	OutgoingPlayerName      string `json:"outgoing_player_name"`
	IncomingPlayerID        string `json:"incoming_player_id"`
	IncomingPlayerName      string `json:"incoming_player_name"`
	IncomingPosition        string `json:"incoming_position"`
	IncomingRating          uint8  `json:"incoming_rating"`
	IncomingSalary          int    `json:"incoming_salary"`
	AcceptedAdditionalAsset string `json:"accepted_additional_asset,omitempty"`
}

func BuildPostTradeNarrative(event TradeAcceptedEvent) NarrativeEvent {
	assetLine := "sin sumar activos extra."
	if event.AcceptedAdditionalAsset != "" {
		assetLine = "incluyendo " + event.AcceptedAdditionalAsset + " como concesion adicional."
	}

	return NarrativeEvent{
		EventID: "post-trade-" + event.ProposalID,
		GameID:  event.GameID,
		Type:    "narrative.event",
		Subject: SubjectNarrativeEventGenerated,
		Emitter: "director_player_personnel",
		Kind:    "post_trade",
		Urgency: "normal",
		Title:   "Trade cerrado",
		Body: fmt.Sprintf(
			"Player Personnel confirma el cierre: %s sale de PulseCity y %s llega para cubrir %s, %s La lectura interna es simple: el roster cambia hoy, y el vestuario va a necesitar unas horas para absorberlo.",
			event.OutgoingPlayerName,
			event.IncomingPlayerName,
			event.IncomingPosition,
			assetLine,
		),
		Metadata: map[string]string{
			"proposal_id":          event.ProposalID,
			"source_event_id":      event.EventID,
			"source_subject":       SubjectTradeAccepted,
			"simulated_date":       event.SimulatedDate,
			"rival_team_id":        event.RivalTeamID,
			"outgoing_player_id":   event.OutgoingPlayerID,
			"incoming_player_id":   event.IncomingPlayerID,
			"incoming_position":    event.IncomingPosition,
			"incoming_rating":      fmt.Sprint(event.IncomingRating),
			"incoming_salary":      fmt.Sprint(event.IncomingSalary),
			"accepted_extra_asset": event.AcceptedAdditionalAsset,
		},
		Choices: []NarrativeChoice{
			{ID: "acknowledge", Label: "Tomar nota"},
		},
	}
}
