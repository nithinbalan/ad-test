-- +goose Up
-- +goose StatementBegin
INSERT INTO violation_report_channels
    (title_en, title_ar, description_en, description_ar, status)
VALUES
    ('Website', 'الموقع الإلكتروني', NULL, NULL, 'active'),
    ('Telephone', 'الهاتف', NULL, NULL, 'active'),
    ('In-Person', 'حضورياً', NULL, NULL, 'active');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM violation_report_channels
WHERE title_en IN ('Website', 'Telephone', 'In-Person');
-- +goose StatementEnd
