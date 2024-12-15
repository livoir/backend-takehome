-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL
);
CREATE INDEX index_email_table_users ON users (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX index_email_table_users ON users;
DROP TABLE users;
-- +goose StatementEnd
