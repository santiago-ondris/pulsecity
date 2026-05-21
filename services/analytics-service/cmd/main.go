package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pulsecity/services/analytics-service/internal/domain"
	natsclient "github.com/pulsecity/services/analytics-service/internal/nats"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	natsURL := envOrDefault("NATS_URL", "nats://localhost:4222")

	bus, err := natsclient.New(natsURL)
	if err != nil {
		return err
	}
	defer bus.Close()

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
