-- +goose Up
-- +goose StatementBegin

-- Document module mapping
DECLARE @documentModuleID INT;
SELECT @documentModuleID = id FROM modules WHERE identifier = 'document';

UPDATE reminder_contexts
SET module_id = @documentModuleID
WHERE type = 'policy-acknowledgement-due';

-- Compliance module mapping
DECLARE @complianceModuleID INT;
SELECT @complianceModuleID = id FROM modules WHERE identifier = 'compliance';

UPDATE reminder_contexts
SET module_id = @complianceModuleID
WHERE type = 'compliance-review-due';

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

UPDATE reminder_contexts
SET module_id = NULL
WHERE type IN ('policy-acknowledgement-due', 'compliance-review-due');

-- +goose StatementEnd
