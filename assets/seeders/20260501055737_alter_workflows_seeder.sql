-- +goose Up
-- +goose StatementBegin
DECLARE @contextID INT;

SET @contextID = (SELECT id FROM workflow_contexts WHERE [type] = 'risk_response');

INSERT INTO workflows (status, created_by_id, created_at, title_en, title_ar, workflow_context_id)
VALUES 
    ('active', 1, GETDATE(), 'Risk Response Workflow', 'سير عمل استجابة المخاطر', @contextID);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DECLARE @contextID INT;

SET @contextID = (SELECT id FROM workflow_contexts WHERE [type] = 'risk_response');

DELETE FROM workflows WHERE workflow_context_id = @contextID;
-- +goose StatementEnd