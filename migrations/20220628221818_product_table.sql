-- +goose Up
-- +goose StatementBegin
CREATE TABLE product (
product_id BIGSERIAL PRIMARY KEY,
product_name VARCHAR(400),
product_description VARCHAR(500),
product_brand BIGINT REFERENCES product_branding(brand_id),
product_processor BIGINT REFERENCES product_processor(processor_id),
product_category BIGINT REFERENCES product_category(category_id),
product_price NUMERIC,
product_rating NUMERIC,
product_created_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product;
-- +goose StatementEnd 
