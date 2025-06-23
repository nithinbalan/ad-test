-- +goose Up
-- +goose StatementBegin
INSERT INTO auditor_risk_ratings
    (title_en, title_ar, type, color_code, status)
VALUES
    ('Low', 'منخفض', 'low', '#28A745', 'active'),
    ('Medium', 'متوسط', 'medium', '#FFC107', 'active'),
    ('High', 'عالي', 'high', '#DC3545', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE auditor_risk_ratings;
-- +goose StatementEnd
