-- +goose Up
-- +goose StatementBegin
DECLARE @policyManagementModuleID INT;
DECLARE @documentModuleID INT;
DECLARE @auditModuleID INT;

SET @policyManagementModuleID = (SELECT id FROM modules WHERE identifier = 'policy_management');
SET @documentModuleID = (SELECT id FROM modules WHERE identifier = 'document');
SET @auditModuleID = (SELECT id FROM modules WHERE identifier = 'audit');
 
INSERT INTO workflow_contexts (module_id, title_en, title_ar, type) VALUES
(@documentModuleID, 'Document Workflow', 'سير عمل الوثيقة', 'document_register'),
(@auditModuleID, 'Audit Final Report Workflow', 'سير عمل التقرير النهائي للتدقيق', 'audit_final_report'),
(@policyManagementModuleID, 'Policy Register Workflow', 'سير عمل تسجيل السياسة', 'policy_register');
 
-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
 
TRUNCATE TABLE workflow_contexts;
 
-- +goose StatementEnd