-- +goose Up
-- +goose StatementBegin
INSERT INTO document_categories (title_en, title_ar, status) VALUES
('Information Security', 'أمن المعلومات', 'active'),
('Human Resources', 'الموارد البشرية', 'active'),
('Compliance', 'الامتثال', 'active'),
('Risk Management', 'إدارة المخاطر', 'active'),
('Financial Management', 'الإدارة المالية', 'active'),
('Data Protection', 'حماية البيانات', 'active'),
('Health and Safety', 'الصحة والسلامة', 'active'),
('Environmental', 'البيئة', 'active'),
('Governance', 'الحوكمة', 'active'),
('IT & Cybersecurity', 'تقنية المعلومات والأمن السيبراني', 'active'),
('Operational', 'العمليات', 'active'),
('Quality Control', 'مراقبة الجودة', 'active'),
('Customer Service', 'خدمة العملاء', 'active'),
('Supply Chain', 'سلسلة التوريد', 'active'),
('Crisis Management', 'إدارة الأزمات', 'active'),
('Employee Onboarding', 'تأهيل الموظفين الجدد', 'active'),
('Audit & Assurance', 'التدقيق والضمان', 'active'),
('Vendor Evaluation', 'تقييم البائعين', 'active'),
('Project Management', 'إدارة المشاريع', 'active'),
('Training Materials', 'مواد التدريب', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE document_categories;
-- +goose StatementEnd
