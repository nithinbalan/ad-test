-- +goose Up
-- +goose StatementBegin
DECLARE @TPRMModuleID INT;
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;

-- Get TPRM parent module ID
SELECT @TPRMModuleID = id 
FROM modules 
WHERE identifier = 'tprm';

-- Insert new submodule: Report
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
('tprm_report', 'Report', 'تقرير', @TPRMModuleID, 5, NULL, '/tprm/report');

-- Get Super Admin role ID and permission types
SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';

-- Get new submodule ID
DECLARE @reportSubmoduleID INT;
SELECT @reportSubmoduleID = id FROM modules WHERE identifier = 'tprm_report';

-- Assign permissions to Super Admin
INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
VALUES
(@superAdminRoleID, @reportSubmoduleID, @viewPermissionTypeID, 1, @superAdminRoleID, GETDATE()),
(@superAdminRoleID, @reportSubmoduleID, @modifyPermissionTypeID, 1, @superAdminRoleID, GETDATE()),
(@superAdminRoleID, @reportSubmoduleID, @managePermissionTypeID, 1, @superAdminRoleID, GETDATE());
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
-- Remove permissions for 'tprm_report' submodule
DELETE FROM role_permissions
WHERE module_id IN (SELECT id FROM modules WHERE identifier = 'tprm_report');

-- Remove the submodule
DELETE FROM modules
WHERE identifier = 'tprm_report';
-- +goose StatementEnd
