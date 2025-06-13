-- +goose Up
-- +goose StatementBegin
ALTER TABLE categories DROP CONSTRAINT IF EXISTS categories_name_key;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'categories_name_key'
  ) THEN
    ALTER TABLE categories
    ADD CONSTRAINT categories_name_key UNIQUE (name);
  END IF;
END
$$;
-- +goose StatementEnd
