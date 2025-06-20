-- +goose Up
-- +goose StatementBegin
INSERT INTO control_types (title_en, title_ar, status) VALUES
('Preventive', 'الوقائية', 'active'),
('Detective', 'الكشفية', 'active'),
('Directive', 'التوجيهية', 'active'),
('Corrective', 'التصحيحية', 'active'),
('Compensating', 'التعويضية', 'active'),
('Physical', 'الفيزيائية', 'active'),
('Administrative', 'الإدارية', 'active'),
('Technical', 'التقنية', 'active'),
('Feedback', 'الارتجاعية', 'active'),
('Predictive', 'التنبؤية', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE control_types;
-- +goose StatementEnd
