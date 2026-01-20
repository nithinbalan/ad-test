-- +goose Up
-- +goose StatementBegin
INSERT INTO violation_report_statuses (title_en, title_ar, color_code, type, status) VALUES
    ('New', 'جديد', '#28A745', 'new', 'active'),
    ('Under Investigation', 'قيد التحقيق', '#17A2B8', 'under_investigation', 'active'),
    ('Needs Clarification', 'بحاجة إلى توضيح', '#FFC107', 'needs_clarification', 'active'),
    ('In Progress', 'قيد التنفيذ', '#007BFF', 'in_progress', 'active'),
    ('Under Review', 'قيد المراجعة', '#17A2B8', 'under_review', 'active'),
    ('On Hold', 'معلق', '#FFA500', 'on_hold', 'active'),
    ('Closed - Rejected', 'مغلق - مرفوض', '#DC3545', 'closed_rejected', 'active'),
    ('Closed - Resolved', 'مغلق - تم الحل', '#6C757D', 'closed_resolved', 'active'),
    ('Closed - Irrelevant', 'مغلق - غير ذي صلة', '#343A40', 'closed_irrelevant', 'active'),
    ('Closed - No Violation Found', 'مغلق - لا يوجد انتهاك', '#6F42C1', 'closed_no_violation_found', 'active'),
    ('Closed - Referred to Entity', 'مغلق - تمت الإحالة إلى الجهة المختصة', '#FD7E14', 'closed_referred_to_entity', 'active');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM violation_report_statuses
WHERE type IN (
    'new',
    'under_investigation',
    'needs_clarification',
    'in_progress',
    'under_review',
    'on_hold',
    'closed_rejected',
    'closed_resolved',
    'closed_irrelevant',
    'closed_no_violation_found',
    'closed_referred_to_entity'
);
-- +goose StatementEnd
