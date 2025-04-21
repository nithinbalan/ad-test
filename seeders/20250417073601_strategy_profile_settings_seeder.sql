-- +goose Up
-- +goose StatementBegin
INSERT INTO strategy_profile_settings (title_en, title_ar, is_enable) VALUES
('Focus Area', 'منطقة التركيز', 1),
('Goals', 'الأهداف', 1),
('Objective', 'الهدف', 1),
('Initiative', 'المبادرة', 1),
('Projects', 'المشاريع', 1),
('Milestone', 'معلم رئيسي', 1),
('Action Plan', 'خطة العمل', 1),
('KPIs', 'مؤشرات الأداء الرئيسية', 1),
('Target', 'الهدف', 1);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE strategy_profile_settings;
-- +goose StatementEnd