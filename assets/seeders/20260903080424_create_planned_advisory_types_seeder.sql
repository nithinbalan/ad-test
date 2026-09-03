-- +goose Up
-- +goose StatementBegin
INSERT INTO planned_advisory_types (title_en, title_ar, status) VALUES
    ('IOC', N'مؤشرات الاختراق', 'active'),
    ('IOB', N'مؤشرات السلوك', 'active'),
    ('Cases', N'الحالات', 'active'),
    ('Threats', N'التهديدات', 'active'),
    ('Update', N'التحديثات', 'active'),
    ('Other', N'أخرى', 'active');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM planned_advisory_types
WHERE title_en IN ('IOC', 'IOB', 'Cases', 'Threats', 'Update', 'Other');
-- +goose StatementEnd
