-- +goose Up
-- +goose StatementBegin
INSERT INTO applicability_statuses (title_en, title_ar, type, status) VALUES
('Applicable', 'ملائم', 'applicable', 'active'),
('Not Applicable', 'غير قابل للتطبيق', 'not-applicable', 'active');
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE applicability_statuses;
-- +goose StatementEnd