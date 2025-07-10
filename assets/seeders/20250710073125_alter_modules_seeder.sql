-- +goose Up
-- +goose StatementBegin

DECLARE @strategyModuleID INT;
DECLARE @strategyIntiativeModuleID INT;
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;

-- Get the strategy module ID
SELECT @strategyModuleID = id FROM modules WHERE identifier = 'strategy';

-- Insert Strategy Intiative module under Strategy
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES ('strategy_intiative', 'Strategy Intiative', 'المشاريع', @strategyModuleID, 5, NULL, '/strategy/strategy-intiative');

-- Ensure KPI is order 6
UPDATE modules SET [order] = 6 WHERE identifier = 'kpi' AND module_id = @strategyModuleID;

-- Get the newly inserted Strategy Intiative module ID
SET @strategyIntiativeModuleID = (SELECT id FROM modules WHERE identifier = 'strategy_intiative');

-- Get IDs for roles and permission types
SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';

-- Assign permissions to Strategy Intiative module
IF NOT EXISTS (
    SELECT 1 FROM role_permissions 
    WHERE role_id = @superAdminRoleID 
      AND module_id = @strategyIntiativeModuleID 
      AND role_permission_type_id = @viewPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @strategyIntiativeModuleID, @viewPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

IF NOT EXISTS (
    SELECT 1 FROM role_permissions 
    WHERE role_id = @superAdminRoleID 
      AND module_id = @strategyIntiativeModuleID 
      AND role_permission_type_id = @modifyPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @strategyIntiativeModuleID, @modifyPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

IF NOT EXISTS (
    SELECT 1 FROM role_permissions 
    WHERE role_id = @superAdminRoleID 
      AND module_id = @strategyIntiativeModuleID 
      AND role_permission_type_id = @managePermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @strategyIntiativeModuleID, @managePermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Remove permissions for Strategy Intiative module
DELETE FROM role_permissions
WHERE module_id = (SELECT id FROM modules WHERE identifier = 'strategy_intiative');

-- Remove Strategy Intiative module
DELETE FROM modules WHERE identifier = 'strategy_intiative';

-- Restore KPI order to 5
UPDATE modules SET [order] = 5 WHERE identifier = 'kpi' AND module_id = (SELECT id FROM modules WHERE identifier = 'strategy');

-- +goose StatementEnd
