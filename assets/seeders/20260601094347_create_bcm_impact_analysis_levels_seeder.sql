-- +goose Up
-- +goose StatementBegin
INSERT INTO bcm_impact_analysis_levels (title_en, title_ar, score, color_code, type, status) VALUES
    ('Non-applicable', 'غير قابل للتطبيق', 0, '#6C757D', 'non_applicable', 'active'),
    ('Very Low', 'منخفض جداً', 1, '#90DD75', 'very_low', 'active'),
    ('Low', 'منخفض', 2, '#28A745', 'low', 'active'),
    ('Moderate', 'متوسط', 3, '#FFC107', 'moderate', 'active'),
    ('High', 'عالي', 4, '#DC3545', 'high', 'active'),
    ('Very High', 'عالي جداً', 5, '#FF0000', 'very_high', 'active');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM bcm_impact_analysis_levels;
-- +goose StatementEnd
