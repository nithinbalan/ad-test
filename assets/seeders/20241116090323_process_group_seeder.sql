-- +goose Up
-- +goose StatementBegin
INSERT INTO process_groups (title_en, title_ar, status) VALUES
('Risk Management', 'إدارة المخاطر', 'active'),
('Compliance Management', 'إدارة الامتثال', 'active'),
('Policy Management', 'إدارة السياسات', 'active'),
('Audit Management', 'إدارة التدقيق', 'active'),
('Incident Management', 'إدارة الحوادث', 'active'),
('Change Management', 'إدارة التغيير', 'active'),
('IT Governance', 'حوكمة تكنولوجيا المعلومات', 'active'),
('Data Protection', 'حماية البيانات', 'active'),
('Access Control', 'التحكم في الوصول', 'active'),
('Business Continuity', 'استمرارية الأعمال', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE process_groups;
-- +goose StatementEnd
