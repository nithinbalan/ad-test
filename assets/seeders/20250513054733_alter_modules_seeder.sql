-- +goose Up
-- +goose StatementBegin
DECLARE @reportModuleID INT;
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;

UPDATE modules
SET [order] = 15
WHERE identifier = 'setting';

-- Create the main Report module
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('report', 'Report', 'تقرير', NULL, 14,  'icon-report.svg', '/report');

-- Fetch the newly inserted Report module ID
SET @reportModuleID = (SELECT id FROM modules WHERE identifier = 'report');

-- Create Submodules
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('grc_report', 'GRC Report', 'تقرير مركز الخليج للأبحاث', @reportModuleID, 1, NULL, '/report/grc-report');

-- Get IDs for roles and permission types
SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';

-- Assign permissions to all Report submodules
DECLARE @subModuleID INT;
DECLARE report_cursor CURSOR FOR 
    SELECT id FROM modules WHERE module_id = @reportModuleID;

OPEN report_cursor;
FETCH NEXT FROM report_cursor INTO @subModuleID;

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

    FETCH NEXT FROM report_cursor INTO @subModuleID;
END

CLOSE report_cursor;
DEALLOCATE report_cursor;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Remove permissions assigned to Report submodules
DELETE FROM role_permissions
WHERE module_id IN (
    SELECT id FROM modules WHERE identifier IN ('grc_report')
);

-- Remove Report submodules
DELETE FROM modules
WHERE identifier IN ('grc_report');

-- Remove Report main module
DELETE FROM modules
WHERE identifier = 'report';

-- +goose StatementEnd
