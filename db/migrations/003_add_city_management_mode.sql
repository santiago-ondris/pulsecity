ALTER TABLE games
ADD COLUMN IF NOT EXISTS city_management_mode TEXT NOT NULL DEFAULT 'owner_influence';
