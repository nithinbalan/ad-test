-- +goose Up
-- +goose StatementBegin
INSERT INTO audit_finding_statuses (title_en, title_ar, type, is_closed, color_code)
VALUES 
('New', 'جديد', 'new', 0, '#007bff'),
('In Progress', 'قيد التنفيذ', 'in-progress', 0, '#ffc107'),
('Resolved', 'تم الحل', 'resolved', 0, '#28a745'),
('Reopened', 'أعيد فتحه', 'reopened', 0, '#17a2b8'),
('Closed', 'مغلق', 'closed', 1, '#6c757d'),
('Deferred', 'مؤجل', 'deferred', 0, '#fd7e14');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE audit_finding_statuses;
-- +goose StatementEnd
