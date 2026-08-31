-- +goose Up
-- +goose StatementBegin
DECLARE @organizationModuleID INT;
DECLARE @orgOrder INT;
DECLARE @projectModuleID INT;
DECLARE @projectRegisterID INT;
DECLARE @oldProjectOrder INT;

DECLARE @superAdminRoleID INT;

SELECT @organizationModuleID = id, @orgOrder = [order]
FROM modules
WHERE identifier = 'organization';

UPDATE modules
SET [order] = [order] + 1
WHERE module_id IS NULL
  AND [order] > @orgOrder;

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
    'project',
    'Projects',
    N'المشاريع',
    NULL,
    @orgOrder + 1,
    'icon-project.svg',
    '/projects'
);

SELECT @projectModuleID = SCOPE_IDENTITY();

SELECT @projectRegisterID = id, @oldProjectOrder = [order]
FROM modules
WHERE identifier = 'project_management';

UPDATE modules
SET identifier = 'project_register',
    title_en = 'Project Register',
    title_ar = N'سجل المشاريع',
    module_id = @projectModuleID,
    [order] = 1,
    icon = NULL,
    route = '/projects/project-register'
WHERE id = @projectRegisterID;

UPDATE modules
SET [order] = [order] - 1
WHERE module_id = @organizationModuleID
  AND [order] > @oldProjectOrder;

SELECT @superAdminRoleID = id
FROM roles
WHERE is_super_admin = 1;

INSERT INTO role_permissions (
    role_id,
    module_id,
    role_permission_type_id,
    is_granted,
    created_by,
    created_at
)
SELECT 
    @superAdminRoleID,
    m.id,
    rpt.id,
    1,
    @superAdminRoleID,
    GETDATE()
FROM modules m
JOIN role_permission_types rpt 
    ON rpt.type IN ('view', 'modify', 'manage')
WHERE m.identifier = 'project_register'
  AND NOT EXISTS (
      SELECT 1 
      FROM role_permissions rp
      WHERE rp.role_id = @superAdminRoleID
        AND rp.module_id = m.id
        AND rp.role_permission_type_id = rpt.id
  );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DECLARE @organizationModuleID INT;
DECLARE @orgOrder INT;
DECLARE @projectModuleID INT;
DECLARE @projectRegisterID INT;
DECLARE @maxOrgOrder INT;

SELECT @organizationModuleID = id, @orgOrder = [order]
FROM modules
WHERE identifier = 'organization';

SELECT @projectModuleID = id
FROM modules
WHERE identifier = 'project';

SELECT @projectRegisterID = id
FROM modules
WHERE identifier = 'project_register';

-- 1. Restore project_register back under Organization as project_management
SELECT @maxOrgOrder = ISNULL(MAX([order]), 0) + 1
FROM modules
WHERE module_id = @organizationModuleID;

UPDATE modules
SET identifier = 'project_management',
    title_en = 'Projects',
    title_ar = N'المشاريع',
    module_id = @organizationModuleID,
    [order] = @maxOrgOrder,
    icon = NULL,
    route = '/organization/projects'
WHERE id = @projectRegisterID;

-- 2. Delete Project main module
DELETE FROM modules
WHERE id = @projectModuleID;

-- 3. Shift all main modules after Organization back by -1
UPDATE modules
SET [order] = [order] - 1
WHERE module_id IS NULL
  AND [order] > @orgOrder;
-- +goose StatementEnd