-- +goose Up
-- +goose StatementBegin
INSERT INTO dbo.action_plan_statuses (title_en, title_ar, color_code, [type], [status]) VALUES
('New', 'نشط', '#007BFF', 'new', 'active'),
('In Progress', 'قيد التنفيذ', '#17A2B8', 'in_progress', 'active'),
('Completed', 'مكتمل', '#28A745', 'completed', 'active'),
('On Hold', 'معلق', '#FFC107', 'on_hold', 'active'),
('Terminated', 'منتهية', '#DC3545', 'terminated', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE action_plan_statuses;
-- +goose StatementEnd
