# Sesión — 2026-05-08

## Estado de la sesión

Primera sesión real de codeo de PulseCity.

Se trabajó sobre el corazón técnico de `Milestone 1`, empezando desde un repo casi vacío y dejando una base ya jugable/tocable a nivel técnico:

- `gateway` en Go
- `map-service` en Rust
- flujo por NATS
- WebSocket con `snapshot` y `patch`
- rehidratación
- persistencia básica en PostgreSQL
- página debug en el `gateway`
- frontend real mínimo en React + TypeScript + Vite

---

## Qué se hizo

### 1. Contexto y documentación

Se leyó y fijó contexto desde:

- `docs/00_Start_Here/PulseCity_v2.md`
- `docs/01_Canon/PulseCity_Arquitectura_Tecnica.md`
- `docs/01_Canon/PulseCity_Experiencia_Jugador.md`

También se creó:

- `docs/00_Start_Here/Dudas_Varias.md`
- varias notas en `docs/03_Aprendizaje/`
- progreso técnico en `docs/02_Progress/2026-05-08_Vertical_Slice_Gateway_Map_Service.md`

### 2. Vertical slice inicial

Se implementó el flujo:

1. frontend o debug UI hace `POST /api/v1/games`
2. `gateway` publica `mapa.generacion_iniciada`
3. `map-service` genera el mapa
4. publica progreso por etapas
5. `gateway` traduce eso a `map.snapshot` / `map.patch`
6. cliente consume el estado y lo renderiza

### 3. Contrato de frontend

Se ordenó el contrato WebSocket:

- `map.snapshot`
- `map.patch`

Se separó correctamente:

- snapshot inicial
- cambios incrementales
- rehidratación tardía

### 4. Rehidratación

Se resolvió el caso de cliente que “entra tarde”:

- cache en memoria por `game_id`
- `GET /api/v1/games/{gameID}/snapshot`
- `GET /ws?game_id=...`

### 5. Persistencia

Se agregó persistencia básica en PostgreSQL para:

- creación de partida
- último snapshot del mapa

Quedó:

- `services/gateway/internal/persistence/postgres.go`
- `db/migrations/001_create_games.sql`

### 6. Generación procedural del mapa

La generación evolucionó desde una grilla rígida a:

- `Perlin 2D` real para elevación
- partición Voronoi simple y determinística para distritos
- reglas de zonas sobre esa base
- ubicación automática del estadio

### 7. Visualización

Se hicieron dos capas:

- mejora de la debug page del `gateway`
- frontend real mínimo en React

El frontend nuevo:

- vive en `frontend/`
- usa `Vite + React + TypeScript`
- crea partidas reales
- consume snapshot/patch
- muestra el mapa 2D mínimo

### 8. Fixes finales

Se corrigió un problema real del frontend:

- faltaba soporte CORS en el `gateway`
- el botón `Nueva partida` podía quedar clavado en “Creando partida...”

Eso quedó resuelto y validado.

---

## Milestone 1 — estimación actual

### Porcentaje completado

Estimación al cierre de esta sesión:

**55% de Milestone 1**

### Por qué 55%

Ya existe una base bastante seria de la parte técnica:

- arquitectura inicial viva
- evento de creación de partida
- generación de mapa funcional
- contrato frontend/backend estable
- persistencia inicial
- cliente real mínimo

Todavía falta bastante de la experiencia real del milestone:

- render visual más cercano a la ceremonia final
- frontend más pulido
- integración de narrativa
- primer evento del Owner
- identidad de franquicia y creación de partida más completa
- mapa visualmente más expresivo

---

## Qué falta para Milestone 1

### Backend / dominio

- enriquecer más la generación procedural del mapa
- decidir si se sigue profundizando Perlin/Voronoi o si ya alcanza para pasar al render serio
- formalizar mejor migraciones y bootstrap de DB
- modelar mejor la entidad `partida`

### Frontend

- evolucionar el frontend mínimo hacia una experiencia real
- reemplazar el render 2D básico por una escena más rica
- eventualmente migrar esa representación a Three.js

### Experiencia de juego

- flujo más completo de creación de franquicia
- ceremonia visual del nacimiento del mundo
- primer evento narrativo obligatorio del Owner

---

## Punto exacto para retomar la próxima sesión

El sistema hoy puede:

- crear partida
- generar mapa
- emitir eventos por NATS
- traducirlos a snapshot/patch
- persistir snapshot
- rehidratar clientes tardíos
- mostrar el mapa en frontend React mínimo

### Lo más lógico para la próxima sesión

Opción recomendada:

- seguir con el frontend real
- empezar a reemplazar la visual 2D mínima por una representación más cercana a la definitiva

El siguiente salto natural sería:

- mantener React para estado y UI
- introducir una primera escena más expresiva
- decidir si ese paso ya entra con Three.js o si conviene una transición intermedia

---

## Consideraciones importantes para no olvidar

### Arquitectura

- no romper `snapshot` / `patch`
- no mandar estado completo arbitrariamente por WebSocket
- mantener `gateway` como traductor hacia frontend
- mantener `map-service` como dueño de la generación del mundo

### Persistencia

- hoy la rehidratación usa memoria + PostgreSQL
- no hay todavía historial de eventos ni replay
- el snapshot persistido es el último estado conocido

### Frontend

- `frontend/` ya compila con `npm run build`
- `make run-frontend` ya existe
- el `gateway` ya expone CORS para `http://localhost:5173`

### VSCode / TypeScript

Si vuelven falsos errores del editor:

- usar TypeScript del workspace

- reiniciar TS server
- revisar la nota:
  - `docs/03_Aprendizaje/2026-05-08_VSCode_y_TypeScript_React.md`

---

## Comandos útiles para retomar

### Infra

```bash
make up
```

### Servicios

```bash
make run-map-service
make run-gateway
make run-frontend
```

### Frontend

Abrir:

```text
http://localhost:5173/
```

### Verificación

```bash
make test
cd frontend && npm run build
```

---

## Resultado de la sesión

La sesión terminó con un sistema mucho más avanzado de lo que había al empezar:

- backend inicial real
- algoritmo procedural ya serio
- persistencia básica
- cliente real

La base del proyecto dejó de ser solo documentación y pasó a ser software ejecutable de punta a punta.
