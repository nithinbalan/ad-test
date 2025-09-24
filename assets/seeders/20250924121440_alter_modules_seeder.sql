-- +goose Up
-- +goose StatementBegin
DECLARE @dashboardModuleID INT;
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;
 
-- Get parent module ID
SELECT @dashboardModuleID = id FROM modules WHERE identifier = 'dashboard';
 
-- Insert new submodule: compliance_calender
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES
    ('erm_dashboard', 'ERM Dashboard', 'ERM لوحة القيادة', @dashboardModuleID, 2, NULL, 'dashboards/erm');
 
-- Get role and permission type IDs
SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';
 
-- Get the ID of the new submodule
DECLARE @ermDashboardModuleID INT;
SELECT @ermDashboardModuleID = id FROM modules WHERE identifier = 'erm_dashboard';
 
-- Assign view, modify, and manage permissions to super admin
IF NOT EXISTS (
    SELECT 1 FROM role_permissions
    WHERE role_id = @superAdminRoleID
    AND module_id = @ermDashboardModuleID
    AND role_permission_type_id = @viewPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @ermDashboardModuleID, @viewPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END
 
IF NOT EXISTS (
    SELECT 1 FROM role_permissions
    WHERE role_id = @superAdminRoleID
    AND module_id = @ermDashboardModuleID
    AND role_permission_type_id = @modifyPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @ermDashboardModuleID, @modifyPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END
 
IF NOT EXISTS (
    SELECT 1 FROM role_permissions
    WHERE role_id = @superAdminRoleID
    AND module_id = @ermDashboardModuleID
    AND role_permission_type_id = @managePermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @ermDashboardModuleID, @managePermissionTypeID, 1, @superAdminRoleID, GETDATE());
END
-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
-- Remove permissions for 'erm_dashboard' submodule
DELETE FROM role_permissions
WHERE module_id IN (
    SELECT id FROM modules WHERE identifier = 'erm_dashboard'
);
 
-- Remove 'erm_dashboard' submodule
DELETE FROM modules
WHERE identifier = 'erm_dashboard';
-- +goose StatementEnd