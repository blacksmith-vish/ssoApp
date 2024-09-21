-- +goose Up
CREATE TABLE IF NOT EXISTS apps
(
    id     TEXT NOT NULL UNIQUE,
    name   TEXT NOT NULL UNIQUE,
    secret TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE IF EXISTS apps;