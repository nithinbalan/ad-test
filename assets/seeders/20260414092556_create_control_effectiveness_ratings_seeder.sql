-- +goose Up
-- +goose StatementBegin
INSERT INTO control_effectiveness_ratings
(title_en, title_ar, type, min_score, max_score, reduction_percent, status)
VALUES
('Control missing', 'التحكم مفقود', 'control_missing', 0.0, 0.1, 0, 'active'),
('Ineffective', 'غير فعال', 'ineffective', 0.1, 2.0, 25, 'active'),
('Improvement Needed', 'يحتاج إلى تحسين', 'improvement_needed', 2.0, 3.5, 55, 'active'),
('Partially Effective', 'فعال جزئياً', 'partially_effective', 3.5, 5.0, 65, 'active'),
('Very Effective', 'فعال جداً', 'very_effective', 5.0, 5.0, 75, 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE control_effectiveness_ratings;
-- +goose StatementEnd
