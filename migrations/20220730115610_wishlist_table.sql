-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS wishlist (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    product_code VARCHAR(200) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- DROP TABLE IF EXISTS wishlist;
-- +goose StatementEnd
