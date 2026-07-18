# M3.P2f-a — Sesiones guest y WebSocket del gateway

## Objetivo de la sesion

Iniciar la particion de `gateway/internal/handlers/http.go` con una frontera por recurso, sin cambiar rutas ni comportamiento.

## Analisis previo

Los handlers de registro, login, sesion actual y upgrade guest ya vivian en `auth.go`. Por eso se descarto crear otro archivo de autenticacion solo para `createGuestSession`.

El corte elegido agrupa dos responsabilidades relacionadas:

- `auth.go`: creacion y lectura de identidad guest
- `websocket.go`: validacion de la identidad de conexion y ciclo de sesion activa de simulacion

## Cambios realizados

- Se creo `services/gateway/internal/handlers/websocket.go`.
- Se movieron desde `http.go`:
  - `serveWebSocket`
  - `clientIDFromRequest`
  - `guestOwnsGame`
- Se movieron a `auth.go`:
  - `createGuestSession`
  - `guestTokenFromRequest`
- Los metodos de `auth.go` quedaron antes que sus funciones de validacion y utilidades, siguiendo orden idiomatico de Go.
- `RegisterRoutes` no cambio; los handlers siguen disponibles dentro del mismo paquete.

## Decision tomada

La particion se realiza por archivos dentro de `package handlers`, no mediante subpaquetes. Los archivos comparten `Dependencies`, autorizacion y helpers internos; separarlos en paquetes agregaria dependencias y API sin aportar aislamiento real.

## Resultado

- `http.go`: 1641 -> 1492 lineas.
- `websocket.go`: 135 lineas.
- `auth.go`: 395 lineas, incluyendo ahora el flujo guest completo.
- Sin cambios en rutas HTTP, eventos NATS, payloads, ownership o comportamiento.
- Suite, build y vet del gateway pasan.

## Verificacion

```bash
gofmt -w services/gateway/internal/handlers/http.go \
  services/gateway/internal/handlers/auth.go \
  services/gateway/internal/handlers/websocket.go
GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...
GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway vet ./...
GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway build ./...
git diff --check
```

## Pendiente siguiente

`M3.P2` y la particion del gateway continuan pendientes. El siguiente corte recomendado es separar los handlers de games/snapshot de `http.go` antes de abordar chat, medicina y trades.
