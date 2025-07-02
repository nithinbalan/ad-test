-- +goose Up
-- +goose StatementBegin
INSERT INTO dbo.risk_levels ([type], title_en, title_ar, color_code, [status]) VALUES
('high', 'High', 'عالي', '#FF0000', 'active'),
('medium', 'Medium', 'متوسط', '#FFFF00', 'active'),
('low', 'Low', 'منخفض', '#00FF00', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE risk_levels;
-- +goose StatementEnd
