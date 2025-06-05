-- +goose Up
-- +goose StatementBegin
ALTER TABLE phone_verification_codes
DROP CONSTRAINT phone_verification_codes_user_id_key;
-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin
ALTER TABLE phone_verification_codes
ADD CONSTRAINT phone_verification_codes_user_id_key UNIQUE (user_id)
-- +goose StatementEnd
