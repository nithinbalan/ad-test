-- +goose Up
-- +goose StatementBegin
CREATE TABLE emails (
    id BIGINT PRIMARY KEY IDENTITY(1,1),
    module_id INT NULL,
    sub_module_id INT NULL,
    subject NVARCHAR(500) NOT NULL,
    body NVARCHAR(MAX) NOT NULL,
    email_status_id INT NOT NULL,
    is_permanent_failed BIT NOT NULL DEFAULT 0,
    attempt_count INT NOT NULL DEFAULT 0,
    max_attempts INT NOT NULL DEFAULT 5,
    next_retry_at DATETIME NULL,
    last_attempt_at DATETIME NULL,
    locked_at DATETIME NULL,
    locked_by NVARCHAR(64) NULL,
    smtp_remote_response NVARCHAR(MAX) NULL,
    error_message NVARCHAR(MAX) NULL,
    sent_at DATETIME NULL,
    created_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_emails_email_statuses FOREIGN KEY (email_status_id) REFERENCES email_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_emails_users FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_emails_modules FOREIGN KEY (module_id) REFERENCES modules(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_emails_sub_modules FOREIGN KEY (sub_module_id) REFERENCES modules(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE emails;
-- +goose StatementEnd
