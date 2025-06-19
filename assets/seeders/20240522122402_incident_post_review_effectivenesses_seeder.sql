-- +goose Up
-- +goose StatementBegin
INSERT INTO incident_post_review_effectivenesses
    (title_en, title_ar, color_code, status)
VALUES
    ('High', 'عالي', '#00FF00', 'active'),
    ('Moderate', 'متوسط', '#FFFF00', 'active'),
    ('Low', 'منخفض', '#FFA500', 'active'),
    ('Needs Improvement', 'يحتاج لتحسين', '#FF0000', 'active'),
    ('Not Effective', 'غير فعال', '#000000', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE incident_post_review_effectivenesses;
-- +goose StatementEnd
