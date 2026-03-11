-- +goose Up
-- +goose StatementBegin
DECLARE @organizationModuleId INT;
DECLARE @committeeManagementModuleID INT;
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT;
DECLARE @modifyPermissionTypeID INT;
DECLARE @managePermissionTypeID INT;

-- Get parent module
SELECT @organizationModuleId = id FROM modules WHERE identifier = 'organization';

-- Insert Organization Chart module
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES ('committee_management', 'Committee Management', 'إدارة اللجان', @organizationModuleId, 4, NULL, '/organization/committee');


-- Ensure Organization Report is order 6
UPDATE modules SET [order] = 5 WHERE identifier = 'policies' AND module_id = @organizationModuleId;
UPDATE modules SET [order] = 6 WHERE identifier = 'process' AND module_id = @organizationModuleId;
UPDATE modules SET [order] = 7 WHERE identifier = 'project_management' AND module_id = @organizationModuleId;
UPDATE modules SET [order] = 8 WHERE identifier = 'organization_report' AND module_id = @organizationModuleId;
-- Get newly inserted module ID
SELECT @committeeManagementModuleID = SCOPE_IDENTITY();

-- Super Admin & permission types
SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';

-- Assign permissions
INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
VALUES
(@superAdminRoleID, @committeeManagementModuleID, @viewPermissionTypeID, 1, @superAdminRoleID, GETDATE()),
(@superAdminRoleID, @committeeManagementModuleID, @modifyPermissionTypeID, 1, @superAdminRoleID, GETDATE()),
(@superAdminRoleID, @committeeManagementModuleID, @managePermissionTypeID, 1, @superAdminRoleID, GETDATE());
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM role_permissions 
WHERE module_id IN (SELECT id FROM modules WHERE identifier = 'committee_management');

DELETE FROM modules WHERE identifier = 'committee_management';

-- Restore Organization Report order to 5
UPDATE modules SET [order] = 4 WHERE identifier = 'policies' AND module_id = (SELECT id FROM modules WHERE identifier = 'organization');
UPDATE modules SET [order] = 5 WHERE identifier = 'process' AND module_id = (SELECT id FROM modules WHERE identifier = 'organization');
UPDATE modules SET [order] = 6 WHERE identifier = 'project_management' AND module_id = (SELECT id FROM modules WHERE identifier = 'organization');
UPDATE modules SET [order] = 7 WHERE identifier = 'organization_report' AND module_id = (SELECT id FROM modules WHERE identifier = 'organization');
-- +goose StatementEnd
