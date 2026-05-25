package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pulsecity/services/analytics-service/internal/domain"
	natsclient "github.com/pulsecity/services/analytics-service/internal/nats"
	"github.com/pulsecity/services/analytics-service/internal/persistence"
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
			log.Printf("decode partido.terminado: %v", err)
			return
		}

		go ingest("partido.terminado", func(ctx context.Context) error {
			return store.IngestMatchFinished(ctx, event)
		})
	}); err != nil {
		return err
	}

	if _, err := bus.Subscribe(domain.SubjectCityEconomyChange, func(_ string, data []byte) {
		var event domain.CityEconomyChangeEvent
		if err := json.Unmarshal(data, &event); err != nil {
			log.Printf("decode ciudad.economia_cambio: %v", err)
			return
		}

		go ingest("ciudad.economia_cambio", func(ctx context.Context) error {
			return store.IngestCityEconomyChange(ctx, event)
		})
	}); err != nil {
		return err
	}

	if _, err := bus.Subscribe(domain.SubjectCityLandUpdated, func(_ string, data []byte) {
		var event domain.CityLandUpdatedEvent
		if err := json.Unmarshal(data, &event); err != nil {
			log.Printf("decode ciudad.suelo_actualizado: %v", err)
			return
		}

		go ingest("ciudad.suelo_actualizado", func(ctx context.Context) error {
			return store.IngestCityLandUpdated(ctx, event)
		})
	}); err != nil {
		return err
	}

	if _, err := bus.Subscribe(domain.SubjectAgentStateChanged, func(_ string, data []byte) {
		var event domain.AgentStateChangedEvent
		if err := json.Unmarshal(data, &event); err != nil {
			log.Printf("decode agente.estado_cambio: %v", err)
			return
		}

		go ingest("agente.estado_cambio", func(ctx context.Context) error {
			return store.IngestAgentStateChanged(ctx, event)
		})
	}); err != nil {
		return err
	}

	log.Printf("%s connected to nats at %s", domain.ServiceName, natsURL)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	return nil
}

func ingest(subject string, fn func(context.Context) error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := fn(ctx); err != nil {
		log.Printf("ingest analytics %s: %v", subject, err)
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
