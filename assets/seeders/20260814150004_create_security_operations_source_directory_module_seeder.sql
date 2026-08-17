-- +goose Up
-- +goose StatementBegin
DECLARE @securityOperationsModuleID INT;
DECLARE @superAdminRoleID INT;

SELECT @securityOperationsModuleID = id
FROM modules
WHERE identifier = 'security_operations';

IF @securityOperationsModuleID IS NOT NULL
   AND NOT EXISTS (
       SELECT 1
       FROM modules
       WHERE identifier = 'security_operations_source_directory'
   )
BEGIN
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES
('security_operations_source_directory', 'Source Directory', N'دليل المصادر', @securityOperationsModuleID, 1, NULL, '/security-operations/source-directory');
END;

SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;

INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
SELECT @superAdminRoleID, m.id, rpt.id, 1, @superAdminRoleID, GETDATE()
FROM modules m
JOIN role_permission_types rpt
    ON rpt.type IN ('view', 'modify', 'manage')
WHERE m.identifier IN (
    'security_operations_source_directory'
)
AND @superAdminRoleID IS NOT NULL
AND NOT EXISTS (
    SELECT 1 FROM role_permissions rp
    WHERE rp.role_id = @superAdminRoleID
      AND rp.module_id = m.id
      AND rp.role_permission_type_id = rpt.id
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM role_permissions WHERE module_id IN (
    SELECT id FROM modules WHERE identifier IN ('security_operations_source_directory'));
DELETE FROM activity_logs WHERE sub_module_id IN (
    SELECT id FROM modules WHERE identifier IN ('security_operations_source_directory'));
DELETE FROM modules WHERE identifier IN ('security_operations_source_directory');
-- +goose StatementEnd
