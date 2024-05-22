-- +goose Up
-- +goose StatementBegin
INSERT INTO meeting_statuses
    (title_en, title_ar, color_code, type, status)
VALUES
    ('Draft', 'مسودة', '#A0A0A0', 'draft', 'active'),
    ('Scheduled', 'مجدول', '#0000FF', 'scheduled', 'active'),
    ('Completed', 'مكتمل', '#008000', 'completed', 'active'),
    ('In Progress', 'قيد التنفيذ', '#FFFF00', 'in_progress', 'active'),
    ('Postponed', 'مؤجل', '#FFA500', 'postponed', 'active'),
    ('Cancelled', 'ملغى', '#FF0000', 'cancelled', 'active'),
    ('Rescheduled', 'أعيد جدولتها', '#800080', 'rescheduled', 'active'),
    ('Pending Approval', 'في انتظار الموافقة', '#000080', 'pending_approval', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE meeting_statuses;
-- +goose StatementEnd
