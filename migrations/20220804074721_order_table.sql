-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_details (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    shipping_address_id BIGINT REFERENCES address(id),
    coupon_id BIGINT REFERENCES coupons(id),
    is_paid  BOOLEAN DEFAULT false,
    payment_id BIGINT REFERENCES payment(id),
    status TEXT DEFAULT 'waiting',
    total NUMERIC(20,2) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- DROP TABLE IF EXISTS order_details ;
-- +goose StatementEnd
