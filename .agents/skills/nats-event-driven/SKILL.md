---
name: nats-event-driven
description: >
  Patrones de arquitectura event-driven con NATS en Go y Rust.
  Úsala cuando estés implementando publishers, subscribers, o cualquier
  comunicación entre servicios via event bus. Cubre: conexión a NATS,
  publish/subscribe, JetStream, idempotencia, manejo de errores,
  y patrones de choreography sin orquestador central.
---

# NATS — Event-Driven Architecture para Go y Rust

## Principios de diseño para PulseCity

- Los servicios NUNCA se llaman directamente entre sí para escribir estado.
- Toda comunicación de escritura va via eventos NATS.
- Cada servicio es dueño exclusivo de su estado. Reacciona a eventos de otros, no los instruye.
- Los eventos son hechos pasados e inmutables: `jugador.firmado`, no `firmar_jugador`.
- Todos los consumers son idempotentes: procesar el mismo evento dos veces no rompe el estado.

## NATS en Go

### Conexión y setup básico

```go
import "github.com/nats-io/nats.go"

// Conectar con reconexión automática
nc, err := nats.Connect(
    nats.DefaultURL,
    nats.RetryOnFailedConnect(true),
    nats.MaxReconnects(-1),           // reconexión infinita
    nats.ReconnectWait(2*time.Second),
    nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
        log.Printf("NATS disconnected: %v", err)
    }),
    nats.ReconnectHandler(func(nc *nats.Conn) {
        log.Printf("NATS reconnected to %s", nc.ConnectedUrl())
    }),
)
if err != nil {
    return fmt.Errorf("nats connect: %w", err)
}
defer nc.Drain() // drain antes de cerrar, no Close()
```

### Publicar un evento

```go
type JugadorFirmadoEvent struct {
    PlayerID     string  `json:"player_id"`
    FranquiciaID string  `json:"franquicia_id"`
    Salario      float64 `json:"salario"`
    Anos         int     `json:"anos"`
    Timestamp    int64   `json:"timestamp"` // unix nanos — siempre incluir
}

func PublishJugadorFirmado(nc *nats.Conn, event JugadorFirmadoEvent) error {
    event.Timestamp = time.Now().UnixNano()
    data, err := json.Marshal(event)
    if err != nil {
        return fmt.Errorf("marshal event: %w", err)
    }
    return nc.Publish("jugador.firmado", data)
}
```

### Suscribirse a un evento

```go
// Subscribe siempre en goroutine propia, nunca bloquear el handler
sub, err := nc.Subscribe("jugador.firmado", func(msg *nats.Msg) {
    var event JugadorFirmadoEvent
    if err := json.Unmarshal(msg.Data, &event); err != nil {
        log.Printf("ERROR unmarshal jugador.firmado: %v", err)
        return // nunca re-publicar el error como otro evento NATS desde acá
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := handleJugadorFirmado(ctx, event); err != nil {
        log.Printf("ERROR handling jugador.firmado player=%s: %v", event.PlayerID, err)
        // log y seguir — no panics, no retries automáticos sin JetStream
    }
})
if err != nil {
    return fmt.Errorf("subscribe jugador.firmado: %w", err)
}
defer sub.Unsubscribe()
```

### Queue groups (load balancing entre réplicas)

```go
// Si hay múltiples réplicas del servicio, usar QueueSubscribe
// NATS entrega cada mensaje a solo UNA réplica del grupo
sub, err := nc.QueueSubscribe("jugador.firmado", "city-service-workers", handler)
```

### Wildcard subscriptions

```go
// Escuchar todos los eventos de jugadores
sub, err := nc.Subscribe("jugador.*", handler)

// Escuchar toda la jerarquía de eventos
sub, err := nc.Subscribe("ciudad.>", handler)
```

## NATS en Rust

### Dependency (Cargo.toml)

```toml
[dependencies]
async-nats = "0.35"
tokio = { version = "1", features = ["full"] }
serde = { version = "1", features = ["derive"] }
serde_json = "1"
```

### Conexión

```rust
use async_nats::Client;

async fn connect_nats(url: &str) -> Result<Client, async_nats::ConnectError> {
    async_nats::connect(url).await
}
```

### Publicar desde Rust

```rust
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize)]
struct TiempoDiaAvanzadoEvent {
    fecha_simulada: String,
    velocidad: f32,
    timestamp: i64,
}

async fn publish_dia_avanzado(client: &Client, event: TiempoDiaAvanzadoEvent) -> Result<()> {
    let data = serde_json::to_vec(&event)?;
    client.publish("tiempo.dia_avanzado", data.into()).await?;
    Ok(())
}
```

### Suscribirse desde Rust

```rust
async fn subscribe_partido_terminado(client: &Client) -> Result<()> {
    let mut sub = client.subscribe("partido.terminado").await?;

    while let Some(msg) = sub.next().await {
        let event: PartidoTerminadoEvent = match serde_json::from_slice(&msg.payload) {
            Ok(e) => e,
            Err(e) => {
                tracing::error!("Failed to parse partido.terminado: {}", e);
                continue; // no panics
            }
        };

        if let Err(e) = handle_partido_terminado(&event).await {
            tracing::error!("Error handling partido.terminado: {}", e);
        }
    }
    Ok(())
}
```

## Idempotencia — obligatoria en todos los handlers

```go
// Patrón: verificar si ya fue procesado antes de actuar
func handleJugadorFirmado(ctx context.Context, event JugadorFirmadoEvent) error {
    // Verificar si ya procesamos este evento (por event_id o por estado resultante)
    already, err := db.AlreadyProcessed(ctx, "jugador.firmado", event.PlayerID, event.Timestamp)
    if err != nil {
        return err
    }
    if already {
        return nil // ya procesado, no hacer nada
    }

    // Procesar...
    // Marcar como procesado (idealmente en la misma tx que la actualización de estado)
}
```

## Patrones a evitar

```go
// ❌ MAL — servicios llamándose directamente para escribir
func (s *CityService) OnJugadorFirmado(playerID string) {
    s.agentService.UpdateEmotionalState(playerID) // NUNCA
}

// ✅ BIEN — cada servicio reacciona independientemente al evento
func (s *CityService) Subscribe() {
    nc.Subscribe("jugador.firmado", func(msg *nats.Msg) {
        // city-service actualiza su propio estado
        s.updateLandValue(...)
    })
}

// ❌ MAL — publicar eventos de error como eventos de dominio
nc.Publish("jugador.firma_fallida", ...) // genera cascadas de error difíciles de rastrear

// ✅ BIEN — loggear el error, el estado no cambia, el evento original puede re-procesarse

// ❌ MAL — handlers que bloquean con operaciones lentas
nc.Subscribe("partido.terminado", func(msg *nats.Msg) {
    time.Sleep(2 * time.Second) // bloquea el dispatcher de NATS
    generateNarrative()
})

// ✅ BIEN — delegar a goroutine/task
nc.Subscribe("partido.terminado", func(msg *nats.Msg) {
    go func() {
        // procesar en background
        // el delay de 250-500ms de narrative-service va acá
    }()
})
```

## El delay intencional de narrative-service

narrative-service espera 250-500ms después de recibir un evento antes de consultar
el estado de los agentes al LLM. Esto garantiza que agent-service ya procesó el mismo
evento. No es coordinación, es timing. Implementación:

```go
nc.Subscribe("jugador.firmado", func(msg *nats.Msg) {
    go func() {
        time.Sleep(350 * time.Millisecond) // delay intencional
        // ahora el estado de agent-service ya está actualizado
        ctx := context.Background()
        agentState := agentClient.GetState(ctx, event.PlayerID)
        narrative := llm.Generate(ctx, event, agentState)
        nc.Publish("narrativa.evento_generado", narrative)
    }()
})
```

## Subject naming conventions para PulseCity

```
<dominio>.<accion_pasada>

tiempo.dia_avanzado
jugador.firmado
jugador.traspasado
jugador.lesionado
partido.terminado
agente.estado_cambio
ciudad.suelo_actualizado
narrativa.evento_generado
mapa.generacion_completa
```

Nunca usar verbos en presente o imperativo: `firmar_jugador` ❌, `jugador.firmado` ✅
