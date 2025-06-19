-- +goose Up
-- +goose StatementBegin
DECLARE @tprmModuleID INT;
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;

UPDATE modules
SET [order] = 14
WHERE identifier = 'setting';

-- Create the main TPRM module
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('tprm', 'TPRM', 'إدارة مخاطر الطرف الثالث', NULL, 13,  'icon-tprm.svg', '/tprm');

-- Fetch the newly inserted TPRM module ID
SET @tprmModuleID = (SELECT id FROM modules WHERE identifier = 'tprm');

-- Create Submodules
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('third_party_register', 'Third Party Register', 'سجل الطرف الثالث', @tprmModuleID, 1, NULL, '/tprm/third-party-register'),
    ('question_library', 'Question Library', 'مكتبة الأسئلة', @tprmModuleID, 2, NULL, '/tprm/question-library'),
    ('question_set', 'Question Set', 'مجموعة الأسئلة', @tprmModuleID, 3, NULL, '/tprm/question-set'),
    ('assessment', 'Assessments', 'التقييم', @tprmModuleID, 4, NULL, '/tprm/assessment');

-- Get IDs for roles and permission types
SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';

-- Assign permissions to all TPRM submodules
DECLARE @subModuleID INT;
DECLARE tprm_cursor CURSOR FOR 
    SELECT id FROM modules WHERE module_id = @tprmModuleID;

OPEN tprm_cursor;
FETCH NEXT FROM tprm_cursor INTO @subModuleID;

WHILE @@FETCH_STATUS = 0
BEGIN
    -- View Permission
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

    -- Modify Permission
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

    -- Manage Permission
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

    FETCH NEXT FROM tprm_cursor INTO @subModuleID;
END

CLOSE tprm_cursor;
DEALLOCATE tprm_cursor;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Remove permissions assigned to TPRM submodules
DELETE FROM role_permissions
WHERE module_id IN (
    SELECT id FROM modules WHERE identifier IN ('third_party_register', 'question_library', 'question_set', 'assessment')
);

-- Remove TPRM submodules
DELETE FROM modules
WHERE identifier IN ('third_party_register', 'question_library', 'question_set', 'assessment');

-- Remove TPRM main module
DELETE FROM modules
WHERE identifier = 'tprm';

-- +goose StatementEnd
