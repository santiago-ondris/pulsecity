# PulseCity

## Entorno local jugable

Para levantar infraestructura y servicios de app en un solo comando:

```bash
make dev-app
```

Puertos principales:

- frontend: `http://localhost:5173`
- gateway: `http://localhost:8080`
- Postgres/TimescaleDB: `localhost:5433`
- NATS: `localhost:4222`

Flujo recomendado:

1. Ejecutar `make dev-app`.
2. Abrir `http://localhost:5173`.
3. Crear o cargar una partida.
4. Entrar a la pantalla viva de M2 y usar pausa/x1/x5/x20 desde el HUD.

`agent-service` ya no necesita `GAME_ID` para el flujo normal de desarrollo:
cuando el gateway publica `tiempo.sesion_iniciada`, el servicio carga o crea el
estado de simulacion de esa partida y empieza a escuchar sus controles de tiempo.

Para depurar una partida concreta, el modo antiguo sigue disponible:

```bash
GAME_ID=<id-de-la-partida> make run-agent-service
```
