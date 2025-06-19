-- +goose Up
-- +goose StatementBegin
INSERT INTO kpi_statuses
    (title_en, title_ar, color_code, type, status)
VALUES
    ('Active', 'نشيط', '#28A745', 'active', 'active'),
    ('Inactive', 'غير نشط', '#FFA500', 'inactive', 'active'),
    ('Archived', 'مؤرشف', '#DC3545', 'archived', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE kpi_statuses;
-- +goose StatementEnd
