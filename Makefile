.PHONY: up down build test logs run-agent-service run-analytics-service run-city-service run-gateway run-map-service run-match-service run-narrative-service run-team-service run-frontend dev-app

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
	cargo build --manifest-path services/agent-service/Cargo.toml
	cargo build --manifest-path services/match-service/Cargo.toml

build-go:
	GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway build ./...
	GOCACHE=/tmp/pulsecity-narrative-gocache go -C services/narrative-service build ./...
	GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service build ./...
	GOCACHE=/tmp/pulsecity-city-gocache go -C services/city-service build ./...
	GOCACHE=/tmp/pulsecity-analytics-gocache go -C services/analytics-service build ./...

build: build-rust build-go

# Tests
test-rust:
	cargo test --manifest-path services/map-service/Cargo.toml
	cargo test --manifest-path services/agent-service/Cargo.toml
	cargo test --manifest-path services/match-service/Cargo.toml

test-go:
	GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...
	GOCACHE=/tmp/pulsecity-narrative-gocache go -C services/narrative-service test ./...
	GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...
	GOCACHE=/tmp/pulsecity-city-gocache go -C services/city-service test ./...
	GOCACHE=/tmp/pulsecity-analytics-gocache go -C services/analytics-service test ./...

test: test-rust test-go

# Run
run-gateway:
	GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway run ./cmd/main.go

run-agent-service:
	cargo run --manifest-path services/agent-service/Cargo.toml

run-map-service:
	cargo run --manifest-path services/map-service/Cargo.toml

run-match-service:
	cargo run --manifest-path services/match-service/Cargo.toml

run-narrative-service:
	GOCACHE=/tmp/pulsecity-narrative-gocache go -C services/narrative-service run ./cmd/main.go

run-team-service:
	GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service run ./cmd/main.go

run-city-service:
	GOCACHE=/tmp/pulsecity-city-gocache go -C services/city-service run ./cmd/main.go

run-analytics-service:
	GOCACHE=/tmp/pulsecity-analytics-gocache go -C services/analytics-service run ./cmd/main.go

run-frontend:
	npm run dev --prefix frontend

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

# Dev app — levanta infra + servicios de app para probar desde el browser.
dev-app: up
	@echo "PulseCity dev app"
	@echo "Gateway:  http://localhost:8080"
	@echo "Frontend: http://localhost:5173"
	@echo "Postgres: localhost:5433"
	@echo "NATS:     localhost:4222"
	@trap 'for pid in $$(jobs -p); do kill $$pid; done' INT TERM EXIT; \
	$(MAKE) run-gateway & \
	$(MAKE) run-map-service & \
	$(MAKE) run-agent-service & \
	$(MAKE) run-team-service & \
	$(MAKE) run-match-service & \
	$(MAKE) run-city-service & \
	$(MAKE) run-narrative-service & \
	$(MAKE) run-analytics-service & \
	$(MAKE) run-frontend & \
	wait
