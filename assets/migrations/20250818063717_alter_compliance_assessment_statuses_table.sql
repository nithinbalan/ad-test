-- +goose Up
-- +goose StatementBegin
-- Drop the old constraint on the 'type' field (if needed)
ALTER TABLE compliance_assessment_statuses
DROP CONSTRAINT CHK_compliance_assessment_statuses_type;

-- Modify the 'type' field to accept both the old and new values
ALTER TABLE compliance_assessment_statuses
ADD CONSTRAINT CHK_compliance_assessment_statuses_type
CHECK (type IN ('draft', 'in_review', 'approved','in_progress'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop the modified constraint on the 'type' field
ALTER TABLE compliance_assessment_statuses
DROP CONSTRAINT CHK_compliance_assessment_statuses_type;

-- Restore the original 'type' constraint with old values
ALTER TABLE compliance_assessment_statuses
ADD CONSTRAINT CHK_compliance_assessment_statuses_type
CHECK (type IN ('draft', 'in_review', 'approved'));
-- +goose StatementEnd
