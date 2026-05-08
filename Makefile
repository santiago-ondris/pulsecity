.PHONY: up down build test logs run-gateway run-map-service

# Servicios
up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f

# Build
build-rust:
	cargo build --manifest-path services/map-service/Cargo.toml

build-go:
	go -C services/gateway build ./...

build: build-rust build-go

# Tests
test-rust:
	cargo test --manifest-path services/map-service/Cargo.toml

test-go:
	go -C services/gateway test ./...

test: test-rust test-go

# Run
run-gateway:
	go -C services/gateway run ./cmd/main.go

run-map-service:
	cargo run --manifest-path services/map-service/Cargo.toml

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
