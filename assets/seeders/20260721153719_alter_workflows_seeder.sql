-- +goose Up
-- +goose StatementBegin
DECLARE @contextID INT;

SET @contextID = (
    SELECT id
    FROM workflow_contexts
    WHERE type = 'vulnerability-closure'
);

UPDATE workflows
SET
    title_en = 'Vulnerability Closure Workflow',
    title_ar = N'سير عمل إغلاق الثغرات'
WHERE workflow_context_id = @contextID;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DECLARE @contextID INT;

SET @contextID = (
    SELECT id
    FROM workflow_contexts
    WHERE type = 'vulnerability-closure'
);

UPDATE workflows
SET
    title_en = 'Vulnerability Reject Workflow',
    title_ar = N'سير عمل رفض الثغرات'
WHERE workflow_context_id = @contextID;
-- +goose StatementEnd
