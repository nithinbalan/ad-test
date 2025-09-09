-- +goose Up
-- +goose StatementBegin
INSERT INTO project_statuses (
    title_en, title_ar, description_en, description_ar, color_code, type, status
) VALUES
('Not Started', 'لم يبدأ', 'Project has been created but work has not yet begun.', 'تم إنشاء المشروع ولكن لم يبدأ العمل بعد.', '#6C757D', 'not_started', 'active'),
('Planning', 'تخطيط', 'Project is in the planning phase.', 'المشروع في مرحلة التخطيط.', '#17A2B8', 'planning', 'active'),
('In Progress', 'قيد التنفيذ', 'Project work is currently underway.', 'العمل على المشروع جارٍ حالياً.', '#007BFF', 'in_progress', 'active'),
('On Hold', 'معلق', 'Project is temporarily paused.', 'المشروع متوقف مؤقتاً.', '#FFC107', 'on_hold', 'active'),
('Cancelled', 'ملغى', 'Project has been cancelled.', 'تم إلغاء المشروع.', '#DC3545', 'cancelled', 'active'),
('Completed', 'مكتمل', 'Project has been successfully completed.', 'تم الانتهاء من المشروع بنجاح.', '#28A745', 'completed', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE project_statuses;
-- +goose StatementEnd
