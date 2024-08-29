-- +goose Up
-- +goose StatementBegin
INSERT INTO audit_finding_statuses (title_en, title_ar, type, is_closed, color_code)
VALUES 
('New', 'جديد', 'new', 0, '#007BFF'),
('In Progress', 'قيد التنفيذ', 'in-progress', 0, '#FFC107'),
('Resolved', 'تم الحل', 'resolved', 1, '#28A745'),
('Reopened', 'أعيد فتحه', 'reopened', 0, '#17A2B8'),
('Closed', 'مغلق', 'closed', 1, '#6C757D'),
('Deferred', 'مؤجل', 'deferred', 1, '#FD7E14');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE audit_finding_statuses;
-- +goose StatementEnd
