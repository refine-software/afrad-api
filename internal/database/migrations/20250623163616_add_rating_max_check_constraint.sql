-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_constraint
    WHERE conname = 'rating_max_check'
  ) THEN
    ALTER TABLE rating_review
    ADD CONSTRAINT rating_max_check CHECK (rating BETWEEN 1 AND 5);
  END IF;
END$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM pg_constraint
    WHERE conname = 'rating_max_check'
  ) THEN
    ALTER TABLE rating_review
    DROP CONSTRAINT rating_max_check;
  END IF;
END$$;
-- +goose StatementEnd
