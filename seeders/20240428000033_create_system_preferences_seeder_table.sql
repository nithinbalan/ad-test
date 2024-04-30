-- +goose Up
-- +goose StatementBegin
INSERT INTO system_preferences (currency_id,timezone_id, date_format, time_format, doc_upload_size, logo_upload_size, idle_timeout, is_maintenance_mode_on, created_by, created_at)
VALUES
    (1, 1, 'YYYY-MM-DD', '12', 10, 20, 30, 1, 1, GETDATE());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE system_preferences;
-- +goose StatementEnd
