-- +goose Up
-- +goose StatementBegin
DECLARE @organizationModuleId INT;
DECLARE @organizationChartModuleID INT;
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT;
DECLARE @modifyPermissionTypeID INT;
DECLARE @managePermissionTypeID INT;

-- Get parent module
SELECT @organizationModuleId = id FROM modules WHERE identifier = 'organization';

-- Insert Organization Chart module
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES ('organization_chart', 'Organization Chart', 'الهيكل التنظيمي', @organizationModuleId, 3, NULL, '/organization/organization-chart');

-- Get newly inserted module ID
SELECT @organizationChartModuleID = SCOPE_IDENTITY();

-- Super Admin & permission types
SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';

-- Assign permissions
INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
VALUES
(@superAdminRoleID, @organizationChartModuleID, @viewPermissionTypeID, 1, @superAdminRoleID, GETDATE()),
(@superAdminRoleID, @organizationChartModuleID, @modifyPermissionTypeID, 1, @superAdminRoleID, GETDATE()),
(@superAdminRoleID, @organizationChartModuleID, @managePermissionTypeID, 1, @superAdminRoleID, GETDATE());
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM role_permissions 
WHERE module_id IN (SELECT id FROM modules WHERE identifier = 'organization_chart');

DELETE FROM modules WHERE identifier = 'organization_chart';
-- +goose StatementEnd
