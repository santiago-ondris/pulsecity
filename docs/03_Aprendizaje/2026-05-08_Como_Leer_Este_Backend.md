# 2026-05-08 — Cómo leer este backend viniendo de ASP.NET Core o NestJS

## Objetivo de esta nota

Esta nota existe para bajar a tierra cómo pensar PulseCity si venís de frameworks web más opinionated como ASP.NET Core o NestJS.

La idea no es explicar solo el slice actual, sino dejar una forma mental de leer el proyecto que después se pueda extrapolar al resto de la app.

---

## La idea principal

En este proyecto no estamos armando una única app backend grande. Estamos armando varios procesos separados, cada uno con una responsabilidad clara.

En este primer slice:

- `gateway` recibe tráfico HTTP y WebSocket
- `map-service` hace trabajo específico de generación de mapa
- **NATS** conecta ambos servicios mediante eventos

La lógica no viaja por llamadas HTTP entre servicios. Viaja por publicación y consumo de mensajes.

---

## Traducción mental desde ASP.NET Core o NestJS

La equivalencia aproximada sería esta:

- `gateway` se parece a tu API principal
- `map-service` se parece a un worker o microservicio separado
- **NATS** se parece a un message broker liviano
- el **WebSocket** del `gateway` empuja cambios al frontend

Entonces, en vez de pensar:

`controller -> service -> repository`

muchas veces en PulseCity vas a pensar:

`HTTP request -> publish event -> otro servicio procesa -> publish result -> gateway reenvía delta`

---

## El flujo del vertical slice actual

El flujo completo implementado hoy es:

1. El cliente hace `POST /api/v1/games` al `gateway`
2. El `gateway` no genera el mapa por su cuenta
3. Publica en NATS el evento `mapa.generacion_iniciada`
4. `map-service` está suscripto a ese evento
5. Cuando lo recibe, simula cuatro pasos de generación
6. Publica:
   - `mapa.terreno_listo`
   - `mapa.zonas_calculadas`
   - `mapa.estadio_ubicado`
   - `mapa.generacion_completa`
7. El `gateway` escucha esos eventos
8. Cada vez que llega uno, lo reenvía por WebSocket al frontend como un delta

Esto ya sigue el canon de PulseCity:

- HTTP para entrada del jugador
- NATS para comunicación entre servicios
- WebSocket para cambios incrementales

---

## Cómo está pensado el `gateway`

Archivo principal:

- [main.go](/workspace/services/gateway/cmd/main.go)

El `gateway` cumple el rol de puerta de entrada. Si vinieras de ASP.NET Core o NestJS, lo podés pensar como el backend principal que expone endpoints al frontend.

Hoy hace cuatro cosas:

- levanta el servidor HTTP
- se conecta a NATS
- acepta conexiones WebSocket
- escucha eventos `mapa.*` y los retransmite al cliente

No genera el mapa ni resuelve esa lógica localmente. Solo coordina entrada y salida.

### Rutas HTTP

Archivo:

- [http.go](/workspace/services/gateway/internal/handlers/http.go)

Rutas actuales:

- `GET /healthz`
- `GET /ws`
- `POST /api/v1/games`

#### `POST /api/v1/games`

Este endpoint:

- lee opcionalmente `city_name`
- crea un `game_id`
- publica `mapa.generacion_iniciada` en NATS
- responde `202 Accepted`

Eso se parece a un controller action que no ejecuta toda la operación en el mismo proceso, sino que despacha un comando al bus.

### Cliente NATS

Archivo:

- [client.go](/workspace/services/gateway/internal/nats/client.go)

Es un wrapper pequeño para encapsular:

- publicación de JSON
- suscripción a subjects
- cierre ordenado de la conexión

Pensalo como un service de infraestructura inyectable, pero sin framework pesado alrededor.

### Hub WebSocket

Archivo:

- [hub.go](/workspace/services/gateway/internal/ws/hub.go)

Mantiene un conjunto de conexiones WebSocket activas y permite hacer broadcast de mensajes JSON.

Hoy no tiene:

- autenticación
- rooms
- sesiones por partida
- reconexión avanzada

Solo existe para probar el patrón base: evento que llega por NATS, delta que sale por WebSocket.

---

## Cómo está pensado el `map-service`

Archivo principal:

- [main.rs](/workspace/services/map-service/src/main.rs)

El `map-service` es un proceso separado. No expone HTTP en este slice. Su trabajo es escuchar eventos, hacer trabajo y publicar nuevos eventos.

Eso se parece más a:

- un `BackgroundService` en .NET
- un microservicio consumer en NestJS
- un worker desacoplado

### Flujo interno actual

`map-service`:

- se conecta a NATS
- se suscribe a `mapa.generacion_iniciada`
- espera mensajes
- procesa cada mensaje
- publica progreso

La función central es `process_generation(...)`.

Hoy no hay generación procedural real. En vez de eso:

- publica `terrain`
- espera unos milisegundos
- publica `zoning`
- espera
- publica `stadium`
- espera
- publica `complete`

Eso nos permite validar arquitectura y flujo antes de meter Perlin Noise, Voronoi y estado más complejo.

---

## Contrato compartido entre servicios

Archivos:

- [map_generation.go](/workspace/services/gateway/internal/domain/map_generation.go)
- [events.rs](/workspace/services/map-service/src/events.rs)
- [map_generation.md](/workspace/shared/events/map_generation.md)

Aunque un servicio esté en Go y el otro en Rust, ambos hablan el mismo contrato JSON.

Ese contrato define:

- el nombre de los eventos
- el payload del evento inicial
- el payload de cada evento de progreso
- el envelope que el `gateway` envía al frontend por WebSocket

La idea importante acá es que los servicios no comparten memoria ni llamadas directas. Comparten contratos.

Ese patrón se va a repetir mucho en PulseCity.

---

## La estructura de carpetas en Go

En `gateway` usé esta organización:

- `cmd/main.go`
- `internal/handlers`
- `internal/nats`
- `internal/ws`
- `internal/domain`

Una traducción mental aproximada sería:

- `cmd/main.go` → bootstrap / `Program.cs` / `main.ts`
- `handlers` → controllers
- `nats` → infraestructura de mensajería
- `ws` → gateway de sockets
- `domain` → tipos, payloads y conceptos de negocio livianos

No es una copia exacta de NestJS o ASP.NET Core, pero sirve como mapa mental inicial.

---

## Por qué esto puede sentirse alienígena al principio

Hay varias razones:

- Go no usa clases ni decorators
- Rust tampoco usa OOP tradicional
- hay menos magia de framework
- el wiring suele hacerse explícitamente
- muchas decisiones que un framework toma por vos acá están visibles

La contracara es que, una vez que te acostumbrás, el flujo es bastante directo:

- entra un request
- se publica un evento
- otro servicio reacciona
- vuelve un resultado por NATS
- el `gateway` lo empuja al cliente

---

## Cómo conviene leer el código

Para no abrumarte, este orden sirve bastante:

1. [http.go](/workspace/services/gateway/internal/handlers/http.go)
2. [main.go](/workspace/services/gateway/cmd/main.go)
3. [hub.go](/workspace/services/gateway/internal/ws/hub.go)
4. [main.rs](/workspace/services/map-service/src/main.rs)
5. [map_generation.md](/workspace/shared/events/map_generation.md)

Ese recorrido va de lo más familiar a lo menos familiar:

- primero endpoint HTTP
- después wiring de app
- después WebSocket
- después consumer en Rust
- al final contrato compartido

---

## Idea clave para el resto del proyecto

Esta forma de pensar se va a repetir mucho en PulseCity.

Ejemplos futuros:

- `gateway` recibe una acción del jugador
- publica un evento en NATS
- `team-service` actualiza roster
- `agent-service` actualiza emociones
- `city-service` ajusta economía
- `narrative-service` genera texto
- `gateway` escucha resultados y envía deltas al frontend

La implementación cambiará, pero el patrón mental será muy parecido.

---

## Resumen corto

Si venís de ASP.NET Core o NestJS, la mejor forma de leer PulseCity es esta:

- no pensar en una sola app
- pensar en varios servicios chicos
- pensar en eventos como la forma principal de coordinación
- pensar en el `gateway` como puente entre frontend y el resto del sistema
- pensar en WebSocket como canal de deltas, no de estado completo

---

## Próximas notas útiles para esta carpeta

Más adelante conviene agregar notas similares sobre:

- cómo leer servicios en Go dentro de PulseCity
- cómo leer servicios en Rust dentro de PulseCity
- cómo pensar NATS en este proyecto
- cómo pensar state ownership
- cómo pensar el modelo de deltas por WebSocket
