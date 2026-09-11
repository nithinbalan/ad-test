-- +goose Up
-- +goose StatementBegin
DECLARE @contextID INT;

SET @contextID = (SELECT id FROM workflow_contexts WHERE [type] = 'project_register');

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
    'project-register-workflow',
    'active',
    1,
    GETDATE(),
    'Project Register Workflow',
    N'سير عمل سجل المشاريع',
    @contextID
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DECLARE @contextID INT;

SET @contextID = (SELECT id FROM workflow_contexts WHERE [type] = 'project_register');

DELETE FROM workflows WHERE workflow_context_id = @contextID;
-- +goose StatementEnd
