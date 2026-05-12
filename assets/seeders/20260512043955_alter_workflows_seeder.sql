-- +goose Up
-- +goose StatementBegin
DECLARE @contextID INT;

SET @contextID = (SELECT id FROM workflow_contexts WHERE [type] = 'risk_register_update_plan');

INSERT INTO workflows (status, created_by_id, created_at, title_en, title_ar, workflow_context_id)
VALUES 
    ('active', 1, GETDATE(), 'Risk Register Update Plan Workflow', 'سير عمل تحديث خطة المخاطر', @contextID);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DECLARE @contextID INT;

SET @contextID = (SELECT id FROM workflow_contexts WHERE [type] = 'risk_register_update_plan');

DELETE FROM workflows WHERE workflow_context_id = @contextID;
-- +goose StatementEnd