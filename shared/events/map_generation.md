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
  "city_name": "string opcional"
}
```

## Payload de progreso

```json
{
  "game_id": "uuid",
  "stage": "terrain|zoning|stadium|complete",
  "progress": 25,
  "message": "descripcion legible"
}
```

## Delta WebSocket emitido por el gateway

```json
{
  "type": "map.delta",
  "subject": "mapa.terreno_listo",
  "payload": {
    "game_id": "uuid",
    "stage": "terrain",
    "progress": 25,
    "message": "Terreno base generado"
  }
}
```
