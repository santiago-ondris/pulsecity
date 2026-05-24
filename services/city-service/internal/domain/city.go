package domain

import (
	"fmt"
	"time"
)

const (
	OwnTeamID             = "pulsecity"
	StadiumDistrictZoneID = "stadium_district"
)

type CityMetrics struct {
	GameID                   string
	FanSentiment             float64
	TicketSalesIndex         float64
	LocalEconomyIndex        float64
	StadiumDistrictLandValue float64
	WinStreak                uint16
	LossStreak               uint16
	LastMatchID              string
}

type CityReaction struct {
	Metrics      CityMetrics
	EconomyEvent CityEconomyChangeEvent
	LandEvent    CityLandUpdatedEvent
	PatchEvent   CityPatchEvent
}

func DefaultCityMetrics(gameID string) CityMetrics {
	return CityMetrics{
		GameID:                   gameID,
		FanSentiment:             50,
		TicketSalesIndex:         50,
		LocalEconomyIndex:        50,
		StadiumDistrictLandValue: 100,
	}
}

func ApplyMatchFinished(current CityMetrics, event MatchFinishedEvent) CityReaction {
	if current.GameID == "" {
		current = DefaultCityMetrics(event.GameID)
	}

	won := event.WinnerTeamID == OwnTeamID
	homeGame := event.HomeTeam.TeamID == OwnTeamID
	reason := resultReason(won, homeGame)

	fanDelta, ticketDelta, economyDelta := cityDeltas(won, homeGame)
	next := current
	next.FanSentiment = clampFloat(current.FanSentiment+fanDelta, 0, 100)
	next.TicketSalesIndex = clampFloat(current.TicketSalesIndex+ticketDelta, 0, 150)
	next.LocalEconomyIndex = clampFloat(current.LocalEconomyIndex+economyDelta, 0, 150)
	next.LastMatchID = event.MatchID

	if won {
		next.WinStreak = current.WinStreak + 1
		next.LossStreak = 0
	} else {
		next.LossStreak = current.LossStreak + 1
		next.WinStreak = 0
	}

	landDelta := landValueDelta(next.WinStreak, next.LossStreak)
	if landDelta != 0 {
		reason = streakReason(reason, won)
	}
	next.StadiumDistrictLandValue = clampFloat(current.StadiumDistrictLandValue+landDelta, 40, 180)

	occurredAt := event.OccurredAt
	if occurredAt == "" {
		occurredAt = time.Now().UTC().Format(time.RFC3339)
	}
	sourceEventID := event.EventID
	if sourceEventID == "" {
		sourceEventID = fmt.Sprintf("match-finished-%s", event.MatchID)
	}

	reaction := CityReaction{
		Metrics: next,
		EconomyEvent: CityEconomyChangeEvent{
			EventMeta: EventMeta{
				EventID:       fmt.Sprintf("city-economy-%s", event.MatchID),
				GameID:        event.GameID,
				OccurredAt:    occurredAt,
				SchemaVersion: 1,
			},
			SimulatedDate:     event.SimulatedDate,
			SourceEventID:     sourceEventID,
			SourceSubject:     SubjectMatchFinished,
			FanSentimentDelta: fanDelta,
			TicketSalesDelta:  ticketDelta,
			LocalEconomyDelta: economyDelta,
			FanSentiment:      next.FanSentiment,
			TicketSalesIndex:  next.TicketSalesIndex,
			LocalEconomyIndex: next.LocalEconomyIndex,
			WinStreak:         next.WinStreak,
			LossStreak:        next.LossStreak,
			Reason:            reason,
		},
		LandEvent: CityLandUpdatedEvent{
			EventMeta: EventMeta{
				EventID:       fmt.Sprintf("city-land-%s", event.MatchID),
				GameID:        event.GameID,
				OccurredAt:    occurredAt,
				SchemaVersion: 1,
			},
			SimulatedDate:  event.SimulatedDate,
			ZoneID:         StadiumDistrictZoneID,
			LandValueDelta: landDelta,
			NewLandValue:   next.StadiumDistrictLandValue,
			SourceEventID:  sourceEventID,
			Reason:         reason,
		},
	}
	reaction.PatchEvent = CityPatchFromReaction(reaction)

	return reaction
}

func CityPatchFromReaction(reaction CityReaction) CityPatchEvent {
	return CityPatchEvent{
		Type:    SubjectCityPatchDelta,
		Subject: SubjectCityPatchDelta,
		GameID:  reaction.Metrics.GameID,
		Patch: CityStatePatch{
			FanSentiment:             reaction.Metrics.FanSentiment,
			TicketSalesIndex:         reaction.Metrics.TicketSalesIndex,
			LocalEconomyIndex:        reaction.Metrics.LocalEconomyIndex,
			StadiumDistrictLandValue: reaction.Metrics.StadiumDistrictLandValue,
			WinStreak:                reaction.Metrics.WinStreak,
			LossStreak:               reaction.Metrics.LossStreak,
			LastMatchID:              reaction.Metrics.LastMatchID,
			Reason:                   reaction.EconomyEvent.Reason,
		},
	}
}

func cityDeltas(won, homeGame bool) (float64, float64, float64) {
	if won && homeGame {
		return 8, 7, 3.5
	}
	if won {
		return 6, 4, 1.5
	}
	if homeGame {
		return -5, -5, -2.5
	}

	return -4, -3, -1.5
}

func landValueDelta(winStreak, lossStreak uint16) float64 {
	if winStreak >= 3 {
		return 1.5 + float64(winStreak-3)*0.4
	}
	if lossStreak >= 3 {
		return -1.2 - float64(lossStreak-3)*0.3
	}

	return 0
}

func resultReason(won, homeGame bool) string {
	if won && homeGame {
		return "home_win"
	}
	if won {
		return "away_win"
	}
	if homeGame {
		return "home_loss"
	}

	return "away_loss"
}

func streakReason(base string, won bool) string {
	if won {
		return base + "_winning_streak"
	}

	return base + "_losing_streak"
}

func clampFloat(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}

	return value
}
