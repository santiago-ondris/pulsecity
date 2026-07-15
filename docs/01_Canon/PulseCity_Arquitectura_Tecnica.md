# PulseCity — Arquitectura Técnica

> Documento de referencia técnico. Complementa PulseCity_v2.md, PulseCity_Agentes.md y PulseCity_CicloJugable.md.
> Contiene las decisiones de arquitectura definidas en sesión de diseño — Mayo 2026.
> Este es el documento al que volvés cuando estás codeando y necesitás saber quién hace qué.

---

## Modelo de tiempo simulado

### Unidad base

La unidad mínima de tiempo simulado es **un día**. Toda la lógica del sistema se expresa en días simulados — los partidos ocurren en días específicos, los agentes reaccionan en días específicos, los efectos en la ciudad se calculan por día.

### Velocidades

| Velocidad | Tiempo real por día simulado | Tiempo real por temporada (~365 días) |
|---|---|---|
| Pausa | ∞ | ∞ |
| x1 | ~1.6 segundos | ~10 minutos |
| x5 | ~0.32 segundos | ~2 minutos |
| x20 | ~0.08 segundos | ~30 segundos |

La temporada completa (Draft → Playoffs → Cierre) cubre aproximadamente 365 días simulados.

### Loop híbrido

El backend corre en modo **híbrido**:

- **Con sesión activa** — el loop de simulación avanza el tiempo continuamente según la velocidad elegida por el jugador. El mundo vive.
- **Sin sesión activa** — el loop entra en sleep. El tiempo no avanza. Cuando el jugador vuelve, el tiempo reanuda exactamente donde quedó.

Una "sesión activa" se define por la presencia de una conexión WebSocket abierta desde el frontend. El gateway notifica al simulation-loop cuando abre y cierra una sesión.

### Tick del simulation loop

El `simulation-loop` (componente interno del `agent-service` o servicio propio según evolución) corre cada 100ms en tiempo real. En cada tick calcula cuánto tiempo simulado avanzó desde el tick anterior según la velocidad activa y avanza el estado del mundo ese delta. Si el delta acumulado supera un día simulado, procesa ese día — eventos programados, reacciones de agentes, efectos en la ciudad.

```
cada 100ms:
  delta_real = tiempo desde ultimo tick
  dias_simulados_a_procesar = delta_real * velocidad_actual / 1.6s
  
  si dias_simulados_a_procesar >= 1:
    procesar_dia(fecha_actual)
    publicar partido.dia_avanzado en NATS
    fecha_actual += 1
```

---

## State ownership entre servicios

### Principio fundamental

Cada dato tiene exactamente un servicio dueño. Ningún otro servicio puede escribir ese dato — solo puede leerlo via eventos de NATS o consultas directas al servicio dueño. Cuando el dato cambia, el servicio dueño publica un evento en NATS y los demás actualizan su propia vista si la necesitan.

### División de responsabilidades — jugadores

Los jugadores son la entidad más compleja del sistema porque tienen datos en dos dominios muy distintos.

**`team-service` es dueño de:**
- Contrato (salario, años restantes, tipo de contrato, cláusulas especiales)
- Posición en el roster (activo / G-League / inactivo / lesionado)
- Salary cap implications
- Historial de traspasos y fichajes
- Stats de rendimiento en cancha (box score acumulado)

**`agent-service` es dueño de:**
- Estado emocional actual
- Confianza con el GM (`confianza`: -1.0 a 1.0)
- Satisfacción con su situación (`satisfaccion`: -1.0 a 1.0)
- Lealtad a la franquicia (`lealtad`: 0.0 a 1.0)
- Relaciones con otros agentes (tabla de pares)
- Agenda propia (ambición, deseos de protagonismo, etc.)
- Historial de eventos emocionales significativos

**Compartido solo via `player_id`** — ambos servicios referencian al mismo jugador por ID pero nunca se llaman directamente entre sí.

### División de responsabilidades — staff

Misma lógica aplicada al staff de Basketball Ops y Business Ops.

**`team-service` es dueño de:**
- Contrato del staff (salario, años restantes)
- Rol y responsabilidades formales
- Métricas de rendimiento objetivo (precisión histórica del scouting, retention rate de sponsors, etc.)

**`agent-service` es dueño de:**
- Estado emocional y relacional de cada miembro del staff
- Confianza, satisfacción, lealtad con el GM
- Relaciones con otros agentes del staff

### Coordinación entre servicios

Los servicios nunca se llaman directamente. Se comunican exclusivamente via NATS.

**Ejemplo — se firma un jugador estrella:**

```
1. GM toma la decisión en el frontend
2. gateway recibe la acción via REST
3. gateway publica jugador.firma_iniciada en NATS con el payload completo
4. team-service escucha → actualiza contrato y roster → publica jugador.firmado
5. agent-service escucha jugador.firmado → actualiza estado emocional del Director de Marketing, CFO, Head Coach, etc.
6. city-service escucha jugador.firmado → sube valor del suelo near estadio, actualiza ticket sales proyectadas
7. narrative-service escucha jugador.firmado → después de un delay intencional (para dar tiempo a agent-service), genera el evento narrativo con el contexto actualizado
8. gateway escucha jugador.firmado → envía delta al frontend via WebSocket
```

El delay intencional del `narrative-service` (250-500ms) garantiza que cuando genera el texto, el estado de los agentes ya fue actualizado por `agent-service`. No es coordinación — es timing.

### El gateway como combinador para el frontend

Cuando el frontend necesita mostrar el perfil completo de un jugador (datos de contrato + estado emocional), el gateway hace dos consultas en paralelo y combina el resultado:

```
GET /players/{id}/profile
  → query team-service: datos contractuales
  → query agent-service: estado emocional y relacional
  → merge y devuelve al frontend
```

Esto es la única excepción al patrón NATS-first — las consultas de lectura para el frontend pueden ser REST directas al gateway, que internamente consulta a los servicios correspondientes.

---

## Modelo de partidos

### Inputs del `match-service`

El `match-service` recibe un payload completo antes de simular cada partido. Nunca consulta otros servicios durante la simulación — todo lo que necesita llega en el input.

```json
{
  "partido_id": "uuid",
  "fecha_simulada": "2026-01-15",
  "equipo_local": {
    "franquicia_id": "uuid",
    "jugadores": [
      {
        "player_id": "uuid",
        "nombre": "string",
        "posicion": "PG|SG|SF|PF|C",
        "ratings": {
          "general": 0.0,
          "anotacion": 0.0,
          "defensa": 0.0,
          "playmaking": 0.0,
          "atletismo": 0.0,
          "tiro": 0.0,
          "clutch": 0.0
        },
        "estado_dia": {
          "forma_fisica": 0.0,
          "cansancio_acumulado": 0.0,
          "estado_emocional": 0.0
        },
        "minutos_asignados": 0
      }
    ],
    "contexto_tactico": {
      "sistema": "pace_and_space|defense_first|iso_ball|team_system|hybrid",
      "flexibilidad_tactica": 0.0,
      "rotacion": ["player_id_1", "player_id_2"]
    }
  },
  "equipo_visitante": {
    // misma estructura
  }
}
```

El payload es ensamblado por el `gateway` o un orquestador interno que consulta `team-service` y `agent-service` antes de enviarlo al `match-service`. El `match-service` en sí es stateless — recibe, simula, devuelve.

### Output del `match-service`

El `match-service` publica `partido.terminado` en NATS con el siguiente payload:

```json
{
  "partido_id": "uuid",
  "fecha_simulada": "2026-01-15",
  "resultado": {
    "local": { "franquicia_id": "uuid", "puntos": 0 },
    "visitante": { "franquicia_id": "uuid", "puntos": 0 }
  },
  "box_score": {
    "local": [
      {
        "player_id": "uuid",
        "minutos": 0,
        "puntos": 0,
        "rebotes": 0,
        "asistencias": 0,
        "robos": 0,
        "tapones": 0,
        "perdidas": 0,
        "fg_intentados": 0,
        "fg_convertidos": 0,
        "tres_intentados": 0,
        "tres_convertidos": 0,
        "tl_intentados": 0,
        "tl_convertidos": 0,
        "plus_minus": 0,
        "eficiencia_ofensiva": 0.0,
        "eficiencia_defensiva": 0.0
      }
    ],
    "visitante": []
  },
  "momentos_clave": [
    {
      "tipo": "clutch_anotacion|clutch_fallo|actuacion_sorpresa|colapso_individual|liderazgo",
      "player_id": "uuid",
      "descripcion_corta": "string",
      "impacto_emocional": 0.0
    }
  ]
}
```

Los `momentos_clave` son 3 a 5 por partido. Son el input principal del `narrative-service` para generar los eventos narrativos post-partido, y del `agent-service` para actualizar el estado emocional de los jugadores involucrados.

### Efectos en cadena de `partido.terminado`

Cuando `partido.terminado` se publica en NATS, los siguientes servicios reaccionan:

- **`agent-service`** — actualiza estado emocional de jugadores según su actuación y el resultado. Un jugador que falló en clutch acumula presión. Un jugador con actuación sorpresa sube su satisfacción.
- **`city-service`** — actualiza economía urbana. Victoria en casa → ciudadanos activos al día siguiente, ticket sales suben. Racha mala → valor del suelo empieza a bajar lentamente.
- **`analytics-service`** — ingesta el box score completo en TimescaleDB para series temporales y análisis histórico.
- **`narrative-service`** — después de delay intencional, genera el evento narrativo post-partido para la bandeja de entrada del jugador.
- **`team-service`** — actualiza stats acumuladas de la temporada, récord del equipo.

---

## Catálogo de eventos NATS

Todos los eventos siguen la convención `entidad.accion` en minúsculas con guión bajo.

### Eventos de tiempo

| Evento | Publicado por | Consumido por | Payload clave |
|---|---|---|---|
| `tiempo.dia_avanzado` | simulation-loop | todos | `fecha_simulada`, `velocidad` |
| `tiempo.sesion_iniciada` | gateway | simulation-loop | `partida_id` |
| `tiempo.sesion_terminada` | gateway | simulation-loop | `partida_id` |
| `tiempo.pausa_activada` | gateway | simulation-loop | `partida_id` |
| `tiempo.velocidad_cambiada` | gateway | simulation-loop | `velocidad` |

### Eventos de jugadores

| Evento | Publicado por | Consumido por | Payload clave |
|---|---|---|---|
| `jugador.firmado` | team-service | agent-service, city-service, narrative-service, gateway | `player_id`, `franquicia_id`, `salario`, `anos` |
| `jugador.traspasado` | team-service | agent-service, city-service, narrative-service, gateway | `player_id`, `franquicia_origen`, `franquicia_destino` |
| `jugador.lesionado` | team-service | agent-service, narrative-service, gateway | `player_id`, `severidad`, `dias_estimados` |
| `jugador.recuperado` | team-service | agent-service, narrative-service, gateway | `player_id` |
| `jugador.solicita_traspaso` | agent-service | narrative-service, gateway | `player_id`, `motivo` |
| `jugador.contrato_expira_pronto` | team-service | narrative-service, gateway | `player_id`, `dias_restantes` |

### Eventos de partidos

| Evento | Publicado por | Consumido por | Payload clave |
|---|---|---|---|
| `partido.programado` | team-service | match-service | `partido_id`, `fecha`, `equipos` |
| `partido.iniciando` | match-service | gateway | `partido_id` |
| `partido.terminado` | match-service | agent-service, city-service, analytics-service, narrative-service, team-service, gateway | payload completo definido arriba |

### Eventos de agentes

| Evento | Publicado por | Consumido por | Payload clave |
|---|---|---|---|
| `agente.estado_cambio` | agent-service | narrative-service, gateway | `agente_id`, `variable`, `valor_anterior`, `valor_nuevo` |
| `agente.relacion_cambio` | agent-service | narrative-service | `agente_a`, `agente_b`, `confianza_nueva`, `motivo` |
| `agente.evento_critico` | agent-service | narrative-service, gateway | `agente_id`, `tipo`, `urgencia` |

### Eventos de ciudad

| Evento | Publicado por | Consumido por | Payload clave |
|---|---|---|---|
| `ciudad.suelo_actualizado` | city-service | analytics-service, gateway | `zona_id`, `valor_anterior`, `valor_nuevo` |
| `ciudad.edificio_construido` | city-service | agent-service, narrative-service, gateway | `edificio_tipo`, `zona_id` |
| `ciudad.economia_cambio` | city-service | analytics-service, gateway | `indicador`, `valor`, `tendencia` |

### Eventos de narrativa

| Evento | Publicado por | Consumido por | Payload clave |
|---|---|---|---|
| `narrativa.evento_generado` | narrative-service | gateway | `evento_id`, `agente_emisor`, `tipo`, `urgencia`, `texto`, `opciones_respuesta` |
| `narrativa.respuesta_gm` | gateway | narrative-service, agent-service | `evento_id`, `opcion_elegida` |

### Eventos de mapa (solo durante generación inicial)

| Evento | Publicado por | Consumido por | Payload clave |
|---|---|---|---|
| `mapa.terreno_listo` | map-service | gateway | `partida_id`, `datos_terreno` |
| `mapa.zonas_calculadas` | map-service | gateway | `partida_id`, `zonas` |
| `mapa.estadio_ubicado` | map-service | city-service, gateway | `partida_id`, `coordenadas_estadio` |
| `mapa.generacion_completa` | map-service | gateway | `partida_id` |

---

## Milestones de desarrollo

El proyecto se construye en milestones incrementales. Cada milestone es jugable, mostrable, y construye sobre el anterior sin tirar nada.

### Milestone 1 — El mundo nace

**Objetivo:** La primera experiencia emocional del jugador. Ver nacer su ciudad.

**Qué funciona:**
- Autenticación (cuenta + guest token)
- Creación de franquicia: nombre, colores, logo básico, uniforme
- Selección de estado inicial y modo de ciudad
- Generación del mapa con ceremonia visual en tiempo real (WebSocket con eventos `mapa.*`)
- La llamada del dueño como primer evento LLM

**Servicios activos:** `map-service`, `gateway`, `team-service` (solo creación), `narrative-service` (solo el primer evento)

**Criterio de done:** Un jugador puede crear una franquicia, ver su ciudad generarse en tiempo real, y recibir la llamada del dueño.

---

### Milestone 2 — El sistema late

**Objetivo:** El loop central funcionando. Partidos, tiempo, ciudad reaccionando.

**Qué funciona:**
- Loop de tiempo híbrido con velocidades x1/x5/x20/pausa
- Calendario de temporada regular con 82 partidos simulados
- `match-service` simulando partidos con box scores completos y momentos clave
- `city-service` básico reaccionando a resultados (valor del suelo, ticket sales)
- 5 agentes core con estado y personalidad: Owner, Head Coach, CFO, Director de Scouting, Sports Psychologist
- Eventos narrativos básicos post-partido
- Frontend mostrando el mapa con la ciudad reaccionando en tiempo real

**Servicios activos:** todos los del Milestone 1 + `match-service`, `city-service`, `agent-service` (5 agentes core), `analytics-service` (básico)

**Criterio de done:** Se puede simular una temporada completa y ver la ciudad reaccionar a los resultados partido a partido.

---

### Milestone 3 — Los agentes viven

**Objetivo:** El corazón narrativo del sistema funcionando en su totalidad.

**Qué funciona:**
- Los 30 agentes individuales completos con estado, personalidad y agenda propia
- El roster completo de 15 jugadores como agentes
- Sistema de relaciones entre agentes (Coach vs Médico, Analytics vs Scouting, etc.)
- Causalidad real — el sistema recuerda decisiones pasadas y las conecta con consecuencias futuras
- Canal de consulta directa con cada agente (chat in-game anclado a dominio)
- Ciclo jugable completo: Draft, Free Agency, Traspasos, Pretemporada, Temporada Regular, Playoffs, Cierre
- Eventos narrativos con toda la complejidad del sistema de agentes

**Servicios activos:** todos los anteriores con `agent-service` completo y `narrative-service` completo

**Criterio de done:** Se puede jugar una temporada completa con todos los agentes activos, relaciones evolucionando, y causalidad funcionando entre decisiones y consecuencias.

---

### Milestone 4 — La ciudad respira

**Objetivo:** El loop ciudad-franquicia en toda su profundidad.

**Qué funciona:**
- Edificios especiales one-time con todos sus trade-offs
- Modo "Dueño con influencia" — el alcalde con agenda propia, puede decir que no
- Efectos cruzados completos entre decisiones de ciudad y decisiones de franquicia
- Todos los agentes de ciudad activos (Alcalde, Jefe de Policía, Cámara de Comercio, Director de Urbanismo)
- Analytics avanzados y visualizaciones históricas de series temporales

**Criterio de done:** Ambos modos de ciudad son jugables con toda la profundidad diseñada. Los edificios especiales tienen impacto real y visible.

---

### Juego completo

- Equipo rival (equipo 32) manejado por IA con su propia ciudad y ciclo de decisiones
- Liga completa con 32 equipos simulados
- Vista isométrica como segunda capa visual
- Logo builder completo y uniforme builder completo
- Ciudades reales via OpenStreetMap
- Modo multijugador (dos GMs, dos ciudades, misma liga)
- Sistema de legado — varios GMs a lo largo de décadas en la misma franquicia

---

## Notas de implementación

### El delay intencional del narrative-service

Cuando `narrative-service` escucha un evento en NATS, espera 250-500ms antes de generar el texto narrativo. Este delay garantiza que `agent-service` ya procesó el mismo evento y actualizó el estado de los agentes. Cuando el LLM recibe el contexto del agente para generar el texto, el estado ya refleja la consecuencia del evento.

No es coordinación — es timing. Simple y suficiente para el volumen de PulseCity.

### WebSocket — solo deltas

El backend nunca envía estado completo al frontend via WebSocket. Solo envía deltas — lo que cambió desde el último update. El frontend mantiene su propia copia del estado y la actualiza con cada delta recibido.

Throttling inteligente: el frontend no procesa más de un update por segundo en velocidad x1. A velocidades más altas, el throttling se relaja proporcionalmente.

### NATS vs Kafka — decisión tomada

NATS es el event bus de PulseCity. La decisión se tomó priorizando simplicidad operacional sobre features avanzadas. PulseCity no necesita replay de eventos (el estado está persistido en TimescaleDB), no necesita retención de mensajes días después, y no tiene el volumen que justifica la complejidad operacional de Kafka. NATS es un binario único sin dependencias externas — arranca en segundos y hace exactamente lo que PulseCity necesita.

### Coordinación de actualizaciones multi-servicio

Cuando un evento afecta múltiples servicios simultáneamente, cada servicio escucha el evento en NATS y actualiza su propio estado de forma independiente. No hay orquestador central. Esta es la arquitectura event-driven pura — servicios desacoplados que reaccionan a hechos, no a instrucciones.

---

*Documento generado en sesión de diseño — Mayo 2026.*
*Complementa PulseCity_v2.md, PulseCity_Agentes.md y PulseCity_CicloJugable.md — leer en conjunto.*
