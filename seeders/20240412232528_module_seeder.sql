-- +goose Up
-- +goose StatementBegin
DECLARE @dashboardsModuleID INT;
DECLARE @orgModuleID INT;
DECLARE @compModuleID INT;
DECLARE @docsModuleID INT;
DECLARE @incidentsModuleID INT;
DECLARE @mgmtReviewModuleID INT;
DECLARE @riskModuleID INT;
DECLARE @taskModuleID INT;
DECLARE @settingsModuleID INT;

-- Insert Dashboards Module
INSERT INTO modules
    (identifier, title_en, title_ar, [order])
VALUES
    ('dashboards', 'Dashboards', 'لوحات التحكم', 1);
SELECT @dashboardsModuleID = SCOPE_IDENTITY();

-- Insert Organization Module
INSERT INTO modules
    (identifier, title_en, title_ar, [order])
VALUES
    ('organization', 'Organization', 'المنظمة', 2);
SELECT @orgModuleID = SCOPE_IDENTITY();

-- Insert Organization Module
INSERT INTO modules
    (identifier, title_en, title_ar, module_id, [order])
VALUES
    ('strategic_profile', 'Strategic Profile', 'الملف الاستراتيجي', @orgModuleID, 2),
    ('users', 'Users', 'المستخدمين', @orgModuleID, 3),
    ('process', 'Process', 'العمليات', @orgModuleID, 4),
    ('controls', 'Controls', 'ضوابط', @orgModuleID, 5);

-- Insert Compliance Module
INSERT INTO modules
    (identifier, title_en, title_ar, [order])
VALUES
    ('compliance', 'Compliance', 'الامتثال', 3);
SELECT @compModuleID = SCOPE_IDENTITY();

-- Insert Compliance Sub-Modules
INSERT INTO modules
    (identifier, title_en, title_ar, module_id, [order])
VALUES
    ('compliance_register', 'Compliance Register', 'سجل الامتثال', @compModuleID, 2),
    ('sla_contracts', 'SLA & Contracts', 'SLA والعقود', @compModuleID, 3);

-- Insert Documents Module
INSERT INTO modules
    (identifier, title_en, title_ar, [order])
VALUES
    ('documents', 'Documents', 'المستندات', 4);
SELECT @docsModuleID = SCOPE_IDENTITY();

-- Insert Documents Sub-Modules
INSERT INTO modules
    (identifier, title_en, title_ar, module_id, [order])
VALUES
    ('documents', 'Documents', 'المستندات', @docsModuleID, 2),
    ('change_requests', 'Change Requests', 'طلبات التغيير', @docsModuleID, 3);

-- Insert Incidents Module
INSERT INTO modules
    (identifier, title_en, title_ar, [order])
VALUES
    ('incidents', 'Incidents', 'الحوادث', 5);
SELECT @incidentsModuleID = SCOPE_IDENTITY();

-- Insert Incidents Sub-Modules
INSERT INTO modules
    (identifier, title_en, title_ar, module_id, [order])
VALUES
    ('incident_register', 'Incident Register', 'سجل الحوادث', @incidentsModuleID, 2),
    ('action_plans', 'Action Plans', 'خطط العمل', @incidentsModuleID, 3);

-- Insert Management Review Module
INSERT INTO modules
    (identifier, title_en, title_ar, [order])
VALUES
    ('management_review', 'Management Review', 'استعراض الإدارة', 6);
SELECT @mgmtReviewModuleID = SCOPE_IDENTITY();

-- Insert Management Review Sub-Modules
INSERT INTO modules
    (identifier, title_en, title_ar, module_id, [order])
VALUES
    ('meetings', 'Meetings', 'الاجتماعات', @mgmtReviewModuleID, 2),
    ('action_plans', 'Action Plans', 'خطط العمل', @mgmtReviewModuleID, 3);

-- Insert Risk Module
INSERT INTO modules
    (identifier, title_en, title_ar, [order])
VALUES
    ('risk', 'Risk', 'المخاطر', 7);
SELECT @riskModuleID = SCOPE_IDENTITY();

-- Insert Risk Sub-Modules
INSERT INTO modules
    (identifier, title_en, title_ar, module_id, [order])
VALUES
    ('risk_matrix', 'Risk Matrix', 'مصفوفة المخاطر', @riskModuleID, 2),
    ('risk_appetite', 'Risk Appetite', 'ميول المخاطر', @riskModuleID, 3),
    ('potential_risk', 'Potential Risk', 'المخاطر المحتملة', @riskModuleID, 4),
    ('risk_register', 'Risk Register', 'سجل المخاطر', @riskModuleID, 5),
    ('risk_mitigations', 'Risk Mitigations', 'تخفيف المخاطر', @riskModuleID, 6),
    ('kri', 'KRI', 'مؤشرات الأداء الرئيسية', @riskModuleID, 7),
    ('risk_heatmap', 'Risk Heatmap', 'خريطة الخطر', @riskModuleID, 8);

-- Insert Task Module
INSERT INTO modules
    (identifier, title_en, title_ar, [order])
VALUES
    ('task', 'Task', 'المهام', 8);
SELECT @taskModuleID = SCOPE_IDENTITY();

-- Insert Settings Module
INSERT INTO modules
    (identifier, title_en, title_ar, [order])
VALUES
    ('settings', 'Settings', 'الإعدادات', 9);
SELECT @settingsModuleID = SCOPE_IDENTITY();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE modules;
-- +goose StatementEnd
