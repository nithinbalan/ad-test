-- +goose Up
-- +goose StatementBegin
INSERT INTO dbo.action_plan_statuses (title_en, title_ar, color_code, [type], [status], is_closed) VALUES
('New', 'نشط', '#007BFF', 'new', 'active', 0),
('In Progress', 'قيد التنفيذ', '#17A2B8', 'in_progress', 'active', 0),
('Resolved', 'حلها', '#28A745', 'completed', 'active', 1),
('On Hold', 'معلق', '#FFC107', 'on_hold', 'active', 0),
('Terminated', 'منتهية', '#DC3545', 'terminated', 'active', 1),
('Close', 'مغلق', '#6C757D', 'close', 'active', 1),
('Reopen', 'إعادة فتح', '#6610F2', 'reopen', 'active', 0);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE action_plan_statuses;
-- +goose StatementEnd
