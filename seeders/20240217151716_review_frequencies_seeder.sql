-- +goose Up
-- +goose StatementBegin
-- seed_review_frequencies.sql
INSERT INTO review_frequencies (title_en, title_ar, type) VALUES
    ('Daily', 'يومي', 'daily'),
    ('Weekly', 'أسبوعي', 'weekly'),
    ('Monthly', 'شهري', 'monthly'),
    ('Quarterly', 'ربع سنوي', 'quarterly'),
    ('Half-Yearly', 'نصف سنوي', 'half-yearly'),
    ('Annually', 'سنوي', 'annually');

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE review_frequencies;

-- +goose StatementEnd