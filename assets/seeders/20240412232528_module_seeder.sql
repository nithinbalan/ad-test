-- +goose Up
-- +goose StatementBegin

DECLARE @dashboardModuleID INT;
DECLARE @organizationModuleID INT;
DECLARE @complianceModuleID INT;
DECLARE @documentModuleID INT;
DECLARE @incidentModuleID INT;
DECLARE @managementReviewModuleID INT;
DECLARE @riskModuleID INT;
DECLARE @auditModuleID INT;
DECLARE @auditInsightsModuleID INT;

-- Insert Dashboard Module
INSERT INTO modules (identifier, title_en, title_ar, [order], icon, route)
VALUES ('dashboard', 'Dashboards', 'لوحة التحكم', 1, 'icon-dashboard.svg', '/dashboard');
SELECT @dashboardModuleID = SCOPE_IDENTITY();

-- Insert Dashboard Sub-Modules (children)
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('grc_dashboard', 'GRC Dashboard', 'لوحة تحكم GRC', @dashboardModuleID, 1, NULL, '/dashboard/my-dashboard');

-- Insert Organization Module
INSERT INTO modules (identifier, title_en, title_ar, [order], icon, route)
VALUES ('organization', 'Organization', 'المنظمة', 2, 'icon-organization.svg', '/organization');
SELECT @organizationModuleID = SCOPE_IDENTITY();

-- Insert Organization Sub-Modules
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('strategic_profile', 'Strategic Profile', 'الملف الاستراتيجي', @organizationModuleID, 1, NULL, '/organization/strategic-profile'),
    ('user_management', 'User Management', 'إدارة المستخدم', @organizationModuleID, 2, NULL, '/organization/users'),
    ('policies', 'Policies', 'السياسات', @organizationModuleID, 3, NULL, '/organization/policies'),
    ('process', 'Process', 'العملية', @organizationModuleID, 4, NULL, '/organization/process');

-- Insert Compliance Module
INSERT INTO modules (identifier, title_en, title_ar, [order], icon, route)
VALUES ('compliance', 'Compliance', 'الامتثال', 3, 'icon-compliance.svg', '/compliance');
SELECT @complianceModuleID = SCOPE_IDENTITY();

-- Insert Compliance Sub-Modules
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('compliance_register', 'Compliance Register', 'سجل الامتثال', @complianceModuleID, 1, NULL, '/compliance/register'),
    ('sla_contracts', 'SLA & Contracts', 'SLA والعقود', @complianceModuleID, 2, NULL, '/compliance/sla-contracts'),
    ('controls', 'Controls', 'الضوابط', @complianceModuleID, 3, NULL, '/compliance/controls');

-- Insert Risk Module
INSERT INTO modules (identifier, title_en, title_ar, [order], icon, route)
VALUES ('risk', 'Risk', 'المخاطر', 4, 'icon-risk.svg', '/risk');
SELECT @riskModuleID = SCOPE_IDENTITY();

-- Insert Risk Sub-Modules
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('risk_matrix', 'Risk Matrix', 'مصفوفة المخاطر', @riskModuleID, 1, NULL, '/risk/risk-matrix'),
    ('risk_appetite_statement', 'Risk Appetite Statement', 'ميول المخاطر', @riskModuleID, 2, NULL, '/risk/risk-appetite-statement'),
    ('potential_risk', 'Potential Risk', 'المخاطر المحتملة', @riskModuleID, 3, NULL, '/risk/potential-risk'),
    ('kri', 'KRI', 'مؤشر الأداء الرئيسي', @riskModuleID, 4, NULL, '/risk/kri-register'),
    ('risk_register', 'Risk Register', 'سجل المخاطر', @riskModuleID, 5, NULL, '/risk/risk-register'),
    ('risk_mitigation', 'Risk Mitigation', 'تخفيف المخاطر', @riskModuleID, 6, NULL, '/risk/risk-mitigation'),
    ('risk_heatmap', 'Risk Heatmap', 'خريطة الخطر', @riskModuleID, 7, NULL, '/risk/risk-heatmaps');

-- Insert Audit Module
INSERT INTO modules (identifier, title_en, title_ar, [order], icon, route)
VALUES ('audit', 'Audit', 'التدقيق', 5, 'icon-risk.svg', '/audit');
SELECT @auditModuleID = SCOPE_IDENTITY();

-- Insert Audit Sub-Modules
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('audit_universe', 'Audit Universe', 'عالم التدقيق', @auditModuleID, 1, NULL, '/audit/audit-universe'),
    ('audit_plans', 'Audit Plans', 'خطط التدقيق', @auditModuleID, 2, NULL, '/audit/audit-plans'),
    ('audits', 'Audits', 'التدقيقات', @auditModuleID, 3, NULL, '/audit/audits'),
    ('audit_findings', 'Audit Findings', 'نتائج التدقيق', @auditModuleID, 4, NULL, '/audit/audit-findings');

-- Insert Audit Insights Module
INSERT INTO modules (identifier, title_en, title_ar, [order], icon, route)
VALUES ('audit_insights', 'Audit Insights', 'إحصاءات التدقيق', 6, 'icon-risk.svg', '/audit-insights');
SELECT @auditInsightsModuleID = SCOPE_IDENTITY();

-- Insert Audit Insights Sub-Modules
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES
    ('audit_insight_plans', 'Audit Plans', 'خطط التدقيق', @auditInsightsModuleID, 1, NULL, '/audit-insight/audit-plans'),
    ('audit_insight_audits', 'Audits', 'التدقيقات', @auditInsightsModuleID, 2, NULL, '/audit-insight/audits'),
    ('document_requests', 'Document Requests', 'طلبات المستندات', @auditInsightsModuleID, 3, NULL, '/audit-insight/document-requests'),
    ('draft_reports', 'Draft Reports', 'التقارير الأولية', @auditInsightsModuleID, 4, NULL, '/audit-insight/draft-reports'),
    ('findings', 'Findings', 'النتائج', @auditInsightsModuleID, 5, NULL, '/audit-insight/findings'),
    ('action_plans', 'Action Plans', 'خطط العمل', @auditInsightsModuleID, 6, NULL, '/audit-insight/action-plans');

-- Insert Document Module
INSERT INTO modules (identifier, title_en, title_ar, [order], icon, route)
VALUES ('document', 'Document', 'المستند', 7, 'icon-document.svg', '/document');
SELECT @documentModuleID = SCOPE_IDENTITY();

-- Insert Document Sub-Modules
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('document_register', 'Document Register', 'سجل المستندات', @documentModuleID, 1, NULL, '/document/register'),
    ('change_request', 'Change Request', 'طلب التغيير', @documentModuleID, 2, NULL, '/document/change-request');

-- Insert Incident Module
INSERT INTO modules (identifier, title_en, title_ar, [order], icon, route)
VALUES ('incident', 'Incident', 'الحادث', 8, 'icon-incident.svg', '/incident');
SELECT @incidentModuleID = SCOPE_IDENTITY();

-- Insert Incident Sub-Modules
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('incident_register', 'Incident Register', 'سجل الحوادث', @incidentModuleID, 1, NULL, '/incident/register'),
    ('incident_action_plan', 'Action Plan', 'خطة العمل', @incidentModuleID, 2, NULL, '/incident/action-plan');

-- Insert Management Review Module
INSERT INTO modules (identifier, title_en, title_ar, [order], icon, route)
VALUES ('management_review', 'Management Review', 'استعراض الإدارة', 9, 'icon-management-review.svg', '/management-review');
SELECT @managementReviewModuleID = SCOPE_IDENTITY();

-- Insert Management Review Sub-Modules
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('meeting', 'Meeting', 'الاجتماع', @managementReviewModuleID, 1, NULL, '/review-meeting/meetings'),
    ('meeting_action_plan', 'Action Plan', 'خطة العمل', @managementReviewModuleID, 2, NULL, '/review-meeting/action-plan-list');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE modules;
-- +goose StatementEnd
