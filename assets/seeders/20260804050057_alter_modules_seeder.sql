-- +goose Up
-- +goose StatementBegin
 
DECLARE @documentModuleID INT;
DECLARE @documentDashboardModuleID INT;
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;
 
-- Get parent module ID
SELECT @documentModuleID = id
FROM modules
WHERE identifier = 'document';
 
-- Shift all existing Document submodules down by one
UPDATE modules
SET [order] = [order] + 1
WHERE module_id = @documentModuleID;
 
-- Insert new Dashboard submodule
INSERT INTO modules (
    identifier,
    title_en,
    title_ar,
    module_id,
    [order],
    icon,
    route
)
VALUES (
    'document_dashboard',
    'Dashboard',
    N'لوحة المعلومات',
    @documentModuleID,
    1,
    NULL,
    '/document/dashboard'
);
 
-- Get role and permission type IDs
SELECT @superAdminRoleID = id
FROM roles
WHERE is_super_admin = 1;
 
SELECT @viewPermissionTypeID = id
FROM role_permission_types
WHERE type = 'view';
 
SELECT @modifyPermissionTypeID = id
FROM role_permission_types
WHERE type = 'modify';
 
SELECT @managePermissionTypeID = id
FROM role_permission_types
WHERE type = 'manage';
 
-- Get Dashboard module ID
SELECT @documentDashboardModuleID = id
FROM modules
WHERE identifier = 'document_dashboard';
 
-- Assign View permission
IF NOT EXISTS (
    SELECT 1
    FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @documentDashboardModuleID
      AND role_permission_type_id = @viewPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (
        role_id,
        module_id,
        role_permission_type_id,
        is_granted,
        created_by,
        created_at
    )
    VALUES (
        @superAdminRoleID,
        @documentDashboardModuleID,
        @viewPermissionTypeID,
        1,
        @superAdminRoleID,
        GETDATE()
    );
END
 
-- Assign Modify permission
IF NOT EXISTS (
    SELECT 1
    FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @documentDashboardModuleID
      AND role_permission_type_id = @modifyPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (
        role_id,
        module_id,
        role_permission_type_id,
        is_granted,
        created_by,
        created_at
    )
    VALUES (
        @superAdminRoleID,
        @documentDashboardModuleID,
        @modifyPermissionTypeID,
        1,
        @superAdminRoleID,
        GETDATE()
    );
END
 
-- Assign Manage permission
IF NOT EXISTS (
    SELECT 1
    FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @documentDashboardModuleID
      AND role_permission_type_id = @managePermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (
        role_id,
        module_id,
        role_permission_type_id,
        is_granted,
        created_by,
        created_at
    )
    VALUES (
        @superAdminRoleID,
        @documentDashboardModuleID,
        @managePermissionTypeID,
        1,
        @superAdminRoleID,
        GETDATE()
    );
END
 
-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
 
DECLARE @documentModuleID INT;
 
SELECT @documentModuleID = id
FROM modules
WHERE identifier = 'document';
 
-- Remove Dashboard permissions
DELETE FROM role_permissions
WHERE module_id IN (
    SELECT id
    FROM modules
    WHERE identifier = 'document_dashboard'
);
 
-- Remove Dashboard module
DELETE FROM modules
WHERE identifier = 'document_dashboard';
 
-- Restore original ordering
UPDATE modules
SET [order] = [order] - 1
WHERE module_id = @documentModuleID;
 
-- +goose StatementEnd