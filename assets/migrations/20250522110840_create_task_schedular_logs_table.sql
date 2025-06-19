-- +goose Up
-- +goose StatementBegin
CREATE TABLE task_schedular_logs (
    id INT IDENTITY(1,1) PRIMARY KEY,
    task_name NVARCHAR(255) NOT NULL,
    status NVARCHAR(50) NOT NULL,         -- e.g., STARTED, COMPLETED, ERROR
    message NVARCHAR(MAX) NULL,
    run_at DATETIME NOT NULL DEFAULT GETDATE(),
    completed_at DATETIME NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE task_schedular_logs;
-- +goose StatementEnd
