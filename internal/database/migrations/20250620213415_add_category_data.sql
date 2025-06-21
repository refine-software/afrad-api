-- +goose Up
-- +goose StatementBegin
DO $$
DECLARE
  gender_names TEXT[] := ARRAY['رجالي', 'نسائي', 'كلا الجنسين'];
  gender_name TEXT;
  gender_id INT;
  top_id INT;
  bottom_id INT;
  shoe_id INT;
  one_piece_id INT;
  accessories_id INT;
  pant_id INT;
  pajama_id INT;

BEGIN
  FOREACH gender_name IN ARRAY gender_names LOOP
    INSERT INTO categories (name, parent_id) VALUES (gender_name, NULL)
    RETURNING id INTO gender_id;


    INSERT INTO categories (name, parent_id) VALUES ('احذية', gender_id)
    RETURNING id INTO shoe_id;
    INSERT INTO categories (name, parent_id) VALUES 
      ('رسمي', shoe_id),
      ('رياضي', shoe_id),
      ('كاجوال', shoe_id);

    INSERT INTO categories (name, parent_id) VALUES ('فوقي', gender_id)
    RETURNING id INTO top_id;
    INSERT INTO categories (name, parent_id) VALUES 
      ('تيشيرتات', top_id),
      ('قمصان', top_id),
      ('هوديات', top_id),
      ('سويترات', top_id),
      ('بلوزات', top_id),
      ('قماصل', top_id),
      ('كوتات', top_id);


    INSERT INTO categories (name, parent_id) VALUES ('سفلي', gender_id)
    RETURNING id INTO bottom_id;
    
    INSERT INTO categories (name, parent_id) VALUES ('بناطير', bottom_id)
    RETURNING id INTO pant_id;
    INSERT INTO categories (name, parent_id) VALUES 
      ('جينز', pant_id),
      ('رسمي', pant_id),
      ('كارغو', pant_id),
      ('جلد', pant_id);

    INSERT INTO categories (name, parent_id) VALUES 
      ('شورتات', bottom_id);

    IF gender_name = 'نسائي'  THEN
       INSERT INTO categories (name, parent_id) VALUES ('تنانير', bottom_id);
    END IF;


    INSERT INTO categories (name, parent_id) VALUES ('بجامات', bottom_id)
    RETURNING id INTO pajama_id;
    INSERT INTO categories (name, parent_id) VALUES 
      ('رياضي', pajama_id),
      ('نوم', pajama_id);

    
    INSERT INTO categories (name, parent_id) VALUES ('كامل', gender_id)
    RETURNING id INTO one_piece_id;

    IF gender_name = 'رجالي' THEN 
      INSERT INTO categories (name, parent_id) VALUES ('دشداشة', one_piece_id);
    END IF;

    IF gender_name = 'نسائي'  THEN
       INSERT INTO categories (name, parent_id) VALUES ('تنانير', one_piece_id);
    END IF;

    INSERT INTO categories (name, parent_id) VALUES ('اكسسوارات', gender_id)
    RETURNING id INTO accessories_id;
    INSERT INTO categories (name, parent_id) VALUES 
      ('شفقات', accessories_id),
      ('احزمة', accessories_id),
      ('نظارات', accessories_id),
      ('سوارات', accessories_id),
      ('محابس', accessories_id),
      ('قلادات', accessories_id);
  END LOOP;
END $$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM categories WHERE name IN (
  'رجالي', 'نسائي', 'كلا الجنسين',
  'احذية', 'رسمي', 'رياضي', 'كاجوال',
  'فوقي', 'تيشيرتات', 'قمصان', 'هوديات', 'سويترات', 'بلوزات', 'قماصل', 'كوتات',
  'سفلي', 'بناطير', 'جينز', 'كارغو', 'جلد', 'تنانير', 'بجامات', 'رياضي', 'نوم', 'شورتات',
  'كامل', 'فستان', 'دشداشة',
  'اكسسوارات', 'شفقات', 'احزمة', 'نظارات', 'سوارات', 'محابس', 'قلادات'
);
-- +goose StatementEnd
