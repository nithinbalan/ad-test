-- +goose Up
-- +goose StatementBegin
DECLARE @documentWorkflowContextID INT;
DECLARE @auditWorkflowContextID INT;
DECLARE @policyWorkflowContextID INT;
DECLARE @policyDocumentTypeID INT;
DECLARE @procedureDocumentTypeID INT;


SET @documentWorkflowContextID = (SELECT id FROM workflow_contexts WHERE type = 'document_register');
SET @auditWorkflowContextID = (SELECT id FROM workflow_contexts WHERE type = 'audit_final_report');
SET @policyWorkflowContextID = (SELECT id FROM workflow_contexts WHERE type = 'policy_register');
SET @policyDocumentTypeID = (SELECT id FROM document_types WHERE type = 'policy');
SET @procedureDocumentTypeID = (SELECT id FROM document_types WHERE type = 'procedure');

UPDATE workflows
SET workflow_context_id = @documentWorkflowContextID
WHERE type = 'document' AND document_type_id IS NOT NULL;


UPDATE workflows
SET workflow_context_id = @auditWorkflowContextID
WHERE type = 'final-report';

INSERT INTO workflows (document_type_id,status, created_by_id, created_at, title_en, title_ar,workflow_context_id)
VALUES 
    (@policyDocumentTypeID,'active', 1, GETDATE(), 'Policy Workflow', 'سير عمل السياسة', @policyWorkflowContextID),
    (@procedureDocumentTypeID,'active', 1, GETDATE(), 'Procedure Workflow', 'سير العمل الإجراء', @policyWorkflowContextID);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DECLARE @documentWorkflowContextID INT;
DECLARE @auditWorkflowContextID INT;
DECLARE @policyWorkflowContextID INT;

SET @documentWorkflowContextID = (SELECT id FROM workflow_contexts WHERE type = 'document_register');
SET @auditWorkflowContextID = (SELECT id FROM workflow_contexts WHERE type = 'audit_final_report');
SET @policyWorkflowContextID = (SELECT id FROM workflow_contexts WHERE type = 'policy_register');

DELETE FROM workflows WHERE workflow_context_id = @policyWorkflowContextID;

UPDATE workflows
SET workflow_context_id = NULL
WHERE type = 'document' AND document_type_id IS NOT NULL;

UPDATE workflows
SET workflow_context_id = NULL
WHERE type = 'final-report';
-- +goose StatementEnd
