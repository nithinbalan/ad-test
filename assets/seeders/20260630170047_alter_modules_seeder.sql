-- +goose Up
-- +goose StatementBegin

UPDATE modules
SET
    title_en = 'Dashboard',
    title_ar = N'لوحة المعلومات'
WHERE identifier = 'awareness_reporting_dashboard';
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
UPDATE modules
SET
    title_en = 'Awareness Reporting Dashboard',
    title_ar = N'لوحة معلومات تقارير التوعية'
WHERE identifier = 'awareness_reporting_dashboard';
-- +goose StatementEnd