DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'games_exactly_one_owner'
    ) THEN
        ALTER TABLE games
            ADD CONSTRAINT games_exactly_one_owner
            CHECK (
                ((CASE WHEN guest_token <> '' THEN 1 ELSE 0 END) +
                (CASE WHEN user_id IS NOT NULL AND user_id <> '' THEN 1 ELSE 0 END)) = 1
            ) NOT VALID;
    END IF;
END $$;
