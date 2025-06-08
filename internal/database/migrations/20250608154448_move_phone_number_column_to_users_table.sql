-- +goose Up
-- +goose StatementBegin
ALTER TABLE local_auth
DROP COLUMN phone_number;

ALTER TABLE users
ADD COLUMN phone_number VARCHAR UNIQUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE local_auth
ADD COLUMN phone_number VARCHAR UNIQUE;

ALTER TABLE users
DROP COLUMN phone_number;
-- +goose StatementEnd
