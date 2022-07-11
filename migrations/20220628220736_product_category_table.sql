-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS product_category(
category_id BIGSERIAL PRIMARY KEY,
category_name VARCHAR(200), 
category_desc VARCHAR(500),
category_created_at timestamp DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE  IF EXISTS product_category;
-- +goose StatementEnd
