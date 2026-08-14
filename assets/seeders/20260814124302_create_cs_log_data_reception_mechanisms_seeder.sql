-- +goose Up
-- +goose StatementBegin
INSERT INTO cs_log_data_reception_mechanisms
    (title_en, title_ar, status)
VALUES
    ('Email', N'البريد الإلكتروني', 'active'),
    ('Interface', N'الواجهة', 'active'),
    ('Portal', N'البوابة', 'active'),
    ('Manual Delivery', N'التسليم اليدوي', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM cs_log_data_reception_mechanisms
WHERE title_en IN ('Email', 'Interface', 'Portal', 'Manual Delivery');
-- +goose StatementEnd
