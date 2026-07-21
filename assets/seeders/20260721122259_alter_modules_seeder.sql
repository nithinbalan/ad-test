-- +goose Up
-- +goose StatementBegin
 
DECLARE @assetModuleID INT;
DECLARE @assetDashboardModuleID INT;
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;
 
-- Get parent module ID
SELECT @assetModuleID = id
FROM modules
WHERE identifier = 'asset';
 
-- Shift all existing Asset submodules down by one
UPDATE modules
SET [order] = [order] + 1
WHERE module_id = @assetModuleID;
 
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
    'asset_dashboard',
    'Dashboard',
    N'لوحة المعلومات',
    @assetModuleID,
    1,
    NULL,
    '/asset/dashboard'
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
SELECT @assetDashboardModuleID = id
FROM modules
WHERE identifier = 'asset_dashboard';
 
-- Assign View permission
IF NOT EXISTS (
    SELECT 1
    FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @assetDashboardModuleID
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
        @assetDashboardModuleID,
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
      AND module_id = @assetDashboardModuleID
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
        @assetDashboardModuleID,
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
      AND module_id = @assetDashboardModuleID
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
        @assetDashboardModuleID,
        @managePermissionTypeID,
        1,
        @superAdminRoleID,
        GETDATE()
    );
END
 
-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
 
DECLARE @assetModuleID INT;
 
SELECT @assetModuleID = id
FROM modules
WHERE identifier = 'asset';
 
-- Remove Dashboard permissions
DELETE FROM role_permissions
WHERE module_id IN (
    SELECT id
    FROM modules
    WHERE identifier = 'asset_dashboard'
);
 
-- Remove Dashboard module
DELETE FROM modules
WHERE identifier = 'asset_dashboard';
 
-- Restore original ordering
UPDATE modules
SET [order] = [order] - 1
WHERE module_id = @assetModuleID;
 
-- +goose StatementEnd