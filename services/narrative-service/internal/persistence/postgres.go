package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pulsecity/services/narrative-service/internal/domain"
)

type Store struct {
	pool *pgxpool.Pool
}

func NewStore(ctx context.Context, databaseURL string) (*Store, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	return &Store{pool: pool}, nil
}

func (s *Store) Close() {
	s.pool.Close()
}

func (s *Store) EnsureSchema(ctx context.Context) error {
	const query = `
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
`

	_, err := s.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("ensure narrative schema: %w", err)
	}

	return nil
}

func (s *Store) SetOwnerIntroEventIfEmpty(ctx context.Context, gameID string, event domain.NarrativeEvent) (bool, error) {
	payload, err := json.Marshal(event)
	if err != nil {
		return false, fmt.Errorf("marshal owner intro event: %w", err)
	}

	const query = `
UPDATE games
SET owner_intro_event = $2::jsonb,
	updated_at = $3
WHERE game_id = $1
	AND owner_intro_event IS NULL;
`

	commandTag, err := s.pool.Exec(ctx, query, gameID, payload, time.Now().UTC())
	if err != nil {
		return false, err
	}

	return commandTag.RowsAffected() > 0, nil
}

func (s *Store) InsertNarrativeEventIfNew(ctx context.Context, event domain.NarrativeEvent) (bool, error) {
	metadata, err := json.Marshal(event.Metadata)
	if err != nil {
		return false, fmt.Errorf("marshal narrative metadata: %w", err)
	}
	choices, err := json.Marshal(event.Choices)
	if err != nil {
		return false, fmt.Errorf("marshal narrative choices: %w", err)
	}

	sourceSubject := event.Metadata["source_subject"]
	sourceEventID := event.Metadata["source_event_id"]
	sourceMatchID := event.Metadata["match_id"]
	simulatedDate := event.Metadata["simulated_date"]

	const query = `
INSERT INTO narrative_events (
	event_id,
	game_id,
	kind,
	emitter,
	urgency,
	title,
	body,
	metadata_json,
	choices_json,
	source_subject,
	source_event_id,
	source_match_id,
	simulated_date,
	created_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NULLIF($10, ''), NULLIF($11, ''), NULLIF($12, ''), NULLIF($13, ''), NOW())
ON CONFLICT (game_id, source_match_id, kind)
WHERE source_match_id IS NOT NULL
DO NOTHING;
`

	commandTag, err := s.pool.Exec(ctx, query,
		event.EventID,
		event.GameID,
		event.Kind,
		event.Emitter,
		event.Urgency,
		event.Title,
		event.Body,
		string(metadata),
		string(choices),
		sourceSubject,
		sourceEventID,
		sourceMatchID,
		simulatedDate,
	)
	if err != nil {
		return false, fmt.Errorf("insert narrative event: %w", err)
	}

	return commandTag.RowsAffected() > 0, nil
}

func (s *Store) LoadNarrativeEventByMatch(ctx context.Context, gameID, matchID, kind string) (domain.NarrativeEvent, bool, error) {
	var event domain.NarrativeEvent
	var metadataJSON string
	var choicesJSON string

	err := s.pool.QueryRow(ctx, `
SELECT event_id, game_id, emitter, kind, urgency, title, body, metadata_json, choices_json
FROM narrative_events
WHERE game_id = $1 AND source_match_id = $2 AND kind = $3;
`, gameID, matchID, kind).Scan(
		&event.EventID,
		&event.GameID,
		&event.Emitter,
		&event.Kind,
		&event.Urgency,
		&event.Title,
		&event.Body,
		&metadataJSON,
		&choicesJSON,
	)
	if err == pgx.ErrNoRows {
		return domain.NarrativeEvent{}, false, nil
	}
	if err != nil {
		return domain.NarrativeEvent{}, false, fmt.Errorf("load narrative event by match: %w", err)
	}

	event.Type = "narrative.event"
	event.Subject = domain.SubjectNarrativeEventGenerated
	if err := json.Unmarshal([]byte(metadataJSON), &event.Metadata); err != nil {
		return domain.NarrativeEvent{}, false, fmt.Errorf("decode narrative metadata: %w", err)
	}
	if err := json.Unmarshal([]byte(choicesJSON), &event.Choices); err != nil {
		return domain.NarrativeEvent{}, false, fmt.Errorf("decode narrative choices: %w", err)
	}

	return event, true, nil
}

func (s *Store) LoadNarrativeContext(ctx context.Context, gameID string) (domain.NarrativeContext, error) {
	var narrativeContext domain.NarrativeContext
	err := s.pool.QueryRow(ctx, `
SELECT win_streak, loss_streak
FROM city_metrics
WHERE game_id = $1;
`, gameID).Scan(&narrativeContext.WinStreak, &narrativeContext.LossStreak)
	if err == pgx.ErrNoRows {
		return domain.NarrativeContext{}, nil
	}
	if err != nil {
		return domain.NarrativeContext{}, fmt.Errorf("load narrative context: %w", err)
	}

	return narrativeContext, nil
}
