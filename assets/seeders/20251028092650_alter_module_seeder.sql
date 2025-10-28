-- +goose Up
-- +goose StatementBegin
DECLARE @dashboardModuleID INT;
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;

-- Get parent module ID
SELECT @dashboardModuleID = id FROM modules WHERE identifier = 'dashboard';

-- Update order of existing ERM Dashboard to 3
UPDATE modules
SET [order] = 3
WHERE identifier = 'erm_dashboard';

INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES
    ('compliance_dashboard', 'Compliance Dashboard', 'لوحة القيادة للامتثال', @dashboardModuleID, 2, NULL, 'dashboards/compliance');

-- Get role and permission type IDs
SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';

-- Get the ID of the new submodule
DECLARE @complianceDashboardModuleID INT;
SELECT @complianceDashboardModuleID = id FROM modules WHERE identifier = 'compliance_dashboard';

-- Assign view, modify, and manage permissions to super admin
IF NOT EXISTS (
    SELECT 1 FROM role_permissions
    WHERE role_id = @superAdminRoleID
    AND module_id = @complianceDashboardModuleID
    AND role_permission_type_id = @viewPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @complianceDashboardModuleID, @viewPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END;

IF NOT EXISTS (
    SELECT 1 FROM role_permissions
    WHERE role_id = @superAdminRoleID
    AND module_id = @complianceDashboardModuleID
    AND role_permission_type_id = @modifyPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @complianceDashboardModuleID, @modifyPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END;

IF NOT EXISTS (
    SELECT 1 FROM role_permissions
    WHERE role_id = @superAdminRoleID
    AND module_id = @complianceDashboardModuleID
    AND role_permission_type_id = @managePermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @complianceDashboardModuleID, @managePermissionTypeID, 1, @superAdminRoleID, GETDATE());
END;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
-- Remove permissions for 'compliance_dashboard' submodule
DELETE FROM role_permissions
WHERE module_id IN (
    SELECT id FROM modules WHERE identifier = 'compliance_dashboard'
);

-- Remove 'compliance_dashboard' submodule
DELETE FROM modules
WHERE identifier = 'compliance_dashboard';

-- Revert ERM Dashboard order back to 2
UPDATE modules
SET [order] = 2
WHERE identifier = 'erm_dashboard';
-- +goose StatementEnd
