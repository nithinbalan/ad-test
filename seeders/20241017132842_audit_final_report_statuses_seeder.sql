-- +goose Up
-- +goose StatementBegin
INSERT INTO audit_final_report_statuses
    (title_en, title_ar, type, color_code, status)
VALUES
    ('Draft', 'مسودة', 'draft', '#FFC107', 'active'),
    ('In Review', 'قيد المراجعة', 'in-review', '#90DD75', 'active'),
    ('Submit', 'إرسال', 'submit', '#17A2B8', 'active'),
    ('Revert', 'إرجاع', 'revert', '#FD7E14', 'active'),
    ('Reject', 'رفض', 'reject', '#DC3545', 'active'),
    ('Published', 'نشر', 'published', '#28A745', 'active'),
    ('Completed', 'مكتمل', 'completed', '#DC3545', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE audit_final_report_statuses;
-- +goose StatementEnd
