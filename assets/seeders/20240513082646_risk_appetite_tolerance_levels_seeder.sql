-- +goose Up
-- +goose StatementBegin
INSERT INTO risk_appetite_tolerance_levels (title_en, title_ar, color_code, type, status) 
VALUES 
('Above', 'فوق', '#FF0000', 'above', 'active'),
('Within', 'في داخل', '#FFFF00', 'within', 'active'),
('Below', 'تحت', '#00FF00', 'below', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE risk_appetite_tolerance_levels;
-- +goose StatementEnd
