-- +goose Up
-- +goose StatementBegin
DECLARE @violationModuleID INT;
DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;
DECLARE @settingModuleID INT;

-- 1. REMOVE OLD violation_report UNDER INCIDENT (IF EXISTS)
DELETE FROM role_permissions
WHERE module_id IN (
    SELECT id FROM modules WHERE identifier = 'violation_report'
);

-- Remove activity logs referencing the module
DELETE FROM activity_logs
WHERE sub_module_id IN (
    SELECT id FROM modules WHERE identifier IN (
        'violation_report'
    )
);

DELETE FROM modules
WHERE identifier = 'violation_report';

-- Restore original order values
UPDATE modules
SET [order] = CASE
    WHEN identifier = 'incident_action_plan' THEN 2
    WHEN identifier = 'incident_report' THEN 3
END
WHERE identifier IN ('incident_action_plan', 'incident_report');

-- 2. INSERT MAIN MODULE: VIOLATION REPORT (NO PERMISSIONS)
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
('violation_report', 'Violation Report', 'تقرير الانتهاك', NULL, 18, 'icon-violation-report.svg', '/violation-report');

SELECT @violationModuleID = id 
FROM modules 
WHERE identifier = 'violation_report';

-- 3. INSERT SUB MODULE: VIOLATION REPORT REGISTER
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
('violation_report_register', 'Violation Report Register', 'سجل تقارير الانتهاك', @violationModuleID, 1, NULL, '/violation-report/register');

---------------------------------------------------
-- 3. GET SETTING MODULE ID
---------------------------------------------------
SELECT @settingModuleID = id
FROM modules
WHERE identifier = 'setting';

-- 4. INSERT SUB MODULE: VIOLATION REPORT SETTING
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
('violation_report_setting', 'Violation Report', 'تقرير الانتهاك', @settingModuleID, 17, NULL, '/setting/violation-report-setting');

-- 5. GET ROLE & PERMISSION TYPES
SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';

-- 6. ASSIGN PERMISSIONS ONLY TO SUB MODULES
INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
SELECT @superAdminRoleID, m.id, rpt.id, 1, @superAdminRoleID, GETDATE()
FROM modules m
JOIN role_permission_types rpt 
    ON rpt.type IN ('view', 'modify', 'manage')
WHERE m.identifier IN (
    'violation_report_register',
    'violation_report_setting'
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

-- REMOVE ROLE PERMISSIONS (SUB MODULES ONLY)
DELETE FROM role_permissions
WHERE module_id IN (
    SELECT id FROM modules
    WHERE identifier IN (
        'violation_report_register',
        'violation_report_setting'
    )
);

-- Remove activity logs referencing the sub-modules
DELETE FROM activity_logs
WHERE sub_module_id IN (
    SELECT id FROM modules WHERE identifier IN (
        'violation_report_register',
        'violation_report_setting'
    )
);

-- REMOVE SUB MODULES
DELETE FROM modules
WHERE identifier IN (
    'violation_report_register',
    'violation_report_setting'
);

-- REMOVE MAIN MODULE
DELETE FROM modules
WHERE identifier = 'violation_report';

-- +goose StatementEnd