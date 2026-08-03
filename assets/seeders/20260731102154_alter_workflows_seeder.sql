-- +goose Up
-- +goose StatementBegin
DECLARE @contextID INT;

SET @contextID = (SELECT id FROM workflow_contexts WHERE [type] = 'cyber_bia');

INSERT INTO workflows (status, created_by_id, created_at, title_en, title_ar, workflow_context_id)
VALUES 
    ('active', 1, GETDATE(), 'Cyber BIA Workflow', 'سير عمل تقييم الأثر على الأعمال السيبراني', @contextID);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DECLARE @contextID INT;

SET @contextID = (SELECT id FROM workflow_contexts WHERE [type] = 'cyber_bia');

DELETE FROM workflows WHERE workflow_context_id = @contextID;
-- +goose StatementEnd