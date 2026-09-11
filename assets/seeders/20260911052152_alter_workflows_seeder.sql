-- +goose Up
-- +goose StatementBegin
DECLARE @contextID INT;

SET @contextID = (SELECT id FROM workflow_contexts WHERE [type] = 'project');

INSERT INTO workflows (
    type,
    status,
    created_by_id,
    created_at,
    title_en,
    title_ar,
    workflow_context_id
)
VALUES (
    'project-workflow',
    'active',
    1,
    GETDATE(),
    'Project Workflow',
    N'سير عمل المشروع',
    @contextID
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DECLARE @contextID INT;

SET @contextID = (SELECT id FROM workflow_contexts WHERE [type] = 'project');

DELETE FROM workflows WHERE workflow_context_id = @contextID;
-- +goose StatementEnd
