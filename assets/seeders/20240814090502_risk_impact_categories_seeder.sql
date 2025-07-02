-- +goose Up
-- +goose StatementBegin
INSERT INTO risk_impact_categories (title_en, title_ar, status) VALUES
('Financial Impact', 'التأثير المالي', 'active'),
('Operational Impact', 'التأثير التشغيلي', 'active'),
('Reputational Impact', 'التأثير على السمعة', 'active'),
('Compliance Impact', 'تأثير الامتثال', 'active'),
('Legal Impact', 'التأثير القانوني', 'active'),
('Health and Safety Impact', 'التأثير على الصحة والسلامة', 'active'),
('Environmental Impact', 'التأثير البيئي', 'active'),
('Customer Impact', 'التأثير على العملاء', 'active'),
('Technological Impact', 'التأثير التكنولوجي', 'active'),
('Strategic Impact', 'التأثير الاستراتيجي', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE risk_impact_categories;
-- +goose StatementEnd
