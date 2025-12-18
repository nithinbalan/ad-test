-- +goose Up
-- +goose StatementBegin
ALTER TABLE tenable_sync_jobs
ALTER COLUMN created_by INT NULL;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE tenable_sync_jobs
ALTER COLUMN created_by INT NOT NULL;
-- +goose StatementEnd
