# NOTAFINAL

## Milestone 1 — cierre practico

`Milestone 1` queda considerado **cerrado funcionalmente**.

El vertical slice fundacional ya permite:

- resolver identidad de entrada con:
  - cuenta
  - invitado
- crear una nueva franquicia
- recuperar partidas asociadas
- recorrer el onboarding inicial completo
- disparar la ceremonia de generacion del mapa en tiempo real via WebSocket
- recibir la primera llamada narrativa del Owner al finalizar la fundacion

## Estado tecnico alcanzado

Durante `Milestone 1` quedaron resueltos estos bloques:

- `gateway` como punto de entrada real
- `map-service` publicando progreso de generacion
- frontend aplicando `map.snapshot` y `map.patch`
- flujo visual inicial rehcho con direccion consistente
- `SessionPage` separada de `LandingPage`
- `guest auth` y `user auth` funcionales
- `narrative-service` minimo creado
- `owner_intro` movido fuera de `gateway` para quedar bajo ownership narrativo correcto

## Excepcion consciente para el cierre

El unico punto no completado en sentido estricto del canon es este:

- la primera llamada del Owner **todavia no usa un proveedor LLM real**

Hoy esa narrativa:

- nace en `narrative-service`
- se persiste correctamente
- se publica correctamente al frontend
- pero el texto sigue siendo templateado

La integracion con OpenAI, Anthropic o Azure queda **diferida por costo**, no por limitacion tecnica ni por falta de arquitectura.

## Decision tomada

Se considera que `Milestone 1` queda cerrado a nivel de:

- experiencia jugable
- arquitectura base
- flujo fundacional completo

Y se deja explicitamente asentado que el reemplazo del texto templateado por un LLM real queda para una etapa posterior.

## Lo que habilita ahora

Con esto, PulseCity puede pasar a `Milestone 2` con una base ya ordenada para trabajar:

- tiempo simulado
- loop central
- partidos
- reaccion de ciudad
- agentes core

En otras palabras: el mundo ya nace. Ahora tiene que empezar a latir.
