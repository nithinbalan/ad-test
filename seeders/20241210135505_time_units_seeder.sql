-- +goose Up
-- +goose StatementBegin
INSERT INTO time_units (title_en, title_ar, status) VALUES
('Hours', 'ساعات', 'active'),
('Minutes', 'دقائق', 'active'),
('Days', 'أيام', 'active'),
('Weeks', 'أسابيع', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE time_units;
-- +goose StatementEnd
