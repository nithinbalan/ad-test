-- +goose Up
-- +goose StatementBegin
ALTER TABLE tenable_sync_jobs
ALTER COLUMN created_by BIGINT NULL;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE tenable_sync_jobs
ALTER COLUMN created_by BIGINT NOT NULL;
-- +goose StatementEnd
