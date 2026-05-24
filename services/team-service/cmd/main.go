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

		scheduled, found, err := store.DispatchScheduledMatchForDate(ctx, event)
		if err != nil {
			log.Printf("dispatch scheduled match game=%s date=%s: %v", event.GameID, event.SimulatedDate, err)
			return
		}
		if !found {
			return
		}

		if err := bus.PublishJSON(domain.SubjectMatchScheduled, scheduled); err != nil {
			log.Printf("publish scheduled match game=%s match=%s: %v", scheduled.GameID, scheduled.MatchID, err)
			return
		}

		log.Printf("match scheduled game=%s match=%s date=%s", scheduled.GameID, scheduled.MatchID, scheduled.SimulatedDate)
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

		record, changed, err := store.ApplyMatchFinished(ctx, event)
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
