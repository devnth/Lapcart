-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS payment(
    id BIGSERIAL PRIMARY KEY,
    payment_type  VARCHAR(200) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payment;
-- +goose StatementEnd
