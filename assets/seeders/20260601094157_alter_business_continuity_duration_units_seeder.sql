-- +goose Up
-- +goose StatementBegin
INSERT INTO business_continuity_duration_units (title_en, title_ar, minutes, type, status) VALUES
('Month', 'شهر', 43200, 'months', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM business_continuity_duration_units WHERE type = 'months';
-- +goose StatementEnd
