-- +goose Up
-- +goose StatementBegin
ALTER TABLE phone_verification_codes
RENAME TO account_verification_codes;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE account_verification_codes
RENAME TO phone_verification_codes;
-- +goose StatementEnd
