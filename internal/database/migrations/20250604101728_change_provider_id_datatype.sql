-- +goose Up
ALTER TABLE oauth
ALTER COLUMN provider_id TYPE VARCHAR USING provider_id::VARCHAR;

-- +goose Down
ALTER TABLE oauth
ALTER COLUMN provider_id TYPE UUID USING provider_id::UUID;
