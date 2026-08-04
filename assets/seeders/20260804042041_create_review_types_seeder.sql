-- +goose Up
-- +goose StatementBegin
INSERT INTO review_types (title_en, title_ar, type, status)
VALUES
    ('User', N'مستخدم', 'user', 'active'),
    ('Privileged', N'ذو صلاحيات مميزة', 'privileged', 'active'),
    ('Service', N'خدمة', 'service', 'active'),
    ('External', N'خارجي', 'external', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM review_types
WHERE type IN ('user', 'privileged', 'service', 'external');
-- +goose StatementEnd

