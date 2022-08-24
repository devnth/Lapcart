-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS product (
id BIGSERIAL PRIMARY KEY,
code VARCHAR(200) NOT NULL,
name VARCHAR(400) NOT NULL,
description VARCHAR(500) DEFAULT 'product description',
color VARCHAR(200) NOT NULL,
brand_id BIGINT REFERENCES brand(id),
processor_id BIGINT REFERENCES processor(id),
category_id BIGINT REFERENCES category(id),
price NUMERIC NOT NULL,
rating NUMERIC DEFAULT 3.5,
"image" VARCHAR(300)  DEFAULT 'hhdfkhkdsah.jpg',
stock BIGINT DEFAULT 1,
discount_id BIGINT  ,
is_deleted BOOLEAN DEFAULT FALSE,
created_at TIMESTAMP DEFAULT NOW(),
updated_at TIMESTAMP DEFAULT NOW(),
deleted_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- DROP TABLE IF EXISTS product CASCADE;
-- +goose StatementEnd 
