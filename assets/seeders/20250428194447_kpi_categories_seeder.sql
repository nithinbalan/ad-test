-- +goose Up
-- +goose StatementBegin

INSERT INTO kpi_categories (title_en, title_ar, status)
VALUES
    ('Financial', 'مالي', 'active'),
    ('Operational', 'التشغيلية', 'active'),
    ('Compliance', 'امتثال', 'active'),
    ('Customer', 'عميل', 'active');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

TRUNCATE TABLE kpi_categories;

-- +goose StatementEnd
