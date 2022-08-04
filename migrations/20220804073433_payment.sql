-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS payement(
    id BIGSERIAL PRIMARY KEY,
    payment_type  VARCHAR(200) NOT NULL,
    amount BIGINT NOT NULL,
    status BOOLEAN DEFAULT false,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payement;
-- +goose StatementEnd
