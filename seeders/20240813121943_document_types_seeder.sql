-- +goose Up
-- +goose StatementBegin
INSERT INTO document_types (title_en, title_ar, type, status) VALUES
('Policy', 'سياسة', 'policy', 'active'),
('Procedure', 'إجراء', 'procedure', 'active'),
('Checklist', 'قائمة فحص', 'others', 'active'),
('Manuals', 'أدلة', 'others', 'active'),
('Standards', 'معايير', 'others', 'active'),
('Guidelines', 'إرشادات', 'others', 'active'),
('Forms', 'نماذج', 'others', 'active'),
('Training Materials', 'مواد التدريب', 'others', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE document_types;
-- +goose StatementEnd
