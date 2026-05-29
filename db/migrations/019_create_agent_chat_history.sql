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
