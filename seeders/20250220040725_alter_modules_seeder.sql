-- +goose Up
-- +goose StatementBegin
DECLARE @strategyModuleID INT;

UPDATE modules
SET [order] = 13
WHERE identifier = 'setting';


-- Insert Asset Module
INSERT INTO modules (identifier, title_en, title_ar, [order], icon, route)
VALUES ('strategy', 'Strategy', 'الإستراتيجية', 12, 'icon-strategy.svg', '/strategy');
SELECT @strategyModuleID = SCOPE_IDENTITY();

-- Insert Asset Sub-Modules
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('swot', 'SWOT', 'كدح', @strategyModuleID, 1, NULL, '/strategy/swot'),
    ('pestle', 'PESTLE', 'مدقة', @strategyModuleID, 2, NULL, '/strategy/pestle'),
    ('strategy_profile', 'Strategy Profile', 'ملف الإستراتيجية', @strategyModuleID, 3, NULL, '/strategy/strategy-profile'),
    ('strategy_objective', 'Strategy Objectives', 'أهداف الإستراتيجية', @strategyModuleID, 4, NULL, '/strategy/strategy-objectives'),
    ('kpi', 'KPI', 'مؤشرات الأداء الرئيسية', @strategyModuleID, 5, NULL, '/strategy/kpi');

-- Declare role and permission type IDs
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;

-- Get role and permission type IDs
SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';

-- Assign permissions to all submodules for the Super Admin role
DECLARE @subModuleID INT;
DECLARE submodule_cursor CURSOR FOR 
    SELECT id FROM modules WHERE module_id = @strategyModuleID;

OPEN submodule_cursor;
FETCH NEXT FROM submodule_cursor INTO @subModuleID;

WHILE @@FETCH_STATUS = 0
BEGIN
    -- Ensure View Permission is not duplicated
    IF NOT EXISTS (
        SELECT 1 FROM role_permissions 
        WHERE role_id = @superAdminRoleID 
          AND module_id = @subModuleID 
          AND role_permission_type_id = @viewPermissionTypeID
    )
    BEGIN
        INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
        VALUES (@superAdminRoleID, @subModuleID, @viewPermissionTypeID, 1, @superAdminRoleID, GETDATE());
    END

    -- Ensure Modify Permission is not duplicated
    IF NOT EXISTS (
        SELECT 1 FROM role_permissions 
        WHERE role_id = @superAdminRoleID 
          AND module_id = @subModuleID 
          AND role_permission_type_id = @modifyPermissionTypeID
    )
    BEGIN
        INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
        VALUES (@superAdminRoleID, @subModuleID, @modifyPermissionTypeID, 1, @superAdminRoleID, GETDATE());
    END

    -- Ensure Manage Permission is not duplicated
    IF NOT EXISTS (
        SELECT 1 FROM role_permissions 
        WHERE role_id = @superAdminRoleID 
          AND module_id = @subModuleID 
          AND role_permission_type_id = @managePermissionTypeID
    )
    BEGIN
        INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
        VALUES (@superAdminRoleID, @subModuleID, @managePermissionTypeID, 1, @superAdminRoleID, GETDATE());
    END

    FETCH NEXT FROM submodule_cursor INTO @subModuleID;
END

CLOSE submodule_cursor;
DEALLOCATE submodule_cursor;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DECLARE @strategyModuleID INT;

SELECT @strategyModuleID = id FROM modules WHERE identifier = 'strategy';

UPDATE modules 
SET [order] = 12 
WHERE identifier = 'setting';

DELETE FROM modules
WHERE module_id = @strategyModuleID;

DELETE FROM modules
WHERE identifier = 'strategy';

DELETE FROM role_permissions
WHERE module_id IN (
    SELECT id FROM modules WHERE identifier = 'strategy'
);
-- +goose StatementEnd
