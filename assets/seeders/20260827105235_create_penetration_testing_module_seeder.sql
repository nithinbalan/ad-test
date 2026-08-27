-- +goose Up
-- +goose StatementBegin
DECLARE @moduleID INT;
DECLARE @settingModuleID INT;
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;

SELECT @moduleID = id FROM modules WHERE identifier = 'penetration_testing';

IF @moduleID IS NULL
BEGIN
    INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
    VALUES ('penetration_testing', 'Penetration Testing', N'اختبار الاختراق', NULL, 0, NULL, '/penetration-testing');
    SET @moduleID = SCOPE_IDENTITY();
END

SELECT @settingModuleID = id FROM modules WHERE identifier = 'penetration_testing_setting';

IF @settingModuleID IS NULL
BEGIN
    INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
    VALUES ('penetration_testing_setting', 'Penetration Testing Masters', N'البيانات الرئيسية لاختبار الاختراق', @moduleID, 1, NULL, '/penetration-testing');
    SET @settingModuleID = SCOPE_IDENTITY();
END

SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';

IF NOT EXISTS (
    SELECT 1 FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @settingModuleID
      AND role_permission_type_id = @viewPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @settingModuleID, @viewPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

IF NOT EXISTS (
    SELECT 1 FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @settingModuleID
      AND role_permission_type_id = @modifyPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @settingModuleID, @modifyPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

IF NOT EXISTS (
    SELECT 1 FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @settingModuleID
      AND role_permission_type_id = @managePermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @settingModuleID, @managePermissionTypeID, 1, @superAdminRoleID, GETDATE());
END
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM role_permissions
WHERE module_id IN (SELECT id FROM modules WHERE identifier IN ('penetration_testing_setting', 'penetration_testing'));

DELETE FROM modules WHERE identifier = 'penetration_testing_setting';
DELETE FROM modules WHERE identifier = 'penetration_testing';
-- +goose StatementEnd
