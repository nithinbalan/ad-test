-- +goose Up
-- +goose StatementBegin
-- Get the ID of the module with identifier 'audit'
DECLARE @auditInsightsModuleID INT;
SET @auditInsightsModuleID = (SELECT id FROM modules WHERE identifier = 'audit_insights');

-- Insert new module records for Audit Insights
INSERT INTO modules (module_id, identifier, title_en, title_ar)
VALUES
    (@auditInsightsModuleID,'audit_insight_plans', 'Audit Insight Audit Plans', 'خطط رؤى التدقيق'),
    (@auditInsightsModuleID,'audit_insight_audits', 'Audit Insight Audits', 'تدقيقات رؤى التدقيق'),
    (@auditInsightsModuleID,'document_requests', 'Audit Document Requests', 'طلبات مستندات التدقيق'),
    (@auditInsightsModuleID,'draft_reports', 'Audit Draft Reports', 'تقارير التدقيق الأولية'),
    (@auditInsightsModuleID,'findings', 'Audit Findings', 'نتائج التدقيق'),
    (@auditInsightsModuleID,'action_plans', 'Audit Action Plans', 'خطط عمل التدقيق');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove newly added modules
DELETE FROM modules
WHERE identifier IN (
    'audit_insight_plans',
    'audit_insight_audits',
    'document_requests',
    'draft_reports',
    'findings',
    'action_plans'
);

-- Truncate the notification_masters table
TRUNCATE TABLE notification_masters;
-- +goose StatementEnd
