-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS processor (
id BIGSERIAL PRIMARY KEY,
name VARCHAR(200),
description VARCHAR (500) DEFAULT 'processor description'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- DROP TABLE IF EXISTS processor;
-- +goose StatementEnd
