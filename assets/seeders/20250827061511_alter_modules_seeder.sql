-- +goose Up
-- +goose StatementBegin
DECLARE @knowledgeCenterModuleID INT;
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;

-- Get parent module ID
SELECT @knowledgeCenterModuleID = id FROM modules WHERE identifier = 'knowledge_center';

-- Insert new submodule: articles
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('articles', 'Articles', 'مقالات', @knowledgeCenterModuleID, 2, NULL, '/knowledge-center/articles');

-- Get role and permission type IDs
SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';

-- Get the ID of the new submodule
DECLARE @articlesModuleID INT;
SELECT @articlesModuleID = id FROM modules WHERE identifier = 'articles';

-- Assign view, modify, and manage permissions to super admin
IF NOT EXISTS (
    SELECT 1 FROM role_permissions 
    WHERE role_id = @superAdminRoleID 
    AND module_id = @articlesModuleID 
    AND role_permission_type_id = @viewPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @articlesModuleID, @viewPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

IF NOT EXISTS (
    SELECT 1 FROM role_permissions 
    WHERE role_id = @superAdminRoleID 
    AND module_id = @articlesModuleID 
    AND role_permission_type_id = @modifyPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @articlesModuleID, @modifyPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

IF NOT EXISTS (
    SELECT 1 FROM role_permissions 
    WHERE role_id = @superAdminRoleID 
    AND module_id = @articlesModuleID 
    AND role_permission_type_id = @managePermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @articlesModuleID, @managePermissionTypeID, 1, @superAdminRoleID, GETDATE());
END
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove permissions for 'articles' submodule
DELETE FROM role_permissions
WHERE module_id IN (
    SELECT id FROM modules WHERE identifier = 'articles'
);

-- Remove 'articles' submodule
DELETE FROM modules
WHERE identifier = 'articles';
-- +goose StatementEnd
