-- +goose Up
-- +goose StatementBegin
DECLARE @projectModuleID INT;
DECLARE @projectDashboardModuleID INT;
DECLARE @projectReportModuleID INT;
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;

SELECT @projectModuleID = id
FROM modules
WHERE identifier = 'project';

UPDATE modules
SET [order] = [order] + 1
WHERE module_id = @projectModuleID;

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
    'project_dashboard',
    'Dashboard',
    N'لوحة المعلومات',
    @projectModuleID,
    1,
    NULL,
    '/projects/dashboard'
);

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
    'project_report',
    'Report',
    N'تقرير',
    @projectModuleID,
    4,
    NULL,
    '/projects/report'
);

SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';

SELECT @projectDashboardModuleID = id FROM modules WHERE identifier = 'project_dashboard';
SELECT @projectReportModuleID = id FROM modules WHERE identifier = 'project_report';

IF NOT EXISTS (
    SELECT 1
    FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @projectDashboardModuleID
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
        @projectDashboardModuleID,
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
      AND module_id = @projectDashboardModuleID
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
        @projectDashboardModuleID,
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
      AND module_id = @projectDashboardModuleID
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
        @projectDashboardModuleID,
        @managePermissionTypeID,
        1,
        @superAdminRoleID,
        GETDATE()
    );
END

IF NOT EXISTS (
    SELECT 1
    FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @projectReportModuleID
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
        @projectReportModuleID,
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
      AND module_id = @projectReportModuleID
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
        @projectReportModuleID,
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
      AND module_id = @projectReportModuleID
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
        @projectReportModuleID,
        @managePermissionTypeID,
        1,
        @superAdminRoleID,
        GETDATE()
    );
END
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DECLARE @projectModuleID INT;

SELECT @projectModuleID = id
FROM modules
WHERE identifier = 'project';

DELETE FROM role_permissions
WHERE module_id IN (
    SELECT id
    FROM modules
    WHERE identifier IN ('project_dashboard', 'project_report')
);

DELETE FROM modules
WHERE identifier IN ('project_dashboard', 'project_report');

UPDATE modules
SET [order] = [order] - 1
WHERE module_id = @projectModuleID;
-- +goose StatementEnd
