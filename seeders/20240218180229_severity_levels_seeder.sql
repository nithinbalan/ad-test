-- +goose Up
-- +goose StatementBegin
INSERT INTO severity_levels (title_en, title_ar, color_code, type, status) 
VALUES 
('High', 'عالي', '#FF0000', 'high', 'active'),
('Medium', 'متوسط', '#FFFF00', 'medium', 'active'),
('Low', 'منخفض', '#00FF00', 'low', 'active');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE severity_levels;
-- +goose StatementEnd
