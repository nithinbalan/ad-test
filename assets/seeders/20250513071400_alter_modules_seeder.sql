-- +goose Up
-- +goose StatementBegin
DECLARE @knowledgeCenterModuleID INT;
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;

-- Create the main Knowledge Center module
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('knowledge_center', 'Knowledge Center', 'مركز المعرفة', NULL, 16, 'icon-knowledge.svg', '/knowledge-center');

-- Fetch the newly inserted Knowledge Center module ID
SET @knowledgeCenterModuleID = (SELECT id FROM modules WHERE identifier = 'knowledge_center');

-- Create Library submodule
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('library', 'Library', 'المكتبة', @knowledgeCenterModuleID, 1, NULL, '/knowledge-center/library');

-- Get IDs for roles and permission types
SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';

-- Assign permissions to all Knowledge Center submodules
DECLARE @subModuleID INT;
DECLARE kc_cursor CURSOR FOR 
    SELECT id FROM modules WHERE module_id = @knowledgeCenterModuleID;

OPEN kc_cursor;
FETCH NEXT FROM kc_cursor INTO @subModuleID;

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

    FETCH NEXT FROM kc_cursor INTO @subModuleID;
END

CLOSE kc_cursor;
DEALLOCATE kc_cursor;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove permissions for Library submodule
DELETE FROM role_permissions
WHERE module_id IN (
    SELECT id FROM modules WHERE identifier IN ('library')
);

-- Remove Library submodule
DELETE FROM modules
WHERE identifier IN ('library');

-- Remove Knowledge Center main module
DELETE FROM modules
WHERE identifier = 'knowledge_center';
-- +goose StatementEnd
