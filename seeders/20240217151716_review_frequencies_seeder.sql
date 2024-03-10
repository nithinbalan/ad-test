-- +goose Up
-- +goose StatementBegin
-- seed_review_frequencies.sql
INSERT INTO review_frequencies (title_en, title_ar, type, status ) VALUES
    ('Daily', 'يومي', 'daily', 'active'),
    ('Weekly', 'أسبوعي', 'weekly', 'active'),
    ('Monthly', 'شهري', 'monthly', 'active'),
    ('Quarterly', 'ربع سنوي', 'quarterly', 'active'),
    ('Half-Yearly', 'نصف سنوي', 'half-yearly', 'active'),
    ('Annually', 'سنوي', 'annually', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE review_frequencies;
-- +goose StatementEnd