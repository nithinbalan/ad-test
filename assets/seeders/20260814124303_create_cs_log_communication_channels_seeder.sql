-- +goose Up
-- +goose StatementBegin
INSERT INTO cs_log_communication_channels
    (title_en, title_ar, status)
VALUES
    ('Email', N'البريد الإلكتروني', 'active'),
    ('Phone', N'الهاتف', 'active'),
    ('SMS', N'الرسائل النصية', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM cs_log_communication_channels
WHERE title_en IN ('Email', 'Phone', 'SMS');
-- +goose StatementEnd
