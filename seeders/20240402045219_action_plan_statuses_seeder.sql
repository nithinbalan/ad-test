-- +goose Up
-- +goose StatementBegin
INSERT INTO dbo.action_plan_statuses (title_en, title_ar, color_code, [type], [status]) VALUES
('Active', 'نشط', '#00FF00', 'new', 'active'),
('In Progress', 'قيد التنفيذ', '#FFFF00', 'in_progress', 'active'),
('Completed', 'مكتمل', '#008000', 'completed', 'active'),
('On Hold', 'معلق', '#FFA500', 'on_hold', 'active'),
('Terminated', 'منتهية', '#FF0000', 'terminated', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE action_plan_statuses;
-- +goose StatementEnd
