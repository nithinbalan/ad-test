-- +goose Up
-- +goose StatementBegin
DECLARE @AwarenessReportingModuleID INT;
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;

SELECT @AwarenessReportingModuleID = id
FROM modules
WHERE identifier = 'awareness_reporting';

UPDATE modules
SET [order] = [order] + 1
WHERE module_id = @AwarenessReportingModuleID;

INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES ('awareness_reporting_dashboard',  'Awareness Reporting Dashboard', N'لوحة معلومات تقارير التوعية', @AwarenessReportingModuleID, 1, NULL,'/awareness-reporting/dashboard');

SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';

DECLARE @AwarenessReportingDashboardModuleID INT;
SELECT @AwarenessReportingDashboardModuleID = id FROM modules WHERE identifier = 'awareness_reporting_dashboard';

IF NOT EXISTS (
    SELECT 1
    FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @AwarenessReportingDashboardModuleID
      AND role_permission_type_id = @viewPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions ( role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES ( @superAdminRoleID, @AwarenessReportingDashboardModuleID, @viewPermissionTypeID, 1, @superAdminRoleID, GETDATE() );
END;

IF NOT EXISTS (
    SELECT 1
    FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @AwarenessReportingDashboardModuleID
      AND role_permission_type_id = @modifyPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id,module_id, role_permission_type_id, is_granted, created_by, created_at
    )
    VALUES ( @superAdminRoleID, @AwarenessReportingDashboardModuleID, @modifyPermissionTypeID, 1, @superAdminRoleID, GETDATE()
 );
END;

IF NOT EXISTS (
    SELECT 1
    FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @AwarenessReportingDashboardModuleID
      AND role_permission_type_id = @managePermissionTypeID
)
BEGIN
    INSERT INTO role_permissions ( role_id, module_id, role_permission_type_id, is_granted,created_by, created_at
    )
    VALUES ( @superAdminRoleID, @AwarenessReportingDashboardModuleID, @managePermissionTypeID, 1, @superAdminRoleID,GETDATE()
    );
END;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM role_permissions
WHERE module_id IN (
    SELECT id
    FROM modules
    WHERE identifier = 'awareness_reporting_dashboard'
);

DELETE FROM modules
WHERE identifier = 'awareness_reporting_dashboard';

UPDATE modules
SET [order] = [order] - 1
WHERE module_id = (
    SELECT id
    FROM modules
    WHERE identifier = 'awareness_reporting'
);
-- +goose StatementEnd