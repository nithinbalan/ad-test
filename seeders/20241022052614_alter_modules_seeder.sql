-- +goose Up
-- +goose StatementBegin

DECLARE @BCM INT;
DECLARE @BCMPolicy INT;
DECLARE @settingModuleID INT;

SELECT @BCM = id FROM modules WHERE identifier = 'business_continuity';
SELECT @BCMPolicy = id FROM modules WHERE identifier = 'policies' and module_id=@BCM;
SELECT @settingModuleID = id FROM modules WHERE identifier = 'setting';

-- Update identifier for BCM Policies
UPDATE modules
SET identifier = 'bcm_policies'
WHERE id = @BCMPolicy;

-- Insert Master Setting Sub-Module under Settings
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('common_setting', 'Masters', 'الماجستير', @settingModuleID, 1, NULL, '/setting/master-setting');

--assign permission to new modules
-- Declarations for role and permission type IDs
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;

SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';

-- Assign permissions to all submodules for the Super Admin role
DECLARE @subModuleID INT;
DECLARE submodule_cursor CURSOR FOR 
    SELECT id FROM modules WHERE module_id IS NOT NULL;

OPEN submodule_cursor;
FETCH NEXT FROM submodule_cursor INTO @subModuleID;

WHILE @@FETCH_STATUS = 0
BEGIN
    -- Assign View Permission
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @subModuleID, @viewPermissionTypeID, 1, @superAdminRoleID, GETDATE());

    -- Assign Modify Permission
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @subModuleID, @modifyPermissionTypeID, 1, @superAdminRoleID, GETDATE());

    -- Assign Manage Permission
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @subModuleID, @managePermissionTypeID, 1, @superAdminRoleID, GETDATE());

    FETCH NEXT FROM submodule_cursor INTO @subModuleID;
END

CLOSE submodule_cursor;
DEALLOCATE submodule_cursor;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DECLARE @BCM INT;
DECLARE @BCMPolicy INT;
DECLARE @settingModuleID INT;

-- Get IDs again to make sure they're correctly referenced
SELECT @BCM = id FROM modules WHERE identifier = 'business_continuity';
SELECT @BCMPolicy = id FROM modules WHERE identifier = 'bcm_policies' AND module_id = @BCM;
SELECT @settingModuleID = id FROM modules WHERE identifier = 'setting';

-- Revert identifier for BCM Policies
UPDATE modules
SET identifier = 'policies'
WHERE id = @BCMPolicy;

-- Remove the Master Setting sub-module from Settings
DELETE FROM modules
WHERE identifier = 'common_setting';

-- Revoke permissions assigned to the new modules for the Super Admin role
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;

SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';

-- Remove permissions for all submodules
DECLARE @subModuleID INT;
DECLARE submodule_cursor CURSOR FOR 
    SELECT id FROM modules WHERE module_id IS NOT NULL;

OPEN submodule_cursor;
FETCH NEXT FROM submodule_cursor INTO @subModuleID;

WHILE @@FETCH_STATUS = 0
BEGIN
    -- Remove View Permission
    DELETE FROM role_permissions
    WHERE role_id = @superAdminRoleID AND module_id = @subModuleID AND role_permission_type_id = @viewPermissionTypeID;

    -- Remove Modify Permission
    DELETE FROM role_permissions
    WHERE role_id = @superAdminRoleID AND module_id = @subModuleID AND role_permission_type_id = @modifyPermissionTypeID;

    -- Remove Manage Permission
    DELETE FROM role_permissions
    WHERE role_id = @superAdminRoleID AND module_id = @subModuleID AND role_permission_type_id = @managePermissionTypeID;

    FETCH NEXT FROM submodule_cursor INTO @subModuleID;
END

CLOSE submodule_cursor;
DEALLOCATE submodule_cursor;

-- +goose StatementEnd
