-- +goose Up
-- +goose StatementBegin
-- Add the new audit_final_report_id column as INT 
ALTER TABLE tasks
ADD audit_final_report_id INT NULL;

ALTER TABLE tasks
ADD audit_report_workflow_id INT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the audit_final_report_id column
ALTER TABLE tasks
DROP COLUMN audit_final_report_id;

ALTER TABLE tasks
DROP COLUMN audit_report_workflow_id;
-- +goose StatementEnd