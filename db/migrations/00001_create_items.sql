-- +goose Up
CREATE TABLE items (
    id   BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT   NOT NULL CHECK (length(name) BETWEEN 1 AND 100)
);

-- +goose Down
DROP TABLE items;
