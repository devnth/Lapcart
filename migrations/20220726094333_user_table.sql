-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
      id BIGSERIAL PRIMARY KEY,
      first_name VARCHAR(200) NOT NULL,
      last_name VARCHAR(200) NOT NULL,
      password VARCHAR(200) NOT NULL,
      email VARCHAR(200) NOT NULL,
      phone_number BIGINT,
      is_active BOOLEAN DEFAULT TRUE,
      is_verified BOOLEAN DEFAULT FALSE,
      created_at TIMESTAMP DEFAULT NOW(),
      updated_at TIMESTAMP DEFAULT NOW()
   );
INSERT INTO users (first_name, last_name, password, email, phone_number)
VALUES 
 ('Dev', 'Anil', '5f4dcc3b5aa765d61d8327deb882cf99', 'dev@email.com', 9999999898), 
 ('Tanu', 'Anil', '5f4dcc3b5aa765d61d8327deb882cf99', 'Tanu@email.com', 9999999898);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- DROP TABLE IF EXISTS users CASCADE;
-- +goose StatementEnd
