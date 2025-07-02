-- +goose Up
-- +goose StatementBegin
INSERT INTO roles (title_en, title_ar, status, is_super_admin) VALUES
('Super Admin', 'المشرف الفائق', 'active', 1),
('Auditor', 'محرر', 'active', 0),
('Risk Manager', 'مشاهد', 'active', 0),
('Risk Executive', 'ضيف', 'active', 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE roles;
-- +goose StatementEnd
