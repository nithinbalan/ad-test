-- +goose Up
-- +goose StatementBegin
INSERT INTO severity_levels (title_en, title_ar, color_code, type, status) 
VALUES 
('High', 'عالي', '#DC3545', 'high', 'active'),
('Medium', 'متوسط', '#FFC107', 'medium', 'active'),
('Low', 'منخفض', '#28A745', 'low', 'active');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE severity_levels;
-- +goose StatementEnd
