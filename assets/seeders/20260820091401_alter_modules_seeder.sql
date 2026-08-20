-- +goose Up
-- +goose StatementBegin
DECLARE @incidentModuleID INT;
DECLARE @incidentDashboardModuleID INT;
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;

SELECT @incidentModuleID = id
FROM modules
WHERE identifier = 'incident';
 
UPDATE modules
SET [order] = [order] + 1
WHERE module_id = @incidentModuleID;
 
INSERT INTO modules (
    identifier,
    title_en,
    title_ar,
    module_id,
    [order],
    icon,
    route
)
VALUES (
    'incident_dashboard',
    'Dashboard',
    N'لوحة المعلومات',
    @incidentModuleID,
    1,
    NULL,
    '/incident/dashboard'
);
 
SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view'; 
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';
SELECT @incidentDashboardModuleID = id FROM modules WHERE identifier = 'incident_dashboard';
 
IF NOT EXISTS (
    SELECT 1
    FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @incidentDashboardModuleID
      AND role_permission_type_id = @viewPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (
        role_id,
        module_id,
        role_permission_type_id,
        is_granted,
        created_by,
        created_at
    )
    VALUES (
        @superAdminRoleID,
        @incidentDashboardModuleID,
        @viewPermissionTypeID,
        1,
        @superAdminRoleID,
        GETDATE()
    );
END
 
IF NOT EXISTS (
    SELECT 1
    FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @incidentDashboardModuleID
      AND role_permission_type_id = @modifyPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (
        role_id,
        module_id,
        role_permission_type_id,
        is_granted,
        created_by,
        created_at
    )
    VALUES (
        @superAdminRoleID,
        @incidentDashboardModuleID,
        @modifyPermissionTypeID,
        1,
        @superAdminRoleID,
        GETDATE()
    );
END
 
IF NOT EXISTS (
    SELECT 1
    FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @incidentDashboardModuleID
      AND role_permission_type_id = @managePermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (
        role_id,
        module_id,
        role_permission_type_id,
        is_granted,
        created_by,
        created_at
    )
    VALUES (
        @superAdminRoleID,
        @incidentDashboardModuleID,
        @managePermissionTypeID,
        1,
        @superAdminRoleID,
        GETDATE()
    );
END
-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
DECLARE @incidentModuleID INT;
 
SELECT @incidentModuleID = id
FROM modules
WHERE identifier = 'incident';
 
DELETE FROM role_permissions
WHERE module_id IN (
    SELECT id
    FROM modules
    WHERE identifier = 'incident_dashboard'
);
 
DELETE FROM modules
WHERE identifier = 'incident_dashboard';
 
UPDATE modules
SET [order] = [order] - 1
WHERE module_id = @incidentModuleID;
-- +goose StatementEnd