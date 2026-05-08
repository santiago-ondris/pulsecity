# AGENTS.md — PulseCity

> Este archivo es el punto de entrada para cualquier agente de IA (Cursor, Claude Code, Copilot, etc.) que trabaje en este repositorio.
> Leé esto antes de tocar cualquier archivo. Acá están las decisiones que ya se tomaron — no se discuten de nuevo.

---

## Qué es PulseCity

PulseCity es una simulación interactiva browser-based de una ciudad ficticia generada proceduralmente, cuya economía, crecimiento urbano y vida social están directamente ligados al rendimiento de su franquicia de basketball ficticia en la NBA.

El jugador es el General Manager del equipo y alcalde de facto de la ciudad. Cada decisión tiene consecuencias reales y diferidas. Los agentes del sistema (coaches, jugadores, alcalde, prensa, sponsors) viven, reaccionan y acumulan estado de forma independiente al jugador.

Inspiraciones fundacionales: **Cities Skylines 2** (city builder, gestión urbana, emergencia de sistemas) + **NBA 2K** (franquicia, negociaciones, gestión de roster).

**Este es un proyecto personal de aprendizaje, sin deadline, sin usuarios externos, sin presión.** La complejidad no es un problema — es el punto.

---

## Estructura del repositorio

```
pulsecity/
├── services/
│   ├── map-service/        # Rust — generación procedural del mapa
│   ├── agent-service/      # Rust — motor de agentes y simulation loop
│   ├── match-service/      # Rust — simulación de partidos (stateless)
│   ├── gateway/            # Go  — API Gateway, WebSocket, auth, routing
│   ├── city-service/       # Go  — economía urbana, zonificación, edificios
│   ├── team-service/       # Go  — franquicia, roster, contratos, salary cap
│   ├── narrative-service/  # Go  — eventos narrativos via LLM
│   └── analytics-service/  # Go  — series temporales, TimescaleDB
├── shared/
│   ├── events/             # Definiciones de eventos NATS (Go + Rust)
│   └── types/              # Tipos compartidos entre servicios
├── frontend/               # React + WebGL (Three.js) Todo en TypeScript
├── infra/
│   └── k8s/                # Kubernetes manifests (k3s en desarrollo)
├── db/
│   └── migrations/         # Migraciones de TimescaleDB (PostgreSQL)
├── docker-compose.yml      # Entorno local completo
└── Makefile                # make dev, make build-all, make test-all
```

### Estructura interna — servicios Go

```
services/{nombre}/
├── cmd/
│   └── main.go             # Entry point: solo arranca el servidor
├── internal/
│   ├── handlers/           # HTTP + WebSocket handlers
│   ├── domain/             # Lógica de negocio
│   └── nats/               # Publicadores y suscriptores de eventos
├── go.mod
└── Dockerfile
```

### Estructura interna — servicios Rust

```
services/{nombre}/
├── src/
│   ├── main.rs             # Entry point
│   ├── lib.rs              # Módulo raíz
│   └── {módulos}/          # Particionados por dominio
├── Cargo.toml
└── Dockerfile
```

---

## Stack tecnológico

| Componente | Tecnología | Razón de la elección |
|---|---|---|
| Servicios CPU-intensivos | Rust | `map-service`, `agent-service`, `match-service` |
| Servicios de negocio / I/O | Go | `gateway`, `city-service`, `team-service`, `narrative-service`, `analytics-service` |
| Event Bus | NATS | Liviano, sin dependencias externas, no necesitamos replay ni retención |
| Base de datos | TimescaleDB (PostgreSQL) | Series temporales + familiaridad |
| Frontend | React + WebGL (Three.js) | Browser-based, visualización en tiempo real |
| Orquestación | Kubernetes (k3s local) | El punto del proyecto |
| CI/CD | GitHub Actions | |
| Deploy público | Railway | |
| LLM para narrativa | API Claude / GPT-4o mini | Texto emergente para eventos de agentes |

**NATS, no Kafka.** Decisión tomada. PulseCity no necesita replay de eventos ni retención de mensajes — el estado vive en TimescaleDB. No reabrir esta discusión.

---

## Arquitectura — principios que no se negocian

### 1. State ownership estricto

Cada dato tiene exactamente un servicio dueño. Ningún otro servicio escribe ese dato — solo lo lee via eventos NATS o queries directas al dueño.

Cuando el dato cambia, el dueño publica un evento en NATS. Los demás actualizan su propia vista si la necesitan.

**Ejemplo concreto — los jugadores tienen dos dueños distintos:**

| Dato | Dueño |
|---|---|
| Contrato, salary cap, stats en cancha, posición en roster | `team-service` |
| Estado emocional, confianza con GM, satisfacción, relaciones con otros agentes | `agent-service` |

Ambos servicios se referencian por `player_id`. Nunca se llaman directamente entre sí.

### 2. Comunicación solo via NATS (salvo consultas de lectura del frontend)

Los servicios no se llaman entre sí via HTTP. Se comunican exclusivamente publicando y consumiendo eventos NATS.

**Excepción:** el `gateway` puede hacer queries REST directas a servicios internos para ensamblar respuestas al frontend (ej: perfil completo de jugador = datos de `team-service` + estado emocional de `agent-service`).

### 3. Convención de nombres de eventos NATS

```
entidad.accion
```

Minúsculas, guión bajo como separador de palabras dentro de cada parte. Ejemplos:

```
tiempo.dia_avanzado
tiempo.sesion_iniciada
jugador.firmado
jugador.lesionado
jugador.solicita_traspaso
partido.terminado
partido.iniciando
agente.estado_cambio
agente.relacion_cambio
ciudad.suelo_actualizado
```

### 4. WebSocket — solo deltas

El backend nunca envía estado completo al frontend via WebSocket. Solo envía lo que cambió desde el último update. El frontend mantiene su propia copia del estado y la actualiza con cada delta.

Throttling: máximo un update por segundo en velocidad x1. Se relaja proporcionalmente a velocidades más altas.

### 5. `match-service` es completamente stateless

Recibe un payload completo con el estado de ambos equipos (ratings, cansancio, estado emocional de cada jugador, contexto táctico), simula el partido, y publica `partido.terminado` en NATS. Nunca consulta otros servicios durante la simulación. Todo lo que necesita llega en el input.

### 6. El `narrative-service` espera intencionalmente

Cuando `narrative-service` escucha un evento en NATS, espera **250-500ms** antes de generar el texto narrativo. Este delay garantiza que `agent-service` ya procesó el mismo evento y actualizó el estado. Cuando el LLM recibe el contexto del agente, el estado ya refleja las consecuencias. No es coordinación — es timing.

---

## Modelo de tiempo simulado

La unidad mínima es **un día simulado**.

| Velocidad | Tiempo real por día simulado | Tiempo real por temporada (~365 días) |
|---|---|---|
| Pausa | ∞ | ∞ |
| x1 | ~1.6 segundos | ~10 minutos |
| x5 | ~0.32 segundos | ~2 minutos |
| x20 | ~0.08 segundos | ~30 segundos |

### Simulation loop

Corre cada 100ms en tiempo real (componente interno de `agent-service`):

```
cada 100ms:
  delta_real = tiempo desde ultimo tick
  dias_a_procesar = delta_real * velocidad_actual / 1.6s

  si dias_a_procesar >= 1:
    procesar_dia(fecha_actual)
    publicar tiempo.dia_avanzado en NATS
    fecha_actual += 1
```

### Loop híbrido

- **Con sesión activa** (WebSocket abierto desde el frontend): el loop avanza el tiempo.
- **Sin sesión activa**: el loop duerme. El tiempo no avanza. Al volver, reanuda exactamente donde quedó.

El `gateway` publica `tiempo.sesion_iniciada` y `tiempo.sesion_terminada` cuando abre/cierra una conexión WebSocket.

---

## Milestones de desarrollo

PulseCity no tiene MVP. Tiene milestones incrementales — cada uno es jugable, mostrable, y construye sobre el anterior sin tirar nada.

### Milestone 1 — El mundo nace
El jugador crea su franquicia, ve su ciudad generarse en tiempo real, y recibe la primera llamada del dueño.

**Servicios activos:** `gateway`, `map-service`, `narrative-service` (primer evento LLM obligatorio).

**Criterio de done:** El mapa se genera, llega via WebSocket al browser con la ceremonia visual, y aparece el primer evento narrativo del Owner.

### Milestone 2 — El sistema late
El loop central funcionando. Partidos simulándose, tiempo avanzando, ciudad reaccionando.

**Servicios activos:** todos los anteriores + `agent-service` (5 agentes core), `team-service`, `match-service`, `city-service` básico, `analytics-service`.

**Criterio de done:** 82 partidos se simulan en una temporada completa con box scores y eventos narrativos post-partido. Los 5 agentes core (Owner, Head Coach, CFO, Director de Scouting, Sports Psychologist) tienen personalidad y estado real.

### Milestone 3 — Los agentes viven
El corazón narrativo del sistema en su totalidad.

**Qué funciona:** los 30 agentes individuales completos + roster de 15 jugadores como agentes, sistema de relaciones, causalidad real (el sistema recuerda decisiones pasadas), canal de consulta directa con cada agente, ciclo jugable completo (Draft → Playoffs → Cierre).

### Milestone 4 — La ciudad respira
El loop ciudad-franquicia en toda su profundidad.

**Qué funciona:** edificios especiales one-time con trade-offs reales, modo "Dueño con influencia" con el alcalde como agente real con agenda propia, todos los agentes de ciudad activos, analytics avanzados con visualizaciones históricas de series temporales.

---

## Catálogo de servicios

### `map-service` (Rust)
Generación procedural del mapa usando Perlin Noise y Voronoi. Solo se ejecuta al crear una nueva partida. Emite progreso via WebSocket durante la generación (ceremonia visual). CPU-intensivo en la inicialización — stateless después.

### `agent-service` (Rust)
Motor de agentes. Contiene el simulation loop. Gestiona:
- ~50 agentes individuales con memoria, estado emocional y relaciones
- Población agregada (ciudadanos como stats por zona, no entidades individuales)
- Agentes visibles en el mapa (sprites con paths precalculados)

Los agentes individuales reaccionan a eventos — no se evalúan entre sí en cada tick.

### `match-service` (Rust)
Simula partidos. Completamente stateless — recibe payload, simula, publica `partido.terminado`. Puede simular cientos de partidos en segundos.

### `gateway` (Go)
Punto de entrada único. Maneja WebSocket, REST, autenticación (login, registro, guest tokens), y routing a servicios internos. Es el único servicio que el frontend conoce.

### `city-service` (Go)
Economía urbana, zonificación, edificios especiales, presupuesto municipal, valor del suelo. Reacciona a resultados de partidos y decisiones del GM.

### `team-service` (Go)
Franquicia, roster, jugadores ficticios, finanzas del equipo, salary cap, staff técnico y de negocios. Dueño de todos los datos contractuales.

### `narrative-service` (Go + LLM API)
Genera el texto de los eventos narrativos. Escucha eventos NATS, espera 250-500ms, consulta el contexto actualizado de los agentes involucrados, y llama al LLM. Gestiona la bandeja de entrada del jugador y el historial de eventos.

### `analytics-service` (Go + TimescaleDB)
Ingesta y sirve series temporales: valor del suelo por tick, asistencia por partido, ingresos semanales, rendimiento de jugadores por temporada. TimescaleDB es PostgreSQL — el schema y las queries son SQL estándar.

---

## Flujo de un evento típico — firma de jugador

```
1. GM confirma la firma en el frontend
2. gateway recibe POST /actions/sign-player
3. gateway publica jugador.firma_iniciada en NATS
4. team-service → actualiza roster y contrato → publica jugador.firmado
5. agent-service → actualiza estado emocional de Coach, CFO, Director de Marketing, etc.
6. city-service → sube valor del suelo cerca del estadio, actualiza ticket sales proyectadas
7. narrative-service → espera 250-500ms → genera evento narrativo con contexto actualizado
8. gateway → escucha jugador.firmado → envía delta al frontend via WebSocket
```

---

## Lo que este proyecto NO es

- No es un CRUD API con algunos endpoints. Es un sistema distribuido con simulación en tiempo real.
- No tiene "MVP" ni fecha de entrega. Cada milestone es un estado jugable completo.
- No usamos Kafka. NATS es suficiente y la decisión está tomada.
- No enviamos estado completo por WebSocket. Solo deltas.
- Los servicios no se llaman entre sí via HTTP (salvo el gateway para ensamblar respuestas al frontend).

---

## Documentos de diseño

Para contexto más profundo, ver los documentos en la raíz del proyecto:

- `PulseCity_v2.md` — Diseño completo: concepto, mundo, mecánicas, stack
- `PulseCity_Arquitectura_Tecnica.md` — Decisiones técnicas detalladas, catálogo de eventos NATS, modelo de partidos
- `PulseCity_Agentes.md` — Los ~50 agentes individuales: quiénes son, su personalidad, sus motivaciones
- `PulseCity_Mecanicas_Agentes.md` — Qué hacen mecánicamente los agentes, sus outputs concretos
- `PulseCity_CicloJugable.md` — El ciclo de una temporada completa fase por fase
- `PulseCity_Experiencia_Jugador.md` — La experiencia desde el punto de vista del jugador

---

## Sobre Obsidian y las notas

Luego de cada sesion y avance, registrar en /docs, los cambios hechos en las carpetas donde amerite. Si es necesario crear nueva carpeta.
El objetivo es siempre tener mi vault de obsidian actualizado con info pertinente y bien organizada.

*Generado en sesión de diseño — Mayo 2026.*
*Actualizar este archivo cuando cambien decisiones de arquitectura o estructura de carpetas.*
