-- +goose Up
-- +goose StatementBegin
INSERT INTO roles (title_en, title_ar, status) VALUES
('Admin', 'مدير النظام', 'active'),
('Auditor', 'محرر', 'active'),
('Risk Manager', 'مشاهد', 'active'),
('Risk Exicutive', 'ضيف', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE roles;
-- +goose StatementEnd
