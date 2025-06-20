-- +goose Up
-- +goose StatementBegin
ALTER TABLE process_dependencies
ALTER COLUMN title NVARCHAR(255) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE process_dependencies
ALTER COLUMN title NVARCHAR(255) NOT NULL;
-- +goose StatementEnd
