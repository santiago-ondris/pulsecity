# Sesion 2026-05-20

## Mini milestone

Rehacer desde cero la landing page del frontend para que deje de sentirse como un panel tecnico mezclado con login y biblioteca, y pase a ser una entrada visual fuerte, centrada en decision de ingreso.

## Cambios realizados

- Se reemplazo la `LandingPage` anterior por una composicion cinematografica fullscreen.
- Se separo visualmente la entrada principal de dos piezas secundarias:
  - `Cuenta / acceso`
  - `Biblioteca / cargar partida`
- Se mantuvo la logica existente para:
  - iniciar nueva partida cuando ya existe identidad activa
  - crear sesion invitada
  - continuar una partida guardada si existe seleccion
  - login y registro desde el panel de acceso
- Se agrego un fondo ilustrado local en SVG con skyline nocturno y estadio central para salir del fondo negro plano.
- Se incorporo `src/vite-env.d.ts` para permitir imports tipados de assets `.svg`.
- Se descarto la primera reimplementacion de `IdentityPage` porque seguia arrastrando la logica de `form + preview`.
- Se reemplazo `IdentityPage` por una version nueva basada en UX directa:
  - mismo fondo visual de la landing
  - topbar compacto
  - preview contenida
  - editor principal claro para ciudad, franquicia, sigla y paleta
  - CTA visible dentro del flujo de edicion
- Se rehizo `ScenarioPage` desde cero con la misma direccion:
  - fondo compartido con landing/identity
  - resumen claro del escenario seleccionado
  - lista central de opciones grandes
  - seleccion visual evidente
  - acciones `Volver` y `Continuar` dentro del flujo principal
- Se rehizo `ManagementPage` desde cero:
  - decision presentada como modelo de poder
  - resumen activo del modo elegido
  - dos opciones grandes y comparables
  - misma estructura visual full-screen que identity/scenario
- Se rehizo `LaunchPage` desde cero:
  - pantalla final de confirmacion, no dossier tecnico
  - identidad de franquicia resumida
  - lectura compacta de escenario, gobierno, owner y estado
  - CTA principal dominante para fundar ciudad y generar mapa
- Se rehizo `CeremonyPage` desde cero:
  - mapa como pieza principal de la pantalla
  - topbar compacto con etapa y progreso
  - pipeline, estado y eventos en columna operativa
  - se mantuvo el contrato existente de `map.snapshot` / `map.patch`
- Se actualizo `docs/00_Start_Here/PulseCity_DesignSystem.md` a v1.1 para dejar como canon la direccion visual nueva:
  - pantallas de flujo full-screen
  - anti-dashboard
  - fondos visuales con overlay
  - una decision dominante por pantalla
  - checklist de aceptacion visual

## Decisiones UX tomadas

- La landing principal queda reducida a tres decisiones:
  - `Nueva partida`
  - `Cargar partida`
  - `Jugar como invitado`
- `Cargar partida` se trata como experiencia separada de la landing emocional.
- `Cuenta` tambien queda separada del momento de entrada, aunque por ahora siga embebida como panel lateral hasta construir su pagina dedicada.

## Placeholder consciente

- `Biblioteca` todavia no es una pagina propia: por ahora vive en el panel lateral nuevo.
- `Cuenta / login` todavia no es una pagina independiente: por ahora vive en el panel lateral nuevo.
- El fondo visual nuevo es un asset puente para fijar direccion estetica; despues puede ser reemplazado por arte final generado externamente.
- El flujo inicial completo (`IdentityPage`, `ScenarioPage`, `ManagementPage`, `LaunchPage`, `CeremonyPage`) ya comparte el nuevo lenguaje visual.

## Verificacion

- `npm run build` en `frontend` compilo correctamente.

---

## Mini milestone adicional

Separar la entrada de cuenta de la landing principal para que identidad y eleccion de partida dejen de convivir en la misma pantalla.

## Cambios realizados

- Se agrego una `SessionPage` dedicada como primera pantalla cuando no existe sesion activa.
- La `SessionPage` mantiene la misma direccion visual cinematografica del onboarding ya rehcho:
  - fondo compartido
  - topbar compacta
  - bloque editorial de contexto
  - panel de acceso dedicado
- El acceso de cuenta ahora vive exclusivamente en esa pantalla:
  - `Login`
  - `Registro`
  - `Olvide mi contraseña` como placeholder visible de producto
  - `Jugar como invitado`
- La `LandingPage` quedo reducida al problema correcto:
  - `Nueva partida`
  - `Cargar partida`
- Se elimino de la landing la accion de `Jugar como invitado`, porque pertenece a identidad y no a biblioteca o arranque de partida.
- Se ajusto el routing del flujo para que:
  - sin sesion activa se fuerce `SessionPage`
  - con sesion activa se redirija a la landing
  - la ruta `/session` ya no quede accesible una vez autenticado o con invitado restaurado
- Se mantuvo la logica existente de:
  - restaurar sesion de usuario
  - restaurar sesion invitada
  - migrar guest -> user
  - logout y vuelta al invitado guardado si existe

## Decision UX tomada

- La identidad del jugador pasa a resolverse antes de la landing.
- La landing ya no mezcla decisiones de acceso con decisiones de juego.
- `Cuenta` y `Landing` quedan como experiencias separadas con responsabilidades distintas.

## Verificacion adicional

- `npm run build` en `frontend` compilo correctamente despues de separar `SessionPage` y `LandingPage`.

## Deuda consciente de auth

Se valida como suficiente para `Milestone 1`:

- registro funcional
- login funcional
- logout funcional
- acceso invitado funcional
- restauracion de sesion y biblioteca asociada

Queda explicitamente diferido para mas adelante, fuera del cierre de `Milestone 1`:

- ampliar los datos requeridos en registro
- verificacion de email
- politica de contraseña mas estricta
- `Olvide mi contraseña` funcional

### Criterio tomado

Estas mejoras endurecen el sistema de cuenta, pero hoy no destraban el objetivo fundacional del milestone.

Se prioriza mantener baja friccion de entrada mientras se completa el cierre canonico de `Milestone 1`:

- creacion o recuperacion de identidad
- fundacion de la franquicia
- ceremonia de generacion del mapa
- primer evento narrativo real del Owner servido por `narrative-service`
