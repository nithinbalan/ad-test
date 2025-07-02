-- +goose Up
-- +goose StatementBegin
INSERT INTO grc_report_comment_statuses (title_en, title_ar, color_code, type, status) VALUES
('New', 'جديد', '#0000FF', 'new', 'active'),
('Resolved', 'تم الحل', '#FFD700', 'resolved', 'active'),
('Closed', 'مغلق', '#008000', 'closed', 'active');
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE grc_report_comment_statuses;
-- +goose StatementEnd