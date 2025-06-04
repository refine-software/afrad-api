-- +goose Up
-- +goose StatementBegin

ALTER TABLE sessions ADD constraint sessions_user_id_user_agent_key UNIQUE (
    user_agent, user_id
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE sessions DROP CONSTRAINT IF EXISTS sessions_user_id_user_agent_key;
-- +goose StatementEnd
