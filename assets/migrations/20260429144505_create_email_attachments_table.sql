-- +goose Up
-- +goose StatementBegin
CREATE TABLE email_attachments (
    id BIGINT PRIMARY KEY IDENTITY(1,1),
    email_id BIGINT NOT NULL,
    storage_key NVARCHAR(512) NOT NULL,
    file_name NVARCHAR(255) NOT NULL,
    content_type NVARCHAR(128) NULL,
    size_bytes BIGINT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),

    CONSTRAINT FK_email_attachments_emails FOREIGN KEY (email_id) REFERENCES emails(id) ON DELETE CASCADE ON UPDATE NO ACTION
);

CREATE INDEX IX_email_attachments_email_id ON email_attachments(email_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE email_attachments;
-- +goose StatementEnd
