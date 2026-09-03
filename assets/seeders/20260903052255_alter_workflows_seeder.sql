-- +goose Up
-- +goose StatementBegin
DECLARE @contextID INT;

SET @contextID = (SELECT id FROM workflow_contexts WHERE [type] = 'project_change_request');

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
    'project-change-request-workflow',
    'active',
    1,
    GETDATE(),
    'Project Change Request Workflow',
    N'سير عمل طلب تغيير المشروع',
    @contextID
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DECLARE @contextID INT;

SET @contextID = (SELECT id FROM workflow_contexts WHERE [type] = 'project_change_request');

DELETE FROM workflows WHERE workflow_context_id = @contextID;
-- +goose StatementEnd
