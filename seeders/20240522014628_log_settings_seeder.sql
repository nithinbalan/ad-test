-- +goose Up
-- +goose StatementBegin
INSERT INTO log_settings (audit_log_retention_days, email_log_retention_days, system_log_retention_days, updated_by, updated_at)
VALUES (30, 30, 30, 1, GETDATE());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE log_settings;
-- +goose StatementEnd
