-- +goose Up
-- +goose StatementBegin
INSERT INTO audit_risk_ratings
    (title_en, title_ar, type, color_code, status)
VALUES
    ('Very Low', 'منخفض جدا', 'very_low', '#00FF00', 'active'),
    ('Low', 'منخفض', 'low', '#FFFF00', 'active'),
    ('Medium', 'متوسط', 'medium', '#FFA500', 'active'),
    ('High', 'عالي', 'high', '#FF0000', 'active'),
    ('Very High', 'عالي جدا', 'very_high', '#800080', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE audit_risk_ratings;
-- +goose StatementEnd
