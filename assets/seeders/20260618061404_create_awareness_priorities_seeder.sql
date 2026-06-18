-- +goose Up
-- +goose StatementBegin
INSERT INTO awareness_priorities
(title_en, title_ar, type, status)
VALUES
    ('Very High', N'عالية جداً', 'very_high', 'active'),
    ('High', N'عالية', 'high', 'active'),
    ('Medium', N'متوسطة', 'medium', 'active'),
    ('Low', N'منخفضة', 'low', 'active'),
    ('Very Low', N'منخفضة جداً', 'very_low', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM awareness_priorities;
-- +goose StatementEnd