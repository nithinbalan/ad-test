-- +goose Up
-- +goose StatementBegin
DECLARE @riskModuleID INT;
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;
 
-- Get parent module ID
SELECT @riskModuleID = id FROM modules WHERE identifier = 'risk';
 
-- Insert new submodule: risk_register_update_plan
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES
    ('risk_register_update_plan', 'Risk Register Update Plan', 'خطة تحديث سجل المخاطر', @riskModuleID, 9, NULL, '/risk/risk-register-update-plan');
 
-- Get role and permission type IDs
SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';
 
-- Get the ID of the new submodule
DECLARE @riskRegisterUpdatePlanModuleID INT;
SELECT @riskRegisterUpdatePlanModuleID = id FROM modules WHERE identifier = 'risk_register_update_plan';
 
-- Assign view, modify, and manage permissions to super admin
IF NOT EXISTS (
    SELECT 1 FROM role_permissions
    WHERE role_id = @superAdminRoleID
    AND module_id = @riskRegisterUpdatePlanModuleID
    AND role_permission_type_id = @viewPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @riskRegisterUpdatePlanModuleID, @viewPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END
 
IF NOT EXISTS (
    SELECT 1 FROM role_permissions
    WHERE role_id = @superAdminRoleID
    AND module_id = @riskRegisterUpdatePlanModuleID
    AND role_permission_type_id = @modifyPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @riskRegisterUpdatePlanModuleID, @modifyPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END
 
IF NOT EXISTS (
    SELECT 1 FROM role_permissions
    WHERE role_id = @superAdminRoleID
    AND module_id = @riskRegisterUpdatePlanModuleID
    AND role_permission_type_id = @managePermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @riskRegisterUpdatePlanModuleID, @managePermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

-- Update the order for existing submodules
UPDATE modules SET [order] = 10 WHERE identifier = 'risk_report';


-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
-- Remove permissions for 'risk_register_update_plan' submodule
DELETE FROM role_permissions
WHERE module_id IN (
    SELECT id FROM modules WHERE identifier = 'risk_register_update_plan'
);
 
-- Remove 'risk_register_update_plan' submodule
DELETE FROM modules
WHERE identifier = 'risk_register_update_plan';
-- revert module orders back
UPDATE modules SET [order] = 9 WHERE identifier = 'risk_report';
-- +goose StatementEnd