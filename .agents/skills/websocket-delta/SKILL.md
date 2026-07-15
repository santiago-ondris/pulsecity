---
name: websocket-delta
description: >
  Patrones de WebSocket con arquitectura delta-only para PulseCity.
  Úsala cuando trabajes en el gateway (Go), el frontend React, o cualquier
  código que maneje la conexión WebSocket entre backend y browser.
  Cubre: envío de deltas, throttling por velocidad de simulación,
  gestión de sesión activa/inactiva, y reconexión en el cliente.
---

# WebSocket — Delta-Only Architecture

## Principio fundamental

El backend NUNCA envía estado completo al frontend via WebSocket.
Solo envía deltas — lo que cambió desde el último update.
El frontend mantiene su propia copia del estado y aplica cada delta recibido.

Esto permite:
- Mensajes pequeños aunque el estado global sea grande
- Velocidades de simulación altas (x20) sin saturar la conexión
- El frontend puede reconectarse y recibir solo lo que faltó

## Estructura de un delta

```go
// En el gateway (Go)
type Delta struct {
    Type      string          `json:"type"`       // "jugador.firmado", "ciudad.suelo_actualizado", etc.
    Payload   json.RawMessage `json:"payload"`    // datos específicos del cambio
    Timestamp int64           `json:"ts"`         // unix nanos
    SimDay    string          `json:"sim_day"`    // fecha simulada "2026-03-15"
}

// Ejemplos de payloads por tipo
type JugadorFirmadoDelta struct {
    PlayerID string  `json:"player_id"`
    Salario  float64 `json:"salario"`
    Anos     int     `json:"anos"`
}

type SueloActualizadoDelta struct {
    ZonaID        string  `json:"zona_id"`
    ValorAnterior float64 `json:"valor_anterior"`
    ValorNuevo    float64 `json:"valor_nuevo"`
}
```

## Gateway en Go — Hub de conexiones WebSocket

```go
import "github.com/gorilla/websocket"

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 4096,
    CheckOrigin: func(r *http.Request) bool {
        return true // configurar CORS apropiadamente en producción
    },
}

type Client struct {
    conn     *websocket.Conn
    send     chan []byte
    partidaID string
}

type Hub struct {
    clients    map[string]*Client // partidaID → client
    broadcast  chan Delta
    register   chan *Client
    unregister chan *Client
    mu         sync.RWMutex
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client.partidaID] = client
            h.mu.Unlock()
            // Notificar al simulation-loop que hay sesión activa
            natsClient.Publish("tiempo.sesion_iniciada", ...)

        case client := <-h.unregister:
            h.mu.Lock()
            delete(h.clients, client.partidaID)
            h.mu.Unlock()
            close(client.send)
            // Notificar al simulation-loop que no hay sesión
            natsClient.Publish("tiempo.sesion_terminada", ...)

        case delta := <-h.broadcast:
            data, _ := json.Marshal(delta)
            h.mu.RLock()
            for _, client := range h.clients {
                select {
                case client.send <- data:
                default:
                    // buffer lleno: cliente lento, desconectar
                    close(client.send)
                    delete(h.clients, client.partidaID)
                }
            }
            h.mu.RUnlock()
        }
    }
}

// Writer por cliente — un goroutine por conexión
func (c *Client) writePump() {
    ticker := time.NewTicker(30 * time.Second) // ping keepalive
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()

    for {
        select {
        case msg, ok := <-c.send:
            c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if !ok {
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
                return
            }
        case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}
```

## Throttling por velocidad de simulación

El gateway aplica throttling inteligente: a mayor velocidad de simulación,
más updates por segundo. El throttle limita la frecuencia máxima de envío
de deltas para no saturar el WebSocket con cientos de eventos por segundo en x20.

```go
type ThrottledHub struct {
    *Hub
    velocidad    float64       // 1.0, 5.0, 20.0
    lastSent     time.Time
    pendingDeltas []Delta
    mu           sync.Mutex
}

func (th *ThrottledHub) MinInterval() time.Duration {
    // x1: max 1 update/segundo → 1000ms
    // x5: max 5 updates/segundo → 200ms
    // x20: max 20 updates/segundo → 50ms
    return time.Duration(1000/th.velocidad) * time.Millisecond
}

func (th *ThrottledHub) EnqueueDelta(d Delta) {
    th.mu.Lock()
    th.pendingDeltas = append(th.pendingDeltas, d)
    th.mu.Unlock()
}

func (th *ThrottledHub) FlushLoop() {
    ticker := time.NewTicker(50 * time.Millisecond) // check cada 50ms
    for range ticker.C {
        if time.Since(th.lastSent) < th.MinInterval() {
            continue
        }
        th.mu.Lock()
        if len(th.pendingDeltas) == 0 {
            th.mu.Unlock()
            continue
        }
        batch := th.pendingDeltas
        th.pendingDeltas = nil
        th.lastSent = time.Now()
        th.mu.Unlock()

        // Enviar batch como array de deltas en un solo mensaje WebSocket
        th.Hub.broadcast <- BatchDelta{Deltas: batch}
    }
}

// Actualizar velocidad cuando cambia (evento NATS tiempo.velocidad_cambiada)
func (th *ThrottledHub) SetVelocidad(v float64) {
    th.mu.Lock()
    th.velocidad = v
    th.mu.Unlock()
}
```

## Frontend React — aplicar deltas al estado local

```typescript
// Estado local del frontend — copia del estado del backend
interface GameState {
  players: Record<string, Player>
  cityZones: Record<string, Zone>
  agents: Record<string, AgentState>
  simDay: string
}

// Hook de WebSocket con reconexión automática
function useGameWebSocket(partidaId: string) {
  const [state, setState] = useState<GameState>(initialState)
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectTimeout = useRef<ReturnType<typeof setTimeout>>()

  const connect = useCallback(() => {
    const ws = new WebSocket(`wss://api.pulsecity.dev/ws/${partidaId}`)

    ws.onmessage = (event) => {
      const deltas: Delta[] = JSON.parse(event.data)
      setState(prev => applyDeltas(prev, deltas))
    }

    ws.onclose = () => {
      // Reconexión exponencial con jitter
      const delay = Math.min(1000 * 2 ** retryCount.current, 30000)
      reconnectTimeout.current = setTimeout(connect, delay + Math.random() * 1000)
      retryCount.current++
    }

    ws.onopen = () => {
      retryCount.current = 0
    }

    wsRef.current = ws
  }, [partidaId])

  useEffect(() => {
    connect()
    return () => {
      wsRef.current?.close()
      clearTimeout(reconnectTimeout.current)
    }
  }, [connect])

  return state
}

// Aplicar deltas al estado local — función pura
function applyDeltas(state: GameState, deltas: Delta[]): GameState {
  let next = { ...state }

  for (const delta of deltas) {
    switch (delta.type) {
      case 'jugador.firmado': {
        const p = delta.payload as JugadorFirmadoDelta
        next.players = {
          ...next.players,
          [p.player_id]: { ...next.players[p.player_id], salario: p.salario }
        }
        break
      }
      case 'ciudad.suelo_actualizado': {
        const z = delta.payload as SueloActualizadoDelta
        next.cityZones = {
          ...next.cityZones,
          [z.zona_id]: { ...next.cityZones[z.zona_id], valor: z.valor_nuevo }
        }
        break
      }
      // ... resto de tipos de delta
    }
    next.simDay = delta.sim_day
  }

  return next
}
```

## Sesión activa / inactiva — ciclo de vida

```
WebSocket abre  → gateway publica "tiempo.sesion_iniciada"  → simulation-loop activa el loop
WebSocket cierra → gateway publica "tiempo.sesion_terminada" → simulation-loop duerme
```

El simulation-loop SOLO corre cuando hay al menos una conexión WebSocket activa.
Cuando el jugador vuelve después de desconectarse, el tiempo retoma exactamente donde quedó.

## Errores frecuentes a evitar

```go
// ❌ MAL — enviar estado completo
ws.WriteJSON(fullGameState) // puede ser megabytes

// ✅ BIEN — solo el delta
ws.WriteJSON(Delta{Type: "ciudad.suelo_actualizado", Payload: ...})

// ❌ MAL — un goroutine de escritura compartido
go func() {
    for delta := range hub.broadcast {
        conn.WriteMessage(...) // data race si múltiples senders
    }
}()

// ✅ BIEN — un writePump por cliente, channel buffered
// (ver writePump arriba)

// ❌ MAL — bloquear en el handler de NATS para enviar WebSocket
natsConn.Subscribe("partido.terminado", func(msg *nats.Msg) {
    hub.broadcast <- delta // puede bloquear si el channel está lleno
})

// ✅ BIEN — channel con buffer o select con default
select {
case hub.broadcast <- delta:
default:
    log.Warn("broadcast channel full, dropping delta") // nunca bloquear NATS
}
```
