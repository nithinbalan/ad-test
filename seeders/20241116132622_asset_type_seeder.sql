-- +goose Up
-- +goose StatementBegin

INSERT INTO asset_types (title_en, title_ar, status)
VALUES
    ('IT Asset', 'أصل تكنولوجيا المعلومات', 'active'),
    ('Physical Asset', 'الأصل المادي', 'active'),
    ('Human Resource', 'الموارد البشرية', 'active'),
    ('Intellectual Property', 'الملكية الفكرية', 'active'),
    ('Financial Asset', 'الأصل المالي', 'active'),
    ('Infrastructure', 'البنية التحتية', 'active'),
    ('Data Asset', 'أصل البيانات', 'active'),
    ('Customer Information', 'معلومات العميل', 'active'),
    ('Software Asset', 'أصل البرمجيات', 'active'),
    ('Operational Equipment', 'المعدات التشغيلية', 'active');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE asset_types;

-- +goose StatementEnd
