-- +goose Up
-- +goose StatementBegin
DECLARE @dashboardModuleID INT;
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;

-- Get parent module ID
SELECT @dashboardModuleID = id FROM modules WHERE identifier = 'dashboard';

-- Insert new submodule: kri_dashboard
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES
    ('kri_dashboard', 'KRI Dashboard', 'لوحة مؤشرات المخاطر الرئيسية', @dashboardModuleID, 4, NULL, 'dashboards/kri');

-- Get role and permission type IDs
SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';

-- Get the ID of the new submodule
DECLARE @kriDashboardModuleID INT;
SELECT @kriDashboardModuleID = id FROM modules WHERE identifier = 'kri_dashboard';

-- Assign view permission
IF NOT EXISTS (
    SELECT 1 FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @kriDashboardModuleID
      AND role_permission_type_id = @viewPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @kriDashboardModuleID, @viewPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

-- Assign modify permission
IF NOT EXISTS (
    SELECT 1 FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @kriDashboardModuleID
      AND role_permission_type_id = @modifyPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @kriDashboardModuleID, @modifyPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

-- Assign manage permission
IF NOT EXISTS (
    SELECT 1 FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @kriDashboardModuleID
      AND role_permission_type_id = @managePermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @kriDashboardModuleID, @managePermissionTypeID, 1, @superAdminRoleID, GETDATE());
END
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
-- Remove permissions for 'kri_dashboard' submodule
DELETE FROM role_permissions
WHERE module_id IN (SELECT id FROM modules WHERE identifier = 'kri_dashboard');

-- Remove 'kri_dashboard' module
DELETE FROM modules WHERE identifier = 'kri_dashboard';
-- +goose StatementEnd
