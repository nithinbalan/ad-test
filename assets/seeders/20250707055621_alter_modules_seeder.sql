-- +goose Up
-- +goose StatementBegin

DECLARE @organizationModuleID INT;
DECLARE @projectManagementModuleID INT;
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;

-- Get the organization module ID
SELECT @organizationModuleID = id FROM modules WHERE identifier = 'organization';

-- Insert Project Management module under Organization
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES ('project_management', 'Projects', 'المشاريع', @organizationModuleID, 5, NULL, '/organization/projects');

-- Ensure Organization Report is order 6
UPDATE modules SET [order] = 6 WHERE identifier = 'organization_report' AND module_id = @organizationModuleID;

-- Get the newly inserted Project Management module ID
SET @projectManagementModuleID = (SELECT id FROM modules WHERE identifier = 'project_management');

-- Get IDs for roles and permission types
SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';

-- Assign permissions to Project Management module
IF NOT EXISTS (
    SELECT 1 FROM role_permissions 
    WHERE role_id = @superAdminRoleID 
      AND module_id = @projectManagementModuleID 
      AND role_permission_type_id = @viewPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @projectManagementModuleID, @viewPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

IF NOT EXISTS (
    SELECT 1 FROM role_permissions 
    WHERE role_id = @superAdminRoleID 
      AND module_id = @projectManagementModuleID 
      AND role_permission_type_id = @modifyPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @projectManagementModuleID, @modifyPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

IF NOT EXISTS (
    SELECT 1 FROM role_permissions 
    WHERE role_id = @superAdminRoleID 
      AND module_id = @projectManagementModuleID 
      AND role_permission_type_id = @managePermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @projectManagementModuleID, @managePermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Remove permissions for Project Management module
DELETE FROM role_permissions
WHERE module_id = (SELECT id FROM modules WHERE identifier = 'project_management');

-- Remove Project Management module
DELETE FROM modules WHERE identifier = 'project_management';

-- Restore Organization Report order to 5
UPDATE modules SET [order] = 5 WHERE identifier = 'organization_report' AND module_id = (SELECT id FROM modules WHERE identifier = 'organization');

-- +goose StatementEnd
