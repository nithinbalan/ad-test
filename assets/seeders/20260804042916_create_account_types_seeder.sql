-- +goose Up
-- +goose StatementBegin
INSERT INTO account_types (title_en, title_ar, type, status)
VALUES
    ('Human', N'بشري', 'human', 'active'),
    ('Admin', N'مسؤول', 'admin', 'active'),
    ('Service', N'خدمة', 'service', 'active'),
    ('Vendor', N'مورد', 'vendor', 'active'),
    ('Shared', N'مشترك', 'shared', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM account_types
WHERE type IN ('human', 'admin', 'service', 'vendor', 'shared');
-- +goose StatementEnd

    