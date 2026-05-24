package domain

import "testing"

func TestApplyMatchFinishedHomeWinImprovesCityMetrics(t *testing.T) {
	event := matchFinishedEvent("match-001", true, true)

	reaction := ApplyMatchFinished(DefaultCityMetrics("game-1"), event)

	if reaction.Metrics.FanSentiment != 58 {
		t.Errorf("ApplyMatchFinished(%q) fan sentiment = %.1f, want 58.0", event.MatchID, reaction.Metrics.FanSentiment)
	}
	if reaction.Metrics.TicketSalesIndex != 57 {
		t.Errorf("ApplyMatchFinished(%q) ticket sales = %.1f, want 57.0", event.MatchID, reaction.Metrics.TicketSalesIndex)
	}
	if reaction.Metrics.LocalEconomyIndex != 53.5 {
		t.Errorf("ApplyMatchFinished(%q) economy = %.1f, want 53.5", event.MatchID, reaction.Metrics.LocalEconomyIndex)
	}
	if reaction.Metrics.WinStreak != 1 || reaction.Metrics.LossStreak != 0 {
		t.Errorf("ApplyMatchFinished(%q) streak = %d-%d, want 1-0", event.MatchID, reaction.Metrics.WinStreak, reaction.Metrics.LossStreak)
	}
	if reaction.EconomyEvent.Reason != "home_win" {
		t.Errorf("ApplyMatchFinished(%q) reason = %q, want home_win", event.MatchID, reaction.EconomyEvent.Reason)
	}
}

func TestApplyMatchFinishedWinningStreakRaisesStadiumLandValue(t *testing.T) {
	current := DefaultCityMetrics("game-1")
	current.WinStreak = 2
	event := matchFinishedEvent("match-003", true, false)

	reaction := ApplyMatchFinished(current, event)

	if reaction.Metrics.WinStreak != 3 {
		t.Errorf("ApplyMatchFinished(%q) win streak = %d, want 3", event.MatchID, reaction.Metrics.WinStreak)
	}
	if reaction.LandEvent.LandValueDelta != 1.5 {
		t.Errorf("ApplyMatchFinished(%q) land delta = %.1f, want 1.5", event.MatchID, reaction.LandEvent.LandValueDelta)
	}
	if reaction.Metrics.StadiumDistrictLandValue != 101.5 {
		t.Errorf("ApplyMatchFinished(%q) land value = %.1f, want 101.5", event.MatchID, reaction.Metrics.StadiumDistrictLandValue)
	}
	if reaction.PatchEvent.Patch.Reason != "away_win_winning_streak" {
		t.Errorf("ApplyMatchFinished(%q) patch reason = %q, want away_win_winning_streak", event.MatchID, reaction.PatchEvent.Patch.Reason)
	}
}

func TestApplyMatchFinishedLosingStreakCoolsLocalEconomy(t *testing.T) {
	current := DefaultCityMetrics("game-1")
	current.LossStreak = 2
	event := matchFinishedEvent("match-004", false, true)

	reaction := ApplyMatchFinished(current, event)

	if reaction.Metrics.LossStreak != 3 {
		t.Errorf("ApplyMatchFinished(%q) loss streak = %d, want 3", event.MatchID, reaction.Metrics.LossStreak)
	}
	if reaction.LandEvent.LandValueDelta != -1.2 {
		t.Errorf("ApplyMatchFinished(%q) land delta = %.1f, want -1.2", event.MatchID, reaction.LandEvent.LandValueDelta)
	}
	if reaction.Metrics.LocalEconomyIndex != 47.5 {
		t.Errorf("ApplyMatchFinished(%q) economy = %.1f, want 47.5", event.MatchID, reaction.Metrics.LocalEconomyIndex)
	}
	if reaction.PatchEvent.Patch.Reason != "home_loss_losing_streak" {
		t.Errorf("ApplyMatchFinished(%q) patch reason = %q, want home_loss_losing_streak", event.MatchID, reaction.PatchEvent.Patch.Reason)
	}
}

func matchFinishedEvent(matchID string, ownTeamWon bool, homeGame bool) MatchFinishedEvent {
	homeTeam := MatchTeam{TeamID: "rival-1", Name: "Rival One", Abbreviation: "RIV"}
	awayTeam := MatchTeam{TeamID: OwnTeamID, Name: "PulseCity", Abbreviation: "PUL"}
	if homeGame {
		homeTeam, awayTeam = awayTeam, homeTeam
	}

	winnerTeamID := "rival-1"
	if ownTeamWon {
		winnerTeamID = OwnTeamID
	}

	return MatchFinishedEvent{
		EventMeta: EventMeta{
			EventID:       "event-" + matchID,
			GameID:        "game-1",
			OccurredAt:    "2026-10-22T00:00:00Z",
			SchemaVersion: 1,
		},
		MatchID:       matchID,
		SimulatedDate: "2026-10-22",
		HomeTeam:      homeTeam,
		AwayTeam:      awayTeam,
		HomeScore:     101,
		AwayScore:     99,
		WinnerTeamID:  winnerTeamID,
		Seed:          1,
	}
}
