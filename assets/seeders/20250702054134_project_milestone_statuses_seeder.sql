-- +goose Up
-- +goose StatementBegin
INSERT INTO project_milestone_statuses (
    title_en, title_ar, description_en, description_ar, color_code, type, status
) VALUES
('Pending', 'قيد الانتظار', 'Milestone is pending managerial or client approval.', 'المرحلة بانتظار موافقة الإدارة أو العميل.', '#17A2B8', 'pending', 'active'),
('In Progress', 'قيد التنفيذ', 'Work is actively ongoing for this milestone.', 'يتم العمل حالياً على هذه المرحلة.', '#007BFF', 'in_progress', 'active'),
('Completed', 'مكتمل', 'Milestone has been completed successfully.', 'تم إكمال المرحلة بنجاح.', '#28A745', 'completed', 'active'),
('Cancelled', 'ملغى', 'Milestone has been cancelled and will not proceed.', 'تم إلغاء المرحلة ولن تستمر.', '#DC3545', 'cancelled', 'active')
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE project_milestone_statuses;
-- +goose StatementEnd

