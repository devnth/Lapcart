-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS verify_email (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(200) NOT NULL, 
    code BIGINT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS verify_email;
-- +goose StatementEnd
