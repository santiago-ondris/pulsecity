# M2.1 — Scaffolding de servicios M2

## Objetivo de la sesion

Crear la estructura minima compilable para los servicios nuevos de Milestone 2.

## Mini milestone

M2.1 — Scaffolding de servicios M2.

## Archivos tocados

- `services/agent-service`
- `services/match-service`
- `services/team-service`
- `services/city-service`
- `services/analytics-service`
- `Makefile`
- `docs/Sesiones/MILESTONE2/INICIOM2.MD`

## Decision tomada

Los servicios nuevos arrancan con conexion real a NATS y estructura minima, pero sin logica de gameplay. Los contratos de eventos quedan para M2.2 para no mezclar scaffolding con definicion de payloads canonicos.

## Pruebas corridas

```bash
make test-rust
make test-go
make build
```

## Pendiente siguiente

M2.2 — definir contratos compartidos M2 para tiempo, partidos, ciudad, agentes y narrativa.
