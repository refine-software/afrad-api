-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'sessions_user_id_user_agent_key'
    ) THEN
        ALTER TABLE sessions
        ADD CONSTRAINT sessions_user_id_user_agent_key UNIQUE (user_agent, user_id);
    END IF;
END
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE sessions DROP CONSTRAINT IF EXISTS sessions_user_id_user_agent_key;
-- +goose StatementEnd
