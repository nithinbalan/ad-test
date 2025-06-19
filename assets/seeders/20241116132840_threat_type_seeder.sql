-- +goose Up
-- +goose StatementBegin

INSERT INTO threat_types (title_en, title_ar, status)
VALUES
    ('Cyber Threat', 'تهديد إلكتروني', 'active'),
    ('Physical Threat', 'تهديد مادي', 'active'),
    ('Environmental Threat', 'تهديد بيئي', 'active'),
    ('Insider Threat', 'تهديد داخلي', 'active'),
    ('Operational Threat', 'تهديد تشغيلي', 'active'),
    ('Financial Threat', 'تهديد مالي', 'active'),
    ('Regulatory Threat', 'تهديد تنظيمي', 'active'),
    ('Reputational Threat', 'تهديد سمعة', 'active'),
    ('Health and Safety Threat', 'تهديد الصحة والسلامة', 'active'),
    ('Supply Chain Threat', 'تهديد سلسلة التوريد', 'active');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE threat_types;

-- +goose StatementEnd
