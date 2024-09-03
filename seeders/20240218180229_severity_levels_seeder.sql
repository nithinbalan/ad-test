-- +goose Up
-- +goose StatementBegin
INSERT INTO severity_levels (title_en, title_ar, color_code, type, status, weight) 
VALUES 
('High', 'عالي', '#DC3545', 'high', 'active', 10),
('Medium', 'متوسط', '#FFC107', 'medium', 'active', 5),
('Low', 'منخفض', '#28A745', 'low', 'active', 2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE severity_levels;
-- +goose StatementEnd
