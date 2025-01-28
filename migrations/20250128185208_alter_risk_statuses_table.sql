-- +goose Up
-- +goose StatementBegin
-- Drop the existing CHECK constraint
ALTER TABLE risk_statuses
DROP CONSTRAINT CHK_risk_statuses_type;

-- Add a new CHECK constraint with the 'identified' status included
ALTER TABLE risk_statuses
ADD CONSTRAINT CHK_risk_statuses_type 
CHECK (type IN ('open', 'identified', 'assessed', 'mitigated', 'accepted', 'closed', 'pending'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop the new CHECK constraint
ALTER TABLE risk_statuses
DROP CONSTRAINT CHK_risk_statuses_type;

-- Add the original CHECK constraint without 'identified'
ALTER TABLE risk_statuses
ADD CONSTRAINT CHK_risk_statuses_type 
CHECK (type IN ('open', 'assessed', 'mitigated', 'accepted', 'closed', 'pending'));
-- +goose StatementEnd