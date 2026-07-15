# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

> Read `AGENTS.md` before any session — it contains every architectural decision already made and the current milestone state. Do not reopen closed decisions.

---

## Commands

```bash
# Infrastructure (TimescaleDB + NATS)
make up / make down / make logs

# Build
make build-rust      # map-service, agent-service, match-service
make build-go        # gateway, narrative, team, city, analytics
make build           # both

# Test
make test-rust
make test-go
make test

# Run individual services
make run-gateway
make run-agent-service
make run-map-service
make run-match-service
make run-narrative-service
make run-team-service
make run-city-service
make run-analytics-service
make run-frontend     # npm run dev --prefix frontend

# Monitor NATS events in real time
make nats-eventos     # all events
make nats-jugadores   # jugador.*
make nats-tiempo      # tiempo.*
```

Single Rust test: `cargo test --manifest-path services/<name>/Cargo.toml <test_name>`  
Single Go test: `GOCACHE=/tmp/pulsecity-<name>-gocache go -C services/<name> test ./... -run TestName`

---

## Architecture

PulseCity is a browser-based city simulation where a basketball franchise's performance drives the city's economy. The stack splits CPU-intensive work to Rust and I/O-bound work to Go, communicating exclusively via NATS.

### Services

| Service | Lang | Owns |
|---|---|---|
| `gateway` | Go | WebSocket, REST, auth — only service the frontend talks to |
| `agent-service` | Rust | Simulation loop (100ms ticks), ~50 agents with emotional state & relationships |
| `map-service` | Rust | Procedural map generation (Perlin + Voronoi), runs once per game |
| `match-service` | Rust | Stateless match simulation — receives full payload, publishes `partido.terminado` |
| `team-service` | Go | Contracts, roster, salary cap, player stats |
| `city-service` | Go | Urban economy, zoning, land value, municipal budget |
| `narrative-service` | Go | LLM-generated event narratives (Claude API / GPT-4o mini) |
| `analytics-service` | Go | Time series ingestion/queries via TimescaleDB |

### Internal layout

Go services: `cmd/main.go` → `internal/handlers/` + `internal/domain/` + `internal/nats/`  
Rust services: `src/main.rs` → `src/lib.rs` → domain modules

### Principles (non-negotiable)

1. **State ownership is strict.** Each piece of data has exactly one owning service. Others read via events or gateway-assembled queries, never write.
2. **NATS-only inter-service communication.** No HTTP between services. Exception: `gateway` may query services directly to assemble frontend responses.
3. **WebSocket sends deltas only.** Frontend maintains local state; backend sends only what changed. Max 1 update/second at x1 speed.
4. **`match-service` is stateless.** It receives everything it needs in the input payload.
5. **`narrative-service` waits 250–500ms** after receiving a NATS event before querying agent state — this ensures `agent-service` has already processed the same event.
6. **NATS, not Kafka.** PulseCity doesn't need event replay or retention. Decision is final.

### Event naming

```
entidad.accion         # lowercase, underscores within each part
tiempo.dia_avanzado
jugador.firmado
partido.terminado
ciudad.suelo_actualizado
```

### Simulation loop (agent-service)

Runs every 100ms. Pauses when no WebSocket session is active (`tiempo.sesion_iniciada` / `tiempo.sesion_terminada` from gateway). Three speeds: x1 (~1.6s/day), x5, x20.

---

## Development workflow

- Work in **mini-milestones** — agree on a small, finishable goal per session.
- Before starting: read `docs/Sesiones/MILESTONE3/INICIOM3.MD` for current M3 state and where to continue.
- After finishing: mark the mini-milestone as done in `INICIOM3.MD` and update the relevant doc in `/docs/`.
- For any visual work: consult `docs/00_start_here/pulsecity_designsystem.md` before writing CSS or components.
- For game lore and world questions: read `docs/01_Canon/`.

**Current milestone:** M3

---

## Database

TimescaleDB (PostgreSQL) on port 5433. User/db: `pulsecity` / `pulsecity_dev`. Migrations live in `db/migrations/`. Schema is standard SQL — TimescaleDB adds time-series hypertables on top.

Key tables: `games`, `agent_simulation_state`, `team_franchises`, `guest_sessions`, `users`.

# CLAUDE.md

Behavioral guidelines to reduce common LLM coding mistakes. Merge with project-specific instructions as needed.

**Tradeoff:** These guidelines bias toward caution over speed. For trivial tasks, use judgment.

## 1. Think Before Coding

**Don't assume. Don't hide confusion. Surface tradeoffs.**

Before implementing:
- State your assumptions explicitly. If uncertain, ask.
- If multiple interpretations exist, present them - don't pick silently.
- If a simpler approach exists, say so. Push back when warranted.
- If something is unclear, stop. Name what's confusing. Ask.

## 2. Simplicity First

**Minimum code that solves the problem. Nothing speculative.**

- No features beyond what was asked.
- No abstractions for single-use code.
- No "flexibility" or "configurability" that wasn't requested.
- No error handling for impossible scenarios.
- If you write 200 lines and it could be 50, rewrite it.

Ask yourself: "Would a senior engineer say this is overcomplicated?" If yes, simplify.

## 3. Surgical Changes

**Touch only what you must. Clean up only your own mess.**

When editing existing code:
- Don't "improve" adjacent code, comments, or formatting.
- Don't refactor things that aren't broken.
- Match existing style, even if you'd do it differently.
- If you notice unrelated dead code, mention it - don't delete it.

When your changes create orphans:
- Remove imports/variables/functions that YOUR changes made unused.
- Don't remove pre-existing dead code unless asked.

The test: Every changed line should trace directly to the user's request.

## 4. Goal-Driven Execution

**Define success criteria. Loop until verified.**

Transform tasks into verifiable goals:
- "Add validation" → "Write tests for invalid inputs, then make them pass"
- "Fix the bug" → "Write a test that reproduces it, then make it pass"
- "Refactor X" → "Ensure tests pass before and after"

For multi-step tasks, state a brief plan:
```
1. [Step] → verify: [check]
2. [Step] → verify: [check]
3. [Step] → verify: [check]
```

Strong success criteria let you loop independently. Weak criteria ("make it work") require constant clarification.

---

**These guidelines are working if:** fewer unnecessary changes in diffs, fewer rewrites due to overcomplication, and clarifying questions come before implementation rather than after mistakes.
