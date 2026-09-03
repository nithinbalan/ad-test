-- +goose Up
-- +goose StatementBegin
DECLARE @threatIntelligenceModuleID INT;
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;
DECLARE @settingModuleID INT;

SELECT @settingModuleID = id
FROM modules
WHERE identifier = 'setting';

INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
('threat_intelligence', 'Threat Intelligence', 'تتبع التهديدات', NULL, 16, 'icon-threat-intelligence.svg', '/threat-intelligence');

SELECT @threatIntelligenceModuleID = id 
FROM modules 
WHERE identifier = 'threat_intelligence';

INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
('planned_advisory', 'Planned Advisory', 'الإرشادات المخطط لها', @threatIntelligenceModuleID, 1, NULL, '/threat-intelligence/planned-advisory');  

INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
('threat_intelligence_setting', 'Threat Intelligence', 'تتبع التهديدات', @settingModuleID, 15, NULL, '/settings/threat-intelligence');

SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';

INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
SELECT @superAdminRoleID, m.id, rpt.id, 1, @superAdminRoleID, GETDATE()
FROM modules m
JOIN role_permission_types rpt 
    ON rpt.type IN ('view', 'modify', 'manage')
WHERE m.identifier IN (
    'planned_advisory',
    'threat_intelligence_setting'
)
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
    SELECT id FROM modules WHERE identifier IN ('planned_advisory','threat_intelligence_setting'));
DELETE FROM activity_logs WHERE sub_module_id IN (
    SELECT id FROM modules WHERE identifier IN ('planned_advisory','threat_intelligence_setting'));
DELETE FROM modules WHERE identifier IN ('planned_advisory','threat_intelligence_setting');
DELETE FROM modules WHERE identifier = 'threat_intelligence';
-- +goose StatementEnd
