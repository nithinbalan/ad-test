-- +goose Up
-- +goose StatementBegin
ALTER TABLE activity_logs
ALTER COLUMN table_name VARCHAR(200) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE activity_logs
ALTER COLUMN table_name VARCHAR(45) NOT NULL;
-- +goose StatementEnd
