-- +goose Up
-- Seeder for audit_draft_report_statuses

INSERT INTO audit_draft_report_statuses (title_en, title_ar, type, color_code, status)
VALUES 
    ('Draft', 'مسودة', 'draft', '#FFC107', 'active'),                   -- Yellow
    ('In Review', 'قيد المراجعة', 'in-review', '#90DD75', 'active'),     -- Light green
    ('Approved', 'موافَق عليه', 'approved', '#28A745', 'active'),       --  Green
    ('In Department Review', 'مراجعة القسم', 'in-department-review', '#17A2B8', 'active'), -- Light Blue
    ('Completed', 'مكتمل', 'completed', '#DC3545', 'active');           -- Red

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE audit_draft_report_statuses;
-- +goose StatementEnd
