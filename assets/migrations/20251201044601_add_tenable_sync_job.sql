-- +goose Up
-- +goose StatementBegin
CREATE TABLE tenable_sync_jobs (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    lookback_days INT NOT NULL,
    status VARCHAR(20) NOT NULL,
    start_offset INT NOT NULL,
    current_offset INT NOT NULL,
    total_processed INT NOT NULL,
    error_message NVARCHAR(MAX) NULL,
    created_by BIGINT NOT NULL,
    started_at DATETIME NULL,
    finished_at DATETIME NULL,
    last_heartbeat DATETIME NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NOT NULL DEFAULT GETDATE()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tenable_sync_jobs;
-- +goose StatementEnd
