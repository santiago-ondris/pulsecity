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

	if err := store.EnsureSchema(context.Background()); err != nil {
		log.Fatalf("ensure postgres schema: %v", err)
	}

	if _, err := bus.Subscribe(domain.SubjectNarrativeOwnerIntroRequested, func(_ string, data []byte) {
		var request domain.OwnerIntroRequestedEvent
		if err := json.Unmarshal(data, &request); err != nil {
			log.Printf("decode owner intro request: %v", err)
			return
		}

		go processOwnerIntroRequest(bus, store, request)
	}); err != nil {
		log.Fatalf("subscribe owner intro requests: %v", err)
	}

	if _, err := bus.Subscribe(domain.SubjectMatchFinished, func(_ string, data []byte) {
		var event domain.MatchFinishedEvent
		if err := json.Unmarshal(data, &event); err != nil {
			log.Printf("decode match finished: %v", err)
			return
		}

		go processMatchFinished(bus, store, event)
	}); err != nil {
		log.Fatalf("subscribe match finished events: %v", err)
	}

	if _, err := bus.Subscribe(domain.SubjectAgentConsultationStarted, func(_ string, data []byte) {
		var event domain.AgentConsultationStartedEvent
		if err := json.Unmarshal(data, &event); err != nil {
			log.Printf("decode agent consultation: %v", err)
			return
		}

		go processAgentConsultation(bus, store, event)
	}); err != nil {
		log.Fatalf("subscribe agent consultation events: %v", err)
	}

	log.Printf(
		"narrative-service listening on %s, %s and %s",
		domain.SubjectNarrativeOwnerIntroRequested,
		domain.SubjectMatchFinished,
		domain.SubjectAgentConsultationStarted,
	)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}

func processOwnerIntroRequest(bus *natsclient.Client, store *persistence.Store, request domain.OwnerIntroRequestedEvent) {
	time.Sleep(randomNarrativeDelay())

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

func processMatchFinished(bus *natsclient.Client, store *persistence.Store, match domain.MatchFinishedEvent) {
	time.Sleep(deterministicNarrativeDelay(match.EventID, match.MatchID))

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	narrativeContext, err := store.LoadNarrativeContext(ctx, match.GameID)
	if err != nil {
		log.Printf("load narrative context game=%s match=%s: %v", match.GameID, match.MatchID, err)
	}

	event := domain.BuildPostMatchNarrativeWithContext(match, narrativeContext)

	stored, err := store.InsertNarrativeEventIfNew(ctx, event)
	if err != nil {
		log.Printf("persist post-match narrative game=%s match=%s: %v", match.GameID, match.MatchID, err)
		return
	}
	if !stored {
		return
	}

	if err := bus.PublishJSON(domain.SubjectNarrativeEventGenerated, event); err != nil {
		log.Printf("publish post-match narrative game=%s match=%s: %v", match.GameID, match.MatchID, err)
	}
}

func processAgentConsultation(bus *natsclient.Client, store *persistence.Store, event domain.AgentConsultationStartedEvent) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	chatContext, err := store.LoadAgentChatContext(ctx, event.GameID, event.AgentID)
	if err != nil {
		log.Printf("load agent chat context game=%s agent=%s: %v", event.GameID, event.AgentID, err)
		return
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)
	userMessage := domain.BuildUserChatMessage(event, now)
	if _, err := store.InsertChatMessageIfNew(ctx, userMessage, chatContext); err != nil {
		log.Printf("persist user chat message game=%s conversation=%s: %v", event.GameID, event.ConversationID, err)
		return
	}

	agentMessage := domain.BuildStubAgentChatMessage(event, chatContext, time.Now().UTC().Format(time.RFC3339Nano))
	stored, err := store.InsertChatMessageIfNew(ctx, agentMessage, chatContext)
	if err != nil {
		log.Printf("persist agent chat message game=%s conversation=%s: %v", event.GameID, event.ConversationID, err)
		return
	}
	if !stored {
		return
	}

	if err := bus.PublishJSON(domain.SubjectChatMessageDelta, domain.ChatEnvelopeFromMessage(agentMessage)); err != nil {
		log.Printf("publish chat message game=%s conversation=%s: %v", event.GameID, event.ConversationID, err)
	}
}

func randomNarrativeDelay() time.Duration {
	return time.Duration(250+rand.Intn(251)) * time.Millisecond
}

func deterministicNarrativeDelay(parts ...string) time.Duration {
	var total uint16
	for _, part := range parts {
		for _, character := range part {
			total += uint16(character)
		}
	}

	return time.Duration(250+int(total%251)) * time.Millisecond
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
