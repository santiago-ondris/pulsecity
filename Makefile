.PHONY: up down build test logs

# Servicios
up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f

# Build
build-rust:
	cargo build

build-go:
	go build ./...

build: build-rust build-go

# Tests
test-rust:
	cargo test

test-go:
	go test ./...

test: test-rust test-go

# NATS — ver eventos en tiempo real
nats-eventos:
	nats sub ">"

nats-jugadores:
	nats sub "jugador.*"

nats-tiempo:
	nats sub "tiempo.*"

# Dev — levanta todo y muestra logs
dev: up
	docker compose logs -f