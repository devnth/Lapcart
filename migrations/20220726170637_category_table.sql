-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS category(
id BIGSERIAL PRIMARY KEY,
name VARCHAR(400), 
description VARCHAR(500) DEFAULT 'category description'
is_deleted BOOLEAN DEFAULT FALSE,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- DROP TABLE  IF EXISTS category;
-- +goose StatementEnd
