CREATE TABLE IF NOT EXISTS narrative_events (
    event_id TEXT PRIMARY KEY,
    game_id TEXT NOT NULL,
    kind TEXT NOT NULL,
    emitter TEXT NOT NULL,
    urgency TEXT NOT NULL,
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    metadata_json TEXT NOT NULL,
    choices_json TEXT NOT NULL,
    source_subject TEXT,
    source_event_id TEXT,
    source_match_id TEXT,
    simulated_date TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_narrative_events_game_match_kind
ON narrative_events (game_id, source_match_id, kind)
WHERE source_match_id IS NOT NULL;
