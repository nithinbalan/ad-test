-- +goose Up
-- Seeder for audit_draft_report_statuses

INSERT INTO audit_draft_report_statuses (title_en, title_ar, type, color_code, status)
VALUES 
    ('Draft', 'مسودة', 'draft', '#28A745', 'active'),
    ('In Auditee Review', 'مراجعة المدقق', 'auditee-review', '#FFC107', 'active'),
    ('Finalized', 'مكتمل', 'finalized', '#DC3545', 'active');

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE audit_draft_report_statuses;
-- +goose StatementEnd