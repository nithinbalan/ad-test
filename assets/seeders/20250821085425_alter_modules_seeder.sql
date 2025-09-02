-- +goose Up
-- +goose StatementBegin

DECLARE @organizationModuleID INT;
DECLARE @organizationStructureModuleID INT;
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;

-- Get the organization module ID
SELECT @organizationModuleID = id FROM modules WHERE identifier = 'organization';

-- Insert Organization Structure module under Organization
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES ('organization_structure', 'Organization Structure', 'الهيكل التنظيمي', @organizationModuleID, 3, NULL, '/organization/organization-structure');

-- Ensure Organization  order
UPDATE modules SET [order] = 4 WHERE identifier = 'policies' AND module_id = @organizationModuleID;
UPDATE modules SET [order] = 5 WHERE identifier = 'process' AND module_id = @organizationModuleID;
UPDATE modules SET [order] = 6 WHERE identifier = 'project_management' AND module_id = @organizationModuleID;
UPDATE modules SET [order] = 7 WHERE identifier = 'organization_report' AND module_id = @organizationModuleID;

-- Get the newly inserted Organization Structure module ID
SET @organizationStructureModuleID = (SELECT id FROM modules WHERE identifier = 'organization_structure');

-- Get IDs for roles and permission types
SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';

-- Assign permissions to Organization Structure module
IF NOT EXISTS (
    SELECT 1 FROM role_permissions 
    WHERE role_id = @superAdminRoleID 
      AND module_id = @organizationStructureModuleID 
      AND role_permission_type_id = @viewPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @organizationStructureModuleID, @viewPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

IF NOT EXISTS (
    SELECT 1 FROM role_permissions 
    WHERE role_id = @superAdminRoleID 
      AND module_id = @organizationStructureModuleID 
      AND role_permission_type_id = @modifyPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @organizationStructureModuleID, @modifyPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

IF NOT EXISTS (
    SELECT 1 FROM role_permissions 
    WHERE role_id = @superAdminRoleID 
      AND module_id = @organizationStructureModuleID 
      AND role_permission_type_id = @managePermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @organizationStructureModuleID, @managePermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Remove permissions for Organization Structure module
DELETE FROM role_permissions
WHERE module_id = (SELECT id FROM modules WHERE identifier = 'organization_structure');

-- Remove Organization Structure module
DELETE FROM modules WHERE identifier = 'organization_structure';

-- Restore Organization  order 
UPDATE modules SET [order] = 3 WHERE identifier = 'policies' AND module_id = (SELECT id FROM modules WHERE identifier = 'organization');
UPDATE modules SET [order] = 4 WHERE identifier = 'process' AND module_id = (SELECT id FROM modules WHERE identifier = 'organization');
UPDATE modules SET [order] = 5 WHERE identifier = 'project_management' AND module_id = (SELECT id FROM modules WHERE identifier = 'organization');
UPDATE modules SET [order] = 6 WHERE identifier = 'organization_report' AND module_id = (SELECT id FROM modules WHERE identifier = 'organization');

-- +goose StatementEnd
