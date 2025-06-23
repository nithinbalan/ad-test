-- +goose Up
-- +goose StatementBegin

INSERT INTO calculation_types (title_en, title_ar,type, status)
VALUES
    ('Sum', 'مجموع', 'sum', 'active'),
    ('Average', 'متوسط', 'average', 'active'),
    ('Max', 'الأعلى', 'max', 'active'),
    ('Min', 'دقيقة', 'min', 'active'),
    ('Count', 'عدد', 'count', 'active'),
    ('Custom Formula', 'صيغة مخصصة', 'custom-formula', 'active');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

TRUNCATE TABLE calculation_types;

-- +goose StatementEnd
