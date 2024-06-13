-- +goose Up
CREATE TABLE accounts (
                          id SERIAL PRIMARY KEY,
                          email VARCHAR NOT NULL,
                          password VARCHAR NOT NULL,
                          created_at TIMESTAMP NOT NULL,
                          logged_in_at TIMESTAMP
);
-- +goose Down
DROP TABLE accounts;
