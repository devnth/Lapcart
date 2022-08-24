-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT  REFERENCES order_details(id),
    product_id BIGINT  REFERENCES product(id),
    discount_id BIGINT  REFERENCES discount(id),
    quantity BIGINT NOT NULL,
    created_at TIMESTAMP
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- 
-- +goose StatementEnd
