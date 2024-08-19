-- +goose Up
-- +goose StatementBegin
INSERT INTO audit_schedule_statuses (title_en, title_ar, type, color_code, status)
VALUES 
('Scheduled', 'مجدول', 'scheduled', '#DC3545', 'active'),
('In Progress', 'قيد التنفيذ', 'in-progress', '#FFC107', 'active'),
('Completed', 'مكتمل', 'completed', '#28A745', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE audit_schedule_statuses;
-- +goose StatementEnd
