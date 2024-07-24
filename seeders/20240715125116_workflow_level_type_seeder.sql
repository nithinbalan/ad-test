-- +goose Up
-- +goose StatementBegin
INSERT INTO workflow_level_types (type, title_en, title_ar, status) VALUES 
('individual', 'Individual', 'فردي', 'active'),
('multiple-users', 'Multiple Users', 'مستخدمون متعددون', 'active'),
('role', 'Role', 'دور', 'inactive'),
('hod', 'Head of Department', 'رئيس القسم', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE workflow_level_types;
-- +goose StatementEnd
