# 2026-05-08 — Vertical Slice Inicial `gateway` + `map-service`

## Objetivo de la sesión

Dejar un primer corte ejecutable del `Milestone 1` que pruebe el flujo base:

1. El `gateway` recibe una acción HTTP para crear partida
2. Publica `mapa.generacion_iniciada` en NATS
3. El `map-service` consume ese evento
4. Publica progreso simulado de generación de mapa
5. El `gateway` escucha esos eventos y los reenvía por WebSocket como deltas

El slice nacio para fijar estructura, contratos y flujo entre servicios. Luego se fue extendiendo con un mapa procedural-lite, snapshots, rehidratacion y persistencia basica.

---

## Qué quedó implementado

### `services/gateway` (Go)

- `POST /api/v1/games`
  - genera un `game_id`
  - publica `mapa.generacion_iniciada` en NATS
  - responde `202 Accepted`
- `GET /ws`
  - abre WebSocket
  - traduce eventos `mapa.*` a mensajes frontend-oriented
  - emite `map.snapshot` para estado inicial
  - emite `map.patch` para cambios incrementales
- `GET /`
  - sirve una pagina de debug minima
  - permite crear partida desde browser
  - muestra deltas WebSocket en vivo
  - renderiza una grilla simple de mapa con colores por terreno/zona
  - comunica mejor costa, relieve y composicion general del mapa
- `GET /healthz`
  - healthcheck simple
- `GET /api/v1/games/{gameID}/snapshot`
  - devuelve el ultimo snapshot en memoria para una partida
- `GET /ws?game_id=...`
  - si existe snapshot en memoria, lo rehidrata apenas conecta
  - si no existe en memoria, intenta cargarlo desde PostgreSQL

### `frontend` (React + TypeScript + Vite)

- Scaffold manual de `Vite + React + TypeScript`
- Pantalla mínima de creación de partida
- Conexión real a `POST /api/v1/games`
- Reconexión WebSocket por `game_id`
- Carga de snapshot por HTTP
- Estado local aplicado con `map.snapshot` y `map.patch`
- Render 2D mínimo del mapa como paso previo a Three.js

### `services/map-service` (Rust)

- Suscripción a `mapa.generacion_iniciada`
- Publicación secuencial de:
  - `mapa.terreno_listo`
  - `mapa.zonas_calculadas`
  - `mapa.estadio_ubicado`
  - `mapa.generacion_completa`
- Generación procedural-lite determinística de una grilla `28x28`
- Seed derivado del `city_name`
- Terreno basado en `Perlin 2D`, costa y edge falloff
- Terrenos base: `water`, `plain`, `forest`, `hill`
- Zonas base asignadas por humedad + particion Voronoi simple con seeds deterministicas
- Ubicación automática del estadio según score de centralidad y zona
- Delays cortos entre etapas para simular trabajo real

### Contrato compartido

Se documentó en:

- `shared/events/map_generation.md`

Incluye:

- subjects del flujo
- payload de inicio
- payload de progreso
- contrato WebSocket con `snapshot` y `patch`

### Persistencia

- `gateway` ahora se conecta a PostgreSQL
- persiste la partida al crear `game_id`
- persiste el ultimo snapshot completo de mapa por partida
- usa cache en memoria + fallback a base para rehidratacion

Se agrego:

- migracion [001_create_games.sql](/workspace/db/migrations/001_create_games.sql)

---

## Comandos útiles

### Infra

```bash
make up
```

Levanta:

- TimescaleDB en `localhost:5433`
- NATS en `localhost:4222`

### Correr servicios

```bash
make run-map-service
make run-gateway
make run-frontend
```

### Probar desde el navegador

Con ambos servicios arriba, abrir:

```text
http://localhost:8080/
```

La pagina de debug:

- abre WebSocket automaticamente
- permite disparar `POST /api/v1/games` con un boton
- muestra los eventos `mapa.*` recibidos en tiempo real
- dibuja la grilla del mapa y marca el estadio
- aplica `map.snapshot` y `map.patch` como lo haria un frontend real
- permite reconectar el socket a un `game_id` puntual
- permite pedir snapshot por HTTP para simular entrada tardia
- muestra mejor agua, bosque, llanos y costas
- resume porcentajes aproximados de composicion del terreno

### Probar frontend real minimo

Con `gateway`, `map-service` y `frontend` arriba, abrir:

```text
http://localhost:5173/
```

El frontend actual:

- crea partidas reales contra el `gateway`
- escucha `snapshot` y `patch`
- mantiene estado local del mapa
- renderiza una version 2D minima del mundo
- sirve como base de migracion futura a Three.js

### Verificar compilación

```bash
make build
make test
```

---

## Flujo manual esperado

1. Abrir una conexión WebSocket a `ws://localhost:8080/ws`
2. Hacer `POST` a `http://localhost:8080/api/v1/games`
3. Recibir deltas en este orden:
   - `mapa.terreno_listo`
   - `mapa.zonas_calculadas`
   - `mapa.estadio_ubicado`
   - `mapa.generacion_completa`

Ejemplo de request:

```bash
curl -X POST http://localhost:8080/api/v1/games \
  -H "Content-Type: application/json" \
  -d '{"city_name":"Nueva Aurora"}'
```

Alternativamente, ahora se puede probar el mismo flujo desde la pagina de debug en `/`.

### Probar rehidratacion

1. Crear una partida desde la pagina debug
2. Esperar a que termine la generacion
3. Copiar el `game_id`
4. Abrir otra pestaña con:

```text
http://localhost:8080/?game_id=<uuid>
```

Resultado esperado:

- el WebSocket se conecta con `game_id`
- el `gateway` envia un `map.snapshot` rehidratado
- la pagina vuelve a dibujar el mapa aunque haya llegado tarde

Alternativa:

- usar el boton `Cargar snapshot` para pedir `GET /api/v1/games/{gameID}/snapshot`

### Probar persistencia basica

Con la base levantada por `make up`, el `gateway` usa por defecto:

```text
postgres://pulsecity:pulsecity@localhost:5433/pulsecity_dev?sslmode=disable
```

Flujo esperado:

1. Crear una partida
2. Esperar a que llegue al menos el primer snapshot
3. Refrescar la pagina o abrir otra pestaña con el mismo `game_id`
4. El `gateway` debe poder rehidratar desde memoria o desde PostgreSQL

---

## Qué falta después de este slice

- Enriquecer aun mas la generacion procedural sobre la base Perlin/Voronoi actual
- Modelar mejor la entidad `partida`
- Persistir estado de creación de partida
- Conectar un frontend que consuma el WebSocket
- Refinar el contrato de deltas para mapas más ricos
- Introducir resincronizacion para clientes que se conecten tarde
- Persistir snapshots fuera del proceso `gateway`
- Separar formalmente bootstrap de schema y ejecucion de migraciones
- Empezar a separar comandos/eventos compartidos si varios servicios los reutilizan

---

## Validación realizada

- `go test ./...` en `services/gateway` ✅
- `cargo test --manifest-path services/map-service/Cargo.toml` ✅
- `npm run build` en `frontend` ✅
