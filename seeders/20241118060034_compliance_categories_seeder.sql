-- +goose Up
-- +goose StatementBegin

INSERT INTO compliance_categories (title_en, title_ar, status)
VALUES
    ('Regulatory Compliance', 'الامتثال التنظيمي', 'active'),
    ('Data Protection', 'حماية البيانات', 'active'),
    ('Environmental Compliance', 'الامتثال البيئي', 'active'),
    ('Health and Safety Compliance', 'امتثال الصحة والسلامة', 'active'),
    ('Financial Compliance', 'الامتثال المالي', 'active'),
    ('Cybersecurity Compliance', 'الامتثال الأمني السيبراني', 'active'),
    ('Industry Standards Compliance', 'امتثال معايير الصناعة', 'active'),
    ('Supplier Compliance', 'امتثال الموردين', 'active'),
    ('Product Quality Compliance', 'امتثال جودة المنتج', 'active'),
    ('Anti-Corruption Compliance', 'امتثال مكافحة الفساد', 'active');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

TRUNCATE TABLE compliance_categories;

-- +goose StatementEnd
