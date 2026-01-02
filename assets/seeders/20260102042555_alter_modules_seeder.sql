-- +goose Up
-- +goose StatementBegin
DECLARE @tprmModuleID INT;
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;

-- Get parent module ID
SELECT @tprmModuleID = id FROM modules WHERE identifier = 'tprm';

-- Insert new submodule: Dashboard
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES
    ('tprm_dashboard', 'Dashboard', 'لوحة القيادة', @tprmModuleID, 1, NULL, '/tprm/dashboard');


UPDATE modules SET [order] = 2 WHERE identifier = 'third_party_register' AND module_id = @tprmModuleID;
UPDATE modules SET [order] = 3 WHERE identifier = 'question_library' AND module_id = @tprmModuleID;
UPDATE modules SET [order] = 4 WHERE identifier = 'question_set' AND module_id = @tprmModuleID;
UPDATE modules SET [order] = 5 WHERE identifier = 'assessment' AND module_id = @tprmModuleID;
UPDATE modules SET [order] = 6 WHERE identifier = 'tprm_report' AND module_id = @tprmModuleID;

-- Get role and permission type IDs
SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';

-- Get the ID of the new submodule
DECLARE @tprmDashboardModuleID INT;
SELECT @tprmDashboardModuleID = id FROM modules WHERE identifier = 'tprm_dashboard';

-- Assign view permission
IF NOT EXISTS (
    SELECT 1 FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @tprmDashboardModuleID
      AND role_permission_type_id = @viewPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @tprmDashboardModuleID, @viewPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END;

-- Assign modify permission
IF NOT EXISTS (
    SELECT 1 FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @tprmDashboardModuleID
      AND role_permission_type_id = @modifyPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @tprmDashboardModuleID, @modifyPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END;

-- Assign manage permission
IF NOT EXISTS (
    SELECT 1 FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @tprmDashboardModuleID
      AND role_permission_type_id = @managePermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @tprmDashboardModuleID, @managePermissionTypeID, 1, @superAdminRoleID, GETDATE());
END;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
-- Remove permissions for 'tprm_dashboard' submodule
DELETE FROM role_permissions
WHERE module_id IN (
    SELECT id FROM modules WHERE identifier = 'tprm_dashboard'
);

-- Remove 'tprm_dashboard' submodule
DELETE FROM modules
WHERE identifier = 'tprm_dashboard';

UPDATE modules SET [order] = 1 WHERE identifier = 'third_party_register' AND module_id = @tprmModuleID;
UPDATE modules SET [order] = 2 WHERE identifier = 'question_library' AND module_id = @tprmModuleID;
UPDATE modules SET [order] = 3 WHERE identifier = 'question_set' AND module_id = @tprmModuleID;
UPDATE modules SET [order] = 4 WHERE identifier = 'assessment' AND module_id = @tprmModuleID;
UPDATE modules SET [order] = 5 WHERE identifier = 'tprm_report' AND module_id = @tprmModuleID;
-- +goose StatementEnd
