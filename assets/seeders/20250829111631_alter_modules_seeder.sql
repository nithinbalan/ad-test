-- +goose Up
-- +goose StatementBegin
DECLARE @policyManagementModuleID INT;
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;

-- Create the main Policies module
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('policy_management', 'Policy', 'سياسة', NULL, 16, 'icon-policy.svg', '/policies');

-- Fetch the newly inserted Policies module ID
SET @policyManagementModuleID = (SELECT id FROM modules WHERE identifier = 'policy_management');

-- Create Policy Register submodule
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('policy_register', 'Policy Register', 'المكتبة', @policyManagementModuleID, 1, NULL, '/policies/policy-register');

-- Get IDs for roles and permission types
SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';

-- Assign permissions to all Policies submodules
DECLARE @subModuleID INT;
DECLARE policy_cursor CURSOR FOR 
    SELECT id FROM modules WHERE module_id = @policyManagementModuleID;

OPEN policy_cursor;
FETCH NEXT FROM policy_cursor INTO @subModuleID;

WHILE @@FETCH_STATUS = 0
BEGIN
    -- View Permission
    IF NOT EXISTS (
        SELECT 1 FROM role_permissions 
        WHERE role_id = @superAdminRoleID 
          AND module_id = @subModuleID 
          AND role_permission_type_id = @viewPermissionTypeID
    )
    BEGIN
        INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
        VALUES (@superAdminRoleID, @subModuleID, @viewPermissionTypeID, 1, @superAdminRoleID, GETDATE());
    END

    -- Modify Permission
    IF NOT EXISTS (
        SELECT 1 FROM role_permissions 
        WHERE role_id = @superAdminRoleID 
          AND module_id = @subModuleID 
          AND role_permission_type_id = @modifyPermissionTypeID
    )
    BEGIN
        INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
        VALUES (@superAdminRoleID, @subModuleID, @modifyPermissionTypeID, 1, @superAdminRoleID, GETDATE());
    END

    -- Manage Permission
    IF NOT EXISTS (
        SELECT 1 FROM role_permissions 
        WHERE role_id = @superAdminRoleID 
          AND module_id = @subModuleID 
          AND role_permission_type_id = @managePermissionTypeID
    )
    BEGIN
        INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
        VALUES (@superAdminRoleID, @subModuleID, @managePermissionTypeID, 1, @superAdminRoleID, GETDATE());
    END

    FETCH NEXT FROM policy_cursor INTO @subModuleID;
END

CLOSE policy_cursor;
DEALLOCATE policy_cursor;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove permissions for Policy Register submodule
DELETE FROM role_permissions
WHERE module_id IN (
    SELECT id FROM modules WHERE identifier IN ('policy_register')
);

-- Remove Policy Register submodule
DELETE FROM modules
WHERE identifier IN ('policy_register');

-- Remove Policies main module
DELETE FROM modules
WHERE identifier = 'policy_management';
-- +goose StatementEnd
