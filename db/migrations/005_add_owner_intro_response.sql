ALTER TABLE games
ADD COLUMN IF NOT EXISTS owner_intro_response JSONB;
