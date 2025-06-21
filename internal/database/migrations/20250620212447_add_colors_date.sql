-- +goose Up
-- +goose StatementBegin
INSERT INTO colors(color)
VALUES 
  ('اسود'),
  ('ابيض'),
  ('احمر'),
  ('اصفر'),
  ('اخضر'),
  ('ازرق'),
  ('وردي'),
  ('بنفسجي'),
  ('جوزي'),
  ('ماروني'),
  ('رصاصي'),
  ('سمائي'),
  ('زيتوني'),
  ('برتقالي')
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM colors
WHERE color in ('برتقالي', 'زيتوني', 'سمائي', 'رصاصي', 'ماروني', 'جوزي', 'بنفسجي', 'وردي', 'ازرق', 'اخضر', 'اصفر', 'ابيض', 'احمر', 'اسود')
-- +goose StatementEnd
