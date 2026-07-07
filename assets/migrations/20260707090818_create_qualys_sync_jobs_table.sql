-- +goose Up
-- +goose StatementBegin
CREATE TABLE qualys_sync_jobs (
    id               INT IDENTITY(1,1) PRIMARY KEY,
    status           VARCHAR(20)  NOT NULL ,
    start_offset     INT          NOT NULL ,
    current_offset   INT          NOT NULL ,
    total_processed  INT          NOT NULL ,
    started_at       DATETIME     NULL,
    finished_at      DATETIME     NULL,
    last_heartbeat   DATETIME     NULL,
    error_message    NVARCHAR(MAX) NULL,
    created_at       DATETIME     NOT NULL DEFAULT GETDATE(),
    updated_at       DATETIME     NULL 
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS qualys_sync_jobs;
-- +goose StatementEnd
