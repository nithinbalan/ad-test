-- +goose Up
-- +goose StatementBegin
DECLARE @dashboardModuleID INT;
DECLARE @organizationModuleID INT;
DECLARE @complianceModuleID INT;
DECLARE @documentModuleID INT;
DECLARE @incidentModuleID INT;
DECLARE @managementReviewModuleID INT;
DECLARE @riskModuleID INT;
DECLARE @taskModuleID INT;
DECLARE @settingModuleID INT;

-- Insert Dashboard Module
INSERT INTO modules
    (identifier, title_en, title_ar, [order])
VALUES
    ('dashboard', 'Dashboard', 'لوحة التحكم', 1);
SELECT @dashboardModuleID = SCOPE_IDENTITY();

-- Insert Organization Module
INSERT INTO modules
    (identifier, title_en, title_ar, [order])
VALUES
    ('organization', 'Organization', 'المنظمة', 2);
SELECT @organizationModuleID = SCOPE_IDENTITY();

-- Insert Organization Sub-Modules
INSERT INTO modules
    (identifier, title_en, title_ar, module_id, [order])
VALUES
    ('strategic_profile', 'Strategic Profile', 'الملف الاستراتيجي', @organizationModuleID, 2),
    ('user', 'User', 'المستخدم', @organizationModuleID, 3),
    ('process', 'Process', 'العملية', @organizationModuleID, 4),
    ('control', 'Control', 'الضابط', @organizationModuleID, 5);

-- Insert Compliance Module
INSERT INTO modules
    (identifier, title_en, title_ar, [order])
VALUES
    ('compliance', 'Compliance', 'الامتثال', 3);
SELECT @complianceModuleID = SCOPE_IDENTITY();

-- Insert Compliance Sub-Modules
INSERT INTO modules
    (identifier, title_en, title_ar, module_id, [order])
VALUES
    ('compliance_register', 'Compliance Register', 'سجل الامتثال', @complianceModuleID, 2),
    ('sla_contract', 'SLA & Contract', 'SLA والعقد', @complianceModuleID, 3),
    ('compliance_action_plan', 'Compliance Action Plan', 'خطة عمل الامتثال', @incidentModuleID, 4),
    ('sla_action_plan', 'SLA Action Plan', 'خطة عمل جيش تحرير السودان', @incidentModuleID, 5);


-- Insert Document Module
INSERT INTO modules
    (identifier, title_en, title_ar, [order])
VALUES
    ('document', 'Document', 'المستند', 4);
SELECT @documentModuleID = SCOPE_IDENTITY();

-- Insert Document Sub-Modules
INSERT INTO modules
    (identifier, title_en, title_ar, module_id, [order])
VALUES
    ('change_request', 'Change Request', 'طلب التغيير', @documentModuleID, 2);

-- Insert Incident Module
INSERT INTO modules
    (identifier, title_en, title_ar, [order])
VALUES
    ('incident', 'Incident', 'الحادث', 5);
SELECT @incidentModuleID = SCOPE_IDENTITY();

-- Insert Incident Sub-Modules
INSERT INTO modules
    (identifier, title_en, title_ar, module_id, [order])
VALUES
    ('incident_register', 'Incident Register', 'سجل الحوادث', @incidentModuleID, 2),
    ('incident_action_plan', 'Action Plan', 'خطة العمل', @incidentModuleID, 3);

-- Insert Management Review Module
INSERT INTO modules
    (identifier, title_en, title_ar, [order])
VALUES
    ('management_review', 'Management Review', 'استعراض الإدارة', 6);
SELECT @managementReviewModuleID = SCOPE_IDENTITY();

-- Insert Management Review Sub-Modules
INSERT INTO modules
    (identifier, title_en, title_ar, module_id, [order])
VALUES
    ('meeting', 'Meeting', 'الاجتماع', @managementReviewModuleID, 2),
    ('meeting_action_plan', 'Action Plan', 'خطة العمل', @managementReviewModuleID, 3);

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
    ('risk_mitigation', 'Risk Mitigation', 'تخفيف المخاطر', @riskModuleID, 6),
    ('kri', 'KRI', 'مؤشر الأداء الرئيسي', @riskModuleID, 7),
    ('risk_heatmap', 'Risk Heatmap', 'خريطة الخطر', @riskModuleID, 8);

-- Insert Task Module
INSERT INTO modules
    (identifier, title_en, title_ar, [order])
VALUES
    ('task', 'Task', 'المهمة', 8);
SELECT @taskModuleID = SCOPE_IDENTITY();

-- Insert Setting Module
INSERT INTO modules
    (identifier, title_en, title_ar, [order])
VALUES
    ('setting', 'Setting', 'الإعداد', 9);
SELECT @settingModuleID = SCOPE_IDENTITY();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE modules;
-- +goose StatementEnd
