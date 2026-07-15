# SESION21MAYO

## Mini milestone 2

Mover la primera llamada del Owner fuera de `gateway` y levantar un `narrative-service` minimo para acercar el cierre canonico de `Milestone 1`.

## Objetivo

- dejar de generar narrativa dentro de `gateway`
- hacer que `gateway` publique una solicitud de narrativa
- hacer que `narrative-service` genere la `owner_intro`
- mantener persistencia, rehidratacion y envio por WebSocket sin romper el flujo actual

## Notas de alcance

- este corte solo cubre un tipo de evento narrativo: `owner_intro`
- no abre todavia bandeja completa, multiples agentes ni cadenas narrativas
- el foco es corregir ownership y flujo entre servicios

## Cambios realizados

- Se agrego `services/narrative-service` como servicio Go minimo con:
  - conexion a NATS
  - conexion a PostgreSQL
  - suscripcion a `narrativa.owner_intro_solicitada`
  - delay intencional de `250-500ms` antes de generar la narrativa
- Se movio la construccion del `owner_intro` fuera de `gateway`.
- `gateway` ahora:
  - sigue escuchando `mapa.*`
  - cuando el mapa llega a `complete`, publica `narrativa.owner_intro_solicitada`
  - escucha `narrativa.evento_generado`
  - reenvia el `NarrativeEvent` al frontend por WebSocket
- `narrative-service` ahora:
  - recibe la solicitud de intro del Owner
  - genera el `NarrativeEvent`
  - persiste `owner_intro_event` en `games` solo si todavia no existe
  - publica `narrativa.evento_generado`
- Se agregaron tipos de dominio nuevos en `gateway` para formalizar el contrato de solicitud narrativa.
- Se elimino de `gateway` la responsabilidad de redactar narrativa.

## Estado real del corte

- El ownership narrativo ya quedo mejor alineado con el canon.
- La primera narrativa fundacional ya no nace en `gateway`.
- El texto del Owner sigue siendo **templateado**, pero ahora vive en el lugar correcto: `narrative-service`.

## Verificacion

- `go test ./...` en `services/gateway` OK
- `GOCACHE=/tmp/pulsecity-narrative-gocache go test ./...` en `services/narrative-service` OK

## Decision sobre LLM real

Por ahora se decide **no integrar un proveedor LLM real** en `Milestone 1`.

### Motivo

- la complejidad tecnica del cableado no es alta
- pero hoy el costo no se justifica para un entorno de desarrollo
- pagar billing real solo para probar la intro del Owner no aporta suficiente valor en esta etapa

### Conclusion

- `narrative-service` queda ya en el lugar correcto de arquitectura
- la narrativa fundacional se sigue resolviendo con texto templateado
- la integracion con proveedor real queda diferida por costo, no por limitacion tecnica

## Proveedores considerados a futuro

### Opcion 1 — OpenAI directo

- camino simple para conectar rapido un modelo por API
- buena opcion cuando se quiera baja friccion operativa

### Opcion 2 — Anthropic directo

- alternativa similar en complejidad a OpenAI directo
- valida si mas adelante se prioriza el tono de Claude para narrativa

### Opcion 3 — Azure / Azure AI Foundry

- opcion explicitamente valorada para PulseCity
- experiencia previa positiva del usuario con despliegues baratos y suficientes para texto simple
- mas friccion operativa inicial que OpenAI o Anthropic directos
- puede ser una muy buena opcion cuando se quiera:
  - centralizar infraestructura
  - controlar deployments por entorno
  - aprovechar un modelo economico desplegado solo para narrativa ligera

## Decision abierta para mas adelante

Cuando llegue el momento de conectar un LLM real, Azure queda registrado como opcion totalmente valida y deseable para evaluar primero, junto con OpenAI directo.
