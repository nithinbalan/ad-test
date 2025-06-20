-- +goose Up
-- +goose StatementBegin
CREATE TABLE email_logs (
    id BIGINT PRIMARY KEY IDENTITY(1,1),
    recipient_email NVARCHAR(255) NOT NULL,
    subject NVARCHAR(255) NOT NULL,
    template_path NVARCHAR(255) NOT NULL,
    email_data NVARCHAR(MAX),
    send_status NVARCHAR(50) NOT NULL,
    error_message NVARCHAR(MAX),
    sent_at DATETIME NULL,
    created_by_id INT,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_email_logs_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE email_logs;
-- +goose StatementEnd
