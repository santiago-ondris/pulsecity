package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pulsecity/services/team-service/internal/domain"
	natsclient "github.com/pulsecity/services/team-service/internal/nats"
	"github.com/pulsecity/services/team-service/internal/persistence"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	natsURL := envOrDefault("NATS_URL", "nats://localhost:4222")
	databaseURL := envOrDefault("DATABASE_URL", "postgres://pulsecity:pulsecity@localhost:5433/pulsecity_dev?sslmode=disable")

	bus, err := natsclient.New(natsURL)
	if err != nil {
		return err
	}
	defer bus.Close()

	store, err := persistence.NewStore(context.Background(), databaseURL)
	if err != nil {
		return err
	}
	defer store.Close()

	if err := store.EnsureSchema(context.Background()); err != nil {
		return err
	}

	if _, err := bus.Subscribe(domain.SubjectMapGenerationStarted, func(_ string, data []byte) {
		var event domain.GameStartedEvent
		if err := json.Unmarshal(data, &event); err != nil {
			log.Printf("decode game start event: %v", err)
			return
		}

		season, err := domain.GenerateInitialSeason(event)
		if err != nil {
			log.Printf("generate initial season %s: %v", event.GameID, err)
			return
		}

		if err := store.SaveInitialSeason(context.Background(), season); err != nil {
			log.Printf("save initial season %s: %v", event.GameID, err)
			return
		}

		log.Printf("initial season ready game=%s roster=%d opponents=%d schedule=%d", season.GameID, len(season.Roster), len(season.Opponents), len(season.Schedule))
	}); err != nil {
		return err
	}

	if _, err := bus.Subscribe(domain.SubjectTimeDayAdvanced, func(_ string, data []byte) {
		var event domain.DayAdvancedEvent
		if err := json.Unmarshal(data, &event); err != nil {
			log.Printf("decode day advanced event: %v", err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		days, err := domain.CoveredDayAdvancedEvents(event)
		if err != nil {
			log.Printf("expand advanced days game=%s date=%s: %v", event.GameID, event.SimulatedDate, err)
			return
		}
		for _, day := range days {
			recovered, err := store.RecoverPlayersForDate(ctx, day)
			if err != nil {
				log.Printf("recover players game=%s date=%s: %v", day.GameID, day.SimulatedDate, err)
				return
			}
			for _, event := range recovered {
				if err := bus.PublishJSON(domain.SubjectPlayerRecovered, event); err != nil {
					log.Printf("publish recovered player game=%s player=%s: %v", event.GameID, event.PlayerID, err)
					return
				}
				log.Printf("player recovered game=%s player=%s injury=%s", event.GameID, event.PlayerID, event.InjuryID)
			}

			scheduled, found, err := store.DispatchScheduledMatchForDate(ctx, day)
			if err != nil {
				log.Printf("dispatch scheduled match game=%s date=%s: %v", day.GameID, day.SimulatedDate, err)
				return
			}
			if !found {
				continue
			}

			if err := bus.PublishJSON(domain.SubjectMatchScheduled, scheduled); err != nil {
				log.Printf("publish scheduled match game=%s match=%s: %v", scheduled.GameID, scheduled.MatchID, err)
				return
			}

			log.Printf("match scheduled game=%s match=%s date=%s", scheduled.GameID, scheduled.MatchID, scheduled.SimulatedDate)
		}
	}); err != nil {
		return err
	}

	if _, err := bus.Subscribe(domain.SubjectMatchFinished, func(_ string, data []byte) {
		var event domain.MatchFinishedEvent
		if err := json.Unmarshal(data, &event); err != nil {
			log.Printf("decode match finished event: %v", err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		record, injuries, changed, err := store.ApplyMatchFinished(ctx, event)
		if err != nil {
			log.Printf("apply match finished game=%s match=%s: %v", event.GameID, event.MatchID, err)
			return
		}
		if !changed {
			return
		}

		patch := domain.SeasonPatchFromRecord(record)
		if err := bus.PublishJSON(domain.SubjectSeasonPatchDelta, patch); err != nil {
			log.Printf("publish season patch game=%s match=%s: %v", event.GameID, event.MatchID, err)
			return
		}

		log.Printf("record updated game=%s record=%d-%d match=%s", record.GameID, record.Wins, record.Losses, event.MatchID)
		for _, injury := range injuries {
			if err := bus.PublishJSON(domain.SubjectPlayerInjured, injury); err != nil {
				log.Printf("publish injured player game=%s player=%s: %v", injury.GameID, injury.PlayerID, err)
				return
			}
			log.Printf("player injured game=%s player=%s severity=%s days=%d", injury.GameID, injury.PlayerID, injury.Severity, injury.EstimatedDaysOut)
		}
	}); err != nil {
		return err
	}

	if _, err := bus.Subscribe(domain.SubjectRosterPatchDelta, func(_ string, data []byte) {
		var event domain.RosterPatchEnvelope
		if err := json.Unmarshal(data, &event); err != nil {
			log.Printf("decode roster patch event: %v", err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := store.ApplyRosterPatch(ctx, event); err != nil {
			log.Printf("apply roster patch game=%s source=%s: %v", event.GameID, event.Patch.SourceEventID, err)
			return
		}

		log.Printf("roster match state updated game=%s players=%d source=%s", event.GameID, len(event.Patch.Players), event.Patch.SourceEventID)
	}); err != nil {
		return err
	}

	if _, err := bus.Subscribe(domain.SubjectGMDecision, func(_ string, data []byte) {
		var event domain.GMDecisionRegisteredEvent
		if err := json.Unmarshal(data, &event); err != nil {
			log.Printf("decode gm decision event: %v", err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		recovered, err := store.ApplyMedicalDecision(ctx, event)
		if err != nil {
			log.Printf("apply medical decision game=%s decision=%s: %v", event.GameID, event.DecisionID, err)
			return
		}
		for _, recoveredEvent := range recovered {
			if err := bus.PublishJSON(domain.SubjectPlayerRecovered, recoveredEvent); err != nil {
				log.Printf("publish forced recovered player game=%s player=%s: %v", recoveredEvent.GameID, recoveredEvent.PlayerID, err)
				return
			}
		}
	}); err != nil {
		return err
	}

	log.Printf("%s connected to nats at %s", domain.ServiceName, natsURL)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	return nil
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
