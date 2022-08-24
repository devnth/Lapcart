-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS brand(
id BIGSERIAL PRIMARY KEY,
name VARCHAR(200),
description VARCHAR(500) DEFAULT 'brand description'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- DROP TABLE IF EXISTS brand;
-- +goose StatementEnd
