-- +goose Up
-- +goose StatementBegin
ALTER TABLE phone_verification_codes
DROP CONSTRAINT IF EXISTS phone_verification_codes_user_id_key;
-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'phone_verification_codes_user_id_key'
  ) THEN
    ALTER TABLE phone_verification_codes
    ADD CONSTRAINT phone_verification_codes_user_id_key UNIQUE (user_id);
  END IF;
END
$$;
-- +goose StatementEnd
