-- +goose Up
-- +goose StatementBegin

DECLARE @AssetModuleID INT;

SELECT @AssetModuleID = id 
FROM modules 
WHERE identifier = 'asset';

INSERT INTO modules (
    identifier,
    title_en,
    title_ar,
    module_id,
    [order],
    icon,
    route,
    is_passive
)
VALUES (
    'vulnerability_action_plan',
    'Vulnerability Action Plans',
    N'خطط عمل الثغرات الأمنية',
    @AssetModuleID,
    8,
    NULL,
    '/asset/vulnerabilities/action-plans',
    1
);

DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;
DECLARE @newSubModuleID INT;

SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';

SELECT @newSubModuleID = id 
FROM modules 
WHERE identifier = 'vulnerability_action_plan';

IF NOT EXISTS (
    SELECT 1 FROM role_permissions 
    WHERE role_id = @superAdminRoleID 
      AND module_id = @newSubModuleID 
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
        @newSubModuleID,
        @viewPermissionTypeID,
        1,
        @superAdminRoleID,
        GETDATE()
    );
END

IF NOT EXISTS (
    SELECT 1 FROM role_permissions 
    WHERE role_id = @superAdminRoleID 
      AND module_id = @newSubModuleID 
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
        @newSubModuleID,
        @modifyPermissionTypeID,
        1,
        @superAdminRoleID,
        GETDATE()
    );
END

IF NOT EXISTS (
    SELECT 1 FROM role_permissions 
    WHERE role_id = @superAdminRoleID 
      AND module_id = @newSubModuleID 
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
        @newSubModuleID,
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
    WHERE identifier = 'vulnerability_action_plan'
);

DELETE FROM modules
WHERE identifier = 'vulnerability_action_plan';

-- +goose StatementEnd