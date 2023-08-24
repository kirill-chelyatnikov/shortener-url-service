-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS links (
    id varchar(10) PRIMARY KEY NOT NULL UNIQUE,
    baseURL text NOT NULL UNIQUE,
    hash varchar(64)[] NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE links
-- +goose StatementEnd
