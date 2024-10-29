-- +goose Up
-- +goose StatementBegin

-- Insert the new 'approve' status into the audit_final_report_statuses table
INSERT INTO audit_final_report_statuses (title_en, title_ar, color_code, type, status)
VALUES ('Approved', 'حالة الموافقة', '#90DD75', 'approve', 'active');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Remove the 'approve' status from the audit_final_report_statuses table in case of rollback
DELETE FROM audit_final_report_statuses
WHERE type = 'approve';

-- +goose StatementEnd
