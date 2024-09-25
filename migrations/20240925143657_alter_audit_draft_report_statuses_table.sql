-- +goose Up
-- +goose StatementBegin
-- First, drop the existing constraint
SELECT * FROM audit_draft_report_statuses
WHERE type NOT IN ('in-review', 'approved', 'in-department-review', 'completed');

UPDATE audit_draft_report_statuses
SET type = 'draft'
WHERE type NOT IN ('draft', 'in-review', 'approved', 'in-department-review', 'completed');


ALTER TABLE audit_draft_report_statuses
DROP CONSTRAINT CHK_audit_draft_report_statuses_type;

-- Add the new constraint with updated ENUM values
ALTER TABLE audit_draft_report_statuses
ADD CONSTRAINT CHK_audit_draft_report_statuses_type 
CHECK (type IN ('draft', 'in-review', 'approved', 'in-department-review', 'completed'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Revert the change by dropping the updated constraint
ALTER TABLE audit_draft_report_statuses
DROP CONSTRAINT CHK_audit_draft_report_statuses_type;

-- Re-add the original constraint
ALTER TABLE audit_draft_report_statuses
ADD CONSTRAINT CHK_audit_draft_report_statuses_type 
CHECK (type IN ('draft', 'auditee-review', 'finalized'));
-- +goose StatementEnd
