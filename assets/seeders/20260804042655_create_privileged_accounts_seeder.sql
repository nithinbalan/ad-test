-- +goose Up
-- +goose StatementBegin
INSERT INTO privileged_accounts (title_en, title_ar, status)
VALUES
    ('Domain Admin', N'مسؤول المجال', 'active'),
    ('Database Admin', N'مسؤول قاعدة البيانات', 'active'),
    ('Firewall Admin', N'مسؤول جدار الحماية', 'active'),
    ('Server Root', N'المستخدم الجذر للخادم', 'active'),
    ('Global Administrator', N'المسؤول العام', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM privileged_accounts
WHERE title_en IN (
    'Domain Admin',
    'Database Admin',
    'Firewall Admin',
    'Server Root',
    'Global Administrator'
);
-- +goose StatementEnd
