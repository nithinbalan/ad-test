-- +goose Up
-- +goose StatementBegin
INSERT INTO integration_categories (
    title_en,
    title_ar,
    description_ar,
    description_en,
    status
) VALUES (
    'SIEM',
    'سييم',
    'نظام لإدارة معلومات وأحداث الأمان لجمع وتحليل سجلات الأمان',
    'A Security Information and Event Management system for collecting and analyzing security logs',
    'active'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM integration_categories WHERE title_en = 'SIEM';
-- +goose StatementEnd
