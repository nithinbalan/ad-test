-- +goose Up
-- +goose StatementBegin
INSERT INTO awareness_channels 
(title_en, title_ar, status)
VALUES
    ('Email', 'البريد الإلكتروني', 'active'),
    ('SMS', 'الرسائل النصية', 'active'),
    ('Internal Portal', 'البوابة الداخلية', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM awareness_channels
-- +goose StatementEnd
