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
	"github.com/pulsecity/services/gateway/internal/persistence"
	"github.com/pulsecity/services/gateway/internal/state"
	"github.com/pulsecity/services/gateway/internal/ws"
)

func main() {
	port := envOrDefault("PORT", "8080")
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
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		if err := store.UpsertSnapshot(ctx, currentState); err != nil {
			log.Printf("persist snapshot %s: %v", currentState.GameID, err)
		}
		cancel()

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
		Store:     store,
		Snapshots: snapshots,
	})

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           withCORS(mux),
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

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
