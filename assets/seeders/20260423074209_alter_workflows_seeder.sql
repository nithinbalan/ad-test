-- +goose Up
-- +goose StatementBegin
DECLARE @contextID INT;

SET @contextID = (SELECT id FROM workflow_contexts WHERE [type] = 'kri_register');

INSERT INTO workflows (status, created_by_id, created_at, title_en, title_ar, workflow_context_id)
VALUES 
    ('active', 1, GETDATE(), 'KRI Register Workflow', 'سير عمل سجل مؤشرات المخاطر الرئيسية', @contextID);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DECLARE @contextID INT;

SET @contextID = (SELECT id FROM workflow_contexts WHERE [type] = 'kri_register');

DELETE FROM workflows WHERE workflow_context_id = @contextID;
-- +goose StatementEnd