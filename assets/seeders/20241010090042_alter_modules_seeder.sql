-- +goose Up
-- +goose StatementBegin

DECLARE @bcmModuleID INT;
DECLARE @settingModuleID INT;

UPDATE modules
SET [order] = 10
WHERE identifier = 'management_review';
-- Insert BCM Module
INSERT INTO modules (identifier, title_en, title_ar, [order], icon, route)
VALUES ('business_continuity', 'Business Continuity', 'استمرارية الأعمال', 9, 'icon-incident.svg', '/business-continuity');
SELECT @bcmModuleID = SCOPE_IDENTITY();

-- Insert BCM Sub-Modules
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('policies', 'Policies', 'السياسات', @bcmModuleID, 1, NULL, '/business-continuity/policies'),
    ('business_functions', 'Business Functions', 'وظائف الأعمال', @bcmModuleID, 2, NULL, '/business-continuity/functions'),
    ('bia', 'BIA', 'بيا', @bcmModuleID, 3, NULL, '/business-continuity/business-impact-analysis'),
    ('bia_matrics', 'BIA Matrics', 'مصفوفات BIA', @bcmModuleID, 4, NULL, '/business-continuity/bia-matrics'),
    ('bcp', 'BCP', 'BCP', @bcmModuleID, 5, NULL, '/business-continuity/bcp');

-- Insert Setting Module
INSERT INTO modules (identifier, title_en, title_ar, [order], icon, route)
VALUES ('setting', 'Setting', 'جلسة', 11, 'icon-settings.svg', '/settings');
SELECT @settingModuleID = SCOPE_IDENTITY();
-- Insert Setting Sub-Modules
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('organization_setting', 'Organization Setting', 'إعداد المنظمة', @settingModuleID, 1, NULL, '/setting/organization-setting'),
    ('compliance_setting', 'Compliance Setting', 'إعداد الامتثال', @settingModuleID, 2, NULL, '/setting/compliance-setting'),
    ('risk_setting', 'Risk Setting', 'تحديد المخاطر', @settingModuleID, 3, NULL, '/setting/risk-setting'),
    ('audit_setting', 'Audit Setting', 'إعداد التدقيق', @settingModuleID, 4, NULL, '/setting/audit-setting'),
    ('document_setting', 'Document Setting', 'إعداد الوثيقة', @settingModuleID, 5, NULL, '/setting/document-setting'),
    ('incident_setting', 'Incident Setting', 'إعداد الحادث', @settingModuleID, 6, NULL, '/setting/incident-setting'),
    ('review_meeting_setting', 'Review Meeting Setting', 'مراجعة إعدادات الاجتماع', @settingModuleID, 7, NULL, '/setting/review-meeting-setting'),
    ('audit_insight_setting', 'Audit Insight Setting', 'إعداد رؤية التدقيق', @settingModuleID, 8, NULL, '/setting/audit-insight-setting'),
    ('business_continuity_setting', 'Business Continuity Setting', 'إعداد استمرارية الأعمال', @settingModuleID, 9, NULL, '/setting/business-continuity-setting');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE modules 
SET [order] = 10 
WHERE identifier = 'management_review';

TRUNCATE TABLE modules;
-- +goose StatementEnd
