-- +goose Up
-- +goose StatementBegin
DECLARE @contextID INT;

SET @contextID = (SELECT id FROM workflow_contexts WHERE [type] = 'project_security_gate_exception');

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
    'project-security-gate-exception-workflow',
    'active',
    1,
    GETDATE(),
    'Project Security Gate Exception Workflow',
    N'سير عمل استثناء بوابة أمان المشروع',
    @contextID
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DECLARE @contextID INT;

SET @contextID = (SELECT id FROM workflow_contexts WHERE [type] = 'project_security_gate_exception');

DELETE FROM workflows WHERE workflow_context_id = @contextID;
-- +goose StatementEnd
