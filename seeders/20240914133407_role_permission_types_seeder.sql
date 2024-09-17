-- +goose Up
-- +goose StatementBegin
INSERT INTO role_permission_types
    (title_en, title_ar, type, status)
VALUES
    ('Manage', 'إدارة', 'manage', 'active'),
    ('Modify', 'تعديل', 'modify', 'active'),
    ('View', 'عرض', 'view', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE role_permission_types;
-- +goose StatementEnd
