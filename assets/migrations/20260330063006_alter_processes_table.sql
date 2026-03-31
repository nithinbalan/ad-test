-- +goose Up
-- +goose StatementBegin
-- Alter column to allow NULL
ALTER TABLE processes
ALTER COLUMN responsible_user_id INT NULL;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
-- Alter column back to NOT NULL
ALTER TABLE processes
ALTER COLUMN responsible_user_id INT NULL;

-- +goose StatementEnd