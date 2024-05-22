-- +goose Up
-- +goose StatementBegin
-- seed_meeting_frequencies.sql
INSERT INTO meeting_frequencies (title_en, title_ar, type, status ) VALUES
    ('Weekly', 'أسبوعي', 'weekly', 'active'),
    ('Monthly', 'شهري', 'monthly', 'active'),
    ('Quarterly', 'ربع سنوي', 'quarterly', 'active'),
    ('Half-Yearly', 'نصف سنوي', 'half-yearly', 'active'),
    ('Annually', 'سنوي', 'annually', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE meeting_frequencies;
-- +goose StatementEnd