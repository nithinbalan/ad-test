-- +goose Up
-- +goose StatementBegin

DECLARE @auditModuleID INT;
DECLARE @auditDashboardModuleID INT;
DECLARE @superAdminRoleID INT;

DECLARE @viewPermissionTypeID INT;
DECLARE @modifyPermissionTypeID INT;
DECLARE @managePermissionTypeID INT;

SELECT @auditModuleID = id
FROM modules
WHERE identifier = 'audit';

-- Shift existing audit modules
UPDATE modules
SET [order] = [order] + 1
WHERE module_id = @auditModuleID;

-- Add Audit Dashboard
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
    'audit_dashboard',
    'Dashboard',
    N'لوحة المعلومات',
    @auditModuleID,
    1,
    NULL,
    '/audit/dashboard'
);

-- Get required IDs
SELECT @superAdminRoleID = id
FROM roles
WHERE is_super_admin = 1;

SELECT @viewPermissionTypeID = id
FROM role_permission_types
WHERE type = 'view';

SELECT @modifyPermissionTypeID = id
FROM role_permission_types
WHERE type = 'modify';

SELECT @managePermissionTypeID = id
FROM role_permission_types
WHERE type = 'manage';

SELECT @auditDashboardModuleID = id
FROM modules
WHERE identifier = 'audit_dashboard';

-- View permission
IF NOT EXISTS (
    SELECT 1
    FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @auditDashboardModuleID
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
        @auditDashboardModuleID,
        @viewPermissionTypeID,
        1,
        @superAdminRoleID,
        GETDATE()
    );
END

-- Modify permission
IF NOT EXISTS (
    SELECT 1
    FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @auditDashboardModuleID
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
        @auditDashboardModuleID,
        @modifyPermissionTypeID,
        1,
        @superAdminRoleID,
        GETDATE()
    );
END

-- Manage permission
IF NOT EXISTS (
    SELECT 1
    FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @auditDashboardModuleID
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
        @auditDashboardModuleID,
        @managePermissionTypeID,
        1,
        @superAdminRoleID,
        GETDATE()
    );
END

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DECLARE @auditModuleID INT;

SELECT @auditModuleID = id
FROM modules
WHERE identifier = 'audit';

-- Remove Audit Dashboard permissions
DELETE FROM role_permissions
WHERE module_id IN (
    SELECT id
    FROM modules
    WHERE identifier = 'audit_dashboard'
);

-- Remove Audit Dashboard
DELETE FROM modules
WHERE identifier = 'audit_dashboard';

-- Restore existing audit module orders
UPDATE modules
SET [order] = [order] - 1
WHERE module_id = @auditModuleID;

-- +goose StatementEnd