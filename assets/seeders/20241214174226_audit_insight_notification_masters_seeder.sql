-- +goose Up
-- +goose StatementBegin
-- Get the ID of the module with identifier 'audit'
DECLARE @module_id INT;
SET @module_id = (SELECT id FROM modules WHERE identifier = 'audit_insights');

-- Risk module notifications
INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES
    (@module_id, 'audit_insight_plans', 'Notification on Audit Insight Plans', 'إشعار بشأن خطط رؤى التدقيق'),
    (@module_id, 'audit_insight_audits', 'Notification on Audit Insight Audits', 'إشعار بشأن تدقيقات رؤى التدقيق'),
    (@module_id, 'audit_insight_document_requests', 'Notification on Audit Insight Document Requests', 'إشعار بشأن طلبات مستندات رؤى التدقيق'),
    (@module_id, 'audit_insight_draft_reports', 'Notification on Audit Insight Draft Reports', 'إشعار بشأن تقارير رؤى التدقيق الأولية'),
    (@module_id, 'audit_insight_findings', 'Notification on Audit Insight Findings', 'إشعار بشأن نتائج رؤى التدقيق'),
    (@module_id, 'audit_insight_action_plans', 'Notification on Audit Insight Action Plans', 'إشعار بشأن خطط عمل رؤى التدقيق');
    
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE notification_masters;
-- +goose StatementEnd
