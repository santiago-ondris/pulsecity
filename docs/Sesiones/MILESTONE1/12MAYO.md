# Sesión — 2026-05-12

## Estado de la sesión

Tercera sesión de trabajo sobre `Milestone 1`.

El foco real terminó siendo doble:

- volver real la primera decisión narrativa del Owner
- rehacer por completo el inicio de partida para dejar atrás la pantalla única y empezar a ordenar el frontend por features

---

## Mini milestone 1 — Primera respuesta real al Owner

Se cerró el primer corte narrativo interactivo con persistencia real.

### Qué cambió

En `services/gateway`:

- se agregó `POST /api/v1/games/{gameID}/owner-intro-response`
- el `owner_intro_event` ahora ofrece 3 opciones reales:
  - `build_culture`
  - `win_now`
  - `city_first`
- la respuesta elegida por el GM se persiste como `owner_intro_response`
- al responder, el `gateway` publica `narrativa.respuesta_gm` en NATS

En persistencia:

- se agregó `db/migrations/005_add_owner_intro_response.sql`
- `EnsureSchema` contempla también este nuevo campo para no romper entornos ya creados

En `frontend`:

- el modal del Owner dejó de ser un dismiss visual
- cada elección hace `POST` real al `gateway`
- la respuesta queda reflejada en la partida
- si la partida se rehidrata, el evento no vuelve a abrirse cuando ya fue respondido

### Resultado funcional

La fundación de la partida ahora deja una postura inicial guardada del GM.
El arranque ya tiene una consecuencia persistida y un evento NATS reutilizable más adelante.

---

## Mini milestone 2 — Rework completo del onboarding inicial

Se descartó el intento anterior de “wizard dentro de una sola pantalla” y se rehízo la entrada desde cero.

### Qué quedó implementado

En `frontend`:

- ahora existe una pantalla inicial de `PulseCity` con un único CTA:
  - `Nueva partida`
- el flujo de setup pasa a ser multipágina real dentro de la SPA:
  - landing
  - identidad visual
  - escenario inicial
  - gestión de ciudad
  - revisión/fundación
  - ceremonia del mapa
- el usuario ya no puede saltarse pasos ni avanzar fuera de orden
- cada pantalla resuelve una sola tarea clara
- la ceremonia quedó separada como página propia y ya no compite con el setup previo

### Reorganización de carpetas

También se empezó a sacar lógica de `App.tsx`.

Se creó un feature dedicado:

```text
frontend/src/features/new-game/
```

con separación por:

- `components/`
- `hooks/`
- `constants.ts`
- `helpers.ts`
- `types.ts`
- `newGame.css`

`App.tsx` ahora queda como punto de entrada mínimo que compone el feature.

### Decisión de arquitectura frontend

No se incorporó todavía un router externo.

Se resolvió con navegación interna basada en paths y guardas de progreso porque:

- ya da una experiencia multipágina real
- permite imponer orden secuencial desde ahora
- evita abrir todavía otra dependencia estructural
- deja el frontend mucho mejor preparado para pasar a routing formal más adelante si hace falta

---

## Verificación

Se validó con:

```bash
cd services/gateway
GOCACHE=/tmp/pulsecity-gocache go test ./...

cd frontend
npm run build
```

Resultado:

- `go test` OK en `gateway`
- build OK en `frontend`

---

## Qué habilita después

Este cierre deja terreno limpio para el siguiente tramo:

- seguir separando features sin inflar `App.tsx`
- mejorar visualmente cada página del onboarding sin mezclar todo otra vez
- usar `narrativa.respuesta_gm` como input real para agentes, tono narrativo y lógica de ciudad

---

## Lectura actual de Milestone 1

Después de revisar el canon y el alcance real de `Milestone 1`, la conclusión de esta sesión es:

### Lo que ya está razonablemente encaminado

- creación de franquicia con datos fundacionales persistidos
- flujo inicial de `Nueva Partida` mucho más claro que antes
- ceremonia de generación del mapa conectada al backend real
- recepción de eventos de mapa vía WebSocket
- primera llamada del Owner integrada al flujo
- respuesta inicial del GM persistida y publicada como `narrativa.respuesta_gm`

### Lo que todavía falta para considerar cerrado `Milestone 1`

- autenticación real
  - iniciar sesión
  - registrarse
  - jugar como invitado con `guest token`
  - asociación de partidas al usuario o invitado
- primer evento narrativo LLM real
  - hoy la llamada del Owner existe como estructura funcional
  - pero todavía no existe `narrative-service` generando ese mensaje vía LLM como define el canon
- entrada completa al juego según el flujo canon
  - pantalla de sesión
  - menú principal mínimo
  - carga de partidas asociadas a cuenta o guest token

### Conclusión importante

Si se deja de lado el tema visual por ahora, `Milestone 1` no está lejos, pero todavía no está cerrado.

La columna vertebral del inicio ya existe, pero faltan dos piezas grandes y concretas para poder decir que “El mundo nace” está completo de verdad:

- autenticación
- primer evento LLM real

---

## Próximo foco recomendado

Para la próxima sesión, el orden recomendado es:

### 1. Autenticación mínima completa

Resolver primero:

- landing de entrada con:
  - iniciar sesión
  - registrarse
  - jugar como invitado
- `guest token` real
- persistencia/asociación de partida a cuenta o invitado
- base mínima para `Cargar Partida`

La lógica de esta prioridad es simple:

- primero se define **quién es el jugador**
- después se define **qué partida está creando o retomando**

### 2. Primer evento LLM real

Después de auth, el siguiente paso natural es:

- levantar `narrative-service`
- mover la llamada del Owner a ese servicio
- generar el texto vía LLM según:
  - escenario inicial
  - modo de gestión
  - contexto de la partida

Esto cerraría el componente emocional real de `Milestone 1` según canon.

---

## Siguiente mini milestone sugerido

La próxima sesión debería enfocarse en:

**Autenticación mínima de `Milestone 1`**

Objetivo concreto:

- entrada al juego con cuenta o invitado
- token persistido
- partidas asociadas correctamente
- dejar listo el terreno para `Cargar Partida`

Una vez resuelto eso, el siguiente mini milestone ya debería ser:

**primer evento del Owner servido por `narrative-service` con LLM real**
