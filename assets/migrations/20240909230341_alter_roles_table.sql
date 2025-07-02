-- +goose Up
-- +goose StatementBegin
ALTER TABLE roles
ADD is_super_admin BIT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE roles
DROP COLUMN is_super_admin;
-- +goose StatementEnd
