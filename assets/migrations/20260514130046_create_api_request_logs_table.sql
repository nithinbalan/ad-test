-- +goose Up
-- +goose StatementBegin
CREATE TABLE api_request_logs (
    id BIGINT PRIMARY KEY IDENTITY(1,1),
    method NVARCHAR(10) NOT NULL,
    url NVARCHAR(MAX) NOT NULL,
    request_id NVARCHAR(64) NOT NULL,
    user_id INT NULL,
    request_body NVARCHAR(MAX) NULL,
    response_body NVARCHAR(MAX) NULL,
    error_message NVARCHAR(MAX) NULL,
    stack_trace   NVARCHAR(MAX) NULL,
    requested_at DATETIME NOT NULL DEFAULT GETDATE(),
    completed_at DATETIME NULL,
    status_code INT NULL,
    duration_ms INT NULL,
    log_status NVARCHAR(20) NOT NULL DEFAULT 'pending',
    CONSTRAINT chk_api_logs_status
    CHECK (log_status IN ('pending','completed','failed','aborted')),
    CONSTRAINT FK_api_request_logs_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE api_request_logs;
-- +goose StatementEnd