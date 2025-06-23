-- +goose Up
-- +goose StatementBegin
INSERT INTO conformances
    (title_en, title_ar, type, color_code, status)
VALUES
    ('Yes', 'نعم', 'yes', '#28A745', 'active'),
    ('No', 'لا', 'no', '#DC3545', 'active'),
    ('Not Applicable', 'غير قابل للتطبيق', 'not-applicable', '#FFC107', 'active');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE conformances;
-- +goose StatementEnd
