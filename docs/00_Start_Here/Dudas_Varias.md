# Dudas Varias

> Espacio vivo para registrar preguntas sueltas, aclaraciones y respuestas que surjan durante el desarrollo de PulseCity.
> La idea es que funcione como referencia rápida dentro del vault de Obsidian, sin mezclar estas notas con los documentos canónicos del proyecto.

---

## 2026-05-08 — ¿Qué comandos se usan para correr y probar la app?

### Pregunta

Viniendo de apps web más típicas, la duda fue qué equivalente tendría PulseCity a comandos como `npm run dev`, `npx prisma migrate dev` o `dotnet run`.

### Respuesta

PulseCity no va a girar alrededor de un único comando de desarrollo, porque el proyecto está pensado como un sistema distribuido con múltiples servicios y stacks distintos.

La separación esperada es:

- **Infraestructura compartida** con `make up`, `make down`, `make dev` y `make logs`
- **Frontend** con un comando propio tipo `npm run dev`
- **Servicios Go** con comandos tipo `go run ./cmd/main.go`
- **Servicios Rust** con comandos tipo `cargo run`
- **Builds y tests globales** con `make build` y `make test`

### Estado actual del repo

Hoy el repositorio ya tiene estos comandos en el `Makefile`:

- `make up` → levanta Docker Compose
- `make down` → baja Docker Compose
- `make dev` → levanta infra y sigue logs
- `make logs` → muestra logs
- `make build` → compila Go y Rust
- `make test` → corre tests de Go y Rust
- `make nats-eventos` → escucha todos los eventos de NATS
- `make nats-tiempo` → escucha `tiempo.*`
- `make nats-jugadores` → escucha `jugador.*`

### Infra local actual

- **TimescaleDB / PostgreSQL** en `localhost:5433`
- **NATS** en `localhost:4222`
- **Monitor de NATS** en `localhost:8222`

### Criterio acordado

Conviene mantener dos niveles de comandos:

- **Nivel repo** para orquestación general con `make`
- **Nivel servicio** para desarrollo puntual con herramientas nativas de cada stack

Más adelante, cuando existan `gateway`, `map-service` y `frontend`, conviene agregar targets concretos como:

- `make run-gateway`
- `make run-map-service`
- `make run-frontend`

---

## Cómo usar este archivo

- Agregar nuevas dudas con fecha
- Mantener respuestas cortas y concretas
- Si una duda se convierte en decisión de arquitectura, moverla luego a un documento canónico
