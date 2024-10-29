-- +goose Up
-- +goose StatementBegin
-- Drop the existing type constraint on workflow_statuses
ALTER TABLE workflow_statuses
DROP CONSTRAINT CHK_workflow_statuses_type;

-- Add the new constraint with updated ENUM values, including 'pending'
ALTER TABLE workflow_statuses
ADD CONSTRAINT CHK_workflow_statuses_type 
CHECK (type IN ('submit', 'revert', 'reject', 'approve', 'publish', 'pending'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Revert by dropping the modified constraint
ALTER TABLE workflow_statuses
DROP CONSTRAINT CHK_workflow_statuses_type;

-- Re-add the original constraint without 'pending'
ALTER TABLE workflow_statuses
ADD CONSTRAINT CHK_workflow_statuses_type 
CHECK (type IN ('submit', 'revert', 'reject', 'approve', 'publish'));
-- +goose StatementEnd
