-- +goose Up
-- +goose StatementBegin
-- Drop the existing type constraint on audit_final_report_statuses
ALTER TABLE audit_final_report_statuses
DROP CONSTRAINT CHK_audit_final_report_statuses_type;

-- Add the new constraint with updated ENUM values, including 'approve'
ALTER TABLE audit_final_report_statuses
ADD CONSTRAINT CHK_audit_final_report_statuses_type 
CHECK (type IN ('draft','in-review', 'submit', 'revert','reject', 'published', 'completed','approve'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Revert by dropping the modified constraint
ALTER TABLE audit_final_report_statuses
DROP CONSTRAINT CHK_audit_final_report_statuses_type;

-- Re-add the original constraint without 'pending'
ALTER TABLE audit_final_report_statuses
ADD CONSTRAINT CHK_audit_final_report_statuses_type 
CHECK (type IN ('draft','in-review', 'submit', 'revert','reject', 'published', 'completed'));
-- +goose StatementEnd
