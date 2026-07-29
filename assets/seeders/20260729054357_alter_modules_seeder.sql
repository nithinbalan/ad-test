-- +goose Up
-- +goose StatementBegin

DECLARE @assetModuleID INT;
DECLARE @assetRiskMitigationModuleID INT;
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;

-- Get parent module ID
SELECT @assetModuleID = id
FROM modules
WHERE identifier = 'asset';

-- Shift existing modules from order 7 onwards
UPDATE modules
SET [order] = [order] + 1
WHERE module_id = @assetModuleID
  AND [order] >= 7;

-- Insert new submodule
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
    'asset_risk_mitigation',
    'Risk Mitigations',
    N'معالجات المخاطر',
    @assetModuleID,
    7,
    NULL,
    '/asset/risk-mitigations',
    1
);

-- Get role and permission type IDs
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

-- Get new module ID
SELECT @assetRiskMitigationModuleID = id
FROM modules
WHERE identifier = 'asset_risk_mitigation';

-- Assign View permission
IF NOT EXISTS (
    SELECT 1
    FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @assetRiskMitigationModuleID
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
        @assetRiskMitigationModuleID,
        @viewPermissionTypeID,
        1,
        @superAdminRoleID,
        GETDATE()
    );
END

-- Assign Modify permission
IF NOT EXISTS (
    SELECT 1
    FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @assetRiskMitigationModuleID
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
        @assetRiskMitigationModuleID,
        @modifyPermissionTypeID,
        1,
        @superAdminRoleID,
        GETDATE()
    );
END

-- Assign Manage permission
IF NOT EXISTS (
    SELECT 1
    FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @assetRiskMitigationModuleID
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
        @assetRiskMitigationModuleID,
        @managePermissionTypeID,
        1,
        @superAdminRoleID,
        GETDATE()
    );
END

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DECLARE @assetModuleID INT;

SELECT @assetModuleID = id
FROM modules
WHERE identifier = 'asset';

-- Remove permissions
DELETE FROM role_permissions
WHERE module_id IN (
    SELECT id
    FROM modules
    WHERE identifier = 'asset_risk_mitigation'
);

-- Remove module
DELETE FROM modules
WHERE identifier = 'asset_risk_mitigation';

-- Restore original ordering
UPDATE modules
SET [order] = [order] - 1
WHERE module_id = @assetModuleID
  AND [order] >= 8;

-- +goose StatementEnd