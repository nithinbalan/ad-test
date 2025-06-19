-- +goose Up
-- +goose StatementBegin

-- Fetch Compliance module ID
DECLARE @complianceModuleID INT;
SET @complianceModuleID = (SELECT id FROM modules WHERE identifier = 'compliance');

-- Insert new submodule: Policy Exception
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('policy_exception', 'Policy Exceptions', 'استثناء السياسات', @complianceModuleID, 5, NULL, '/compliance/policy-exception');

-- Declarations for role and permission type IDs
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;
DECLARE @policyExceptionModuleID INT;

SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';

SET @policyExceptionModuleID = (SELECT id FROM modules WHERE identifier = 'policy_exception');

-- View Permission
IF NOT EXISTS (
    SELECT 1 FROM role_permissions 
    WHERE role_id = @superAdminRoleID 
      AND module_id = @policyExceptionModuleID 
      AND role_permission_type_id = @viewPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @policyExceptionModuleID, @viewPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

-- Modify Permission
IF NOT EXISTS (
    SELECT 1 FROM role_permissions 
    WHERE role_id = @superAdminRoleID 
      AND module_id = @policyExceptionModuleID 
      AND role_permission_type_id = @modifyPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @policyExceptionModuleID, @modifyPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

-- Manage Permission
IF NOT EXISTS (
    SELECT 1 FROM role_permissions 
    WHERE role_id = @superAdminRoleID 
      AND module_id = @policyExceptionModuleID 
      AND role_permission_type_id = @managePermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @policyExceptionModuleID, @managePermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Remove role permissions for the Policy Exception submodule
DELETE FROM role_permissions
WHERE module_id IN (
    SELECT id FROM modules WHERE identifier = 'policy_exception'
);

-- Remove Policy Exception module
DELETE FROM modules
WHERE identifier = 'policy_exception';

-- +goose StatementEnd
