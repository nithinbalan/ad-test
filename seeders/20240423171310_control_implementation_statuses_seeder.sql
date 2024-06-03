-- +goose Up
-- +goose StatementBegin
INSERT INTO control_implementation_statuses (title_en, title_ar, color_code, type, status) VALUES
('Implemented', 'منفذ', '#00FF00', 'implemented', 'active'),
('Not Implemented', 'غير منفذ', '#FF0000', 'not-implemented', 'active'),
('Partially Implemented', 'منفذ جزئيا', '#FFFF00', 'partially-implemented', 'active'),
('Not Applicable', 'غير قابل للتطبيق', '#0000FF', 'not-applicable', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE control_implementation_statuses;
-- +goose StatementEnd
