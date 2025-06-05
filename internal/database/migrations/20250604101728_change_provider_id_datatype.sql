-- +goose Up
-- +goose StatementBegin
ALTER TABLE oauth
ALTER COLUMN provider_id TYPE VARCHAR USING provider_id::VARCHAR;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE oauth
ALTER COLUMN provider_id TYPE UUID USING provider_id::UUID;
-- +goose StatementEnd
