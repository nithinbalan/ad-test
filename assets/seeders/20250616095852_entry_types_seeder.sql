-- +goose Up
-- +goose StatementBegin
INSERT INTO entry_types (title_en, title_ar,type, status)
VALUES
    ('Manual Entry', 'الإدخال اليدوي', 'manual', 'active'),
    ('Automated Entry', 'الإدخال الآلي', 'automated', 'active'),
    ('Semi-Automated', 'شبه آلي', 'semi_automated', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE entry_types;
-- +goose StatementEnd
