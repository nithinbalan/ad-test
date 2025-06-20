-- +goose Up
-- +goose StatementBegin
DECLARE @settingModuleID INT;
SET @settingModuleID = (SELECT id FROM modules WHERE identifier = 'setting');

UPDATE modules
SET [order] = 13
WHERE identifier = 'log';

UPDATE modules
SET [order] = 14
WHERE identifier = 'common_setting';

INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('strategy_setting', 'Strategy', 'الإستراتيجية', @settingModuleID, 12, NULL, '/setting/strategy-setting');

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
    -- Ensure View Permission is not duplicated
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

    -- Ensure Modify Permission is not duplicated
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

    -- Ensure Manage Permission is not duplicated
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

    FETCH NEXT FROM submodule_cursor INTO @subModuleID;
END

CLOSE submodule_cursor;
DEALLOCATE submodule_cursor;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Delete permissions added for the new submodule
DELETE FROM role_permissions
WHERE module_id IN (
    SELECT id FROM modules WHERE identifier = 'strategy_setting'
);

-- Delete the new submodule from the modules table
DELETE FROM modules
WHERE identifier = 'strategy_setting';

-- +goose StatementEnd

