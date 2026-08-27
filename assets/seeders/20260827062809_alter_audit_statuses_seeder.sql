-- +goose Up

-- +goose StatementBegin

INSERT INTO audit_statuses
    (title_en, title_ar, type, color_code, status)
VALUES
    ('Planning',  N'التخطيط',         'planning',  '#3B82F6', 'active'),
    ('Fieldwork', N'العمل الميداني', 'fieldwork', '#F59E0B', 'active'),
    ('Reporting', N'إعداد التقارير',  'reporting', '#10B981', 'active');

-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin

DELETE FROM audit_statuses
WHERE type IN ('planning', 'fieldwork', 'reporting');

-- +goose StatementEnd