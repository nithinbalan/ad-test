-- +goose Up
-- +goose StatementBegin
ALTER TABLE tenable_sync_jobs
DROP COLUMN lookback_days;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE tenable_sync_jobs
ADD lookback_days INT NULL;
-- +goose StatementEnd
