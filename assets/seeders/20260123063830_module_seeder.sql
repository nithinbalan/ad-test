-- +goose Up
-- +goose StatementBegin

DECLARE @settingModuleID INT;
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;
DECLARE @consultationRequestSettingModuleID INT;

SELECT @settingModuleID = id 
FROM modules 
WHERE identifier = 'setting';

INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES ('consultation_request_setting', 'Consultation Request', 'طلب التشاور', @settingModuleID, 15, NULL, '/settings/consultation-request');

SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';

SELECT @consultationRequestSettingModuleID = id 
FROM modules 
WHERE identifier = 'consultation_request_setting';

IF NOT EXISTS (
    SELECT 1 FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @consultationRequestSettingModuleID
      AND role_permission_type_id = @viewPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions
        (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES
        (@superAdminRoleID, @consultationRequestSettingModuleID, @viewPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

IF NOT EXISTS (
    SELECT 1 FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @consultationRequestSettingModuleID
      AND role_permission_type_id = @modifyPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions
        (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES
        (@superAdminRoleID, @consultationRequestSettingModuleID, @modifyPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

IF NOT EXISTS (
    SELECT 1 FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @consultationRequestSettingModuleID
      AND role_permission_type_id = @managePermissionTypeID
)
BEGIN
    INSERT INTO role_permissions
        (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES
        (@superAdminRoleID, @consultationRequestSettingModuleID, @managePermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM role_permissions 
WHERE module_id IN (
    SELECT id FROM modules WHERE identifier = 'consultation_request_setting'
);

DELETE FROM modules 
WHERE identifier = 'consultation_request_setting';

-- +goose StatementEnd
