---
name: pulsecity-context
description: >
  Contexto completo de arquitectura y diseño de PulseCity. Úsala siempre que estés
  trabajando en cualquier servicio del proyecto: map-service, agent-service, team-service,
  match-service, city-service, narrative-service, gateway, o el frontend React.
  Actívala cuando el usuario mencione PulseCity, cualquier servicio del sistema,
  eventos NATS del juego, agentes, franquicia, simulación de tiempo, o milestones.
---

# PulseCity — Contexto de Arquitectura

PulseCity es una simulación browser-based de gestión de franquicia NBA + city builder.
El jugador es GM de una franquicia expansión (equipo 33). La ciudad reacciona a las
decisiones del GM y los resultados en cancha en tiempo real.

## Stack completo

| Componente | Tecnología | Rol |
|---|---|---|
| map-service | Rust | Generación procedural (Perlin Noise + Voronoi) |
| agent-service | Rust | Simulación de ~50 agentes + simulation loop |
| match-service | Rust | Simulación stateless de partidos |
| city-service | Go | Estado de la ciudad, economía, edificios |
| narrative-service | Go + LLM API | Generación de texto narrativo |
| gateway | Go | API REST + WebSocket hub |
| analytics-service | Go + TimescaleDB | Series temporales, stats históricas |
| frontend | React + WebGL | Display de deltas via WebSocket |
| Event bus | NATS | Comunicación entre servicios (NO Kafka) |
| Base de datos | TimescaleDB (PostgreSQL) | Estado persistido + series temporales |
| Orquestación | k3s (local) / Railway (prod) | |
| CI/CD | GitHub Actions | |

## Regla de oro: state ownership

Cada dato tiene exactamente UN servicio dueño. Los demás solo leen via NATS o query directa.
Nunca escribir el dato de otro servicio. Nunca llamar directamente a otro servicio para escribir.

### Quién es dueño de qué

**team-service** (Go):
- Contratos (salario, años, cláusulas)
- Roster (activo / G-League / inactivo / lesionado)
- Salary cap
- Stats de cancha (box score acumulado)
- Contratos del staff + métricas objetivas

**agent-service** (Rust):
- Estado emocional de todos los agentes
- `confianza` (-1.0 a 1.0), `satisfaccion` (-1.0 a 1.0), `lealtad` (0.0 a 1.0)
- Relaciones entre agentes (tabla de pares)
- Agenda propia de cada agente
- Historial emocional

**city-service** (Go):
- Valor del suelo por zona
- Economía urbana (indicadores agregados)
- Edificios construidos y su estado

**map-service** (Rust):
- Geometría del mapa (solo escritura durante generación inicial)
- Zonas y barrios generados

**match-service** (Rust):
- Stateless: recibe input, devuelve output, no persiste nada

## Catálogo de eventos NATS

### Tiempo
- `tiempo.dia_avanzado` → simulation-loop publica; todos escuchan
- `tiempo.sesion_iniciada` / `tiempo.sesion_terminada` → gateway publica; simulation-loop escucha
- `tiempo.pausa_activada` / `tiempo.velocidad_cambiada` → gateway publica

### Jugadores
- `jugador.firmado` → team-service publica; agent-service, city-service, narrative-service, gateway escuchan
- `jugador.traspasado` → team-service publica
- `jugador.lesionado` / `jugador.recuperado` → team-service publica
- `jugador.solicita_traspaso` → agent-service publica
- `jugador.contrato_expira_pronto` → team-service publica

### Partidos
- `partido.programado` → team-service publica; match-service escucha
- `partido.iniciando` → match-service publica; gateway escucha
- `partido.terminado` → match-service publica; agent-service, city-service, analytics-service, narrative-service, team-service, gateway escuchan

### Agentes
- `agente.estado_cambio` → agent-service publica; narrative-service, gateway escuchan
- `agente.relacion_cambio` → agent-service publica; narrative-service escucha
- `agente.evento_critico` → agent-service publica; narrative-service, gateway escuchan

### Ciudad
- `ciudad.suelo_actualizado` → city-service publica
- `ciudad.edificio_construido` → city-service publica; agent-service, narrative-service, gateway escuchan
- `ciudad.economia_cambio` → city-service publica

### Narrativa
- `narrativa.evento_generado` → narrative-service publica; gateway escucha
- `narrativa.respuesta_gm` → gateway publica; narrative-service, agent-service escuchan

### Mapa (solo durante generación inicial)
- `mapa.terreno_listo` / `mapa.zonas_calculadas` / `mapa.estadio_ubicado` / `mapa.generacion_completa`

## Modelo de tiempo simulado

- Unidad mínima: **1 día simulado**
- x1 speed: 1.6s real = 1 día simulado → temporada completa en ~10 min
- x5: ~0.32s/día | x20: ~0.08s/día
- Tick del simulation-loop: cada 100ms en tiempo real
- Loop **híbrido**: corre con WebSocket activo, duerme sin sesión

```
cada 100ms:
  dias = delta_real * velocidad / 1.6s
  si dias >= 1:
    procesar_dia(fecha_actual)
    NATS.publish("tiempo.dia_avanzado")
    fecha_actual += 1
```

## WebSocket: solo deltas

El gateway NUNCA envía estado completo. Solo envía lo que cambió desde el último update.
El frontend mantiene su propia copia del estado y aplica los deltas recibidos.
Throttling: máximo 1 update/segundo en x1, se relaja proporcionalmente a mayor velocidad.

## Patrón de flujo estándar (ejemplo: firma de jugador)

```
1. GM decide en frontend
2. gateway recibe REST → publica jugador.firma_iniciada en NATS
3. team-service escucha → actualiza → publica jugador.firmado
4. agent-service escucha jugador.firmado → actualiza estado emocional
5. city-service escucha jugador.firmado → actualiza economía
6. narrative-service escucha jugador.firmado → espera 250-500ms → genera texto con LLM
7. gateway escucha jugador.firmado → envía delta al frontend via WebSocket
```

El delay de 250-500ms en narrative-service NO es coordinación, es timing:
garantiza que agent-service ya actualizó cuando el LLM recibe el contexto.

## Gateway como combinador (única excepción REST interna)

```
GET /players/{id}/profile
  → query team-service (datos contractuales)  ← en paralelo
  → query agent-service (estado emocional)    ←
  → merge → devuelve al frontend
```

Esta es la ÚNICA excepción al patrón NATS-first. Las lecturas para el frontend
pueden ir directo al gateway por REST, que consulta internamente en paralelo.

## match-service: stateless puro

Nunca consulta otros servicios durante simulación. Recibe payload completo:
- Ratings de jugadores (general, anotacion, defensa, playmaking, atletismo, tiro, clutch)
- Estado del día (forma_fisica, cansancio_acumulado, estado_emocional)
- Contexto táctico (sistema, flexibilidad, rotación)

Devuelve: box score completo + 3-5 momentos narrativos clave.

## Los ~50 agentes individuales

**Basketball Ops (~15):** Owner, President of Basketball Ops, Assistant GMs,
Director de Scouting, Director de Player Personnel, Head of Analytics, Head Coach,
asistentes del coach (defensa/ataque/desarrollo), Director de Player Development,
Médico, Fisioterapeuta/S&C Coach, Sports Psychologist, Video Coordinator, International Scout.

**Business Ops (~8):** CEO/President Business Ops, CFO, Director Marketing & Brand,
Director Ticket Sales, Director Corporate Partnerships & Sponsors,
Director Community Relations, Director Arena Operations, Head of Digital & Social Media.

**Roster (15 jugadores):** cada jugador es también un agente con estado emocional.

**Ciudad (~7):** Alcalde, Jefe de Policía, Presidenta Cámara de Comercio,
Director de Urbanismo y Planeamiento, Líder Comunitario Barrio del Estadio,
Director del Fondo de Inversión Inmobiliaria, Jefe de Sindicato de Trabajadores.

**The Press (1 agente colectivo):** reacciona a resultados y decisiones públicas.

**30 GMs rivales:** uno por cada equipo NBA real.

## Milestones de desarrollo

- **M1** — El mundo nace: mapa generado en tiempo real via WebSocket, identidad visual
- **M2** — El sistema late: loop de tiempo, 82 partidos, ciudad reacciona, 5 agentes core
- **M3** — Los agentes viven: 50 agentes completos, relaciones, causalidad, ciclo completo
- **M4** — La ciudad respira: edificios especiales, alcalde con agenda propia, analytics avanzados

Cada milestone es jugable y mostrable independientemente.

## Decisiones de arquitectura tomadas (no reabrir)

- **NATS sobre Kafka**: no se necesita replay de eventos ni retención. Estado en TimescaleDB.
- **Rust para simulación**: map-service, agent-service, match-service son CPU-intensivos.
- **Go para coordinación**: city-service, narrative-service, gateway — I/O y HTTP.
- **LLM solo para lenguaje**: la lógica de los agentes es rule-based. El LLM genera el texto.
- **Milestones, no MVP**: cada milestone es jugable. No hay "versión mínima" sin valor.
