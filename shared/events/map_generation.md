# Map Generation Events

Contrato mínimo del vertical slice entre `gateway` y `map-service`.

## Subjects

- `mapa.generacion_iniciada`
- `mapa.terreno_listo`
- `mapa.zonas_calculadas`
- `mapa.estadio_ubicado`
- `mapa.generacion_completa`

## Payload de inicio

```json
{
  "game_id": "uuid",
  "city_name": "string opcional",
  "franchise_name": "string opcional",
  "abbreviation": "string opcional"
}
```

Notas:

- `map-service` solo necesita `game_id` y `city_name`.
- `team-service` consume el mismo evento para inicializar franquicia, roster y temporada. Por eso `franchise_name` y `abbreviation` viajan en el payload aunque no afecten la generacion del mapa.

## Payload de progreso

```json
{
  "game_id": "uuid",
  "stage": "terrain|zoning|stadium|complete",
  "progress": 25,
  "message": "descripcion legible",
  "map_data": {
    "width": 20,
    "height": 20,
    "cells": [
      [
        { "terrain": "water" },
        { "terrain": "plain", "zone": "residential" }
      ]
    ]
  },
  "stadium": { "x": 10, "y": 10 }
}
```

Notas:

- `map_data` aparece cuando hay datos de grilla disponibles
- `stadium` aparece desde la etapa `stadium`
- el mapa actual usa una generacion procedural-lite deterministica basada en el nombre de ciudad
- el terreno ya usa `Perlin 2D` real para elevacion base
- la zonificacion ya usa particion Voronoi simple con seeds deterministicas por ciudad

## Contrato WebSocket emitido por el gateway

El `gateway` ya no reenvia el payload crudo de NATS. Lo traduce a un contrato pensado para frontend:

- `map.snapshot` → estado inicial completo
- `map.patch` → cambios parciales posteriores

### Snapshot inicial

```json
{
  "type": "map.snapshot",
  "subject": "mapa.terreno_listo",
  "state": {
    "game_id": "uuid",
    "stage": "terrain",
    "progress": 25,
    "message": "Terreno base generado",
    "map_data": {
      "width": 20,
      "height": 20,
      "cells": []
    }
  }
}
```

### Patch incremental

```json
{
  "type": "map.patch",
  "subject": "mapa.estadio_ubicado",
  "game_id": "uuid",
  "patch": {
    "stage": "stadium",
    "progress": 80,
    "message": "Estadio ubicado en el distrito central",
    "stadium": { "x": 10, "y": 10 }
  }
}
```

### Regla actual del slice

- el primer evento con `map_data` se convierte en `map.snapshot`
- los eventos siguientes se convierten en `map.patch`
- hoy los patches pueden incluir:
  - `stage`
  - `progress`
  - `message`
  - `map_data`
  - `stadium`

La idea es que el frontend mantenga un estado local y aplique patches sobre el snapshot inicial.

## Rehidratacion actual

Para resolver clientes que se conectan tarde, el `gateway` guarda en memoria el ultimo estado conocido por `game_id`.

Eso permite dos caminos:

### Snapshot por HTTP

```text
GET /api/v1/games/{game_id}/snapshot
```

Respuesta:

```json
{
  "type": "map.snapshot",
  "subject": "gateway.snapshot_http",
  "state": {
    "game_id": "uuid",
    "stage": "complete",
    "progress": 100,
    "message": "Generacion de mapa completada"
  }
}
```

### Snapshot al conectar WebSocket

Si el cliente abre:

```text
/ws?game_id={game_id}
```

y el `gateway` ya tiene snapshot en memoria, envia inmediatamente:

```json
{
  "type": "map.snapshot",
  "subject": "gateway.snapshot_rehidratado",
  "state": {
    "game_id": "uuid"
  }
}
```

## Limitacion actual

- el `gateway` usa cache en memoria para acceso rapido
- si el snapshot no esta en memoria, intenta rehidratarlo desde PostgreSQL
- la persistencia actual guarda el ultimo snapshot completo por `game_id`
- no hay todavia historial de versiones ni eventos reproducibles desde base
