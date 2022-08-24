-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS discount (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL UNIQUE,
    percentage NUMERIC(5,2) NOT NULL,
    created_at TIMESTAMP,
    valid_till TIMESTAMP,
    updated_at TIMESTAMP,
    status BOOLEAN DEFAULT true
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- DROP TABLE IF EXISTS discount;
-- +goose StatementEnd
