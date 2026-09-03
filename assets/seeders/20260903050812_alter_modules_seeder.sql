-- +goose Up
-- +goose StatementBegin

DECLARE @projectModuleID INT;
DECLARE @projectSecurityGateModuleID INT;

DECLARE @superAdminRoleID INT;

DECLARE @viewPermissionTypeID INT;
DECLARE @modifyPermissionTypeID INT;
DECLARE @managePermissionTypeID INT;

SELECT @projectModuleID = id
FROM modules
WHERE identifier = 'project';

-- Add Project Security Gate

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
    'project_security_gate',
    'Project Security Gate',
    N'بوابة أمان المشروع',
    @projectModuleID,
    2,
    NULL,
    '/projects/security-gate'
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

SELECT @projectSecurityGateModuleID = id
FROM modules
WHERE identifier = 'project_security_gate';

-- View permission

IF NOT EXISTS (
    SELECT 1
    FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @projectSecurityGateModuleID
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
        @projectSecurityGateModuleID,
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
      AND module_id = @projectSecurityGateModuleID
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
        @projectSecurityGateModuleID,
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
      AND module_id = @projectSecurityGateModuleID
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
        @projectSecurityGateModuleID,
        @managePermissionTypeID,
        1,
        @superAdminRoleID,
        GETDATE()
    );
END

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM role_permissions
WHERE module_id IN (
    SELECT id
    FROM modules
    WHERE identifier = 'project_security_gate'
);

DELETE FROM modules
WHERE identifier = 'project_security_gate'
    AND module_id = (
      SELECT id
      FROM modules
      WHERE identifier = 'project'
  );

-- +goose StatementEnd