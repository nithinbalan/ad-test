-- +goose Up
-- +goose StatementBegin
INSERT INTO control_types (title_en, title_ar, status) VALUES
('Financial', 'المالية', 'active'),
('Operational', 'التشغيلية', 'active'),
('Strategic', 'الاستراتيجية', 'active'),
('Compliance', 'الامتثال', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE control_types;
-- +goose StatementEnd
