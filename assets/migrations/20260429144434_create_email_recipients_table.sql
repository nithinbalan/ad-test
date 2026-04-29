-- +goose Up
-- +goose StatementBegin
CREATE TABLE email_recipients (
    id BIGINT PRIMARY KEY IDENTITY(1,1),
    email_id BIGINT NOT NULL,
    user_id INT NULL,
    recipient_type NVARCHAR(10) NOT NULL DEFAULT 'to' 
    CONSTRAINT CK_email_recipients_type CHECK (recipient_type IN ('to', 'cc', 'bcc')),
    email_address NVARCHAR(500) NOT NULL,
    display_name NVARCHAR(500) NULL,
    is_external BIT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),

    CONSTRAINT FK_email_recipients_emails FOREIGN KEY (email_id) REFERENCES emails(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_email_recipients_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE email_recipients;
-- +goose StatementEnd
