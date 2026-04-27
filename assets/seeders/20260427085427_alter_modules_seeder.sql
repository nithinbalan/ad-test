-- +goose Up
-- +goose StatementBegin

DECLARE @ComplianceModuleID INT;

-- Get Compliance module ID
SELECT @ComplianceModuleID = id
FROM modules
WHERE identifier = 'compliance';

-- Shift existing compliance modules down by 1
UPDATE modules
SET [order] = [order] + 1
WHERE module_id = @ComplianceModuleID;

-- Update compliance_dashboard
UPDATE modules
SET 
    module_id = @ComplianceModuleID,
    route = '/compliance/dashboard',
    [order] = 1
WHERE identifier = 'compliance_dashboard';

-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin

DECLARE @DashboardModuleID INT;
DECLARE @ComplianceModuleID INT;

-- Get module IDs
SELECT @DashboardModuleID = id FROM modules WHERE identifier = 'dashboard';
SELECT @ComplianceModuleID = id FROM modules WHERE identifier = 'compliance';

-- Revert compliance_dashboard
UPDATE modules
SET 
    module_id = @DashboardModuleID,
    route = '/dashboards/compliance',
    [order] = 2
WHERE identifier = 'compliance_dashboard';

-- Shift compliance modules back (-1)
UPDATE modules
SET [order] = [order] - 1
WHERE module_id = @ComplianceModuleID;

-- +goose StatementEnd