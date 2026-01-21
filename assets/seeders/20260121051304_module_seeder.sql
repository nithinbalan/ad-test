
-- +goose Up
-- +goose StatementBegin
DECLARE @complianceModuleID INT;
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;
DECLARE @consultationRequestModuleID INT;

-- Get parent module ID (Compliance)
SELECT @complianceModuleID = id FROM modules WHERE identifier = 'compliance';
-- Insert new submodule: consultation_request

    INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
    VALUES  ( 'consultation_request', 'Consultation Request','طلب استشارة',@complianceModuleID,6, NULL,'compliance/consultation-request');

-- Get role and permission type IDs
SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';
-- Get the ID of the new submodule
SELECT @consultationRequestModuleID = id FROM modules WHERE identifier = 'consultation_request';
-- Assign view permission
IF NOT EXISTS (
    SELECT 1 FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @consultationRequestModuleID
      AND role_permission_type_id = @viewPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions
        (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES
        (@superAdminRoleID, @consultationRequestModuleID, @viewPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END
-- Assign modify permission
IF NOT EXISTS (
    SELECT 1 FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @consultationRequestModuleID
      AND role_permission_type_id = @modifyPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions
        (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES
        (@superAdminRoleID, @consultationRequestModuleID, @modifyPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END
-- Assign manage permission
IF NOT EXISTS (
    SELECT 1 FROM role_permissions
    WHERE role_id = @superAdminRoleID
      AND module_id = @consultationRequestModuleID
      AND role_permission_type_id = @managePermissionTypeID
)
BEGIN
    INSERT INTO role_permissions
        (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES
        (@superAdminRoleID, @consultationRequestModuleID, @managePermissionTypeID, 1, @superAdminRoleID, GETDATE());
END
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
-- Remove permissions for 'consultation_request' submodule
DELETE FROM role_permissions WHERE module_id IN (SELECT id FROM modules WHERE identifier = 'consultation_request');

-- Remove 'consultation_request' module
DELETE FROM modules WHERE identifier = 'consultation_request';
-- +goose StatementEnd
