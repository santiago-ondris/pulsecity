package domain

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const (
	SubjectMatchFinished = "partido.terminado"
	OwnTeamID            = "pulsecity"
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
	KeyMoments    []KeyMoment      `json:"key_moments"`
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

type KeyMoment struct {
	Quarter     uint8  `json:"quarter"`
	Clock       string `json:"clock"`
	Kind        string `json:"kind"`
	Description string `json:"description"`
	TeamID      string `json:"team_id"`
	PlayerID    string `json:"player_id,omitempty"`
}

func BuildPostMatchNarrative(event MatchFinishedEvent) NarrativeEvent {
	return BuildPostMatchNarrativeWithContext(event, NarrativeContext{})
}

type NarrativeContext struct {
	WinStreak  uint16
	LossStreak uint16
}

func BuildPostMatchNarrativeWithContext(event MatchFinishedEvent, narrativeContext NarrativeContext) NarrativeEvent {
	context := postMatchContextFromEvent(event)
	context.winStreak = narrativeContext.WinStreak
	context.lossStreak = narrativeContext.LossStreak
	emitter := postMatchEmitter(context)
	urgency := postMatchUrgency(context)
	metadata := map[string]string{
		"match_id":        event.MatchID,
		"source_event_id": event.EventID,
		"source_subject":  SubjectMatchFinished,
		"simulated_date":  event.SimulatedDate,
		"home_team_id":    event.HomeTeam.TeamID,
		"away_team_id":    event.AwayTeam.TeamID,
		"home_score":      fmt.Sprint(event.HomeScore),
		"away_score":      fmt.Sprint(event.AwayScore),
		"winner_team_id":  event.WinnerTeamID,
		"margin":          fmt.Sprint(context.margin),
	}
	if context.winStreak > 0 {
		metadata["win_streak"] = fmt.Sprint(context.winStreak)
	}
	if context.lossStreak > 0 {
		metadata["loss_streak"] = fmt.Sprint(context.lossStreak)
	}

	return NarrativeEvent{
		EventID:  "post-match-" + uuid.NewString(),
		GameID:   event.GameID,
		Type:     "narrative.event",
		Subject:  SubjectNarrativeEventGenerated,
		Emitter:  emitter,
		Kind:     "post_match",
		Urgency:  urgency,
		Title:    postMatchTitle(context),
		Body:     postMatchBody(context),
		Metadata: metadata,
		Choices: []NarrativeChoice{
			{ID: "acknowledge", Label: "Tomar nota"},
		},
	}
}

type postMatchContext struct {
	matchID       string
	ownTeamName   string
	opponentName  string
	ownScore      uint16
	opponentScore uint16
	margin        uint16
	won           bool
	homeGame      bool
	closeGame     bool
	blowout       bool
	keyMoment     string
	topScorerID   string
	topScorerPts  uint16
	winStreak     uint16
	lossStreak    uint16
}

func postMatchContextFromEvent(event MatchFinishedEvent) postMatchContext {
	ownTeam := event.AwayTeam
	opponent := event.HomeTeam
	ownScore := event.AwayScore
	opponentScore := event.HomeScore
	homeGame := false
	if event.HomeTeam.TeamID == OwnTeamID {
		ownTeam = event.HomeTeam
		opponent = event.AwayTeam
		ownScore = event.HomeScore
		opponentScore = event.AwayScore
		homeGame = true
	}

	margin := ownScore
	if opponentScore > ownScore {
		margin = opponentScore - ownScore
	} else {
		margin = ownScore - opponentScore
	}

	topScorerID, topScorerPts := topOwnScorer(event.BoxScore)
	return postMatchContext{
		matchID:       event.MatchID,
		ownTeamName:   fallbackTeamName(ownTeam),
		opponentName:  fallbackTeamName(opponent),
		ownScore:      ownScore,
		opponentScore: opponentScore,
		margin:        margin,
		won:           event.WinnerTeamID == OwnTeamID,
		homeGame:      homeGame,
		closeGame:     margin <= 5,
		blowout:       margin >= 15,
		keyMoment:     firstOwnKeyMoment(event.KeyMoments),
		topScorerID:   topScorerID,
		topScorerPts:  topScorerPts,
	}
}

func postMatchEmitter(context postMatchContext) string {
	if !context.won && context.blowout {
		return "head_coach"
	}
	if context.won && context.homeGame {
		return "owner"
	}
	if context.closeGame {
		return "sports_psychologist"
	}
	return "head_coach"
}

func postMatchUrgency(context postMatchContext) string {
	if !context.won && context.blowout {
		return "high"
	}
	if !context.won {
		return "medium"
	}
	return "low"
}

func postMatchTitle(context postMatchContext) string {
	result := "Victoria"
	if !context.won {
		result = "Derrota"
	}
	if context.blowout {
		return fmt.Sprintf("%s amplia ante %s", result, context.opponentName)
	}
	if context.closeGame {
		return fmt.Sprintf("%s cerrada ante %s", result, context.opponentName)
	}
	return fmt.Sprintf("%s ante %s", result, context.opponentName)
}

func postMatchBody(context postMatchContext) string {
	scoreLine := fmt.Sprintf("%s terminó %d-%d ante %s.", context.ownTeamName, context.ownScore, context.opponentScore, context.opponentName)
	if context.won && context.homeGame {
		scoreLine += " Ganar en casa sostiene el clima interno y le da aire a la ciudad alrededor del estadio."
	} else if context.won {
		scoreLine += " Ganar como visitante fortalece la idea de que el proyecto puede competir fuera de su zona cómoda."
	} else if context.homeGame {
		scoreLine += " Perder en casa enfria el edificio y deja mas preguntas que respuestas para el staff."
	} else {
		scoreLine += " La derrota de visitante no rompe el plan, pero sube la presion del proximo partido."
	}

	if context.topScorerID != "" {
		scoreLine += fmt.Sprintf(" La referencia estadistica fue %s con %d puntos.", context.topScorerID, context.topScorerPts)
	}
	if context.keyMoment != "" {
		scoreLine += " Momento clave: " + context.keyMoment
	}
	if context.winStreak >= 3 {
		scoreLine += fmt.Sprintf(" La racha positiva ya llega a %d victorias y empieza a sentirse fuera de la cancha.", context.winStreak)
	}
	if context.lossStreak >= 3 {
		scoreLine += fmt.Sprintf(" La racha negativa ya llega a %d derrotas y el margen emocional se achica.", context.lossStreak)
	}

	return scoreLine
}

func topOwnScorer(boxScore []PlayerBoxScore) (string, uint16) {
	var playerID string
	var points uint16
	for _, line := range boxScore {
		if line.TeamID != OwnTeamID || line.Points <= points {
			continue
		}
		playerID = line.PlayerID
		points = line.Points
	}

	return playerID, points
}

func firstOwnKeyMoment(moments []KeyMoment) string {
	for _, moment := range moments {
		if moment.TeamID == OwnTeamID && strings.TrimSpace(moment.Description) != "" {
			return moment.Description
		}
	}
	for _, moment := range moments {
		if strings.TrimSpace(moment.Description) != "" {
			return moment.Description
		}
	}

	return ""
}

func fallbackTeamName(team MatchTeam) string {
	if strings.TrimSpace(team.Name) != "" {
		return team.Name
	}
	if strings.TrimSpace(team.Abbreviation) != "" {
		return team.Abbreviation
	}
	return team.TeamID
}
