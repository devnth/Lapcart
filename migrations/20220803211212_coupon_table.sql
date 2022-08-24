-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS coupons(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL UNIQUE,
    code VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    minimum_amount  BIGINT NOT NULL,
    value BIGINT NOT NULL,
    created_at TIMESTAMP,
    valid_till TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- DROP TABLE IF EXISTS coupons;
-- +goose StatementEnd
