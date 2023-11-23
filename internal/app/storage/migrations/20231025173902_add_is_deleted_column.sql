-- +goose Up
-- +goose StatementBegin
ALTER TABLE links
    ADD is_deleted bool NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE links
    DROP COLUMN is_deleted;
-- +goose StatementEnd
