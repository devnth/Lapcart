-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cart(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES "users"(id),
    product_id BIGINT NOT NULL REFERENCES product(id),
    count BIGINT DEFAULT 1,
    created_at TIMESTAMP,
    updated_at TIMESTAMP 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- DROP TABLE IF EXISTS cart;
-- +goose StatementEnd
