-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_update_discount_updated_at ON discounts;
DROP TRIGGER IF EXISTS  trigger_update_order_updated_at ON orders;
DROP TRIGGER IF EXISTS trigger_update_product_variant_updated_at ON product_variants;
DROP TRIGGER IF EXISTS trigger_update_review_updated_at ON rating_review;
DROP TRIGGER IF EXISTS trigger_update_session_updated_at ON sessions;
DROP TRIGGER IF EXISTS trigger_update_user_updated_at ON users;

CREATE TRIGGER trigger_update_discount_updated_at
BEFORE UPDATE ON discounts
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trigger_update_order_updated_at
BEFORE UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();


CREATE TRIGGER trigger_update_product_variant_updated_at
BEFORE UPDATE ON product_variants
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();


CREATE TRIGGER trigger_update_review_updated_at
BEFORE UPDATE ON rating_review
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();


CREATE TRIGGER trigger_update_session_updated_at
BEFORE UPDATE ON sessions
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();


CREATE TRIGGER trigger_update_user_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- +goose StatementEnd


-- +goose Down

-- +goose StatementBegin
DROP TRIGGER IF EXISTS trigger_update_discount_updated_at ON discounts;
DROP TRIGGER IF EXISTS  trigger_update_order_updated_at ON orders;
DROP TRIGGER IF EXISTS trigger_update_product_variant_updated_at ON product_variants;
DROP TRIGGER IF EXISTS trigger_update_review_updated_at ON rating_review;
DROP TRIGGER IF EXISTS trigger_update_session_updated_at ON sessions;
DROP TRIGGER IF EXISTS trigger_update_user_updated_at ON users;
DROP FUNCTION IF EXISTS set_updated_at;
-- +goose StatementEnd
