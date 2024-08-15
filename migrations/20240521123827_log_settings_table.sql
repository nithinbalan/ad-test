-- +goose Up
-- +goose StatementBegin
CREATE TABLE log_settings (
    id INT PRIMARY KEY IDENTITY(1,1),
    audit_log_retention_days INT NOT NULL DEFAULT 0,
    email_log_retention_days INT NOT NULL DEFAULT 0,
    system_log_retention_days INT NOT NULL DEFAULT 0,
    updated_by_id INT,
    updated_at DATETIME,
    CONSTRAINT FK_log_settings_users_1 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS log_settings;
-- +goose StatementEnd
