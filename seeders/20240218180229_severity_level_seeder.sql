-- +goose Up
-- +goose StatementBegin
INSERT INTO dbo.severity_levels (title_en, title_ar, color_code, type) 
VALUES 
('High', 'عالي', '#FF0000', 'high'),
('Medium', 'متوسط', '#FFFF00', 'medium'),
('Low', 'منخفض', '#00FF00', 'low');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE severity_levels;
-- +goose StatementEnd
