-- +goose Up
-- +goose StatementBegin
ALTER TABLE local_auth
ALTER COLUMN phone_number DROP NOT NULL;

ALTER TABLE users
ALTER COLUMN email SET NOT NULL;

ALTER TABLE local_auth
RENAME COLUMN is_phone_verified to is_account_verified;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE local_auth
ALTER COLUMN phone_number SET NOT NULL;


ALTER TABLE users
ALTER COLUMN email DROP NOT NULL;


ALTER TABLE local_auth
RENAME COLUMN is_account_verified to is_phone_verified;
-- +goose StatementEnd
