-- +goose Up
-- +goose StatementBegin
DECLARE @organizationalContextModuleID INT;

DECLARE @settingModuleID INT;
SET @settingModuleID = (SELECT id FROM modules WHERE identifier = 'setting');

UPDATE modules
SET [order] = 14
WHERE identifier = 'setting';

UPDATE modules
SET [order] = 14
WHERE identifier = 'log';

UPDATE modules
SET [order] = 15
WHERE identifier = 'common_setting';


-- Insert Asset Module
INSERT INTO modules (identifier, title_en, title_ar, [order], icon, route)
VALUES ('organizational_context', 'Organizational Context', 'السياق التنظيمي', 13, 'icon-organization-context.svg', '/organizational-context');
SELECT @organizationalContextModuleID = SCOPE_IDENTITY();

-- Insert Asset Sub-Modules
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('swot_analysis', 'SWOT Analysis', 'تحليل نقاط القوة والضعف والفرص والتهديدات (SWOT).', @organizationalContextModuleID, 1, NULL, '/organizational-context/swot-analysis'),
    ('pestel_analysis', 'PESTEL Analysis', 'تحليل بيستل', @organizationalContextModuleID, 2, NULL, '/organizational-context/pestle-analysis'),
    ('stakeholder_analysis', 'Stakeholder Analysis', 'تحليل أصحاب المصلحة', @organizationalContextModuleID, 3, NULL, '/organizational-context/stakeholder-analysis'),
    ('organizational_context_setting', 'Organizational Context','السياق التنظيمي', @settingModuleID, 13, NULL, '/setting/organizational-context-setting');

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
    SELECT id FROM modules WHERE module_id = @organizationalContextModuleID;

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
DECLARE @organizationalContextModuleID INT;

SELECT @organizationalContextModuleID = id FROM modules WHERE identifier = 'organizational_context';

UPDATE modules 
SET [order] = 12 
WHERE identifier = 'setting';

DELETE FROM modules
WHERE module_id = @organizationalContextModuleID;

DELETE FROM modules
WHERE identifier = 'organizational_context';

DELETE FROM role_permissions
WHERE module_id IN (
    SELECT id FROM modules WHERE identifier = 'organizational_context'
);
-- +goose StatementEnd
