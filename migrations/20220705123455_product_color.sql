-- +goose Up
-- +goose StatementBegin
CREATE TABLE product_color (
    color_id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES product(product_id),
    color_name VARCHAR(200) NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product_color;
-- +goose StatementEnd
