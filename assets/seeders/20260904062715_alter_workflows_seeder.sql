-- +goose Up
-- +goose StatementBegin
DECLARE @contextID INT;

SET @contextID = (SELECT id FROM workflow_contexts WHERE [type] = 'project_decision_log');

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
    'project-decision-workflow',
    'active',
    1,
    GETDATE(),
    'Project Decision Workflow',
    N'سير عمل قرارات المشروع',
    @contextID
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DECLARE @contextID INT;

SET @contextID = (SELECT id FROM workflow_contexts WHERE [type] = 'project_decision_log');

DELETE FROM workflows WHERE workflow_context_id = @contextID;
-- +goose StatementEnd
