-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS product_processor (
processor_id BIGSERIAL PRIMARY KEY,
processor_name VARCHAR(200),
processor_desc VARCHAR (500),
processor_created_at TIMESTAMP DEFAULT NOW()
	-- product_processor_modified_at timestamp,
	-- product_processor_deleted_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product_processor;
-- +goose StatementEnd
