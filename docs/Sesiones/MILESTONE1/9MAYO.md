# Sesión — 2026-05-09

## Estado de la sesión

Continuación de `Milestone 1`, ya con una base técnica viva del día anterior.

El foco de hoy pasó a ser el frontend como experiencia:

- mantener el flujo real `gateway` + `map-service`
- elevar la ceremonia visual del nacimiento del mundo
- empezar a alinear el frontend con el nuevo `Design System`

---

## Avance registrado hoy

### 1. Ceremonia visual del nacimiento del mundo

Se avanzó sobre el `frontend/` para que deje de sentirse como una debug page bonita y empiece a parecerse a la experiencia real de `Milestone 1`.

### Qué cambió

- se rediseñó la pantalla principal como una ceremonia de fundación del mapa
- cada `stage` real emitido por backend (`terrain`, `zoning`, `stadium`, `complete`) ahora tiene identidad visual y copy propio
- el mapa quedó presentado dentro de una escena más atmosférica y más cercana al tono final de PulseCity
- se agregó timeline visible de etapas del backend para que el nacimiento del mundo se entienda como pipeline real
- se mejoró la lectura del estado local del cliente y la traza de eventos recientes
- se mantuvo intacto el contrato actual:
  - `map.snapshot`
  - `map.patch`

### Alcance de este mini milestone

No se agregaron nuevas dependencias ni se introdujo todavía `Three.js`.
La mejora se resolvió con:

- React
- TypeScript
- CSS

para empujar la experiencia sin abrir complejidad técnica innecesaria en esta sesión.

### Verificación

Se validó con:

```bash
cd frontend
npm run build
```

Resultado:

- build OK
- bundle generado correctamente por Vite

---

## Nuevo insumo incorporado hoy

Se agregó `docs/00_Start_Here/PulseCity_DesignSystem.md` como referencia formal para cualquier trabajo de frontend.

Impacto inmediato:

- ya no conviene seguir mejorando UI “por intuición”
- las próximas iteraciones visuales deberían alinearse con:
  - paleta oficial
  - tipografías `Bona Nova SC` + `DM Sans`
  - densidad media-alta
  - panels oscuros con bordes sutiles
  - animación funcional, no decorativa

### Conclusión de dirección

El siguiente mini milestone más lógico pasa a ser:

- flujo de creación de franquicia con identidad visual

pero ya implementado desde el lenguaje visual oficial de PulseCity, no como una landing suelta.

---

## Mini milestone adicional — Nueva Partida con identidad visual y estado inicial

Se reemplazó el foco del frontend: en lugar de arrancar desde una pantalla centrada en la ceremonia técnica, ahora la entrada principal es un flujo de `Nueva Partida` más cercano al canon del juego.

### Qué quedó implementado

- pantalla de `Nueva Partida` separada y más alineada al tono real de PulseCity
- bloque de `identidad visual de la franquicia` con:
  - nombre de ciudad
  - nombre de franquicia
  - abreviatura de 3 letras
  - color primario
  - color secundario
  - color de acento
- preview viva de identidad de franquicia usando la paleta elegida
- bloque de `estado inicial de la franquicia` con las cuatro opciones canon:
  - Reconstrucción
  - Ventana de contención
  - Histórica en declive
  - Expansión pura
- panel de confirmación que dispara la ceremonia real del mapa
- la ceremonia y el mapa siguen funcionando sobre el pipeline ya existente de `snapshot` / `patch`

### Decisión importante de alcance

Todavía no se cambiaron contratos backend para persistir:

- nombre de franquicia
- abreviatura
- colores
- escenario inicial

Por ahora esos datos viven del lado del frontend y alimentan la experiencia local de creación de partida. El `POST /api/v1/games` sigue usando únicamente `city_name`.

Esto fue deliberado para mantener el mini milestone acotado y no abrir todavía modelado de dominio adicional en `gateway` y persistencia.

### Alineación con el Design System

Se reorientó la UI para respetar mucho mejor el documento `PulseCity_DesignSystem.md`:

- fondos oscuros fijos por capas
- tipografías `Bona Nova SC` + `DM Sans`
- densidad media-alta
- radios chicos
- bordes sutiles sin depender de sombras pesadas
- verde PulseCity como acento principal de producto

### Verificación

Se validó con:

```bash
cd frontend
npm run build
```

Resultado:

- build OK
- bundle generado correctamente por Vite

### Qué habilita después

Este paso deja mejor preparada la siguiente evolución del frontend:

- persistir realmente la identidad de franquicia en backend
- sumar `modo de gestión de ciudad`
- avanzar hacia el primer evento narrativo obligatorio del Owner

---

## Mini milestone adicional — Persistencia real de identidad de franquicia

Se conectó la nueva pantalla de `Nueva Partida` con el dominio real del `gateway`, para que la partida ya no persista solamente `city_name`.

### Qué quedó implementado

En `services/gateway`:

- el `POST /api/v1/games` ahora acepta y persiste:
  - `city_name`
  - `franchise_name`
  - `abbreviation`
  - `primary_color`
  - `secondary_color`
  - `accent_color`
  - `initial_scenario`
- se agregaron defaults y normalización básica del lado del handler
- `map-service` no se tocó: sigue recibiendo únicamente lo que necesita para generar el mapa

En persistencia:

- se extendió el schema de `games`
- se agregó la migración:
  - `db/migrations/002_add_game_setup_fields.sql`
- `EnsureSchema` ahora contempla también estos nuevos campos para no romper entornos ya existentes

Además:

- se agregó `GET /api/v1/games/{gameID}` para leer la metadata completa de una partida

### Decisión de diseño importante

No se convirtió todavía la creación de partida en un onboarding multipaso real.

Se mantuvo el enfoque actual:

- una pantalla única como slice funcional
- persistencia real del contenido importante
- refactor a onboarding real más adelante, reutilizando estos bloques

Esto evita abrir ahora mismo navegación multipaso, validaciones intermedias y guardado parcial.

### Verificación

Se validó con:

```bash
cd services/gateway
go test ./...

cd frontend
npm run build
```

Resultado:

- `go test` OK en `gateway`
- frontend build OK

### Qué habilita después

Este paso deja el terreno listo para cualquiera de estos próximos saltos:

- sumar `modo de gestión de ciudad` al dominio persistido
- rehidratar también la metadata completa de la partida en frontend
- empezar a modelar la transición hacia el onboarding por pasos definitivo

---

## Mini milestone adicional — Rehidratación de metadata completa en frontend

Se cerró el loop de lectura/escritura de la nueva información de partida.

Hasta este punto, el frontend ya persistía:

- identidad de franquicia
- colores
- escenario inicial

pero al reabrir una partida solo rehidrataba el snapshot del mapa. Faltaba volver a traer también la metadata fundacional.

### Qué quedó implementado

En `frontend/`:

- se agregó el tipo `GameSetup`
- el frontend ahora consume `GET /api/v1/games/{gameID}`
- al rehidratar una partida por `game_id`, se restauran también:
  - `city_name`
  - `franchise_name`
  - `abbreviation`
  - `primary_color`
  - `secondary_color`
  - `accent_color`
  - `initial_scenario`

### Comportamiento nuevo

- si el usuario carga una partida existente, no vuelve solamente el mapa
- también vuelve la configuración fundacional de la franquicia en la UI
- eso deja la pantalla de `Nueva Partida` mucho más cerca de una futura experiencia de onboarding real, porque ya puede leer y reflejar estado persistido

### Verificación

Se validó con:

```bash
cd frontend
npm run build
```

Resultado:

- build OK

### Qué habilita después

El próximo salto natural ahora sí puede ser cualquiera de estos dos:

- agregar `modo de gestión de ciudad` a persistencia y frontend
- empezar a separar esta pantalla única en pasos reales de onboarding, pero ya sobre datos persistidos y rehidratables

---

## Mini milestone adicional — Persistencia y rehidratación de modo de gestión de ciudad

Se completó el último bloque fundacional importante que faltaba en la creación de partida actual: `modo de gestión de ciudad`.

### Qué quedó implementado

En backend:

- el `POST /api/v1/games` ahora acepta `city_management_mode`
- se agregó normalización del campo en el `gateway`
- `GameSetup` quedó extendido para incluirlo tanto en escritura como en lectura
- se agregó la migración:
  - `db/migrations/003_add_city_management_mode.sql`
- la tabla `games` ahora persiste este dato junto con el resto de la identidad fundacional

En frontend:

- se agregó un nuevo bloque visual para elegir entre:
  - `Dueño con influencia`
  - `Figura dual`
- la elección se envía al crear partida
- al rehidratar una partida por `game_id`, también vuelve este valor y se refleja en la UI
- el preview lateral ahora resume también el modo de gestión elegido

### Decisión de alcance

Todavía no hay comportamiento de gameplay atado a este modo.

Por ahora el objetivo fue:

- fijar el dato en el dominio
- persistirlo
- rehidratarlo
- dejarlo visible en el frontend

Eso es suficiente para preparar más adelante la lógica real de acceso a controles de ciudad, negociación con el alcalde o control dual.

### Verificación

Se validó con:

```bash
cd services/gateway
go test ./...

cd frontend
npm run build
```

Resultado:

- `go test` OK en `gateway`
- frontend build OK

### Qué habilita después

Con esto, la pantalla fundacional actual ya cubre casi todo el contenido central de la creación de partida temprana:

- identidad de franquicia
- estado inicial
- modo de gestión de ciudad

El siguiente salto lógico ya no parece ser “más campos”, sino uno de estos:

- empezar a partir esta pantalla en pasos reales de onboarding
- o conectar la siguiente consecuencia del flujo: primer evento narrativo obligatorio del Owner

---

## Mini milestone adicional — Primera llamada del Owner

Se implementó la primera consecuencia narrativa real del flujo fundacional: al completarse la generación del mapa, aparece la llamada inicial del Owner.

### Enfoque elegido

No se creó todavía `narrative-service` ni integración con LLM.

En cambio, se hizo un punto medio deliberado:

- estructura de evento narrativa ya preparada
- texto todavía hardcodeado
- contenido variable según `initial_scenario`

Esto permite avanzar en la experiencia sin abrir todavía toda la infraestructura narrativa final.

### Qué quedó implementado

En `gateway`:

- cuando el mapa llega a `stage = complete`, se genera un `NarrativeEvent`
- ese evento representa la `owner_intro`
- se persiste en base como parte de la partida
- se emite al frontend por WebSocket con una forma ya más cercana a la narrativa futura

Se agregó la migración:

- `db/migrations/004_add_owner_intro_event.sql`

El texto cambia según el `initial_scenario`:

- `rebuild`
- `contention`
- `decline`
- `expansion`

Además, el texto también hace una referencia ligera al `city_management_mode` elegido.

### En frontend

- se agregó soporte para `narrative.event`
- la llamada del Owner aparece como modal obligatorio sobre la UI
- si la partida ya tiene ese evento persistido, también se rehidrata con `GET /api/v1/games/{gameID}`
- la traza de eventos recientes ya puede mostrar también eventos narrativos

### Qué NO hace todavía

- no hay respuestas del GM persistidas
- no existe todavía `narrative-service`
- no hay LLM
- no hay cadena completa de eventos narrativos del sistema

Pero sí quedó resuelto el primer momento obligatorio del milestone con una forma técnica reutilizable.

### Verificación

Se validó con:

```bash
cd services/gateway
go test ./...

cd frontend
npm run build
```

Resultado:

- `go test` OK en `gateway`
- frontend build OK

### Qué habilita después

Con esto, `Milestone 1` ya está mucho más cerca de su experiencia objetivo:

- creación de franquicia
- ceremonia del mapa
- primera llamada del Owner

Los siguientes caminos naturales pasan a ser:

- persistir respuesta del GM a este evento
- o empezar a aislar este sistema como base del futuro `narrative-service`
