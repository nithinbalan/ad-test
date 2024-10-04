-- +goose Up
-- +goose StatementBegin
-- First, ensure there are no entries with invalid types
SELECT * FROM audit_draft_finding_conversations
WHERE type NOT IN ('DR', 'AC', 'IR');

-- Optionally update any existing invalid types (if necessary)
-- UPDATE audit_draft_finding_conversations
-- SET type = 'DR'  -- or appropriate value
-- WHERE type NOT IN ('DR', 'AC', 'IR');

-- Drop the existing constraint
ALTER TABLE audit_draft_finding_conversations
DROP CONSTRAINT CHK_audit_draft_finding_conversations_type;

-- Add the new constraint with updated ENUM values including 'IR'
ALTER TABLE audit_draft_finding_conversations
ADD CONSTRAINT CHK_audit_draft_finding_conversations_type 
CHECK (type IN ('DR', 'AC', 'IR'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Revert the change by dropping the updated constraint
ALTER TABLE audit_draft_finding_conversations
DROP CONSTRAINT CHK_audit_draft_finding_conversations_type;

-- Re-add the original constraint without 'IR'
ALTER TABLE audit_draft_finding_conversations
ADD CONSTRAINT CHK_audit_draft_finding_conversations_type 
CHECK (type IN ('DR', 'AC'));
-- +goose StatementEnd
