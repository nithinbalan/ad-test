-- +goose Up
-- +goose StatementBegin
DECLARE @riskModuleID INT;
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;
 
-- Get parent module ID
SELECT @riskModuleID = id FROM modules WHERE identifier = 'risk';
 
-- Insert new submodule: identified_risk
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES
    ('identified_risk', 'Identified Risk', 'المخاطر المحددة', @riskModuleID, 3, NULL, '/risk/identified-risk');
 
-- Get role and permission type IDs
SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';
 
-- Get the ID of the new submodule
DECLARE @identifiedModuleID INT;
SELECT @identifiedModuleID = id FROM modules WHERE identifier = 'identified_risk';
 
-- Assign view, modify, and manage permissions to super admin
IF NOT EXISTS (
    SELECT 1 FROM role_permissions
    WHERE role_id = @superAdminRoleID
    AND module_id = @identifiedModuleID
    AND role_permission_type_id = @viewPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @identifiedModuleID, @viewPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END
 
IF NOT EXISTS (
    SELECT 1 FROM role_permissions
    WHERE role_id = @superAdminRoleID
    AND module_id = @identifiedModuleID
    AND role_permission_type_id = @modifyPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @identifiedModuleID, @modifyPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END
 
IF NOT EXISTS (
    SELECT 1 FROM role_permissions
    WHERE role_id = @superAdminRoleID
    AND module_id = @identifiedModuleID
    AND role_permission_type_id = @managePermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @identifiedModuleID, @managePermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

-- Update the order for existing submodules
UPDATE modules SET [order] = 1 WHERE identifier = 'risk_matrix';
UPDATE modules SET [order] = 2 WHERE identifier = 'potential_risk';
UPDATE modules SET [order] = 4 WHERE identifier = 'risk_register';
UPDATE modules SET [order] = 5 WHERE identifier = 'risk_mitigation';
UPDATE modules SET [order] = 6 WHERE identifier = 'kri';
UPDATE modules SET [order] = 7 WHERE identifier = 'risk_appetite_statement';
UPDATE modules SET [order] = 8 WHERE identifier = 'risk_heatmap';
UPDATE modules SET [order] = 9 WHERE identifier = 'risk_report';


-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
-- Remove permissions for 'identified_risk' submodule
DELETE FROM role_permissions
WHERE module_id IN (
    SELECT id FROM modules WHERE identifier = 'identified_risk'
);
 
-- Remove 'identified_risk' submodule
DELETE FROM modules
WHERE identifier = 'identified_risk';
-- revert module orders back
UPDATE modules SET [order] = 1 WHERE identifier = 'risk_matrix';
UPDATE modules SET [order] = 2 WHERE identifier = 'risk_appetite_statement';
UPDATE modules SET [order] = 3 WHERE identifier = 'potential_risk';
UPDATE modules SET [order] = 4 WHERE identifier = 'kri';
UPDATE modules SET [order] = 5 WHERE identifier = 'risk_register';
UPDATE modules SET [order] = 6 WHERE identifier = 'risk_mitigation';
UPDATE modules SET [order] = 7 WHERE identifier = 'risk_heatmap';
UPDATE modules SET [order] = 8 WHERE identifier = 'risk_report';
-- +goose StatementEnd