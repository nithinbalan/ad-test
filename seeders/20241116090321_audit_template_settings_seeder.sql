-- +goose Up
-- +goose StatementBegin
INSERT INTO audit_template_settings (logo_image_path, commencement_letter_header_bg_image_path, draft_report_cover_page_bg_image_path, final_report_cover_page_bg_image_path, created_by_id, created_at) VALUES
(NULL, NULL, NULL, NULL, 1, GETDATE());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE audit_template_settings;
-- +goose StatementEnd
