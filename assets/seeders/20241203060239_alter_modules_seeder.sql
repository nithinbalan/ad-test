-- +goose Up
-- +goose StatementBegin
DECLARE @organizationModuleID INT;
DECLARE @assetModuleID INT;
DECLARE @complianceModuleID INT;
DECLARE @riskModuleID INT;
DECLARE @auditModuleID INT;
DECLARE @auditInsightsModuleID INT;
DECLARE @documentModuleID INT;
DECLARE @bcmModuleID INT;
DECLARE @incidentModuleID INT;
DECLARE @managementReviewModuleID INT;

DECLARE @assetRiskModuleID INT;

SELECT @organizationModuleID = id FROM modules WHERE identifier = 'organization';
SELECT @assetModuleID = id FROM modules WHERE identifier = 'asset';
SELECT @complianceModuleID = id FROM modules WHERE identifier = 'compliance';
SELECT @riskModuleID = id FROM modules WHERE identifier = 'risk';
SELECT @auditModuleID = id FROM modules WHERE identifier = 'audit';
SELECT @auditInsightsModuleID = id FROM modules WHERE identifier = 'audit_insights';
SELECT @documentModuleID = id FROM modules WHERE identifier = 'document';
SELECT @bcmModuleID = id FROM modules WHERE identifier = 'business_continuity';
SELECT @incidentModuleID = id FROM modules WHERE identifier = 'incident';
SELECT @managementReviewModuleID = id FROM modules WHERE identifier = 'management_review';

SELECT @assetRiskModuleID = id FROM modules WHERE identifier = 'risk' and module_id=@assetModuleID;

-- Update identifier for BCM Policies
UPDATE modules
SET identifier = 'asset_risk'
WHERE id = @assetRiskModuleID;

-- adding reports to all modules

INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('organization_report', 'Report', 'تقرير', @organizationModuleID, 5, NULL, '/organization/report');

INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('asset_report', 'Report', 'تقرير', @assetModuleID, 6, NULL, '/asset/report');

INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('compliance_report', 'Report', 'تقرير', @complianceModuleID, 4, NULL, '/compliance/report');

INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('risk_report', 'Report', 'تقرير', @riskModuleID, 8, NULL, '/risk/report');

INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('audit_report', 'Report', 'تقرير', @auditModuleID, 5, NULL, '/audit/report');

INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('audit_insight_report', 'Report', 'تقرير', @auditInsightsModuleID, 7, NULL, '/audit-insight/report');

INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('document_report', 'Report', 'تقرير', @documentModuleID, 3, NULL, '/document/report');

INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('bcm_report', 'Report', 'تقرير', @bcmModuleID, 7, NULL, '/business-continuity/report');

INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('incident_report', 'Report', 'تقرير', @incidentModuleID, 3, NULL, '/incident/report');

INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('meeting_report', 'Report', 'تقرير', @managementReviewModuleID, 3, NULL, '/review-meeting/report');
-- Update arabic language for all modules
UPDATE modules
SET title_ar = 'لوحات المعلومات'
WHERE identifier = 'dashboard';

UPDATE modules
SET title_ar = 'لوحة معلومات الحوكمة وإدارة المخاطر والامتثال'
WHERE identifier = 'grc_dashboard';

UPDATE modules
SET title_ar = 'المنظمة'
WHERE identifier = 'organization';

UPDATE modules
SET title_ar = 'الملف الاستراتيجي'
WHERE identifier = 'strategic_profile';

UPDATE modules
SET title_en = 'Users', title_ar = 'المستخدمون'
WHERE identifier = 'user_management';

UPDATE modules
SET title_ar = 'السياسات'
WHERE identifier = 'policies';

UPDATE modules
SET title_ar = 'العمليات'
WHERE identifier = 'process';

UPDATE modules
SET title_ar = 'أصل'
WHERE identifier = 'asset';

UPDATE modules
SET title_ar = 'سجل الأصول'
WHERE identifier = 'asset_register';

UPDATE modules
SET title_ar = 'مصفوفة أهمية الأصول'
WHERE identifier = 'asset_criticality_matrix';

UPDATE modules
SET title_ar = 'نقاط الضعف'
WHERE identifier = 'vulnerability';

UPDATE modules
SET title_ar = 'التهديدات'
WHERE identifier = 'threat';

UPDATE modules
SET title_ar = 'المخاطر'
WHERE identifier = 'asset_risk';

UPDATE modules
SET title_ar = 'الامتثال'
WHERE identifier = 'compliance';

UPDATE modules
SET title_ar = 'سجل الامتثال'
WHERE identifier = 'compliance_register';

UPDATE modules
SET title_ar = 'اتفاقيات مستوى الخدمة والعقود'
WHERE identifier = 'sla_contracts';

UPDATE modules
SET title_ar = 'الضوابط'
WHERE identifier = 'controls';

UPDATE modules
SET title_ar = 'المخاطر'
WHERE identifier = 'risk';

UPDATE modules
SET title_ar = 'مؤشرات تحليل المخاطر'
WHERE identifier = 'risk_matrix';

UPDATE modules
SET title_ar = 'بيان تحمل المخاطر'
WHERE identifier = 'risk_appetite_statement';

UPDATE modules
SET title_ar = 'المخاطر المحتملة'
WHERE identifier = 'potential_risk';

UPDATE modules
SET title_ar = 'مؤشرات المخاطر الرئيسية'
WHERE identifier = 'kri';

UPDATE modules
SET title_ar = 'سجل المخاطر'
WHERE identifier = 'risk_register';

UPDATE modules
SET title_ar = 'إدارة المخاطر'
WHERE identifier = 'risk_mitigation';

UPDATE modules
SET title_ar = 'خريطة توزيع المخاطر'
WHERE identifier = 'risk_heatmap';

UPDATE modules
SET title_ar = 'التدقيق'
WHERE identifier = 'audit';

UPDATE modules
SET title_ar = 'نطاق التدقيق'
WHERE identifier = 'audit_universe';

UPDATE modules
SET title_ar = 'خطط التدقيق'
WHERE identifier = 'audit_plans';

UPDATE modules
SET title_ar = 'عمليات التدقيق'
WHERE identifier = 'audits';

UPDATE modules
SET title_ar = 'نتائج التدقيق'
WHERE identifier = 'audit_findings';

UPDATE modules
SET title_ar = 'تحليلات التدقيق'
WHERE identifier = 'audit_insights';

UPDATE modules
SET title_ar = 'خطط التدقيق'
WHERE identifier = 'audit_insight_plans';

UPDATE modules
SET title_ar = 'عمليات التدقيق'
WHERE identifier = 'audit_insight_audits';

UPDATE modules
SET title_ar = 'التقارير المبدئية'
WHERE identifier = 'document_requests';

UPDATE modules
SET title_ar = 'التقارير المبدئية'
WHERE identifier = 'draft_reports';

UPDATE modules
SET title_ar = 'الاستنتاجات'
WHERE identifier = 'findings';

UPDATE modules
SET title_ar = 'خطط الإجراءات'
WHERE identifier = 'action_plans';

UPDATE modules
SET title_ar = 'الوثائق'
WHERE identifier = 'document';

UPDATE modules
SET title_ar = 'سجل الوثائق'
WHERE identifier = 'document_register';

UPDATE modules
SET title_ar = 'طلب التغيير'
WHERE identifier = 'change_request';

UPDATE modules
SET title_ar = 'استمرارية الأعمال'
WHERE identifier = 'business_continuity';

UPDATE modules
SET title_ar = 'السياسات'
WHERE identifier = 'bcm_policies';

UPDATE modules
SET title_ar = 'الوظائف التشغيلية'
WHERE identifier = 'business_functions';

UPDATE modules
SET title_ar = 'تحليل تأثير الأعمال'
WHERE identifier = 'bia';

UPDATE modules
SET title_ar = 'مؤشرات تحليل تأثير الأعمال'
WHERE identifier = 'bia_matrics';

UPDATE modules
SET title_ar = 'خطة استمرارية الأعمال'
WHERE identifier = 'bcp';

UPDATE modules
SET title_ar = 'الأحداث'
WHERE identifier = 'incident';

UPDATE modules
SET title_ar = 'سجل الأحداث'
WHERE identifier = 'incident_register';

UPDATE modules
SET title_ar = 'خطة الإجراءات'
WHERE identifier = 'incident_action_plan';

UPDATE modules
SET title_ar = 'مراجعة الإدارة'
WHERE identifier = 'management_review';

UPDATE modules
SET title_ar = 'اجتماع'
WHERE identifier = 'meeting';

UPDATE modules
SET title_ar = 'خطة الإجراءات'
WHERE identifier = 'meeting_action_plan';

UPDATE modules
SET title_ar = 'الإعدادات'
WHERE identifier = 'setting';

UPDATE modules
SET title_ar = 'عام'
WHERE identifier = 'general_setting';

UPDATE modules
SET title_ar = 'المنظمة'
WHERE identifier = 'organization_setting';

UPDATE modules
SET title_ar = 'الامتثال'
WHERE identifier = 'compliance_setting';

UPDATE modules
SET title_ar = 'المخاطر'
WHERE identifier = 'risk_setting';

UPDATE modules
SET title_ar = 'التدقيق'
WHERE identifier = 'audit_setting';

UPDATE modules
SET title_ar = 'الوثائق'
WHERE identifier = 'document_setting';

UPDATE modules
SET title_ar = 'الأحداث'
WHERE identifier = 'incident_setting';

UPDATE modules
SET title_ar = 'مراجعة الإدارة'
WHERE identifier = 'review_meeting_setting';

UPDATE modules
SET title_ar = 'تحليلات التدقيق'
WHERE identifier = 'audit_insight_setting';

UPDATE modules
SET title_ar = 'استمرارية الأعمال'
WHERE identifier = 'business_continuity_setting';

UPDATE modules
SET title_ar = 'أصل'
WHERE identifier = 'asset_setting';

UPDATE modules
SET title_ar = 'سجلات'
WHERE identifier = 'log';

UPDATE modules
SET title_ar = 'المدراء'
WHERE identifier = 'common_setting';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE modules;
-- +goose StatementEnd
