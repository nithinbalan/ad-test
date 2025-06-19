-- +goose Up
-- +goose StatementBegin
INSERT INTO bcp_statuses (title_en, title_ar, color_code, type, status) VALUES
('Draft', 'مسودة', '#90DD75', 'draft', 'active'),
('In Review', 'قيد المراجعة', '#17A2B8', 'in-review', 'active'),
('Published', 'منشور', '#6C757D', 'published', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE bcp_statuses;
-- +goose StatementEnd
