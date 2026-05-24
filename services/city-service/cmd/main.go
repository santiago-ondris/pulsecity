package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pulsecity/services/city-service/internal/domain"
	natsclient "github.com/pulsecity/services/city-service/internal/nats"
	"github.com/pulsecity/services/city-service/internal/persistence"
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

	if _, err := bus.Subscribe(domain.SubjectMatchFinished, func(_ string, data []byte) {
		var event domain.MatchFinishedEvent
		if err := json.Unmarshal(data, &event); err != nil {
			log.Printf("decode match finished event: %v", err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		reaction, changed, err := store.ApplyMatchFinished(ctx, event)
		if err != nil {
			log.Printf("apply city reaction game=%s match=%s: %v", event.GameID, event.MatchID, err)
			return
		}
		if !changed {
			return
		}

		if err := bus.PublishJSON(domain.SubjectCityEconomyChange, reaction.EconomyEvent); err != nil {
			log.Printf("publish city economy change game=%s match=%s: %v", event.GameID, event.MatchID, err)
			return
		}
		if err := bus.PublishJSON(domain.SubjectCityLandUpdated, reaction.LandEvent); err != nil {
			log.Printf("publish city land update game=%s match=%s: %v", event.GameID, event.MatchID, err)
			return
		}
		if err := bus.PublishJSON(domain.SubjectCityPatchDelta, reaction.PatchEvent); err != nil {
			log.Printf("publish city patch game=%s match=%s: %v", event.GameID, event.MatchID, err)
			return
		}

		log.Printf(
			"city updated game=%s match=%s fan=%.1f tickets=%.1f land=%.1f",
			event.GameID,
			event.MatchID,
			reaction.Metrics.FanSentiment,
			reaction.Metrics.TicketSalesIndex,
			reaction.Metrics.StadiumDistrictLandValue,
		)
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
