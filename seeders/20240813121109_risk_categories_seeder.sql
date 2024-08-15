-- +goose Up
-- +goose StatementBegin
INSERT INTO risk_categories (title_en, title_ar, status) VALUES
('Strategic Risk', 'المخاطر الاستراتيجية', 'active'),
('Operational Risk', 'المخاطر التشغيلية', 'active'),
('Compliance Risk', 'مخاطر الامتثال', 'active'),
('Financial Risk', 'المخاطر المالية', 'active'),
('Reputational Risk', 'مخاطر السمعة', 'active'),
('Cybersecurity Risk', 'مخاطر الأمن السيبراني', 'active'),
('Environmental Risk', 'المخاطر البيئية', 'active'),
('Market Risk', 'مخاطر السوق', 'active'),
('Legal Risk', 'المخاطر القانونية', 'active'),
('Health and Safety Risk', 'مخاطر الصحة والسلامة', 'active'),
('Technological Risk', 'المخاطر التكنولوجية', 'active'),
('Project Risk', 'مخاطر المشروع', 'active'),
('Human Resources Risk', 'مخاطر الموارد البشرية', 'active'),
('Supply Chain Risk', 'مخاطر سلسلة التوريد', 'active'),
('Data Privacy Risk', 'مخاطر خصوصية البيانات', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE risk_categories;
-- +goose StatementEnd
