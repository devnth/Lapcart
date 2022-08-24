-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS address(
    id BIGSERIAL PRIMARY KEY,
    type VARCHAR DEFAULT 'Home Address',
    user_id INT NOT NULL REFERENCES users(id),
    house_name VARCHAR(300) NOT NULL,
    street_name VARCHAR(300) NOT NULL,
    landmark VARCHAR(300) NOT NULL,
    district VARCHAR(300) NOT NULL,
    state VARCHAR(300) NOT NULL,
    country VARCHAR(300) NOT NULL,
    pincode  BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- DROP TABLE IF EXISTS "address";
-- +goose StatementEnd
