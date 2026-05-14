-- +goose Up
-- +goose StatementBegin
DECLARE @assetModuleID INT;
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;
 
-- Get parent module ID
SELECT @assetModuleID = id FROM modules WHERE identifier = 'asset';
 
-- Insert new submodule: system_registe
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES
    ('system_register', 'System Register', 'سجل النظام', @assetModuleID, 6, NULL, '/asset/system-register');
 
-- Get role and permission type IDs
SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';
 
-- Get the ID of the new submodule
DECLARE @systemRegisterModuleID INT;
SELECT @systemRegisterModuleID = id FROM modules WHERE identifier = 'system_register';
 
-- Assign view, modify, and manage permissions to super admin
IF NOT EXISTS (
    SELECT 1 FROM role_permissions
    WHERE role_id = @superAdminRoleID
    AND module_id = @systemRegisterModuleID
    AND role_permission_type_id = @viewPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @systemRegisterModuleID, @viewPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END
 
IF NOT EXISTS (
    SELECT 1 FROM role_permissions
    WHERE role_id = @superAdminRoleID
    AND module_id = @systemRegisterModuleID
    AND role_permission_type_id = @modifyPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @systemRegisterModuleID, @modifyPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END
 
IF NOT EXISTS (
    SELECT 1 FROM role_permissions
    WHERE role_id = @superAdminRoleID
    AND module_id = @systemRegisterModuleID
    AND role_permission_type_id = @managePermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @systemRegisterModuleID, @managePermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

-- Update the order for existing submodules
UPDATE modules SET [order] = 7 WHERE identifier = 'asset_report';


-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
-- Remove permissions for 'system_register' submodule
DELETE FROM role_permissions
WHERE module_id IN (
    SELECT id FROM modules WHERE identifier = 'system_register'
);
 
-- Remove 'system_register' submodule
DELETE FROM modules
WHERE identifier = 'system_register';
-- revert module orders back
UPDATE modules SET [order] = 6 WHERE identifier = 'asset_report';
-- +goose StatementEnd