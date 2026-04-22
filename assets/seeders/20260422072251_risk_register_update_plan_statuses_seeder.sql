-- +goose Up
-- +goose StatementBegin

INSERT INTO risk_register_update_plan_statuses (title_en, title_ar, color_code, type, status)
VALUES
    ('Not Started', N'لم يبدأ', '#17A2B8', 'not_started', 'active'),
    ('In Progress', N'قيد التنفيذ', '#FD7E14', 'in_progress', 'active'),
    ('Overdue', N'متأخر', '#DC3545', 'overdue', 'active'),
    ('Completed', N'مكتمل', '#28A745', 'completed', 'active');

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DELETE FROM risk_register_update_plan_statuses
WHERE type IN ('not_started', 'in_progress', 'overdue', 'completed');

-- +goose StatementEnd