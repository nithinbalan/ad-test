-- +goose Up
-- +goose StatementBegin
INSERT INTO impact_categories (title_en, title_ar, status) VALUES
('Financial', 'مالي', 'active'),
('Operational', 'التشغيلية', 'active'),
('Reputational', 'السمعة', 'active'),
('Legal/Compliance', 'القانونية/الامتثال', 'active'),
('Human Resources', 'الموارد البشرية', 'active'),
('IT/System', 'تكنولوجيا المعلومات/النظام', 'active'),
('Supply Chain', 'الموردين', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE impact_categories;
-- +goose StatementEnd
