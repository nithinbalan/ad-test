-- +goose Up
-- +goose StatementBegin
INSERT INTO third_party_question_categories (title_en, title_ar, status) VALUES
('Compliance', 'الامتثال', 'active'),
('Technical Security', 'الأمن التقني', 'active'),
('Data Protection & Privacy', 'حماية البيانات والخصوصية', 'active'),
('Financial Stability', 'الاستقرار المالي', 'active'),
('Business Continuity', 'استمرارية الأعمال', 'active'),
('Legal & Regulatory', 'القوانين والتنظيمات', 'active'),
('Contract Management', 'إدارة العقود', 'active'),
('Operational Risk', 'المخاطر التشغيلية', 'active'),
('Ethical Practices', 'الممارسات الأخلاقية', 'active'),
('Cybersecurity', 'الأمن السيبراني', 'active'),
('Third-Party Governance', 'حوكمة الطرف الثالث', 'active'),
('Risk Management', 'إدارة المخاطر', 'active'),
('Incident Management', 'إدارة الحوادث', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM third_party_question_categories;
-- +goose StatementEnd
