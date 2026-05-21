package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pulsecity/services/narrative-service/internal/domain"
	natsclient "github.com/pulsecity/services/narrative-service/internal/nats"
	"github.com/pulsecity/services/narrative-service/internal/persistence"
)

func main() {
	natsURL := envOrDefault("NATS_URL", "nats://localhost:4222")
	databaseURL := envOrDefault("DATABASE_URL", "postgres://pulsecity:pulsecity@localhost:5433/pulsecity_dev?sslmode=disable")

	bus, err := natsclient.New(natsURL)
	if err != nil {
		log.Fatalf("connect nats: %v", err)
	}
	defer bus.Close()

	store, err := persistence.NewStore(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("connect postgres: %v", err)
	}
	defer store.Close()

	if _, err := bus.Subscribe(domain.SubjectNarrativeOwnerIntroRequested, func(_ string, data []byte) {
		var request domain.OwnerIntroRequestedEvent
		if err := json.Unmarshal(data, &request); err != nil {
			log.Printf("decode owner intro request: %v", err)
			return
		}

		processOwnerIntroRequest(bus, store, request)
	}); err != nil {
		log.Fatalf("subscribe owner intro requests: %v", err)
	}

	log.Printf("narrative-service listening on %s", domain.SubjectNarrativeOwnerIntroRequested)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}

func processOwnerIntroRequest(bus *natsclient.Client, store *persistence.Store, request domain.OwnerIntroRequestedEvent) {
	delay := time.Duration(250+rand.Intn(251)) * time.Millisecond
	time.Sleep(delay)

	event := domain.BuildOwnerIntroEvent(request)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	stored, err := store.SetOwnerIntroEventIfEmpty(ctx, request.GameID, event)
	if err != nil {
		log.Printf("persist owner intro event %s: %v", request.GameID, err)
		return
	}
	if !stored {
		return
	}

	if err := bus.PublishJSON(domain.SubjectNarrativeEventGenerated, event); err != nil {
		log.Printf("publish narrative event %s: %v", request.GameID, err)
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
