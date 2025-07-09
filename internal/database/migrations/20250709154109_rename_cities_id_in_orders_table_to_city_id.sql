-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders
RENAME COLUMN cities_id TO city_id
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders
RENAME COLUMN city_id TO cities_id
-- +goose StatementEnd
