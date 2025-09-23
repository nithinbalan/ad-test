-- +goose Up
-- +goose StatementBegin
DECLARE @riskWorkflowContextID INT;

SET @riskWorkflowContextID = (SELECT id FROM workflow_contexts WHERE type = 'risk_register');

INSERT INTO workflows (status, created_by_id, created_at, title_en, title_ar,workflow_context_id)
VALUES 
    ('active', 1, GETDATE(), 'Risk Workflow', 'سير عمل المخاطر', @riskWorkflowContextID);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM workflows WHERE workflow_context_id = @riskWorkflowContextID;
-- +goose StatementEnd
