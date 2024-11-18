-- +goose Up
-- +goose StatementBegin

INSERT INTO compliance_requirement_categories (title_en, title_ar, status)
VALUES
    ('Legal Requirements', 'المتطلبات القانونية', 'active'),
    ('Regulatory Requirements', 'المتطلبات التنظيمية', 'active'),
    ('Contractual Obligations', 'الالتزامات التعاقدية', 'active'),
    ('Standards and Guidelines', 'المعايير والإرشادات', 'active'),
    ('Internal Policies', 'السياسات الداخلية', 'active'),
    ('Audit Requirements', 'متطلبات التدقيق', 'active'),
    ('Data Privacy Requirements', 'متطلبات خصوصية البيانات', 'active'),
    ('Cybersecurity Controls', 'ضوابط الأمن السيبراني', 'active'),
    ('Environmental Obligations', 'الالتزامات البيئية', 'active'),
    ('Health and Safety Obligations', 'الالتزامات الصحية والسلامة', 'active');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

TRUNCATE TABLE compliance_requirement_categories;

-- +goose StatementEnd
