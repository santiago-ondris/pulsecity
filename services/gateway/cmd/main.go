package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pulsecity/services/gateway/internal/domain"
	httphandlers "github.com/pulsecity/services/gateway/internal/handlers"
	natsclient "github.com/pulsecity/services/gateway/internal/nats"
	"github.com/pulsecity/services/gateway/internal/state"
	"github.com/pulsecity/services/gateway/internal/ws"
)

func main() {
	port := envOrDefault("PORT", "8080")
	natsURL := envOrDefault("NATS_URL", "nats://localhost:4222")

	bus, err := natsclient.New(natsURL)
	if err != nil {
		log.Fatalf("connect nats: %v", err)
	}
	defer bus.Close()

	hub := ws.NewHub()
	snapshots := state.NewMapSnapshots()
	if _, err := bus.Subscribe(domain.SubjectMapWildcard, func(subject string, data []byte) {
		if subject == domain.SubjectMapGenerationStarted {
			return
		}

		var progress domain.MapGenerationProgress
		if err := json.Unmarshal(data, &progress); err != nil {
			log.Printf("decode map event %s: %v", subject, err)
			return
		}

		currentState, existed := snapshots.ApplyProgress(progress)
		if !existed && currentState.MapData != nil {
			hub.Broadcast(domain.MapSnapshotEnvelope{
				Type:    "map.snapshot",
				Subject: subject,
				State:   currentState,
			})
			return
		}

		stage := currentState.Stage
		progressValue := currentState.Progress
		message := currentState.Message
		hub.Broadcast(domain.MapPatchEnvelope{
			Type:    "map.patch",
			Subject: subject,
			GameID:  progress.GameID,
			Patch: domain.MapStatePatch{
				Stage:    &stage,
				Progress: &progressValue,
				Message:  &message,
				MapData:  progress.MapData,
				Stadium:  progress.Stadium,
			},
		})
	}); err != nil {
		log.Fatalf("subscribe map events: %v", err)
	}

	mux := http.NewServeMux()
	httphandlers.RegisterRoutes(mux, httphandlers.Dependencies{
		Bus:       bus,
		Hub:       hub,
		Snapshots: snapshots,
	})

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("gateway listening on http://localhost:%s", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("serve http: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("shutdown http: %v", err)
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
