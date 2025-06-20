-- +goose Up
-- +goose StatementBegin
INSERT INTO kri_categories (title_en, title_ar, status) VALUES
('Financial Risk', 'المخاطر المالية', 'active'),
('Operational Risk', 'المخاطر التشغيلية', 'active'),
('Compliance Risk', 'مخاطر الامتثال', 'active'),
('Customer Risk', 'مخاطر العملاء', 'active'),
('Supply Chain Risk', 'مخاطر سلسلة التوريد', 'active'),
('Technology Risk', 'مخاطر التكنولوجيا', 'active'),
('Market Risk', 'مخاطر السوق', 'active'),
('Health and Safety Risk', 'مخاطر الصحة والسلامة', 'active'),
('Reputational Risk', 'مخاطر السمعة', 'active'),
('Environmental Risk', 'مخاطر البيئة', 'active'),
('Project Risk', 'مخاطر المشروع', 'active'),
('Human Resources Risk', 'مخاطر الموارد البشرية', 'active'),
('Strategic Risk', 'المخاطر الاستراتيجية', 'active'),
('Legal Risk', 'المخاطر القانونية', 'active'),
('Third-Party Risk', 'مخاطر الأطراف الثالثة', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE kri_categories;
-- +goose StatementEnd
