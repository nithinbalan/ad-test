-- +goose Up
-- +goose StatementBegin
INSERT INTO report_template_settings (logo_image_path, cover_page_bg_image_path, created_by_id, created_at) VALUES
(NULL, NULL, 1, GETDATE());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE report_template_settings;
-- +goose StatementEnd
