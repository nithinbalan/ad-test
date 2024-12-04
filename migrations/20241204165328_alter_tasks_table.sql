-- +goose Up
-- +goose StatementBegin
-- Add the new audit_draft_report_id column as INT 
ALTER TABLE tasks
ADD audit_draft_report_id INT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the audit_draft_report_id column
ALTER TABLE tasks
DROP COLUMN audit_draft_report_id;
-- +goose StatementEnd