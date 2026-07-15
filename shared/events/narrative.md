# Narrative Events

Contratos narrativos para Milestone 2. El dueño de la bandeja narrativa es `narrative-service`.

## Subjects

- `narrativa.evento_generado`
- `narrativa.respuesta_gm`

M1 ya usa `narrativa.owner_intro_solicitada` como evento interno de arranque. Se mantiene.

## `narrativa.evento_generado`

Publicado por `narrative-service`.

```json
{
  "event_id": "uuid",
  "game_id": "uuid",
  "type": "narrative.event",
  "subject": "narrativa.evento_generado",
  "emitter": "owner|head_coach|cfo|scouting_director|sports_psychologist|press",
  "kind": "owner_intro|post_match|agent_warning|city_reaction",
  "urgency": "low|medium|high|critical",
  "title": "Victoria que calma el vestuario",
  "body": "El coach cree que el equipo encontró una base competitiva.",
  "metadata": {
    "match_id": "uuid",
    "source_event_id": "uuid",
    "source_subject": "partido.terminado",
    "simulated_date": "2026-10-22",
    "home_team_id": "pulsecity",
    "away_team_id": "rival-generated-id",
    "home_score": "112",
    "away_score": "106",
    "winner_team_id": "pulsecity",
    "margin": "6",
    "win_streak": "3"
  },
  "choices": [
    {
      "id": "acknowledge",
      "label": "Tomar nota"
    }
  ]
}
```

## `narrativa.respuesta_gm`

Publicado por `gateway` cuando el jugador responde una opcion narrativa.

```json
{
  "event_id": "uuid",
  "game_id": "uuid",
  "occurred_at": "2026-05-21T00:00:10Z",
  "schema_version": 1,
  "narrative_event_id": "uuid",
  "choice_id": "acknowledge",
  "metadata": {}
}
```

Notas:

- en M2 la narrativa es templateada/rule-based.
- para post-partido, `narrative-service` escucha `partido.terminado`, espera 250-500ms y genera `narrativa.evento_generado`.
- para post-trade, `narrative-service` escucha `trade.aceptada`, espera 250-500ms y genera `narrativa.evento_generado` con `kind = post_trade`.

## `post_trade`

Evento narrativo templateado generado cuando se cierra un trade.

```json
{
  "event_id": "post-trade-trade-uuid",
  "game_id": "game-1",
  "type": "narrative.event",
  "subject": "narrativa.evento_generado",
  "emitter": "director_player_personnel",
  "kind": "post_trade",
  "urgency": "normal",
  "title": "Trade cerrado",
  "body": "Player Personnel confirma el cierre: Adrian Vale sale de PulseCity y Jalen Warren llega para cubrir PG...",
  "metadata": {
    "proposal_id": "trade-uuid",
    "source_event_id": "trade-accepted-trade-uuid",
    "source_subject": "trade.aceptada",
    "simulated_date": "2026-11-01",
    "rival_team_id": "bos",
    "outgoing_player_id": "game-1-player-06",
    "incoming_player_id": "trade-uuid-incoming",
    "incoming_position": "PG",
    "incoming_rating": "76",
    "incoming_salary": "12000000"
  },
  "choices": [
    { "id": "acknowledge", "label": "Tomar nota" }
  ]
}
```
- despues del delay puede leer contexto actualizado de la partida para incluir racha si esta disponible.
- el `gateway` reenvia este evento al frontend como `narrative.event`.
