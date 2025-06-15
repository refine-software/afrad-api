-- +goose Up
-- +goose StatementBegin
ALTER TABLE products
ADD COLUMN IF NOT EXISTS thumbnail TEXT NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE products
DROP COLUMN IF EXISTS thumbnail;
-- +goose StatementEnd
