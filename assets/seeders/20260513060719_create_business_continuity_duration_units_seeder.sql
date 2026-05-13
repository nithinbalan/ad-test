-- +goose Up
-- +goose StatementBegin

INSERT INTO business_continuity_duration_units (title_en, title_ar, minutes, type, status) VALUES
('Minute', 'دقيقة', 1, 'minutes', 'active'),
('Hour', 'ساعة', 60, 'hours', 'active'),
('Day', 'يوم', 1440, 'days', 'active'),
('Week', 'أسبوع', 10080, 'weeks', 'active');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM business_continuity_duration_units;

-- +goose StatementEnd
