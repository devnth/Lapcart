-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS product_branding (
brand_id BIGSERIAL PRIMARY KEY,
brand_name VARCHAR(200),
brand_created_at timestamp default now()
-- product_brand_modified_at timestamp,
-- product_brand_deleted_at timestamp
	);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product_branding;
-- +goose StatementEnd
