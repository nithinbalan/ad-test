-- +goose Up
-- +goose StatementBegin
DECLARE @incidentModuleID INT;
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;

UPDATE modules
SET [order] = CASE
    WHEN identifier = 'incident_action_plan' THEN 3
    WHEN identifier = 'incident_report' THEN 4
END
WHERE identifier IN ('incident_action_plan', 'incident_report')

-- Get parent module ID
SELECT @incidentModuleID = id FROM modules WHERE identifier = 'incident';

-- Insert new submodule: articles
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('violation_report', 'Violation Report', 'تقرير الانتهاك', @incidentModuleID, 2, NULL, '/incident/violation-report');

-- Get role and permission type IDs
SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';

-- Get the ID of the new submodule
DECLARE @violationReportModuleId INT;
SELECT @violationReportModuleId = id FROM modules WHERE identifier = 'violation_report';

-- Assign view, modify, and manage permissions to super admin
IF NOT EXISTS (
    SELECT 1 FROM role_permissions 
    WHERE role_id = @superAdminRoleID 
    AND module_id = @violationReportModuleId 
    AND role_permission_type_id = @viewPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @violationReportModuleId, @viewPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

IF NOT EXISTS (
    SELECT 1 FROM role_permissions 
    WHERE role_id = @superAdminRoleID 
    AND module_id = @violationReportModuleId 
    AND role_permission_type_id = @modifyPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @violationReportModuleId, @modifyPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

IF NOT EXISTS (
    SELECT 1 FROM role_permissions 
    WHERE role_id = @superAdminRoleID 
    AND module_id = @violationReportModuleId 
    AND role_permission_type_id = @managePermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @violationReportModuleId, @managePermissionTypeID, 1, @superAdminRoleID, GETDATE());
END
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove permissions for 'annual audit plan' submodule
DELETE FROM role_permissions    
WHERE module_id IN (
    SELECT id FROM modules WHERE identifier = 'violation_report'
);

-- Remove 'annual audit plan' submodule
DELETE FROM modules
WHERE identifier = 'violation_report';

-- Restore original order values
UPDATE modules
SET [order] = CASE
    WHEN identifier = 'incident_action_plan' THEN 2
    WHEN identifier = 'incident_report' THEN 3
END
WHERE identifier IN ('incident_action_plan', 'incident_report');
-- +goose StatementEnd
