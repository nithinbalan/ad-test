-- +goose Up
-- +goose StatementBegin
INSERT INTO control_categories (title_en, title_ar, status) VALUES
('Access Control', 'التحكم في الوصول', 'active'),
('Audit and Accountability', 'التدقيق والمساءلة', 'active'),
('Risk Assessment', 'تقييم المخاطر', 'active'),
('Incident Response', 'الاستجابة للحوادث', 'active'),
('Compliance Management', 'إدارة الامتثال', 'active'),
('Data Protection', 'حماية البيانات', 'active'),
('Policy Management', 'إدارة السياسات', 'active'),
('Change Management', 'إدارة التغيير', 'active'),
('Vendor Risk Management', 'إدارة مخاطر الموردين', 'active'),
('Business Continuity', 'استمرارية الأعمال', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE control_categories;
-- +goose StatementEnd
