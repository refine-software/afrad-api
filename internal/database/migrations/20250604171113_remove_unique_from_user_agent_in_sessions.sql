-- +goose Up
-- +goose StatementBegin
ALTER TABLE sessions DROP CONSTRAINT IF EXISTS sessions_user_agent_key;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE sessions ADD CONSTRAINT sessions_user_agent_key UNIQUE (user_agent);
-- +goose StatementEnd
