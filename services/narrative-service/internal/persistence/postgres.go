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

CREATE TABLE IF NOT EXISTS agent_chat_history (
	message_id TEXT PRIMARY KEY,
	game_id TEXT NOT NULL,
	conversation_id TEXT NOT NULL,
	agent_id TEXT NOT NULL,
	sender TEXT NOT NULL,
	body TEXT NOT NULL,
	metadata_json TEXT NOT NULL,
	context_json TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_agent_chat_history_conversation
ON agent_chat_history (game_id, agent_id, conversation_id, created_at);
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

func (s *Store) LoadAgentChatContext(ctx context.Context, gameID, agentID string) (domain.AgentChatContext, error) {
	context, err := s.loadIndividualAgentContext(ctx, gameID, agentID)
	if err != nil {
		return domain.AgentChatContext{}, err
	}
	if context.AgentID == "" {
		context, err = s.loadPlayerAgentContext(ctx, gameID, agentID)
		if err != nil {
			return domain.AgentChatContext{}, err
		}
	}
	if context.AgentID == "" {
		context.AgentID = agentID
		context.DisplayName = agentID
		context.Role = "agente"
		context.Domain = "estado general de la franquicia"
		context.EmotionalState = "unknown"
	}

	relationship, err := s.loadRelationshipWithGM(ctx, gameID, agentID)
	if err != nil {
		return domain.AgentChatContext{}, err
	}
	context.Relationship = relationship

	decisions, err := s.loadLatestGMDecisions(ctx, gameID, 5)
	if err != nil {
		return domain.AgentChatContext{}, err
	}
	context.Decisions = decisions

	return context, nil
}

func (s *Store) InsertChatMessageIfNew(ctx context.Context, message domain.ChatMessage, chatContext domain.AgentChatContext) (bool, error) {
	metadata, err := json.Marshal(message.Metadata)
	if err != nil {
		return false, fmt.Errorf("marshal chat metadata: %w", err)
	}
	contextPayload, err := json.Marshal(chatContext)
	if err != nil {
		return false, fmt.Errorf("marshal chat context: %w", err)
	}

	const query = `
INSERT INTO agent_chat_history (
	message_id,
	game_id,
	conversation_id,
	agent_id,
	sender,
	body,
	metadata_json,
	context_json,
	created_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (message_id) DO NOTHING;
`

	commandTag, err := s.pool.Exec(ctx, query,
		message.MessageID,
		message.GameID,
		message.ConversationID,
		message.AgentID,
		message.Sender,
		message.Body,
		string(metadata),
		string(contextPayload),
		parseTimeOrNow(message.CreatedAt),
	)
	if err != nil {
		return false, fmt.Errorf("insert chat message: %w", err)
	}

	return commandTag.RowsAffected() > 0, nil
}

func (s *Store) loadIndividualAgentContext(ctx context.Context, gameID, agentID string) (domain.AgentChatContext, error) {
	var chatContext domain.AgentChatContext
	err := s.pool.QueryRow(ctx, `
SELECT agent_id, display_name, role, domain, emotional_state, confidence, satisfaction, loyalty
FROM agent_individual_states
WHERE game_id = $1 AND agent_id = $2;
`, gameID, agentID).Scan(
		&chatContext.AgentID,
		&chatContext.DisplayName,
		&chatContext.Role,
		&chatContext.Domain,
		&chatContext.EmotionalState,
		&chatContext.Confidence,
		&chatContext.Satisfaction,
		&chatContext.Loyalty,
	)
	if err == pgx.ErrNoRows {
		return domain.AgentChatContext{}, nil
	}
	if err != nil {
		return domain.AgentChatContext{}, fmt.Errorf("load individual agent context: %w", err)
	}

	return chatContext, nil
}

func (s *Store) loadPlayerAgentContext(ctx context.Context, gameID, playerID string) (domain.AgentChatContext, error) {
	var chatContext domain.AgentChatContext
	err := s.pool.QueryRow(ctx, `
SELECT player_id, full_name, position, emotional_state, satisfaction, loyalty
FROM agent_player_states
WHERE game_id = $1 AND player_id = $2;
`, gameID, playerID).Scan(
		&chatContext.AgentID,
		&chatContext.DisplayName,
		&chatContext.Role,
		&chatContext.EmotionalState,
		&chatContext.Satisfaction,
		&chatContext.Loyalty,
	)
	if err == pgx.ErrNoRows {
		return domain.AgentChatContext{}, nil
	}
	if err != nil {
		return domain.AgentChatContext{}, fmt.Errorf("load player agent context: %w", err)
	}

	chatContext.Domain = "rendimiento, rol, vestuario y experiencia personal dentro del roster"
	return chatContext, nil
}

func (s *Store) loadRelationshipWithGM(ctx context.Context, gameID, agentID string) (domain.AgentRelationshipContext, error) {
	var relationship domain.AgentRelationshipContext
	err := s.pool.QueryRow(ctx, `
SELECT trust, trend, last_event
FROM agent_relationships
WHERE game_id = $1
	AND (
		(agent_a_id = $2 AND agent_b_id = 'gm')
		OR (agent_a_id = 'gm' AND agent_b_id = $2)
	)
ORDER BY updated_at DESC
LIMIT 1;
`, gameID, agentID).Scan(
		&relationship.Trust,
		&relationship.Trend,
		&relationship.LastEvent,
	)
	if err == pgx.ErrNoRows {
		return domain.AgentRelationshipContext{}, nil
	}
	if err != nil {
		return domain.AgentRelationshipContext{}, fmt.Errorf("load relationship with gm: %w", err)
	}

	return relationship, nil
}

func (s *Store) loadLatestGMDecisions(ctx context.Context, gameID string, limit int) ([]domain.GMDecisionContext, error) {
	rows, err := s.pool.Query(ctx, `
SELECT decision_id, kind, simulated_date
FROM gm_decisions_log
WHERE game_id = $1
ORDER BY created_at DESC
LIMIT $2;
`, gameID, limit)
	if err != nil {
		return nil, fmt.Errorf("load latest gm decisions: %w", err)
	}
	defer rows.Close()

	decisions := make([]domain.GMDecisionContext, 0)
	for rows.Next() {
		var decision domain.GMDecisionContext
		if err := rows.Scan(&decision.DecisionID, &decision.Kind, &decision.SimulatedDate); err != nil {
			return nil, err
		}
		decisions = append(decisions, decision)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return decisions, nil
}

func parseTimeOrNow(value string) time.Time {
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err == nil {
		return parsed.UTC()
	}

	parsed, err = time.Parse(time.RFC3339, value)
	if err == nil {
		return parsed.UTC()
	}

	return time.Now().UTC()
}
