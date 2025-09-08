-- +goose Up
-- +goose StatementBegin

-- Update the order for the Compliance Report submodule to 6
DECLARE @complianceModuleID INT;
SET @complianceModuleID = (SELECT id FROM modules WHERE identifier = 'compliance');

UPDATE modules
SET [order] = 6
WHERE identifier = 'compliance_report'
  AND module_id = @complianceModuleID;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- (Optional) Revert it back to 4 if needed
DECLARE @complianceModuleID INT;
SET @complianceModuleID = (SELECT id FROM modules WHERE identifier = 'compliance');

UPDATE modules
SET [order] = 5
WHERE identifier = 'compliance_report'
  AND module_id = @complianceModuleID;

-- +goose StatementEnd
