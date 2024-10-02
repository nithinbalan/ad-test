-- +goose Up
-- +goose StatementBegin
-- First, drop the existing constraint
SELECT * FROM audit_statuses
WHERE type NOT IN ('draft',  'published', 'in-progress', 'completed');

UPDATE audit_statuses
SET type = 'in-progress'
WHERE type NOT IN ('draft',  'published', 'in-progress', 'completed');

-- First, drop the existing constraint
ALTER TABLE audit_statuses
DROP CONSTRAINT CHK_audit_statuses_type;

-- Update the ENUM values for type
ALTER TABLE audit_statuses
ADD CONSTRAINT CHK_audit_statuses_type 
CHECK (type IN ('draft', 'published', 'in-progress', 'completed'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Revert the changes by dropping the updated constraint
ALTER TABLE audit_statuses
DROP CONSTRAINT CHK_audit_statuses_type;

-- Re-add the original constraint with the old values
ALTER TABLE audit_statuses
ADD CONSTRAINT CHK_audit_statuses_type 
CHECK (type IN ('pending', 'planned', 'in-progress', 'completed'));
-- +goose StatementEnd

