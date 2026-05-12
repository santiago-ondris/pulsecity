package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"
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
		} else {
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
		}

		if currentState.Stage != "complete" {
			return
		}

		game, found, err := store.GetGame(context.Background(), currentState.GameID)
		if err != nil {
			log.Printf("load game setup %s: %v", currentState.GameID, err)
			return
		}
		if !found || game.OwnerIntroEvent != nil {
			return
		}

		event := buildOwnerIntroEvent(game)
		ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
		if err := store.SetOwnerIntroEvent(ctx, currentState.GameID, event); err != nil {
			log.Printf("persist owner intro event %s: %v", currentState.GameID, err)
			cancel()
			return
		}
		cancel()

		hub.Broadcast(event)
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

func buildOwnerIntroEvent(game domain.GameSetup) domain.NarrativeEvent {
	title := "Llamada del Owner"
	body := ownerIntroBody(game)

	return domain.NarrativeEvent{
		EventID: "owner-intro-" + uuid.NewString(),
		GameID:  game.GameID,
		Type:    "narrative.event",
		Subject: "narrativa.owner_intro_generado",
		Emitter: "owner",
		Kind:    "owner_intro",
		Urgency: "critical",
		Title:   title,
		Body:    body,
		Metadata: map[string]string{
			"city_name":            game.CityName,
			"franchise_name":       game.FranchiseName,
			"initial_scenario":     game.InitialScenario,
			"city_management_mode": game.CityManagementMode,
		},
		Choices: []domain.NarrativeChoice{
			{ID: "build_culture", Label: "Empezá por identidad y cultura"},
			{ID: "win_now", Label: "Acelerá para competir rapido"},
			{ID: "city_first", Label: "Usá la franquicia para activar la ciudad"},
		},
	}
}

func ownerIntroBody(game domain.GameSetup) string {
	franchise := strings.TrimSpace(game.FranchiseName)
	city := strings.TrimSpace(game.CityName)
	modeLine := "Quiero que entiendas algo desde el primer dia: esta franquicia y esta ciudad van a empujarse entre si."
	if game.CityManagementMode == "dual_figure" {
		modeLine = "Vas a llevar dos sombreros desde el primer dia: la franquicia y la ciudad. No quiero que uses ese poder a medias."
	}

	switch game.InitialScenario {
	case "rebuild":
		return "Te traje a " + city + " para construir con paciencia. " + franchise + " no necesita humo, necesita direccion. " + modeLine + " Si hacemos bien las bases, el resto llega."
	case "contention":
		return "No te contraté para aprender en el cargo. " + franchise + " tiene talento, gasto y expectativas encima desde el primer dia. " + modeLine + " Quiero resultados rapido y no pienso disfrazarlo."
	case "decline":
		return "Acá todavía pesa demasiado el pasado. " + franchise + " vive bajo la sombra de lo que fue, y la ciudad lo siente. " + modeLine + " Tu trabajo es devolver autoridad antes de que esto se vuelva costumbre."
	default:
		return city + " no tiene historia previa que la sostenga; eso es una ventaja y una responsabilidad. " + franchise + " empieza desde cero, y con eso viene libertad real. " + modeLine + " Quiero vision, criterio y una identidad que se note desde el primer movimiento."
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
