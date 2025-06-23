-- +goose Up
-- +goose StatementBegin
UPDATE modules
SET title_ar = 'لوحه المعلومات'
WHERE identifier = 'grc_dashboard';

UPDATE modules
SET [order] = 6
WHERE identifier = 'compliance_report';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE modules
SET title_ar = 'لوحة معلومات الحوكمة وإدارة المخاطر والامتثال'
WHERE identifier = 'grc_dashboard';

UPDATE modules
SET [order] = 5
WHERE identifier = 'compliance_report';
-- +goose StatementEnd
