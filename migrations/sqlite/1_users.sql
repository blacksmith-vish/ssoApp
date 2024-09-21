-- +goose Up
CREATE TABLE IF NOT EXISTS users
(
    id        TEXT NOT NULL UNIQUE,
    nickname  TEXT NOT NULL,
    email     TEXT NOT NULL UNIQUE,
    pass_hash BLOB NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_email ON users (email);

-- +goose Down
DROP TABLE IF EXISTS users;