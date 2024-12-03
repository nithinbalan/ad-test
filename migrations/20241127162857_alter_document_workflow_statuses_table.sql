-- +goose Up
-- +goose StatementBegin
-- Drop the existing constraint
ALTER TABLE document_workflow_statuses
DROP CONSTRAINT CHK_document_workflow_statuses_type;

-- Add the updated constraint with the new "resubmit" status
ALTER TABLE document_workflow_statuses
ADD CONSTRAINT CHK_document_workflow_statuses_type CHECK (type IN ('comment', 'revert', 'reject', 'approve', 'pending', 'resubmit'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop the updated constraint
ALTER TABLE document_workflow_statuses
DROP CONSTRAINT CHK_document_workflow_statuses_type;

-- Re-add the original constraint without "resubmit"
ALTER TABLE document_workflow_statuses
ADD CONSTRAINT CHK_document_workflow_statuses_type CHECK (type IN ('comment', 'revert', 'reject', 'approve', 'pending'));
-- +goose StatementEnd
