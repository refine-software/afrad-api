-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'categories_name_parent_id_key'
  ) THEN
    ALTER TABLE categories
    ADD CONSTRAINT categories_name_parent_id_key UNIQUE (name, parent_id);
  END IF;
END
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE categories
DROP CONSTRAINT IF EXISTS categories_name_parent_id_key;
-- +goose StatementEnd
