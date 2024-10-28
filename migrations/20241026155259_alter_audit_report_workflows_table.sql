-- +goose Up
-- +goose StatementBegin
-- Add the audit_final_report_id column with a foreign key to audit_final_reports table
ALTER TABLE audit_report_workflows
ADD audit_final_report_id INT NULL;

-- Add the foreign key constraint for audit_final_report_id to audit_final_reports
ALTER TABLE audit_report_workflows
ADD CONSTRAINT FK_audit_report_workflows_final_reports
FOREIGN KEY (audit_final_report_id) REFERENCES audit_final_reports(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the foreign key constraint first
ALTER TABLE audit_report_workflows
DROP CONSTRAINT FK_audit_report_workflows_final_reports;

-- Drop the audit_final_report_id column
ALTER TABLE audit_report_workflows
DROP COLUMN audit_final_report_id;
-- +goose StatementEnd
