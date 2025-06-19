-- +goose Up
-- +goose StatementBegin
DECLARE @settingModuleID INT;
SET @settingModuleID = (SELECT id FROM modules WHERE identifier = 'setting');

-- Update existing module orders
UPDATE modules
SET [order] = 14
WHERE identifier = 'log';

UPDATE modules
SET [order] = 15
WHERE identifier = 'common_setting';

-- Insert new TPRM Setting submodule
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('tprm_setting', 'TPRM', 'إدارة مخاطر الطرف الثالث', @settingModuleID, 13, NULL, '/setting/tprm-setting');

-- Declarations for role and permission type IDs
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;
DECLARE @newSubModuleID INT;

-- Get Super Admin role and permission type IDs
SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';

-- Get the new inserted module ID
SELECT @newSubModuleID = id FROM modules WHERE identifier = 'tprm_setting';

-- Assign View Permission
IF NOT EXISTS (
    SELECT 1 FROM role_permissions 
    WHERE role_id = @superAdminRoleID 
      AND module_id = @newSubModuleID 
      AND role_permission_type_id = @viewPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @newSubModuleID, @viewPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

-- Assign Modify Permission
IF NOT EXISTS (
    SELECT 1 FROM role_permissions 
    WHERE role_id = @superAdminRoleID 
      AND module_id = @newSubModuleID 
      AND role_permission_type_id = @modifyPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @newSubModuleID, @modifyPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

-- Assign Manage Permission
IF NOT EXISTS (
    SELECT 1 FROM role_permissions 
    WHERE role_id = @superAdminRoleID 
      AND module_id = @newSubModuleID 
      AND role_permission_type_id = @managePermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @newSubModuleID, @managePermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Remove permissions for TPRM Setting
DELETE FROM role_permissions
WHERE module_id IN (
    SELECT id FROM modules WHERE identifier = 'tprm_setting'
);

-- Remove TPRM Setting module
DELETE FROM modules
WHERE identifier = 'tprm_setting';

-- Revert module orders back
UPDATE modules
SET [order] = 13
WHERE identifier = 'log';

UPDATE modules
SET [order] = 14
WHERE identifier = 'common_setting';

-- +goose StatementEnd
